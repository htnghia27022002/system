package queue

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	natsio "github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type Handler interface {
	Name() string
	Handle(ctx context.Context, payload []byte) error
}

type HandlerRegistry interface {
	Get(name string) (Handler, error)
}

type Client struct {
	cfg      Config
	conn     *natsio.Conn
	js       jetstream.JetStream
	consumes []jetstream.ConsumeContext
}

func New(cfg Config) (*Client, error) {
	if !cfg.Connection.Enabled {
		return &Client{cfg: cfg}, nil
	}

	reconnectWait, err := time.ParseDuration(cfg.Connection.ReconnectWait)
	if err != nil || reconnectWait <= 0 {
		reconnectWait = 2 * time.Second
	}

	opts := []natsio.Option{
		natsio.Name(cfg.Connection.Name),
		natsio.MaxReconnects(cfg.Connection.MaxReconnects),
		natsio.ReconnectWait(reconnectWait),
		natsio.DisconnectErrHandler(func(_ *natsio.Conn, err error) {
			if err != nil {
				log.Printf("queue disconnected: %v", err)
			}
		}),
	}

	conn, err := natsio.Connect(cfg.Connection.URL, opts...)
	if err != nil {
		return nil, err
	}

	js, err := jetstream.New(conn)
	if err != nil {
		conn.Close()
		return nil, err
	}

	return &Client{cfg: cfg, conn: conn, js: js}, nil
}

func (c *Client) Enabled() bool {
	return c != nil && c.conn != nil && c.conn.IsConnected()
}

func (c *Client) Close() {
	for _, consume := range c.consumes {
		consume.Stop()
	}
	if c.conn != nil {
		c.conn.Close()
	}
}

func (c *Client) Publish(ctx context.Context, subject string, data []byte) error {
	if !c.Enabled() {
		return nil
	}

	_, err := c.js.Publish(ctx, subject, data)
	return err
}

func (c *Client) EnsureInfrastructure(ctx context.Context) error {
	if !c.Enabled() {
		return nil
	}

	for streamKey, streamOpts := range c.cfg.Streams {
		streamDef, err := resolveStream(streamKey)
		if err != nil {
			return err
		}

		streamCfg, err := buildStreamConfig(streamDef, streamOpts)
		if err != nil {
			return fmt.Errorf("stream %q: %w", streamKey, err)
		}

		if _, err := c.js.CreateOrUpdateStream(ctx, streamCfg); err != nil {
			return fmt.Errorf("ensure stream %s: %w", streamDef.Name, err)
		}
		log.Printf("queue stream ready name=%s subjects=%v", streamDef.Name, streamDef.Subjects)

		for consumerKey, consumerOpts := range streamOpts.Consumers {
			consumerDef, err := resolveConsumer(streamKey, consumerKey)
			if err != nil {
				return err
			}

			consumerCfg, err := buildConsumerConfig(consumerDef, consumerOpts)
			if err != nil {
				return fmt.Errorf("consumer %q: %w", consumerKey, err)
			}

			if _, err := c.js.CreateOrUpdateConsumer(ctx, streamDef.Name, consumerCfg); err != nil {
				return fmt.Errorf("ensure consumer %s: %w", consumerDef.Name, err)
			}
			log.Printf("queue consumer ready stream=%s consumer=%s subject=%s handler=%s",
				streamDef.Name, consumerDef.Name, consumerDef.FilterSubject, consumerDef.Handler)
		}
	}

	return nil
}

func (c *Client) StartConsumers(ctx context.Context, registry HandlerRegistry) error {
	if !c.Enabled() {
		return nil
	}
	if registry == nil {
		return fmt.Errorf("queue: handler registry is required")
	}

	for streamKey, streamOpts := range c.cfg.Streams {
		streamDef, err := resolveStream(streamKey)
		if err != nil {
			return err
		}

		for consumerKey := range streamOpts.Consumers {
			consumerDef, err := resolveConsumer(streamKey, consumerKey)
			if err != nil {
				return err
			}

			handler, err := registry.Get(consumerDef.Handler)
			if err != nil {
				return err
			}

			consumer, err := c.js.Consumer(ctx, streamDef.Name, consumerDef.Name)
			if err != nil {
				return fmt.Errorf("load consumer %s: %w", consumerDef.Name, err)
			}

			handlerName := handler.Name()
			consumeCtx, err := consumer.Consume(func(msg jetstream.Msg) {
				if err := handler.Handle(ctx, msg.Data()); err != nil {
					log.Printf("queue handler %s failed: %v", handlerName, err)
					_ = msg.Nak()
					return
				}
				_ = msg.Ack()
			})
			if err != nil {
				return fmt.Errorf("consume %s: %w", consumerDef.Name, err)
			}

			c.consumes = append(c.consumes, consumeCtx)
			log.Printf("queue consumer started stream=%s consumer=%s handler=%s",
				streamDef.Name, consumerDef.Name, handlerName)
		}
	}

	return nil
}

func buildStreamConfig(def StreamDefinition, opts StreamOptions) (jetstream.StreamConfig, error) {
	cfg := jetstream.StreamConfig{
		Name:     def.Name,
		Subjects: def.Subjects,
	}

	switch strings.ToLower(strings.TrimSpace(opts.Retention)) {
	case "", "limits":
		cfg.Retention = jetstream.LimitsPolicy
	case "interest":
		cfg.Retention = jetstream.InterestPolicy
	case "workqueue", "work_queue":
		cfg.Retention = jetstream.WorkQueuePolicy
	default:
		return jetstream.StreamConfig{}, fmt.Errorf("unsupported retention %q", opts.Retention)
	}

	switch strings.ToLower(strings.TrimSpace(opts.Storage)) {
	case "", "file":
		cfg.Storage = jetstream.FileStorage
	case "memory":
		cfg.Storage = jetstream.MemoryStorage
	default:
		return jetstream.StreamConfig{}, fmt.Errorf("unsupported storage %q", opts.Storage)
	}

	if opts.MaxAge != "" {
		maxAge, err := time.ParseDuration(opts.MaxAge)
		if err != nil {
			return jetstream.StreamConfig{}, fmt.Errorf("invalid max_age: %w", err)
		}
		cfg.MaxAge = maxAge
	}
	if opts.MaxMsgs > 0 {
		cfg.MaxMsgs = opts.MaxMsgs
	}
	if opts.MaxBytes > 0 {
		cfg.MaxBytes = opts.MaxBytes
	}

	switch strings.ToLower(strings.TrimSpace(opts.Discard)) {
	case "", "old":
		cfg.Discard = jetstream.DiscardOld
	case "new":
		cfg.Discard = jetstream.DiscardNew
	default:
		return jetstream.StreamConfig{}, fmt.Errorf("unsupported discard %q", opts.Discard)
	}

	if opts.DuplicateWindow != "" {
		dupWindow, err := time.ParseDuration(opts.DuplicateWindow)
		if err != nil {
			return jetstream.StreamConfig{}, fmt.Errorf("invalid duplicate_window: %w", err)
		}
		cfg.Duplicates = dupWindow
	}

	return cfg, nil
}

func buildConsumerConfig(def ConsumerDefinition, opts ConsumerOptions) (jetstream.ConsumerConfig, error) {
	cfg := jetstream.ConsumerConfig{
		Name:          def.Name,
		Durable:       def.Name,
		FilterSubject: def.FilterSubject,
	}

	switch strings.ToLower(strings.TrimSpace(opts.AckPolicy)) {
	case "", "explicit":
		cfg.AckPolicy = jetstream.AckExplicitPolicy
	case "all":
		cfg.AckPolicy = jetstream.AckAllPolicy
	case "none":
		cfg.AckPolicy = jetstream.AckNonePolicy
	default:
		return jetstream.ConsumerConfig{}, fmt.Errorf("unsupported ack_policy %q", opts.AckPolicy)
	}

	if opts.AckWait != "" {
		ackWait, err := time.ParseDuration(opts.AckWait)
		if err != nil {
			return jetstream.ConsumerConfig{}, fmt.Errorf("invalid ack_wait: %w", err)
		}
		cfg.AckWait = ackWait
	}
	if opts.MaxDeliver > 0 {
		cfg.MaxDeliver = opts.MaxDeliver
	}
	if opts.MaxAckPending > 0 {
		cfg.MaxAckPending = opts.MaxAckPending
	}

	switch strings.ToLower(strings.TrimSpace(opts.DeliverPolicy)) {
	case "", "all":
		cfg.DeliverPolicy = jetstream.DeliverAllPolicy
	case "last":
		cfg.DeliverPolicy = jetstream.DeliverLastPolicy
	case "new":
		cfg.DeliverPolicy = jetstream.DeliverNewPolicy
	case "last_per_subject":
		cfg.DeliverPolicy = jetstream.DeliverLastPerSubjectPolicy
	default:
		return jetstream.ConsumerConfig{}, fmt.Errorf("unsupported deliver_policy %q", opts.DeliverPolicy)
	}

	switch strings.ToLower(strings.TrimSpace(opts.ReplayPolicy)) {
	case "", "instant":
		cfg.ReplayPolicy = jetstream.ReplayInstantPolicy
	case "original":
		cfg.ReplayPolicy = jetstream.ReplayOriginalPolicy
	default:
		return jetstream.ConsumerConfig{}, fmt.Errorf("unsupported replay_policy %q", opts.ReplayPolicy)
	}

	return cfg, nil
}

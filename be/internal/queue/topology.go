package queue

import "context"

// Handler processes one message delivered to a consumer.
type Handler interface {
	Name() string
	Handle(ctx context.Context, payload []byte) error
}

// Stream is a JetStream stream declared by a module, together with its consumers.
// Key matches an optional override entry in nats.json ("streams": {"<key>": {...}}).
type Stream struct {
	Key       string
	Name      string
	Subjects  []string
	Options   StreamOptions
	Consumers []Consumer
}

// Consumer is a durable pull consumer on a Stream. Handler is required only in the
// worker process (cmd/queue); the API process only needs the stream to publish.
type Consumer struct {
	Key           string
	Name          string
	FilterSubject string
	Options       ConsumerOptions
	Handler       Handler
}

// streamOptions applies a nats.json override (whole-object replace, as before
// modules existed) on top of the stream's own defaults.
func (c *Client) streamOptions(s Stream) StreamOptions {
	if override, ok := c.cfg.Streams[s.Key]; ok {
		return override
	}
	return s.Options
}

func (c *Client) consumerOptions(s Stream, consumer Consumer) ConsumerOptions {
	if override, ok := c.cfg.Streams[s.Key]; ok {
		if opts, ok := override.Consumers[consumer.Key]; ok {
			return opts
		}
	}
	return consumer.Options
}

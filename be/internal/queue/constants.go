package queue

// JSON stream keys (internal/queue/nats.json → streams map).
const (
	StreamKeySearch = "search"
)

// JSON consumer keys nested under a stream in nats.json.
const (
	ConsumerKeySearchOutbox = "search_outbox"
)

// JetStream stream names.
const (
	StreamSearch = "SYSTEM_SEARCH"
)

// NATS subjects — use these everywhere instead of string literals.
const (
	SubjectSearchPrefix   = "system.search"
	SubjectSearchOutbox   = "system.search.outbox"
	SubjectSearchWildcard = "system.search.>"
)

// JetStream durable consumer names.
const (
	ConsumerSearchOutbox = "process_search_outbox"
)

// Handler registry keys wired in subscribers.
const (
	HandlerSearchOutbox = "process_search_outbox"
)

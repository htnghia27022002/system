package queue

// JSON stream keys (internal/queue/nats.json → streams map).
const (
	StreamKeySearch = "search"
	StreamKeyMaps   = "maps"
)

// JSON consumer keys nested under a stream in nats.json.
const (
	ConsumerKeySearchOutbox = "search_outbox"
	ConsumerKeyMapsIngest   = "maps_ingest"
)

// JetStream stream names.
const (
	StreamSearch = "SYSTEM_SEARCH"
	StreamMaps   = "SYSTEM_MAPS"
)

// NATS subjects — use these everywhere instead of string literals.
const (
	SubjectSearchPrefix   = "system.search"
	SubjectSearchOutbox   = "system.search.outbox"
	SubjectSearchWildcard = "system.search.>"
	SubjectMapsPrefix     = "system.maps"
	SubjectMapsIngest     = "system.maps.ingest"
	SubjectMapsWildcard   = "system.maps.>"
)

// JetStream durable consumer names.
const (
	ConsumerSearchOutbox = "process_search_outbox"
	ConsumerMapsIngest   = "process_maps_ingest"
)

// Handler registry keys wired in subscribers.
const (
	HandlerSearchOutbox = "process_search_outbox"
	HandlerMapsIngest   = "process_maps_ingest"
)

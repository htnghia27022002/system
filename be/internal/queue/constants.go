package queue

// JSON stream keys (internal/queue/nats.json → streams map).
const (
	StreamKeyMaps = "maps"
)

// JSON consumer keys nested under a stream in nats.json.
const (
	ConsumerKeyMapsIngest = "maps_ingest"
)

// JetStream stream names.
const (
	StreamMaps = "SYSTEM_MAPS"
)

// NATS subjects — use these everywhere instead of string literals.
const (
	SubjectMapsPrefix   = "system.maps"
	SubjectMapsIngest   = "system.maps.ingest"
	SubjectMapsWildcard = "system.maps.>"
)

// JetStream durable consumer names.
const (
	ConsumerMapsIngest = "process_maps_ingest"
)

// Handler registry keys wired in subscribers.
const (
	HandlerMapsIngest = "process_maps_ingest"
)

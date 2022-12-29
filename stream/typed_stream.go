package stream

// TypedStream is an optional interface that defines a stream type string
// which may be used by the repositories to mark event database records.
type TypedStream interface {
	StreamType() string
}

// GetStreamType returns stream type of generic type if
// it implemented the TypedStream interface.
// Otherwise, the empty string is returned.
func GetStreamType(stream any) string {
	streamType := ""

	st, ok := stream.(TypedStream)
	if ok {
		streamType = st.StreamType()
	}

	return streamType
}

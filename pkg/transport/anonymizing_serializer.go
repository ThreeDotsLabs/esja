package transport

import (
	"github.com/ThreeDotsLabs/esja/pkg/pii"
	"github.com/ThreeDotsLabs/esja/stream"
)

type AESAnonymizingSerializer[T any] struct {
	serializer EventSerializer[T]
	anonymizer pii.StructAnonymizer[stream.ID]
}

func NewAESAnonymizingSerializer[T any](
	serializer EventSerializer[T],
	secretProvider pii.SecretProvider[stream.ID],
) *AESAnonymizingSerializer[T] {
	return &AESAnonymizingSerializer[T]{
		serializer: serializer,
		anonymizer: pii.NewStructAnonymizer[stream.ID](
			pii.NewAESAnonymizer[stream.ID](secretProvider),
		),
	}
}

func (s *AESAnonymizingSerializer[T]) Serialize(streamID stream.ID, event stream.Event[T]) ([]byte, error) {
	err := s.anonymizer.Anonymize(streamID, event)
	if err != nil {
		return nil, err
	}

	return s.serializer.Serialize(streamID, event)
}

func (s *AESAnonymizingSerializer[T]) Deserialize(streamID stream.ID, name stream.EventName, payload []byte) (stream.Event[T], error) {
	event, err := s.serializer.Deserialize(streamID, name, payload)
	if err != nil {
		return nil, err
	}

	err = s.anonymizer.Deanonymize(streamID, event)
	if err != nil {
		return nil, err
	}

	return event, nil
}

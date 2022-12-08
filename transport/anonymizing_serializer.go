package transport

import (
	pii2 "github.com/ThreeDotsLabs/esja/pii"
	"github.com/ThreeDotsLabs/esja/stream"
)

type AESAnonymizingSerializer[T any] struct {
	serializer EventSerializer[T]
	anonymizer pii2.StructAnonymizer[stream.ID, stream.Event[T]]
}

func NewAESAnonymizingSerializer[T any](
	serializer EventSerializer[T],
	secretProvider pii2.SecretProvider[stream.ID],
) *AESAnonymizingSerializer[T] {
	return &AESAnonymizingSerializer[T]{
		serializer: serializer,
		anonymizer: pii2.NewStructAnonymizer[stream.ID, stream.Event[T]](
			pii2.NewAESAnonymizer[stream.ID](secretProvider),
		),
	}
}

func (s *AESAnonymizingSerializer[T]) Serialize(streamID stream.ID, event stream.Event[T]) ([]byte, error) {
	anonymized, err := s.anonymizer.Anonymize(streamID, event)
	if err != nil {
		return nil, err
	}

	return s.serializer.Serialize(streamID, anonymized)
}

func (s *AESAnonymizingSerializer[T]) Deserialize(streamID stream.ID, name stream.EventName, payload []byte) (stream.Event[T], error) {
	event, err := s.serializer.Deserialize(streamID, name, payload)
	if err != nil {
		return nil, err
	}

	deanonymized, err := s.anonymizer.Deanonymize(streamID, event)
	if err != nil {
		return nil, err
	}

	return deanonymized, nil
}

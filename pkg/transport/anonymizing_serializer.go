package transport

import (
	"github.com/ThreeDotsLabs/esja/pkg/aggregate"
	"github.com/ThreeDotsLabs/esja/pkg/pii"
)

type AESAnonymizingSerializer[T any] struct {
	serializer EventSerializer[T]
	anonymizer pii.StructAnonymizer[aggregate.ID]
}

func NewAESAnonymizingSerializer[T any](
	serializer EventSerializer[T],
	secretProvider pii.SecretProvider[aggregate.ID],
) *AESAnonymizingSerializer[T] {
	return &AESAnonymizingSerializer[T]{
		serializer: serializer,
		anonymizer: pii.NewStructAnonymizer[aggregate.ID](
			pii.NewAESAnonymizer[aggregate.ID](secretProvider),
		),
	}
}

func (s *AESAnonymizingSerializer[T]) Serialize(aggregateID aggregate.ID, event aggregate.Event[T]) ([]byte, error) {
	err := s.anonymizer.Anonymize(aggregateID, event)
	if err != nil {
		return nil, err
	}

	return s.serializer.Serialize(aggregateID, event)
}

func (s *AESAnonymizingSerializer[T]) Deserialize(aggregateID aggregate.ID, name aggregate.EventName, payload []byte) (aggregate.Event[T], error) {
	event, err := s.serializer.Deserialize(aggregateID, name, payload)
	if err != nil {
		return nil, err
	}

	err = s.anonymizer.Deanonymize(aggregateID, event)
	if err != nil {
		return nil, err
	}

	return event, nil
}

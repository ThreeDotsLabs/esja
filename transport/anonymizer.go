package transport

import (
	"github.com/ThreeDotsLabs/esja/stream"
	"github.com/ThreeDotsLabs/pii"
)

type AESAnonymizer[T any] struct {
	mapper     Mapper[T]
	anonymizer pii.StructAnonymizer[stream.ID, any]
}

func NewAESAnonymizer[T any](
	mapper Mapper[T],
	secretProvider pii.SecretProvider[stream.ID],
) *AESAnonymizer[T] {
	return &AESAnonymizer[T]{
		mapper: mapper,
		anonymizer: pii.NewStructAnonymizer[stream.ID, any](
			pii.NewAESAnonymizer[stream.ID](secretProvider),
		),
	}
}

func (s *AESAnonymizer[T]) FromStorage(
	streamID stream.ID,
	name stream.EventName,
	payload interface{},
) (stream.Event[T], error) {
	payload, err := s.anonymizer.Deanonymize(streamID, payload)
	if err != nil {
		return nil, err
	}

	event, err := s.mapper.FromStorage(streamID, name, payload)
	if err != nil {
		return nil, err
	}

	return event, nil
}

func (s *AESAnonymizer[T]) ToStorage(
	streamID stream.ID,
	event stream.Event[T],
) (interface{}, error) {
	e, err := s.mapper.ToStorage(streamID, event)
	if err != nil {
		return nil, err
	}

	payload, err := s.anonymizer.Anonymize(streamID, e)
	if err != nil {
		return nil, err
	}

	return payload, nil
}

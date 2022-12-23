package transport

import (
	"github.com/ThreeDotsLabs/esja/stream"
	"github.com/ThreeDotsLabs/pii"
)

// AESAnonymizer is a wrapper to any transport.Mapper instance.
// AESAnonymizer will anonymize transport model properties
// if those were tagged with `anonymize:true` tag.
type AESAnonymizer[T any] struct {
	mapper     Mapper[T]
	anonymizer pii.StructAnonymizer[stream.ID, any]
}

// NewAESAnonymizer returns a new instance of AESAnonymizer.
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

func (a *AESAnonymizer[T]) New(name stream.EventName) (any, error) {
	return a.mapper.New(name)
}

func (a *AESAnonymizer[T]) FromStorage(
	streamID stream.ID,
	payload any,
) (stream.Event[T], error) {
	payload, err := a.anonymizer.Deanonymize(streamID, payload)
	if err != nil {
		return nil, err
	}

	event, err := a.mapper.FromStorage(streamID, payload)
	if err != nil {
		return nil, err
	}

	return event, nil
}

func (a *AESAnonymizer[T]) ToStorage(
	streamID stream.ID,
	event stream.Event[T],
) (any, error) {
	e, err := a.mapper.ToStorage(streamID, event)
	if err != nil {
		return nil, err
	}

	payload, err := a.anonymizer.Anonymize(streamID, &e)
	if err != nil {
		return nil, err
	}

	return payload, nil
}

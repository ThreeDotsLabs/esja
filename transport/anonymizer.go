package transport

import (
	"context"

	"github.com/ThreeDotsLabs/esja"
)

// StructAnonymizer is an interface of the anonymizer component.
type StructAnonymizer interface {
	// Anonymize encrypts struct properties using secrets
	// correlated with a provided stream ID.
	Anonymize(ctx context.Context, key string, data any) (any, error)

	// Deanonymize decrypts struct properties using secrets
	// correlated with a provided stream ID.
	Deanonymize(ctx context.Context, key string, data any) (any, error)
}

// Anonymizer is a wrapper to any transport.Mapper instance.
// Anonymizer will anonymize transport model properties
// using provided StructAnonymizer implementation.
type Anonymizer[T any] struct {
	mapper     Mapper[T]
	anonymizer StructAnonymizer
}

// NewAnonymizer returns a new instance of Anonymizer.
func NewAnonymizer[T any](
	mapper Mapper[T],
	anonymizer StructAnonymizer,
) *Anonymizer[T] {
	return &Anonymizer[T]{
		mapper:     mapper,
		anonymizer: anonymizer,
	}
}

func (a *Anonymizer[T]) New(eventName string) (any, error) {
	return a.mapper.New(eventName)
}

func (a *Anonymizer[T]) FromTransport(
	ctx context.Context,
	streamID string,
	transportEvent any,
) (esja.Event[T], error) {
	transportEvent, err := a.anonymizer.Deanonymize(ctx, streamID, transportEvent)
	if err != nil {
		return nil, err
	}

	event, err := a.mapper.FromTransport(ctx, streamID, transportEvent)
	if err != nil {
		return nil, err
	}

	return event, nil
}

func (a *Anonymizer[T]) ToTransport(
	ctx context.Context,
	streamID string,
	event esja.Event[T],
) (any, error) {
	e, err := a.mapper.ToTransport(ctx, streamID, event)
	if err != nil {
		return nil, err
	}

	payload, err := a.anonymizer.Anonymize(ctx, streamID, e)
	if err != nil {
		return nil, err
	}

	return payload, nil
}

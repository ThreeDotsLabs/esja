package pii

import "fmt"

// MaskingAnonymizer masks struct fields, losing the original value in progress.
// This might be useful for things like avoiding logging of sensitive information.
type MaskingAnonymizer[T any] struct{}

func NewMaskingAnonymizer[T any](mask T) MaskingAnonymizer[T] {
	return MaskingAnonymizer[T]{}
}

func (m MaskingAnonymizer[T]) AnonymizeString(key T, value string) (string, error) {
	return fmt.Sprintf("%v", key), nil
}

func (m MaskingAnonymizer[T]) DeanonymizeString(_ T, value string) (string, error) {
	return value, nil
}

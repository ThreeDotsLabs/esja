package pii

import (
	"reflect"
)

const (
	anonymizeTag = "anonymize"
)

// StringAnonymizer anonymizes and deanonymizes strings.
// K is the type of key used to anonymize the string.
type StringAnonymizer[K any] interface {
	AnonymizeString(key K, value string) (string, error)
	DeanonymizeString(key K, value string) (string, error)
}

// StructAnonymizer anonymizes and deanonymizes structs.
// K is the type of key used to anonymize the struct.
// T is the type of struct to be anonymized.
type StructAnonymizer[K any, T any] struct {
	stringAnonymizer StringAnonymizer[K]
}

func NewStructAnonymizer[K any, T any](
	stringAnonymizer StringAnonymizer[K],
) StructAnonymizer[K, T] {
	return StructAnonymizer[K, T]{
		stringAnonymizer: stringAnonymizer,
	}
}

func (a StructAnonymizer[K, T]) Anonymize(key K, data T) (T, error) {
	t := reflect.TypeOf(data)
	cp := reflect.New(t).Elem()
	cp.Set(reflect.ValueOf(data))

	err := a.anonymize(key, cp)
	if err != nil {
		var empty T
		return empty, err
	}

	return cp.Interface().(T), nil
}

func (a StructAnonymizer[K, T]) anonymize(key K, v reflect.Value) error {
	tv := v
	t := v.Type()
	for tv.Kind() == reflect.Ptr {
		tv = v.Elem()
		t = tv.Type()
	}

	for i := 0; i < t.NumField(); i++ {
		field := tv.Field(i)
		fieldType := t.Field(i)

		if field.Kind() == reflect.String {
			_, ok := fieldType.Tag.Lookup(anonymizeTag)
			if ok {
				anonymized, err := a.stringAnonymizer.AnonymizeString(key, field.String())
				if err != nil {
					return err
				}
				field.SetString(anonymized)
			}
		} else if field.Kind() == reflect.Struct {
			err := a.anonymize(key, field)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (a StructAnonymizer[K, T]) Deanonymize(key K, data T) (T, error) {
	t := reflect.TypeOf(data)
	cp := reflect.New(t).Elem()
	cp.Set(reflect.ValueOf(data))

	err := a.deanonymize(key, cp)

	if err != nil {
		var empty T
		return empty, err
	}

	return cp.Interface().(T), nil
}

func (a StructAnonymizer[K, T]) deanonymize(key K, v reflect.Value) error {
	tv := v
	t := v.Type()
	for tv.Kind() == reflect.Ptr {
		tv = v.Elem()
		t = tv.Type()
	}

	for i := 0; i < t.NumField(); i++ {
		field := tv.Field(i)
		fieldType := t.Field(i)

		if field.Kind() == reflect.String {
			_, ok := fieldType.Tag.Lookup(anonymizeTag)
			if ok {
				deanonymized, err := a.stringAnonymizer.DeanonymizeString(key, field.String())
				if err != nil {
					return err
				}
				field.SetString(deanonymized)
			}
		} else if field.Kind() == reflect.Struct {
			err := a.deanonymize(key, field)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

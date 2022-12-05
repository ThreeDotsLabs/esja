package pii

import (
	"reflect"
)

const (
	anonymizeTag = "anonymize"
)

type StringAnonymizer[K any] interface {
	AnonymizeString(key K, value string) (string, error)
	DeanonymizeString(key K, value string) (string, error)
}

type StructAnonymizer[K any] struct {
	stringAnonymizer StringAnonymizer[K]
}

func NewStructAnonymizer[K any](
	stringAnonymizer StringAnonymizer[K],
) StructAnonymizer[K] {
	return StructAnonymizer[K]{
		stringAnonymizer: stringAnonymizer,
	}
}

func (a StructAnonymizer[K]) Anonymize(key K, data any) error {
	v := reflect.ValueOf(data)
	return a.anonymize(key, v)
}

func (a StructAnonymizer[K]) anonymize(key K, v reflect.Value) error {
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

func (a StructAnonymizer[K]) Deanonymize(key K, data any) error {
	v := reflect.ValueOf(data)
	return a.deanonymize(key, v)
}

func (a StructAnonymizer[K]) deanonymize(key K, v reflect.Value) error {
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

package sql

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"reflect"

	"github.com/ThreeDotsLabs/esja/pkg/aggregate"
)

type Anonymizer interface {
	AnonymizeString(aggregateID aggregate.ID, value string) (string, error)
	DeanonymizeString(aggregateID aggregate.ID, value string) (string, error)
}

type AnonymizingMarshaler struct {
	marshaler  Marshaler
	anonymizer Anonymizer
}

func NewAnonymizingMarshaler(
	marshaler Marshaler,
	anonymizer Anonymizer,
) AnonymizingMarshaler {
	return AnonymizingMarshaler{
		marshaler:  marshaler,
		anonymizer: anonymizer,
	}
}

func (a AnonymizingMarshaler) Marshal(aggregateID aggregate.ID, data interface{}) ([]byte, error) {
	v := reflect.ValueOf(data)

	err := a.anonymize(v, aggregateID)
	if err != nil {
		return nil, err
	}

	return a.marshaler.Marshal(aggregateID, data)
}

func (a AnonymizingMarshaler) anonymize(v reflect.Value, aggregateID aggregate.ID) error {
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
			_, ok := fieldType.Tag.Lookup("anonymize")
			if ok {
				anonymized, err := a.anonymizer.AnonymizeString(aggregateID, field.String())
				if err != nil {
					return err
				}
				field.SetString(anonymized)
			}
		} else if field.Kind() == reflect.Struct {
			err := a.anonymize(field, aggregateID)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (a AnonymizingMarshaler) Unmarshal(aggregateID aggregate.ID, bytes []byte, target interface{}) error {
	err := a.marshaler.Unmarshal(aggregateID, bytes, target)
	if err != nil {
		return err
	}

	v := reflect.ValueOf(target)

	return a.deanonymize(v, aggregateID)
}

func (a AnonymizingMarshaler) deanonymize(v reflect.Value, aggregateID aggregate.ID) error {
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
			_, ok := fieldType.Tag.Lookup("anonymize")
			if ok {
				deanonymized, err := a.anonymizer.DeanonymizeString(aggregateID, field.String())
				if err != nil {
					return err
				}
				field.SetString(deanonymized)
			}
		} else if field.Kind() == reflect.Struct {
			err := a.deanonymize(field, aggregateID)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

type SecretProvider interface {
	SecretForAggregate(aggregateID aggregate.ID) ([]byte, error)
}

type AESAnonymizer struct {
	secretProvider SecretProvider
}

func NewAESAnonymizer(secretProvider SecretProvider) AESAnonymizer {
	return AESAnonymizer{
		secretProvider: secretProvider,
	}
}

func (a AESAnonymizer) AnonymizeString(aggregateID aggregate.ID, value string) (string, error) {
	secret, err := a.secretProvider.SecretForAggregate(aggregateID)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(secret)
	if err != nil {
		return "", err
	}

	aead, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, aead.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := aead.Seal(nonce, nonce, []byte(value), nil)
	return fmt.Sprintf("%x", ciphertext), nil
}

func (a AESAnonymizer) DeanonymizeString(aggregateID aggregate.ID, value string) (string, error) {
	secret, err := a.secretProvider.SecretForAggregate(aggregateID)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(secret)
	if err != nil {
		return "", err
	}

	aead, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	decoded, err := hex.DecodeString(value)
	if err != nil {
		return "", err
	}

	nonceSize := aead.NonceSize()
	if len(decoded) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	nonce, cipherText := decoded[:nonceSize], decoded[nonceSize:]
	data, err := aead.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

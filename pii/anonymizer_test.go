package pii_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ThreeDotsLabs/esja/pii"
)

type testStruct struct {
	FirstName string `anonymize:"true"`
	LastName  string `anonymize:"true"`
	Company   string
}

func TestStructAnonymizer(t *testing.T) {
	a := pii.NewStructAnonymizer[string, testStruct](testStringAnonymizer{})

	s := testStruct{
		FirstName: "John",
		LastName:  "Doe",
		Company:   "ThreeDotsLabs",
	}

	anonymized, err := a.Anonymize("id", s)
	require.NoError(t, err)

	assert.Equal(t, "anonymized.id.John", anonymized.FirstName)
	assert.Equal(t, "anonymized.id.Doe", anonymized.LastName)
	assert.Equal(t, "ThreeDotsLabs", anonymized.Company)

	deanonymized, err := a.Deanonymize("id", anonymized)
	require.NoError(t, err)

	assert.Equal(t, "John", deanonymized.FirstName)
	assert.Equal(t, "Doe", deanonymized.LastName)
	assert.Equal(t, "ThreeDotsLabs", deanonymized.Company)
}

type testStringAnonymizer struct{}

func (t testStringAnonymizer) AnonymizeString(key string, value string) (string, error) {
	return fmt.Sprintf("anonymized.%s.%s", key, value), nil
}

func (t testStringAnonymizer) DeanonymizeString(key string, value string) (string, error) {
	parts := strings.Split(value, ".")
	if len(parts) != 3 {
		return "", fmt.Errorf("invalid value")
	}
	if parts[1] != key {
		return "", fmt.Errorf("invalid key")
	}
	return parts[2], nil
}

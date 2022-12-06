package pii_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"

	"github.com/ThreeDotsLabs/esja/pkg/pii"
)

type testStruct struct {
	FirstName string `anonymize:"true"`
	LastName  string `anonymize:"true"`
	Company   string
}

func TestStructAnonymizer(t *testing.T) {
	a := pii.NewStructAnonymizer[string](testStringAnonymizer{})

	s := &testStruct{
		FirstName: "John",
		LastName:  "Doe",
		Company:   "ThreeDotsLabs",
	}

	err := a.Anonymize("id", s)
	require.NoError(t, err)

	assert.Equal(t, "anonymized.id.John", s.FirstName)
	assert.Equal(t, "anonymized.id.Doe", s.LastName)
	assert.Equal(t, "ThreeDotsLabs", s.Company)

	err = a.Deanonymize("id", s)
	require.NoError(t, err)

	assert.Equal(t, "John", s.FirstName)
	assert.Equal(t, "Doe", s.LastName)
	assert.Equal(t, "ThreeDotsLabs", s.Company)
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

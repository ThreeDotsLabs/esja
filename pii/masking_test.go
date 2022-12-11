package pii_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ThreeDotsLabs/esja/pii"
)

func TestMaskingAnonymizer(t *testing.T) {
	a := pii.NewStructAnonymizer[string, testStruct](pii.MaskingAnonymizer[string]{})

	s := testStruct{
		FirstName: "John",
		LastName:  "Doe",
		Company:   "ThreeDotsLabs",
	}

	anonymized, err := a.Anonymize("***", s)
	require.NoError(t, err)

	assert.Equal(t, "***", anonymized.FirstName)
	assert.Equal(t, "***", anonymized.LastName)
	assert.Equal(t, "ThreeDotsLabs", anonymized.Company)

	deanonymized, err := a.Deanonymize("", anonymized)
	require.NoError(t, err)

	assert.Equal(t, "***", deanonymized.FirstName)
	assert.Equal(t, "***", deanonymized.LastName)
	assert.Equal(t, "ThreeDotsLabs", deanonymized.Company)
}

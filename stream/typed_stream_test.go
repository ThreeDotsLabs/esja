package stream_test

import (
	"testing"

	"github.com/ThreeDotsLabs/esja/stream"
	"github.com/stretchr/testify/assert"
)

type testStructA struct {
}

type testStructB struct {
}

func (t *testStructB) StreamType() string {
	return "TestStructB"
}

type testStructC struct {
}

func (t *testStructC) StreamType() string {
	return "TestStructC"
}

func TestGetStreamType(t *testing.T) {
	assert.Equal(t, "", stream.GetStreamType[testStructA]())
	assert.Equal(t, "TestStructB", stream.GetStreamType[testStructB]())
	assert.Equal(t, "TestStructC", stream.GetStreamType[testStructC]())
}

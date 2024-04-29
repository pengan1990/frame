package core

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestReader struct {
}

func (t *TestReader) Execute() (interface{}, error) {
	return 1, nil
}

func (t *TestReader) Next(ctx context.Context) {
	select {
	case <-ctx.Done():
		return
	default:
		return
	}
}

func (t *TestReader) IsDone() bool {
	return false
}

type TestWriter struct {
}

func (t *TestWriter) Execute(data interface{}) error {
	return nil
}

func (t *TestWriter) Next(ctx context.Context) {
	select {
	case <-ctx.Done():
		return
	default:
		return
	}
}

type TestCheckPointer struct {
}

func (t *TestCheckPointer) Save(data interface{}) error {
	return nil
}

func TestNewFrame(t *testing.T) {
	reader, writer, chk := &TestReader{}, &TestWriter{}, &TestCheckPointer{}
	ctx := context.Background()
	size := 10
	frame, err := NewFrame(reader, []WritePair{writer}, chk, size, ctx)
	assert.NoError(t, err, "new frame error")
	err = frame.Execute()
	assert.NoError(t, err, "frame work execute with error")
}

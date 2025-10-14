package io

import (
	"bufio"
	"fmt"
	"io"

	"github.com/shivakuppa/Go_Redis/internals/resp"
)

type RespWriter struct {
	writer *bufio.Writer
}

func NewRespWriter(w io.Writer) *RespWriter {
	return &RespWriter{writer: bufio.NewWriter(w)}
}

func (w *RespWriter) Write(v *resp.Value) error {
	reply, err := resp.Serialize(v)
	if err != nil {
		return fmt.Errorf("serialize value: %w", err)
	}

	if _, err := w.writer.Write([]byte(reply)); err != nil {
		return fmt.Errorf("write to buffer: %w", err)
	}

	return nil
}

func (w *RespWriter) Flush() error {
	return w.writer.Flush()
}

package server

import (
	"bufio"
	"fmt"
	"io"

	"github.com/shivakuppa/Go_Redis/internals/resp"
)

type Writer struct {
	writer *bufio.Writer
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{writer: bufio.NewWriter(w)}
}

func (w *Writer) Write(v *resp.Value) error {
	reply, err := resp.Serialize(v)
	if err != nil {
		return fmt.Errorf("serialize value: %w", err)
	}

	if _, err := w.writer.Write([]byte(reply)); err != nil {
		return fmt.Errorf("write to buffer: %w", err)
	}

	return nil
}

func (w *Writer) Flush() error {
	return w.writer.Flush()
}

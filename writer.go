package subcmp

import (
	"encoding/json"
	"fmt"
	subtest "github.com/streamingfast/substreams/tools/test"
	"io"
	"os"
)

type Writer struct {
	f    *os.File
	in   chan *subtest.TestConfig
	done chan bool
}

func newWriter(outputfile string) (*Writer, error) {
	f, err := os.Create(outputfile)
	if err != nil {
		return nil, fmt.Errorf("failed to create file %q: %w", outputfile, err)
	}

	return &Writer{
		f:    f,
		in:   make(chan *subtest.TestConfig, 10000),
		done: make(chan bool, 1),
	}, nil
}

func (w *Writer) wait() {
	<-w.done
}

func (w *Writer) closeTest() {
	close(w.in)
}

func (w *Writer) Write() error {
	defer func() {
		w.done <- true
	}()
	for {
		test, ok := <-w.in
		if !ok {
			w.close()
			return io.EOF
		}

		cnt, err := json.Marshal(test)
		if err != nil {
			return fmt.Errorf("failed to marshall test: %w", err)
		}
		if _, err := w.f.Write(cnt); err != nil {
			return fmt.Errorf("failed to write test: %w", err)
		}
		if _, err := w.f.Write([]byte{'\n'}); err != nil {
			return fmt.Errorf("failed to write test: %w", err)
		}
	}
}

func (w *Writer) close() {
	w.f.Close()
}

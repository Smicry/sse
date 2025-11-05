// Copyright 2025 SSE Project Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sse

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
)

// Event holds all of the event source fields
type Event struct {
	ID    []byte
	Data  []byte
	Event []byte
	Retry []byte
}

// WriteEvent encodes an event to the sse format, and writes it to the writer.
func WriteEvent(writer io.Writer, event Event) error {
	var buf bytes.Buffer

	if err := writeId(&buf, event.ID); err != nil {
		return fmt.Errorf("write id: %w", err)
	}
	if err := writeEvent(&buf, event.Event); err != nil {
		return fmt.Errorf("write event: %w", err)
	}
	if err := writeRetry(&buf, event.Retry); err != nil {
		return fmt.Errorf("write retry: %w", err)
	}
	if err := writeData(&buf, event.Data); err != nil {
		return fmt.Errorf("write data: %w", err)
	}

	buf.WriteString("\n")
	_, err := writer.Write(buf.Bytes())
	return err
}

func writeId(w io.Writer, id []byte) error {
	if len(id) == 0 {
		return nil
	}
	if _, err := w.Write([]byte("id:")); err != nil {
		return err
	}
	if _, err := w.Write(id); err != nil {
		return err
	}
	_, err := w.Write([]byte("\n"))
	return err
}

func writeEvent(w io.Writer, event []byte) error {
	if len(event) == 0 {
		return nil
	}
	if _, err := w.Write([]byte("event:")); err != nil {
		return err
	}
	if _, err := w.Write(event); err != nil {
		return err
	}
	_, err := w.Write([]byte("\n"))
	return err
}

func writeRetry(w io.Writer, retry []byte) error {
	retryUint, err := strconv.ParseUint(string(retry), 10, 64)
	if err != nil {
		return nil
	}
	if retryUint == 0 {
		return nil
	}
	if _, err := w.Write([]byte("retry:")); err != nil {
		return err
	}
	if _, err := w.Write(retry); err != nil {
		return err
	}
	_, err = w.Write([]byte("\n"))
	return err
}

func writeData(w io.Writer, data []byte) error {
	if _, err := w.Write([]byte("data:")); err != nil {
		return err
	}
	if _, err := w.Write(data); err != nil {
		return err
	}
	_, err := w.Write([]byte("\n"))
	return err
}

package mem_constrained_buffer

import (
	"bytes"
	"io"
	"io/ioutil"
	"mime/multipart"
	"os"
)

const (
	DefaultMemorySize int64 = 1 << 18
	FilenamePrefix          = "mem-buf-"
)

type MemoryConstrainedBuffer struct {
	b             *bytes.Buffer
	tmpfile       string
	max           int64
	size          int64
	removeOnClose bool
	file          multipart.File
}

func NewMemoryConstrainedBuffer(maxMemory int64, removeOnClose bool) *MemoryConstrainedBuffer {
	return &MemoryConstrainedBuffer{
		max:           maxMemory,
		removeOnClose: removeOnClose,
	}
}

func NewDefaultMemoryConstrainedBuffer() *MemoryConstrainedBuffer {
	return NewMemoryConstrainedBuffer(DefaultMemorySize, true)
}

func (m *MemoryConstrainedBuffer) Write(p []byte) (int, error) {
	n, err := m.ReadFrom(bytes.NewReader(p))
	return int(n), err
}

func (m *MemoryConstrainedBuffer) ReadFrom(r io.Reader) (int64, error) {
	var (
		n   int64
		err error
		buf bytes.Buffer
	)

	for {
		n, err = io.CopyN(&buf, r, m.max+1)
		if err != nil && err != io.EOF {
			return 0, err
		}

		m.size += n

		if err == io.EOF {
			err = nil
			break
		}

		if n > m.max {
			// too big, write to disk and flush buffer
			file, err := ioutil.TempFile("", FilenamePrefix)
			if err != nil {
				return 0, err
			}
			defer file.Close()
			n, err = io.Copy(file, io.MultiReader(&buf, r))
			if err != nil {
				os.Remove(file.Name())
				return 0, err
			}
			m.tmpfile = file.Name()
			m.size = n
			break

		} else {
			m.max -= n
		}
	}

	if m.tmpfile == "" {
		m.b = &buf
	}

	return n, err
}

func (m *MemoryConstrainedBuffer) Len() int64 {
	return m.size
}

func (m *MemoryConstrainedBuffer) Remove() (err error) {
	if m.file == nil && m.tmpfile == "" {
		return nil
	}
	if m.file != nil {
		err = m.file.Close()
	}
	if m.tmpfile != "" {
		err = os.Remove(m.tmpfile)
	}
	return err
}

func (m *MemoryConstrainedBuffer) open() error {
	if m.file != nil {
		return nil
	}
	if b := m.b; b != nil {
		m.file = &sectionReadCloser{io.NewSectionReader(bytes.NewReader(m.b.Bytes()), 0, int64(b.Len()))}
		return nil
	}
	f, err := os.Open(m.tmpfile)
	m.file = f
	return err
}

func (m *MemoryConstrainedBuffer) Read(p []byte) (int, error) {
	if m.file == nil {
		if err := m.open(); err != nil {
			return 0, err
		}
	}
	return m.file.Read(p)
}

func (m *MemoryConstrainedBuffer) ReadAt(p []byte, off int64) (int, error) {
	if m.file == nil {
		if err := m.open(); err != nil {
			return 0, err
		}
	}
	return m.file.ReadAt(p, off)
}

func (m *MemoryConstrainedBuffer) Seek(offset int64, whence int) (int64, error) {
	if m.file == nil {
		if err := m.open(); err != nil {
			return 0, err
		}
	}
	return m.file.Seek(offset, whence)
}

func (m *MemoryConstrainedBuffer) Close() error {
	if m.removeOnClose {
		defer m.Remove()
	}
	if m.file == nil {
		return nil
	}
	return m.file.Close()
}

type sectionReadCloser struct {
	*io.SectionReader
}

func (rc sectionReadCloser) Close() error {
	return nil
}

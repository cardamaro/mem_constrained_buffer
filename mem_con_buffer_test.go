package mem_constrained_buffer

import (
	"os"
	"strings"
	"testing"
)

func TestBufferWrite(t *testing.T) {
	var err error
	var n int

	buf := NewWithSize(10, true)

	n, err = buf.Write([]byte("1234567890abcdefghijklmnopqrstuvwxyz"))
	if n != 36 {
		t.Errorf("wrong read bytes: %d", n)
	}
	if err != nil {
		t.Error(err)
	}

	err = buf.Close()
	if err != nil {
		t.Error(err)
	}

	buf2 := NewWithSize(40, true)

	n, err = buf2.Write([]byte("1234567890abcdefghijklmnopqrstuvwxyz"))
	if n != 36 {
		t.Errorf("wrong read bytes: %d", n)
	}
	if err != nil {
		t.Error(err)
	}

	err = buf2.Close()
	if err != nil {
		t.Error(err)
	}
}

func TestBufferReadFrom(t *testing.T) {
	var err error
	var n int64

	buf := NewWithSize(10, true)

	n, err = buf.ReadFrom(strings.NewReader("1234567890abcdefghijklmnopqrstuvwxyz"))
	if n != 36 {
		t.Errorf("wrong read bytes: %d", n)
	}
	if err != nil {
		t.Error(err)
	}

	err = buf.Close()
	if err != nil {
		t.Error(err)
	}

	buf2 := NewWithSize(40, true)

	n, err = buf2.ReadFrom(strings.NewReader("1234567890abcdefghijklmnopqrstuvwxyz"))
	if n != 36 {
		t.Errorf("wrong read bytes: %d", n)
	}
	if err != nil {
		t.Error(err)
	}

	err = buf2.Close()
	if err != nil {
		t.Error(err)
	}
}

func TestBufferLen(t *testing.T) {
	buf := NewWithSize(10, true)

	_, err := buf.ReadFrom(strings.NewReader("1234567890abcdefghijklmnopqrstuvwxyz"))
	if err != nil {
		t.Error(err)
	}

	if buf.Len() != 36 {
		t.Errorf("Len() is wrong: %d", buf.Len())
	}

	err = buf.Close()
	if err != nil {
		t.Error(err)
	}
}

func TestBufferReadSmall(t *testing.T) {
	str := "1234567890abcdefghijklmnopqrstuvwxyz"

	buf := NewWithSize(10, true)

	_, err := buf.ReadFrom(strings.NewReader(str))
	if err != nil {
		t.Error(err)
	}

	b := make([]byte, 36)
	n2, err := buf.Read(b)
	if n2 != 36 {
		t.Errorf("wrong read bytes: %d", n2)
	}
	if err != nil {
		t.Error(err)
	}

	if string(b) != str {
		t.Errorf("string does not match: %s", b)
	}

	err = buf.Close()
	if err != nil {
		t.Error(err)
	}
}

func TestBufferReadLarge(t *testing.T) {
	str := "1234567890abcdefghijklmnopqrstuvwxyz"

	buf := NewWithSize(100, true)

	_, err := buf.ReadFrom(strings.NewReader(str))
	if err != nil {
		t.Error(err)
	}

	b := make([]byte, 36)
	n2, err := buf.Read(b)
	if n2 != 36 {
		t.Errorf("wrong read bytes: %d", n2)
	}
	if err != nil {
		t.Error(err)
	}

	if string(b) != str {
		t.Errorf("string does not match: %s", b)
	}

	err = buf.Close()
	if err != nil {
		t.Error(err)
	}
}

func TestBufferCloseTempCleanup(t *testing.T) {
	buf := NewWithSize(4, true)

	_, err := buf.ReadFrom(strings.NewReader("12345"))
	if err != nil {
		t.Error(err)
	}

	err = buf.Remove()
	if err != nil {
		t.Error(err)
	}

	if _, err := os.Stat(buf.tmpfile); err == nil {
		t.Errorf("tmpfile not removed")
	}
}

package colorable

import (
	"bytes"
	"io"
)

// NonColorable holds writer but removes escape sequence.
type NonColorable struct {
	out io.Writer
}

// NewNonColorable returns new instance of Writer which removes escape sequence from Writer.
func NewNonColorable(w io.Writer) io.Writer {
	return &NonColorable{out: w}
}

// Write writes data on console
func (w *NonColorable) Write(data []byte) (n int, err error) {
	er := bytes.NewReader(data)
	var plaintext bytes.Buffer
loop:
	for {
		c1, err := er.ReadByte()
		if err != nil {
			plaintext.WriteTo(w.out)
			break loop
		}
		if c1 != 0x1b {
			plaintext.WriteByte(c1)
			continue
		}
		_, err = plaintext.WriteTo(w.out)
		if err != nil {
			break loop
		}
		c2, err := er.ReadByte()
		if err != nil {
			break loop
		}
		if c2 != 0x5b {
			continue
		}

<<<<<<< HEAD
=======
		var buf bytes.Buffer
>>>>>>> 1cb7c9a8c04b7de79c2dd46f84bd5239eed4ee16
		for {
			c, err := er.ReadByte()
			if err != nil {
				break loop
			}
			if ('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z') || c == '@' {
				break
			}
<<<<<<< HEAD
=======
			buf.Write([]byte(string(c)))
>>>>>>> 1cb7c9a8c04b7de79c2dd46f84bd5239eed4ee16
		}
	}

	return len(data), nil
}

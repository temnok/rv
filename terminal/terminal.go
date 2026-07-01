package terminal

import (
	"io"
)

type Terminal struct {
	stdin chan byte
	out   io.Writer

	Closed bool
}

const ctrlC = 3

func New(in io.Reader, out io.Writer) *Terminal {
	t := &Terminal{
		stdin: make(chan byte),
		out:   out,
	}

	go func() {
		for {
			buf := []byte{0}
			if n, _ := in.Read(buf); n > 0 {
				t.stdin <- buf[0]

				if buf[0] == ctrlC {
					break
				}
			}
		}
	}()

	return t
}

func (t *Terminal) PutChar(char byte) {
	t.out.Write([]byte{char})
}

func (t *Terminal) GetChar() (byte, bool) {
	select {
	case ch := <-t.stdin:
		if ch == ctrlC {
			t.Closed = true
		}

		return ch, true

	default:
		return 0, false
	}
}

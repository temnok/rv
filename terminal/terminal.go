package terminal

import (
	"golang.org/x/term"
	"io"
	"os"
)

type Terminal struct {
	stdin chan byte
	out   io.Writer

	Closed bool
}

const ctrlD = 4

func WithRaw(f func(stdin io.Reader, stdout io.Writer)) {
	state, err := term.MakeRaw(0)
	if err != nil {
		return
	}

	defer func() {
		term.Restore(0, state)
	}()

	f(os.Stdin, os.Stdout)
}

func New(in io.Reader, out io.Writer) *Terminal {
	t := &Terminal{
		stdin: make(chan byte),
		out:   out,
	}

	go func() {
		buf := []byte{0}

		for {
			if n, _ := in.Read(buf); n > 0 {
				t.stdin <- buf[0]

				if buf[0] == ctrlD {
					break
				}
			}
		}
	}()

	return t
}

func (t *Terminal) PutChar(char int) {
	t.out.Write([]byte{byte(char)})
}

func (t *Terminal) GetChar() int {
	select {
	case ch := <-t.stdin:
		if ch == ctrlD {
			t.Closed = true
		}

		return int(ch)

	default:
		return -1
	}
}

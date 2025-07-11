package rv

import (
	"io"
)

type Terminal struct {
	stdin chan byte
	out   io.Writer

	Closed bool
}

const ctrlC = 3

func newTerminal(in io.Reader, out io.Writer) *Terminal {
	t := &Terminal{
		stdin: make(chan byte),
		out:   out,
	}

	go func() {
		for {
			buf := []byte{0}
			if check1(in.Read(buf)) > 0 {
				t.stdin <- buf[0]

				if buf[0] == ctrlC {
					break
				}
			}
		}
	}()

	return t
}

func (t *Terminal) callback(data *byte, write bool) bool {
	if write {
		check1(t.out.Write([]byte{*data}))
		return true
	}

	select {
	case ch := <-t.stdin:
		if ch == ctrlC {
			t.Closed = true
		}

		*data = ch
		return true

	default:
		return false
	}
}

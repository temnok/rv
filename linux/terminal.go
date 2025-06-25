package main

import "os"

type Terminal struct {
	stdin chan byte

	Closed bool
}

const ctrlC = 3

func newTerminal() *Terminal {
	t := &Terminal{
		stdin: make(chan byte),
	}

	go func() {
		for {
			buf := []byte{0}
			if check1(os.Stdin.Read(buf)) > 0 {
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
		os.Stdout.Write([]byte{*data})
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

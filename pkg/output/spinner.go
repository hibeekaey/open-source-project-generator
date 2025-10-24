package output

import (
	"fmt"
	"time"
)

type Spinner struct {
	message string
	done    chan bool
}

func NewSpinner(message string) *Spinner {
	return &Spinner{
		message: message,
		done:    make(chan bool),
	}
}

func (s *Spinner) Start() {
	go func() {
		spinner := []rune{'⠙', '⠹', '⠸', '⠼', '⠴', '⠦', '⠧', '⠇', '⠏'}
		i := 0
		fmt.Print("\n")
		for {
			select {
			case <-s.done:
				fmt.Print("\r\033[K")
				return
			default:
				fmt.Printf("\r"+ColorCyan+"%c %s"+ColorReset, spinner[i%len(spinner)], s.message)
				i++
				time.Sleep(80 * time.Millisecond)
			}
		}
	}()
}

func (s *Spinner) Stop() {
	s.done <- true
	time.Sleep(100 * time.Millisecond)
}

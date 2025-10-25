package input

import (
	"fmt"
	"os"

	"github.com/cuesoftinc/open-source-project-generator/pkg/output"
	"golang.org/x/term"
)

func MultiSelect(prompt string, options []string) ([]string, error) {
	selected := make([]bool, len(options))
	for i := range selected {
		selected[i] = true
	}
	cursor := 0

	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = term.Restore(int(os.Stdin.Fd()), oldState)
	}()

	lines := len(options) + 3

	render := func() {
		fmt.Print("\033[" + fmt.Sprintf("%d", lines) + "A")
		fmt.Print("\r" + output.ColorCyan + prompt + output.ColorReset + "\033[K\r\n")
		fmt.Print("\r" + output.ColorYellow + "↑/↓: navigate | Space: toggle | a: toggle all | Enter: confirm" + output.ColorReset + "\033[K\r\n")
		fmt.Print("\r\033[K\r\n")

		for i, option := range options {
			checkbox := "[ ]"
			if selected[i] {
				checkbox = "[✓]"
			}

			if i == cursor {
				fmt.Print("\r" + output.ColorCyan + "❯ " + checkbox + " " + option + output.ColorReset + "\033[K\r\n")
			} else {
				fmt.Print("\r  " + checkbox + " " + option + "\033[K\r\n")
			}
		}
	}

	fmt.Print("\r\n")
	fmt.Print(output.ColorCyan + prompt + output.ColorReset + "\r\n")
	fmt.Print(output.ColorYellow + "↑/↓: navigate | Space: toggle | a: toggle all | Enter: confirm" + output.ColorReset + "\r\n")
	fmt.Print("\r\n")

	for i, option := range options {
		checkbox := "[ ]"
		if selected[i] {
			checkbox = "[✓]"
		}

		if i == cursor {
			fmt.Print(output.ColorCyan + "❯ " + checkbox + " " + option + output.ColorReset + "\r\n")
		} else {
			fmt.Print("  " + checkbox + " " + option + "\r\n")
		}
	}

	buf := make([]byte, 3)
	for {
		n, err := os.Stdin.Read(buf)
		if err != nil {
			return nil, err
		}

		if n == 0 {
			continue
		}

		switch {
		case buf[0] == 13 || buf[0] == 10:
			_ = term.Restore(int(os.Stdin.Fd()), oldState)

			fmt.Print("\033[" + fmt.Sprintf("%d", lines) + "A")

			for i := 0; i < lines; i++ {
				fmt.Print("\r\033[K\r\n")
			}

			fmt.Print("\033[" + fmt.Sprintf("%d", lines) + "A")
			fmt.Print("\r" + output.ColorCyan + prompt + output.ColorReset + "\r\n")
			fmt.Print("\r\n")

			result := []string{}
			for i, sel := range selected {
				if sel {
					result = append(result, options[i])
				}
			}

			for _, item := range result {
				fmt.Print("\r" + output.ColorGreen + "  ✓ " + item + output.ColorReset + "\r\n")
			}

			return result, nil

		case buf[0] == 32:
			selected[cursor] = !selected[cursor]
			render()

		case buf[0] == 'a' || buf[0] == 'A':
			allSelected := true
			for _, sel := range selected {
				if !sel {
					allSelected = false
					break
				}
			}
			for i := range selected {
				selected[i] = !allSelected
			}
			render()

		case buf[0] == 27 && n == 3 && buf[1] == 91:
			switch buf[2] {
			case 65:
				if cursor > 0 {
					cursor--
				}
				render()
			case 66:
				if cursor < len(options)-1 {
					cursor++
				}
				render()
			}
		}
	}
}

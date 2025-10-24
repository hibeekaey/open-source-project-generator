package input

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/cuesoftinc/open-source-project-generator/pkg/constants"
	"github.com/cuesoftinc/open-source-project-generator/pkg/output"
)

func SelectApps(reader *bufio.Reader) ([]string, error) {
	fmt.Println()
	fmt.Print(output.ColorCyan + "Select Next.js apps to create (comma-separated, or press Enter for all): " + output.ColorReset)
	fmt.Println()

	for i, app := range constants.NextJSApps {
		fmt.Printf("  %d. %s\n", i+1, app)
	}

	fmt.Print(output.ColorCyan + "Your choice (e.g., 1,2,3 or press Enter for all): " + output.ColorReset)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input == "" {
		return constants.NextJSApps, nil
	}

	choices := strings.Split(input, ",")
	selectedApps := []string{}

	for _, choice := range choices {
		choice = strings.TrimSpace(choice)
		var index int
		_, err := fmt.Sscanf(choice, "%d", &index)
		if err != nil || index < 1 || index > len(constants.NextJSApps) {
			return nil, output.NewError("invalid choice: %s", choice)
		}
		selectedApps = append(selectedApps, constants.NextJSApps[index-1])
	}

	return selectedApps, nil
}

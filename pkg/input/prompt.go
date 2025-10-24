package input

import (
	"bufio"
	"fmt"
	"regexp"
	"strings"

	"github.com/cuesoftinc/open-source-project-generator/pkg/output"
)

type ProjectInput struct {
	Name         string
	OutputFolder string
}

func ReadProjectInput(reader *bufio.Reader, defaultOutputFolder string) (*ProjectInput, error) {
	fmt.Print(output.ColorCyan + "Project name: " + output.ColorReset)
	projectName, _ := reader.ReadString('\n')
	projectName = strings.TrimSpace(projectName)

	if projectName == "" {
		return nil, output.NewError("project name cannot be empty")
	}

	validName := regexp.MustCompile(`^[a-z0-9-]+$`)
	if !validName.MatchString(projectName) {
		return nil, output.NewError("project name must contain only lowercase letters, numbers, and hyphens")
	}

	fmt.Print(output.ColorCyan + "Output folder (default: " + defaultOutputFolder + "): " + output.ColorReset)
	outputFolder, _ := reader.ReadString('\n')
	outputFolder = strings.TrimSpace(outputFolder)

	if outputFolder == "" {
		outputFolder = defaultOutputFolder
	}

	return &ProjectInput{
		Name:         projectName,
		OutputFolder: outputFolder,
	}, nil
}

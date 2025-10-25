package input

import (
	"bufio"
	"fmt"
	"regexp"
	"strings"

	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
	"github.com/cuesoftinc/open-source-project-generator/pkg/output"
)

type ProjectInput struct {
	Name            string
	OutputFolder    string
	SelectedFolders []string
}

func ReadProjectInput(reader *bufio.Reader, defaultOutputFolder string) (*ProjectInput, models.Apps, error) {
	fmt.Print(output.ColorCyan + "Project name: " + output.ColorReset)
	projectName, _ := reader.ReadString('\n')
	projectName = strings.TrimSpace(projectName)

	if projectName == "" {
		return nil, models.Apps{}, fmt.Errorf("project name cannot be empty")
	}

	validName := regexp.MustCompile(`^[a-z0-9-]+$`)
	if !validName.MatchString(projectName) {
		return nil, models.Apps{}, fmt.Errorf("project name must contain only lowercase letters, numbers, and hyphens")
	}

	componentSelection, err := ReadComponentSelection()
	if err != nil {
		return nil, models.Apps{}, fmt.Errorf("failed to read component selection: %w", err)
	}

	fmt.Print("\n" + output.ColorCyan + "Output folder (default: " + defaultOutputFolder + "): " + output.ColorReset)
	outputFolder, _ := reader.ReadString('\n')
	outputFolder = strings.TrimSpace(outputFolder)

	fmt.Println()

	if outputFolder == "" {
		outputFolder = defaultOutputFolder
	}

	projectInput := &ProjectInput{
		Name:            projectName,
		OutputFolder:    outputFolder,
		SelectedFolders: componentSelection.Folders,
	}

	return projectInput, componentSelection.Apps, nil
}

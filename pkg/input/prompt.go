package input

import (
	"bufio"
	"fmt"
	"regexp"
	"strings"

	"github.com/cuesoftinc/open-source-project-generator/pkg/constants"
	"github.com/cuesoftinc/open-source-project-generator/pkg/output"
)

type ProjectInput struct {
	Name            string
	OutputFolder    string
	SelectedFolders []string
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

	fmt.Println()
	fmt.Print(output.ColorCyan + "Select folders to create (comma-separated, or press Enter for all): " + output.ColorReset)
	fmt.Println()

	for i, folder := range constants.ProjectFolders {
		fmt.Printf("  %d. %s\n", i+1, folder)
	}

	fmt.Print(output.ColorCyan + "Your choice (e.g., 1,2,3 or press Enter for all): " + output.ColorReset)
	folderInput, _ := reader.ReadString('\n')
	folderInput = strings.TrimSpace(folderInput)

	var selectedFolders []string
	if folderInput == "" {
		selectedFolders = constants.ProjectFolders
	} else {
		choices := strings.Split(folderInput, ",")
		for _, choice := range choices {
			choice = strings.TrimSpace(choice)
			var index int
			_, err := fmt.Sscanf(choice, "%d", &index)
			if err != nil || index < 1 || index > len(constants.ProjectFolders) {
				return nil, output.NewError("invalid choice: %s", choice)
			}
			selectedFolders = append(selectedFolders, constants.ProjectFolders[index-1])
		}
	}

	fmt.Println()
	fmt.Print(output.ColorCyan + "Output folder (default: " + defaultOutputFolder + "): " + output.ColorReset)
	outputFolder, _ := reader.ReadString('\n')
	outputFolder = strings.TrimSpace(outputFolder)

	if outputFolder == "" {
		outputFolder = defaultOutputFolder
	}

	return &ProjectInput{
		Name:            projectName,
		OutputFolder:    outputFolder,
		SelectedFolders: selectedFolders,
	}, nil
}

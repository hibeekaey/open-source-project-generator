package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"slices"

	"github.com/cuesoftinc/open-source-project-generator/internal/generator"
	"github.com/cuesoftinc/open-source-project-generator/pkg/config"
	"github.com/cuesoftinc/open-source-project-generator/pkg/constants"
	"github.com/cuesoftinc/open-source-project-generator/pkg/filesystem"
	"github.com/cuesoftinc/open-source-project-generator/pkg/input"
	"github.com/cuesoftinc/open-source-project-generator/pkg/models"
	"github.com/cuesoftinc/open-source-project-generator/pkg/output"
)

func runGenerate() {
	versions, err := config.LoadVersions()
	if err != nil {
		fmt.Fprintf(os.Stderr, output.ColorRed+"Error: %v\n"+output.ColorReset, err)
		os.Exit(1)
	}

	var configFile string
	for i := 2; i < len(os.Args); i++ {
		if os.Args[i] == "--config-file" || os.Args[i] == "-c" {
			if i+1 < len(os.Args) {
				configFile = os.Args[i+1]
			} else {
				fmt.Fprintf(os.Stderr, output.ColorRed+"Error: %s requires a file path\n"+output.ColorReset, os.Args[i])
				fmt.Fprintln(os.Stderr, "Run 'generator generate --help' for usage.")
				os.Exit(1)
			}
			break
		}
	}

	var projectInput *input.ProjectInput
	var selectedApps models.Apps

	if configFile != "" {
		fmt.Printf(output.ColorCyan+"Using config file: %s"+output.ColorReset+"\n\n", configFile)

		project, err := config.LoadProject(configFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, output.ColorRed+"Error: %v\n"+output.ColorReset, err)
			os.Exit(1)
		}

		projectInput = &input.ProjectInput{
			Name:            project.ProjectName,
			OutputFolder:    project.OutputFolder,
			SelectedFolders: project.Folders,
		}
		selectedApps = project.Apps
	} else {
		reader := bufio.NewReader(os.Stdin)
		projectInput, selectedApps, err = input.ReadProjectInput(reader, constants.DefaultOutputFolder)
		if err != nil {
			fmt.Fprintf(os.Stderr, output.ColorRed+"Error: %v\n"+output.ColorReset, err)
			os.Exit(1)
		}
	}

	projectPath := filepath.Join(projectInput.OutputFolder, projectInput.Name)

	if err := filesystem.CreateProjectStructure(projectPath, projectInput.SelectedFolders); err != nil {
		fmt.Fprintf(os.Stderr, output.ColorRed+"Error: %v\n"+output.ColorReset, err)
		os.Exit(1)
	}

	if slices.Contains(projectInput.SelectedFolders, constants.ComponentFrontend) && len(selectedApps.Frontend) > 0 {
		frontendGen := &generator.FrontendGenerator{
			Version:    versions.Frontend.NextJS.Version,
			ProjectDir: projectPath,
			Component:  constants.ComponentFrontend,
			Apps:       selectedApps.Frontend,
		}

		if err := frontendGen.Generate(projectInput.Name); err != nil {
			fmt.Fprintf(os.Stderr, output.ColorRed+"Error: %v\n"+output.ColorReset, err)
			os.Exit(1)
		}
	}

	output.PrintSuccess(projectInput.Name, projectPath)
}

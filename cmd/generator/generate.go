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
	"github.com/cuesoftinc/open-source-project-generator/pkg/output"
)

func runGenerate() {
	versions, err := config.LoadVersions()
	if err != nil {
		fmt.Printf("%v\n", output.NewError("error loading versions: %v", err))
		os.Exit(1)
	}

	var configFile string
	for i := 2; i < len(os.Args); i++ {
		if os.Args[i] == "--config-file" || os.Args[i] == "-c" {
			if i+1 < len(os.Args) {
				configFile = os.Args[i+1]
			}
			break
		}
	}

	var projectInput *input.ProjectInput
	var selectedApps []string

	if configFile != "" {
		project, err := config.LoadProject(configFile)
		if err != nil {
			fmt.Printf("%v\n", output.NewError("failed to load config: %v", err))
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
		projectInput, err = input.ReadProjectInput(reader, constants.DefaultOutputFolder)
		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}

		if slices.Contains(projectInput.SelectedFolders, constants.FolderApp) {
			selectedApps, err = input.MultiSelect("Select Next.js apps to create:", constants.NextJSApps)
			if err != nil {
				fmt.Printf("%v\n", err)
				os.Exit(1)
			}
		}
	}

	projectPath := filepath.Join(projectInput.OutputFolder, projectInput.Name)

	if err := filesystem.CreateProjectStructure(projectPath, projectInput.SelectedFolders); err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

	if slices.Contains(projectInput.SelectedFolders, constants.FolderApp) && len(selectedApps) > 0 {
		nextjsGen := &generator.NextJSGenerator{
			Version:    versions.Frontend.NextJS.Version,
			ProjectDir: projectPath,
			AppFolder:  constants.FolderApp,
			Apps:       selectedApps,
		}

		if err := nextjsGen.Generate(projectInput.Name); err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}
	}

	output.PrintSuccess(projectInput.Name, projectPath)
}

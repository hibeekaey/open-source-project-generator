package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"

	"github.com/cuesoftinc/open-source-project-generator/internal/generator"
	"github.com/cuesoftinc/open-source-project-generator/pkg/config"
	"github.com/cuesoftinc/open-source-project-generator/pkg/constants"
	"github.com/cuesoftinc/open-source-project-generator/pkg/filesystem"
	"github.com/cuesoftinc/open-source-project-generator/pkg/input"
	"github.com/cuesoftinc/open-source-project-generator/pkg/output"
)

func main() {
	versions, err := config.LoadVersions(constants.VersionsConfigPath)
	if err != nil {
		fmt.Printf("%v\n", output.NewError("error loading versions: %v", err))
		os.Exit(1)
	}

	reader := bufio.NewReader(os.Stdin)
	projectInput, err := input.ReadProjectInput(reader, constants.DefaultOutputFolder)
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

	projectPath := filepath.Join(projectInput.OutputFolder, projectInput.Name)

	if err := filesystem.CreateProjectStructure(projectPath, constants.ProjectFolders); err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

	nextjsGen := &generator.NextJSGenerator{
		Version:    versions.Frontend.NextJS.Version,
		ProjectDir: projectPath,
		AppFolder:  constants.FolderApp,
		Apps:       []string{"main", "admin", "home"},
	}

	if err := nextjsGen.Generate(projectInput.Name); err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

	output.PrintSuccess(projectInput.Name, projectPath)
}

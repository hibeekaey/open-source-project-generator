package generator

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/cuesoftinc/open-source-project-generator/pkg/output"
)

type NextJSGenerator struct {
	Version    string
	ProjectDir string
	AppFolder  string
	Apps       []string
}

func (g *NextJSGenerator) Generate(projectName string) error {
	appPath := filepath.Join(g.ProjectDir, g.AppFolder)

	for _, app := range g.Apps {
		appName := projectName + "-" + app

		cmd := exec.Command("npx", fmt.Sprintf("create-next-app@%s", g.Version), appName,
			"--typescript",
			"--tailwind",
			"--app",
			"--disable-git",
			"--eslint",
			"--turbopack",
			"--src-dir",
			"--import-alias", "@/*",
			"--react-compiler",
			"--skip-install")
		cmd.Dir = appPath
		cmd.Stderr = os.Stderr

		spinner := output.NewSpinner(fmt.Sprintf("Setting up %s app...", app))
		spinner.Start()

		err := cmd.Run()
		spinner.Stop()

		if err != nil {
			return output.NewError("error creating %s app: %v", app, err)
		}

		oldPath := filepath.Join(appPath, appName)
		newPath := filepath.Join(appPath, app)
		if err := os.Rename(oldPath, newPath); err != nil {
			return output.NewError("error renaming %s folder: %v", app, err)
		}

		fmt.Printf(output.ColorGreen+"✔"+output.ColorReset+" Done setting up %s app\n", app)
	}

	return nil
}

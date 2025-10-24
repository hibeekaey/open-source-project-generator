package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

const (
	colorReset  = "\033[0m"
	colorGreen  = "\033[32m"
	colorCyan   = "\033[36m"
	colorRed    = "\033[31m"
	colorYellow = "\033[33m"

	defaultOutputFolder = "output/generated"
	configPath          = "configs/versions.yaml"
)

type Versions struct {
	Frontend struct {
		NextJS struct {
			Version string `yaml:"version"`
		} `yaml:"nextjs"`
	} `yaml:"frontend"`
}

func loadVersions() (*Versions, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var versions Versions
	if err := yaml.Unmarshal(data, &versions); err != nil {
		return nil, err
	}

	return &versions, nil
}

func main() {
	versions, err := loadVersions()
	if err != nil {
		fmt.Printf(colorRed+"Error loading versions: %v\n"+colorReset, err)
		os.Exit(1)
	}

	reader := bufio.NewReader(os.Stdin)

	fmt.Print(colorCyan + "Project name: " + colorReset)
	projectName, _ := reader.ReadString('\n')
	projectName = strings.TrimSpace(projectName)

	if projectName == "" {
		fmt.Println(colorRed + "Error: Project name cannot be empty" + colorReset)
		os.Exit(1)
	}

	validName := regexp.MustCompile(`^[a-z0-9-]+$`)
	if !validName.MatchString(projectName) {
		fmt.Println(colorRed + "Error: Project name must contain only lowercase letters, numbers, and hyphens" + colorReset)
		os.Exit(1)
	}

	fmt.Print(colorCyan + "Output folder (default: " + defaultOutputFolder + "): " + colorReset)
	outputFolder, _ := reader.ReadString('\n')
	outputFolder = strings.TrimSpace(outputFolder)

	if outputFolder == "" {
		outputFolder = defaultOutputFolder
	}

	projectPath := filepath.Join(outputFolder, projectName)

	if _, err := os.Stat(projectPath); err == nil {
		if err := os.RemoveAll(projectPath); err != nil {
			fmt.Printf(colorRed+"Error removing existing project: %v\n"+colorReset, err)
			os.Exit(1)
		}
	}

	folders := []string{"App", "CommonServer", "Mobile", "Deploy", "Docs", "Scripts", ".github"}
	for _, folder := range folders {
		folderPath := filepath.Join(projectPath, folder)
		if err := os.MkdirAll(folderPath, 0755); err != nil {
			fmt.Printf(colorRed+"Error creating folder %s: %v\n"+colorReset, folder, err)
			os.Exit(1)
		}
	}

	appPath := filepath.Join(projectPath, "App")
	apps := []string{"main", "admin", "home"}

	for _, app := range apps {
		appName := projectName + "-" + app

		cmd := exec.Command("npx", fmt.Sprintf("create-next-app@%s", versions.Frontend.NextJS.Version), appName,
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

		done := make(chan bool)
		go func() {
			spinner := []rune{'⠙', '⠹', '⠸', '⠼', '⠴', '⠦', '⠧', '⠇', '⠏'}
			i := 0
			fmt.Print("\n")
			for {
				select {
				case <-done:
					fmt.Print("\r\033[K")
					return
				default:
					fmt.Printf("\r"+colorCyan+"%c Setting up %s app..."+colorReset, spinner[i%len(spinner)], app)
					i++
					time.Sleep(80 * time.Millisecond)
				}
			}
		}()

		err := cmd.Run()
		done <- true
		time.Sleep(100 * time.Millisecond)

		if err != nil {
			fmt.Printf(colorRed+"Error creating %s app: %v\n"+colorReset, app, err)
			os.Exit(1)
		}

		oldPath := filepath.Join(appPath, appName)
		newPath := filepath.Join(appPath, app)
		if err := os.Rename(oldPath, newPath); err != nil {
			fmt.Printf(colorRed+"Error renaming %s folder: %v\n"+colorReset, app, err)
			os.Exit(1)
		}

		fmt.Printf(colorGreen+"✔"+colorReset+" Done setting up %s app\n", app)
	}

	fmt.Printf("\n"+colorGreen+"✔"+colorReset+" Created project "+colorCyan+"'%s'"+colorReset+" at "+colorYellow+"%s"+colorReset+"\n", projectName, projectPath)
}

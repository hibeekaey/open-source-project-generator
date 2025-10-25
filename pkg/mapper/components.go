package mapper

import "github.com/cuesoftinc/open-source-project-generator/pkg/constants"

func ComponentToFolder(component string) string {
	switch component {
	case "frontend":
		return constants.ComponentFrontend
	case "backend":
		return constants.ComponentBackend
	case "mobile":
		return constants.ComponentMobile
	case "deploy":
		return constants.ComponentDeploy
	case "docs":
		return constants.ComponentDocs
	case "scripts":
		return constants.ComponentScripts
	case "github":
		return constants.ComponentGithub
	default:
		return ""
	}
}

var ComponentOptions = []string{
	"frontend",
	"backend",
	"mobile",
	"deploy",
	"docs",
	"scripts",
	"github",
}

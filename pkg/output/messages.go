package output

import "fmt"

func PrintSuccess(projectName, projectPath string) {
	fmt.Printf("\n"+ColorGreen+"✔"+ColorReset+" Created project "+ColorCyan+"'%s'"+ColorReset+" at "+ColorYellow+"%s"+ColorReset+"\n", projectName, projectPath)
}

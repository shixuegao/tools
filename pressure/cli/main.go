package main

import (
	"umx/tools/pressure/cli/cmd"
)

func main() {
	err := cmd.RootCmd.Execute()
	if nil != err {
		return
	}
}

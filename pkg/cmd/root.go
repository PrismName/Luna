package cmd

import (
	"github.com/abiosoft/ishell"
	"github.com/fatih/color"
	"github.com/seaung/luna/pkg/utils"
)

var rootCmd = ishell.New()

func RunShell() {
	color.Blue(utils.ShowBanner())
	color.Red(utils.Author, utils.Version)

	rootCmd.SetPrompt("Luna > ")
	rootCmd.Run()
}

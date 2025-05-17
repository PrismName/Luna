package main

import (
	"github.com/seaung/Luna/internal/cli"
)

func main() {
	// 创建并运行交互式shell
	shell := cli.NewShell()
	shell.Run()
}

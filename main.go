package main

import (
	"os"

	"github.com/mkideal/cbuild/commands"
	"github.com/mkideal/cli"
	"github.com/mkideal/log"
)

func main() {
	defer log.Uninit(log.InitColoredConsole(log.LvWARN))
	err := cli.Root(commands.Root()).Run(os.Args[1:])
	log.If(err != nil).Error("Error: %v", err)
}

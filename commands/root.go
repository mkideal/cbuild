package commands

import (
	"github.com/mkideal/cbuild/commands/internal"
	"github.com/mkideal/cbuild/etc"
	"github.com/mkideal/cli"
	"github.com/mkideal/log"
)

func Root() *cli.Command { return root }

type rootT struct {
	cli.Helper
	etc.BuildEnv
	etc.Config
}

var root = &cli.Command{
	Name: "cbuild",
	Desc: "c/c++ program builder",
	Argv: func() interface{} {
		argv := new(rootT)
		argv.Config.Init()
		return argv
	},

	Fn: func(ctx *cli.Context) error {
		argv := ctx.Argv().(*rootT)
		internal.SetLogLevel(argv.Config.LogLevel)
		log.WithJSON(argv).Trace("argv")
		return internal.CreateMakefile(argv.Config, argv.BuildEnv)
	},
}

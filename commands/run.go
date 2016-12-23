package commands

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/mkideal/cbuild/commands/internal"
	"github.com/mkideal/cbuild/etc"
	"github.com/mkideal/cli"
	"github.com/mkideal/log"
)

func Run() *cli.Command { return run }

type runT struct {
	cli.Helper
	etc.BuildEnv
	etc.Config
	Args []string `cli:"t" usage:"arguments for built program"`
}

var run = &cli.Command{
	Name: "run",
	Desc: "immediately run program",
	Argv: func() interface{} { return new(runT) },

	Fn: func(ctx *cli.Context) error {
		argv := ctx.Argv().(*runT)
		if err := argv.Config.Load(ctx); err != nil {
			return err
		}
		internal.SetLogLevel(argv.Config.LogLevel)
		log.WithJSON(argv).Trace("argv")
		argv.BuildEnv.Stdout = ioutil.Discard
		makefile, err := internal.CreateMakefile(argv.Config, argv.BuildEnv)
		if makefile == nil {
			return err
		}
		defer os.RemoveAll(makefile.BuildDir)
		defer os.Remove(makefile.Target)
		bin := makefile.Target
		if !filepath.IsAbs(bin) {
			bin = "./" + bin
		}
		cmd := exec.Command(bin, argv.Args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	},
}

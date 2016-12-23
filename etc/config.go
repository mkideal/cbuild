package etc

import (
	"encoding/json"
	"io"
	"os"

	"github.com/mkideal/cli"
	"github.com/mkideal/log/logger"
)

type Config struct {
	Filename       string       `cli:"f" usage:"cbuild config filename" dft:"cbuild.json" json:"-"`
	LogLevel       logger.Level `cli:"v,log-level" usage:"verbose level(log level): trace,debug,info,warn,error or fatal" dft:"warn" json:"verbose"`
	C              string       `cli:"c" usage:"specific c compiler" dft:"gcc" json:"c"`
	Cpp            string       `cli:"c++,cpp" usage:"specific c++ compiler" dft:"g++" json:"c++"`
	Make           string       `cli:"make-cmd" usage:"make command" dft:"make" json:"make"`
	ObjectsDir     string       `cli:"d,objects-dir" usage:"built objects directory" dft:"objects" json:"objects_dir"`
	IncludeDirs    []string     `cli:"I,include-dir" usage:"included directories" json:"include_dirs"`
	LibDirs        []string     `cli:"L,library-dir" usage:"library directories" json:"lib_dirs"`
	Libs           []string     `cli:"l,lib" usage:"libs which should be linked" json:"libs"`
	ExcludeSources []string     `cli:"exclude-source" usage:"excluded source directories or/and files" json:"exclude_dirs"`
}

type BuildEnv struct {
	Release bool      `cli:"r,release" usage:"build release version?" dft:"false"`
	Stdout  io.Writer `cli:"-"`
	Stderr  io.Writer `cli:"-"`
}

func (c *Config) Load(ctx *cli.Context) error {
	if len(c.IncludeDirs) == 0 {
		c.IncludeDirs = append(c.IncludeDirs, ".")
	}
	file, err := os.Open(c.Filename)
	if err != nil {
		if !ctx.IsSet("-f") {
			return nil
		}
		return err
	}
	defer file.Close()
	return json.NewDecoder(file).Decode(c)
}

func (c Config) IsExcluded(path string) bool {
	for _, f := range c.ExcludeSources {
		if f == path {
			return true
		}
	}
	return false
}

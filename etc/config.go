package etc

import (
	"github.com/mkideal/log/logger"
)

type Config struct {
	Self           *Config      `cli:"f" usage:"cbuild config filename" dft:"cbuild.json" json:"-" parser:"jsonfile"`
	LogLevel       logger.Level `cli:"v,log-level" usage:"verbose level(log level): trace,debug,info,warn,error or fatal" dft:"warn" json:"verbose"`
	C              string       `cli:"c" usage:"specific c compiler" dft:"gcc" json:"c"`
	Cpp            string       `cli:"c++,cpp" usage:"specific c++ compiler" dft:"g++" json:"c++"`
	Make           string       `cli:"make-cmd" usage:"make command" dft:"make" json:"make"`
	TargetsDir     string       `cli:"d,targets-dir" usage:"targets building directory" dft:"targets" json:"targets_dir"`
	IncludeDirs    []string     `cli:"I,include-dir" usage:"included directories" json:"include_dirs"`
	LibDirs        []string     `cli:"L,library-dir" usage:"library directories" json:"lib_dirs"`
	Libs           []string     `cli:"l,lib" usage:"libs which should be linked" json:"libs"`
	ExcludeSources []string     `cli:"exclude-source" usage:"excluded source directories or/and files" json:"exclude_dirs"`
}

type BuildEnv struct {
	Release bool `cli:"r,release" usage:"build release version?" dft:"false"`
}

func (c *Config) Init() {
	c.Self = c
	c.IncludeDirs = append(c.IncludeDirs, ".")
}

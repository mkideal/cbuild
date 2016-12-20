package internal

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"

	"github.com/mkideal/cbuild/etc"
	"github.com/mkideal/log"
	"github.com/mkideal/log/logger"
	"github.com/mkideal/pkg/errors"
)

type Makefile struct {
	etc.Config
	SourceFiles []string
	Target      string
	ReleaseTag  string
	Includes    string
	LibraryDirs string
	Libs        string
}

const T_makefile = `build:{{range $source := .SourceFiles}} {{$source}}.o{{end}}
	{{.Cpp}}{{range $source := .SourceFiles}} {{$source}}.o{{end}} {{.ReleaseTag}} {{.Includes}} {{.LibraryDirs}} {{.Libs}} -o {{.Target}}

{{range $source := .SourceFiles}}
{{$source}}.o: $source
	{{.Cpp}} -c {{$source}} {{.ReleaseTag}} {{.Includes}} -o {{$source}}.o
{{end}}
`

func SetLogLevel(level logger.Level) {
	log.SetLevel(level)
	if !level.MoreVerboseThan(log.LvINFO) {
		log.NoHeader()
	}
}

func GetSourceFiles(projectDir string, c etc.Config) ([]os.FileInfo, error) {
	var files []os.FileInfo
	filepath.Walk(projectDir, func(path string, info os.FileInfo, err error) error {
		return nil
	})
	return files, nil
}

func CreateMakefile(c etc.Config, env etc.BuildEnv) error {
	makefile := Makefile{Config: c}
	return GenMakefile(makefile)
}

func GenMakefile(makefile Makefile) error {
	t := template.New("makefile")
	t, err := t.Parse(T_makefile)
	if err != nil {
		return errors.Throw(err.Error())
	}
	var buf bytes.Buffer
	err = t.Execute(&buf, makefile)
	if err != nil {
		return errors.Throw(err.Error())
	}
	err = ioutil.WriteFile("Makefile", buf.Bytes(), 0666)
	if err != nil {
		return errors.Throw(err.Error())
	}
	return nil
}

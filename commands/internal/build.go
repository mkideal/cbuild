package internal

import (
	"bytes"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/mkideal/cbuild/etc"
	"github.com/mkideal/log"
	"github.com/mkideal/log/logger"
	"github.com/mkideal/pkg/errors"
)

var sourceFileSuffixList = [...]string{
	"cpp",
	"cxx",
	"hpp",
	"hxx",
	"cc",
	"c",
}

func hasSourceSuffix(name string) bool {
	splitted := strings.Split(name, ".")
	if len(splitted) == 0 {
		return false
	}
	suffix := strings.ToLower(splitted[len(splitted)-1])
	for _, s := range sourceFileSuffixList {
		if suffix == s {
			return true
		}
	}
	return false
}

type Makefile struct {
	Filename string
	etc.Config
	SourceFiles []string
	BuildDir    string
	Target      string
	ReleaseTag  string
	Includes    string
	LibraryDirs string
	Libs        string
}

func newMakefile(c etc.Config) *Makefile {
	makefile := &Makefile{
		Config:   c,
		Filename: "Makefile",
	}
	makefile.Includes = strings.Join(c.IncludeDirs, " -I")
	if len(c.IncludeDirs) > 0 {
		makefile.Includes = "-I" + makefile.Includes
	}
	makefile.LibraryDirs = strings.Join(c.LibDirs, "-L")
	if len(c.LibDirs) > 0 {
		makefile.LibraryDirs = "-L" + makefile.LibraryDirs
	}
	makefile.Libs = strings.Join(c.Libs, "-l")
	if len(c.Libs) > 0 {
		makefile.Libs = "-l" + makefile.Libs
	}
	return makefile
}

const T_makefile = `build:{{with $m := .}}{{range $source := .SourceFiles}} {{$m.BuildDir}}{{$source}}.o{{end}}{{end}}
	{{.Cpp}}{{with $m := .}}{{range $source := .SourceFiles}} {{$m.BuildDir}}{{$source}}.o{{end}}{{end}} {{.ReleaseTag}} {{.Includes}} {{.LibraryDirs}} {{.Libs}} -o {{.Target}}

	{{with $m := .}}{{range $source := .SourceFiles}}
{{$m.BuildDir}}{{$source}}.o: {{$source}}
	{{$m.Cpp}} -c {{$source}} {{$m.ReleaseTag}} {{$m.Includes}} -o {{$m.BuildDir}}{{$source}}.o
{{end}}{{end}}
`

func SetLogLevel(level logger.Level) {
	log.SetLevel(level)
	if !level.MoreVerboseThan(log.LvINFO) {
		log.NoHeader()
	}
}

func GetSourceFiles(projectDir string, c etc.Config) []string {
	var files []string
	filepath.Walk(projectDir, func(path string, info os.FileInfo, err error) error {
		if info == nil {
			return err
		}
		log.Debug("path=%s", path)
		if info.IsDir() && c.IsExcluded(path) {
			return filepath.SkipDir
		}
		if !info.IsDir() && hasSourceSuffix(info.Name()) {
			files = append(files, path)
		}
		return nil
	})
	return files
}

func CreateMakefile(c etc.Config, env etc.BuildEnv) (*Makefile, error) {
	if env.Stdout == nil {
		env.Stdout = os.Stdout
	}
	if env.Stderr == nil {
		env.Stderr = os.Stderr
	}
	projectDir := "."
	abs, err := filepath.Abs(projectDir)
	if err != nil {
		return nil, err
	}
	dir := filepath.Base(abs)
	makefile := newMakefile(c)
	makefile.SourceFiles = GetSourceFiles(projectDir, c)
	if len(makefile.SourceFiles) == 0 {
		log.Warn("no source file")
		return nil, nil
	}
	makefile.Target = dir
	if env.Release {
		makefile.BuildDir = filepath.Join(c.ObjectsDir, "release")
	} else {
		makefile.BuildDir = filepath.Join(c.ObjectsDir, "debug")
	}
	if err := InitBuildDir(makefile); err != nil {
		os.RemoveAll(makefile.BuildDir)
		return nil, err
	}
	makefile.BuildDir += "/"
	err = GenMakefile(makefile)
	if err != nil {
		return nil, err
	}
	defer os.Remove(makefile.Filename)
	cmd := exec.Command(makefile.Config.Make)
	cmd.Stdout = env.Stdout
	cmd.Stderr = env.Stderr
	if err := cmd.Run(); err != nil {
		return nil, err
	}
	return makefile, nil
}

func GenMakefile(makefile *Makefile) error {
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
	err = ioutil.WriteFile(makefile.Filename, buf.Bytes(), 0666)
	if err != nil {
		return errors.Throw(err.Error())
	}
	return nil
}

func InitBuildDir(makefile *Makefile) error {
	for _, name := range makefile.SourceFiles {
		dir, _ := filepath.Split(name)
		dir = filepath.Join(makefile.BuildDir, dir)
		log.Debug("dir=%s", dir)
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}

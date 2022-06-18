// В проекте представлена утилита для приведения кода проекта на go от snake_case к camelCase.
package main

import (
	"CamelCase/camelCase"
	"CamelCase/fileEditRuner"
	"CamelCase/fileEditor"
	"flag"
	"regexp"
)

var (
	path   = flag.String("path", ".", "folder with code files")
	recurs = flag.Bool("recurs", false, "execute recursively for nested directories")
	regExp = flag.String("regExp", `.*\.go$`, "regular expression for filenames")
)

func main() {
	flag.Parse()

	// Конфигурация запуска редактора.
	r := fileEditRuner.New()
	r.Path = *path
	r.Recurs = *recurs
	r.Reg = regexp.MustCompile(*regExp)

	// Конфигурация редактора
	fe := fileEditor.New()
	fe.Ed = camelCase.New()

	r.FileEditRun(fe)
}

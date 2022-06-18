// fes содержит реализацию механизма параллельного запуска произвольного FileEditer.
package fileEditRuner

import (
	"io/fs"
	"log"
	"os"
	"regexp"
	"strings"
	"sync"
)

// FileEditer - контракт с редактором файла
type FileEditer interface {
	FileEdit(fullFileName string) error
}

func New() *FileEditRuner {
	return &FileEditRuner{}
}

// FileEditStarter хранит параметры запуска FileEditer.
type FileEditRuner struct {
	Path   string         // корневая папка с кодом
	Recurs bool           // флаг рекурсивного запуска для подпапок
	Reg    *regexp.Regexp // регулярное выражение имен редактируемых файлов
}

// FileEditStart запускает переданный FileEditer для
// каждого файла, в соответствии с параметрами fes.
func (fes *FileEditRuner) FileEditRun(fe FileEditer) {
	fes.dirWalker(fes.Path, fe)
}

// dirWalker запускает FileEdit, переданного FileEditer для всех файлов в директрории
// path. Если fes.Recurs==true, то FileEdit запускается рекурсивно для всех
// поддиректорий. Ошибки открытия папок и ошибки, возвращаемые FileEditer,
// записываются в лог. Все запуски параллельны.
func (fes *FileEditRuner) dirWalker(path string, fe FileEditer) {
	list, err := os.ReadDir(path)
	if err != nil {
		log.Println(err)
	}

	var wg sync.WaitGroup

	fullName := func(el fs.DirEntry) string {
		return strings.Join([]string{path, el.Name()}, "/")
	}

	for _, el := range list {
		if el.IsDir() {
			if *&fes.Recurs {
				wg.Add(1)
				go func(fullName string) {
					if fes.dirWalker(fullName, fe); err != nil {
						log.Println(err)
					}
					wg.Done()
				}(fullName(el))
			}
		} else {
			if fes.Reg.MatchString(el.Name()) {
				wg.Add(1)
				go func(fullName string) {
					if err := fe.FileEdit(fullName); err != nil {
						log.Println(err)
					}
					wg.Done()
				}(fullName(el))
			}
		}
	}
	wg.Wait()
}

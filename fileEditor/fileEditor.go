// Пакет реализует редактор файла FileEditor, предназначенный для
// редактирования текстового файла с помощью произвольного редактора,
// удовлетворяющего интерфейсу Editer.
package fileEditor

import (
	"bufio"
	"log"
	"os"
)

// Editor является контрактом редактора, который принемает руны по
// каналу input и возвращает измененный текст в канал output.
type Editer interface {
	Edit(input, output chan rune)
}

// FileEditor предназначен для редактирования текстового файла с
// помощью произвольного редактора, удовлетворяющего интерфейсу Editer.
type FileEditor struct {
	Ed      Editer
	BufSize int // размер буфера каналов input и output передаваемых редактору.
}

func New() *FileEditor {
	return &FileEditor{}
}

// FileEdit редактирует текстовый файл с помощью произвольного
// редактора, удовлетворяющего интерфейсу Editer.
func (fe *FileEditor) FileEdit(fullFileName string) error {
	input := make(chan rune, fe.BufSize)
	output := make(chan rune, fe.BufSize)

	if err := os.Rename(fullFileName, fullFileName+".tmp"); err != nil {
		return err
	}

	outF, err := os.Create(fullFileName)
	if err != nil {
		os.Rename(fullFileName+".tmp", fullFileName)
		return err
	}
	defer outF.Close()

	inF, err := os.Open(fullFileName + ".tmp")
	if err != nil {
		os.Remove(fullFileName)
		os.Rename(fullFileName+".tmp", fullFileName)
		return err
	}

	sc := bufio.NewScanner(inF)
	sc.Split(bufio.ScanRunes)

	go fe.Ed.Edit(input, output)

	go func() {
		for sc.Scan() {
			input <- []rune(sc.Text())[0]
		}
		close(input)
	}()

	for ch := range output {
		if _, err := outF.WriteString(string([]rune{ch})); err != nil {
			return err
		}
	}
	inF.Close()

	if err := os.Remove(fullFileName + ".tmp"); err != nil {
		return err
	}

	log.Printf("%q camelCaseFormatted\n", fullFileName)
	return nil
}

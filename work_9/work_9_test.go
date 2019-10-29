package main

import (
	// "fmt"
	// "os"
	// log "github.com/sirupsen/logrus"

	"bufio"
	"bytes"
	"crypto/rand"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/cheggaaa/pb"
	"go.uber.org/zap"
	// "github.com/urfave/cli"
	// "io/ioutil"
	// "regexp"
	// "strconv"
	// "strings"
	// "time"
)

func TestSetLogger(t *testing.T) {
	var err error
	if log, err = zap.NewDevelopment(); err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
}

func createFileForTest(filePath string, fileSize int64) error {

	f, err := os.Create(filePath)
	defer f.Close()
	if err != nil {
		log.Error(`Ошибка создания тестового файла`, zap.Error(err))
		return err
	}

	fb := bufio.NewWriter(f)
	defer fb.Flush()
	buf := make([]byte, defaultIBS)
	bar := pb.StartNew(int(fileSize))
	reader := bar.NewProxyReader(rand.Reader)
	for i := fileSize; i > 0; i -= defaultIBS {
		if _, err = reader.Read(buf); err != nil {
			fmt.Printf("error occurred during random: %s\n", err)
			break
		}
		bR := bytes.NewReader(buf)
		if _, err = io.Copy(fb, bR); err != nil {
			fmt.Printf("failed during copy: %s\n", err)
			break
		}
	}
	bar.Finish()
	return nil
}

func PrepareForTest(fileSizeForTest int64) {

	if f, err := os.Open(defaultInFile); err == nil {
		if fInfo, err2 := f.Stat(); err2 == nil {
			if fInfo.Size() == fileSizeForTest {
				// "Файл размером уже создан
				return
			}
		}
	}
	fmt.Printf("Создаем файл для тестов, из которого будем читать %s. Размер %f Мбайт\n", defaultInFile, float64(fileSizeForTest/1024/1024))
	if err := createFileForTest(defaultInFile, fileSizeForTest); err != nil {
		log.Error(`Ошибка при создании файла для тестирования`, zap.String(`File`, defaultInFile))
		return
	}
}

const inFileSize = defaultIBS * 1024 * 1024

func TestWork9Prepare(t *testing.T) {
	PrepareForTest(inFileSize)

}

func TestWork9(t *testing.T) {
	// t.Skip()

	CopyFile(defaultInFile, `/tmp/out_file_test_1`, 0, defaultIBS, defaultOBS)
	CopyFile(defaultInFile, `/tmp/out_file_test_2`, inFileSize/2, defaultIBS, defaultOBS)
	CopyFile(defaultInFile, `/tmp/out_file_test_3`, inFileSize/2+10240, defaultIBS, defaultOBS)

}

func TestFeedBack(t *testing.T) {
	CopyFile(`Makefile`, `file.txt`, 0, 512, 512)
}

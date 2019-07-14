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
	log "github.com/sirupsen/logrus"
	// "github.com/urfave/cli"
	// "io/ioutil"
	// "regexp"
	// "strconv"
	// "strings"
	// "time"
)

func createFileForTest(filePath string, fileSize int64) error {

	f, err := os.Create(filePath)
	defer f.Close()
	if err != nil {
		log.Error(err)
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
				fmt.Printf("Файл %s размером %d байт уже создан\n", defaultInFile, fileSizeForTest)
				return
			}
		}
	}
	fmt.Printf("Создаем файл для тестов, из которого будем читать %s\n", defaultInFile)
	if err := createFileForTest(defaultInFile, fileSizeForTest); err != nil {
		log.Errorf(`Ошибка при создании файла для тестирования %s`, defaultInFile)
		return
	}
}

func TestWork11(t *testing.T) {

	// fmt.Printf("The date is %s\n", out)
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
	var inFileSize int64
	inFileSize = defaultIBS * 1000000
	PrepareForTest(inFileSize)
	CopyFile(defaultInFile, `/tmp/out_file_test_1`, 0, defaultIBS, defaultOBS)
	CopyFile(defaultInFile, `/tmp/out_file_test_2`, inFileSize/2, defaultIBS, defaultOBS)
	CopyFile(defaultInFile, `/tmp/out_file_test_3`, inFileSize/2+10240, defaultIBS, defaultOBS)

}

package main

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/schollz/progressbar"
	"go.uber.org/zap"
	// "github.com/spf13/viper"
)

// CopyFile копирует данные из файла по пути inPath в файл по пути outPath
// Используя смещение в исходном файле offset
// читая за раз ibs и пишет за раз obs
func CopyFile(inPath, outPath string, offset int64, ibs, obs int) error {
	var (
		in, out *os.File
		err     error
	)

	if in, err = os.OpenFile(inPath, os.O_RDONLY, 0600); err != nil {
		fmt.Fprintf(os.Stderr, "Не удалось открыть файл %s\n", inPath)
		return err
	}
	defer in.Close()

	var fInfo os.FileInfo
	if fInfo, err = in.Stat(); err != nil {
		fmt.Fprintf(os.Stderr, "Не удалось пуолучить информацию о файле %s\n", inPath)
		return err
	}
	var inFileSize int64
	inFileSize = fInfo.Size()
	if _, err = in.Seek(offset, 0); err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка смещения по файлу %s размером %d на %d байт от начала\n", inPath, inFileSize, offset)
		return err
	}

	sizeForCopy := inFileSize - offset
	fmt.Printf("\t\tБудем копировать %d байт из файла %s размером %d\n", sizeForCopy, inPath, inFileSize)

	if out, err = os.OpenFile(outPath, os.O_RDWR|os.O_CREATE, 0600); err != nil {
		fmt.Fprintf(os.Stderr, "Не удалось открыть файл %s\n", outPath)
		return err
	}
	defer out.Close()
	// bar := progressbar.NewOptions(
	// 	int(sizeForCopy),
	// 	progressbar.OptionSetBytes(int(sizeForCopy)),
	// )
	writer := bufio.NewWriterSize(out, obs)
	if err = Copy(in, make([]byte, ibs), writer, sizeForCopy); err != nil {
		return err
	}
	if err = writer.Flush(); err != nil {
		return err
	}
	fmt.Printf("\t\tФайл %s закончил копирование\n", outPath)
	return nil
}

//Copy производим копирование используя полученные буферы
func Copy(reader io.ReadWriteSeeker, readBuf []byte, writer *bufio.Writer, sizeForCopy int64) error {
	var err error
	var readBytes int
	// reader := io.LimitReader(in, sizeForCopy)
	// bar := progressbar.New(int(sizeForCopy / int64(len(readBuf))))

	bar := progressbar.NewOptions(int(sizeForCopy/int64(len(readBuf))),
		progressbar.OptionSetPredictTime(false),
	)
	for {
		readBytes, err = reader.Read(readBuf)
		if err == io.EOF {
			if _, err = writer.Write(readBuf[0:readBytes]); err != nil {
				return err
			}
			return nil
		}
		if err != nil {
			return err
		}

		if _, err = writer.Write(readBuf[0:readBytes]); err != nil {
			return err
		}
		if err = bar.Add(1); err != nil {
			log.Warn("Ошибка progressbar", zap.Error(err))
		}
	}
}

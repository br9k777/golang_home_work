package main

import (
	"bufio"
	"fmt"
	"io"
	"os"

	pb "github.com/cheggaaa/pb/v3"
	log "github.com/sirupsen/logrus"
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

	reader := io.LimitReader(in, sizeForCopy)
	bar := pb.Full.Start64(sizeForCopy)
	barReader := bar.NewProxyReader(reader)

	readBuf := make([]byte, ibs)
	var read int
	outBuf := bufio.NewWriterSize(out, obs)
	for {

		read, err = barReader.Read(readBuf)
		if err == io.EOF {
			outBuf.Write(readBuf[0:read])
			bar.Finish()
			readPosition, _ := in.Seek(0, 1)
			outPosition, _ := out.Seek(0, 1)
			fmt.Printf("Файл закончился при четнии со смещения %d с размером буфера %d\n", readPosition, ibs)
			fmt.Printf("При последенем чтении успешно прочитано %d байт\n", read)
			fmt.Printf("Записываемый файл закончил писатся при смещении %d с размером буфера %d\n", outPosition, obs)
			log.Debugf(`Cам буфер %#v`, readBuf)
			break
		}
		outBuf.Write(readBuf)
		if err != nil {
			readPosition, _ := in.Seek(0, 1)
			log.Errorf(`Ошибка при смещении %d с буфром %#v`, readPosition, readBuf)
			fmt.Printf("Failed to read: %v", err)
			return err
		}
	}
	bar.Finish()

	return nil
}

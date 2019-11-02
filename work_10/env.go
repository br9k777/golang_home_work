package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

// GetEnviromentsFromDir читает файлы из указанной директории
// и возвращает map[имя_файла]содержимое_файла.
func GetEnviromentsFromDir(pathToDir string) (map[string]string, error) {
	// чистим путь до файлов, в теории поможет при использовании на windows
	clearDirPath := filepath.Clean(pathToDir)
	if c, err := os.Stat(clearDirPath); os.IsNotExist(err) || !c.Mode().IsDir() {
		return nil, fmt.Errorf("Can't reach directory %s", clearDirPath)
	}
	var err error
	var files []os.FileInfo
	if files, err = ioutil.ReadDir(clearDirPath); err != nil {
		return nil, err
	}
	newEnv := make(map[string]string)
	var bytesInFile []byte
	for _, f := range files {
		if f.Mode().IsRegular() {
			if bytesInFile, err = ioutil.ReadFile(filepath.Clean(fmt.Sprintf("%s/%s/", clearDirPath, f.Name()))); err != nil {
				return newEnv, err
			}

			if len(bytesInFile) > 0 {
				// если последний байт = 10 LF (Line Feed) убираем его
				if bytesInFile[len(bytesInFile)-1] == 10 {
					newEnv[f.Name()] = string(bytesInFile[0 : len(bytesInFile)-1])
				} else {
					newEnv[f.Name()] = string(bytesInFile)
				}
			}
		}
	}
	return newEnv, nil
}

// RunProgragWirhEnviroments запустить программу  program
// используя переменные из map enviroments.
func RunProgragWirhEnviroments(enviroments map[string]string, program string) error {
	command := exec.Command(program)
	newEnivroments := make([]string, 0, len(enviroments))
	for envName, envValue := range enviroments {
		newEnivroments = append(newEnivroments, fmt.Sprintf("%s=%s", envName, envValue))
	}
	command.Env = newEnivroments
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	var err error
	if err = command.Run(); err != nil {
		return err
	}
	return nil
}

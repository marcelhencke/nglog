package main

import (
	"bufio"
	"errors"
	"io"
	"nglog/config"
	debug "nglog/debugPrint"
	"nglog/nginxErrorLogFormatter"
	"os"
)

func LoadData() error {
	if isInputFromPipe() {
		debug.Println("Reading data from pipe")
		return processData(os.Stdin, os.Stdout)
	} else {
		debug.Println("Loading data from file: " + config.Flags.FilePath)
		file, e := getFile()
		if e != nil {
			return e
		}
		defer file.Close()
		return processData(file, os.Stdout)
	}
}

func processData(r io.Reader, w io.Writer) error {
	formatter := nginxErrorLogFormatter.New(w)

	scanner := bufio.NewScanner(bufio.NewReader(r))
	rowCounter := 0
	for scanner.Scan() {
		rowCounter++
		e := formatter.ReadBufferLine(scanner.Text())
		if e != nil {
			return e
		}
		if config.Flags.DebugMode && config.Flags.DebugMaxLines > 0 && rowCounter >= config.Flags.DebugMaxLines {
			return nil
		}
	}
	return nil
}

func isInputFromPipe() bool {
	fi, _ := os.Stdin.Stat()
	return fi.Mode()&os.ModeCharDevice == 0
}

func getFile() (*os.File, error) {
	if config.Flags.FilePath == "" {
		return nil, errors.New("please input a file path")
	}
	if !fileExists(config.Flags.FilePath) {
		return nil, errors.New("the file provided does not exist")
	}
	file, e := os.Open(config.Flags.FilePath)
	if e != nil {
		return nil, errors.New("unable to read the file " + config.Flags.FilePath + " : " + e.Error())
	}
	return file, nil
}

func fileExists(filepath string) bool {
	info, e := os.Stat(filepath)
	if os.IsNotExist(e) {
		return false
	}
	return !info.IsDir()
}

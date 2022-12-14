package io

import (
	"errors"
	"os"
	"path/filepath"
	"regexp"
	"time"
)

func ReadLastLine(filePath string, charactersToRead int) (line string, err error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return "", err
	}

	fileSize := fileInfo.Size()
	bufferSize := int64(charactersToRead)
	if fileSize < bufferSize {
		return "", errors.New("File is too small.")
	}

	offset, err := file.Seek(0, 2)
	if err != nil {
		return "", err
	}

	start := offset - bufferSize

	lineBuf := make([]byte, bufferSize)
	count, err := file.ReadAt(lineBuf, start)
	if err != nil || bufferSize != int64(count) {
		return "", err
	}

	return lastCompleteLine(lineBuf)
}

func lastCompleteLine(lineBuf []byte) (line string, err error) {
	endIndex := -1
	for i := len(lineBuf) - 1; i >= 0; i-- {
		value := lineBuf[i]
		if !isLineBreak(value) {
			endIndex = i + 1
			break
		}
	}
	startIndex := -1
	for i := endIndex - 1; i >= 0; i-- {
		value := lineBuf[i]
		if isLineBreak(value) {
			startIndex = i + 1
			break
		}
	}

	if startIndex < 0 || endIndex < 0 {
		return "", errors.New("Line is longer than expected!")
	}

	return string(lineBuf[startIndex:endIndex]), nil
}

func isLineBreak(char byte) bool {
	return char == 10 || char == 13
}

func DetermineFileByDate(directory string) (*File, error) {
	time := time.Now()
	fileName := time.Format("2006-01-02") + ".txt"
	filePath := filepath.Join(directory, fileName)
	_, err := os.Stat(filePath)
	if err != nil {
		return nil, err
	} else {
		return &File{FilePath: filePath, FileName: fileName}, nil
	}
}

func DetermineFileByOrder(directory string) (*File, error) {
	files, err := os.ReadDir(directory)
	if err != nil {
		return nil, err
	}

	fileNameFormat, _ := regexp.Compile("\\d{4}-\\d{2}-\\d{2}\\.txt")

	last := len(files) - 1
	for i := range files {
		file := files[last-i]
		fileName := file.Name()
		if file.Type().IsRegular() && fileNameFormat.MatchString(fileName) {
			filePath := filepath.Join(directory, fileName)
			return &File{FilePath: filePath, FileName: fileName}, nil
		}
	}

	return nil, errors.New("No regular file found in directory!")
}

package utils

import (
	"errors"
	"little-auth/config"
	"os"
	"regexp"
	"strings"
)

var (
	ErrDataTooLarge = errors.New("data is too large")
)

func WriteBytesToFile(dir string, filename string, data []byte) error {
	if int64(len(data)) > config.MAX_FILE_BYTE_LENGTH {
		return ErrDataTooLarge
	}
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return err
		}
	}

	if dir[len(dir)-1] != '/' {
		dir += "/"
	}
	file, err := os.Create(dir + filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		return err
	}
	return nil
}

func ReadBytesFromFile(dir string, filename string) ([]byte, error) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return nil, err
	}

	if dir[len(dir)-1] != '/' {
		dir += "/"
	}
	file, err := os.Open(dir + filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}

	fileSize := fileInfo.Size()
	if fileSize > config.MAX_FILE_BYTE_LENGTH {
		return nil, ErrDataTooLarge
	}

	data := make([]byte, fileSize)
	if _, err = file.Read(data); err != nil {
		return nil, err
	}
	return data, nil
}

func DeleteFile(dir string, filename string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return err
	}
	if dir[len(dir)-1] != '/' {
		dir += "/"
	}
	return os.Remove(dir + filename)
}

func ListFiles(dir string) ([]string, error) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return nil, err
	}
	if dir[len(dir)-1] != '/' {
		dir += "/"
	}
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var fileNames []string
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		fileNames = append(fileNames, file.Name())
	}
	return fileNames, nil
}

func NormalizeFileName(fileName string) string {
	// use regex to replace all non-alphanumeric characters with a dash
	fileName = strings.ToLower(fileName)
	re := regexp.MustCompile("[^a-z0-9_]+")
	fileName = re.ReplaceAllString(fileName, "-")
	fileName = strings.Trim(fileName, "-")
	return fileName
}

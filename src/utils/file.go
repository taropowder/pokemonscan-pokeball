package utils

import (
	"bytes"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

const CacheDir = "/tmp"

func GetResultTmp(plugin, target string) string {
	target = strings.Replace(target, ".", "_", -1)
	return filepath.Join(CacheDir, plugin+"_"+target)
}

func WriteFileFromBase64(fileDir, fileName, base64Text string) (string, error) {
	//pwd, _ := os.Getwd()
	//configDir := path.Join(pwd, fileDir)
	if _, err := os.Stat(fileDir); err != nil {
		if err := os.MkdirAll(fileDir, os.ModePerm); err != nil {
			return "", err
		}
	}

	writeFile := path.Join(fileDir, fileName)

	apiContext, err := DecodeB64(base64Text)
	if err != nil {
		log.Errorf("error writr error %v", err)
	}

	err = ioutil.WriteFile(writeFile, bytes.Trim([]byte(apiContext), "\x00"), 0644)
	if err != nil {
		log.Errorf("error writr error %v", err)
	}

	return writeFile, nil
}

func RemoveFileIfExist(filePath string) {
	if _, err := os.Stat(filePath); err != nil {
		return
	}
	err := os.Remove(filePath)
	if err != nil {
		log.Errorf("remove file %s error %v", filePath, err)
	}
}

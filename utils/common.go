package utils

import (
	"encoding/json"
	"github.com/subchen/go-log"
	"math/rand"
	"os"
	"strconv"
	"time"
)

// Jsonify to convert data into json without cumbersome error handling
func Jsonify(v interface{}) string {
	var j []byte
	var err error
	if j, err = json.Marshal(v); err != nil {
		log.Error(err)
	}
	return string(j)
}

// GetPodName used to get pod name when service run on top of Kubernetes
func GetPodName() string {
	str := ""
	exist := false
	if str, exist = os.LookupEnv("POD_NAME"); !exist {
		rand.Seed(time.Now().UnixNano())
		str = time.Now().Format("2006-01-02") + "-" + strconv.Itoa(rand.Int())
	}
	return str
}

// Mkdir used to make a dir without cumbersome error handling (default mode is 777)
func Mkdir(location string) {
	err := os.Mkdir(location, os.ModeDir)
	if err != nil {
		log.Error(err)
	}
}

// CreateFile used to create file without cumbersome error handling
func CreateFile(location string) *os.File {
	file, errFileCreate := os.Create(location)
	if errFileCreate != nil {
		log.Error("Create file error: ", errFileCreate)
		return nil
	}
	return file
}

func GetWorkingDir() string {
	workDir, err := os.Getwd()
	if err != nil {
		log.Error(err)
	}
	return workDir
}

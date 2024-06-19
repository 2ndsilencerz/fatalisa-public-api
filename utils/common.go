package utils

import (
	"encoding/json"
	utils2 "fatalisa-public-api/service/web/utils"
	"github.com/subchen/go-log"
	"os"
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
//func GetPodName() string {
//	str := ""
//	exist := false
//	if str, exist = os.LookupEnv("POD_NAME"); !exist {
//		rand.New(rand.NewSource(time.Now().UnixNano()))
//		str = time.Now().Format("2006-01-02") + "-" + strconv.Itoa(rand.Int())
//	}
//	return str
//}

func CheckAndCreateDir(path string) {
	if _, err := os.Stat(path); err != nil {
		log.Warn(err)
		workDir := GetWorkingDir()
		createDir := workDir + path
		Mkdir(createDir)
	}
}

// Mkdir used to make a dir without cumbersome error handling (default mode is 777)
func Mkdir(location string) {
	log.Info("Creating dir ", location)
	err := os.Mkdir(location, os.ModePerm)
	utils2.ErrorHandler(err)
}

// CreateFile used to create file without cumbersome error handling
func CreateFile(location string) *os.File {
	log.Info("Creating file ", location)
	file, errFileCreate := os.Create(location)
	if err, _ := utils2.ErrorHandler(errFileCreate); err {
		return nil
	}
	return file
}

func GetWorkingDir() string {
	workDir, err := os.Getwd()
	utils2.ErrorHandler(err)
	return workDir + "/"
}

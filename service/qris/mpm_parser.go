package qris

import (
	"fmt"
	"github.com/pieterclaerhout/go-log"
	"strconv"
	"sync"
)

var waitGroup = sync.WaitGroup{}
var totalGoRoutine = 0

// SEPARATOR only used as key in map for subId
const SEPARATOR = ""
const MaxIndex = 65

var contents map[string]string

func init() {
	contents = make(map[string]string)
}

func parseAsMap(rawData string) {
	parseMPM(rawData, "")
}

func parseMPM(rawData string, rootId string) {
	indexId := 0
	for {
		if len(rawData) == 0 {
			break
		}

		currentId := rawData[0:2]
		expectedId := fmt.Sprintf("%02d", indexId)
		if currentId == expectedId {
			data := getContentMpm(rawData)
			putContentMpm(rootId+SEPARATOR+expectedId, data)

			if rootId == "" && isRootIdHaveSubId(currentId) {
				totalGoRoutine++
				waitGroup.Add(1)
				go getSubContentMpm(data, currentId)
				waitGroup.Wait()
			}
			rawData = stripContent(rawData, len(data))
		} else {
			putContentMpm(rootId+expectedId, "")
		}

		indexId++
		if indexId == MaxIndex {
			break
		}
	}

	if totalGoRoutine > 0 {
		waitGroup.Done()
		totalGoRoutine--
	}
}

func getContentMpm(rawData string) string {
	strResult := ""
	if lengthData, err := strconv.Atoi(rawData[2:4]); err != nil {
		log.Error(err)
	} else {
		strResult = rawData[4 : 4+lengthData]
	}
	return strResult
}

func getSubContentMpm(rawData string, rootId string) {
	parseMPM(rawData, rootId+SEPARATOR)
}

func putContentMpm(key string, data string) {
	contents[key] = data
}

func stripContent(rawData string, length int) string {
	return rawData[4+length:]
}

func isRootIdHaveSubId(rootId string) bool {
	rootIdInt := 0
	var err error
	if rootIdInt, err = strconv.Atoi(rootId); err != nil {
		log.Error(err)
	} else if rootIdInt >= 2 && rootIdInt <= 51 {
		return true
	} else if rootIdInt == 62 {
		return true
	}
	return rootIdInt == 64
}

func getQrisDataWithoutCrc(mapContent map[string]string) string {
	var qrisData string
	for i := 0; i < 63; i++ {
		key := fmt.Sprintf("%02d", i)
		if len(mapContent[key]) > 0 && mapContent[key] != "" {
			qrisData += key
			contentLength := fmt.Sprintf("%02d", len(mapContent[key]))
			qrisData += contentLength
			qrisData += mapContent[key]
		}
	}
	qrisData += "6304"
	return qrisData
}

func CompareCrc(mapContent map[string]string, crc string) bool {
	calculatedSum := CheckSum(getQrisDataWithoutCrc(mapContent))
	return calculatedSum == crc
}

func GetResultMpm(raw string) map[string]string {
	parseAsMap(raw)
	return contents
}

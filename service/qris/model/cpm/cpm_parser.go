package cpm

import (
	"encoding/base64"
	"encoding/hex"
	"fatalisa-public-api/service/qris/model/mpm"
	"github.com/subchen/go-log"
	"regexp"
	"strconv"
	"strings"
)

var contents map[string]string

func init() {
	contents = make(map[string]string)
}

var idList = []string{"4F", "50", "57", "5A", "5F", "9F", "63"}
var id5FList = []string{"5F20", "5F2D", "5F50", "5F55"}
var id9FList = []string{"9F08", "9F19", "9F24", "9F25", "9F76"}
var id63List = []string{"9F26", "9F27", "9F10", "9F36", "82", "9F37", "9F74"}
var keysOfNumberData = []string{
	"61", "4F", "57", "90F8", "9F19", "9F25",
	"63", "9F26", "9F27", "9F10", "9F36", "82", "9F37",
}

func Parse(base64Str string) {
	if decodedBytesFromBase64, errDecodeBase64 := base64.StdEncoding.DecodeString(base64Str); errDecodeBase64 != nil {
		log.Error(errDecodeBase64)
	} else {
		rawHex := strings.ToUpper(hex.EncodeToString(decodedBytesFromBase64))
		id85Content := getContentCpm(rawHex, false)
		putContentCpm("85", id85Content)

		id61Content := getContentCpmCustom(rawHex, "4F")
		putContentCpm("64", id61Content)

		getLoopContent(idList, id61Content, false)
	}
}

func getContentCpm(rawHex string, is4Digit bool) string {
	var dataLength, start int64
	var err error
	if is4Digit {
		start = 6
		if dataLength, err = strconv.ParseInt(rawHex[4:6], 16, 64); err != nil {
			log.Error(err)
		}
	} else {
		start = 4
		if dataLength, err = strconv.ParseInt(rawHex[2:4], 16, 64); err != nil {
			log.Error(err)
		}
	}
	return rawHex[start : start+(dataLength*2)]
}

func getContentCpmCustom(rawHex string, patternToFind string) string {
	var result string
	pattern := regexp.MustCompile(patternToFind + "(.*)")
	//pattern := regexp.MustCompilePOSIX(patternToFind + "(.*)")
	res := pattern.FindAllString(rawHex, -1)
	for _, v := range res {
		result = v
	}
	return result
}

func putContentCpm(key string, hexData string) {
	number := isNumberData(key)
	numericAcc := isNumericAccount(key)
	data := hexData
	if !number && !numericAcc {
		if decodedHex, err := hex.DecodeString(hexData); err != nil {
			log.Error(err)
		} else {
			data = string(decodedHex)
		}
	} else if numericAcc {
		data = strings.ReplaceAll(data, "F", "")
	}
	contents[key] = data
}

func isNumberData(key string) bool {
	for _, v := range keysOfNumberData {
		if key == v {
			return true
		}
	}
	return false
}

func isNumericAccount(key string) bool {
	return key == "5A"
}

func getLoopContent(idList []string, raw string, isSub bool) string {
	for _, currentId := range idList {
		if len(raw) == 0 {
			break
		}

		id := raw[:len(currentId)]
		if id == currentId {
			if id == "5F" && !isSub {
				raw = getLoopContent(id5FList, raw, true)
			} else if id == "9F" && !isSub {
				raw = getLoopContent(id9FList, raw, true)
			} else if id == "63" && !isSub {
				raw = getContentCpmCustom(raw, "9F")
				raw = getLoopContent(id63List, raw, true)
			} else if id == "82" || !isSub {
				content := getContentCpm(raw, false)
				putContentCpm(currentId, content)
				raw = mpm.StripContent(raw, len(content))
			} else {
				content := getContentCpm(raw, true)
				putContentCpm(currentId, content)
				raw = mpm.StripContent(raw, len(content)+2)
			}
		} else {
			putContentCpm(currentId, "")
		}
	}
	return raw
}

func GetResultCpm(raw string) map[string]string {
	Parse(raw)
	return contents
}

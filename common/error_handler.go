package common

import "github.com/pieterclaerhout/go-log"

func ErrPrint(err error, isFatal bool) {
	if err != nil {
		if isFatal {
			log.Fatal(err)
		} else {
			log.Error(err)
		}
	}
}

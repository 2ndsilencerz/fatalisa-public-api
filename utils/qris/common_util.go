package qris

import (
	"fmt"
	"strings"
)

func CheckSum(str string) string {
	output := ""
	crc := 0xFFFF
	polynomial := 0x1021

	bytes := []byte(str)
	for _, v := range bytes {
		for i := 0; i < 8; i++ {
			bit := func() bool {
				return (v >> (7 - i) & 1) == 1
			}()
			c15 := func() bool {
				return (crc >> 15 & 1) == 1
			}()
			if bit != c15 {
				crc ^= polynomial
			}
		}
	}

	crc &= 0xffff
	output = fmt.Sprintf("%04x", crc)
	output = strings.ToUpper(output)
	return output
}

package qris

import (
	"fatalisa-public-api/utils"
	"github.com/pieterclaerhout/go-log"
	"testing"
)

type cpmTestData struct {
	raw string
	*CpmData
}

var testDataCpm *cpmTestData

func init() {
	testDataCpm = &cpmTestData{
		raw: "hQVDUFYwMWGTTwegAAAGAiAgUAdxcmlzY3BtWgqTYAUDMSNFZ4mfXyALUmlraSBEZXJpYW5fLQRpZGVuX1AXcmlraS5kZXJpYW5AcXJpc2NwbS5jb22fJQJ4mWM/n3Q8MTIzNDU2Nzg5MDEyMzQ1Njc4OTAxMjM0NTY3ODkwMTIzNDU2Nzg5MDEyMzQ1Njc4OTAxMjM0NTY3ODkw",
	}
}

func TestCpmParse(t *testing.T) {
	cpmData := &CpmData{}
	cpmData.GetData(testDataCpm.raw)
	log.Info(utils.Jsonify(cpmData))

	if cpmData.PayloadFormatIndicator == "" || cpmData.ApplicationPAN == "" || cpmData.IssuerURL == "" {
		t.Error()
	}
}

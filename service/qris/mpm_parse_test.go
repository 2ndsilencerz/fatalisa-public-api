package qris

import (
	"fatalisa-public-api/utils"
	"github.com/pieterclaerhout/go-log"
	"reflect"
	"testing"
)

type MpmTestData struct {
	Raw string
	*MpmData
}

var testDataMpm *MpmTestData

func init() {
	testDataMpm = &MpmTestData{
		Raw: "00020101021126670018ID.CO.EXAMPLE2.WWW01159360056701234560215MIDCONTOH1234560303UMI5204123453033605502015802ID5914NamaMerchantC76009NamaKota16110123456789062070703K1963040BE8",
		MpmData: &MpmData{
			GlobalUniqueIdentifier: "ID.CO.EXAMPLE2.WWW",
			MerchantPAN:            "936005670123456",
			MerchantID:             "MIDCONTOH123456",
			MerchantCriteria:       "UMI",
			MerchantCategoryCode:   "1234",
			TransactionCurrency:    "360",
			TipIndicator:           "01",
			CountryCode:            "ID",
			MerchantName:           "NamaMerchantC7",
			MerchantCity:           "NamaKota1",
			PostalCode:             "1234567890",
			AdditionalDataField:    "0703K19",
			TerminalLabel:          "K19",
			Crc:                    "0BE8",
		},
	}
}

func TestMpmParse(t *testing.T) {
	mpmData := &MpmData{}
	mpmData.GetData(testDataMpm.Raw)
	jsonFormat := utils.Jsonify(mpmData)
	log.Info(jsonFormat)

	if reflect.DeepEqual(&mpmData, &testDataMpm.MpmData) {
		t.Error()
	}
	if CompareCrc(GetResultMpm(testDataMpm.Raw), "0BE8") {
		t.Error()
	}
}

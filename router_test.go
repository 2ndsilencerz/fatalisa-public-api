package main

import (
	"bytes"
	"encoding/json"
	"fatalisa-public-api/router"
	"fatalisa-public-api/service/common/pray-schedule"
	"fatalisa-public-api/service/qris"
	"fatalisa-public-api/utils"
	"github.com/subchen/go-log"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

func setupRouter() *router.Config {
	routerInit := &router.Config{}
	routerInit.Get()
	routerInit.InitRoutes()
	return routerInit
}

func sendAsGet(uri string) []byte {
	var getResult []byte
	routerTest := setupRouter()
	httpRes := httptest.NewRecorder()
	httpReq, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		log.Error(err)
	} else {
		routerTest.Gin.ServeHTTP(httpRes, httpReq)
		if rawRes, err := ioutil.ReadAll(httpRes.Body); err != nil {
			log.Error(err)
		} else {
			getResult = rawRes
		}
	}
	return getResult
}

func sendAsPost(uri string, body []byte) []byte {
	var postResult []byte
	routerTest := setupRouter()
	httpRes := httptest.NewRecorder()
	httpReq, err := http.NewRequest("POST", uri, bytes.NewBuffer(body))
	if err != nil {
		log.Error(err)
	} else {
		routerTest.Gin.ServeHTTP(httpRes, httpReq)
		if rawRes, err := ioutil.ReadAll(httpRes.Body); err != nil {
			log.Error(err)
		} else {
			postResult = rawRes
		}
	}
	return postResult
}

func TestGetSchedule(t *testing.T) {
	pray_schedule.DownloadFile(83)
	current := time.Now()
	city := "jakarta"
	bodyReq := &pray_schedule.PrayScheduleReq{
		City: city,
		Date: current.Format("2006/01/02"),
	}

	if bodyReqJson := utils.Jsonify(bodyReq); len(bodyReqJson) > 0 {
		dataRes := &pray_schedule.PrayScheduleData{}
		rawRes := sendAsPost("/api/pray-schedule", []byte(bodyReqJson))
		if err := json.Unmarshal(rawRes, dataRes); err != nil {
			t.Error(err)
		}
		if current.Format("2006") != dataRes.Year &&
			current.Format("01") != dataRes.Month &&
			current.Format("02") != dataRes.Date {
			t.Errorf("Wrong date, expected %s got %s/%s/%s",
				current.Format("2006/01/02"), dataRes.Year, dataRes.Month, dataRes.Date)
		}
	}
}

func TestParseMpmGet(t *testing.T) {
	testDataMpm := struct {
		Raw     string
		MpmData *qris.MpmData
	}{
		"00020101021126670018ID.CO.EXAMPLE2.WWW01159360056701234560215MIDCONTOH1234560303UMI5204123453033605502015802ID5914NamaMerchantC76009NamaKota16110123456789062070703K1963040BE8",
		&qris.MpmData{
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

	dataRes := &qris.MpmData{}
	rawRes := sendAsGet("/api/qris/mpm/" + testDataMpm.Raw)
	if err := json.Unmarshal(rawRes, dataRes); err != nil {
		t.Error(err)
	} else {
		if reflect.DeepEqual(&dataRes, &testDataMpm.MpmData) {
			t.Error()
		}
		if qris.CompareCrc(qris.GetResultMpm(testDataMpm.Raw), "0BE8") {
			t.Error()
		}
	}
}

func TestParseMpmPost(t *testing.T) {
	testDataMpm := struct {
		Raw     string
		MpmData *qris.MpmData
	}{
		"00020101021126670018ID.CO.EXAMPLE2.WWW01159360056701234560215MIDCONTOH1234560303UMI5204123453033605502015802ID5914NamaMerchantC76009NamaKota16110123456789062070703K1963040BE8",
		&qris.MpmData{
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

	req := qris.MpmRequest{
		Raw: testDataMpm.Raw,
	}
	if jsonBody := utils.Jsonify(req); len(jsonBody) > 0 {
		dataRes := &qris.MpmData{}
		rawRes := sendAsPost("/api/qris/mpm", []byte(jsonBody))
		if err := json.Unmarshal(rawRes, dataRes); err != nil {
			t.Error(err)
		} else {
			if reflect.DeepEqual(&dataRes, &testDataMpm.MpmData) {
				t.Error("MPM Data not match")
			}
			if qris.CompareCrc(qris.GetResultMpm(testDataMpm.Raw), "0BE8") {
				t.Error("CRC not match")
			}
		}
	}
}

func TestParseCpmPost(t *testing.T) {
	testDataCpm := struct {
		raw string
	}{
		raw: "hQVDUFYwMWGTTwegAAAGAiAgUAdxcmlzY3BtWgqTYAUDMSNFZ4mfXyALUmlraSBEZXJpYW5fLQRpZGVuX1AXcmlraS5kZXJpYW5AcXJpc2NwbS5jb22fJQJ4mWM/n3Q8MTIzNDU2Nzg5MDEyMzQ1Njc4OTAxMjM0NTY3ODkwMTIzNDU2Nzg5MDEyMzQ1Njc4OTAxMjM0NTY3ODkw",
	}

	req := qris.CpmRequest{
		Raw: testDataCpm.raw,
	}
	if jsonBody := utils.Jsonify(req); len(jsonBody) > 0 {
		dataRes := &qris.CpmData{}
		rawRes := sendAsPost("/api/qris/cpm", []byte(jsonBody))
		if err := json.Unmarshal(rawRes, dataRes); err != nil {
			t.Error(err)
		} else if dataRes.PayloadFormatIndicator == "" || dataRes.ApplicationPAN == "" || dataRes.IssuerURL == "" {
			t.Error("CPM Data Not Match")
		}
	}
}

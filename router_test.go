package main

import (
	"bytes"
	"encoding/json"
	"fatalisa-public-api/router"
	pray_schedule "fatalisa-public-api/service/pray-schedule"
	"fatalisa-public-api/service/pray-schedule/model"
	prayschedulepkpu "fatalisa-public-api/service/pray-schedule/pkpu"
	"fatalisa-public-api/service/qris/model/cpm"
	"fatalisa-public-api/service/qris/model/mpm"
	utils2 "fatalisa-public-api/service/web/utils"
	"fmt"
	"github.com/subchen/go-log"
	"io"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"testing"
	"time"
)

const (
	apiUri  = "/api"
	qrisUri = apiUri + "/qris"
	mpmUri  = qrisUri + "/mpm"
	cpmUri  = qrisUri + "/cpm"
	//cityTestPkpu            = "jakarta"
	cityTestKemenag         = "KOTA JAKARTA"
	prayScheduleUri         = apiUri + "/pray-schedule"
	prayScheduleCityListUri = prayScheduleUri + "/city-list"
)

type mpmTestData struct {
	Raw     string
	MpmData *mpm.Data
}

func setupRouter() *router.Config {
	routerInit := &router.Config{}
	routerInit.Get()
	routerInit.InitRoutes()
	return routerInit
}

func sendData(method string, uri string, data interface{}) []byte {
	var body []byte
	if data != nil {
		body, _ = json.Marshal(data)
	}
	var postResult []byte
	routerTest := setupRouter()
	//httpRes := httptest.NewRecorder()
	httpReq, err := http.NewRequest(method, uri, bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")
	if err, _ := utils2.ErrorHandler(err); !err {
		log.Info("Request is " + fmt.Sprint(httpReq))
		httpRes, err := routerTest.Fiber.Test(httpReq, -1)
		utils2.ErrorHandler(err)
		if rawRes, err := io.ReadAll(httpRes.Body); err != nil {
			log.Error(err)
		} else {
			postResult = rawRes
		}
	}
	log.Info(string(postResult))
	return postResult
}

func getMpmTestData() mpmTestData {
	return struct {
		Raw     string
		MpmData *mpm.Data
	}{
		"00020101021126670018ID.CO.EXAMPLE2.WWW01159360056701234560215MIDCONTOH1234560303UMI5204123453033605502015802ID5914NamaMerchantC76009NamaKota16110123456789062070703K1963040BE8",
		&mpm.Data{
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

func TestParseMpmGet(t *testing.T) {
	testDataMpm := getMpmTestData()

	dataRes := &mpm.Data{}
	rawRes := sendData(http.MethodGet, mpmUri+"/"+testDataMpm.Raw, nil)
	if err := json.Unmarshal(rawRes, dataRes); err != nil {
		t.Error(err)
	} else {
		if reflect.DeepEqual(&dataRes, &testDataMpm.MpmData) {
			t.Error()
		}
		if mpm.CompareCrc(mpm.GetResultMpm(testDataMpm.Raw), "0BE8") {
			t.Error()
		}
	}
}

func TestParseMpmPost(t *testing.T) {
	testDataMpm := getMpmTestData()

	dataRes := &mpm.Data{}
	rawRes := sendData(http.MethodPost, mpmUri, testDataMpm)
	if err := json.Unmarshal(rawRes, dataRes); err != nil {
		t.Error(err)
	} else {
		if reflect.DeepEqual(&dataRes, &testDataMpm.MpmData) {
			t.Error("MPM Data not match")
		}
		if mpm.CompareCrc(mpm.GetResultMpm(testDataMpm.Raw), "0BE8") {
			t.Error("CRC not match")
		}
	}
}

func TestParseCpmPost(t *testing.T) {
	testDataCpm := struct {
		Raw string
	}{
		Raw: "hQVDUFYwMWGTTwegAAAGAiAgUAdxcmlzY3BtWgqTYAUDMSNFZ4mfXyALUmlraSBEZXJpYW5fLQRpZGVuX1AXcmlraS5kZXJpYW5AcXJpc2NwbS5jb22fJQJ4mWM/n3Q8MTIzNDU2Nzg5MDEyMzQ1Njc4OTAxMjM0NTY3ODkwMTIzNDU2Nzg5MDEyMzQ1Njc4OTAxMjM0NTY3ODkw",
	}

	dataRes := &cpm.Data{}
	rawRes := sendData(http.MethodPost, cpmUri, testDataCpm)
	if err := json.Unmarshal(rawRes, dataRes); err != nil {
		t.Error(err)
	} else if dataRes.PayloadFormatIndicator == "" || dataRes.ApplicationPAN == "" || dataRes.IssuerURL == "" {
		t.Error("CPM Data Not Match")
	}
}

func init() {
	if os.Getenv("PROVIDER") == pray_schedule.ProviderPkpu {
		log.Info("Downloading pray schedule for test")
		prayschedulepkpu.DownloadFile(83)
	}
}

func TestGetSchedulePost(t *testing.T) {
	_ = os.Setenv("PROVIDER", pray_schedule.ProviderKemenag)
	req := &model.Request{
		City: cityTestKemenag,
		Date: time.Now().Format("2006/01/02"),
	}

	dataRes := &model.Response{}
	rawRes := sendData(http.MethodPost, prayScheduleUri, req)
	if err := json.Unmarshal(rawRes, dataRes); err != nil {
		t.Error(err)
	} else {
		tmp, _ := time.Parse("2006/01/02", req.Date)
		if dataRes.Year != strconv.Itoa(tmp.Year()) ||
			dataRes.Month != fmt.Sprintf("%02s", strconv.Itoa(int(tmp.Month()))) ||
			dataRes.Date != fmt.Sprintf("%02s", strconv.Itoa(tmp.Day())) {
			t.Error("Data not match")
		}
	}
}

func TestGetScheduleGet(t *testing.T) {
	_ = os.Setenv("PROVIDER", pray_schedule.ProviderKemenag)
	dataRes := &model.Response{}
	rawRes := sendData(http.MethodGet, prayScheduleUri+"/"+cityTestKemenag, nil)
	if err := json.Unmarshal(rawRes, dataRes); err != nil {
		t.Error(err)
	} else {
		tmp := time.Now()
		if dataRes.Year != strconv.Itoa(tmp.Year()) ||
			dataRes.Month != fmt.Sprintf("%02s", strconv.Itoa(int(tmp.Month()))) ||
			dataRes.Date != fmt.Sprintf("%02s", strconv.Itoa(tmp.Day())) {
			t.Error("Data not match")
		}
	}
}

func TestScheduleCityList(t *testing.T) {
	_ = os.Setenv("PROVIDER", pray_schedule.ProviderKemenag)
	var dataRes model.CityList
	rawRes := sendData(http.MethodGet, prayScheduleCityListUri, nil)
	if err := json.Unmarshal(rawRes, &dataRes); err != nil {
		t.Error(err)
	} else {
		for _, data := range dataRes.List {
			if data.Name == cityTestKemenag {
				return
			}
		}
		t.Error("Data not found")
	}
}

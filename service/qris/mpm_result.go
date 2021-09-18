package qris

import "fmt"

type MpmRequest struct {
	Raw string `json:"raw" binding:"required"`
}

type MpmData struct {
	PayloadFormatIndicator        string `json:"payloadFormatIndicator"`
	PointOfInitiationMethod       string `json:"pointOfInitiationMethod"`
	GlobalUniqueIdentifier        string `json:"globalUniqueIdentifier"`
	MerchantPAN                   string `json:"merchantPAN"`
	MerchantID                    string `json:"merchantID"`
	MerchantCriteria              string `json:"merchantCriteria"`
	MerchantAccountInformation    string `json:"merchantAccountInformation"`
	MerchantCategoryCode          string `json:"merchantCategoryCode"`
	TransactionCurrency           string `json:"transactionCurrency"`
	TransactionAmount             string `json:"transactionAmount"`
	TipIndicator                  string `json:"tipIndicator"`
	TipFixedValue                 string `json:"tipFixedValue"`
	TipPercentageValue            string `json:"tipPercentageValue"`
	CountryCode                   string `json:"countryCode"`
	MerchantName                  string `json:"merchantName"`
	MerchantCity                  string `json:"merchantCity"`
	PostalCode                    string `json:"postalCode"`
	AdditionalDataField           string `json:"additionalDataField"`
	BillNumber                    string `json:"billNumber"`
	MobileNumber                  string `json:"mobileNumber"`
	StoreLabel                    string `json:"storeLabel"`
	LoyaltyNumber                 string `json:"loyaltyNumber"`
	ReferenceLabel                string `json:"referenceLabel"`
	CustomerLabel                 string `json:"customerLabel"`
	TerminalLabel                 string `json:"terminalLabel"`
	PurposeOfTransaction          string `json:"purposeOfTransaction"`
	AdditionalConsumerDataRequest string `json:"additionalConsumerDataRequest"`
	Crc                           string `json:"crc"`
	LanguagePreference            string `json:"languagePreference"`
	MerchantNameAlt               string `json:"merchantNameAlt"`
	MerchantCityAlt               string `json:"merchantCityAlt"`

	//ValidUntil time.Time `json:"validUntil"`
	//QrisStatus string `json:"qrisStatus"`
	//InactiveDate time.Time `json:"inactiveDate"`
}

func (data *MpmData) setContents(qrisParsed map[string]string) {
	data.PayloadFormatIndicator = qrisParsed["00"]
	data.PointOfInitiationMethod = qrisParsed["01"]
	for i := 2; i <= 45; i++ {
		index := fmt.Sprintf("%02d", i)
		if len(qrisParsed[index]) > 0 && qrisParsed[index] != "" {
			guid := qrisParsed[index+SEPARATOR+"00"]
			mpan := qrisParsed[index+SEPARATOR+"01"]
			mid := qrisParsed[index+SEPARATOR+"02"]
			criteria := qrisParsed[index+SEPARATOR+"03"]
			data.GlobalUniqueIdentifier = guid
			data.MerchantPAN = mpan
			data.MerchantID = mid
			data.MerchantCriteria = criteria
			break
		}
	}

	data.MerchantAccountInformation = qrisParsed["51"]
	if len(data.MerchantAccountInformation) > 0 && data.MerchantAccountInformation != "" {
		data.GlobalUniqueIdentifier = qrisParsed["51"+SEPARATOR+"00"]
		data.MerchantID = qrisParsed["51"+SEPARATOR+"02"]
		data.MerchantCriteria = qrisParsed["51"+SEPARATOR+"03"]
	}

	data.MerchantCategoryCode = qrisParsed["52"]
	data.TransactionCurrency = qrisParsed["53"]
	data.TransactionAmount = qrisParsed["54"]
	if len(qrisParsed["55"]) > 0 && qrisParsed["55"] != "" {
		data.TipIndicator = qrisParsed["55"]
		if data.TipIndicator == "02" {
			data.TipFixedValue = qrisParsed["56"]
		} else if data.TipIndicator == "03" {
			data.TipPercentageValue = qrisParsed["57"]
		}
	}

	data.CountryCode = qrisParsed["58"]
	data.MerchantName = qrisParsed["59"]
	data.MerchantCity = qrisParsed["60"]
	data.PostalCode = qrisParsed["61"]
	data.AdditionalDataField = qrisParsed["62"]
	if len(data.AdditionalDataField) > 0 && data.AdditionalDataField != "" {
		data.getBit62Contents(qrisParsed)
		if len(data.TerminalLabel) == 0 || data.TerminalLabel == "" {
			data.TerminalLabel = fmt.Sprintf("%-16s", func() string {
				tmp := qrisParsed["62"+SEPARATOR+"07"]
				if len(tmp) > 0 && tmp != "" {
					return tmp
				}
				return qrisParsed["62"]
			}())
		}
	}
	data.Crc = qrisParsed["63"]
	if len(qrisParsed["64"]) > 0 && qrisParsed["64"] != "" {
		data.getBit64Contents(qrisParsed)
	}
}

func (data *MpmData) getBit62Contents(qrisParsed map[string]string) {
	var contents [9]string
	for i := 1; i < 10; i++ {
		key := "62" + SEPARATOR + fmt.Sprintf("%02d", i)
		contents[i-1] = qrisParsed[key]
	}
	data.BillNumber = contents[0]
	data.MobileNumber = contents[1]
	data.StoreLabel = contents[2]
	data.LoyaltyNumber = contents[3]
	data.ReferenceLabel = contents[4]
	data.CustomerLabel = contents[5]
	data.TerminalLabel = contents[6]
	data.PurposeOfTransaction = contents[7]
	data.AdditionalConsumerDataRequest = contents[8]
}

func (data *MpmData) getBit64Contents(qrisParsed map[string]string) {
	var contents [3]string
	for i := 0; i < 3; i++ {
		key := "64" + SEPARATOR + fmt.Sprintf("%02d", i)
		contents[i] = qrisParsed[key]
	}
	data.LanguagePreference = contents[0]
	data.MerchantNameAlt = contents[1]
	data.MerchantCityAlt = contents[2]
}

func (data *MpmData) GetData(raw string) {
	qrisParsed := GetResultMpm(raw)
	data.setContents(qrisParsed)
}

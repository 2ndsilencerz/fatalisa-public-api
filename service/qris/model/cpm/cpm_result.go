package cpm

import (
	"fatalisa-public-api/database/config"
)

type Data struct {
	PayloadFormatIndicator                 string `json:"payloadFormatIndicator"`
	ApplicationTemplate                    string `json:"applicationTemplate"`
	ApplicationDefinitionFileName          string `json:"applicationDefinitionFileName"`
	ApplicationLabel                       string `json:"applicationLabel"`
	Track2EquivalentData                   string `json:"track2EquivalentData"`
	ApplicationPAN                         string `json:"applicationPAN"`
	CardHolderName                         string `json:"cardHolderName"`
	LanguagePreference                     string `json:"languagePreference"`
	IssuerURL                              string `json:"issuerURL"`
	ApplicationVersionNumber               string `json:"applicationVersionNumber"`
	TokenRequesterID                       string `json:"tokenRequesterID"`
	PaymentAccountReference                string `json:"paymentAccountReference"`
	Last4DigitPAN                          string `json:"last4DigitPAN"`
	ApplicationSpecificTransparentTemplate string `json:"applicationSpecificTransparentTemplate"`
	ApplicationCryptogram                  string `json:"applicationCryptogram"`
	CryptogramInformationData              string `json:"cryptogramInformationData"`
	IssuerApplicationData                  string `json:"issuerApplicationData"`
	ApplicationTransactionCounter          string `json:"applicationTransactionCounter"`
	ApplicationInterchangeProfile          string `json:"applicationInterchangeProfile"`
	UnpredictableNumber                    string `json:"unpredictableNumber"`
	IssuerQRISData                         string `json:"issuerQRISData"`
}

func (data *Data) setContents(qrisParsed map[string]string) {
	data.PayloadFormatIndicator = qrisParsed["85"]
	data.ApplicationTemplate = qrisParsed["61"]
	data.ApplicationDefinitionFileName = qrisParsed["4F"]
	data.ApplicationLabel = qrisParsed["50"]
	data.Track2EquivalentData = qrisParsed["57"]
	data.ApplicationPAN = qrisParsed["5A"]
	data.CardHolderName = qrisParsed["5F20"]
	data.LanguagePreference = qrisParsed["5F2D"]
	data.IssuerURL = qrisParsed["5F50"]
	data.ApplicationVersionNumber = qrisParsed["9F08"]
	data.TokenRequesterID = qrisParsed["9F19"]
	data.PaymentAccountReference = qrisParsed["9F24"]
	data.Last4DigitPAN = qrisParsed["9F25"]
	//data.ApplicationSpecificTransparentTemplate = qrisParsed["63"]
	data.ApplicationCryptogram = qrisParsed["9F26"]
	data.CryptogramInformationData = qrisParsed["9F27"]
	data.IssuerApplicationData = qrisParsed["9F10"]
	data.ApplicationTransactionCounter = qrisParsed["9F36"]
	data.ApplicationInterchangeProfile = qrisParsed["82"]
	data.UnpredictableNumber = qrisParsed["9F37"]
	data.IssuerQRISData = qrisParsed["9F74"]
}

func (data *Data) Parse(raw string) {
	qrisParsed := GetResultCpm(raw)
	data.setContents(qrisParsed)
}

func (data *Data) SaveToDB() {
	mariadb := config.InitMariaDB()
	mariadb.Write(data)

	postgres := config.InitPostgres()
	postgres.Write(data)

	mongo := config.InitMongoDB()
	mongo.InsertOne("CPM", data)
}

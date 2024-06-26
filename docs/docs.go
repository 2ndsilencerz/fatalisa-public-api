// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/pray-schedule": {
            "post": {
                "description": "GetString Schedule By City and Date",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Pray-Schedule"
                ],
                "summary": "PrayScheduleCityPost",
                "parameters": [
                    {
                        "description": "data",
                        "name": "city",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "object"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.Response"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {}
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {}
                    }
                }
            }
        },
        "/api/pray-schedule/:city": {
            "get": {
                "description": "GetString Schedule By City",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Pray-Schedule"
                ],
                "summary": "PrayScheduleCity",
                "parameters": [
                    {
                        "type": "string",
                        "description": "city",
                        "name": "city",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.Response"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {}
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {}
                    }
                }
            }
        },
        "/api/pray-schedule/city-list": {
            "get": {
                "description": "GetString Available City List",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Pray-Schedule"
                ],
                "summary": "PrayScheduleCityList",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/model.CityList"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {}
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {}
                    }
                }
            }
        },
        "/api/qris/cpm": {
            "post": {
                "description": "Parse CPM by raw body",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "CPM"
                ],
                "summary": "ParseCpm",
                "parameters": [
                    {
                        "description": "data",
                        "name": "raw",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/cpm.Request"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/cpm.Data"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {}
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {}
                    }
                }
            }
        },
        "/api/qris/mpm": {
            "post": {
                "description": "Parse MPM by raw body",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "MPM"
                ],
                "summary": "ParseMpm",
                "parameters": [
                    {
                        "description": "data",
                        "name": "raw",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/mpm.Request"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/mpm.Data"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {}
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {}
                    }
                }
            }
        },
        "/api/qris/mpm/:raw": {
            "get": {
                "description": "Parse MPM by parameter",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "MPM"
                ],
                "summary": "ParseMpm",
                "parameters": [
                    {
                        "type": "string",
                        "description": "raw",
                        "name": "raw",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/mpm.Data"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {}
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {}
                    }
                }
            }
        },
        "/health": {
            "get": {
                "description": "Health Check",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Default"
                ],
                "summary": "HealthCheck",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/common.Body"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {}
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {}
                    }
                }
            }
        },
        "/version": {
            "get": {
                "description": "Show Version",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Default"
                ],
                "summary": "VersionChecker",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/common.Body"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {}
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {}
                    }
                }
            }
        }
    },
    "definitions": {
        "common.Body": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                }
            }
        },
        "cpm.Data": {
            "type": "object",
            "properties": {
                "applicationCryptogram": {
                    "type": "string"
                },
                "applicationDefinitionFileName": {
                    "type": "string"
                },
                "applicationInterchangeProfile": {
                    "type": "string"
                },
                "applicationLabel": {
                    "type": "string"
                },
                "applicationPAN": {
                    "type": "string"
                },
                "applicationSpecificTransparentTemplate": {
                    "type": "string"
                },
                "applicationTemplate": {
                    "type": "string"
                },
                "applicationTransactionCounter": {
                    "type": "string"
                },
                "applicationVersionNumber": {
                    "type": "string"
                },
                "cardHolderName": {
                    "type": "string"
                },
                "cryptogramInformationData": {
                    "type": "string"
                },
                "issuerApplicationData": {
                    "type": "string"
                },
                "issuerQRISData": {
                    "type": "string"
                },
                "issuerURL": {
                    "type": "string"
                },
                "languagePreference": {
                    "type": "string"
                },
                "last4DigitPAN": {
                    "type": "string"
                },
                "payloadFormatIndicator": {
                    "type": "string"
                },
                "paymentAccountReference": {
                    "type": "string"
                },
                "tokenRequesterID": {
                    "type": "string"
                },
                "track2EquivalentData": {
                    "type": "string"
                },
                "unpredictableNumber": {
                    "type": "string"
                }
            }
        },
        "cpm.Request": {
            "type": "object",
            "required": [
                "raw"
            ],
            "properties": {
                "raw": {
                    "type": "string"
                }
            }
        },
        "model.City": {
            "type": "object",
            "properties": {
                "cityName": {
                    "type": "string"
                }
            }
        },
        "model.CityList": {
            "type": "object",
            "properties": {
                "list": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/model.City"
                    }
                }
            }
        },
        "model.Response": {
            "type": "object",
            "required": [
                "ashr",
                "date",
                "dzuhur",
                "fajr",
                "imsyak",
                "isha",
                "maghrib",
                "month",
                "syuruq",
                "year"
            ],
            "properties": {
                "ashr": {
                    "type": "string"
                },
                "date": {
                    "type": "string"
                },
                "dzuhur": {
                    "type": "string"
                },
                "fajr": {
                    "type": "string"
                },
                "imsyak": {
                    "type": "string"
                },
                "isha": {
                    "type": "string"
                },
                "maghrib": {
                    "type": "string"
                },
                "month": {
                    "type": "string"
                },
                "syuruq": {
                    "type": "string"
                },
                "year": {
                    "type": "string"
                }
            }
        },
        "mpm.Data": {
            "type": "object",
            "properties": {
                "additionalConsumerDataRequest": {
                    "type": "string"
                },
                "additionalDataField": {
                    "type": "string"
                },
                "billNumber": {
                    "type": "string"
                },
                "countryCode": {
                    "type": "string"
                },
                "crc": {
                    "type": "string"
                },
                "customerLabel": {
                    "type": "string"
                },
                "globalUniqueIdentifier": {
                    "type": "string"
                },
                "languagePreference": {
                    "type": "string"
                },
                "loyaltyNumber": {
                    "type": "string"
                },
                "merchantAccountInformation": {
                    "type": "string"
                },
                "merchantCategoryCode": {
                    "type": "string"
                },
                "merchantCity": {
                    "type": "string"
                },
                "merchantCityAlt": {
                    "type": "string"
                },
                "merchantCriteria": {
                    "type": "string"
                },
                "merchantID": {
                    "type": "string"
                },
                "merchantName": {
                    "type": "string"
                },
                "merchantNameAlt": {
                    "type": "string"
                },
                "merchantPAN": {
                    "type": "string"
                },
                "mobileNumber": {
                    "type": "string"
                },
                "payloadFormatIndicator": {
                    "type": "string"
                },
                "pointOfInitiationMethod": {
                    "type": "string"
                },
                "postalCode": {
                    "type": "string"
                },
                "purposeOfTransaction": {
                    "type": "string"
                },
                "referenceLabel": {
                    "type": "string"
                },
                "storeLabel": {
                    "type": "string"
                },
                "terminalLabel": {
                    "type": "string"
                },
                "tipFixedValue": {
                    "type": "string"
                },
                "tipIndicator": {
                    "type": "string"
                },
                "tipPercentageValue": {
                    "type": "string"
                },
                "transactionAmount": {
                    "type": "string"
                },
                "transactionCurrency": {
                    "type": "string"
                }
            }
        },
        "mpm.Request": {
            "type": "object",
            "required": [
                "raw"
            ],
            "properties": {
                "raw": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}

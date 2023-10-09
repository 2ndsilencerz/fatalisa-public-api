package model

// Header is the header of xml used to parse schedule data downloaded
type Header struct {
	Adzan `xml:"adzan"`
}

// Parameter used to get misc data from schedule in xml
type Parameter struct {
	Longitude string `xml:"longitude"`
	Latitude  string `xml:"latitude"`
	Direction string `xml:"direction"`
	Distance  string `xml:"distance"`
}

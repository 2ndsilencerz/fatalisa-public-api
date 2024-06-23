package model

type Response struct {
	City    string `xml:"city" json:"city" binding:"required"`
	Year    string `xml:"year" json:"year" binding:"required"`
	Month   string `xml:"month" json:"month" binding:"required"`
	Date    string `xml:"date" json:"date" binding:"required"`
	Fajr    string `xml:"fajr" json:"fajr" binding:"required"`
	Dzuhur  string `xml:"dzuhr" json:"dzuhur" binding:"required"`
	Ashr    string `xml:"ashr" json:"ashr" binding:"required"`
	Maghrib string `xml:"maghrib" json:"maghrib" binding:"required"`
	Isha    string `xml:"isha" json:"isha" binding:"required"`
}

type CityList struct {
	List []*City `json:"list"`
}

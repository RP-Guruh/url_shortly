package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/ipinfo/go/v2/ipinfo"
	"github.com/joho/godotenv"
	"github.com/ua-parser/uap-go/uaparser"
)

type IPInfoResponse struct {
	IP          string `json:"ip"`
	City        string `json:"city"`
	Region      string `json:"region"`
	Country     string `json:"country"`
	CountryName string `json:"country_name"`
	CountryFlag struct {
		Emoji   string `json:"emoji"`
		Unicode string `json:"unicode"`
	} `json:"country_flag"`
	CountryCurrency struct {
		Code   string `json:"code"`
		Symbol string `json:"symbol"`
	} `json:"country_currency"`
	Continent struct {
		Code string `json:"code"`
		Name string `json:"name"`
	} `json:"continent"`
	Loc      string `json:"loc"`      // "-7.1955,107.4313"
	Org      string `json:"org"`      // Organization
	Postal   string `json:"postal"`   // Postal code
	Timezone string `json:"timezone"` // Asia/Jakarta
}

type VisitorInfo struct {
	IDURL         string `json:"id_url"`
	IPAddress     string `json:"ip_address"`
	Browser       string `json:"browser"`
	Device        string `json:"device"`
	OS            string `json:"os"`
	City          string `json:"city"`
	Region        string `json:"region"`
	Country       string `json:"country"`        // kode negara, misalnya "ID"
	CountryName   string `json:"country_name"`   // nama negara, misalnya "Indonesia"
	Continent     string `json:"continent"`      // kode benua, misalnya "AS"
	ContinentName string `json:"continent_name"` // nama benua, misalnya "Asia"
	LatLong       string `json:"latlong"`        // format: "latitude,longitude"
	Org           string `json:"org"`            // organisasi/ISP
	Postal        string `json:"postal"`         // kode pos
	Timezone      string `json:"timezone"`       // contoh: "Asia/Jakarta"
}

func initIpAPI() string {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	token := os.Getenv("IPINFO_TOKEN")
	return token
}

func getInformation(c *fiber.Ctx, ip string) VisitorInfo {
	parser := uaparser.NewFromSaved()
	myToken := initIpAPI()
	client := ipinfo.NewClient(nil, nil, myToken)
	info, err := client.GetIPInfo(net.ParseIP(ip))
	if err != nil {
		log.Printf("error getting IP info: %v", err)
		return VisitorInfo{} // Jangan hentikan program
	}

	// ambil User-Agent dari header
	uaString := c.Get("User-Agent")
	clientUA := parser.Parse(uaString)

	return VisitorInfo{
		IPAddress:     info.IP.String(),
		Browser:       fmt.Sprintf("%s %s", clientUA.UserAgent.Family, clientUA.UserAgent.ToVersionString()),
		Device:        clientUA.Device.Family,
		OS:            fmt.Sprintf("%s %s", clientUA.Os.Family, clientUA.Os.ToVersionString()),
		City:          info.City,
		Region:        info.Region,
		Country:       info.Country,
		CountryName:   info.CountryName,
		Continent:     info.Continent.Code,
		ContinentName: info.Continent.Name,
		LatLong:       info.Location,
		Org:           info.Org,
		Postal:        info.Postal,
		Timezone:      info.Timezone,
	}
}

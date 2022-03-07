package main

import (
	"encoding/json"
	"time"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"os/exec"
)

// Date Gathering
func getDate() string {
    out, err := exec.Command("date").Output()
    if err != nil {
	return err.Error()
    }
    return fmt.Sprintf("%s", out);
}

// Internal Network IP Gathering
func getIp() string {
    addrs, err := net.InterfaceAddrs()
    if err != nil {
	return fmt.Sprintf("Oops: %s", err.Error())
    }

    for _, a := range addrs {
	if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
	    if ipnet.IP.To4() != nil {
		return ipnet.IP.String()
	    }
	}
    } 
    return "NOT CONNECTED" 
}

// Crypto HTTP Section
type Crypto []struct {
	ID                           string      `json:"id"`
	Symbol                       string      `json:"symbol"`
	Name                         string      `json:"name"`
	Image                        string      `json:"image"`
	CurrentPrice                 float64     `json:"current_price"`
	MarketCap                    int64       `json:"market_cap"`
	MarketCapRank                int         `json:"market_cap_rank"`
	FullyDilutedValuation        int64       `json:"fully_diluted_valuation"`
	TotalVolume                  int64       `json:"total_volume"`
	High24H                      int         `json:"high_24h"`
	Low24H                       int         `json:"low_24h"`
	PriceChange24H               float64     `json:"price_change_24h"`
	PriceChangePercentage24H     float64     `json:"price_change_percentage_24h"`
	MarketCapChange24H           float64     `json:"market_cap_change_24h"`
	MarketCapChangePercentage24H float64     `json:"market_cap_change_percentage_24h"`
	CirculatingSupply            float64     `json:"circulating_supply"`
	TotalSupply                  float64     `json:"total_supply"`
	MaxSupply                    float64     `json:"max_supply"`
	Ath                          int         `json:"ath"`
	AthChangePercentage          float64     `json:"ath_change_percentage"`
	AthDate                      time.Time   `json:"ath_date"`
	Atl                          float64     `json:"atl"`
	AtlChangePercentage          float64     `json:"atl_change_percentage"`
	AtlDate                      time.Time   `json:"atl_date"`
	Roi                          interface{} `json:"roi"`
	LastUpdated                  time.Time   `json:"last_updated"`
}

func getCrypto() string {
    var crypto Crypto
    var stringify string
    

    resp, err := http.Get("https://api.coingecko.com/api/v3/coins/markets?vs_currency=usd&ids=bitcoin,ethereum,solana")
    if err != nil {
	return err.Error()
    }
    defer resp.Body.Close()
    
    respBytes, _ := ioutil.ReadAll(resp.Body) 
    json.Unmarshal(respBytes, &crypto)

    for _, asset := range crypto {
	stringify += fmt.Sprintf(" %s: %d$ ", strings.ToUpper(asset.Symbol), int64(asset.CurrentPrice)) 
    }
    return stringify 
}

func main () {
    // While I don't use proper routines I will just add some static variables that will only be obtain in the beggining of the program.
    ip := getIp()
    crypto := getCrypto()
    for {
	bar := fmt.Sprintf("| %s | %s | %s |", crypto, ip, getDate())
	exec.Command("xsetroot", "-name", bar).Run()
	time.Sleep(time.Second)
    }
}

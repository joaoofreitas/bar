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


// General Purpose
func createRoutine(f func() string, interval time.Duration, channel chan string) {
    defer close(channel)
    for {
    	channel <- f()
	time.Sleep(interval)
    }
}

// Date Gathering
func getDate() string {
    out, err := exec.Command("date").Output()
    if err != nil {
	return err.Error()
    }
    return fmt.Sprintf("%s", out);
}

func getBattery() string { 
    out, err := exec.Command("acpi").Output()
    if err != nil {
	return err.Error()
    }
    bat := fmt.Sprintf("%s", out);
    if strings.Contains(bat, " 0%") {
	return " - "
    }
    bat = strings.ReplaceAll(bat, ", rate information unavailable", "")
    return bat
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
    return "-" 
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
	return "-"
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
    var bar string

    var date string
    var ip string
    var bat string
    var crypto string 

    // Create channels for each routine
    dateChannel := make(chan string)
    networkChannel := make(chan string)
    powerChannel := make(chan string)
    cryptoChannel := make(chan string)

    // Calling the routines
    go createRoutine(getDate, 1 * time.Second, dateChannel)
    go createRoutine(getIp, 10 * time.Second, networkChannel)
    go createRoutine(getBattery, 2 * time.Minute, powerChannel)
    go createRoutine(getCrypto, time.Minute, cryptoChannel)

    for {
	select {
	    case msg := <- dateChannel:
		date = msg
	    case msg := <- networkChannel:
		ip = msg
	    case msg := <- powerChannel:
		bat = msg
	    case msg := <- cryptoChannel:
		crypto = msg
	    default:
		bar = fmt.Sprintf("| %s | %s | %s | %s |", crypto, bat, ip, date)
		exec.Command("xsetroot", "-name", bar).Run()
		time.Sleep(10 * time.Millisecond)
	}
    }
}

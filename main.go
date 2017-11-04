package main

import (
	"fmt"
	"github.com/getlantern/systray"
	"net/http"
	"encoding/json"
	"log"
	"time"
)

type Crypto struct {
	Id       string  `json:"id"`
	PriceBtc float64 `json:"price_btc"`
	PriceUsd float64 `json:"price_usd"`
	PriceEur float64 `json:"price_eur"`
}

func main() {
	systray.Run(onReady, onExit)
}

func updateTicker(coin string) Crypto {
	url := fmt.Sprintf(fmt.Sprintf("http://coincap.io/page/%s", coin))

	// Build request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println(err)
	}

	// Create a client
	client := &http.Client{}

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}

	// Close
	defer resp.Body.Close()

	var crypto Crypto

	// Decode
	if err := json.NewDecoder(resp.Body).Decode(&crypto); err != nil {
		log.Println(err)
	}

	return crypto
}

func setText(crypto Crypto) {
	systray.SetTitle(fmt.Sprintf("%.2f$", crypto.PriceUsd))
	systray.SetTooltip(fmt.Sprintf("%s is currently at %.2f$", crypto.Id, crypto.PriceUsd))
}

func onReady() {
	// Setup menu items
	menuBtc := systray.AddMenuItem("BTC", "Bitcoin")
	menuEth := systray.AddMenuItem("ETH", "Ethereum")
	menuLtc := systray.AddMenuItem("LTC", "Litecoin")
	menuNeo := systray.AddMenuItem("NEO", "NEO")
	systray.AddSeparator()
	menuQuit := systray.AddMenuItem("Quit", "Quit application")

	// Set btc as default
	var coin string = "BTC"
	menuBtc.Check()

	go func() {
		for {
			setText(updateTicker(coin))
			time.Sleep(1 * time.Second)
		}
	}()

	go func() {
		for {
			select {
			case <-menuBtc.ClickedCh:
				coin = "BTC"
				setText(updateTicker(coin))

				menuBtc.Check()
				menuEth.Uncheck()
				menuLtc.Uncheck()
				menuNeo.Uncheck()
			case <-menuEth.ClickedCh:
				coin = "ETH"
				setText(updateTicker(coin))

				menuEth.Check()
				menuBtc.Uncheck()
				menuLtc.Uncheck()
				menuNeo.Uncheck()
			case <-menuLtc.ClickedCh:
				coin = "LTC"
				setText(updateTicker(coin))

				menuLtc.Check()
				menuBtc.Uncheck()
				menuEth.Uncheck()
				menuNeo.Uncheck()
			case <-menuNeo.ClickedCh:
				coin = "NEO"
				setText(updateTicker(coin))

				menuNeo.Check()
				menuBtc.Uncheck()
				menuEth.Uncheck()
				menuLtc.Uncheck()
			case <-menuQuit.ClickedCh:
				systray.Quit()
				return
			}
		}
	}()
}

func onExit() {
	// Cleaning stuff here
}

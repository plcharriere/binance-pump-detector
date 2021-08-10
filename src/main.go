package main

import (
	"log"
	"strings"
	"time"

	"gopkg.in/ini.v1"
)

func main() {
	log.Print("Binance Pump Detector")

	app := NewApp()

	cfg, err := ini.Load(app.Args.ConfigFilePath)
	if err != nil {
		log.Print(err)
		return
	}

	user, err := app.CreateUser(cfg.Section("user").Key("apiKey").String(), cfg.Section("user").Key("secretKey").String())
	if err != nil {
		log.Panic(err)
		return
	}

	for _, section := range cfg.SectionStrings() {
		if section != "DEFAULT" && section != "user" {
			go func(section string) {
				symbolsStrings := strings.Split(section, "/")
				symbols := NewSymbols(symbolsStrings[0], symbolsStrings[1])
				pair, err := app.GetOrCreatePair(symbols)
				if err != nil {
					log.Print(err)
					return
				}

				percentChange, err := cfg.Section(section).Key("percentChange").Float64()
				if err != nil {
					log.Print(err)
					return
				}
				timeInterval, err := cfg.Section(section).Key("timeInterval").Int()
				if err != nil {
					log.Print(err)
					return
				}
				minimumTradeCount, err := cfg.Section(section).Key("minimumTradeCount").Int()
				if err != nil {
					log.Print(err)
					return
				}
				buyMarket, err := cfg.Section(section).Key("buyMarket").Bool()
				if err != nil {
					log.Print(err)
					return
				}
				buyQuantity, err := cfg.Section(section).Key("buyQuantity").Float64()
				if err != nil {
					log.Print(err)
					return
				}
				sellLimitPriceMultiplier, err := cfg.Section(section).Key("buyQuantity").Float64()
				if err != nil {
					log.Print(err)
					return
				}

				_, err = user.Watch(pair, Pump{
					PercentChange:            percentChange,
					TimeInterval:             time.Second * time.Duration(timeInterval),
					MinimumTradeCount:        minimumTradeCount,
					BuyMarket:                buyMarket,
					BuyQuantity:              buyQuantity,
					SellLimitPriceMultiplier: sellLimitPriceMultiplier,
				})
				if err != nil {
					log.Print(err)
				}
			}(section)
		}
	}

	<-app.Context.Done()
}

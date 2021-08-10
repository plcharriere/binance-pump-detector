package main

import (
	"log"
	"strings"
	"time"
)

var pairs = []string{
	"ETH/BTC",
}

func main() {
	log.Print("Binance Pump Detector")

	app := NewApp()

	user, err := app.CreateUser("", "")
	if err != nil {
		log.Panic(err)
		return
	}

	for _, symbols := range pairs {
		go func(symbols string) {
			ssymbols := strings.Split(symbols, "/")
			symbols2 := NewSymbols(ssymbols[0], ssymbols[1])
			pair, err := app.GetOrCreatePair(symbols2)
			if err != nil {
				log.Print(err)
				return
			}
			_, err = user.Watch(pair, Pump{
				PercentChange:            15,
				TimeInterval:             time.Second * 60,
				MinimumTradeCount:        1,
				BuyMarket:                false,
				BuyQuantity:              0.01,
				SellLimitPriceMultiplier: 1.0002,
			})
			if err != nil {
				log.Print(err)
			}
		}(symbols)
	}

	<-app.Context.Done()
}

package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/adshao/go-binance/v2"
)

type Watch struct {
	User  *User
	Pair  *Pair
	Pump  Pump
	subC  chan Trade
	stopC chan struct{}
}

func NewWatch(user *User, pair *Pair, pump Pump) (*Watch, error) {
	if pump.BuyMarket {
		balance, err := user.GetSymbolBalance(pair.Symbols.SecondSymbol)
		if err != nil {
			return nil, err
		}
		if balance < pump.BuyQuantity {
			return nil, fmt.Errorf("not enough %s to buy %s (%f/%f)", pair.Symbols.SecondSymbol, pair.Symbols.FirstSymbol, balance, pump.BuyQuantity)
		}
	}
	watch := &Watch{
		User:  user,
		Pair:  pair,
		Pump:  pump,
		subC:  make(chan Trade, 1000),
		stopC: make(chan struct{}, 1000),
	}
	go watch.goroutine()
	return watch, nil
}

func (watch *Watch) goroutine() {
	watch.Pair.SubC <- watch.subC
	log.Printf("Started goroutine on pair %s with pump conditions : %.2f%% change in %v with %d minimum trade(s)", watch.Pair.Symbols.ToStringSeparated(), watch.Pump.PercentChange, watch.Pump.TimeInterval, watch.Pump.MinimumTradeCount)
	for {
		select {
		case trade := <-watch.subC:
			change := watch.Pair.Change(watch.Pump.TimeInterval)
			timeDiff := time.Now().UnixNano()/int64(time.Millisecond) - trade.time
			if watch.Pair.App.Args.Verbose {
				log.Printf("%s got traded at %v (%vms ago): %.3f%% change in %v (%v -> %v) [%v -> %v] with %d trades", watch.Pair.Symbols.ToStringSeparated(), trade.price, timeDiff, change.Percent, change.Interval, change.Farthest.price, change.Nearest.price, change.From.Format("15:04:05"), change.To.Format("15:04:05"), change.TradeCount)
			}
			watch.Change(change)
		case <-watch.stopC:
			return
		case <-watch.Pair.WsStopC:
			return
		}
	}
}

func (watch *Watch) Change(change PairChange) {
	if change.Percent >= watch.Pump.PercentChange && change.TradeCount >= watch.Pump.MinimumTradeCount {
		if !watch.Pump.PercentChangeTriggered {
			log.Print("PUMP DETECTED")
			watch.Pump.PercentChangeTriggered = true
			if !watch.Pump.BuyOrderPlaced && watch.Pump.BuyMarket {
				log.Print("Placing MARKET buy order")
				buyOrder := watch.User.Client.NewCreateOrderService().
					Symbol(watch.Pair.Symbols.ToString()).
					Side(binance.SideTypeBuy).
					Type(binance.OrderTypeMarket).
					Quantity(fmt.Sprintf("%f", watch.Pump.BuyQuantity))
				buyOrderResponse, err := buyOrder.Do(context.Background())
				if err != nil {
					log.Print(err)
					return
				}
				log.Print("Bought ", buyOrderResponse.ExecutedQuantity)
				watch.Pump.BuyOrderPlaced = true

				log.Println(change.Nearest.price, fmt.Sprintf("%f", change.Nearest.price*watch.Pump.SellLimitPriceMultiplier))

				sellOrder := watch.User.Client.NewCreateOrderService().
					Symbol(watch.Pair.Symbols.ToString()).
					Side(binance.SideTypeSell).Type(binance.OrderTypeLimit).
					TimeInForce(binance.TimeInForceTypeGTC).
					Price(fmt.Sprintf("%f", change.Nearest.price*watch.Pump.SellLimitPriceMultiplier)).
					Quantity(buyOrderResponse.ExecutedQuantity)
				_, err = sellOrder.Do(context.Background())
				if err != nil {
					log.Print(err)
					return
				}
			}
		}
	}
}

func (watch *Watch) Stop() {
	watch.Pair.UnsubC <- watch.subC
	close(watch.stopC)
	close(watch.subC)
}

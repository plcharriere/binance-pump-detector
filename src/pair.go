package main

import (
	"fmt"
	"log"
	"math"
	"strconv"
	"sync"
	"time"

	"github.com/adshao/go-binance/v2"
)

type Trade struct {
	time  int64
	price float64
}

type Pair struct {
	App     *App
	Symbols Symbols
	trades  []Trade
	tradesM sync.RWMutex
	tradesC chan Trade
	subs    []chan Trade
	SubC    chan chan Trade
	UnsubC  chan chan Trade
	WsDoneC chan struct{}
	WsStopC chan struct{}
}

type Pairs []*Pair

type PairChange struct {
	From       time.Time
	To         time.Time
	Interval   time.Duration
	Percent    float64
	TradeCount int
	Nearest    Trade
	Farthest   Trade
}

func NewPair(app *App, symbols Symbols) (*Pair, error) {
	pair := &Pair{App: app}
	pair.tradesC = make(chan Trade, 1000)
	doneC, stopC, err := binance.WsAggTradeServe(symbols.ToString(), pair.AggTradeHandler, func(err error) {
		fmt.Println(err)
	})
	if err != nil {
		return nil, err
	}
	pair.Symbols = symbols
	pair.WsDoneC = doneC
	pair.WsStopC = stopC
	pair.SubC = make(chan chan Trade, 1000)
	pair.UnsubC = make(chan chan Trade, 1000)
	go pair.SubsGoroutine()
	return pair, nil
}

func (pair *Pair) AggTradeHandler(event *binance.WsAggTradeEvent) {
	floatPrice, err := strconv.ParseFloat(event.Price, 64)
	if err != nil {
		log.Println(err)
		return
	}
	aDayAgo := time.Now().UnixNano()/int64(time.Millisecond) - 86400000
	trade := Trade{event.Time, floatPrice}
	pair.tradesM.Lock()
	pair.trades = append(pair.trades, trade)
	for _, order := range pair.trades {
		if order.time < aDayAgo {
			pair.trades = pair.trades[1:]
		} else {
			break
		}
	}
	pair.tradesM.Unlock()
	pair.tradesC <- trade
}

func (pair *Pair) Stop() {
	close(pair.WsStopC)
	close(pair.SubC)
	close(pair.UnsubC)
	close(pair.tradesC)
}

func (pair *Pair) GetTrades() []Trade {
	pair.tradesM.RLock()
	trades := pair.trades
	pair.tradesM.RUnlock()
	return trades
}

func (pair *Pair) SubsGoroutine() {
	for {
		select {
		case trade := <-pair.tradesC:
			for _, sub := range pair.subs {
				sub <- trade
			}
		case sub := <-pair.SubC:
			pair.subs = append(pair.subs, sub)
		case unsub := <-pair.UnsubC:
			for i, sub := range pair.subs {
				if unsub == sub {
					pair.subs = sliceRemoveTradeChan(pair.subs, i)
					break
				}
			}
			if len(pair.subs) == 0 {
				pair.App.StopPair(pair)
				return
			}
		case <-pair.WsStopC:
			return
		}
	}
}

func (pair *Pair) Change(interval time.Duration) PairChange {
	trades := pair.GetTrades()
	totalTradeCount := len(trades)

	now := time.Now()
	from := now.Add(-interval)

	fromMili := from.UnixNano() / int64(time.Millisecond)

	tradeCount := 0
	farthest := Trade{}
	nearest := Trade{}

	for i := totalTradeCount - 1; i >= 0; i-- {
		trade := trades[i]

		if trade.time >= fromMili {
			if farthest.time == 0 || trade.time < farthest.time {
				farthest = trade
			}
			if nearest.time == 0 || trade.time > nearest.time {
				nearest = trade
			}
			tradeCount++
		} else {
			break
		}
	}

	percent := truncateFloat(math.Abs(farthest.price-nearest.price)/farthest.price*100, 0.001)
	if percent > 0 && nearest.price < farthest.price {
		percent *= -1
	}

	return PairChange{
		From:       from,
		To:         now,
		Interval:   interval,
		Percent:    percent,
		TradeCount: tradeCount,
		Nearest:    nearest,
		Farthest:   farthest,
	}
}

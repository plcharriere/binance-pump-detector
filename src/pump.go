package main

import "time"

type Pump struct {
	PercentChange            float64
	TimeInterval             time.Duration
	MinimumTradeCount        int
	PercentChangeTriggered   bool
	BuyMarket                bool
	BuyQuantity              float64
	BuyOrderPlaced           bool
	SellLimitPriceMultiplier float64
}

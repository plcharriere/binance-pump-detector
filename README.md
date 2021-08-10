# Binance Pump Detector

Binance Pump Detector is a program to monitor multiple coins on Binance and detect when a coin gets pumped with defined conditions.

It can open a Buy Market order when a pump is detected and open a Sell Limit order at a desired price after buying.

## Warnings

**DO NOT USE THIS PROGRAM IF YOU DO NOT KNOW EXACTLY WHAT YOU ARE DOING.**

I am not responsible for any loss.

## Configuration

Create a configuration file with :

`cp config.ini.example config.ini`

Configure your `config.ini` file with your Binance Api Key and Secret Key *(optional)* and coins you want to monitor.

```ini
[user]
apiKey    = ; Your Binance Api Key (optional)
secretKey = ; Your Binance Secret Key (optional)

[ETH/BTC]
percentChange             = 5       ; How much percent change in timeInterval
timeInterval              = 300     ; Time interval in seconds for percentChange
minimumTradeCount         = 500     ; Minimum trade count
buyMarket                 = false   ; Open Buy Market order when pump is detected
buyQuantity               = 0.01    ; How much quantity to buy
sellLimitPriceMultiplier  = 1.0002  ; Open Sell Limit order at (price pump is detected) * sellLimitPriceMultiplier
```

If Binance keys are not defined, the program will still monitor coins but won't be able to open orders.

## Run

`make run`

Runs without building a binary with `-race` flag.

## Build

`make`

or

`make build`

Builds a binary `binance-pump-detector` in `bin/`.
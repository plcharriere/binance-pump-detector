package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/adshao/go-binance/v2"
)

type User struct {
	Client  *binance.Client
	watches []*Watch
}

type Users []*User

func NewUser(apiKey, secretKey string) (*User, error) {
	client := binance.NewClient(apiKey, secretKey)
	return &User{Client: client}, nil
}

func (user *User) Watch(pair *Pair, pump Pump) (*Watch, error) {
	if user.GetWatch(pair) != nil {
		return nil, fmt.Errorf("already watching pair '%s'", pair.Symbols.ToStringSeparated())
	}
	watch, err := NewWatch(user, pair, pump)
	if err != nil {
		return nil, err
	}
	user.watches = append(user.watches, watch)
	return watch, nil
}

func (user *User) GetWatch(pair *Pair) *Watch {
	for _, watch := range user.watches {
		if watch.Pair == pair {
			return watch
		}
	}
	return nil
}

func (user *User) Unwatch(pair *Pair) error {
	for i, watch := range user.watches {
		if watch.Pair == pair {
			watch.Stop()
			user.watches = sliceRemoveWatch(user.watches, i)
			return nil
		}
	}
	return fmt.Errorf("not watching pair '%s'", pair.Symbols.ToStringSeparated())
}

func (user *User) GetSymbolBalance(symbol string) (float64, error) {
	account, err := user.Client.NewGetAccountService().Do(context.Background())
	if err != nil {
		return 0, err
	}
	for _, balance := range account.Balances {
		if balance.Asset == symbol {
			amount, err := strconv.ParseFloat(balance.Free, 64)
			return amount, err
		}
	}
	return 0, fmt.Errorf("asset not found")
}

package main

import (
	"context"
	"sync"
)

type App struct {
	Args Args

	Pairs  Pairs
	PairsM sync.RWMutex

	Users  Users
	UsersM sync.RWMutex

	Context context.Context
}

func NewApp() *App {
	var args Args
	args.Parse()

	return &App{
		Args:    args,
		Context: context.Context(context.Background()),
	}
}

func (app *App) GetPair(symbols Symbols) *Pair {
	app.PairsM.RLock()
	for _, pair := range app.Pairs {
		if pair.Symbols == symbols {
			return pair
		}
	}
	app.PairsM.RUnlock()
	return nil
}

func (app *App) GetOrCreatePair(symbols Symbols) (*Pair, error) {
	pair := app.GetPair(symbols)
	if pair == nil {
		newPair, err := NewPair(app, symbols)
		if err != nil {
			return nil, err
		}
		app.PairsM.Lock()
		app.Pairs = append(app.Pairs, newPair)
		app.PairsM.Unlock()
		pair = newPair
	}
	return pair, nil
}

func (app *App) CreateUser(apiKey, secretKey string) (*User, error) {
	user, err := NewUser(apiKey, secretKey)
	if err != nil {
		return nil, err
	}
	app.Users = append(app.Users, user)
	return user, nil
}

func (app *App) StopPair(pair *Pair) {
	pair.Stop()
	for i, pair2 := range app.Pairs {
		if pair == pair2 {
			app.Pairs = sliceRemovePair(app.Pairs, i)
			break
		}
	}
}

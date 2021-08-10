package main

import (
	"math/big"
)

func sliceRemoveTradeChan(s []chan Trade, i int) []chan Trade {
	s[len(s)-1], s[i] = s[i], s[len(s)-1]
	return s[:len(s)-1]
}

func sliceRemoveWatch(s []*Watch, i int) []*Watch {
	s[len(s)-1], s[i] = s[i], s[len(s)-1]
	return s[:len(s)-1]
}

func sliceRemovePair(s []*Pair, i int) []*Pair {
	s[len(s)-1], s[i] = s[i], s[len(s)-1]
	return s[:len(s)-1]
}

func truncateFloat(f, unit float64) float64 {
	bf := big.NewFloat(0).SetPrec(1000).SetFloat64(f)
	bu := big.NewFloat(0).SetPrec(1000).SetFloat64(unit)

	bf.Quo(bf, bu)

	i := big.NewInt(0)
	bf.Int(i)
	bf.SetInt(i)

	f, _ = bf.Mul(bf, bu).Float64()
	return f
}

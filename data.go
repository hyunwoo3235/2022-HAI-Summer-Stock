package main

import (
	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/stat"
)

func getInput(prices []float64) [35]float64 {
	var input [35]float64

	lastIndex := len(prices) - 1
	windows := [5]int{5, 10, 20, 30, 60}
	for i := 0; i < len(windows); i++ {
		input[i] = prices[lastIndex-windows[i]] / prices[lastIndex]
		input[i+5] = stat.Mean(prices[lastIndex-windows[i]:lastIndex], nil)
		input[i+10] = stat.StdDev(prices[lastIndex-windows[i]:lastIndex], nil)
		input[i+15] = floats.Max(prices[lastIndex-windows[i] : lastIndex])
		input[i+20] = floats.Min(prices[lastIndex-windows[i] : lastIndex])

		c := 0
		for j := lastIndex - windows[i]; j < lastIndex; j++ {
			if prices[j] > prices[j-1] {
				c++
			}
		}
		input[i+25] = float64(c) / float64(windows[i])

		c = 0
		for j := lastIndex - windows[i]; j < lastIndex; j++ {
			if prices[j] < prices[j-1] {
				c++
			}
		}
		input[i+30] = float64(c) / float64(windows[i])

	}

	return input
}

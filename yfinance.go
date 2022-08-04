package main

import (
	"github.com/imroc/req/v3"
)

type ChartQueryResult struct {
	Chart struct {
		Result []struct {
			Meta struct {
				Currency             string  `json:"currency"`
				Symbol               string  `json:"symbol"`
				ExchangeName         string  `json:"exchangeName"`
				InstrumentType       string  `json:"instrumentType"`
				FirstTradeDate       int     `json:"firstTradeDate"`
				RegularMarketTime    int     `json:"regularMarketTime"`
				Gmtoffset            int     `json:"gmtoffset"`
				Timezone             string  `json:"timezone"`
				ExchangeTimezoneName string  `json:"exchangeTimezoneName"`
				RegularMarketPrice   float64 `json:"regularMarketPrice"`
				ChartPreviousClose   float64 `json:"chartPreviousClose"`
				PreviousClose        float64 `json:"previousClose"`
				Scale                int     `json:"scale"`
				PriceHint            int     `json:"priceHint"`
				CurrentTradingPeriod struct {
					Pre struct {
						Timezone  string `json:"timezone"`
						Start     int    `json:"start"`
						End       int    `json:"end"`
						Gmtoffset int    `json:"gmtoffset"`
					} `json:"pre"`
					Regular struct {
						Timezone  string `json:"timezone"`
						Start     int    `json:"start"`
						End       int    `json:"end"`
						Gmtoffset int    `json:"gmtoffset"`
					} `json:"regular"`
					Post struct {
						Timezone  string `json:"timezone"`
						Start     int    `json:"start"`
						End       int    `json:"end"`
						Gmtoffset int    `json:"gmtoffset"`
					} `json:"post"`
				} `json:"currentTradingPeriod"`
				TradingPeriods [][]struct {
					Timezone  string `json:"timezone"`
					Start     int    `json:"start"`
					End       int    `json:"end"`
					Gmtoffset int    `json:"gmtoffset"`
				} `json:"tradingPeriods"`
				DataGranularity string   `json:"dataGranularity"`
				Range           string   `json:"range"`
				ValidRanges     []string `json:"validRanges"`
			} `json:"meta"`
			Timestamp  []int `json:"timestamp"`
			Indicators struct {
				Quote []struct {
					Open   []float64 `json:"open"`
					Volume []int     `json:"volume"`
					Close  []float64 `json:"close"`
					High   []float64 `json:"high"`
					Low    []float64 `json:"low"`
				} `json:"quote"`
			} `json:"indicators"`
		} `json:"result"`
		Error interface{} `json:"error"`
	} `json:"chart"`
}

type Ticker struct {
	Symbol string `json:"symbol"`
	client *req.Client
}

func NewTicker(symbol string) *Ticker {
	return &Ticker{
		Symbol: symbol,
		client: req.C(),
	}
}

func (t *Ticker) GetChart(rangeStr, intervalStr string) (*ChartQueryResult, error) {
	var r ChartQueryResult
	url := "https://query1.finance.yahoo.com/v7/finance/chart/" + t.Symbol
	_, err := t.client.R().
		SetPathParam("range", rangeStr).
		SetPathParam("interval", intervalStr).
		SetResult(&r).
		Get(url)

	if err != nil {
		return nil, err
	}

	return &r, nil
}

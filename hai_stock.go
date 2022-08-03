package main

import (
	"github.com/imroc/req/v3"
)

type HAIStock struct {
	host   string
	token  string
	client *req.Client
}

type AccountInfo struct {
	Token   string `json:"token"`
	Name    string `json:"name"`
	Deposit int    `json:"deposit"`
	Stocks  map[string]struct {
		Ticker string `json:"ticker"`
		Price  int    `json:"price"`
		Share  int    `json:"share"`
	}
}

func NewHAIStock(host, token string) *HAIStock {
	return &HAIStock{
		host:   host,
		token:  token,
		client: req.C().SetCommonHeader("token", token),
	}
}

func (h *HAIStock) SendOrder(orderType string, ticker string, price int, share int) (int, error) {
	var r int
	_, err := h.client.R().
		SetBody(map[string]any{
			"ticker": ticker,
			"price":  price,
			"share":  share,
		}).
		SetResult(&r).
		Post(h.host + "/" + orderType)

	if err != nil {
		return 0, err
	}

	return r, nil
}

func (h *HAIStock) dropOrder(orderId int) (string, error) {
	var r string
	_, err := h.client.R().
		SetBody(map[string]any{
			"order_id": orderId,
		}).
		SetResult(&r).
		Post(h.host + "/drop_order")

	if err != nil {
		return "", err
	}

	return r, nil
}

func (h *HAIStock) AccountInfo() (AccountInfo, error) {
	var r AccountInfo
	_, err := h.client.R().
		SetResult(&r).
		Get(h.host + "/account")

	if err != nil {
		return AccountInfo{}, err
	}

	return r, nil
}

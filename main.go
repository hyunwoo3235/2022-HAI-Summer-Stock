package main

import (
	"github.com/dmitryikh/leaves"
	"github.com/go-co-op/gocron"
	"log"
	"os"
	"time"
)

var (
	token    = os.Getenv("token")
	host     = os.Getenv("host")
	haistock = NewHAIStock(host, token)
	model, _ = leaves.LGEnsembleFromFile("models/lgb_model.txt", true)
	norm     = ZScoreNorm{mean: 331.84802192774686, std: 36.820766125382214}
)

var (
	deposit = 50000000
	stock   = ""
)

func main() {
	s := gocron.NewScheduler(time.FixedZone("America/New_York", -4*60*60))

	s.Cron("59 9-15 * * *").Do(task)
	s.StartAsync()
}

func task() {
	res, err := NewTicker("TQQQ").GetPrices("60h", "1h")
	if err != nil {
		log.Println(err)
		return
	}

	for i := 0; i < len(res); i++ {
		res[i] = norm.Normalize(res[i])
	}
	input := getInput(res)

	last := res[len(res)-1]
	pred := model.PredictSingle(input[:], 0)

	last = norm.Denormalize(last)
	pred = norm.Denormalize(pred)

	log.Printf("Last: %f, Pred: %f", last, pred)

	var tobuy, tosell string
	if pred > last {
		tobuy, tosell = "TQQQ", "SQQQ"
	} else {
		tobuy, tosell = "SQQQ", "TQQQ"
	}

	if stock == tobuy {
		return
	}
	stock = tobuy

	ainfo, _ := haistock.AccountInfo()

	_, _ = haistock.SendOrder("sell", tosell, 1, ainfo.Stocks[tosell].Share)
	deposit += ainfo.Stocks[tosell].Share * ainfo.Stocks[tosell].Price

	share := deposit / (int(last) * 100)
	share = share - share%100

	_, _ = haistock.SendOrder("buy", tobuy, 1000000, share)

	return
}

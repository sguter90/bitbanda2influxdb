package main

import (
	"errors"
	_ "github.com/influxdata/influxdb1-client" // this is important because of the bug in go mod
	influxDb "github.com/influxdata/influxdb1-client/v2"
	bitpanda "github.com/sguter90/bitpanda-public-api-client"
	"strconv"
	"time"
)

type Bitpanda2InfluxDbConfig struct {
	BitpandaApiKey                 string
	InfluxDbAddr                   string
	InfluxDbUser                   string
	InfluxDbPass                   string
	InfluxDbBatchPointConfig       influxDb.BatchPointsConfig
	InfluxDbPointTicker            string
	InfluxDbPointWalletsEurBalance string
}

func NewBitpanda2InfluxDb(c Bitpanda2InfluxDbConfig) (p Bitpanda2InfluxDb, err error) {
	influxDbClient, err := influxDb.NewHTTPClient(influxDb.HTTPConfig{
		Addr:     c.InfluxDbAddr,
		Username: c.InfluxDbUser,
		Password: c.InfluxDbPass,
	})
	if err != nil {
		err = errors.New("Error creating InfluxDB Client: " + err.Error())
		return
	}

	bitpandaConfig := bitpanda.NewConfig(c.BitpandaApiKey)
	p = Bitpanda2InfluxDb{
		Bitpanda: *bitpanda.NewClient(bitpandaConfig),
		InfluxDb: influxDbClient,
		Config:   c,
	}

	return
}

type Bitpanda2InfluxDb struct {
	Bitpanda bitpanda.Client
	InfluxDb influxDb.Client
	Config   Bitpanda2InfluxDbConfig
}

func (c *Bitpanda2InfluxDb) push2InfluxDb(pointName string, values map[string]interface{}) (err error) {
	// Create a new point batch
	bp, _ := influxDb.NewBatchPoints(c.Config.InfluxDbBatchPointConfig)

	// Create a point and add to batch
	pt, err := influxDb.NewPoint(pointName, nil, values, time.Now())
	if err != nil {
		err = errors.New("Failed to create point: " + err.Error())
		return
	}
	bp.AddPoint(pt)

	// Write the batch
	err = c.InfluxDb.Write(bp)
	if err != nil {
		err = errors.New("Failed to write batch: " + err.Error())
		return
	}

	return
}

func (c *Bitpanda2InfluxDb) GetWalletsBalance() (symbolSums map[string]float64, err error) {
	symbolSums = map[string]float64{}

	resp, err := c.Bitpanda.WalletsGet()
	if err != nil {
		return
	}

	for _, wallet := range resp.Data {
		wSymbol := wallet.Attributes.CryptocoinSymbol
		wBalance, parseErr := strconv.ParseFloat(wallet.Attributes.Balance, 64)
		if parseErr != nil {
			err = parseErr
			return
		}

		balance, exists := symbolSums[wSymbol]
		if exists == false {
			balance = 0.0
		}

		symbolSums[wSymbol] = balance + wBalance
	}

	return
}

func (c *Bitpanda2InfluxDb) PushCoinTicker() (err error) {
	coins, err := c.Bitpanda.TickerGet()
	if err != nil {
		return
	}

	fields := map[string]interface{}{}
	for _, coin := range coins {
		value, _ := strconv.ParseFloat(coin.EUR, 64)
		fields[coin.Name] = value
	}

	err = c.push2InfluxDb(c.Config.InfluxDbPointTicker, fields)
	return
}

func (c *Bitpanda2InfluxDb) PushWalletsEur() (err error) {
	wBalances, err := c.GetWalletsBalance()
	if err != nil {
		return
	}

	coins, err := c.Bitpanda.TickerGet()
	if err != nil {
		return
	}

	eurValues := map[string]interface{}{}
	for coinSymbol, balance := range wBalances {
		for _, coin := range coins {
			if coinSymbol != coin.Name {
				continue
			}

			coinPrice, parseErr := strconv.ParseFloat(coin.EUR, 64)
			if parseErr != nil {
				err = parseErr
				return
			}
			eurValues[coin.Name] = balance * coinPrice
		}
	}

	err = c.push2InfluxDb(c.Config.InfluxDbPointWalletsEurBalance, eurValues)
	return
}

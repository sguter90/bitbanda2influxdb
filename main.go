package main

import (
	"fmt"
	influxDb "github.com/influxdata/influxdb1-client/v2"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	c := Bitpanda2InfluxDbConfig{
		BitpandaApiKey: os.Getenv("BITPANDA_API_KEY"),
		InfluxDbAddr:   os.Getenv("INFLUXDB_ADDRESS"),
		InfluxDbUser:   os.Getenv("INFLUXDB_USER"),
		InfluxDbPass:   os.Getenv("INFLUXDB_PASS"),
		InfluxDbBatchPointConfig: influxDb.BatchPointsConfig{
			Database:  os.Getenv("INFLUXDB_DATABASE"),
			Precision: "s",
		},
		InfluxDbPointTicker:            os.Getenv("INFLUXDB_POINT_TICKER"),
		InfluxDbPointWalletsEurBalance: os.Getenv("INFLUXDB_POINT_WALLETS_EUR"),
	}

	p, err := NewBitpanda2InfluxDb(c)
	if err != nil {
		fmt.Println(err)
	}

	err = p.PushCoinTicker()
	if err != nil {
		fmt.Println(err)
	}

	err = p.PushWalletsEur()
	if err != nil {
		fmt.Println(err)
	}
}

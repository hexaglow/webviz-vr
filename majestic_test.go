package main

import (
	"testing"
	"os"
	"log"
)

var testApi MajesticApi

func TestMain(m *testing.M) {
	key := os.Getenv("MAJESTIC_API_KEY")
	if key == "" {
		log.Fatalln("MAJESTIC_API_KEY missing from environment.")
	}
	testApi = MajesticApi{apiKey:key}

	os.Exit(m.Run())
}

func TestGetBackLinkData(t *testing.T) {
	request := testApi.GetBackLinkData()
	request.Item = "majestic.com"
	request.Count = 5
	request.Datasource = Fresh

	if _, err := request.Perform(); err != nil {
		t.Fatal(err.Error())
	}
}

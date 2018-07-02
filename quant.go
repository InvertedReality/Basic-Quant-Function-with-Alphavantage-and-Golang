package main

import (
	"encoding/json"
	"fmt"
	// structs "github.com/fatih/structs"
	"io/ioutil"
	"net/http"
	"sort"
)

type AutoGenerated struct {
	MetaData struct {
		OneInformation     string `json:"1. Information"`
		TwoSymbol          string `json:"2. Symbol"`
		ThreeLastRefreshed string `json:"3. Last Refreshed"`
		FourInterval       string `json:"4. Interval"`
		FiveOutputSize     string `json:"5. Output Size"`
		SixTimeZone        string `json:"6. Time Zone"`
	} `json:"Meta Data"`

	TimeSeries1Min map[string]interface{} `json:"Time Series (1min)"`
}

func (f AutoGenerated) TimeSeries1() (map[string]interface{}, error) {

	return f.TimeSeries1Min, nil
}
func (f AutoGenerated) Meta() (interface{}, error) {

	return f.MetaData, nil
}

// var slice map[string]interface{}

type Timedata struct {
	OneOpen    string `json:"1. open"`
	TwoHigh    string `json:"2. high"`
	ThreeLow   string `json:"3. low"`
	FourClose  string `json:"4. close"`
	FiveVolume string `json:"5. volume"`
}

const (
	empty = ""
	tab   = "\t"
)

// /get/data/
//returns data from computations on received data from alphaadvantage
func getdata(writer http.ResponseWriter, request *http.Request) {

	auto := AutoGenerated{}
	var keys []string
	slice := make(map[string]interface{})
	response, err := http.Get("https://www.alphavantage.co/query?function=TIME_SERIES_INTRADAY&symbol=MSFT&interval=1min&apikey=KLMH2VFJ0LCFNOX5")
	if err != nil {
		fmt.Println("Error getting data from alphavantage")
	}
	responsedata, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error getting data from alphavantage")
	}
	response.Body.Close()
	json.Unmarshal(responsedata, &auto)
	times, err := auto.TimeSeries1()
	//Get the keys of the latest time series data in a sorted fashion
	if err != nil {
		fmt.Println("Error")
	} else {
		for k := range times {
			keys = append(keys, k)
		}

		sort.Strings(keys)
		// get the keys of the five most recent data
		latest_keys := keys[(len(keys) - 5):(len(keys))]
		//Add the keys to a map to retain the latest data in a single variable
		{
			for _, key := range latest_keys {
				slice[key] = times[key]
			}
		}
	}

	if err != nil {
		{
			writer.Header().Set("Content-Type", "application/json")
			writer.WriteHeader(http.StatusInternalServerError)
			encoder := json.NewEncoder(writer)
			encoder.SetIndent(empty, tab)
			fmt.Println(http.StatusInternalServerError)
			fmt.Println(err)
		}
	} else {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)
		encoder := json.NewEncoder(writer)
		encoder.SetIndent(empty, tab)
		encoder.Encode(slice)
		fmt.Println("/get", http.StatusOK)
	}
}

//home
//lists urls
func Home(writer http.ResponseWriter, request *http.Request) {
	urls := map[int]string{
		1: "/get/data",
		2: "/get/meta",
	}

	{
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)
		encoder := json.NewEncoder(writer)
		encoder.SetIndent(empty, tab)
		encoder.Encode(urls)

	}

}

//get meta data
//get/meta/
func GetMetaData(writer http.ResponseWriter, request *http.Request) {
	auto := AutoGenerated{}
	response, err := http.Get("https://www.alphavantage.co/query?function=TIME_SERIES_INTRADAY&symbol=MSFT&interval=1min&apikey=KLMH2VFJ0LCFNOX5")
	if err != nil {
		fmt.Println("Error getting data from alphavantage")
	}
	responsedata, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error getting data from alphavantage")
	}
	response.Body.Close()
	json.Unmarshal(responsedata, &auto)
	meta, err := auto.Meta()
	if err != nil {
		{
			writer.Header().Set("Content-Type", "application/json")
			writer.WriteHeader(http.StatusInternalServerError)
			encoder := json.NewEncoder(writer)
			encoder.SetIndent(empty, tab)
			fmt.Println(http.StatusInternalServerError)
			fmt.Println(err)
		}
	} else {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)
		encoder := json.NewEncoder(writer)
		encoder.SetIndent(empty, tab)
		encoder.Encode(meta)
		fmt.Println("/get/meta", http.StatusOK)
	}
}
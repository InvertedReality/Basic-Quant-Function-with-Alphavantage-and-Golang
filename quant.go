package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"sort"
	"strconv"
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

const (
	empty = ""
	tab   = "\t"
)

// /get/data/
//returns latest data from alphaadvantage
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
	// fmt.Println(auto)
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
		fmt.Println(request.URL.Path, http.StatusOK)
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
			fmt.Println(http.StatusInternalServerError)
			fmt.Println(err)
		}
	} else {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)
		encoder := json.NewEncoder(writer)
		encoder.SetIndent(empty, tab)
		encoder.Encode(meta)
		fmt.Println(request.URL.Path, http.StatusOK)
	}
}

//function for calculating the desired values and returning the co ordinates for plotting the graph
// /get/graph

func getgraph(writer http.ResponseWriter, request *http.Request) {
	auto := AutoGenerated{}
	var keys []string
	average_data := make(map[string]interface{})
	std_dev_data := make([]interface{}, 0)
	graph_data := make(map[string]interface{})
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
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusInternalServerError)
		encoder := json.NewEncoder(writer)
		encoder.SetIndent(empty, tab)
		fmt.Println(http.StatusInternalServerError)
		fmt.Println(err)
	} else {
		for k := range times {
			keys = append(keys, k)
		}

		sort.Strings(keys)
		// get the keys of the five most recent data
		latest_keys := keys[(len(keys) - 5):(len(keys))]
		//loop to get value of ech of the keys and calculate the
		//required values for the client side graph
		{
			opening_float := make([]float64, 0)
			closing_float := make([]float64, 0)
			high_float := make([]float64, 0)
			low_float := make([]float64, 0)
			volume_float := make([]float64, 0)
			x_vs_y := make(map[string][]interface{})
			for _, key := range latest_keys {
				value := times[key]
				y_values := make([]interface{}, 0)
				if value, ok := value.(map[string]interface{}); ok {
					opening_value := value["1. open"]
					open_map := make(map[string]float64)
					if opening_value, ok := opening_value.(string); ok {
						parsed_value, err := strconv.ParseFloat(opening_value, 64)
						if err == nil {
							opening_float = append(opening_float, parsed_value)
							open_map["open"] = parsed_value
							y_values = append(y_values, open_map)
						}
					}
					high_value := value["2. high"]
					high_map := make(map[string]float64)
					if high_value, ok := high_value.(string); ok {
						parsed_value, err := strconv.ParseFloat(high_value, 64)
						if err == nil {
							high_float = append(high_float, parsed_value)
							high_map["high"] = parsed_value
							y_values = append(y_values, high_map)
						}
					}
					low_value := value["3. low"]
					low_map := make(map[string]float64)
					if low_value, ok := low_value.(string); ok {
						parsed_value, err := strconv.ParseFloat(low_value, 64)
						if err == nil {
							low_float = append(low_float, parsed_value)
							low_map["low"] = parsed_value
							y_values = append(y_values, low_map)
						}
					}
					closing_value := value["4. close"]
					close_map := make(map[string]float64)
					if closing_value, ok := closing_value.(string); ok {
						parsed_value, err := strconv.ParseFloat(closing_value, 64)
						if err == nil {
							closing_float = append(closing_float, parsed_value)
							close_map["close"] = parsed_value
							y_values = append(y_values, close_map)
						}
					}
					volume := value["5. volume"]
					volume_map := make(map[string]float64)
					if volume, ok := volume.(string); ok {
						parsed_value, err := strconv.ParseFloat(volume, 64)
						if err == nil {
							volume_float = append(volume_float, parsed_value)
							volume_map["volume"] = parsed_value
							y_values = append(y_values, volume_map)
						}
					}
				} else {
					fmt.Println("errors")
				}
				x_vs_y[key] = y_values
				graph_data["co_ordinates"] = x_vs_y
			}

			//calculate the averages and append to our map
			average_data["average_volume"] = average(volume_float)
			average_data["average_open"] = average(opening_float)
			average_data["average_close"] = average(closing_float)
			average_data["average_low"] = average(low_float)
			average_data["average_high"] = average(high_float)
			graph_data["averages"] = average_data
			//Calculate the variance

			//Step 1: Find the mean.
			// Step 2: For each data point, find the square of its distance to the mean.
			// Step 3: Sum the values from Step 2.
			// Step 4: Divide by the number of data points.
			// Step 5: Take the square root.

			//for the volume
			{
				var sum float64 = 0
				dev_volume := make(map[string]float64)
				for _, volume := range volume_float {
					square := math.Pow(volume, 2.00) / 100
					sum += square
				}
				std_deviation := math.Pow((sum/float64(len(volume_float))), 0.50) / 100
				dev_volume["volume"] = math.Round(std_deviation*10000) / 10000
				std_dev_data = append(std_dev_data, dev_volume)
			}
			//for the open
			{
				var sum float64 = 0
				dev_open := make(map[string]float64)
				for _, opening := range opening_float {
					square := math.Pow(opening, 2.00) / 100
					sum += square
				}
				std_deviation := math.Pow((sum/float64(len(opening_float))), 0.50) / 100
				dev_open["open"] = math.Round(std_deviation*10000) / 10000
				std_dev_data = append(std_dev_data, dev_open)
			}
			//for the close
			{
				var sum float64 = 0
				dev_close := make(map[string]float64)
				for _, closing := range closing_float {
					square := math.Pow(closing, 2.00) / 100
					sum += square
				}
				std_deviation := math.Pow((sum/float64(len(closing_float))), 0.50) / 100
				dev_close["close"] = math.Round(std_deviation*10000) / 10000
				std_dev_data = append(std_dev_data, dev_close)
			}
			//for the high
			{
				var sum float64 = 0
				dev_high := make(map[string]float64)
				for _, high := range high_float {
					square := math.Pow(high, 2.00) / 100
					sum += square
				}
				std_deviation := math.Pow((sum/float64(len(high_float))), 0.50) / 100
				dev_high["high"] = math.Round(std_deviation*10000) / 10000
				std_dev_data = append(std_dev_data, dev_high)
			}
			//for the low
			{
				var sum float64 = 0
				dev_low := make(map[string]float64)
				for _, low := range low_float {
					square := math.Pow(low, 2.00) / 100
					sum += square
				}
				std_deviation := math.Pow((sum/float64(len(low_float))), 0.50) / 100
				dev_low["low"] = math.Round(std_deviation*10000) / 10000
				std_dev_data = append(std_dev_data, dev_low)
			}
			//Appending the data
			graph_data["market_volatility"] = std_dev_data

			//serving the data as json
			writer.Header().Set("Content-Type", "application/json")
			writer.WriteHeader(http.StatusOK)
			encoder := json.NewEncoder(writer)
			encoder.SetIndent(empty, tab)
			encoder.Encode(graph_data)
			fmt.Println(request.URL.Path, http.StatusOK)
		}
	}

}

//returns the average
func average(floats []float64) float64 {
	var total float64 = 0
	for _, value := range floats {

		total += value
	}
	average_computation := math.Round((total/float64(len(floats)))*100) / 100
	return average_computation

}

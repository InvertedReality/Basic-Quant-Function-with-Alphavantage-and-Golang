package main

import (
	"encoding/json"
	"net/http"
	// "time"
	"fmt"
	// "io/ioutil"
)

func main() {
	http.HandleFunc("/get/", getdata)
	http.ListenAndServe("0.0.0.0:8080", nil)
}

func getdata(writer http.ResponseWriter, request *http.Request) {

	resp, err := http.Get("https://www.alphavantage.co/query?function=TIME_SERIES_INTRADAY&symbol=MSFT&interval=1min&apikey=KLMH2VFJ0LCFNOX5")
	dec := json.NewDecoder(resp.Body)
	if dec == nil {
		panic("Failed to start decoding JSON data")
	}

	json_map := make(map[string]interface{})

	err = dec.Decode(&json_map)
	for k, v := range json_map {
		if k == "Meta Data" {
			continue
		}
		enc := json.NewEncoder(os.Stdout)
		new_data := enc.Encode(v)
		// for k_2,v_2 := range v{
		//
		// }
		//fmt.Printf("key[%s] value[%s]\n", k, v)
		if new_data != nil {
			fmt.Println(new_data.Error())
			for k_2, v_2 := range new_data {
				fmt.Printf("key[%s]", k_2)
			}
		}
	}
	if err != nil {
		panic(err)
	}

	d.Data["json"] = &json_map
	d.ServeJSON()

}

package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/tealeg/xlsx"
	"net/http"
	"os"
)

type WeatherInfo struct {
	Status 		string	`json:"status"`
	Count 		string	`json:"count"`
	Info 		string	`json:"info"`
	Infocode 	string	`json:"infocode"`
	Lives 		[]*Live	`json:"lives"`
}

type Live struct {
	Province 		string		`json:"province"`
	City 			string		`json:"city"`
	Adcode 			string		`json:"adcode"`
	Weather 		string		`json:"weather"`
	Temperature 	string		`json:"temperature"`
	Winddirection 	string		`json:"winddirection"`
	Windpower 		string		`json:"windpower"`
	Humidity 		string		`json:"humidity"`
	Reporttime 		string		`json:"reporttime"`
}

const (
	MapUrl = "https://restapi.amap.com/v3/weather/weatherInfo?"
	MyKey = "ba43d9204d56947664c7e17bad1995f7"
	InputHelp = "Please input city name (eg. 北京市，注意有\"市\"；也可以输入区，eg. 海淀区): "
)

func main()  {

	cityCodeList := ParseXlsx("queryWeather/AMap_adcode_citycode_2020_4_10.xlsx")
	input := bufio.NewScanner(os.Stdin)
	fmt.Println(InputHelp)
	for ; input.Scan();  fmt.Printf("\n%s\n", InputHelp) {
		cityCode, ok := cityCodeList[input.Text()]
		if !ok {
			fmt.Printf("Not include %s\n", input.Text())
			continue
		}

		weather := QueryWeather(cityCode)
		if weather == nil {
			fmt.Printf("Query %s weather failed\n", input.Text())
			continue
		}

		for _, item := range weather.Lives {
			fmt.Printf("province: %s\tcity: %s,\tweather: %s,\ttemperature: %s\n",
				 item.Province, item.City, item.Weather, item.Temperature)
		}
	}
}

func ParseXlsx(filename string) map[string]string {
	fd, err := xlsx.OpenFile(filename)
	if err != nil {
		fmt.Printf("xlsx.OpenFile failed --> %s\n", err)
		return nil
	}

	var cityCode = make(map[string]string)
	for i := 1; i < fd.Sheets[0].MaxRow; i++ {
		row, _ := fd.Sheets[0].Row(i)
		cityCode[row.GetCell(0).Value] = row.GetCell(1).Value
	}

	return cityCode
}

func QueryWeather(cityCode string) *WeatherInfo  {
	queryPara := fmt.Sprintf("city=%s&key=%s", cityCode, MyKey)
	queryUrl := MapUrl + queryPara
	
	reps, err := http.Get(queryUrl)
	if err != nil {
		fmt.Printf("http.Get failed --> %s\n", err)
		return nil
	}

	if reps.StatusCode != http.StatusOK {
		fmt.Printf("reps.StatusCode is not ok, num --> %d\n", reps.StatusCode)
		reps.Body.Close()
		return nil
	}

	var weatherInfo WeatherInfo
	err = json.NewDecoder(reps.Body).Decode(&weatherInfo)
	if err != nil {
		fmt.Printf("json.NewDecoder(reps.Body).Decode failed --> %s\n", err)
		reps.Body.Close()
		return nil
	}

	reps.Body.Close()
	return &weatherInfo
}
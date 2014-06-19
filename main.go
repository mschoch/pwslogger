package main

import (
	//"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/dustin/go-jsonpointer"
)

var baseUrl = flag.String("base", "http://weatherstation.wunderground.com/weatherstation/updateweatherstation.php", "base url")
var id = flag.String("id", "", "device id")
var passwd = flag.String("passwd", "", "password")
var tempCPath = flag.String("tempCPath", "", "path to temp in celcius")
var humidityPath = flag.String("humidityPath", "", "path to humidity")
var pressurePath = flag.String("pressurePath", "", "path to pressure value in inches mercury")

func main() {

	flag.Parse()

	if *id == "" {
		log.Fatal("id is required")
	}
	if *passwd == "" {
		log.Fatal("password is required")
	}

	inputBytes, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	urlValues := url.Values{
		"action":       []string{"updateraw"},
		"softwaretype": []string{"pwslogger"},
		"dateutc":      []string{time.Now().UTC().Format("2006-01-02 15:04:05")},
		"ID":           []string{*id},
		"PASSWORD":     []string{*passwd},
	}

	var tempCFloat, humidityFloat float64
	var haveTemp, haveHumidity bool

	if *tempCPath != "" {
		tempCValue, err := jsonpointer.Find(inputBytes, *tempCPath)
		if err == nil && len(tempCValue) > 0 {
			tempCFloat, err = strconv.ParseFloat(string(tempCValue), 64)
			if err == nil {
				haveTemp = true
				tempFFloat := CelciusToFahrenheit(tempCFloat)
				urlValues["tempf"] = []string{fmt.Sprintf("%.1f", tempFFloat)}
			} else {
				log.Printf("error: %v", err)
			}
		}
	}

	if *humidityPath != "" {
		humidityValue, err := jsonpointer.Find(inputBytes, *humidityPath)
		if err == nil && len(humidityValue) > 0 {
			humidityFloat, err = strconv.ParseFloat(string(humidityValue), 64)
			if err == nil {
				haveHumidity = true
				urlValues["humidity"] = []string{fmt.Sprintf("%.1f", humidityFloat)}
			} else {
				log.Printf("error: %v", err)
			}
		}
	}

	if haveTemp && haveHumidity {
		dewpointCFloat := ApproximateDewpoint(tempCFloat, humidityFloat)
		dewpointFFloat := CelciusToFahrenheit(dewpointCFloat)
		urlValues["dewptf"] = []string{fmt.Sprintf("%.1f", dewpointFFloat)}
	}

	if *pressurePath != "" {
		pressureValue, err := jsonpointer.Find(inputBytes, *pressurePath)
		if err == nil && len(pressureValue) > 0 {
			pressureFloat, err := strconv.ParseFloat(string(pressureValue), 64)
			if err == nil {
				urlValues["baromin"] = []string{fmt.Sprintf("%.2f", pressureFloat)}
			} else {
				log.Printf("error: %v", err)
			}
		}
	}

	updateURL := fmt.Sprintf("%s?%s", *baseUrl, urlValues.Encode())
	log.Printf("URL: %s", updateURL)
	resp, err := http.Get(updateURL)
	if err != nil {
		log.Printf("error: %v", err)
	} else {
		log.Printf("response: %v", resp.Status)
	}
}

func CelciusToFahrenheit(c float64) float64 {
	return (9.0 / 5.0 * c) + 32.0
}

func ApproximateDewpoint(tempC, humidity float64) float64 {
	// (1) Saturation Vapor Pressure = ESGG(T)
	ratio := 373.15 / (273.15 + tempC)
	rhs := -7.90298 * (ratio - 1)
	rhs += 5.02808 * math.Log10(ratio)
	rhs += -1.3816e-7 * (math.Pow(10, (11.344*(1-1/ratio))) - 1)
	rhs += 8.1328e-3 * (math.Pow(10, (-3.49149*(ratio-1))) - 1)
	rhs += math.Log10(1013.246)

	// factor -3 is to adjust units - Vapor Pressure SVP * humidity
	vp := math.Pow(10, rhs-3) * humidity

	// (2) DEWPOINT = F(Vapor Pressure)
	t := math.Log(vp / 0.61078) // temp var
	return (241.88 * t) / (17.558 - t)
}

package main

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

// CONFIG LATER
var debugMode = false
var url_alert = "http://52.45.17.177:80/XpertRestApi/api/alert_data"
var url_location = "http://52.45.17.177:80/XpertRestApi/api/location_data"

func main() {
	printStringInDebugMode("Starting Server...")
	readSerialDataAndPost()
}

func readSerialDataAndPost() {
	// ports, err := enumerator.GetDetailedPortsList()
	// if err != nil {
	// 	//log.Fatal(err)
	// }

	timeLastPostedLocation := time.Now()
	locationPostingInterval := time.Second * 5
	var lat float64 = 38.443976 + (rand.Float64() / 100)
	var lon float64 = -78.874720 + (rand.Float64() / 100)

	// mode := &serial.Mode{
	// 	BaudRate: 115200, //115200 tag //230400 antenna
	// }
	// openPort, err := serial.Open(ports[0].Name, mode)
	// if err != nil {
	// 	//log.Fatal(err)
	// }
	// buff := make([]byte, 1024) //100 ?

	for {
		var messageToAPI = "ERROR"
		//n, err := openPort.Read(buff)
		// if err != nil {
		// 	messageToAPI = "ERROR READING"
		// } else if n == 0 {
		// 	messageToAPI = "ERROR 0 BYTES READ"
		// } else {
		// 	messageToAPI = string(buff[:n])
		// }

		// AoA decode
		if time.Now().After(timeLastPostedLocation.Add(locationPostingInterval)) {

			var jsonData = []byte(`{
				"deviceimei": 111112222233333,
				"altitude": 1,
				"latitude": ` + strconv.FormatFloat(lat, 'f', -1, 64) + `,
				"longitude": ` + strconv.FormatFloat(lon, 'f', -1, 64) + `,
				"devicetime": 10,
				"speed": 0,
				"Batterylevel": "85",
				"casefile_id": "string",
				"address": "string",
				"positioningmode": "string",
				"tz": "string",
				"alert_type": "string",
				"alert_message": "` + messageToAPI + `",
				"alert_id": "string",
				"offender_name": "string",
				"offender_id": "string"
			}`)

			printStringInDebugMode("Sending Packet to API")
			timeLastPostedLocation = time.Now()
			_, err := http.Post(url_location, "application/json", bytes.NewBuffer(jsonData))
			//_, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

func printStringInDebugMode(str string) {
	if debugMode == true {
		fmt.Println(str)
	}
}

func printArrayInDebugMode(str []string) {
	if debugMode == true {
		fmt.Println(str)
	}
}

package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"go.bug.st/serial/enumerator"
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
	var portArray []string
	ports, err := enumerator.GetDetailedPortsList()
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	for _, port := range ports {
		portArray = append(portArray, port.Name)
	}
	printArrayInDebugMode(portArray)

	timeLastPostedLocation := time.Now()
	locationPostingInterval := time.Second * 5
	var lat float64 = 38.443976 + (rand.Float64() / 100)
	var lon float64 = -78.874720 + (rand.Float64() / 100)

	var portName = "ERROR"
	if len(portArray) > 0 {
		portName = portArray[0]
	}

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
			"alert_message": "` + portName + `",
			"alert_id": "string",
			"offender_name": "string",
			"offender_id": "string"
	}`)

	// mode := &serial.Mode{
	// 	BaudRate: 230400,
	// }
	// port, err := serial.Open(portArray[0], mode)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// buff := make([]byte, 1024) //100 ?

	for {
		// n, err := port.Read(buff)
		// if err != nil {
		// 	log.Fatal(err)
		// 	break
		// }
		// if n == 0 {
		// 	fmt.Println("\nEOF")
		// 	break
		// }
		// printStringInDebugMode(string(buff[:n]))

		if time.Now().After(timeLastPostedLocation.Add(locationPostingInterval)) {
			printStringInDebugMode("Sending Packet to API")
			timeLastPostedLocation = time.Now()
			resp, err := http.Post(url_location, "application/json", bytes.NewBuffer(jsonData))
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Fatal(err)
			}

			if resp.StatusCode == http.StatusOK {
				printStringInDebugMode(string(bodyBytes))
			} else {
				printStringInDebugMode(string(bodyBytes))
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

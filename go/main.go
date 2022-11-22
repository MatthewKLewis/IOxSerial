package main

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"go.bug.st/serial"
	"go.bug.st/serial/enumerator"
)

func main() {
	fmt.Println("Starting...")
	printPorts()
	readSerialDataAndPost()
}

func printPorts() {
	ports, err := enumerator.GetDetailedPortsList()
	if err != nil {
		log.Fatal(err)
	}
	if len(ports) == 0 {
		fmt.Println("No serial ports found!")
		return
	}
	for _, port := range ports {
		fmt.Printf("Found port: %s\n", port.Name)
		if port.IsUSB {
			fmt.Printf("   USB ID     %s:%s\n", port.VID, port.PID)
			fmt.Printf("   USB serial %s\n", port.SerialNumber)
		}
	}
}

func readSerialDataAndPost() {
	timeLastPostedLocation := time.Now()
	locationPostingInterval := time.Second * 5
	var lat float64 = 38.443976 + (rand.Float64() / 100)
	var lon float64 = -78.874720 + (rand.Float64() / 100)

	//url_alert := "http://52.45.17.177:80/XpertRestApi/api/alert_data"
	url_location := "http://52.45.17.177:80/XpertRestApi/api/location_data"
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
			"alert_message": "` + "Hello" + `",
			"alert_id": "string",
			"offender_name": "string",
			"offender_id": "string"
	}`)

	mode := &serial.Mode{
		BaudRate: 230400,
	}
	port, err := serial.Open("COM5", mode)
	if err != nil {
		log.Fatal(err)
	}

	buff := make([]byte, 100)

	for {
		n, err := port.Read(buff)
		if err != nil {
			log.Fatal(err)
			break
		}
		if n == 0 {
			fmt.Println("\nEOF")
			break
		}
		fmt.Printf("%v", string(buff[:n]))

		if time.Now().After(timeLastPostedLocation.Add(locationPostingInterval)) {
			timeLastPostedLocation = time.Now()
			resp, err := http.Post(url_location, "application/json", bytes.NewBuffer(jsonData))
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(resp.Status)
		}
	}
}

package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"time"

	"go.bug.st/serial"
)

func main() {
	fmt.Println("Starting...")

	ports, err := serial.GetPortsList()
	if err != nil {
		log.Fatal(err)
	}
	if len(ports) == 0 {
		log.Fatal("No serial ports found!")
	}
	for _, port := range ports {
		fmt.Printf("Found port: %v\n", port)
	}

	//url_alert := "http://52.45.17.177:80/XpertRestApi/api/alert_data"
	url_location := "http://52.45.17.177:80/XpertRestApi/api/location_data"
	var jsonData = []byte(`{
		"deviceimei": 111112222233333,
		"altitude": 1,
		"latitude": 38.443976,
		"longitude": -78.874720,
		"devicetime": 10,
		"speed": 0,
		"Batterylevel": "85",
		"casefile_id": "string",
		"address": "string",
		"positioningmode": "string",
		"tz": "string",
		"alert_type": "string",
		"alert_message": "string",
		"alert_id": "string",
		"offender_name": "string",
		"offender_id": "string"
	}`)

	for true {
		time.Sleep(time.Second * 3)
		resp, err := http.Post(url_location, "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(resp.Status)
	}

}

package main

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	"strconv"
	"time"

	"go.bug.st/serial"
	"go.bug.st/serial/enumerator"
)

// CONFIG LATER
var debugMode = false
var url_alert = "http://52.45.17.177:80/XpertRestApi/api/alert_data"
var url_location = "http://52.45.17.177:80/XpertRestApi/api/location_data"
var tcpAddr = "52.45.17.177:24888"

func main() {
	printStringInDebugMode("Starting Server...")
	//readSerialDataAndPost()
	readSerialDataAndFwd()
}

func readSerialDataAndPost() {
	ports, err := enumerator.GetDetailedPortsList()
	if err != nil {
		//log.Fatal(err)
	}

	timeLastPostedLocation := time.Now()
	locationPostingInterval := time.Second * 10
	var lat float64 = 38.443996 + (rand.Float64() / 100)
	var lon float64 = -78.874740 + (rand.Float64() / 100)

	mode := &serial.Mode{
		BaudRate: 921600, //115200 tag //230400 antenna //921600 3-6-2023 BLEAP
	}
	openPort, err := serial.Open(ports[0].Name, mode)
	if err != nil {
		//log.Fatal(err)
	}
	buff := make([]byte, 1024) //100 ?

	for {
		var messageToAPI = "ERROR"
		if len(ports) > 0 {
			messageToAPI = ports[0].Name
		}

		n, err := openPort.Read(buff)
		if err != nil {
			messageToAPI = "ERROR READING"
		} else if n == 0 {
			messageToAPI = "ERROR 0 BYTES READ"
		} else {
			messageToAPI = string(buff[:n])
		}

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

func readSerialDataAndFwd() {
	//Get Port Information
	ports, err := enumerator.GetDetailedPortsList()
	if err != nil {
		log.Fatal(err)
	}
	if len(ports) == 0 {
		log.Fatal("No ports found.")
	} else {
		fmt.Println(ports[0].Name)
	}

	//Set Configs
	mode := &serial.Mode{
		BaudRate: 921600, //115200 tag //230400 antenna //921600 3-6-2023 BLEAP
	}

	//Open a COM Port
	openPort, err := serial.Open(ports[0].Name, mode)
	if err != nil {
		log.Fatal(err)
	}
	buff := make([]byte, 1024) //100 ?

	//Dial to TCP
	tcpAddr, err := net.ResolveTCPAddr("tcp4", tcpAddr)
	if err != nil {
		log.Fatal("Resolve Address failed:", err.Error())
	}
	fmt.Println(tcpAddr.String())
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		log.Fatal("Dial to server failed:", err.Error())
	}

	//Report COM data as soon as read
	for {
		n, err := openPort.Read(buff)
		_, err = conn.Write(buff[:n])
		if err != nil {
			log.Fatal("Write to server failed:", err.Error())
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

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
	readSerialDataAndPost()
	//readSerialDataAndFwd()
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
	timeLastPostedLocation := time.Now()
	locationPostingInterval := time.Second * 10
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

	//Report COM data at a regular interval
	for {
		var messageToAPI = "$IQ,1,2233554,872295,13,-65,7,-48,24,-22,-50,47,-24,12,27,-27,20,-32,-33,35,-33,26,37,-5,0,47,13,26,45,33,18,-61,20,31,14,-59,26,53,4,-26,-8,28,-15,-38,43,28,-37,-7,-46,45,-31,64,-79,19,-62,96,-11,3,-51,21,-50,-12,-42,52,40,-15,-51,63,14,-11,-43,94,-46,-28,-42,80,10,-49,-26,24,-72,-50,3,4,-2,-38,3,-43,-45,-62,10,40,-50,-32,31,14,-61,-17,30,47,12,-19,28,4,-57,-14,51,51,39,14,41,-31,88,12,59,-99,60,25,43,14,52,27,50,-69,14,49,27,-50,64,48,3,-58,68,41,-3,-68,43,58,-12,17,65,49,-28,-4,27,24,-36,62,2,25,-52,23,59,31,-42,16,56,14,-50,-41,14,-18,-53,62,48,-15,-46,-66,33,-28,-43\n"
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

		if time.Now().After(timeLastPostedLocation.Add(locationPostingInterval)) {
			fmt.Println(messageToAPI)

			_, err = conn.Write([]byte(messageToAPI))
			if err != nil {
				log.Fatal("Write to server failed:", err.Error())
			}

			timeLastPostedLocation = time.Now()
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

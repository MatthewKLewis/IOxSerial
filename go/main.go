package main

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
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

	// // Real Functions
	readSerialDataAndFwd()
	//readSerialDataAndPost()

	// // Tester Functions
	//simplestCaseFWD()
	//simplestCasePOST()
}

func readSerialDataAndPost() {
	ports, err := enumerator.GetDetailedPortsList()
	if err != nil {
		//log.Fatal(err)
	}

	//timeLastPostedLocation := time.Now()
	//locationPostingInterval := time.Second * 10
	var lat float64 = 38.443986 + (rand.Float64() / 100)
	var lon float64 = -78.874730 + (rand.Float64() / 100)

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
		_, err = http.Post(url_location, "application/json", bytes.NewBuffer(jsonData))
		//_, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		time.Sleep(5 * time.Millisecond)
	}
}

func readSerialDataAndFwd() {
	//Get Port Information
	ports, err := enumerator.GetDetailedPortsList()
	if err != nil {
		os.Exit(1)
	}
	if len(ports) == 0 {
		os.Exit(1)
	}

	//Set Configs
	mode := &serial.Mode{
		BaudRate: 921600, //115200 tag //230400 antenna //921600 3-6-2023 BLEAP
	}

	//Open the 1st? 2nd? 3rd? COM Port
	openPort, err := serial.Open(ports[0].Name, mode)
	if err != nil {
		printStringInDebugMode("trying next COM.")
		openPort, err = serial.Open(ports[1].Name, mode)
		if err != nil {
			printStringInDebugMode("trying next COM..")
			openPort, err = serial.Open(ports[2].Name, mode)
			if err != nil {
				printStringInDebugMode("trying next COM...")
				os.Exit(1)
			}
		}
	}
	buff := make([]byte, 4096) //100 ?

	//Dial to TCP
	tcpAddr, err := net.ResolveTCPAddr("tcp4", tcpAddr)
	if err != nil {
		os.Exit(1)
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		waitAndRestart("rsdFwd")
	}

	//Report COM data as soon as read
	for {
		n, err := openPort.Read(buff)
		if err != nil {
			waitAndRestart("rsdFwd")
		}

		_, err = conn.Write(buff[:n])
		if err != nil {
			waitAndRestart("rsdFwd")
		}
		time.Sleep(20 * time.Millisecond)
	}
}

func simplestCaseFWD() {
	buff := make([]byte, 4096) //100 ?

	//Dial to TCP
	tcpAddr, err := net.ResolveTCPAddr("tcp4", tcpAddr)
	if err != nil {
		log.Fatal("Resolve Address failed:", err.Error())
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		waitAndRestart("scFwd")
	}

	//Report COM data as soon as read
	printStringInDebugMode("Entering for-loop")
	for {
		buff = []byte("Test message")
		_, err = conn.Write(buff)
		if err != nil {
			waitAndRestart("scFwd")
		}
		time.Sleep(5 * time.Millisecond)
	}
}

func simplestCasePOST() {
	var timeLastPostedLocation = time.Now()
	var locationPostingInterval = time.Second * 1
	var lat float64 = 38.443995
	var lon float64 = -78.874741
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
		"alert_message": "Working!!! given deployment to AP",
		"alert_id": "string",
		"offender_name": "string",
		"offender_id": "string"
	}`)

	for {
		printStringInDebugMode("Loop...")
		if time.Now().After(timeLastPostedLocation.Add(locationPostingInterval)) {
			timeLastPostedLocation = time.Now()
			_, err := http.Post(url_location, "application/json", bytes.NewBuffer(jsonData))
			//_, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Fatal(err)
			}
		}
		time.Sleep(1000 * time.Millisecond)
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

func getComPort() {

}

func waitAndRestart(mode string) {
	printStringInDebugMode("Connection Lost. Trying to reconnect in 1 min, then " + mode)
	time.Sleep(2 * time.Second)

	if mode == "" {
		mode = "rsdFwd"
	}
	switch mode {
	case "rsdFwd":
		readSerialDataAndFwd()
	case "rsdPost":
		readSerialDataAndPost()
	case "scFwd":
		simplestCaseFWD()
	case "scPost":
		simplestCasePOST()
	}
	simplestCaseFWD()
}

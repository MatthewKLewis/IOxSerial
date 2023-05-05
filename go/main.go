package main

import (
	"fmt"
	"net"
	"os"
	"time"

	"go.bug.st/serial"
	"go.bug.st/serial/enumerator"
	"gopkg.in/ini.v1"
)

var debug_local = false
var debug_fwd = true

// 115200 tag // 230400 antenna // 921600 3-6-2023 BLEAP
var mode = &serial.Mode{BaudRate: 921600}

// add a new location for debug messages? Using 170 when it gets to Cisco.
var tcpAddr string = ""
var tcpAddrDebug = ""
var connDebug *net.TCPConn = nil
var connFwd *net.TCPConn = nil
var openPort serial.Port = nil

func main() {
	//configs
	cfg, err := ini.Load("package_config.ini")
	if err != nil {
		tcpAddr = "52.45.17.177:24888"
		tcpAddrDebug = "52.45.17.177:25888"
	} else {
		tcpAddr = cfg.Section("upstream").Key("server_ip").String() + ":" + cfg.Section("upstream").Key("forward_port").String()
		tcpAddrDebug = cfg.Section("upstream").Key("server_ip").String() + ":" + cfg.Section("upstream").Key("debug_port").String()
	}
	readSerialDataAndFwd()
}

func readSerialDataAndFwd() {
	//Dial to TCP for debugging
	debug("resolve debug tcp address")
	tcpAddrDebug, err := net.ResolveTCPAddr("tcp4", tcpAddrDebug)
	if err != nil {
		debug("couldn't resolve address for to debug server")
	}
	debug("connect debug")
	connDebug, err = net.DialTCP("tcp", nil, tcpAddrDebug)
	if err != nil {
		connDebug = nil
		debug("couldn't connect to debug server")
	}

	//Dial to TCP for data forwarding
	debug("resolve aoa tcp address")
	tcpAddr, err := net.ResolveTCPAddr("tcp4", tcpAddr)
	if err != nil {
		debug("couldn't resolve tcp address for forwading.")
		os.Exit(1) //KILL
	}
	debug("connect aoa")
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		debug("couldn't connect to tcp address for forwading, restarting...")
		waitAndRestart()
	}

	//Get Port Information
	ports, err := enumerator.GetDetailedPortsList()
	if err != nil {
		debug("Couldn't get ports list")
		os.Exit(1) //KILL
	}
	if len(ports) == 0 {
		debug("Ports list empty")
		os.Exit(1) //KILL
	}

	for i := 0; i < len(ports); i++ {
		debug(ports[i].Name)
	}

	//Open the 1st? 2nd? 3rd? COM Port
	if openPort == nil {
		openPort, err = serial.Open(ports[0].Name, mode)
		if err != nil {
			debug("couldn't open serial port at dev/ttyUSB0.")
			os.Exit(1) //KILL
		}
	}

	//Report COM data as soon as read
	debug("reading and writing!")
	buff := make([]byte, 4096) //100 ?

	//LOOP
	for {
		n, err := openPort.Read(buff)
		if err != nil {
			debug("couldn't read, restarting...")
			waitAndRestart()
		}

		_, err = conn.Write(buff[:n])
		if err != nil {
			debug("couldn't write, restarting...")
			waitAndRestart()
		}
		time.Sleep(20 * time.Millisecond)
	}
}

func debug(msg string) {
	if debug_local { //if debugging locally, print locally.
		fmt.Println(msg)
	}
	if debug_fwd && connDebug != nil { //if debugging on our cloud server, print there.
		_, err := connDebug.Write([]byte(" --- " + msg))
		if err != nil {
			//os.Exit(1) //KILL?
		}
	}
}

func waitAndRestart() {
	debug("Restarting")
	time.Sleep(2 * time.Second)
	readSerialDataAndFwd()
}

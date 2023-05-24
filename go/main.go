package main

import (
	"net"
	"os"
	"time"

	"github.com/tarm/serial"
	"gopkg.in/ini.v1"
)

var debug_fwd = true

// dump error logs in /data/logs MKL 5-23-2023
// add a new location for debug messages? Using 170 when it gets to Cisco.
var tcpAddr string = ""
var tcpAddrDebug = ""
var connDebug *net.TCPConn = nil
var connFwd *net.TCPConn = nil

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
	hostDev := os.Getenv("HOST_DEV")
	if hostDev == "" {
		hostDev = "COM3"
	}
	serialConfig := &serial.Config{Name: hostDev, Baud: 921600}
	openPort, err := serial.OpenPort(serialConfig)
	if err != nil {
		debug("couldn't connect to port")
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

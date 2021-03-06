package main

import (
	"common"
	"encoding/json"
	"lobbyserver"
	//	"log"
	"net"
	"os"
	"path"
	//"strconv"
	"fmt"
	"syscall"
)

var ipcfg common.LobbyServerCfg

func main() {

	//cfgpath, _ := os.Getwd()
	cfgpath := "/home/hanfeng/golang/src/bzmj/bin"
	cfg, err := os.Open(path.Join(cfgpath, "lobbycfg.json"))

	if err != nil {
		println("can't find lobbycfg.json")
		println(cfgpath)
		return
	}

	defer cfg.Close()

	deccfg := json.NewDecoder(cfg)

	if err := deccfg.Decode(&ipcfg); err != nil {
		println("can't find gscfg.json")
		println(err.Error())
		return
	}

	common.DebugInit(ipcfg.GcTime, ipcfg.DebugHost)

	quitChan := make(chan int)

	lobbyServer := lobbyserver.NewLobbyServer(ipcfg)

	listenerForServer, err := net.Listen("tcp", ipcfg.LobbyIpForServer)
	defer listenerForServer.Close()
	if err != nil {
		println("Listening to: ", listenerForServer.Addr().String())
		return
	}

	go lobbyserver.CreateLobbyServicesForCnserver(lobbyServer, listenerForServer)
	go lobbyserver.CreateLobbyServicesForClient(ipcfg.LobbyIpForClient, ipcfg.ConnectType)

	handler := func(s os.Signal, arg interface{}) {
		fmt.Printf("handle signal: %v\n", s)
		println("lobbyserver close")
		os.Exit(0)
	}

	handlerArray := []os.Signal{syscall.SIGINT,
		syscall.SIGILL,
		syscall.SIGFPE,
		syscall.SIGSEGV,
		syscall.SIGTERM,
		syscall.SIGABRT}

	common.WatchSystemSignal(&handlerArray, handler)

	nQuitCount := 0
	for {
		select {
		case <-quitChan:
			nQuitCount = nQuitCount + 1
		}

		if nQuitCount == 2 {
			break
		}
	}

	println("lobbyserver close")

}

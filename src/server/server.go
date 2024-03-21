/* Server module */

package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/signal"

	"github.com/xdsdmg/jail/lib"
)

var number = 0 /* Used to generate the network config of container */

// handleConn handles the request for fetching network config from container.
func handleConn(conn net.Conn) {
	defer func() { number++ }()

	if number >= 255 {
		panic("ip addressed were used up")
	}

	buf := make([]byte, 512)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Printf("[ERROR] [handleConn] read failed, error: %+v\n", err)
		return
	}

	data := buf[0:n]
	fmt.Printf("[INFO] [handleConn] received: %s, len: %dByte\n", string(data), n)

	resp := lib.NetworkConfig{
		VethName:     fmt.Sprintf("container%d", number),
		VethAddr:     fmt.Sprintf("%s.%d/24", lib.ADDR_PREFIX, number+lib.ADDR_SUFFIX_BEGIN),
		BridgeName:   lib.BRIDGE,
		NS:           fmt.Sprintf("ns%d", number),
		DefaultRoute: lib.DEFAULT_ROUTE,
	}
	bs, err := json.Marshal(resp)
	if err != nil {
		fmt.Printf("[ERROR] [handleConn] marshal failed, error: %+v\n", err)
		return
	}

	n, err = conn.Write(bs)
	if err != nil {
		fmt.Printf("[ERROR] [handleConn] read failed, error: %+v\n", err)
		return
	}

	fmt.Printf("[INFO] [handleConn] write successed, len: %dByte\n", n)
}

func removeSockFile() error {
	isSockFileExist := false
	if _, err := os.Stat(lib.SOCK_FILE); err == nil {
		isSockFileExist = true
	} else if !os.IsNotExist(err) {
		return err
	}

	if !isSockFileExist {
		return nil
	}

	if err := os.Remove(lib.SOCK_FILE); err != nil {
		return err
	}

	return nil
}

func main() {
	/* Handle SIGINT */
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			fmt.Printf("[INFO] [signal_handler] signal %v received", sig)
			if err := removeSockFile(); err != nil {
				fmt.Printf("[ERROR] [signal_handler] clean sock file failed, error: %+v\n", err)
			}
			if err := lib.RemoveBridge(lib.BRIDGE); err != nil {
				fmt.Printf("[ERROR] [signal_handler] clean bridge failed, error: %+v\n", err)
			}
			os.Exit(0)
		}
	}()

	if err := lib.CreateBridge(lib.BRIDGE, lib.BRIDGE_ADDR); err != nil {
		panic(fmt.Sprintf("create bridge failed, error: %+v", err))
	}

	l, err := net.Listen("unix", lib.SOCK_FILE)
	if err != nil {
		panic(fmt.Sprintf("create unix socket failed, error: %+v\n", err))
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Printf("[ERROR] [main] error occurred in Accept() of unix socket, %+v\n", err)
			continue
		}

		handleConn(conn)
	}
}

package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"net"
	"time"
)

func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			panic(err)
		}

		go handleRequest(conn)
	}
}

var fwCmd = &cobra.Command{
	Use:   "fw",
	Short: "Print the version number of Hugo",
	Long:  `All software has versions. This is Hugo's`,
	Run: func(cmd *cobra.Command, args []string) {
		ln, err := net.Listen("tcp", ":8080")
		if err != nil {
			panic(err)
		}
		go func() {
			time.Sleep(time.Second*20)
			ln.Close()
		}()
		for {
			conn, err := ln.Accept()
			if err != nil {

				panic("l")
			}

			go handleRequest(conn)
		}	},
}

func init() {
	rootCmd.AddCommand(fwCmd)
}
func handleRequest(conn net.Conn) {
	fmt.Println("new client")

	proxy, err := net.Dial("tcp", "t97.asuscomm.com:2222")
	if err != nil {
		panic(err)
	}

	fmt.Println("proxy connected")
	go copyIO(conn, proxy)
	go copyIO(proxy, conn)
}

func copyIO(src, dest net.Conn) {
	defer src.Close()
	defer dest.Close()
	io.Copy(src, dest)
}
package main

import (
	"flag"
	"fmt"
	"github.com/hashicorp/yamux"
	"io"
	"log"
	"net"
)

var (
	localAddr string
	remoteAddr string
)

func init()  {
	flag.StringVar(&localAddr, "local", "0.0.0.0:81", "listen address")
	flag.StringVar(&remoteAddr, "remote", "127.0.0.1:80", "remote address")
	flag.Parse()
}


func main() {
	fmt.Println("Starting mux client")
	// Get a TCP connection

	// Accept a TCP connection
	listener, err := net.Listen("tcp", localAddr)
	if err != nil{
		println("listen err:", err)
		return
	}

	for{

		remoteconn, err := net.Dial("tcp", remoteAddr)
		if err != nil {
			println("remote conn err:", err)
		}

		// Setup client side of yamux
		log.Println("creating client session")
		session, err := yamux.Client(remoteconn, nil)
		if err != nil {
			println("creating client session err:", err)
			return
		}

		for{
			localconn, err := listener.Accept()
			if err != nil {
				println("accept err:", err)
				return
			}

			go func(conn net.Conn) {
				defer conn.Close()

				// Open a new stream
				log.Println("opening stream")
				forwardConn, err := session.Open()
				if err != nil {
					println("listen err:", err)
					return
				}

				up, down := make(chan int64), make(chan int64)
				go pipe(conn, forwardConn, up)
				go pipe(forwardConn, conn, down)
				<-up
				<-down
				return
			}(localconn)

		}



	}



}


func pipe(src io.Reader, dst io.WriteCloser, result chan<- int64) {
	defer dst.Close()
	n, _ := io.Copy(dst, src)
	result <- int64(n)
}

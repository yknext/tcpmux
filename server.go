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
	forwardAddr string
)

func init()  {
	flag.StringVar(&localAddr, "local", "0.0.0.0:80", "listen address")
	flag.StringVar(&forwardAddr, "forward", "192.168.200.157:80", "forward address")
	flag.Parse()
}

func main() {
	fmt.Println("Starting mux server")
	// Accept a TCP connection
	listener, err := net.Listen("tcp", localAddr)
	if err != nil{
		println("listen err:", err)
		return
	}

	for{

		muxconn, err := listener.Accept()
		if err != nil {
			println("accept err:", err)
			return
		}

		go func(muxconn net.Conn) {
			// Setup server side of yamux
			log.Println("creating server session")
			session, err := yamux.Server(muxconn, nil)
			if err != nil {
				println("creating server session err:", err)
				return
			}


			for{
				// Accept a stream
				log.Println("accepting stream")
				conn, err := session.Accept()
				if err != nil {
					println("accepting stream err:", err)
					return
				}

				go func(conn net.Conn) {
					defer conn.Close()

					forwardConn,err := net.Dial("tcp", forwardAddr)
					if err!= nil {
						println("forward connect error ! " , err)
						return
					}

					up, down := make(chan int64), make(chan int64)
					go pipe(conn, forwardConn, up)
					go pipe(forwardConn, conn, down)
					<-up
					<-down
					return
				}(conn)
			}

		}(muxconn)
	}

}

func pipe(src io.Reader, dst io.WriteCloser, result chan<- int64) {
	defer dst.Close()
	n, _ := io.Copy(dst, src)
	result <- int64(n)
}

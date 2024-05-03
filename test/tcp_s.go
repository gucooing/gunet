package main

import (
	"log"

	"gunet"
)

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	listener, err := gunet.NewTcpS("10.0.0.16:12996")
	if err != nil {
		log.Println(err)
		return
	}
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Println(err)
				return
			}
			go read(conn)
		}
	}()
	select {}
}

func read(conn *gunet.TcpConn) {
	defer conn.Close()
	for {
		bin, err := conn.Read()
		if err != nil {
			log.Println(err)
			return
		}
		log.Println(string(bin))
	}
}

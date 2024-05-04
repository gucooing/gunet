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
			go readS(conn)
		}
	}()
	select {}
}

func readS(conn *gunet.TcpConn) {
	for {
		bin, err := conn.Read()
		if err != nil {
			log.Println(err)
			return
		}
		go writeS(conn, bin)
	}
}

func writeS(conn *gunet.TcpConn, data []byte) {
	_, err := conn.Write(data)
	if err != nil {
		log.Println(err)
		return
	}
}

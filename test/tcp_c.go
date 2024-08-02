package main

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"

	"github.com/gucooing/gunet"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randStringBytes(n int) []byte {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return b
}

func main() {
	conn, err := gunet.NewTcpC("127.0.0.1:12996")
	if err != nil {
		fmt.Println(err)
		return
	}
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	// defer conn.Close()

	/*
		bin := randStringBytes(1024 * 1024)

		time1 := time.Now().UnixMilli()
		for i := 0; i < 1024; i++ {
			var err error
			data := append([]byte(strconv.Itoa(i)), bin...)
			_, err = conn.Write(data)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
		time2 := time.Now().UnixMilli()
		log.Println(time2-time1, "ms")
	*/
	go readC(conn)
	select {}
}

var i int

func readC(conn *gunet.TcpConn) {
	bin := randStringBytes(1024)
	go writeC(conn, bin)
	for {
		_, err := conn.Read()
		if err != nil {
			log.Println(err)
			return
		}
		data := append([]byte(strconv.Itoa(i)), bin...)
		i++
		log.Println(i)
		go writeC(conn, data)
		// time.Sleep(1 * time.Millisecond)
	}
}

func writeC(conn *gunet.TcpConn, data []byte) {
	_, err := conn.Write(data)
	if err != nil {
		log.Println(err)
		return
	}
}

func forWriteC(conn *gunet.TcpConn, data []byte) {
	for {
		i++
		go log.Println(i)
		_, err := conn.Write(data)
		if err != nil {
			log.Println(err)
			return
		}
		// time.Sleep(20 * time.Microsecond)
	}
}

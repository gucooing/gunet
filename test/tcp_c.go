package main

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"

	"gunet"
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
	conn, err := gunet.NewTcpC("10.0.0.16:12996")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

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
	time.Sleep(time.Second)
}

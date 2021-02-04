package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"tmall_seckill/cmd"
)

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		defer close(c)
		cmd.GlobalCancelFunc()
		log.Fatalf("关闭浏览器成功")

	}()
	cmd.Execute()
}

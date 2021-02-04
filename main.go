package main

import (
	"context"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	"tmall_seckill/cmd"
)

func main() {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		cmd.GlobalCancelFunc()
		os.Exit(0)
	}()
	cmd.Execute()
}

func getQrCode(ctx context.Context) {
	var qrCode []byte
	var iframes []*cdp.Node
	err := chromedp.Run(ctx, chromedp.Tasks{
		chromedp.Navigate(`https://login.tmall.com/`),
		chromedp.Sleep(5 * time.Second),
		chromedp.Nodes(`J_loginIframe`, &iframes, chromedp.ByID),
	})
	if err != nil {
		log.Fatal(err)
	}
	err = chromedp.Run(ctx, chromedp.Click(`#login > div.corner-icon-view.view-type-qrcode > i`, chromedp.ByQuery, chromedp.FromNode(iframes[0])),
		chromedp.Sleep(2*time.Second),
		chromedp.Screenshot(`.qrcode-login`, &qrCode, chromedp.ByQuery, chromedp.FromNode(iframes[0])))
	if err := ioutil.WriteFile("e://qrCode.png", qrCode, os.ModePerm); err != nil {
		log.Fatal(err)
	}
}

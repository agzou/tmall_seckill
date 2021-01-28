package main

import (
	"context"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"time"
	"tmall_seckill/auth"
)

func main() {

	ctx, cc := chromedp.NewExecAllocator(context.Background(), initOptions()...)
	defer cc()
	ctx, cc = chromedp.NewContext(ctx, chromedp.WithLogf(log.Printf))
	defer cc()
	loc, _ := time.LoadLocation("Local")
	t, _ := time.ParseInLocation("2006-01-02 15:04:05", "2021-01-28 20:00:00", loc)
	total := 2
	wg := sync.WaitGroup{}
	wg.Add(total)
	if auth.HasCookies() {
		for i := 0; i < total; i++ {
			go func(wg *sync.WaitGroup) {
				defer wg.Done()
				c, _ := chromedp.NewContext(ctx)
				if err := chromedp.Run(c, chromedp.Tasks{
					auth.SetCookies(),
					chromedp.Navigate("https://www.tmall.com/"),
					auth.GoCar(t),
				}); err != nil {
					log.Fatal(err)
				}
			}(&wg)
		}
	} else {
		if err := chromedp.Run(ctx, chromedp.Tasks{
			auth.Login(),
			auth.SaveCookies(),
			chromedp.Navigate("https://www.tmall.com/"),
			auth.GoCar(t),
		}); err != nil {
			log.Fatal(err)
		}
	}
	wg.Wait()
}

func initOptions() []chromedp.ExecAllocatorOption {
	options := []chromedp.ExecAllocatorOption{
		chromedp.Flag("headless", false),                      //debug使用
		chromedp.Flag("blink-settings", "imagesEnabled=true"), //禁用图片加载
		chromedp.Flag("start-maximized", true),                //最大化窗口
		chromedp.Flag("no-sandbox", true),                     //禁用沙盒, 性能优先
		chromedp.Flag("disable-setuid-sandbox", true),         //禁用setuid沙盒, 性能优先
		chromedp.Flag("no-default-browser-check", true),       //不检查默认浏览器
		chromedp.Flag("disable-plugins", true),                //禁用扩展
	}
	options = append(chromedp.DefaultExecAllocatorOptions[:], options...)
	return options
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

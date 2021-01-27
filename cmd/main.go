package main

import (
	"context"
	"encoding/json"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"io/ioutil"
	"log"
	"os"
	"time"
)

func main() {

	ctx, cc := chromedp.NewExecAllocator(context.Background(), initOptions()...)
	defer cc()
	ctx, cc = chromedp.NewContext(ctx, chromedp.WithLogf(log.Printf))
	defer cc()
	if err := setCookies(ctx); err != nil {
		log.Fatal(err)
	}
	getQrCode(ctx)
	time.Sleep(20 * time.Second)

}
func setCookies(ctx context.Context) error {
	return chromedp.Run(ctx, chromedp.ActionFunc(func(ctx context.Context) error {
		bytes, err := ioutil.ReadFile("e://cookies.txt")
		if err != nil {
			return err
		}
		var cookies []*network.CookieParam
		err = json.Unmarshal(bytes, &cookies)
		if err != nil {
			return err
		}
		err = network.SetCookies(cookies).Do(ctx)
		return err
	}))

}
func getCookies(ctx context.Context) error {
	return chromedp.Run(ctx, chromedp.Tasks{
		chromedp.Navigate("https://github.com/login"),
		chromedp.SendKeys("login_field", "test", chromedp.ByID),
		chromedp.SendKeys("password", "test", chromedp.ByID),
		chromedp.Click(`#login > div.auth-form-body.mt-3 > form > input.btn.btn-primary.btn-block`, chromedp.ByQuery),
		chromedp.Sleep(2 * time.Second),
		chromedp.ActionFunc(func(ctx context.Context) error {
			cookies, err := network.GetAllCookies().Do(ctx)
			if err != nil {
				return err
			}
			bytes, err := json.Marshal(cookies)
			if err != nil {
				return err
			}
			err = ioutil.WriteFile("e://cookies.txt", bytes, os.ModePerm)
			return err
		}), chromedp.Sleep(100 * time.Second),
	})
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
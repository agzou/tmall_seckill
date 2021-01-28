package auth

import (
	"context"
	"encoding/json"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"io/ioutil"
	"os"
	"path/filepath"
)

var cookiesPath string

//判断cookies 是否已经存在
func HasCookies() bool {
	_, err := os.Stat(cookiesPath)
	return err == nil
}

//设置登陆得cookies
func SetCookies() chromedp.Action {
	return chromedp.ActionFunc(func(ctx context.Context) error {
		bytes, err := ioutil.ReadFile(cookiesPath)
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
	})
}
func SaveCookies() chromedp.Action {
	return chromedp.ActionFunc(func(ctx context.Context) error {
		cookies, err := network.GetAllCookies().Do(ctx)
		if err != nil {
			return err
		}
		bytes, err := json.Marshal(cookies)
		if err != nil {
			return err
		}
		err = ioutil.WriteFile(cookiesPath, bytes, os.ModePerm)
		return err
	})
}

func init() {
	p, _ := os.UserHomeDir()
	cookiesPath = filepath.Join(p, "cookies.txt")
}

package auth

import (
	"context"
	"fmt"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/dom"
	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
	"github.com/pkg/errors"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"
)

var sku = "20739895092"

func Login() chromedp.Action {
	return chromedp.ActionFunc(func(ctx context.Context) error {
		var iframes []*cdp.Node
		err := chromedp.Run(ctx, chromedp.Navigate(`https://login.tmall.com/`), chromedp.Sleep(5*time.Second),
			chromedp.Nodes(`J_loginIframe`, &iframes, chromedp.ByID))
		if err != nil {
			return errors.WithStack(err)
		}
		return chromedp.Run(ctx, chromedp.Click(`#login > div.corner-icon-view.view-type-qrcode > i`,
			chromedp.ByQuery,
			chromedp.FromNode(iframes[0])),
			chromedp.WaitNotPresent(`J_loginIframe`, chromedp.ByID))
	})
}
func Buy(targetTime time.Time) chromedp.Action {
	return chromedp.Tasks{
		chromedp.Navigate("https://cart.taobao.com/cart.htm?from=btop"),
		chromedp.QueryAfter(fmt.Sprintf(`[href$="%s"]`, sku), func(ctx context.Context, id runtime.ExecutionContextID, node ...*cdp.Node) error {
			if len(node) == 0 {
				return errors.New("找不到购物车链接")
			}
			n := node[0]
			for ; ; n = n.Parent {
				if n.AttributeValue(`class`) == `item-content clearfix` {
					break
				}
			}
			nodeID, err := dom.QuerySelector(n.NodeID, `.cart-checkbox > label`).Do(ctx)
			if err != nil {
				return errors.WithStack(err)
			}
			err = chromedp.Run(ctx,
				chromedp.Click([]cdp.NodeID{nodeID}, chromedp.ByNodeID),
				chromedp.QueryAfter(`#J_Go`, func(ctx context.Context, id runtime.ExecutionContextID, node ...*cdp.Node) error {
					startTime := targetTime.Sub(time.Now())
					if startTime < 0 {
						return errors.New("抢购时间已过")
					}
					log.Printf("将在%s小时后执行任务", startTime.String())
					err := chromedp.Run(ctx,
						chromedp.Sleep(startTime),
						chromedp.Click(`#J_Go`, chromedp.ByQuery),
						chromedp.WaitVisible(`#submitOrderPC_1 .go-btn`, chromedp.ByQuery),
						chromedp.Sleep(1*time.Second),
						chromedp.Click(`#submitOrderPC_1 .go-btn`, chromedp.ByQuery),
						chromedp.Sleep(15*time.Minute),
					)
					if err != nil {
						return err
					}
					return nil
				}, chromedp.ByQuery))
			if err != nil {
				return errors.WithStack(err)
			}
			return nil
		},
			chromedp.ByQuery),
	}
}

//获取登录二维码
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
	p := filepath.Join(filepath.Dir(cookiesPath), "qrCode.png")
	if err := ioutil.WriteFile(p, qrCode, os.ModePerm); err != nil {
		log.Fatal(err)
	}
}

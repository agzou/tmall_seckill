package auth

import (
	"context"
	"errors"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/dom"
	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
	"log"
	"time"
)

func Login() chromedp.Action {
	return chromedp.ActionFunc(func(ctx context.Context) error {
		var iframes []*cdp.Node
		err := chromedp.Run(ctx, chromedp.Navigate(`https://login.tmall.com/`), chromedp.Sleep(5*time.Second),
			chromedp.Nodes(`J_loginIframe`, &iframes, chromedp.ByID))
		if err != nil {
			return err
		}
		return chromedp.Run(ctx, chromedp.Click(`#login > div.corner-icon-view.view-type-qrcode > i`,
			chromedp.ByQuery,
			chromedp.FromNode(iframes[0])),
			chromedp.WaitNotPresent(`J_loginIframe`, chromedp.ByID))
	})
}
func GoCar(targetTime time.Time) chromedp.Action {
	return chromedp.Tasks{
		chromedp.Navigate("https://cart.taobao.com/cart.htm?from=btop"),
		chromedp.QueryAfter(`[href$="20739895092"]`, func(ctx context.Context, id runtime.ExecutionContextID, node ...*cdp.Node) error {
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
				return err
			}
			err = chromedp.Run(ctx,
				chromedp.Click([]cdp.NodeID{nodeID}, chromedp.ByNodeID),
				chromedp.QueryAfter(`#J_Go`, func(ctx context.Context, id runtime.ExecutionContextID, node ...*cdp.Node) error {
					startTime := targetTime.Sub(time.Now())
					if startTime < 0 {
						return errors.New("抢购时间已过")
					}
					log.Printf("ID[%d]将在%s小时后执行任务", id.Int64(), startTime.String())
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
				return err
			}
			return nil
		},
			chromedp.ByQuery),
	}
}

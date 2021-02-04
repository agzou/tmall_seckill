package cmd

import (
	"github.com/chromedp/chromedp"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"log"
	"time"
	"tmall_seckill/auth"
)

var date string
var seckill = &cobra.Command{
	Use:   "seckill",
	Short: "开始抢购",
	RunE: func(cmd *cobra.Command, args []string) error {
		if !auth.HasCookies() {
			log.Printf("未登陆,请先登陆!\n")
			return nil
		}
		ctx, cancelFunc := NewChromedpCtx()
		defer cancelFunc()
		target, _ := time.ParseInLocation("2006-01-02 15:04:05", date, time.Now().Location())
		log.Printf("抢购日期:%s\n", date)
		err := chromedp.Run(ctx, chromedp.Tasks{
			auth.SetCookies(),
			chromedp.Navigate("https://www.tmall.com/"),
			auth.Buy(target),
		})
		return errors.WithStack(err)

	},
}

func init() {
	today := time.Now()
	today = time.Date(today.Year(), today.Month(), today.Day(), 20, 0, 0, 0, today.Location())
	seckill.Flags().StringVarP(&date, "date", "d", today.Format("2006-01-02 15:04:05"), "指定抢购日期")
}

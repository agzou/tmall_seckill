package cmd

import (
	"github.com/chromedp/chromedp"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"log"
	"tmall_seckill/auth"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "登陆账号",
	RunE: func(c *cobra.Command, args []string) error {
		ctx, cancelFunc := NewChromedpCtx()
		defer cancelFunc()
		err := chromedp.Run(ctx, chromedp.Tasks{
			auth.Login(), auth.SaveCookies(),
		})
		if err != nil {
			return errors.WithStack(err)

		}
		log.Printf("登陆成功,保存Cookies路径[%s]\n", auth.GetCookiesPath())
		return nil
	},
}

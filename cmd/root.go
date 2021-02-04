package cmd

import (
	"context"
	"github.com/chromedp/chromedp"
	"github.com/spf13/cobra"
	"log"
)

var Version = "v0.01"
var Debug = false
var rootCmd = &cobra.Command{Use: "tmall_seckill", Version: Version, SilenceErrors: true, SilenceUsage: true}
var GlobalCtx context.Context
var GlobalCancelFunc context.CancelFunc

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		GlobalCancelFunc()
		if Debug {
			log.Fatalf("Error:%+v", err)
		} else {
			log.Fatalf("Error:%v", err)
		}
	}
}
func NewChromedpCtx() (context.Context, context.CancelFunc) {
	ctx, cc := chromedp.NewExecAllocator(GlobalCtx, initOptions()...)
	ctx, cc = chromedp.NewContext(ctx, chromedp.WithLogf(log.Printf))
	return ctx, cc
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
func init() {
	GlobalCtx, GlobalCancelFunc = context.WithCancel(context.Background())
	rootCmd.PersistentFlags().BoolVar(&Debug, "debug", false, "是否开启Debug模式")
	rootCmd.AddCommand(loginCmd)
	rootCmd.AddCommand(seckill)
	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(logoutCmd)
}

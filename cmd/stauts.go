package cmd

import (
	"github.com/spf13/cobra"
	"log"
	"os"
	"time"
	"tmall_seckill/auth"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "获取当前登录状态",
	Run: func(c *cobra.Command, args []string) {
		if !auth.HasCookies() {
			log.Printf("当前未登陆!\n")
		} else {
			f, _ := os.Stat(auth.GetCookiesPath())
			interval := time.Now().Sub(f.ModTime()).Minutes()
			if interval > (30 * time.Minute).Minutes() {
				log.Printf("登陆时间%v,登陆已超过%f分钟,请重新登陆!\n", f.ModTime().Format("2006-01-02 15:04:05"), interval)
				return
			}
			log.Printf("已登陆,登陆时间%s\n", f.ModTime().Format("2006-01-02 15:04:05"))
		}
	},
}

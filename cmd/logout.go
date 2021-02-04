package cmd

import (
	"github.com/spf13/cobra"
	"log"
	"tmall_seckill/auth"
)

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "清除登陆状态",
	Run: func(cmd *cobra.Command, args []string) {
		if auth.HasCookies() {
			log.Printf("存在Cookies文件[%s],删除中...\n", auth.GetCookiesPath())
			auth.RemoveCookies()
			log.Printf("删除文件[%s]成功\n", auth.GetCookiesPath())
		} else {
			log.Printf("未登陆")
		}
	},
}

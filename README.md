# tmall_seckill
## 前言
**⚠本项目只是简单脚本，并不能提高太多抢购概率！**  
**⚠本仓库发布的项目中涉及的任何代码，仅用于测试和学习研究，禁止用于商业用途，不能保证其合法性，准确性，完整性和有效性，请根据情况自行判断。**
## 安装
0. **提前把茅台加入购物车**
1. 确保已经安装Chrome浏览器.
2. 下载项目代码
3. 构建二进制文件
``` shell
cd tmall_seckill & go build
```
## 登录
``` shell
# 在弹出的浏览器中扫码登录淘宝账号，默认会把cookies文件保存在当前目录
./tmall_seckill login
```
## 抢购
``` shell
# 抢购时间，格式为yyyy-MM-dd HH:mm:ss,不填默认为当天20:00:00
./tmall_seckill seckill --date "2021-02-04 20:00:00"
```
## 登出
``` shell
## 删除cookies文件
./tmall_seckill logout 
```
## 状态
```shell
## 判断是否已经登录
./tmall_seckill stauts
```

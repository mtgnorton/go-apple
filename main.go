package main

import (
	"fmt"
	"gopkg.in/gomail.v2"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const Re = `"seoPrice":([^}]*)}[\s\S]*?"title":"([^"]+)"`

var last = ""

func main() {

	ticker := time.NewTicker(5 * time.Second)
	for _ = range ticker.C {
		err := Grab()
		if err != nil {
			log(err)
		}
	}

}

func Grab() error {
	defer func() {
		if err := recover(); err != nil {
			log("发生panic错误,错误信息为:[%v]", err)

		}
	}()
	mailTo := []string{
		"15726204663@163.com",
	}

	resp, err := http.Get("https://www.apple.com.cn/cn-k12/shop/refurbished/ipad")
	if err != nil {
		log("抓取错误", err);
		return err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log("解析响应错误", err);
		return err
	}

	allNode := ExtractAllString(Re, body)
	var amount = 0
	var info = ""
	for _, singleNode := range allNode {
		price := string(singleNode[1])
		title := string(singleNode[2])
		if (strings.Contains(title, "ro") || strings.Contains(title, "ir")) && !strings.Contains(title, "蜂窝") {
			amount++
			info += "商品为: " + title + "   <br>"
			info += "价格为: " + price + "   <br> "
			info += " <br>"
			info += " <br>"
			info += " <br>"
			info += " <br>"
			info += " <br>"
		}

		//if (strings.Contains(title, "5")) {
		//	amount++
		//	info += "商品为: " + title + "   <br>"
		//	info += "价格为: " + price + "   <br> "
		//	info += " <br>"
		//	info += " <br>"
		//	info += " <br>"
		//	info += " <br>"
		//}
	}
	if (amount > 0) {

		if (info != last) {
			last = info
			err = SendMail(mailTo, fmt.Sprintf("官方可买数量为%d", amount), info)
			if err != nil {
				log("邮箱发送错误")
				return err
			}
			log(fmt.Sprintf("官方可买数量为%d", amount), strings.Replace(info, "<br>", "\n", -1))
		}
	}
	log("完成")
	return nil
}

func log(a ...interface{}) {
	t := time.Now().Format("2006-01-02 15:04:05")
	fmt.Println("当前时间为:" + t)
	fmt.Println(a...)
}

func ExtractAllString(reStr string, content []byte) [][][]byte {
	ReMust := regexp.MustCompile(reStr)
	ReMatches := ReMust.FindAllSubmatch(content, -1)
	return ReMatches
}

func SendMail(mailTo []string, subject string, body string) error {
	//定义邮箱服务器连接信息，如果是阿里邮箱 pass填密码，qq邮箱填授权码
	mailConn := map[string]string{
		"user": "851426308",
		"pass": "rqibjzbwnpbzbfhe",
		"host": "smtp.qq.com",
		"port": "465",
	}

	port, _ := strconv.Atoi(mailConn["port"]) //转换端口类型为int

	m := gomail.NewMessage()
	m.SetHeader("From", m.FormatAddress("851426308@qq.com", "官方翻新")) //这种方式可以添加别名，即“XD Game”， 也可以直接用<code>m.SetHeader("From",mailConn["user"])</code> 读者可以自行实验下效果
	m.SetHeader("To", mailTo...)                                     //发送给多个用户
	m.SetHeader("Subject", subject)                                  //设置邮件主题
	m.SetBody("text/html", body)                                     //设置邮件正文

	d := gomail.NewPlainDialer(mailConn["host"], port, mailConn["user"], mailConn["pass"])

	err := d.DialAndSend(m)
	return err

}

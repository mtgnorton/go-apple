package main

import (
	"fmt"
	"gopkg.in/gomail.v2"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

const Re = `"seoPrice":([^}]*)}[\s\S]*?"title":"([^"]+)"`

func main() {
	mailTo := []string{
		"15726204663@163.com",
	}

	resp, err := http.Get("https://www.apple.com.cn/cn-k12/shop/refurbished/ipad")
	if err != nil {
		fmt.Println("抓取错误", err);
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("解析响应错误", err);
	}

	allNode := ExtractAllString(Re, body)
	var amount = 0
	var info = ""
	for _, singleNode := range allNode {
		price := string(singleNode[1])
		title := string(singleNode[2])
		if (strings.Contains(title, "ro") || strings.Contains(title, "ir") ) && !strings.Contains(title, "蜂窝") {
			amount++
			info += "商品为: " + title + "   <br>"
			info += "价格为: " + price + "   <br> "
			info += " <br>"
			info += " <br>"
			info += " <br>"
			info += " <br>"
			info += " <br>"
		}

		//if (strings.Contains(title, "5")) && strings.Contains(title, "蜂窝") {
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
		err := SendMail(mailTo, fmt.Sprintf("官方可买数量为%d", amount), info)
		fmt.Println(err)
	}
	fmt.Println("完成")
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

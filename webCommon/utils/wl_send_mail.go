package util

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/wanglu119/me-deps/webCommon/lib/smtp"
)

type SendMail struct {
	User     string
	Password string
	Host     string
	Port     string
	Auth     smtp.Auth
}

type Attachment struct {
	Name        []string
	ContentType string
	WithFile    bool
	Content     []byte
}

type Message struct {
	Nickname    string
	From        string
	To          []string
	Cc          []string
	Bcc         []string
	Subject     string
	Body        string
	ContentType string
	Attachment  *Attachment
}

func (mail *SendMail) GenAuth() {
	mail.Auth = smtp.PlainAuth("", mail.User, mail.Password, mail.Host)
}

func (mail *SendMail) writeHeader(buffer *bytes.Buffer, Header map[string]string) string {
	header := ""
	for key, value := range Header {
		header += key + ":" + value + "\r\n"
	}
	header += "\r\n"
	buffer.WriteString(header)
	return header
}

func (mail *SendMail) writeFile(buffer *bytes.Buffer, fileName string) {
	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic(err.Error())
	}

	payload := make([]byte, base64.StdEncoding.EncodedLen(len(file)))
	base64.StdEncoding.Encode(payload, file)
	buffer.WriteString("\r\n")
	for index, line := 0, len(payload); index < line; index++ {
		buffer.WriteByte(payload[index])
		if (index+1)%76 == 0 {
			buffer.WriteString("\r\n")
		}
	}
}

func (mail *SendMail) writeContent(buffer *bytes.Buffer, content []byte) {
	payload := make([]byte, base64.StdEncoding.EncodedLen(len(content)))
	base64.StdEncoding.Encode(payload, content)
	buffer.WriteString("\r\n")
	for index, line := 0, len(payload); index < line; index++ {
		buffer.WriteByte(payload[index])
		if (index+1)%76 == 0 {
			buffer.WriteString("\r\n")
		}
	}
}

func (mail *SendMail) Send(message *Message) error {
	mail.GenAuth()
	buffer := bytes.NewBuffer(nil)
	boundary := "GoBoundary"
	Header := make(map[string]string)
	// From:" + nickname + "<" + user + ">
	if len(message.Nickname) > 0 {
		Header["From"] = message.Nickname + "<" + message.From + ">"
	} else {
		Header["From"] = message.From
	}

	Header["To"] = strings.Join(message.To, ";")
	Header["Cc"] = strings.Join(message.Cc, ";")
	Header["Bcc"] = strings.Join(message.Bcc, ";")
	Header["Subject"] = message.Subject
	Header["Content-Type"] = "multipart/related;boundary=" + boundary
	Header["Date"] = time.Now().String()

	mail.writeHeader(buffer, Header)

	var imgsrc string
	if message.Attachment != nil {
		if message.Attachment.WithFile {
			// multiple image transmission
			for _, graphname := range message.Attachment.Name {
				attachment := "\r\n--" + boundary + "\r\n"
				attachment += "Content-Transfer-Encoding:base64\r\n"
				attachment += "Content-Type:" + message.Attachment.ContentType + ";name=\"" + graphname + "\"\r\n"
				attachment += "Content-ID: <" + graphname + "> \r\n\r\n"
				buffer.WriteString(attachment)

				// stitching into html
				imgsrc += "<p><img src=\"cid:" + graphname + "\" height=200 width=300></p><br>\r\n\t\t\t"

				defer func() {
					if err := recover(); err != nil {
						fmt.Printf(err.(string))
					}
				}()

				mail.writeFile(buffer, graphname)
			}
		} else {
			graphname := strings.Join(message.Attachment.Name, ",")
			attachment := "\r\n--" + boundary + "\r\n"
			attachment += "Content-Transfer-Encoding:base64\r\n"
			attachment += "Content-Type:" + message.Attachment.ContentType + ";name=\"" + graphname + "\"\r\n"
			attachment += "Content-ID: <" + graphname + "> \r\n\r\n"
			buffer.WriteString(attachment)

			// stitching into html
			imgsrc += "<p><img src=\"cid:" + graphname + "\" height=200 width=300></p><br>\r\n\t\t\t"

			mail.writeContent(buffer, message.Attachment.Content)
		}
	}

	// The html format that needs to be displayed in the body text
	var template = `
<html>
	<body>
		<p>%s</p><br>
		%s
	</body>
</html>
`
	var content = fmt.Sprintf(template, message.Body, imgsrc)
	body := "\r\n--" + boundary + "\r\n"
	body += "Content-Type: text/html; charset=UTF-8 \r\n"
	body += content
	buffer.WriteString(body)

	buffer.WriteString("\r\n--" + boundary + "--")

	err := smtp.SendMail(mail.Host+":"+mail.Port, mail.Auth, message.From, message.To, buffer.Bytes())
	if err != nil {
		return err
	}
	return nil

}

// -----------------------------------------------------------------------------

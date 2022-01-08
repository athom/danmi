package main

import (
	"fmt"
	"os"

	"github.com/athom/danmi"
)

func main() {
	accountSid := os.Getenv("DANMI_ACCOUNT_SID")
	authToken := os.Getenv("DANMI_AUTH_TOKEN")
	endpoint := os.Getenv("DANMI_AUTH_ENDPOINT")
	templateId := os.Getenv("DANMI_AUTH_TEMPLATE_ID")
	danmi := danmi.NewDanmi(accountSid, authToken, endpoint)
	danmi.EnableDebug = true
	data, err := danmi.SendOTP("131xxxxxxxx", templateId, "hello")
	if err != nil {
		panic(err)
	}
	fmt.Println(data)
}

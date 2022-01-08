package danmi

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// https://www.danmi.com/developer.html#smsSend

// 	{
//		"respDesc":"请求成功。",
//		"smsId":"ed4bb01827334ccaa769203db69c3240",
//		"failList":[
//    		{
//        		"phone":"152XXXXXXXX",
//        		"respCode":"0098"
//    		}
//		],
//		"respCode":"0000"
//	}

type DanmiSendOTPResponse struct {
	SmsId    string `json:"smsId"`
	RespCode string `json:"respCode"`
}

func (this *DanmiSendOTPResponse) IsSuccess() bool {
	return this.RespCode == "0000"
}

func MD5(text string) string {
	ctx := md5.New()
	ctx.Write([]byte(text))
	return hex.EncodeToString(ctx.Sum(nil))
}

func MD5WithLowerCase(text string) string {
	s := MD5(text)
	return strings.ToLower(s)
}

type Danmi struct {
	AccountSid  string
	AuthToken   string
	Endpoint    string
	EnableDebug bool
}

func NewDanmi(accountSid string, authToken string, endpoint string) (r *Danmi) {
	r = &Danmi{}
	r.AccountSid = accountSid
	r.AuthToken = authToken
	r.Endpoint = endpoint
	return
}

func (this *Danmi) makeRequestSign(ts string) (r string) {
	ss := this.AccountSid + this.AuthToken + ts
	r = MD5WithLowerCase(ss)
	return
}

func (this *Danmi) SendOTP(phoneNumber string, templateId string, content string) (r *DanmiSendOTPResponse, err error) {
	var (
		endpoint = this.Endpoint
		sig      string
	)

	ts := fmt.Sprintf("%v", time.Now().Unix()*1000)
	sig = this.makeRequestSign(ts)
	param := content

	form := url.Values{}
	form.Add("accountSid", this.AccountSid)
	form.Add("to", phoneNumber)
	form.Add("templateid", templateId)
	form.Add("param", param)
	form.Add("timestamp", ts)
	form.Add("sig", sig)
	body := strings.NewReader(form.Encode())
	request, err := http.NewRequest("POST", endpoint, body)
	if err != nil {
		return
	}
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	if this.EnableDebug {
		log.Printf("danmi request:\ncurl -H 'Content-Type: application/x-www-form-urlencoded' -d 'accountSid=%v&to=%v&templateid=%v&param=%v&timestamp=%v&sig=%v' %v \n",
			this.AccountSid,
			phoneNumber,
			templateId,
			param,
			ts,
			sig,
			endpoint)
	}
	httpClient := &http.Client{
		Timeout: 5 * time.Second,
	}
	res, err := httpClient.Do(request)
	if err != nil {
		return
	}
	defer res.Body.Close()

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		s := string(b)
		log.Println(s)
		return
	}
	if this.EnableDebug {
		log.Println("danmi reponse:\n", string(b))
	}
	r = &DanmiSendOTPResponse{}
	err = json.Unmarshal(b, &r)
	if err != nil {
		return
	}
	if !r.IsSuccess() {
		err = fmt.Errorf("danmi return failed, %v", r)
		return
	}
	return
}

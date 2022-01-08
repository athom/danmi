# Danmi Go SDK

### Install 

```bash
go get github.com/athom/danmi
```

### Usage

```go
danmi := danmi.NewDanmi(accountSid, authToken, endpoint)
danmi.SendOTP("13145791314", templateId, "hello yu")
```

Happy coding!

## License

MimiPay is released under the [WTFPL License](http://www.wtfpl.net/txt/copying).


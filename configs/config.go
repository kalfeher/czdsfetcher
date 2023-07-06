package configs

import "strings"

var (
	Bucket         = "czds"
	Prefix         = ""
	LocalDirectory = "/tmp"
	Region         = "ap-southeast-2"
	tlds           = "example, ninja, aws, melbourne, dev"
	CZDSuser       = "dummyuser"
	CZDSpassword   = "mypassword"
	CZDStlds       = strings.Split(tlds, ",")
	CZDSserver     = "https://postman-echo.com/post"
)

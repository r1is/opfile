package common

type OssCfg struct {
	TmpSecretID  string `json:"TmpSecretId,omitempty"`
	TmpSecretKey string `json:"TmpSecretKey,omitempty"`
	SessionToken string `json:"Token,omitempty"`
	BucketURL    string `json:"BucketURL"`
	BatchURL     string `json:"BatchURL"`
}

type Data struct {
	Code string `json:"code"`
}

type Resp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

// 下发OSS、COS 临时AK/SK 的云函数地址
var EncUrl = "7ee90589b0e2cda77765931e3591fe82d103a37756acf0124ffce97a0780b5fed99c25f0d5fa2f6c5610ddb33c6c9ce28edd2656d9091f0b9ee4e2674bc0d03056b3fbba1901f7a572b27df4ac4e1659220080407822c37b5eb73cdfa23e41eb"
var SecretForEncUrl = "speedtest"

// excludeExtensions:jpg,gif,png,docx,doc,xlsx,csv,exe,pdf,pbf,xls,rar,txt,svg,jpeg
// excludeDirectories:image,upload

var ExcludeDirectories = []string{"image", "upload", "img"}
var ExcludeExtensions = []string{"jpg", "gif", "png", "docx", "doc", "xlsx", "csv", "exe", "pdf", "pbf", "xls", "rar", "txt", "svg", "jpeg", "log", "iso"}

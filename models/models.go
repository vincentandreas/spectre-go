package models

type GenSiteParam struct {
	Username   string `json:"username" validate:"required,omitempty"`
	Password   string `json:"password" validate:"required"`
	Site       string `json:"site" validate:"required"`
	KeyCounter int    `json:"keyCounter" validate:"required,min=1,max=4294967295"`
	KeyPurpose string `json:"keyPurpose" validate:"required,oneof=password loginName answer"`
	KeyType    string `json:"keyType" validate:"required,oneof=med long max short basic pin name phrase"`
}

type SiteResult struct {
	hashedKey string
	password  string
}

type ApiResponse struct {
	Result   string `json:"result"`
	RespCode string `json:"responseCode"`
	RespMsg  string `json:"responseMessage"`
}

type SiteResultRepository interface {
	Save(hashedKey string, password string)
	FindSiteResult(hashedKey string) string
}

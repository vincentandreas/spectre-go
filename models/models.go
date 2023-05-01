package models

type GenSiteParam struct {
	Username   string `json:"username" validate:"required"`
	Password   string `json:"password" validate:"required"`
	Site       string `json:"site" validate:"required"`
	KeyCounter int    `json:"keyCounter"`
	KeyPurpose string `json:"keyPurpose" validate:"required"`
	KeyType    string `json:"keyType" validate:"required"`
}

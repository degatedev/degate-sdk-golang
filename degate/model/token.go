package model

type AccessToken struct {
	Token  string `json:"token"`
	Expire int64  `json:"expire"`
}

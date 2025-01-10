package models

type ImMessage struct {
	Code    int                    `json:"code"`
	Message map[string]interface{} `json:"message"`
	UserId  int                    `json:"userId"`
}

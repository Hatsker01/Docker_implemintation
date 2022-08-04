package model

type JwtRequestModel struct {
	Token string `json:"token"`
}
type ResponseError struct{
	Error interface{} `json:"error"`
}

type ServerError struct{
	Status string `json:"status"`
	Message string `json:"message"`
}
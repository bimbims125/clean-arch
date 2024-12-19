package utils

// ResponseData represent the response data struct
type ResponseData struct {
	Data interface{} `json:"data"`
}

// ResponseError represent the response error struct
type ResponseError struct {
	Message string `json:"message"`
}

// ResponseSuccess represent the response success struct
type ResponseSuccess struct {
	Message string `json:"message"`
}

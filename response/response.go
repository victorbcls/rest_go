package response

type ResponseHttp struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func NewResponseHttp(status int, message string) ResponseHttp {
	return ResponseHttp{status, message}
}

type DataResponseHttp struct {
	Status int `json:"status"`
	Data   any `json:"data"`
}

func NewDataResponse(status int, data any) DataResponseHttp {
	return DataResponseHttp{status, data}
}

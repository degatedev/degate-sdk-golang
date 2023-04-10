package binance

type Response struct {
	HttpStatusCode int    `json:"http_status_code"`
	HttpBodyText   string `json:"http_body_text"`
	Code           int    `json:"code"`
	Message        string `json:"message"`
	Header         map[string]string `json:"header"`
}

type ListData struct {
	Total int `json:"total"`
}

func (r *Response) Success() bool {
	return r.HttpStatusCode != 0 && r.HttpStatusCode < 400 && r.Code == 0
}

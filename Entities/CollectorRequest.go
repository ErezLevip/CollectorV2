package Entities

type CollectorRequest struct {
	Gender string `json:"gender"`
	D1 float64 `json:"d1"`
	P1 float64 `json:"p1"`
	App_name string `json:"app_name"`
	Mime_type string `json:"mime_type"`
	N1 float64 `json:"n1"`
	Cc string `json:"cc"`

}

type Request struct{
	Data string `json:"data"`
}
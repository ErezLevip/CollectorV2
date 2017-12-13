package Collector

import (
	"Collector/Entities"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"strconv"
)

type RequestHandler struct {
}

func (r *RequestHandler) Collect(ctx context.Context, w http.ResponseWriter, req *http.Request) {
	var request Entities.Request
	_ = json.NewDecoder(req.Body).Decode(&request)

	collectorRequest, err := ConvertRequestToCollectorRequest(request)
	if (err != nil) {
		log.Println(err)
		return
	}

	ch := ctx.Value("Requests").(chan Entities.CollectorRequest)
	go func() {
		ch <- collectorRequest
	}()
	//log.Println(request)
}

func ConvertRequestToCollectorRequest(req Entities.Request) (res Entities.CollectorRequest, err error) {
	splited := strings.Split(req.Data, ",")
	res = Entities.CollectorRequest{}

	res.Gender = splited[0]
	res.D1, err = strconv.ParseFloat(splited[1], 64)
	res.P1, err = strconv.ParseFloat(splited[2], 64)
	res.App_name = splited[3]
	res.Mime_type = splited[4]
	res.N1, err = strconv.ParseFloat(splited[5], 64)
	res.Cc = splited[6]

	return
}

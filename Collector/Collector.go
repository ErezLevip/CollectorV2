package Collector

import (
	"Collector/Entities"
	"log"
	"net/http"
	"time"
	"context"
	"Collector/DbHandler"
	"github.com/gorilla/mux"
)

const bufferedChannelSize = 5

type Collector struct {
	conStr string
}

func (c Collector) Make(conStr string) *Collector {
	return &Collector{
		conStr: conStr,
	}
}

func (c *Collector) Run() {
	m := mux.NewRouter()

	handler := &RequestHandler{}

	requestCh := make(chan Entities.CollectorRequest, bufferedChannelSize)
	ctx := context.Background()
	ctx = context.WithValue(ctx, "Requests", requestCh)

	m.HandleFunc("/collect", RequestWrapper(ctx, handler.Collect)).Methods("POST")

	go func() {
		c.ProcessRequests(requestCh)
	}()

	port := "8000"
	srv := &http.Server{
		Handler: m,
		Addr:    ":" + port,

		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Println("listening on port", port)
	log.Fatal(srv.ListenAndServe())
}

type requestHandle func(ctx context.Context, w http.ResponseWriter, req *http.Request)

func RequestWrapper(ctx context.Context, handlerFunc requestHandle) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		handlerFunc(ctx, w, req)
	}
}

func (c *Collector) ProcessRequests(requests chan Entities.CollectorRequest) {

	dbHandler := DbHandler.MySqlDbHandler{}.Make(c.conStr)
	/*redisHandler := &RedisHandler.RedisHandler{}
	redisHandler.Init(RedisHandler.RedisConfiguration{
		ConnectionString:"",
		Credentials:"",
		Db:1,
	})*/

	dbChannel := make(chan Entities.CollectorRequest, bufferedChannelSize)
	//redisChannel := make(chan Entities.CollectorRequest, bufferedChannelSize)

	go func() {
		dbHandler.LogRequest(dbChannel)
	}()
	for r := range requests {
		dbChannel <- r
		//redisChannel <- r
	}

	/*go func() {
		redisHandler.Set()
	}()*/

}

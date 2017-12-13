package main

import (
	"net/http"
	"fmt"
	"bytes"
	"encoding/json"
)

func main() {

	options := []string{
		"Male,1,5,Alphazap,text/plain,13.6,FR",
		"Female,0,1,Andalax,application/pdf,3.8,IT",
		"Male,0,5,Redhold,video/mpeg,16.2,FR",
		"Male,1,4,Biodex,application/x-mspowerpoint,4.6,US",
		"Male,1,1,Treeflex,application/x-troff-msvideo,7.8,US",
		"Male,0,4,Konklux,image/x-tiff,5.3,FR",
		"Male,0,5,Keylex,audio/mpeg3,20.6,US",
	}

	type request struct {
		Data string `json:"data"`
	}

	for _, o := range options {

		r := request{
			Data: o,
		}
		url := "http://localhost:8000/collect"
		fmt.Println("URL:>", url)

		b, err := json.Marshal(r)
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(b))
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
	}
}

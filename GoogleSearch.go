package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"strings"
	"sync"
	"time"
)
// GResponse Struct (Model)
type GResponse struct {
	Response  string `json:"response"`
	Time_took string `json:"time_took"`
	Keyword   string `json:"keyword"`
}
var gResponse []GResponse
func createIndividual(keyword string,url string,wg *sync.WaitGroup) {
	var record GResponse
	wg.Done()
	start := time.Now()
	ctx,cl:=context.WithTimeout(context.Background(),3000*time.Millisecond)
	defer cl()
	req,err:=http.NewRequestWithContext(ctx,http.MethodGet,url,nil)
	if err != nil {
		gResponse = append(gResponse,record)
		log.Fatal(err)
	}
	response,err:=http.DefaultClient.Do(req)
	//response, err := http.Get(url)
	elapsed := time.Since(start)
	if err != nil {
		gResponse = append(gResponse,record)
		log.Fatal(err)
	}
	defer response.Body.Close()
	var bodyString string
	if response.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}
		bodyString = string(bodyBytes)
	}
	//json.NewEncoder(w).Encode(bodyString)
	record.Keyword=keyword
	record.Response=bodyString
	record.Time_took=elapsed.String()
	// push the population object down the channel
	gResponse = append(gResponse,record)
	// let the wait group know we finished

}
// Init gResponse var as a slice Google Response struct
func getResponse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	keywords := params["id"]
	keywordArray := strings.Split(keywords, ",")
	//ch := make(chan GResponse, len(keywordArray))
	var wg sync.WaitGroup
	wg.Add(len(keywordArray))
	generalURL := "https://www.google.com/search?sxsrf=ALeKk02javhrMKcd_J3lDJfIa-Wa7gx6ug%3A1610593932726&ei=jLb_X-3kK82srQHHuryQDg&q=ram&oq=ram&gs_lcp=CgZwc3ktYWIQDFAAWABg3a87aABwAXgAgAEAiAEAkgEAmAEAqgEHZ3dzLXdpeg&sclient=psy-ab&ved=0ahUKEwjt2vL5uZruAhVNVisKHUcdD-IQ4dUDCA0"
	for index := range keywordArray {
		url := strings.Replace(generalURL, "ram", keywordArray[index], -1)
		go createIndividual(keywordArray[index],url,&wg)
	}
	wg.Wait()
	for index := range gResponse{
		json.NewEncoder(w).Encode(gResponse[index].Keyword + " " + gResponse[index].Time_took)
	}
	//json.NewEncoder(w).Encode(gResponse)
	gResponse = nil
}
func main() {
	// Init Router
	fmt.Print(math.Sqrt(4.0))
	r := mux.NewRouter()
	r.HandleFunc("/search/keywords={id}", getResponse).Methods("GET")
	log.Fatal(http.ListenAndServe(":3000", r))
}
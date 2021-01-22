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
var (
	generalURL = "https://www.google.com/search?q="
)
// KeywordResponse response for each keyword
type KeywordResponse struct {
	Response string `json:"response"`
	TimeTook string `json:"time_took"`
	Keyword  string `json:"keyword"`
}
// GoogleResponse final response with all keywords
type GoogleResponse struct{
	Responses []KeywordResponse
}
func createIndividual(keyword string,url string,wg *sync.WaitGroup, googleResponse *GoogleResponse) {
	defer wg.Done()
	var keywordResponse KeywordResponse
	start := time.Now()
	ctx, cncl := context.WithTimeout(context.Background(),8000*time.Millisecond)
	defer cncl()
	req,err:=http.NewRequestWithContext(ctx,http.MethodGet,url,nil)
	if err != nil {
		fmt.Println("can not create request")
		return
	}
	response,err := http.DefaultClient.Do(req)
	elapsed := time.Since(start)
	if err != nil {
		fmt.Println("error while searching: " + keyword)
		return
	}
	defer response.Body.Close()
	var bodyString string
	if response.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(response.Body)
		if err != nil {
			//fmt.Println("error 3")
			//log.Fatal(err)
			return
		}
		bodyString = string(bodyBytes)
	}
	keywordResponse.Keyword=keyword
	keywordResponse.TimeTook =elapsed.String()
	keywordResponse.Response = bodyString
	googleResponse.Responses = append(googleResponse.Responses, keywordResponse)
}
// getKeywordsResponse
func getKeywordsResponse(w http.ResponseWriter, r *http.Request) {
	var googleResponse GoogleResponse
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	keywords := strings.Split(params["keywords"], ",")
	var wg sync.WaitGroup
	wg.Add(len(keywords))
	for index := range keywords {
		url := generalURL + keywords[index]
		go createIndividual(keywords[index],url,&wg,&googleResponse)
	}
	wg.Wait()
	json.NewEncoder(w).Encode(googleResponse)
}
func main() {
	// Init Router
	fmt.Print(math.Sqrt(4.0))
	r := mux.NewRouter()
	r.HandleFunc("/search/keywords={keywords}", getKeywordsResponse).Methods("GET")
	log.Fatal(http.ListenAndServe(":3000", r))
}
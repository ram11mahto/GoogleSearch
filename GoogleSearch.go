package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"strings"
	"time"
)
// GResponse Struct (Model)
type GResponse struct {
	Keyword   string `json:"keyword"`
	Response  string `json:"response"`
	Time_took string `json:"time_took"`
}
// Init gResponse var as a slice Google Response struct
var gResponse []GResponse
func getResponse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	keywords := params["id"]
	keywords = keywords[10:]
	json.NewEncoder(w).Encode(keywords)
	keywordArray := strings.Split(keywords, ",")
	generalURL := "https://www.google.com/search?sxsrf=ALeKk02javhrMKcd_J3lDJfIa-Wa7gx6ug%3A1610593932726&ei=jLb_X-3kK82srQHHuryQDg&q=ram&oq=ram&gs_lcp=CgZwc3ktYWIQDFAAWABg3a87aABwAXgAgAEAiAEAkgEAmAEAqgEHZ3dzLXdpeg&sclient=psy-ab&ved=0ahUKEwjt2vL5uZruAhVNVisKHUcdD-IQ4dUDCA0"
	for index := range keywordArray {
		url := strings.Replace(generalURL, "ram", keywordArray[index], -1)
		startTime := time.Now()
		response, err := http.Get(url)
		elapsedTime := time.Since(startTime)
		if err != nil {
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
		gResponse = append(gResponse, GResponse{keywordArray[index], bodyString, elapsedTime.String()})
	}
	json.NewEncoder(w).Encode(gResponse)
}
func main() {
	// Init Router
	fmt.Print(math.Sqrt(4.0))
	r := mux.NewRouter()
	r.HandleFunc("/search/keywords={id}", getResponse).Methods("GET")
	log.Fatal(http.ListenAndServe(":3000", r))
}
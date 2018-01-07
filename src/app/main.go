package main

import (
	"app/transit_realtime"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	proto "github.com/golang/protobuf/proto"

	"github.com/julienschmidt/httprouter"
)

func indexHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprintf(w, "This is the RESTful api")
}

func getBuses(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var client = &http.Client{
		Timeout: time.Second * 10,
	}

	req, _ := http.NewRequest("GET", "https://api.transport.nsw.gov.au/v1/gtfs/vehiclepos/buses", nil)
	req.Header.Add("Authorization", "apikey iiuBRFGdRfFWyswZtHRHRlNt77i10lwpph0H")

	resp, error := client.Do(req)
	if error == nil {
		fmt.Println("All went well!")
		buses := &transit_realtime.FeedMessage{}
		gtfs, _ := ioutil.ReadAll(resp.Body)
		if err := proto.Unmarshal(gtfs, buses); err != nil {
			log.Fatalln("Failed to parse GTFS data!")
		}

		fmt.Printf("%T\n", buses)
	}
	fmt.Println("Getting buses!")
}

func main() {
	router := httprouter.New()
	router.GET("/", indexHandler)
	router.GET("/buses", getBuses)

	// print env
	env := os.Getenv("APP_ENV")
	if env == "production" {
		log.Println("Running api server in production mode")
	} else {
		log.Println("Running api server in dev mode")
	}

	http.ListenAndServe(":8080", router)
}

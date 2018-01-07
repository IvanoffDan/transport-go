package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"app/dtos"
	"app/transit_realtime"

	proto "github.com/golang/protobuf/proto"
)

func GetBuses(w http.ResponseWriter, r *http.Request) {
	var client = &http.Client{
		Timeout: time.Second * 10,
	}

	req, _ := http.NewRequest("GET", "https://api.transport.nsw.gov.au/v1/gtfs/vehiclepos/buses", nil)
	req.Header.Add("Authorization", "apikey iiuBRFGdRfFWyswZtHRHRlNt77i10lwpph0H")

	resp, error := client.Do(req)
	defer resp.Body.Close()
	if error == nil {
		fmt.Println("All went well!")
		feed := &transit_realtime.FeedMessage{}
		gtfs, _ := ioutil.ReadAll(resp.Body)
		if err := proto.Unmarshal(gtfs, feed); err != nil {
			log.Fatalln("Failed to parse GTFS data!")
		}

		buses := feed.GetEntity()

		vehicles := make([]dtos.Vehicle, 0)

		for _, bus := range buses {
			position := bus.Vehicle.GetPosition()
			vehicleId := bus.Vehicle.Vehicle.GetId()
			routeId := bus.Vehicle.Trip.GetRouteId()

			vehicle := &dtos.Vehicle{
				Position: &dtos.Position{position.GetLatitude(), position.GetLongitude()},
				Id:       vehicleId,
				RouteId:  routeId,
			}

			vehicles = append(vehicles, *vehicle)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(vehicles)
	}

	fmt.Println("Getting buses!")
}

package main

import (
	"app/api"
	"app/dtos"
	"app/transit_realtime"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"testing"
	"time"

	proto "github.com/golang/protobuf/proto"

	"github.com/julienschmidt/httprouter"
)

type Store struct {
	Buses chan []dtos.Vehicle
}

func indexHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprintf(w, "This is the RESTful api")
}

func getBusesHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	api.GetBuses(w, r)
}

func getBuses() []*transit_realtime.FeedEntity {

	var result = make([]*transit_realtime.FeedEntity, 0)

	var client = &http.Client{
		Timeout: time.Second * 10,
	}

	req, _ := http.NewRequest("GET", "https://api.transport.nsw.gov.au/v1/gtfs/vehiclepos/buses", nil)
	req.Header.Add("Authorization", "apikey iiuBRFGdRfFWyswZtHRHRlNt77i10lwpph0H")

	resp, error := client.Do(req)
	defer resp.Body.Close()
	if error == nil {
		feed := &transit_realtime.FeedMessage{}
		gtfs, _ := ioutil.ReadAll(resp.Body)
		if err := proto.Unmarshal(gtfs, feed); err != nil {
			log.Fatalln("Failed to parse GTFS data!")
		} else {
			result = feed.GetEntity()
		}
	}

	return result
}

func BenchFormatBusesFromGTFSMostNaive(b *testing.B) {
	buses := getBuses()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		formatBusesFromGTFSMostNaive(buses)
	}
}

func BenchFormatBusesFromGTFS1(b *testing.B) {
	buses := getBuses()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		formatBusesFromGTFS1(buses)
	}
}

func BenchFormatBusesFromGTFS2(b *testing.B) {
	buses := getBuses()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		formatBusesFromGTFS2(buses)
	}
}

func formatBusesFromGTFSMostNaive(buses []*transit_realtime.FeedEntity) []dtos.Vehicle {
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

	return vehicles
}

func formatBusesFromGTFS1(buses []*transit_realtime.FeedEntity) []dtos.Vehicle {
	vehicles := make([]dtos.Vehicle, 0)
	for _, bus := range buses {
		go func() {
			position := bus.Vehicle.GetPosition()
			vehicleId := bus.Vehicle.Vehicle.GetId()
			routeId := bus.Vehicle.Trip.GetRouteId()

			vehicle := &dtos.Vehicle{
				Position: &dtos.Position{position.GetLatitude(), position.GetLongitude()},
				Id:       vehicleId,
				RouteId:  routeId,
			}

			vehicles = append(vehicles, *vehicle)
			// runtime.Gosched()
		}()
	}

	return vehicles
}

func formatBusesFromGTFS2(buses []*transit_realtime.FeedEntity) []dtos.Vehicle {
	vehicles := make([]dtos.Vehicle, 0)
	var wg sync.WaitGroup

	queue := make(chan dtos.Vehicle, 1)

	for _, bus := range buses {
		wg.Add(1)
		go func(bus *transit_realtime.FeedEntity) {
			position := bus.Vehicle.GetPosition()
			vehicleId := bus.Vehicle.Vehicle.GetId()
			routeId := bus.Vehicle.Trip.GetRouteId()

			vehicle := &dtos.Vehicle{
				Position: &dtos.Position{position.GetLatitude(), position.GetLongitude()},
				Id:       vehicleId,
				RouteId:  routeId,
			}

			queue <- *vehicle
		}(bus)
	}

	go func() {
		for v := range queue {
			vehicles = append(vehicles, v)
			wg.Done()
		}
	}()

	wg.Wait()

	return vehicles
}

func main() {

	// store := Store{make(chan []dtos.Vehicle)}
	// ticker := time.Tick(time.Second * 5)

	// go func() {
	// 	for _ = range ticker {

	// 		gtfsEntities := getBuses()

	// 		vehicles := formatBusesFromGTFS(gtfsEntities)

	// 		store.Buses <- vehicles
	// 	}
	// }()

	// go func(c <-chan []dtos.Vehicle) {
	// 	for vehicles := range store.Buses {
	// 		fmt.Println("Total buses:", len(vehicles))
	// 	}
	// }(store.Buses)

	fmt.Println(testing.Benchmark(BenchFormatBusesFromGTFSMostNaive))
	fmt.Println(testing.Benchmark(BenchFormatBusesFromGTFS1))
	fmt.Println(testing.Benchmark(BenchFormatBusesFromGTFS2))

	// router := httprouter.New()
	// router.GET("/", indexHandler)
	// router.GET("/api/buses", getBusesHandler)
	// // router.NotFound = http.FileServer(http.File("public/index.html"))

	// // print env
	// env := os.Getenv("APP_ENV")
	// if env == "production" {
	// 	log.Println("Running api server in production mode")
	// } else {
	// 	log.Println("Running api server in dev mode")
	// }

	// http.ListenAndServe(":8080", router)
}

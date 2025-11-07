package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/S1FFFkA/PET_PROJECT_N1_WEATHER/internal/client/HTTP/geocoding"
	open_meteo "github.com/S1FFFkA/PET_PROJECT_N1_WEATHER/internal/client/HTTP/open-meteo"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-co-op/gocron/v2"
)

const httpPort = ":3000"

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	httpClient := &http.Client{Timeout: time.Second * 10}
	geocodingClient := geocoding.NewClient(httpClient)
	openMeteoClient := open_meteo.NewClient(httpClient)
	r.Get("/{city}", func(w http.ResponseWriter, r *http.Request) {
		city := chi.URLParam(r, "city")
		fmt.Fprintf(w, "Hello, %s\n!", city)
		geores, err := geocodingClient.GetCoords(city)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		openMetRes, err := openMeteoClient.GetTemperature(geores.Latitude, geores.Longitude)
		raw, err := json.Marshal(openMetRes)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		_, err = w.Write(raw)
		if err != nil {
			log.Println(err)
		}
	})

	s, err := gocron.NewScheduler()
	if err != nil {
		panic(err)
	}
	jobs, err := initJobs(s)
	if err != nil {
		panic(err)
	}

	wg := &sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		fmt.Println("starting server on port " + httpPort)
		err := http.ListenAndServe(httpPort, r)
		if err != nil {
			panic(err)
		}

	}()

	go func() {
		defer wg.Done()
		fmt.Printf("starting job : %v\n", jobs[0].ID())
		s.Start()
	}()
	wg.Wait()

}

func initJobs(scheduler gocron.Scheduler) ([]gocron.Job, error) {
	// create a scheduler

	j, err := scheduler.NewJob(
		gocron.DurationJob(
			10*time.Second,
		),
		gocron.NewTask(
			func() {
				fmt.Println("hello world")
			},
		),
	)
	if err != nil {
		return nil, err
	}
	return []gocron.Job{j}, nil
}
func runCron() {

	//	fmt.Println(j.ID())

	// start the scheduler
	//	s.Start()

}

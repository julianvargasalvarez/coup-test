package scooters

import (
	"net/http"
        "encoding/json"
        "time"
        //"io"
        "strings"
        "fmt"
        "sync"
        "strconv"
)

type Scooter struct {
  ID int64 `json:"id"`
  BatteryLevel int `json:"battery_level"`
  AvailableForRent bool `json:"available_for_rent"`
}

const jsonStream = `
    {"id": 1, "battery_level": 30, "available_for_rent": true}
`

//const jsonStream = `
//    {}
//`


func fetchScooter(id int64, scootersChannel chan Scooter, wg *sync.WaitGroup) {
        defer wg.Done()
        var m Scooter
        err := json.NewDecoder(strings.NewReader(jsonStream)).Decode(&m);

        if err != nil {
	  fmt.Printf(err.Error())
        }

        if(m.AvailableForRent && m.BatteryLevel > 20 && len(scootersChannel)<cap(scootersChannel)){
          scootersChannel <- m
          fmt.Println("Scooter added")
        }
}

func fetchAvailableScooters(upTo int) []Scooter {
        var idCounter int64
        idCounter = 1 // id for the first scooter

        scootersChannel := make(chan Scooter, upTo)
        var wg sync.WaitGroup // for syncing the goroutienes

        go func(){
               wg.Wait()
               close(scootersChannel)
        }()

        go func(){
              time.Sleep(1*time.Second) // After 1000 milliseconds the channel is closed
              close(scootersChannel)
        }()

        for i := 1; i<upTo*200; i++{
          wg.Add(1)
          go fetchScooter(idCounter, scootersChannel, &wg)
          idCounter += 1
          if len(scootersChannel) >= upTo {
            fmt.Println("Max reached")
            close(scootersChannel)
            break
          }
        }

        finalResult := []Scooter{}

        for response := range scootersChannel {
          finalResult = append(finalResult, response)
        }

        return finalResult
}

func Index(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")

        max, err := strconv.Atoi(r.URL.Query().Get("max"))

        if err != nil {
          fmt.Println(err.Error())
        }

        finalResult := fetchAvailableScooters(max)
        json.NewEncoder(w).Encode(finalResult)
}

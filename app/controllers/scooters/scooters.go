package scooters

import (
	"net/http"
        "encoding/json"
        "time"
        //"io"
        //"strings"
        "fmt"
        "sync"
        "strconv"
)

type Scooter struct {
  ID int64 `json:"id"`
  BatteryLevel int `json:"battery_level"`
  AvailableForRent bool `json:"available_for_rent"`
}

//const jsonStream = `
//    {"id": 1, "battery_level": 30, "available_for_rent": true}
//`

//const jsonStream = `
//    {}
//`


func ChannelIsClosed(ch <-chan Scooter) bool {
	select {
	case <-ch:
		return true
	default:
	}

	return false
}

func fetchScooter(id int64, scootersChannel chan Scooter, wg *sync.WaitGroup) {
        defer func() {
                if r := recover(); r != nil {
                        fmt.Println(r)
                        wg.Done()
                }
        }()

        var myClient = &http.Client{Timeout: 1 * time.Second}

        response, err := myClient.Get(fmt.Sprintf("%s%d", "https://qc05n0gp78.execute-api.eu-central-1.amazonaws.com/prod/BackendGoChallenge?id=", id))

        if err != nil {
	  fmt.Println(err.Error())
        }
        defer response.Body.Close()

        defer wg.Done()
        var m Scooter
        //err = json.NewDecoder(strings.NewReader(jsonStream)).Decode(&m);
        err = json.NewDecoder(response.Body).Decode(&m);

        if err != nil {
	  fmt.Printf(err.Error())
        }

        if(m.AvailableForRent && m.BatteryLevel > 20 && len(scootersChannel)<cap(scootersChannel) && !ChannelIsClosed(scootersChannel)){
          scootersChannel <- m
          fmt.Println("Scooter added")
        }
}

func fetchAvailableScooters(upTo int) []Scooter {
        var idCounter int64
        idCounter = 1 // id for the first scooter

        scootersChannel := make(chan Scooter, upTo)
        var wg sync.WaitGroup // for syncing the goroutienes

        //closes the channel when the go routines are all done
        go func(){
               wg.Wait()
               if !ChannelIsClosed(scootersChannel) {
                 close(scootersChannel)
               }
        }()

        // After 1000 milliseconds the channel is closed
        go func(){
              time.Sleep(1*time.Second)
               if !ChannelIsClosed(scootersChannel) {
                 close(scootersChannel)
               }
        }()

        for i := 1; i<upTo*200; i++{
          wg.Add(1)
          go fetchScooter(idCounter, scootersChannel, &wg)
          idCounter += 1
          if len(scootersChannel) >= upTo {
            fmt.Println("Maximum reached")
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

// main
package main

import (
	"encoding/json"
	"github.com/codegangsta/martini"
	slopeone "github.com/ginuerzh/go-slope-one"
	"github.com/ginuerzh/recommendsys/models"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	so *slopeone.SlopeOne
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func classic() *martini.ClassicMartini {
	r := martini.NewRouter()
	m := martini.New()
	m.Use(martini.Logger())
	m.Use(martini.Recovery())
	//m.Use(martini.Static("public"))
	m.Action(r.Handle)
	return &martini.ClassicMartini{m, r}
}

func slopeOneHandler(request *http.Request, resp http.ResponseWriter) {
	var rate map[string]float32
	ids := make([]string, 0)

	rdata, err := ioutil.ReadAll(request.Body)
	if err != nil {
		resp.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := json.Unmarshal(rdata, &rate); err != nil {
		log.Println(err)
	}

	result := so.Predict(rate)
	for k, _ := range result {
		if _, ok := rate[k]; !ok {
			ids = append(ids, k)
		}
	}
	//log.Println(ids)
	data, _ := json.Marshal(ids)
	resp.Header().Set("Content-Type", "application/json; charset=utf-8")
	resp.Write(data)
}

func main() {
	prepare()

	go func() {
		ticker := time.NewTicker(time.Minute * 60)
		for {
			select {
			case <-ticker.C:
				prepare()
			}
		}
	}()

	m := classic()
	m.Map(log.New(os.Stdout, "[martini] ", log.LstdFlags))
	//m.Map(controllers.NewRedisLogger())

	m.Post("/slopeone", slopeOneHandler)
	http.ListenAndServe(":8090", m)
}

func prepare() {
	var rates []map[string]float32

	models.IterRate(func(userRate *models.UserRate) {
		rate := make(map[string]float32)
		for i, _ := range userRate.Rates {
			rate[userRate.Rates[i].Article] = float32(userRate.Rates[i].Rate)
		}
		rates = append(rates, rate)
	})
	//log.Println("prepare")
	so = slopeone.NewSlopeOne(rates)
}

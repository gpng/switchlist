// Package switchlist provides API endpoints for interfacing with Nintendo eShop API's, and retrieving games prices and discounts
package switchlist

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/eshop/us", eShopUSHandler).Methods("GET")

	if err := http.ListenAndServe(":5000", r); err != nil {
		log.Fatal(err)
	}
}

// Game : Object from US eShop Response
type Game struct {
	Categories struct {
		Category json.RawMessage `json:"category"`
	} `json:"categories"`
	Slug            string `json:"slug"`
	BuyItNow        string `json:"buyitnow"`
	ReleaseDate     string `json:"release_date"`
	DigitalDownload string `json:"digitaldownload"`
	FreeToStart     string `json:"free_to_start"`
	Title           string `json:"title"`
	System          string `json:"system"`
	ID              string `json:"id"`
	CAPrice         string `json:"ca_price"`
	NumberOfPlayers string `json:"number_of_players"`
	NSUID           string `json:"nsuid"`
	EShopPrice      string `json:"eshop_price"`
	FrontBoxArt     string `json:"front_box_art"`
	GameCode        string `json:"game_code"`
	BuyOnline       string `json:"buyonline"`
}

func eShopUSHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	var allGames []Game
	total := 1
	for len(allGames) < total {
		games, newTotal, err := getEShopUSGames(len(allGames))
		if err != nil {
			log.Println(err.Error())
			http.Error(w, err.Error(), 500)
		}
		allGames = uniqueGames(append(allGames, games...))
		total = newTotal
		fmt.Printf("total: %d, current: %d\n", total, len(allGames))
	}
	json.NewEncoder(w).Encode(allGames)
}

func getEShopUSGames(offset int) ([]Game, int, error) {
	url := "https://www.nintendo.com/json/content/get/filter/game?system=switch&sort=title&direction=asc&shop=ncom&offset=" + strconv.Itoa(offset)
	resp, err := http.Get(url)
	if err != nil {
		return nil, 0, err
	}

	defer resp.Body.Close()

	var d struct {
		Filter struct {
			Total int `json:"total"`
		} `json:"filter"`
		Games struct {
			Game []Game `json:"game"`
		} `json:"games"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&d); err != nil {
		return nil, 0, err
	}

	return d.Games.Game, d.Filter.Total, nil
}

// uniqueGames : Function that filters a slice of Game by 'Slug' and returns slice of unique Game
func uniqueGames(input []Game) []Game {
	u := make([]Game, 0, len(input))
	m := make(map[string]bool)
	for _, val := range input {
		if _, ok := m[val.Slug]; !ok {
			m[val.Slug] = true
			u = append(u, val)
		}
	}
	return u
}

package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/guptarohit/asciigraph"
)

//Item is a WoW Classic Item
type Item struct {
	Server        string      `json:"server"`
	ItemID        int         `json:"itemId"`
	Name          string      `json:"name"`
	UniqueName    string      `json:"uniqueName"`
	Icon          string      `json:"icon"`
	Tags          []string    `json:"tags"`
	RequiredLevel int         `json:"requiredLevel"`
	ItemLevel     int         `json:"itemLevel"`
	SellPrice     int         `json:"sellPrice"`
	VendorPrice   interface{} `json:"vendorPrice"`
	Tooltip       []struct {
		Label  string `json:"label"`
		Format string `json:"format,omitempty"`
	} `json:"tooltip"`
	ItemLink string `json:"itemLink"`
	Stats    struct {
		LastUpdated time.Time `json:"lastUpdated"`
		Current     struct {
			MarketValue     int `json:"marketValue"`
			HistoricalValue int `json:"historicalValue"`
			MinBuyout       int `json:"minBuyout"`
			NumAuctions     int `json:"numAuctions"`
			Quantity        int `json:"quantity"`
		} `json:"current"`
		Previous struct {
			MarketValue     int `json:"marketValue"`
			HistoricalValue int `json:"historicalValue"`
			MinBuyout       int `json:"minBuyout"`
			NumAuctions     int `json:"numAuctions"`
			Quantity        int `json:"quantity"`
		} `json:"previous"`
	} `json:"stats"`
}

//Prices is a struct showing current and old prices of Item
type Prices struct {
	Slug       string `json:"slug"`
	ItemID     int    `json:"itemId"`
	Name       string `json:"name"`
	UniqueName string `json:"uniqueName"`
	Timerange  int    `json:"timerange"`
	Data       []struct {
		MarketValue float64   `json:"marketValue"`
		MinBuyout   float64   `json:"minBuyout"`
		Quantity    int       `json:"quantity"`
		ScannedAt   time.Time `json:"scannedAt"`
	} `json:"data"`
}

//Search is an array of structs containing responds from a search query
type Search []struct {
	ItemID     int    `json:"itemId"`
	Name       string `json:"name"`
	UniqueName string `json:"uniqueName,omitempty"`
	ImgURL     string `json:"imgUrl"`
}

func (i Item) print() {
	if i.Name == "" {
		fmt.Println("No Item Found!")
		return
	}

	for _, tip := range i.Tooltip {
		if strings.HasPrefix(tip.Label, "Sell Price:") {
			continue
		}

		fmt.Println(tip.Label)
	}
	fmt.Println("")

	fmt.Println("Vendor Price: " + intToWoWString(i.SellPrice))

	TimeAgo := time.Since(i.Stats.LastUpdated)
	fmt.Println("Last Updated: " + i.Stats.LastUpdated.String() + ", " + TimeAgo.String() + " ago.")
	fmt.Println("")

	fmt.Println("Market Value: " + intToWoWString(i.Stats.Current.MarketValue) + " (14 day weighted avg, favoring last 3 days)")
	fmt.Println("Historic Value: " + intToWoWString(i.Stats.Current.HistoricalValue) + " (60 day avg)")
	fmt.Println("Current Minimum Buyout: " + intToWoWString(i.Stats.Current.MinBuyout))
	fmt.Println("Num Auctions: " + strconv.Itoa(i.Stats.Current.NumAuctions))
	fmt.Println("Quantity: " + strconv.Itoa(i.Stats.Current.Quantity))
	fmt.Println("")
}

func (p Prices) print() {
	if p.Name == "" {
		return
	}

	if len(p.Data) == 0 {
		fmt.Println("No market data found!")
		return
	}

	data := []float64{}
	dat2 := []float64{}

	scale, graphUnit := graphScale(p.Data[0].MarketValue)

	for _, dat := range p.Data {
		data = append(data, dat.MinBuyout/scale)
		dat2 = append(dat2, dat.MarketValue/scale)
	}

	graph := asciigraph.Plot(data, asciigraph.Height(10), asciigraph.Width(50))
	fmt.Println("Min Buyout Graph " + graphUnit)
	fmt.Println(graph)

	graph2 := asciigraph.Plot(dat2, asciigraph.Height(10), asciigraph.Width(50))
	fmt.Println("Market Value Graph " + graphUnit)
	fmt.Println(graph2)

	TimeAgo := time.Since(p.Data[0].ScannedAt)
	days := fmt.Sprintf("%f", TimeAgo.Hours()/24)
	fmt.Println("Oldest Data Point: " + p.Data[0].ScannedAt.String() + ", " + TimeAgo.String() + " ago. => " + days + "days")

}

func graphScale(price float64) (float64, string) {
	if price >= 10000 {
		return 10000, "(g):"
	} else if price <= 100 {
		return 0, "(c):"
	} else {
		return 100, "(s):"
	}
}

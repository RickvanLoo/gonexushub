package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/guptarohit/asciigraph"
)

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

func (i Item) print() {
	if i.Name == "" {
		fmt.Println("Item not found!")
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

	fmt.Println("Market Value:" + intToWoWString(i.Stats.Current.MarketValue))
	fmt.Println("Minimum Buyout:" + intToWoWString(i.Stats.Current.MinBuyout))
	fmt.Println("Num Auctions:" + strconv.Itoa(i.Stats.Current.NumAuctions))
	fmt.Println("Quantity:" + strconv.Itoa(i.Stats.Current.Quantity))
	fmt.Println("")
}

func (p Prices) print() {
	if p.Name == "" {
		return
	}

	data := []float64{}

	for _, dat := range p.Data {
		data = append(data, dat.MarketValue)
	}

	graph := asciigraph.Plot(data, asciigraph.Height(10), asciigraph.Width(50))

	fmt.Println(graph)
}

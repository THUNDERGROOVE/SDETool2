package graph

import (
	"code.google.com/p/plotinum/plot"
	"code.google.com/p/plotinum/plotter"
	"code.google.com/p/plotinum/vg"
	"fmt"
	"github.com/THUNDERGROOVE/SDETool2/log"
	"github.com/THUNDERGROOVE/SDETool2/market"
	"github.com/lucasb-eyer/go-colorful"
	"math/rand"
	"sort"
	"time"
)

type Cache struct {
	Data map[string]map[int]*Data
}
type Data struct {
	UnitsSold int
	Month     int
	Year      int
}

var MarketCache Cache

func init() {
	log.Info("Graph imported.  Generating MarketCache.")
	MarketCache = Cache{make(map[string]map[int]*Data, 0)}
}

func doCache(data map[string]map[string]market.MarketData) {
	for name, set := range data {
		for _, v := range set {
			for _, vv := range v.Items {
				d, err := time.Parse("2006-1-2T03:04:05", vv.Date)
				if err != nil {
					fmt.Println("Error parsing date from market data", err.Error())
					continue
				}
				if time.Now().Sub(d).Hours() >= 8040 { // Ignore data set.  8040 is how many hours are in 11 months
					log.Info("Ignoring old data set.", time.Now().Sub(d).Hours(), "hours old.")
					continue
				}
				if MarketCache.Data == nil {
					MarketCache.Data = make(map[string]map[int]*Data, 0)
				}
				if MarketCache.Data[name] == nil {
					MarketCache.Data[name] = make(map[int]*Data, 0)
				}
				if MarketCache.Data[name][int(d.Month())] == nil {
					MarketCache.Data[name][int(d.Month())] = &Data{}
				}
				MarketCache.Data[name][int(d.Month())].UnitsSold += int(vv.Volume)
				MarketCache.Data[name][int(d.Month())].Month = int(d.Month())
				MarketCache.Data[name][int(d.Month())].Year = int(d.Month())

			}
		}
	}
}

func PlotSuitData(data map[string]map[string]market.MarketData) { //ew lol
	log.Info("Proccessing Market data")

	doCache(data)

	p, err := plot.New()
	if err != nil {
		log.LogError("Error creating log", err.Error())
		return
	}
	p.Title.Text = "Suit Market Data"
	p.X.Label.Text = "Month"
	p.Y.Label.Text = "Suits Sold"

	rand.Seed(time.Now().UnixNano())
	c, _ := colorful.SoftPalette(len(MarketCache.Data))
	var ci int
	for name, v := range MarketCache.Data {
		var i int
		points := make(plotter.XYs, len(v))

		// Sort keys so we can have the lines in order
		var keys []int
		for k := range v {
			keys = append(keys, k)
		}
		sort.Ints(keys)

		for _, k := range keys {
			points[i].X = float64(k)
			points[i].Y = float64(MarketCache.Data[name][k].UnitsSold)
			i++
		}

		line, err := plotter.NewLine(points)
		if err != nil {
			log.LogError("Error creating line in plot", err.Error())
			continue
		}

		line.Color = c[ci]
		line.Width = vg.Points(5)
		p.Add(line)
		p.Legend.Add(name, line)
		ci++
	}
	if err := p.Save(16, 9, "graph.png"); err != nil {
		log.LogError("Error saving graph as png.")
	}
}

func BarSuitData(data map[string]map[string]market.MarketData) { //ew lol
	log.Info("Proccessing Market data")
	doCache(data)
	p, err := plot.New()
	if err != nil {
		log.LogError("Error creating log", err.Error())
		return
	}

	rand.Seed(time.Now().UnixNano())
	c, _ := colorful.SoftPalette(len(MarketCache.Data))
	var ci int
	for name, v := range MarketCache.Data {
		val := plotter.Values{}
		for _, vv := range v {
			val = append(val, float64(vv.UnitsSold))
		}
		bar, err := plotter.NewBarChart(val, vg.Points(5))
		if err != nil {
			log.LogError("Error creating bar for chart")
		}
		bar.Color = c[ci]
		bar.Offset = vg.Points(float64(5 * ci))
		p.Add(bar)
		p.Legend.Add(name, bar)

		ci++
	}
	p.NominalX("January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December")
	if err := p.Save(16, 9, "graph.png"); err != nil {
		log.LogError("Error saving graph as png.")
	}
}

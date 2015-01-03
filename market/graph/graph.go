package graph

import (
	"code.google.com/p/plotinum/plot"
	"code.google.com/p/plotinum/plotter"
	"code.google.com/p/plotinum/vg"
	"fmt"
	"github.com/THUNDERGROOVE/SDETool2/log"
	"github.com/THUNDERGROOVE/SDETool2/market"
	"github.com/lucasb-eyer/go-colorful"
	"math"
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

func save(p *plot.Plot, count int, isBar bool) {
	baseW := float64(8)
	baseH := float64(4.5)
	log.Info("Saving graph with options count", count, "isBar", isBar, "using size ", saveScale(baseW, count), saveScale(baseH, count))

	if err := p.Save(saveScale(baseW, count), saveScale(baseH, count), "graph.png"); err != nil {
		log.LogError("Error saving graph as png.")
	}
}

func saveScale(f float64, c int) float64 {
	return math.Sqrt(f*float64(c)) * float64(2)
}

func doCache(data map[string]map[string]market.MarketData) {
	for name, set := range data {
		for _, v := range set {
			for _, vv := range v.Items {
				d, err := time.Parse("2006-01-02T03:04:05", vv.Date)
				if err != nil {
					fmt.Println("Error parsing date from market data", err.Error())
					continue
				}
				if time.Now().Sub(d).Hours() >= 4380 { // Ignore data set.  4380 is half of a year
					//log.Info("Ignoring old data set.", time.Now().Sub(d).Hours(), "hours old.", d.String())
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
	save(p, len(MarketCache.Data), false)
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
	p.Legend.Top = true
	p.NominalX("January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December")
	save(p, len(MarketCache.Data), true)
}

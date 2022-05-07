package customalgorithms

import (
	"fmt"
	"math"
	"sort"

	Clustering "bitbucket.org/sjbog/go-dbscan"
)

type (
	S_R struct {
		Val float64
		Idx int
	}
	ImportantPoints struct {
		strength int
		avg      float64
		min      float64
		max      float64
	}
)

func Supports(timeseries []float64, leading int) []S_R {
	supports := make([]S_R, 0)
	timeseries = buffered(timeseries)
	for i := 1; i < len(timeseries)-1; i++ {
		var lowerbound int = 0
		var upperbound int = len(timeseries) - 1
		if i-int(leading) > 0 {
			lowerbound = i - int(leading)
		}
		if i+leading < len(timeseries) {
			upperbound = i + leading
		}
		isMinima := smallerThanAll(timeseries[i+1:upperbound], timeseries[i]) && smallerThanAll(timeseries[lowerbound:i-1], timeseries[i])
		if isMinima {
			supports = append(supports, S_R{Val: timeseries[i], Idx: i})
		}
	}
	return supports
}

func Resistances(timeseries []float64, leading int) []S_R {
	resistances := make([]S_R, 0)
	timeseries = buffered(timeseries)
	for i := 1; i < len(timeseries)-1; i++ {
		var lowerbound int = 0
		var upperbound int = len(timeseries) - 1
		if i-int(leading) > 0 {
			lowerbound = i - int(leading)
		}
		if i+leading < len(timeseries) {
			upperbound = i + leading
		}
		isMaxima := greaterThanAll(timeseries[i+1:upperbound], timeseries[i]) && greaterThanAll(timeseries[lowerbound:i-1], timeseries[i])
		if isMaxima {
			resistances = append(resistances, S_R{Val: timeseries[i], Idx: i})
		}
	}
	return resistances
}

func GetImportantLevels(timeseries []float64, leading int, closeness float64) []ImportantPoints {
	supports := Supports(timeseries, leading)
	resistances := Resistances(timeseries, leading)
	data := []Clustering.ClusterablePoint{}
	for _, support := range supports {
		data = append(data, &Clustering.NamedPoint{
			Name:  fmt.Sprintf("%d", support.Idx),
			Point: []float64{support.Val},
		})
	}
	for _, resistance := range resistances {
		data = append(data, &Clustering.NamedPoint{
			Name:  fmt.Sprintf("%d", resistance.Idx),
			Point: []float64{resistance.Val},
		})
	}
	clusterer := Clustering.NewDBSCANClusterer(closeness, 2)
	importantPoints := clusterer.Cluster(data)
	realSupports := []ImportantPoints{}
	for _, val := range importantPoints {
		var avg float64 = 0
		var min float64 = math.Inf(1)
		var max float64 = math.Inf(-1)
		for _, d := range val {
			value := d.GetPoint()[0]
			avg += value
			if value < min {
				min = value
			}
			if value > max {
				max = value
			}
		}
		avg = avg / float64(len(val))
		realSupports = append(realSupports, ImportantPoints{
			strength: len(val),
			avg:      avg,
			min:      min,
			max:      max,
		})
	}
	sort.SliceStable(realSupports, func(i, j int) bool {
		return realSupports[i].strength > realSupports[j].strength
	})
	return realSupports
}

func buffered(timeseries []float64) []float64 {
	toRet := make([]float64, 0)
	toRet = append(toRet, 0)
	toRet = append(toRet, timeseries...)
	toRet = append(toRet, 0)
	return toRet
}

func smallerThanAll(subSeries []float64, val float64) bool {
	isSmallest := true
	for i := 0; i < len(subSeries); i++ {
		if subSeries[i] < val {
			isSmallest = false
			break
		}
	}
	return isSmallest
}

func greaterThanAll(subSeries []float64, val float64) bool {
	isGreatest := true
	for i := 0; i < len(subSeries); i++ {
		if subSeries[i] > val {
			isGreatest = false
			break
		}
	}
	return isGreatest
}

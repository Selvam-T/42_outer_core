/* stats + normalization */

package model

import (
	"math"
)

type MileageStats struct {
	MeanMileage_train float64 `json:"meanMileage_train"`
 	StdMileage_train float64 `json:"stdMileage_train"`
}

func Mean(data [][]float64) (float64) {
	n := len(data)
	sum := 0.0
	
	for _, row:= range data {
		sum += row[0] //sum of mileage
	}
	return sum / float64(n)
}

func StdDev(data [][]float64, mean float64) (float64) {
	n := len(data)
	sum := 0.0
	
	for _, row:= range data {
		sum += square(row[0] - mean)
	}
	return math.Sqrt( 1.0 / float64(n) * sum)
}

func square(x float64) float64 {
	return x * x
}

func NormalizeMileage(mileage float64, mean float64, std float64) float64 {
	
	if (std <= 0) { return 0.0 }
	
	return (mileage - mean) / std
}

func NormalizeData(data [][]float64, jsonFile string) ([][]float64, error) {
	
	var m MileageStats
	err := LoadJson(jsonFile, &m) 
	if err != nil { return nil, err }
	
	for _,row := range data {
		row[0] = NormalizeMileage(row[0], m.MeanMileage_train, m.StdMileage_train)
	}

	return data, nil		
}

func ComputeTrainStats(jsonFile string, data [][]float64) error {
	var m MileageStats
	m.MeanMileage_train = Mean(data)
	m.StdMileage_train = StdDev(data, m.MeanMileage_train)
	
	e := SaveJson(jsonFile, &m)
	if e != nil { return e }
	
	return nil		
}



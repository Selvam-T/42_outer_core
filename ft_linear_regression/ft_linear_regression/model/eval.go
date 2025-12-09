package model

import (
	"math"
)

func sumSquaredErrors(data [][]float64, n *Model) (float64, error) {
	sse := 0.0
	
	for _, row := range data {		
		mileage := row[0]				
		price := row[1]
				
		pe, err := PredictionError(mileage, price, n)
		if err != nil { return -1, err }
		
		sse += pe * pe	// Σ(ŷ − y)²
	}
	return sse, nil
}

func totalSumErrors(data [][]float64) (float64) {
	price := 0.0
	sst := 0.0
	
	priceAvg := Mean(data)
	
	for _, row := range data {					
		price = row[1]				
		diff := price - priceAvg // Σ(yᵢ − ȳ)²	
		sst += (diff) * (diff)
	}
	return sst
}

func RMSE(data [][]float64, n *Model) (float64, error) {

	sse, err:= sumSquaredErrors(data, n)
	if err != nil {return -1, err}
	
	m := float64 (len(data))
	return math.Sqrt(sse/m), nil
}

func Rsquare(data [][]float64, n *Model) (float64, error) {

	sse, err:= sumSquaredErrors(data, n)
	if err != nil {return -1, err}
	
	sst := totalSumErrors(data)
	
	return 1 - (sse / sst), nil
}

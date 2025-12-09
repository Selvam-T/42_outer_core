/* Predict the price of a car for a given mileage  */

package main

import (
	"fmt"
	"strconv"
	"ft_linear_regression/model"
	"errors"
)

func loadJsons(cfg model.Config) (model.Model, model.MileageStats, error) {

	/* Extract model.json */
	var m model.Model
	err := model.LoadJson(cfg.ModelJson, &m) 
	if err != nil { return model.Model{}, model.MileageStats{}, err}
	
	/* Extract mileagestats.json */
	var n model.MileageStats
	err = model.LoadJson(cfg.MileageJson, &n)
	if err != nil { return model.Model{}, model.MileageStats{}, err}
	
	return m, n, nil
}

func readMileage() (float64, error) {

	/* read user input */
	var input string
	fmt.Print("Enter Mileage: ")
	fmt.Scanln(&input)
	
	/*normalize mileage n := (row[0] - mean) / std */
	mileage, err := strconv.ParseFloat(input, 64)
	if err != nil { return 0, err}
			
	/* Validate mileage */
	if mileage >= 500000 || mileage < 0 { return 0, errors.New("mileage > 500000")}
	
	return mileage, nil
}

func main() {
	
	/* Extract config.json*/ 
	cfg, err := model.LoadConfig("config.json")
	if err != nil { fmt.Println("Error reading config: ", err); return }
	
	m, n, err := loadJsons(cfg)
	if err != nil { fmt.Println("Error reading Json: ", err); return }
	
	mileage, err := readMileage()
	if err != nil { fmt.Println("Input Error: ", err); return }
	
	mileage = model.NormalizeMileage(mileage, n.MeanMileage_train, n.StdMileage_train)
	
	/* estimate price */
	price, err := model.EstimatePrice(mileage, &m)
	if err != nil { fmt.Println("Error estimating price: ", err); return }
	
	fmt.Printf("Estimated Price: %.2f\n", price)
}

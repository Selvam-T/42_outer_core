/* 	Perform linear regression on the model (hypothesis).
	Gradient descent method used to find the best θ₀ and θ₁ for that line. */

package main

import (
	"fmt"
	"ft_linear_regression/model"
	"math/rand"
	"time"
)

/* randomly split data into 80:20 train:test data */
func splitData(data [][]float64, ratio float64) ([][]float64, [][]float64, error) {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(data), func(i, j int) {
		data[i], data[j] = data[j], data[i]
	})
	
	//split by ratio
	split := int(ratio * float64(len(data)))
	
	trainData := data[:split]
	testData := data[split:]
	
	return trainData, testData, nil
}

func main() {
		
	/* Load configuration */ 	
	cfg, err := model.LoadConfig("config.json")
	if err != nil { fmt.Println("Error loading configuration parameters: ", err); return }	

	/* Load model parameters */
	var m model.Model
	err = model.LoadJson(cfg.ModelJson, &m) 
	if err != nil { fmt.Println("Error loading model parameters: ", err); return }
	
	/* Load data from file, convert string to float */
	parsedData, err := model.LoadParseCSVFlt64(cfg)
	if err != nil { fmt.Println("Error parsing data: ", err); return }

	/* split train / test data */
	trainData, testData, err := splitData(parsedData, cfg.TrainRatio)
	if err != nil { fmt.Println("Error splitting data: ", err); return }
	
	/* compute mileage stats on train data */
	err = model.ComputeTrainStats(cfg.MileageJson, trainData)
	if err != nil { fmt.Println("Error computing mileage stats: ", err); return }
	
	/* normalize train data */
	normTrainData, err := model.NormalizeData(trainData, cfg.MileageJson)
	if err != nil { fmt.Println("Error normalizing train data: ", err); return }
		
	/* Find the optimal set of model parameters θ₀ and θ₁ on train data */
	err = model.GradientDescent(cfg, normTrainData, &m)
	if err != nil { fmt.Println("Error gradient descent: ", err); return }

	/* normalize test data */
	normTestData, err := model.NormalizeData(testData, cfg.MileageJson)
	if err != nil { fmt.Println("Error normalizing test data: ", err); return }

	/* RMSE on test Data */
	rmse, err := model.RMSE(normTestData, &m)
	if err != nil { fmt.Println("Error rmse: ", err); return }
	fmt.Printf("RMSE on testData: $ %.2f\n", rmse)
	fmt.Printf("It means my predictions on average is far off by $ %.2f\n", rmse)
	
	/*  R² */
	rsqr, err := model.Rsquare(normTestData, &m)
	if err != nil { fmt.Println("Error rsquare: ", err); return }
	fmt.Printf("\nR² is %.2f\n", rsqr)
	
	/* plot */
	err = model.Plot(normTrainData, &m, cfg.GraphOut)
	if err != nil { fmt.Println("Error plotting graph: ", err); return }
	
	return
}

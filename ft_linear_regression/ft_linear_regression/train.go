/* 	Perform linear regression on the model (hypothesis).
	Gradient descent method used to find the best θ₀ and θ₁ for that line. */

package main

import (
	"fmt"
	"ft_linear_regression/model"
	"encoding/json"
)

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
	
	/* compute mileage stats */
	err = model.ComputeTrainStats(cfg.MileageJson, parsedData)
	if err != nil { fmt.Println("Error computing mileage stats: ", err); return }
	
	/* normalize data */
	normData, err := model.NormalizeData(parsedData, cfg.MileageJson)
	if err != nil { fmt.Println("Error normalizing data: ", err); return }
	
	/* print model parameter before training */
	b, _ := json.MarshalIndent(m, "", "  ")
	fmt.Println("Model parameters before Gradient Descent: ", string(b))
	
	/* Find the optimal set of model parameters θ₀ and θ₁ on data*/
	err = model.GradientDescent(cfg, normData, &m)
	if err != nil { fmt.Println("Error in gradient descent: ", err); return }
	
	/* print model parameter after training */
	b, _ = json.MarshalIndent(m, "", "  ")
	fmt.Println("Model parameters after Gradient Descent: ", string(b))
	fmt.Println("Training data is completed !")
}

/* hypothesis is the formula model uses to make predictions.*/

/* Parameters explained with example:
	θ₀ = 20000 (base price when mileage = 0) 
		    also the y-intercept of the line
	θ₁ = -0.05 (price drops by 0.05 per km)  
		    also the slope of the line */

/* model + prediction */

package model

type Model struct {
	Theta0 float64 `json:"theta0"` /* how I represent Theata0 in model.json */
	Theta1 float64 `json:"theta1"`
}

/* hypothesis */
func EstimatePrice(mileage float64, m *Model) (float64, error) {
	return m.Theta0 + (m.Theta1 * mileage), nil
}

/* for Gradient calculation*/
func PredictionError(mileage float64, price float64, m *Model) (float64, error) {
	
	estPrice, err := EstimatePrice(mileage, m)
	if err != nil { return -1, err }
	return  estPrice - price, nil
}




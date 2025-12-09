/* 
	gradient0 = (1/m) * Σ (prediction − actual)
	gradient1 = (1/m) * Σ (prediction − actual) * mileage
	TmpTheta0 = lr * gradient0
	TmpTheta1 = lr * gradient1 */

/* training / optimization */

package model

import (
	"math"
	//"fmt"
)

func TmpTheta1(lr float64, data [][]float64, n *Model) (float64, error) {
	gradient1 := 0.0
	
	for _, row := range data {		
		mileage := row[0]				
		price := row[1]
		
		pe, err := PredictionError(mileage, price, n)
		if err != nil { return -1, err }
		
		gradient1 += (pe * mileage)
	}
	
	m := float64 (len(data))
	
	return lr * (1/m) * gradient1, nil
}

func TmpTheta0(lr float64, data [][]float64, n *Model) (float64, error) {
	gradient0 := 0.0
	
	for _, row := range data {		
		mileage := row[0]				
		price := row[1]
				
		pe, err := PredictionError(mileage, price, n)
		if err != nil { return -1, err }
		
		gradient0 += pe
	}
	
	m:= float64 (len(data))
	
	return lr * (1/m) * gradient0, nil	
}

/* Iteratively adjusting the model parameters */
func GradientDescent(cfg Config, data [][]float64, n *Model) (error){
	prevLoss := math.Inf(1) // +infinity so first check always fails
	i := 0

	for i < cfg.MaxIterations { //stop rule 1
		tmpTheta0, err := TmpTheta0(cfg.LearningRate, data, n)
		if err != nil { return err}
		
		tmpTheta1, err := TmpTheta1(cfg.LearningRate, data, n)
		if err != nil { return err}	
		
		/* gradient descent update*/
		n.Theta0 -= tmpTheta0
		n.Theta1 -= tmpTheta1
		
		/*update new parameters */
		SaveJson(cfg.ModelJson, n)
		if err != nil { return err}
		
		/* change in loss function */
		curLoss, err := LossConvergence(cfg.LearningRate, data, n)
		if err != nil { return err}
		
		//stop rule 2 Loss Convergence
		if math.Abs(curLoss - prevLoss) < cfg.Tolerance {  
			//fmt.Println("Stop rule 2 reached at iteration [",i, "]") /*** DELETE ***/
			//fmt.Println("CurLoss: ", curLoss, "\nPrevLoss: ", prevLoss, "\nDiff: ", curLoss - prevLoss) /*** DELETE ***/
			break 
		}
		prevLoss = curLoss
		
		i++
	}
	return nil
}

/* loss-convergence after each gradient-descent update, using the new θ₀ and θ₁ */
func LossConvergence(lr float64, data [][]float64, n *Model) (float64, error) {
	sse, err:= sumSquaredErrors(data, n)
	if err != nil {return -1, err}
	
	m:= float64 (len(data))
	
	return sse / (2 * m), nil /* J = (1 / 2m) Σ(ŷ − y)² */	
}


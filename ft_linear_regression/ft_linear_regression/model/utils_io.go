/* IO */

package model

import (
	"os"
	"encoding/json"
	"encoding/csv"
	"strconv"
	"fmt"
)

type Config struct {
	LearningRate	float64 `json:"learningRate"`
	DataFile	string `json:"dataFile"`
	Tolerance	float64 `json:"tolerance"`
	TrainRatio	float64 `json:"trainRatio"`
	MaxIterations	int `json:"maxIterations"`
	ModelJson	string `json:"modelJson"`
	MileageJson	string `json:"mileageJson"`
	GraphOut	string `json:"graphOut"`
}

func LoadJson(modelPath string, v interface{}) error {
	
	b, err := os.ReadFile(modelPath)
	if err != nil { return err }
	
	return json.Unmarshal(b, v)
}

func SaveJson(modelPath string, v interface{}) error{
	
	jsonData, err := json.MarshalIndent(v, "", " ")
	if err != nil { return err}
	
	return os.WriteFile(modelPath, jsonData, 0666)
}

func LoadConfig(path string) (Config, error) {
	b, err := os.ReadFile(path)
	if err != nil { return Config{}, err}

	var cfg Config
	if err:= json.Unmarshal(b, &cfg); err != nil { return Config{}, err }
	
	return cfg, nil
}

/* Load sample data from .csv */
func loadData(path string) ([][]string, error) {
	file, err := os.Open(path)
	if err != nil { return nil, err}
	defer file.Close()
	
	reader := csv.NewReader(file)
	record, err := reader.ReadAll()
	if err != nil { return nil, err }
	
	return record, nil
}

/* clean and convert price and mileage from string to float64 */
func convertValidToFloat(records [][]string) ([][]float64, error) {

	rows := len(records) - 1
	if (rows <= 0) { return nil, fmt.Errorf("no data rows")}
	
	var data [][]float64
	
	for i, record := range records {
		if i == 0 { continue } /* skip header */
		
		mileage, err1 := strconv.ParseFloat(record[0], 64)		
		price, err2 := strconv.ParseFloat(record[1], 64)
		
		// drop invalid rows
		if err1 != nil || err2 != nil || mileage < 0 || price < 0 { continue }
		
		data = append(data, []float64{mileage, price})
	}
	
	if len(data) == 0 { return nil, fmt.Errorf("no valid rows after cleaning")}
	return data, nil
}

//1. load data + convert to str --> before split
func LoadParseCSVFlt64(cfg Config) ([][]float64, error){
	sampleData, err :=  loadData(cfg.DataFile)
	if err != nil { return nil , err }
	
	cleanData, err := convertValidToFloat(sampleData)
	if err != nil { return nil , err }
	
	return cleanData, nil
}



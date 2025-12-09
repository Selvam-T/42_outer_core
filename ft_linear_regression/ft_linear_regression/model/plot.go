/* 1. Plotting the data into a graph to see their repartition. */
/* 2. Plotting the line resulting from your linear regression into the same graph */

package model

import (
	"fmt"
	"image/color"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

func MinMax(data [][]float64) (float64, float64) {
    xMin := data[0][0]
    xMax := data[0][0]

    for _, row := range data {
        x := row[0]
        if x < xMin { xMin = x }
        if x > xMax { xMax = x }
    }
    return xMin, xMax
}

func Plot(trainData [][]float64, n *Model, outFile string) (error){

	// scatter plot on train data
	pts := make(plotter.XYs, len(trainData))
	for i, d := range trainData {
		pts[i].X = d[0]
		pts[i].Y = d[1]
	}
	
	scatterData := pts

	p := plot.New()
	p.Title.Text = "ft_linear_regression"
	p.Title.TextStyle.Font.Size = vg.Points(24)
	p.X.Label.Text = "Mileage"
	p.X.Label.TextStyle.Font.Size = vg.Points(20)
	p.Y.Label.Text = "Price"
	p.Y.Label.TextStyle.Font.Size = vg.Points(20)
	p.Add(plotter.NewGrid())

	s, err := plotter.NewScatter(scatterData)
	if err != nil { return err }
	
	s.GlyphStyle.Shape = draw.CircleGlyph{}
	s.GlyphStyle.Color = color.RGBA{R: 255, A: 255}
	s.GlyphStyle.Radius = vg.Points(4)
	
	// regression line
	
	xMin, xMax := MinMax(trainData)

	yMin := n.Theta0 + (n.Theta1 * xMin)
	yMax := n.Theta0 + (n.Theta1 * xMax)
	
	linePts := plotter.XYs {
		{X: xMin, Y: yMin},
		{X: xMax, Y: yMax},
	}
	
	l, err := plotter.NewLine(linePts)
	if err != nil { return err }
	
	l.LineStyle.Width = vg.Points(1)
	l.LineStyle.Dashes = []vg.Length{vg.Points(5), vg.Points(5)}
	l.LineStyle.Color = color.RGBA{B: 255, A: 255}

	p.Add(s, l)

	err = p.Save(1000, 1000, outFile)
	if err != nil { return err }
	
	fmt.Println("\nGraph saved in file: ", outFile)
	
	return nil
}


package lib

import (
	"os/exec"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

func DrawGraph(gainPts plotter.XYs, capitalPts plotter.XYs, operations plotter.XYs) {
	p, err := plot.New()
	if err != nil {
		panic(err)
	}

	p.Title.Text = "gain"
	p.X.Label.Text = "time"
	p.X.Tick.Marker = plot.TimeTicks{Format: "2006-01-02\n15:04"}
	p.Y.Label.Text = "$"

	err = plotutil.AddLinePoints(p,
		"Gain", gainPts,
		"Capital", capitalPts)
	if err != nil {
		panic(err)
	}

	// Save the plot to a PNG file.
	if err := p.Save(12*vg.Inch, 12*vg.Inch, "points.png"); err != nil {
		panic(err)
	}
	exec.Command("rundll32", "url.dll,FileProtocolHandler", "points.png").Start()
}

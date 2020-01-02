package lib

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"github.com/willfr/goback/model"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/willfr/goback/globals"
)

func SaveRun() {
	f, _ := os.OpenFile("C:\\Users\\Guillaume\\Desktop\\stocks\\run.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	defer f.Close()
	f.WriteString(strings.Join(os.Args, " ") + fmt.Sprintf(" GAIN: %.2f ", globals.Total) + fmt.Sprintf("CAPITAL: %.2f \r\n", globals.Capital))
	f.Sync()
}

func GenerateHistory() {
	sort.Slice(globals.History, func(i, j int) bool {
		return globals.History[i].Date.Before(globals.History[j].Date)
	})
	f2, _ := os.OpenFile("C:\\Users\\Guillaume\\Desktop\\stocks\\history.txt", os.O_CREATE, 0600)
	defer f2.Close()
	f2.Truncate(0)
	f2.Seek(0, 0)
	for _, h := range globals.History {
		f2.WriteString(h.Name + " " + h.Date.Format(time.RFC3339) + " " + fmt.Sprintf("%.0f", h.Quantity) + " " + fmt.Sprintf("%f", h.Price) + " \r\n")
	}
	f2.Sync()
}

type binDataPoint struct {
	Date   int64
	Price  float64
	Low    float64
	Volume float64
}

func WriteToBin(writer *bufio.Writer, dp *model.DataPoint) {
	binary.Write(writer, binary.LittleEndian, binDataPoint{Date: dp.Date.Unix(), Price: dp.Price, Low: dp.Low, Volume: dp.Volume})
}

func ReadFromBin(reader *bufio.Reader) (*model.DataPoint, error) {
	var bdp binDataPoint
	err:= binary.Read(reader, binary.LittleEndian, &bdp)
	dp := &model.DataPoint{Date: time.Unix(bdp.Date, 0).UTC(), Price: bdp.Price, Low: bdp.Low, Volume: bdp.Volume}
	return dp, err
}
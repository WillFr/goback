package lib

import (
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/willfr/goback/model"
)

func ParseLine(line string) model.DataPoint {
	datePart := line[:16]
	data := strings.Split(line[17:], ",")
	price, _ := strconv.ParseFloat(data[0], 64)
	volume, _ := strconv.ParseFloat(data[len(data)-1], 64)
	datetime, _ := time.Parse("01/02/2006,15:04", datePart)
	return model.DataPoint{Date: datetime, Price: price, Volume: volume}
}

func ListTickers(root string) []string {
	var tickers []string
	filepath.Walk(root, func(filePath string, info os.FileInfo, err error) error {
		if path.Ext(filePath) == ".txt" {
			name := strings.TrimSuffix(strings.TrimPrefix(filePath, root), ".txt")
			tickers = append(tickers, name)
		}
		return nil
	})
	return tickers
}

func OpenTicker(name string) *os.File {
	filePath := "C:\\Users\\Guillaume\\Desktop\\stocks\\intraday\\" + name + ".txt"
	file, _ := os.Open(filePath)
	return file
}

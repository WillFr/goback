package lib

import (
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/willfr/goback/model"
)

func ParseLine(line []byte) model.DataPoint {
	datetime := parseDate(line)
	oi := 17
	i := 18
	for line[i] != ',' {
		i++
	}
	price := parseFloat(line[oi:i])
	i = len(line) - 1
	oi = i + 1
	for line[i] != ',' {
		i--
	}
	volume := parseFloat(line[i+1:oi])
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

func parseDate(date []byte) time.Time {
	year := (((int(date[6])-'0')*10+int(date[7])-'0')*10+int(date[8])-'0')*10 + int(date[9]) - '0'
	month := time.Month((int(date[0])-'0')*10 + int(date[1]) - '0')
	day := (int(date[3])-'0')*10 + int(date[4]) - '0'
	hour := (int(date[11])-'0')*10 + int(date[12]) - '0'
	minute := (int(date[14])-'0')*10 + int(date[15]) - '0'
	return time.Date(year, month, day, hour, minute, 0, 0, time.UTC)
}

func parseFloat(str []byte) float64 {
	res := 0
	div := 0
	for _,  c := range str  {
		if c == '.' {
			div = 1
			continue
		}
		res = res*10 + (int(c)   - '0')
		div *= 10
	}
	div = div 
	if div > 0{
		return float64(res)   / float64(div)
	}
	return float64(res)
}


package lib

import (
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/willfr/goback/model"
)

func ParseLine(line []byte) *model.DataPoint {
	if len(line) <= 20 {
		return nil
	}
	datetime := parseDate(line)

	oi := 17
	i := 18
	//open
	for line[i] != ',' {
		i++
	}
	price := parseFloat(line[oi:i])

	//high
	i++
	for line[i] != ',' {
		i++
	}

	//low
	i++
	oi = i
	for line[i] != ',' {
		i++
	}
	low := parseFloat(line[oi:i])

	//close -not done
	i = len(line) - 1
	oi = i + 1
	for line[i] != ',' {
		i--
	}
	volume := parseFloat(line[i+1 : oi])
	return &model.DataPoint{Date: datetime, Price: price, Low: low, Volume: volume}
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

func GetTickerFilePath(name string) string {
	return "C:\\Users\\Guillaume\\Desktop\\stocks\\intraday\\" + name + ".txt"
}
func OpenTicker(name string) *os.File {
	filePath := "C:\\Users\\Guillaume\\Desktop\\stocks\\intraday\\" + name + ".txt"
	file, _ := os.Open(filePath)
	return file
}

func parseDate(date []byte) model.SimplifiedDate {
	year := (((uint16(date[6])-'0')*10+uint16(date[7])-'0')*10+uint16(date[8])-'0')*10 + uint16(date[9]) - '0'
	month := (uint8(date[0])-'0')*10 + uint8(date[1]) - '0'
	day := (uint8(date[3])-'0')*10 + uint8(date[4]) - '0'
	hour := (uint8(date[11])-'0')*10 + uint8(date[12]) - '0'
	minute := (uint8(date[14])-'0')*10 + uint8(date[15]) - '0'
	return model.SimplifiedDate{Year: year, Month: month, Day: day, Hour: hour, Minute: minute}
}

func parseFloat(str []byte) float64 {
	res := 0
	div := 0
	for _, c := range str {
		if c == '.' {
			div = 1
			continue
		}
		res = res*10 + (int(c) - '0')
		div *= 10
	}
	div = div
	if div > 0 {
		return float64(res) / float64(div)
	}
	return float64(res)
}

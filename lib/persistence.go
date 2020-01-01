package lib

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/willfr/goback/globals"
)

func SaveRun() {
	f, _ := os.OpenFile("C:\\Users\\Guillaume\\Desktop\\stocks\\run.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	defer f.Close()
	f.Truncate(0)
	f.Seek(0,0)
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
	f2.Seek(0,0)
	for _, h := range globals.History {
		f2.WriteString(h.Name + " " + h.Date.Format(time.RFC3339) + " " + fmt.Sprintf("%.0f", h.Quantity) + " " + fmt.Sprintf("%f", h.Price) + " \r\n")
	}
	f2.Sync()
}

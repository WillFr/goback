package lib

import (
	"bufio"
	"fmt"
	"github.com/niubaoshu/gotiny"
	"github.com/willfr/goback/globals"
	"github.com/willfr/goback/model"
	"math"
	"os"
	"sort"
	"strings"
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
		f2.WriteString(h.Name + " " + h.Date.Format() + " " + fmt.Sprintf("%.0f", h.Quantity) + " " + fmt.Sprintf("%f", h.Price) + " \r\n")
	}
	f2.Sync()
}

type EncoderDecoder struct {
	encoder    *gotiny.Encoder
	decoder    *gotiny.Decoder
	writer     *bufio.Writer
	reader     *bufio.Reader
	readOffset int
	lastBuff   []byte
}

func NewEncoderDecoder(writer *bufio.Writer, reader *bufio.Reader) *EncoderDecoder {
	return &EncoderDecoder{
		writer: writer,
		reader: reader,
	}
}

func (ed *EncoderDecoder) WriteToBin(dp *model.DataPoint) {
	// YYYY YYMM - MMDD DDDh - hhhh mmmm - mm
	yearMoByte := uint8(dp.Date.Year-1995)<<2 | dp.Date.Month>>2
	nthDayHByte := dp.Date.Month<<6 | dp.Date.Day<<1 | dp.Date.Hour>>4
	ourMinByte := dp.Date.Hour<<4 | dp.Date.Minute>>2

	intPrice := int32(math.Floor(dp.Price * 1000))
	intLow := int32(math.Floor(dp.Low * 1000))
	intVolume := int32(math.Floor(dp.Volume))

	//     1           2           3           4          1          2          3           4
	// PPPP PPPP - PPPP PPPP - PPPP PPLL - LLLL LLLL- LLLL LLLL- LLLL VVVV- VVVV VVVV - VVVV VVMM
	priceLoByte := intPrice<<10 | intLow>>12
	wVolumebyte := intLow<<20 | intVolume<<2 | (int32(dp.Date.Minute) & 0x3)

	var err error
	(*ed).writer.WriteByte(byte(yearMoByte))
	(*ed).writer.WriteByte(byte(nthDayHByte))
	(*ed).writer.WriteByte(byte(ourMinByte))

	(*ed).writer.WriteByte(byte(priceLoByte >> 24))
	(*ed).writer.WriteByte(byte(priceLoByte >> 16))
	(*ed).writer.WriteByte(byte(priceLoByte >> 8))
	(*ed).writer.WriteByte(byte(priceLoByte >> 0))

	(*ed).writer.WriteByte(byte(wVolumebyte >> 24))
	(*ed).writer.WriteByte(byte(wVolumebyte >> 16))
	(*ed).writer.WriteByte(byte(wVolumebyte >> 8))
	(*ed).writer.WriteByte(byte(wVolumebyte >> 0))

	//err:= binary.Write(writer, binary.LittleEndian, dp)
	if err != nil {
		fmt.Println(err)
	}
}

func (ed *EncoderDecoder) ReadFromBin() (*model.DataPoint, error) {
	var returnValue model.DataPoint
	var yearMoByte, nthDayHByte, ourMinByte byte
	var priceLoByte1, priceLoByte2, priceLoByte3, priceLoByte4 byte
	var wVolumebyte1, wVolumebyte2, wVolumebyte3, wVolumebyte4 byte
	var err error
	yearMoByte, err = ed.reader.ReadByte()
	nthDayHByte, err = ed.reader.ReadByte()
	ourMinByte, err = ed.reader.ReadByte()

	priceLoByte1, err = ed.reader.ReadByte()
	priceLoByte2, err = ed.reader.ReadByte()
	priceLoByte3, err = ed.reader.ReadByte()
	priceLoByte4, err = ed.reader.ReadByte()

	wVolumebyte1, err = ed.reader.ReadByte()
	wVolumebyte2, err = ed.reader.ReadByte()
	wVolumebyte3, err = ed.reader.ReadByte()
	wVolumebyte4, err = ed.reader.ReadByte()

	if err != nil {
		return nil, err
	}

	// YYYY YYMM - MMDD DDDh - hhhh mmmm - mm
	minute := uint8((ourMinByte&0x0F)<<2 | (wVolumebyte4 & 0x3))
	hour := uint8(ourMinByte>>4 | (nthDayHByte&0x1)<<4)
	day := uint8((nthDayHByte & 0x3E) >> 1)
	month := uint8(((nthDayHByte & 0xC0) >> 6) | (yearMoByte&3)<<2)
	year := uint16(yearMoByte>>2) + 1995

	//     1           2           3           4          1          2          3           4
	// PPPP PPPP - PPPP PPPP - PPPP PPLL - LLLL LLLL- LLLL LLLL- LLLL VVVV- VVVV VVVV - VVVV VVMM
	priceInt := int32(priceLoByte1)<<14 | int32(priceLoByte2)<<6 | int32(priceLoByte3&0xFC)>>2
	lowInt := int32(priceLoByte3&0x03)<<20 | int32(priceLoByte4)<<12 | int32(wVolumebyte1)<<4 | int32(wVolumebyte2)>>4
	volumeInt := int32(wVolumebyte2&0xF)<<14 | int32(wVolumebyte3)<<6 | int32(wVolumebyte4>>2)

	price := float64(priceInt) / 1000.0
	low := float64(lowInt) / 1000.0
	volume := float64(volumeInt)

	date := model.SimplifiedDate{
		Year:   year,
		Month:  month,
		Day:    day,
		Hour:   hour,
		Minute: minute,
	}
	returnValue = model.DataPoint{
		Price:  price,
		Low:    low,
		Volume: volume,
		Date:   date,
	}

	return &returnValue, err
}

package srt

import (
	"io"
	"time"
	"fmt"
	"math"
)

func writeTime(w io.Writer, dur time.Duration) (nbytes int, err error) {
	hoursToPrint := int(math.Floor(dur.Hours()))
	minutesToPrint := int(math.Floor(dur.Minutes() - (time.Duration(hoursToPrint)*time.Hour).Minutes()))
	secondsToPrint := int(math.Floor(dur.Seconds() - (time.Duration(hoursToPrint)*time.Hour + time.Duration(minutesToPrint)*time.Minute).Seconds()))
	millisecondsToPrint := int(math.Floor(float64(dur/time.Millisecond - (time.Duration(hoursToPrint)*time.Hour + time.Duration(minutesToPrint)*time.Minute + time.Duration(secondsToPrint)*time.Second)/time.Millisecond)))

	nbytes, err = fmt.Fprintf(w, "%02d:%02d:%02d,%03d", hoursToPrint, minutesToPrint, secondsToPrint, millisecondsToPrint)
	return
}

func (s Subtitle)WriteTo(writer io.Writer) (nbytes int, err error) {
	var wlen int

	wlen, err = fmt.Fprintf(writer, "%v\n", s.Number)
	nbytes += wlen
	if err != nil {
		return nbytes, err
	}

	wlen, err = writeTime(writer, s.Start)
	nbytes += wlen
	if err != nil {
		return nbytes, err
	}

	wlen, err = io.WriteString(writer, " --> ")
	nbytes += wlen
	if err != nil {
		return nbytes, err
	}

	wlen, err = writeTime(writer, s.End)
	nbytes += wlen
	if err != nil {
		return nbytes, err
	}

	wlen, err = io.WriteString(writer, "\n")
	nbytes += wlen
	if err != nil {
		return nbytes, err
	}

	wlen, err = io.WriteString(writer, s.Text)
	nbytes += wlen
	if err != nil {
		return nbytes, err
	}

	wlen, err = io.WriteString(writer, "\n\n")
	nbytes += wlen
	if err != nil {
		return nbytes, err
	}

	return nbytes, nil
}


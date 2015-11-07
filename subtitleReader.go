package srt

import (
	"fmt"
	"bufio"
	"bytes"
	"io"
	"strconv"
	"strings"
	"time"
	"regexp"
)

func dropCR(data []byte) []byte {
	if len(data) > 0 && data[len(data)-1] == '\r' {
		return data[0 : len(data)-1]
	}
	return data
}

func scanDoubleNewline(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if i := bytes.Index(data, []byte{'\n', '\n'}); i >= 0 {
		// We have a full double newline-terminated line.
		return i + 2, dropCR(data[0:i]), nil
	} else if i := bytes.Index(data, []byte{'\n', '\r', '\n'}); i >= 0 {
		// We have a full double newline-terminated line.
		return i + 3, dropCR(data[0:i]), nil
	}
	// If we're at EOF, we have a final, non-terminated line. Return it.
	if atEOF {
		return len(data), dropCR(data), nil
	}
	// Request more data.
	return 0, nil, nil
}

type SubtitleScanner struct {
	scanner *bufio.Scanner
	nextSub Subtitle
	err     error
}

func NewScanner(r io.Reader) SubtitleScanner {
	s := bufio.NewScanner(r)
	s.Split(scanDoubleNewline)
	return SubtitleScanner{s, Subtitle{}, nil}
}

func parseTime(input string) (time.Time, error) {
	regex := regexp.MustCompile(`(\d{2}):(\d{2}):(\d{2}),(\d{3})`)
	matches := regex.FindStringSubmatch(input)

	hour, err := strconv.Atoi(matches[1])
	if (err != nil) { return time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC), err }
	minute, err := strconv.Atoi(matches[2])
	if (err != nil) { return time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC), err }
	second, err := strconv.Atoi(matches[3])
	if (err != nil) { return time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC), err }
	millisecond, err := strconv.Atoi(matches[4])
	if (err != nil) { return time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC), err }

	return time.Date(1,1,1,hour,minute,second,millisecond * 1000000, time.Local), nil
}

func (s *SubtitleScanner) Scan() bool {
	if s.scanner.Scan() {
		var (
			nextnum int
			start time.Time
			duration time.Duration
			subtitletext string
		)

		str := strings.Split(s.scanner.Text(), "\n")

		for i := 0; i < len(str); i++ {
			text := strings.TrimRight(str[i], "\r")
			switch i {
			case 0:
				num, err := strconv.Atoi(text)
				if err != nil {
					s.err = err
					return false
				}
				nextnum = num
			case 1:
				elements := strings.Split(text, " ")
				if len(elements) >= 3 {
					startTime, err := parseTime(elements[0])
					if err != nil {
						s.err = err
						return false
					}
					endTime, err := parseTime(elements[2])
					if err != nil {
						s.err = err
						return false
					}
					start = startTime
					duration = endTime.Sub(startTime)
				} else {
					s.err = fmt.Errorf("srt: Invalid timestamp on row: %s", text)
					return false
				}
			default:
				if len(subtitletext) > 0 {
					subtitletext += "\n"
				}
				subtitletext += text
			}
		}

		s.nextSub = Subtitle{nextnum, start, duration, subtitletext}

		return true;
	} else {
		return false
	}
}

func (s *SubtitleScanner) Err() error {
	if s.err != nil {
		return s.err;
	}
	return s.scanner.Err()
}

func (s *SubtitleScanner) Subtitle() Subtitle {
	return s.nextSub
}

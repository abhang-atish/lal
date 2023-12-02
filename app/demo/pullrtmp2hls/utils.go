package main

import (
	"bytes"
	"github.com/grafov/m3u8"
	"time"
)

func GetNonNilSegments(arr []*m3u8.MediaSegment) []*m3u8.MediaSegment {
	var nonNil []*m3u8.MediaSegment
	for _, i := range arr {
		if i != nil {
			nonNil = append(nonNil, i)
		}
	}
	return nonNil
}

func GetLastNonNilSegment(arr []*m3u8.MediaSegment) *m3u8.MediaSegment {
	var lastNonNil *m3u8.MediaSegment
	for _, i := range arr {
		if i != nil {
			lastNonNil = i
		}
	}
	return lastNonNil
}

func ConvertTimeToRFC3339(v time.Time) (time.Time, error) {
	value := v.Format(time.RFC3339)
	return time.Parse(time.RFC3339, value)
}

func ConvertTimeToRFC3339Nano(v time.Time) (time.Time, error) {
	value := v.Format(time.RFC3339Nano)
	return time.Parse(time.RFC3339Nano, value)
}

func ParseM3U8(data []byte, key string) (m3u8.Playlist, m3u8.ListType, bool) {
	p, l, err := m3u8.DecodeFrom(bytes.NewReader(data), true)
	if err != nil {
		return nil, 0, false
	}
	return p, l, true
}

package rvgo

// Some simple Golang helpers

import (
	"os"
	"strings"
	"time"
)

func CloneBytes(a []byte) []byte {
	if a == nil {
		return nil
	}

	b := make([]byte, len(a))
	copy(b, a)
	return b
}

const millisInSecond = time.Second / time.Millisecond
const nanosInMillisecond = time.Millisecond / time.Nanosecond

func TimeToUnixMillis(ts time.Time) int64 {
	return ts.Unix()*int64(millisInSecond) + int64(ts.Nanosecond()/int(nanosInMillisecond))
}

func UnixSecToMillis(sec int64) int64 {
	return sec * int64(millisInSecond)
}

func UnixMillisToTime(millis int64) time.Time {
	return time.Unix(millis/int64(millisInSecond), (millis%int64(millisInSecond))*int64(nanosInMillisecond))
}

func UnixSecToTime(sec int64) time.Time {
	return time.Unix(sec, 0)
}

func TrimmedStringFromFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(data)), nil
}

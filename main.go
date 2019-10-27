package goutils

import (
    "crypto/md5"
    "encoding/base64"
    "fmt"
    "time"
)

func EncodeBase64(msg string) string {
    return base64.StdEncoding.EncodeToString([]byte(msg))
}

func MD5(src string) string {
    return fmt.Sprintf("%x", md5.Sum([]byte(src)))
}

func NowNano() uint64 {
    return uint64(time.Now().UnixNano())
}

func Now() int64 {
    return time.Now().Unix()
}

func random(foo uint64) float64 {
    return float64(( foo * 9301 + 49297 ) % 233280) / 233280
}

func Random() float64 {
    now := NowNano()
    now = uint64(random(now) * (1 << 63)) + NowNano()
    return float64(( now * 9301 + 49297 ) % 233280) / 233280
}

func RandomString(length int) string {
    foo := fmt.Sprintf("%x", int64(Random() * 10000000))
    return EncodeBase64(foo)
}

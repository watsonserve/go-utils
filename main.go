package goutils

import (
    "crypto/md5"
    "encoding/base64"
    "fmt"
    "github.com/watsonserve/goengine"
    "net/http"
    "net/url"
    "time"
)

func cutUri(raw *url.URL) string {
    uri := raw.Path
    if "" != raw.RawQuery {
        uri += "?" + raw.RawQuery
    }
    if "" != raw.Fragment {
        uri += "#" + raw.Fragment
    }
    return uri
}

func EncodeBase64(msg string) string {
    return base64.StdEncoding.EncodeToString([]byte(msg))
}

func MD5(src string) string {
    return fmt.Sprintf("%x", md5.Sum([]byte(src)))
}

func Now() int64 {
    return time.Now().Unix()
}

/**
 * @param {string} authAddr 授权路径，例如：/auth
 * @param {*url.URL} raw 当前路径，将被转换为：%2Fpathname%3Fsearch%23hash
 * @return string /auth?r=%2Fpathname%3Fsearch%23hash
 */
func getAuthAddr(authAddr string, raw *url.URL) string {
    curAddr := cutUri(raw)
    redirect := url.URL {
        Scheme: "http",
        Host: "shopping.watsonserve.com",
        Path: authAddr,
    }
    q := redirect.Query()
    q.Set("rd", curAddr)
    redirect.RawQuery = q.Encode()
    return redirect.String()
}

// 从session检出用户id
func WhoIsUser(session *goengine.Session) map[string]interface{} {
    // 检出数据
    user := session.Get(SESSION_USER_KEY)
    if nil == user {
        return nil
    }
    return user.(map[string]interface{})
}

// 检查ref
func chkReferer(req *http.Request, selfDomain string) *url.URL {
    for {
        referer := req.Header.Get("referer")
        if "" == referer {
            break
        }
        refUri, err := url.Parse(referer)
        if nil != err {
            break
        }
        
        refHost := refUri.Scheme + "://" + refUri.Host
        if selfDomain != refHost {
            break
        }
        return refUri
    }
    return nil
}

func nowNano() uint64 {
    return uint64(time.Now().UnixNano())
}
func random(foo uint64) float64 {
    return float64(( foo * 9301 + 49297 ) % 233280) / 233280
}

func Random() float64 {
    now := nowNano()
    now = uint64(random(now) * (1 << 63)) + nowNano()
    return float64(( now * 9301 + 49297 ) % 233280) / 233280
}

func RandomString(length int) string {
    foo := fmt.Sprintf("%x", int64(Random() * 10000000))
    return EncodeBase64(foo)
}

func passportRPC(ticket string, app string, secret string) (*UserData, error) {
    // @TODO
    return nil, nil
}


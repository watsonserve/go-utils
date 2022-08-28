package goutils

import (
    "bufio"
    "crypto/tls"
    "fmt"
    "io"
    "net"
	"log"
	"net/http"

	http3 "github.com/lucas-clemente/quic-go/http3"
)

/*
 * 这里使用的是每个链接启动一个新的go程的模型
 * 高并发的话，性能取决于go语言的协程能力
 */
func TLSSocket(port string, crt string, key string) (net.Listener, error) {
    cert, err := tls.LoadX509KeyPair(crt, key)
    if nil != err {
        return nil, err
    }
    ln, err := tls.Listen("tcp", port, &tls.Config {
        Certificates: []tls.Certificate{cert},
        CipherSuites: []uint16 {
          tls.TLS_RSA_WITH_AES_256_CBC_SHA,
          tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
          tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
          tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
          tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
        },
        PreferServerCipherSuites: true,
    })
    if nil != err {
        return nil, err
    }
    defer ln.Close()
    return ln, nil
}

/*
 * 这里使用的是每个链接启动一个新的go程的模型
 * 高并发的话，性能取决于go语言的协程能力
 */
func Socket(port string) (net.Listener, error) {
    ln, err := net.Listen("tcp", port)
    if nil != err {
        return nil, err
    }
    defer ln.Close()
    return ln, nil
}

type Stream struct {
    scanner *bufio.Scanner
    sock io.ReadWriteCloser
}

func InitStream(sock io.ReadWriteCloser) *Stream {
    return &Stream {
        scanner: bufio.NewScanner(sock),
        sock: sock,
    }
}

func (this *Stream) ReadLine() (string, error) {
    this.scanner.Scan()
    err := this.scanner.Err()
    if nil != err {
        return "", err
    }
    msg := this.scanner.Text()
    fmt.Printf("c: %s\n", msg)
    return msg, nil
}

// 发送
func (this *Stream) Send(content string) {
    fmt.Printf("s: %s\n", content)
    fmt.Fprint(this.sock, content)
}

// 发送并关闭
func (this *Stream) End(content string) {
    fmt.Fprint(this.sock, content)
    this.sock.Close()
}

func ListenAndHttp(network, addr, crt, key string, handler http.Handler) {
	var err error

	switch network {
	case "quic":
		err = http3.ListenAndServeQUIC(addr, crt, key, handler)
	case "tcp":
		fallthrough
	case "unix":
		ln, _err := net.Listen(network, addr)

		if nil == _err {
			if "" != crt && "" != key {
				_err = http.ServeTLS(ln, handler, crt, key)
			} else {
				_err = http.Serve(ln, handler)
			}
		}
		err = _err
	}
	if nil != err {
		log.Fatal("failed to start server", err)
	}
}

type AddressSet struct {
	Tcp string
	TcpLts string
	Unix string
	Quic string
}

func ServeHttp(addr *AddressSet, handler http.Handler, crt, key string) {
    if nil == addr {
        return
    }
	addrTcp := addr.Tcp
	addrTcpLts := addr.TcpLts
	addrUnix := addr.Unix
	addrQuic := addr.Quic

	if "" != addrTcp {
		fmt.Printf("listen tcp: %s\n", addrTcp)
		go ListenAndHttp("tcp", addrTcp, "", "", handler)
	}

	if "" != addrUnix {
		fmt.Printf("listen unix: %s\n", addrUnix)
		go ListenAndHttp("unix", addrUnix, "", "", handler)
	}

	if "" != addrTcpLts {
		fmt.Printf("listen tls: %s\n", addrTcpLts)
		go ListenAndHttp("tcp", addrTcpLts, crt, key, handler)
	}

	if "" != addrQuic {
		fmt.Printf("listen quic: %s\n", addrQuic)
		go ListenAndHttp("quic", addrQuic, crt, key, handler)
	}
}
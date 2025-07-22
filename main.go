package main

import (
	"bufio"
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/url"
	"strings"
	"sync"
)

func main() {

	var address string
	flag.StringVar(&address, "L", "127.0.0.1:11088", "listen address")
	flag.Parse()

	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal("Listen error:", err)
	}
	log.Println("Listening on " + address)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Accept error:", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {

	// 读取第一个字节
	buf := make([]byte, 1)
	n, err := conn.Read(buf)
	if n != 1 || err != nil {
		log.Println("Read error:", err)
		return
	}
	if buf[0] == 0x05 { // SOCKS5
		handleSocks5(conn)
	} else { // HTTP
		handleHttpProxy(conn, buf[0])
	}
}

func handleSocks5(conn net.Conn) {
	buf := make([]byte, 256)

	// 读取 VER 和 NMETHODS
	n, err := io.ReadFull(conn, buf[:1])
	if n != 1 || err != nil {
		log.Println("reading version err: ", err)
		return
	}
	nMethods := int(buf[1])

	// 读取 METHODS 列表
	n, err = io.ReadFull(conn, buf[:nMethods])
	if n != nMethods || err != nil {
		log.Println("reading methods err: ", err)
		return
	}

	// 响应认证方法，使用无认证
	n, err = conn.Write([]byte{0x05, 0x00})
	if n != 2 || err != nil {
		log.Println("write methods err:", err)
		return
	}

	n, err = io.ReadFull(conn, buf[:4])
	if n != 4 || err != nil {
		log.Println("read connect err:", err)
		return
	}
	ver, cmd, _, atyp := buf[0], buf[1], buf[2], buf[3]
	if ver != 0x05 || cmd != 0x01 {
		log.Println("invalid ver/cmd")
		return
	}

	var addr string
	switch atyp {
	case 0x01: // IPv4
		n, err = io.ReadFull(conn, buf[:4])
		if n != 4 || err != nil {
			log.Println("read ipv4 addr err:", err)
			return
		}
		addr = fmt.Sprintf("%d.%d.%d.%d", buf[0], buf[1], buf[2], buf[3])
	case 0x03: // 域名
		n, err = io.ReadFull(conn, buf[:1])
		if n != 1 || err != nil {
			log.Println("read domain addr len err:", err)
			return
		}
		addrLen := int(buf[0])

		n, err = io.ReadFull(conn, buf[:addrLen])
		if n != addrLen || err != nil {
			log.Println("read domain addr err:", err)
			return
		}
		addr = string(buf[:addrLen])
	case 0x04: // IPv6
		// 处理IPv6地址，这里暂时忽略
		// ...
	default:
		// 错误处理
		return
	}

	n, err = io.ReadFull(conn, buf[:2])
	if n != 2 || err != nil {
		log.Println("read addr port err:", err)
		return
	}
	port := binary.BigEndian.Uint16(buf[:2])

	destAddrPort := fmt.Sprintf("%s:%d", addr, port)

	log.Println("Via socks5 connect to:", destAddrPort)
	// 建立连接
	targetConn, err := net.Dial("tcp", destAddrPort)
	if err != nil {
		log.Println("Dial error:", err)
		return
	}

	n, err = conn.Write([]byte{0x05, 0x00, 0x00, 0x01, 0, 0, 0, 0, 0, 0})
	if err != nil {
		targetConn.Close()
		log.Println("write connect resp error:", err)
		return
	}

	// 使用 sync.WaitGroup 管理双向传输
	var wg sync.WaitGroup
	wg.Add(2)
	go io.Copy(conn, targetConn)
	go io.Copy(targetConn, conn)
	// 等待所有传输完成
	defer conn.Close()
	defer targetConn.Close()
	wg.Wait()
}

func handleHttpProxy(conn net.Conn, firstByte byte) {

	connReader := bufio.NewReader(conn)
	line, _, err := connReader.ReadLine()
	if err != nil {
		log.Println("read http proto err:", err)
		return
	}
	n := len(line)
	buf := make([]byte, n+1)
	buf[0] = firstByte
	copy(buf[1:], line)

	httpReqLine := string(buf)

	// 解析请求行 GET http://example.com/path HTTP/1.1
	httpReqLineParts := strings.Split(strings.TrimSpace(httpReqLine), " ")
	if len(httpReqLineParts) < 3 {
		log.Println("Invalid request")
		return
	}

	destAddrPort := ""
	isConnect := false

	if strings.Contains(strings.ToUpper(httpReqLineParts[0]), "CONNECT") {
		isConnect = true
		destAddrPort = httpReqLineParts[1]
	} else {
		urlParsed, errP := url.Parse(httpReqLineParts[1])
		if errP != nil {
			log.Println("Url parse error:", errP)
			return
		}
		destAddrPort = urlParsed.Host
		if !strings.Contains(destAddrPort, ":") {
			destAddrPort += ":80"
		}
	}

	log.Println("Connect to:", destAddrPort)
	// 建立连接
	targetConn, err := net.Dial("tcp", destAddrPort)
	if err != nil {
		log.Println("Dial error:", err)
		return
	}
	var wg sync.WaitGroup
	wg.Add(2)
	if isConnect {
		_, _ = fmt.Fprintf(conn, "HTTP/1.1 200 Connection established\r\n\r\n")
	} else {
		targetConn.Write([]byte(httpReqLine + "\r\n"))
		// targetConn.Write(buf)
		go io.Copy(targetConn, connReader)
	}
	go io.Copy(conn, targetConn)
	if isConnect {
		go io.Copy(targetConn, conn)
	}
	defer conn.Close()
	defer targetConn.Close()
	wg.Wait()

}

package main

import (
	"fmt"
        "net"
	"log"
        "strings"
        "strconv"

	"github.com/mvanveen/framebuffer"
)

type message_from struct {
    message string
    ip string
    conn *net.UDPConn
}

func handleMessage(fb *framebuffer.Framebuffer, msg string, ip string, conn *net.UDPConn) {
   parts := strings.SplitN(msg, " ", 6)
   cmd := parts[0]


   if cmd == "get" {
      host := net.UDPAddr{
          Port: 6668,
          IP: net.ParseIP(ip),
      }

      x_coord, _ := strconv.Atoi(parts[1])
      y_coord, _ := strconv.Atoi(parts[2])

      color := fb.GetPixel(x_coord, y_coord)
      response := fmt.Sprintf("%d %d %d %d %d", x_coord, y_coord, color[0], color[1], color[2])
      conn.WriteToUDP([]byte(response), &host)

   } else {

      x_coord, _ := strconv.Atoi(parts[1])
      y_coord, _ := strconv.Atoi(parts[2])

      red, _ := strconv.Atoi(parts[3])
      green, _ := strconv.Atoi(parts[4])
      blue, _ := strconv.Atoi(parts[5])

      SendColor(fb, x_coord, y_coord, red, green, blue)
   }
}


func myUDPServer(messages chan message_from) {
    fmt.Println("opening UDP server")

    LISTENING_IP := "0.0.0.0"
    LISTENING_PORT := 6668

    addr := net.UDPAddr{
        Port: LISTENING_PORT,
        IP: net.ParseIP(LISTENING_IP),
    }

    conn, err := net.ListenUDP("udp", &addr)
    defer conn.Close()
    if err != nil {
        panic(err)
    }

    var buf []byte = make([]byte, 3000)

    for {
        rlen, address, err := conn.ReadFromUDP(buf)

        if err != nil {
            fmt.Println("error reading data from connection")
            fmt.Println(err)
            return
        }

        if address != nil {
            if rlen > 0 {
                messages <- message_from{string(buf[0:rlen]), address.IP.String(), conn}
            }
        }
     }
}

func SendColor(fb *framebuffer.Framebuffer, x_coord int, y_coord int, red int, green int, blue int) {
    fb.WritePixel(x_coord, y_coord, red, green, blue, 255)
}

func handleScreen(fb *framebuffer.Framebuffer, messages chan message_from) {
    fb.Clear(0, 0, 0, 0)
    for {
        msg := <- messages
        handleMessage(fb, msg.message, msg.ip, msg.conn)
    }
}


func main() {
	fb, err := framebuffer.Init("/dev/fb0")
	if err != nil {
                log.Fatalln(err)
	}
	defer fb.Close()
        messages := make(chan message_from, 10000)

        go handleScreen(fb, messages)
        myUDPServer(messages)
}

package main

import (
	"fmt"
        "net"
	"log"
        "strings"
        "strconv"

	"github.com/kaey/framebuffer"
)


func handleMessage(fb *framebuffer.Framebuffer, msg string) {
   parts := strings.SplitN(msg, " ", 5)

   x_coord, _ := strconv.Atoi(parts[0])
   y_coord, _ := strconv.Atoi(parts[1])

   red , _ := strconv.Atoi(parts[2])
   blue , _ := strconv.Atoi(parts[3])
   green , _ := strconv.Atoi(parts[4])

   SendColor(fb, x_coord, y_coord, red, blue, green)

   // send message to channel to incrememnt color
   fmt.Println("got message: (", x_coord, ", ", y_coord, ") ", red, blue, green)
}


func myUDPServer(messages chan string) {
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

    var buf []byte = make([]byte, 1500)

    for {
        rlen, address, err := conn.ReadFromUDP(buf)

        if err != nil {
            fmt.Println("error reading data from connection")
            fmt.Println(err)
            return
        }

        if address != nil {
            if rlen > 0 {
                messages <- string(buf[0:rlen])
            }
        }
     }
}

func SendColor(fb *framebuffer.Framebuffer, x_coord int, y_coord int, red int, green int, blue int) {
    fb.WritePixel(x_coord, y_coord, red, green, blue, 255)
}

func handleScreen(fb *framebuffer.Framebuffer, messages chan string) {
    fb.Clear(0, 0, 0, 0)
    for {
        msg := <- messages
        handleMessage(fb, msg)
    }
}

func main() {
	fb, err := framebuffer.Init("/dev/fb0")
	if err != nil {
		log.Fatalln(err)
	}
	defer fb.Close()

        messages := make(chan string, 10000)

        go myUDPServer(messages)
        handleScreen(fb, messages)
}

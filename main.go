package main

import "fmt"
import "net"
import "strings"

func handleMessage(msg string) {
   parts := strings.SplitN(msg, " ", 3)
   x_coord := parts[0]
   y_coord := parts[1]
   color := strings.TrimSpace(parts[2])

   // send message to channel to incrememnt color
   fmt.Println("got message: (", x_coord, ", ", y_coord, ") ", color)
}

func myUDPServer() {
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
                go handleMessage(string(buf[0:rlen]))
            }
        }
     }
}

func main() {
    myUDPServer()
}

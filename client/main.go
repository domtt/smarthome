package main

import (
        "bufio"
        "fmt"
        "net"
        "os"
)

func main() {
        arguments := os.Args
        if len(arguments) == 1 {
                fmt.Println("Please provide host:port.")
                return
        }

        CONNECT := arguments[1]
        c, err := net.Dial("tcp", CONNECT)
        if err != nil {
                fmt.Println(err)
                return
        }

        go func() {for {
                reader := bufio.NewReader(os.Stdin)
                fmt.Print(">> ")
                text, _ := reader.ReadString('\n')
                fmt.Fprintf(c, text+"\n")
        }}()
        for {
                message, _ := bufio.NewReader(c).ReadString('\n')
                fmt.Print("->: " + message)
              }
}
    

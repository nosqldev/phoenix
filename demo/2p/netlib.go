package main

import (
    "fmt"
    "net"
    "strconv"
    "strings"
    //"encoding/binary"
    //"bytes"
)

type Courser struct {
    RegListener net.Listener
}

func Echo(s string) {
    fmt.Println(s)
}

func RunListener(port int) Courser {
    var c Courser
    l, err := net.Listen("tcp4", ":" + strconv.Itoa(port))
    if err != nil {
        fmt.Println(err)
    }
    c.RegListener = l

    go register_thread(c.RegListener)

    return c
}

func RunHeartbeat(port string) {
    go func(port string) {
        udpaddr, _ := net.ResolveUDPAddr("udp4", ":" + port)
        for {
            var request []byte = make([]byte, 1024)
            c, _ := net.ListenUDP("udp4", udpaddr)
            Echo("Heartbeat listenning: " + fmt.Sprint(udpaddr))

            n, _, err := c.ReadFrom(request)
            if err != nil {
                Echo("ReadFrom() err:" + fmt.Sprint(err))
                c.Close()
                continue
            }

            fmt.Print("heartbeat recv " + fmt.Sprint(n) + " bytes")
            n, _, err = c.ReadFrom(request)
            fmt.Print("heartbeat recv " + fmt.Sprint(n) + " bytes")

            /*
             *var remote_ip net.IP
             *buf := bytes.NewBuffer(request)
             *binary.Read(buf, binary.LittleEndian, &remote_ip)
             *fmt.Print(remote_ip)
             */

            c.Close()
        }
    } (port)
}

func register_thread(l net.Listener) {
    Echo("register thread launched")
    for {
        conn, _ := l.Accept()
        var request string
        var response string
        fmt.Fscanf(conn, "%s", &request)

        s_array := strings.Split(request, ":")
        SV.Lock()
        notify_other_masters()
        for addr, _ := range SV.Servers {
            response += addr + "\n"
        }
        SV.Servers[request] = s_array
        SV.Unlock()

        fmt.Fprint(conn, response)
        Echo(SV.Name + " got register request: host->" + s_array[0] + ", port->" + s_array[1])
        conn.Close()
    }
}

func notify_other_masters() {
    for _, addr_info := range SV.Servers {
        fmt.Printf("Notify(Host: %s, Port: %s)\n", addr_info[0], addr_info[1])
    }
}

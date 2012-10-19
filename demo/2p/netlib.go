package main

import (
    "time"
    "fmt"
    "net"
    "strconv"
    "encoding/binary"
    "bytes"
    "hash/fnv"
    "io"
)

type Courser struct {
    RegListener net.Listener
}

func Echo(s string) {
    fmt.Println(s)
}

func RunRegister(port int) Courser {
    var c Courser
    l, err := net.Listen("tcp4", ":" + strconv.Itoa(port))
    if err != nil {
        fmt.Println(err)
    }
    c.RegListener = l

    go register_thread(c.RegListener)

    return c
}

/* HEARTBEAT STRUCTURE
 *  [Normal Protocol]
 *  Root Master IP
 *  Root Master Port
 *  ServerList Checksum
 *  -------------------
 *  [Fetch Server List Protocol]
 *  ServerList(string)
 */

func parse_binary_addr(buffer []byte) (string, string) {
    var ip net.IP = net.IPv4(buffer[0], buffer[1], buffer[2], buffer[3])
    var port int16
    buf := bytes.NewBuffer(buffer[4:6])
    binary.Read(buf, binary.LittleEndian, &port)

    return ip.String(), fmt.Sprint(port)
}

func calc_serverlist_cksum() uint64 {
    h := fnv.New64()

    SV.Lock()
    for i := range SV.ServerList {
        io.WriteString(h, SV.ServerList[i])
    }
    SV.Unlock()

    //fmt.Println("Server cksum:", h.Sum64())

    return h.Sum64()
}

func parse_serverlist_cksum(buffer []byte) (cksum uint64) {
    buf := bytes.NewBuffer(buffer[14:22])
    binary.Read(buf, binary.LittleEndian, &cksum)

    return cksum
}

func parse_heartbeat_id(buffer []byte) (id uint64) {
    buf := bytes.NewBuffer(buffer[6:14])
    binary.Read(buf, binary.LittleEndian, &id)

    return id
}

func build_heartbeat_resp(loadavg float32) []byte {
    buf := new(bytes.Buffer)
    binary.Write(buf, binary.LittleEndian, uint32(1)) // protocol version
    binary.Write(buf, binary.LittleEndian, uint32(0)) // status
    binary.Write(buf, binary.LittleEndian, SV.LocalIP_num)
    binary.Write(buf, binary.LittleEndian, SV.LocalPort_num)
    binary.Write(buf, binary.LittleEndian, loadavg)

    return buf.Bytes()
}

func send_heartbeat(ip, port string, hb_id uint64) {
    p, _ := strconv.Atoi(port)
    laddr := net.UDPAddr { net.ParseIP("0.0.0.0").To4(), 10001 }
    raddr := net.UDPAddr { net.ParseIP(ip).To4(), p }
    conn, err := net.DialUDP("udp4", &laddr, &raddr)
    defer conn.Close()
    if err != nil {
        fmt.Println("[error 3]", err, " | ", hb_id, laddr, raddr)
        return
    }
    srv_cksum := calc_serverlist_cksum()
    buf := new(bytes.Buffer)
    binary.Write(buf, binary.LittleEndian, SV.LocalIP_num)
    binary.Write(buf, binary.LittleEndian, uint16(10000))
    binary.Write(buf, binary.LittleEndian, hb_id)
    binary.Write(buf, binary.LittleEndian, srv_cksum)
    n, err := conn.Write(buf.Bytes())
    if n == 0 || err != nil {
        fmt.Println("[error 1]", n, err)
        return
    }

    /*
     *resp := make([]byte, 1024)
     *conn.SetReadDeadline(time.Now().Add(100 * time.Millisecond))
     *n, err = conn.Read(resp)
     *if n == 0 || err != nil {
     *    fmt.Println("[error 2]", n, err)
     *}
     */
}

func heartbeat_lisener() {
    laddr, _ := net.ResolveUDPAddr("udp4", "0.0.0.0:10000")

    for {
        var request []byte = make([]byte, 1024)
        c, _ := net.ListenUDP("udp4", laddr)
        n, _, err := c.ReadFrom(request)
        if err != nil {
            fmt.Println("[error 4]", err)
            c.Close()
            continue
        }

        fmt.Println("heartbeat recv", n, "bytes")

        c.Close()
    }
}

func heartbeat_server() {
    go heartbeat_lisener()
    for hb_id := uint64(0); ; hb_id ++ {
        SV.Lock()
        if len(SV.Servers) > 0 {
            /*fmt.Print("[heartbeat_server] ")*/
            for _, v := range SV.Servers {
                /*fmt.Print(addr, " ")*/
                /*go send_heartbeat(v[0], v[1])*/
                if (v[0] != SV.LocalIP || v[1] != SV.LocalPort) {
                    go send_heartbeat(v[0], v[1], hb_id)
                }
            }
            /*fmt.Println("")*/
        }

        SV.Unlock()
        time.Sleep(1000 * time.Millisecond)
    }
}

func heartbeat_client(port string) {
    udpaddr, _ := net.ResolveUDPAddr("udp4", "0.0.0.0:" + port)
    for {
        var request []byte = make([]byte, 1024)
        c, _ := net.ListenUDP("udp4", udpaddr)
        //Echo("Heartbeat listenning: " + fmt.Sprint(udpaddr))

        n, _, err := c.ReadFrom(request)
        if err != nil {
            Echo("ReadFrom() err:" + fmt.Sprint(err))
            c.Close()
            continue
        }

        remote_ip, remote_port := parse_binary_addr(request)
        server_cksum := parse_serverlist_cksum(request)
        hb_id := parse_heartbeat_id(request)
        fmt.Print("[heartbeat] recv " + fmt.Sprint(n) + " bytes -> ", remote_ip, ":", remote_port, " server cksum: ", server_cksum, " hb_id: ", hb_id)

        ra, err := net.ResolveUDPAddr("udp4", remote_ip + ":" + remote_port)
        resp := build_heartbeat_resp(0)
        n, err = c.WriteToUDP(resp, ra)
        fmt.Println(", write", n, "bytes, err = ", err)

        c.Close()
    }
}

func RunHeartbeat(port string) {
    go heartbeat_client(port)
    go heartbeat_server()
}

func register_thread(l net.Listener) {
    Echo("register thread launched")
    for {
        conn, _ := l.Accept()
        var request string
        var response string
        fmt.Fscanf(conn, "%s", &request)

        SV.Lock()
        // filter registered same host
        notify_other_masters()
        for addr, _ := range SV.Servers {
            response += addr + "\n"
        }
        SV.AddServer(request)
        SV.Unlock()

        fmt.Fprint(conn, response)
        Echo(SV.Name + " got register request: " + request)
        conn.Close()
    }
}

func notify_other_masters() {
    for addr, addr_tuple := range SV.Servers {
        if addr != SV.LocalIP + ":" + SV.LocalPort {
            fmt.Printf("Notify(Host: %s, Port: %s)\n", addr_tuple[0], addr_tuple[1])
        }
    }
}

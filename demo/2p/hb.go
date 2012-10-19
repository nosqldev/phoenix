/* © Copyright 2012 jingmi. All Rights Reserved.
 *
 * +----------------------------------------------------------------------+
 * | hearbeat mock                                                        |
 * +----------------------------------------------------------------------+
 * | Author: jingmi@gmail.com                                             |
 * +----------------------------------------------------------------------+
 * | Created: 2012-10-14 15:31                                            |
 * +----------------------------------------------------------------------+
 */

package main

import (
    "net"
    "fmt"
    "time"
    /*"encoding/binary"*/
    /*"bytes"*/
    /*"os"*/
)

func main() {
    InitSV("Master0", "127.0.0.1", "20002")

    conn, err := net.Dial("tcp4", "127.0.0.1:10000")
    if err != nil {
        fmt.Println(err)
        return
    }

    fmt.Fprintln(conn, "127.0.0.1:20002")
    for {
        var host_string string
        n, err := fmt.Fscanln(conn, &host_string)
        if n == 0 || err != nil {
            fmt.Println(err)
            break
        }
        fmt.Println("got", host_string)
    }
    conn.Close()

    go heartbeat_client("20002")
    time.Sleep(30 * time.Second)

/*
 *    hb_ip := net.ParseIP("127.0.0.1").To4()
 *    heartbeat_port := int16(20002)
 *    buf := new(bytes.Buffer)
 *    binary.Write(buf, binary.LittleEndian, hb_ip)
 *    binary.Write(buf, binary.LittleEndian, heartbeat_port)
 *
 *    conn, err = net.Dial("udp4", "127.0.0.1:10000")
 *    if err != nil {
 *        fmt.Println("error:", err)
 *    }
 *    conn.Write(buf.Bytes())
 *    fmt.Println(len(buf.Bytes()), len(hb_ip), hb_ip)
 *
 *    var resp []byte = make([]byte, 128)
 *    fmt.Println("ready to read")
 *    n, _, err := conn.(*net.UDPConn).ReadFrom(resp)
 *    fmt.Println("recv", n, "bytes")
 *    if err != nil {
 *        fmt.Println(err)
 *        os.Exit(-1)
 *    }
 *
 *    var ver, status uint32
 *    buf = bytes.NewBuffer(resp)
 *    binary.Read(buf, binary.LittleEndian, &ver)
 *    binary.Read(buf, binary.LittleEndian, &status)
 *    fmt.Println("ver =", ver, "status =", status)
 *
 *    conn.Close()
 */
}

/* vim: set expandtab tabstop=4 shiftwidth=4 foldmethod=marker: */

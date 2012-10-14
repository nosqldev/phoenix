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
    "encoding/binary"
    "bytes"
    //"unsafe"
)

func main() {
    hb_ip := net.ParseIP("127.0.0.99").To4()
    port := int16(1234)
    buf := new(bytes.Buffer)
    binary.Write(buf, binary.LittleEndian, hb_ip)
    binary.Write(buf, binary.LittleEndian, port)

    conn, err := net.Dial("udp4", "127.0.0.1:10001")
    if err != nil {
        fmt.Println("error:", err)
    }
    conn.Write(buf.Bytes())

    fmt.Println(len(buf.Bytes()), len(hb_ip), hb_ip)
}

/* vim: set expandtab tabstop=4 shiftwidth=4 foldmethod=marker: */

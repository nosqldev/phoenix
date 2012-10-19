/* © Copyright 2012 jingmi. All Rights Reserved.
 *
 * +----------------------------------------------------------------------+
 * | share variables                                                      |
 * +----------------------------------------------------------------------+
 * | Author: jingmi@gmail.com                                             |
 * +----------------------------------------------------------------------+
 * | Created: 2012-10-13 23:44                                            |
 * +----------------------------------------------------------------------+
 */

package main

import (
    "sync"
    "net"
    "strconv"
    "bytes"
    "encoding/binary"
    "strings"
    "sort"
)

type ServerInfo struct {
    ID      int
    IP      string
    Port    string
    Role    string
}

type ShareVars struct {
    L sync.Mutex
    Servers map[string][]string // Addr -> [IP,Port]
    ServersTab map[string]ServerInfo // Not used now
    Name string
    LocalIP string
    LocalPort string
    LocalIP_num net.IP
    LocalPort_num uint16 // Master Heartbeat Listen Port
    MasterHeartbeatAddr []byte
    ServerCount int
    ServerList []string
    ServerListCkSum uint64 // Unused right now
}

var SV ShareVars

func InitSV(name, ip, heartbeat_port string) {
    SV.L.Lock()
    SV.Name = name
    SV.LocalIP = ip
    SV.LocalPort = heartbeat_port
    SV.Servers = make(map[string][]string)
    SV.LocalIP_num = net.ParseIP(ip).To4()
    p, _ := strconv.Atoi(heartbeat_port)
    SV.LocalPort_num = uint16(p)
    buf := new(bytes.Buffer)
    binary.Write(buf, binary.LittleEndian, &SV.LocalIP_num)
    binary.Write(buf, binary.LittleEndian, &SV.LocalPort_num)
    SV.MasterHeartbeatAddr = buf.Bytes()
    SV.ServerCount = 0
    SV.ServerList = make([]string, 128)
    SV.AddServer(ip + ":" + heartbeat_port)
    SV.L.Unlock()
}

func (sv *ShareVars) Lock() {
    sv.L.Lock()
}

func (sv *ShareVars) Unlock() {
    sv.L.Unlock()
}

func (sv *ShareVars) AddServer(host string) {
    s_array := strings.Split(host, ":")
    sv.Servers[host] = s_array
    sv.ServerList = append(sv.ServerList, host)
    sv.ServerCount ++
    sort.Strings(sv.ServerList)
}

func (sv *ShareVars) DelServer(host string) {
    delete(sv.Servers, host)
    for i := 0; i < sv.ServerCount; i++ {
        if sv.ServerList[i] == host {
            copy(sv.ServerList[i:], sv.ServerList[i+1:sv.ServerCount])
            sv.ServerCount --
            sv.ServerList = sv.ServerList[0:sv.ServerCount]
            break
        }
    }
}

/* vim: set expandtab tabstop=4 shiftwidth=4 foldmethod=marker: */

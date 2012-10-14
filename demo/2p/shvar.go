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
}

var SV ShareVars

func InitSV(name, ip, port string) {
    SV.L.Lock()
    SV.Name = name
    SV.LocalIP = ip
    SV.LocalPort = port
    SV.Servers = make(map[string][]string)
    SV.L.Unlock()
}

func (sv *ShareVars) Lock() {
    sv.L.Lock()
}

func (sv *ShareVars) Unlock() {
    sv.L.Unlock()
}

/* vim: set expandtab tabstop=4 shiftwidth=4 foldmethod=marker: */

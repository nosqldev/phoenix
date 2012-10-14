/* © Copyright 2012 jingmi. All Rights Reserved.
 *
 * +----------------------------------------------------------------------+
 * | demo                                                                 |
 * +----------------------------------------------------------------------+
 * | Author: jingmi@gmail.com                                             |
 * +----------------------------------------------------------------------+
 * | Created: 2012-10-13 23:12                                            |
 * +----------------------------------------------------------------------+
 */

package main

import (
    "time"
)

func main() {
    Echo("Master0")
    InitSV("Master0", "127.0.0.1", "10000")
    RunListener(10000)
    RunHeartbeat("10001")
    time.Sleep(60 * time.Second)
}

/* vim: set expandtab tabstop=4 shiftwidth=4 foldmethod=marker: */

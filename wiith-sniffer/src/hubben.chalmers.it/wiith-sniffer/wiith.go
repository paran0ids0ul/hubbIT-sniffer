/* This is the WhoIsInTheHubb sniffer application written in Go.
 * The sniffer listens on nearby WiFi traffic and records the
 * visible client MAC-addresses and sends them through a
 * TCP-connection to the backend for data processing.
 *
 * The following signals are handled/recognized:
 * SIGUSR1: Will ignore TCP connection delay and try to reconnect
 *          immediately if the backend is unreachable for some reason.
 *
 * SIGUSR2: Will print the current hitcount. That is the total number
 *          of MAC-addresses seen not excluding duplicates. This is
 *          essentially the number of frames that passed the
 *          capture/display filter.

 * SIGTERM, SIGINT: Will shutdown the sniffer gracefully.
 *
 * Copyright (C) 2013 Emil 'Eda' Edholm (digIT13)
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */
package main

import (
    "flag"
    "fmt"
    "github.com/golang/glog"
    "github.com/ogier/pflag"
    "os"
    "os/signal"
    "os/user"
    "strconv"
    "syscall"
    "encoding/json"
    "net/http"
    "net/url"
)

// The command line flags available
var (
    // Log related flags
    logToStderr  = pflag.BoolP("logtostderr", "e", true, "log to stderr instead of to files")
    logThreshold = pflag.StringP("logthreshold", "t", "INFO", "Log events at or above this severity are logged to standard error as well as to files. Possible values: INFO, WARNING, ERROR and FATAL")
    logdir       = pflag.StringP("logpath", "l", "./logs", "The log files will be written in this directory/path")

    flushInterval = pflag.Int64P("flushinterval", "f", 283, "The flush interval in seconds")
    iface         = pflag.StringP("interface", "i", "mon0", "The capture interface to listen on")
    pcap          = pflag.StringP("pcap", "p", "", "Use a pcap file instead of live capturing")
    server          = pflag.StringP("server", "s", "", "Server to post to")
    hitCount      uint64
)

func init() {
    pflag.Parse()

    // glog Logging options
    flag.Set("logtostderr", strconv.FormatBool(*logToStderr))
    flag.Set("log_dir", *logdir)
    flag.Set("stderrthreshold", *logThreshold)

    if isRoot() {
        glog.Warning("Server run with root privileges! This is uneccessary if tshark has been setup correctly")
    }
}

func main() {
    defer glog.Flush()

    if !InterfaceExists(*iface) && len(*pcap) == 0 {
        glog.Error(*iface + " interface does not exist")
        os.Exit(1)
    }

    // TODO: Init connection to backend first...

    glog.Info("Starting whoIsInTheHubb sniffer")
    defer glog.Info("Shutting down whoIsInTheHubb sniffer...")

    capchan := make(chan *CapturedFrame, 10)
    errchan := make(chan error)

    go func() {
        errchan <- StartTshark(CaptureFilter, capchan)
    }()

    go listenSIGUSR()
    go listenForClients(capchan)
    err := <-errchan // Block until exit...
    if err == nil {
        glog.Info("tshark exited successfully")
        printHitCount()
    }
}

func listenForClients(capchan <-chan *CapturedFrame) {
    for frame := range capchan {
        hitCount++
        b, err := json.Marshal(frame)
         if err != nil {
            fmt.Println("error:", err)
            glog.Error("error: ",err)
        }

        hej := string(b)
        fmt.Println("msg:", hej)
        resp, err := http.Get(*server+"?json="+url.QueryEscape(hej))
       
       //fmt.Println("framen:",resp)
        if err != nil {
             fmt.Println("error:", resp)
            fmt.Println("error:", err)
            glog.Error("error: ", err)
        }else{
            defer resp.Body.Close()
        }
        // TODO: Send through tcp-channel
        // TODO: Maybe impose restrictions on how often "duplicates" are sent
    }
}

// Listen and handle SIGUSR1 and SIGUSR2.
// SIGUSR1 will flush clients and SIGUSR2 will print the clients
// seen since start or last flush to stdout.
func listenSIGUSR() {
    ch := make(chan os.Signal)
    signal.Notify(ch, syscall.SIGUSR1, syscall.SIGUSR2)
    for {
        signal := <-ch
        glog.Info("Caught signal: ", signal)

        switch signal {
        case syscall.SIGUSR1:
            // TODO: Reconnect to backend
            break
        case syscall.SIGUSR2:
            printHitCount()
            break
        }
    }
}

func printHitCount() {
    fmt.Println("Total hitcount:", hitCount)
}

// returns true if run with root privileges, else false.
func isRoot() bool {
    user, err := user.Current()
    if err != nil {
        glog.Error(err.Error())
        return true
    }
    return user.Username == "root"
}

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
    "encoding/json"
    "flag"
    "fmt"
    "github.com/golang/glog"
    "github.com/ogier/pflag"
    "io/ioutil"
    "net/http"
    "os"
    "os/signal"
    "os/user"
    "strconv"
    "strings"
    "syscall"
)

// The command line flags available
var (
    // Log related flags
    logToStderr  = pflag.BoolP("logtostderr", "e", true, "log to stderr instead of to files")
    logThreshold = pflag.StringP("logthreshold", "o", "INFO", "Log events at or above this severity are logged to standard error as well as to files. Possible values: INFO, WARNING, ERROR and FATAL")
    logdir       = pflag.StringP("logpath", "l", "./logs", "The log files will be written in this directory/path")

    flushInterval = pflag.Int64P("flushinterval", "f", 59, "The interval between the PUT of batches of mac-addresses")
    iface         = pflag.StringP("interface", "i", "mon0", "The capture interface to listen on")
    pcap          = pflag.StringP("pcap", "p", "", "Use a pcap file instead of live capturing")
    server        = pflag.StringP("server", "s", "http://localhost:3000/sessions.json", "Server to PUT macs to")
    authToken     = pflag.StringP("token", "t", "", "API-token. (Required)")
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

    if len(*server) == 0 {
        glog.Warning("No server supplied...")
    }

    if len(*authToken) == 0 {
        glog.Warning("No API-token supplied. Will probably fail to PUT mac-addresses to the backend.")
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
        jsonByteBuffer, err := json.Marshal(frame)
        if err != nil {
            glog.Error("Json encoding error:", err)
        }

        jsonMsg := string(jsonByteBuffer)
        fmt.Println("Json:", jsonMsg)

        client := &http.Client{}
        request, err := http.NewRequest("PUT", *server, strings.NewReader(jsonMsg))
        request.ContentLength = int64(len(jsonMsg))
        request.Header.Set("Content-Type", "application/json")

        token := fmt.Sprintf("Token token=\"%v\"", *authToken)
        request.Header.Set("Authorization", token)
        if err != nil {
            glog.Error("PUT request error: ", err)
            continue
        }

        response, err := client.Do(request)
        if err != nil {
            glog.Error("PUT response error: ", err)
            continue
        }
        defer response.Body.Close()

        contents, err := ioutil.ReadAll(response.Body)
        if err != nil {
            glog.Error("Error reading response:", err)
            continue
        }
        fmt.Println("The calculated length is:", len(string(contents)), "for the url:", server)
        fmt.Println("   ", response.StatusCode)
        hdr := response.Header
        for key, value := range hdr {
            fmt.Println("   ", key, ":", value)
        }
        fmt.Println(contents)

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

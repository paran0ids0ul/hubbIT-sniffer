/* This is the WhoIsInTheHubb server application written i Go.
 * The server records nearby WiFi clients (using their MAC-addresses)
 * and stores statistics about them in a database.
 *
 * Will print the current clients on SIGUSR2
 * Will flush and update database of the current
 * known clients on SIGUSR1
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
    "time"
)

//func init() {
//allKnown["b4:07:f9:f3:65:eb"] = User{"Malm", "Phone"}
//allKnown["00:26:08:dd:22:6a"] = User{"Malm", "Computer"}
//allKnown["c8:60:00:3a:53:f3"] = User{"Meddan", "Computer"}
//allKnown["78:d6:f0:df:c1:47"] = User{"Meddan", "Phone"}
//allKnown["20:64:32:5f:33:cc"] = User{"Eda", "Phone"}
//allKnown["08:d4:2b:1a:7e:d4"] = User{"Eda", "Tablet"}
//allKnown["c4:85:08:2c:00:fb"] = User{"Eda", "Computer"}
//allKnown["98:fe:94:4b:d4:f0"] = User{"rekoil", "Computer"}
//allKnown["1c:7b:21:57:4c:40"] = User{"rekoil", "Phone"}
//}

type Client struct {
    FirstSeen time.Time
    LastSeen  time.Time
}

var currClients = make(map[MAC]*Client)

// The command line flags available
var (
    // Log related flags
    logToStderr  = pflag.BoolP("logtostderr", "l", true, "log to stderr instead of to files")
    logThreshold = pflag.StringP("logthreshold", "t", "INFO", "Log events at or above this severity are logged to standard error as well as to files. Possible values: INFO, WARNING, ERROR and FATAL")
    logdir       = pflag.StringP("logpath", "p", "./logs", "The log files will be written in this directory/path")

    flushInterval = pflag.Int64P("flushinterval", "f", 600, "The flush interval in seconds")
    iface         = pflag.StringP("interface", "i", "mon0", "The capture interface to listen on")
)

func init() {
    pflag.Parse()
}

func main() {
    // glog Logging options
    flag.Set("logtostderr", strconv.FormatBool(*logToStderr))
    flag.Set("log_dir", *logdir)
    flag.Set("stderrthreshold", *logThreshold)
    defer glog.Flush()

    if !InterfaceExists(*iface) {
        glog.Error(*iface + " interface does not exist")
        os.Exit(1)
    }

    if isRoot() {
        glog.Warning("Server run with root privileges!")
    }

    glog.Info("Starting whoIsInTheHubb server")
    defer glog.Info("Shutting down whoIsInTheHubb server...")

    capchan := make(chan CapturedFrame, 10)
    errchan := make(chan error)

    go func() {
        errchan <- StartTshark(CaptureFilter, capchan)
    }()

    go listenSIGUSR()
    go flushTimer()
    go listenForClients(capchan)
    err := <-errchan // Block until exit...
    if err == nil {
        glog.Info("tshark exited successfully")
        printClients()
    }
}

func listenForClients(capchan chan CapturedFrame) {
    for frame := range capchan {
        if c, ok := currClients[frame.Mac]; !ok {
            currClients[frame.Mac] = &Client{frame.Timestamp, frame.Timestamp}
        } else {
            // You seem familiar... Update LastSeen
            c.LastSeen = frame.Timestamp
        }
    }
}

// flush the clients after the user specified amount of seconds
func flushTimer() {
    duration := time.Duration(*flushInterval) * time.Second
    for {
        <-time.After(duration)
        flushClients()
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
            flushClients()
            break
        case syscall.SIGUSR2:
            printClients()
            break
        }
    }
}

func printClients() {
    var count int
    for mac, c := range currClients {
        count++
        fmt.Printf("MAC %s {\n\tFirst seen: %v\n\tLast seen:  %v\n}\n", mac, c.FirstSeen, c.LastSeen)

    }
    if count > 0 {
        fmt.Println("Total:", count, "clients")
    }
}

// returns true if run with root priviliges, else false.
func isRoot() bool {
    user, err := user.Current()
    if err != nil {
        glog.Error(err.Error())
        return false
    }
    return user.Username == "root"
}

// Flush means that the amount of seconds seen for each client/mac will be calculated
// and stored in the database.
func flushClients() {
    glog.Info("Flushing current clients...")

    var count uint

    for mac, client := range currClients {
        count++
        // TODO Implement...
        _ = mac
        _ = client
    }

    glog.Info("Flush complete! ", count, " clients flushed.")
}

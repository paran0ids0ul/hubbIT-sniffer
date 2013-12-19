/* This is the WhoIsInTheHubb server application written i Go.
 * The server records nearby WiFi clients (using their MAC-addresses)
 * and stores statistics about them in a database.

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
    "os"
    "os/user"
    "time"
)

const (
    // The interface to monitor
    // TODO: Replace with flag
    iface = "mon0"
)

var allKnown = make(map[MAC]User)

func init() {
    allKnown["b4:07:f9:f3:65:eb"] = User{"Malm", "Phone"}
    allKnown["00:26:08:dd:22:6a"] = User{"Malm", "Computer"}
    allKnown["c8:60:00:3a:53:f3"] = User{"Meddan", "Computer"}
    allKnown["78:d6:f0:df:c1:47"] = User{"Meddan", "Phone"}
    allKnown["20:64:32:5f:33:cc"] = User{"Eda", "Phone"}
    allKnown["08:d4:2b:1a:7e:d4"] = User{"Eda", "Tablet"}
    allKnown["c4:85:08:2c:00:fb"] = User{"Eda", "Computer"}
    allKnown["98:fe:94:4b:d4:f0"] = User{"rekoil", "Computer"}
    allKnown["1c:7b:21:57:4c:40"] = User{"rekoil", "Phone"}
}

type User struct {
    Name, Desc string
}
type Client struct {
    User      *User
    FirstSeen time.Time
    LastSeen  time.Time
}

var (
    currClients = make(map[MAC]*Client)
)

func main() {
    // Logging options
    // TODO: Replace with flags
    flag.Set("logtostderr", "true")
    flag.Set("log_dir", "./logs")
    flag.Set("stderrthreshold", "INFO")
    defer glog.Flush()

    if !InterfaceExists(iface) {
        glog.Error(iface + " interface does not exist")
        os.Exit(1)
    }

    if isRoot() {
        glog.Warning("Server run with root privileges!")
    }

    glog.Info("Starting whoIsInTheHubb server")
    defer glog.Info("Shutting down whoIsInTheHubb server...")
    defer StopTshark()

    capchan := make(chan CapturedFrame)
    errchan := make(chan error)
    // TODO: Handle SIGUSR1, SIGUSR2, flush current and print current clients

    go func() {
        errchan <- StartTshark(DisplayFilter, capchan)
    }()

    go listenForClients(capchan)
    err := <-errchan // Wait for exit...
    if err == nil {
        glog.Info("tshark exited successfully")
        printClients()
    }
}

func listenForClients(capchan chan CapturedFrame) {
    for frame := range capchan {
        if c, ok := currClients[frame.Mac]; !ok {
            client := &Client{nil, frame.Timestamp, frame.Timestamp}
            if user, isKnown := allKnown[frame.Mac]; isKnown {
                client.User = &user
            }
            currClients[frame.Mac] = client
        } else {
            // You seem familiar... Update LastSeen
            c.LastSeen = frame.Timestamp
        }
    }
}

func printClients() {
    var count int
    var first, second string
    for mac, c := range currClients {
        count++
        first = fmt.Sprintf("MAC %s {\n\tFirst seen: %v\n\tLast seen:  %v\n", mac, c.FirstSeen, c.LastSeen)
        if user, exists := allKnown[mac]; exists {
            second = fmt.Sprintf("\tKnown as: %s\n\tDevice:   %s\n}\n", user.Name, user.Desc)
        } else {
            second = "}\n"
        }
        fmt.Print(first + second)

    }
    if count > 0 {
        fmt.Println("Saw", count, "clients totally")
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

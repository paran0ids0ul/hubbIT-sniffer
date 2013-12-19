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
    iface = "wlan0"
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
    clients = make(map[MAC]*Client)
)

func main() {
    // Logging options
    flag.Set("logtostderr", "true")
    flag.Set("log_dir", "./logs")
    flag.Set("stderrthreshold", "INFO")
    defer glog.Flush()

    if !InterfaceExists("mon0") {
        glog.Error("mon0 interface does not exist")
        os.Exit(1)
    }

    if isRoot() {
        glog.Warning("Server run with root privileges!")
    }

    glog.Info("Starting whoIsInTheHubb server")
    defer glog.Info("Stopping whoIsInTheHubb server")
    defer StopTshark()

    capchan := make(chan CapturedFrame)
    errchan := make(chan error)

    go func() {
        errchan <- StartTshark(DisplayFilter, capchan)
    }()

    for frame := range capchan {
        if c, ok := clients[frame.Mac]; !ok {
            client := &Client{nil, frame.Timestamp, frame.Timestamp}
            if user, isKnown := allKnown[frame.Mac]; isKnown {
                client.User = &user
            }
            clients[frame.Mac] = client
            //if client.User != nil {
            //fmt.Println("Saw", client.User.Name, "("+client.User.Desc+")", "first", client.FirstSeen.Format(time.Kitchen), "and last", client.LastSeen.Format(time.Kitchen), "Mac:", frame.Mac)
            //} else {
            //fmt.Println("Saw", frame.Mac)
            //}
        } else {
            // You seem familiar... Update LastSeen
            c.LastSeen = time.Now()
        }

    }
    printClients()

}

func printClients() {
    var count int
    for mac, c := range clients {
        count++
        user := c.User
        if user != nil {
            fmt.Println("Saw", c.User.Name, "("+c.User.Desc+")", "first", c.FirstSeen.Format(time.Kitchen), "and last", c.LastSeen.Format(time.Kitchen), "Mac:", mac)
        } else {
            fmt.Println("Saw", mac, "first", c.FirstSeen, "and last", c.LastSeen)
        }

    }
    fmt.Println("Saw", count, "clients totally")
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

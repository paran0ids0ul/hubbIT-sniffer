package main

import (
    "bufio"
    "fmt"
    "github.com/golang/glog"
    "io"
    "net"
    "os"
    "os/exec"
    "os/signal"
    //"regexp"
    "syscall"
)

type MAC string

var (
    airodumpArgs = []string{"--update", "1",
        "--berlin", "20", "mon0"}
    airodumpCmd = "airodump-ng"
    cmd         *exec.Cmd
    output      io.ReadCloser
)

// Start the airodump-ng process
func StartAirdump() (err error) {
    glog.Info("Starting airodump-ng")

    cmd = exec.Command(airodumpCmd, airodumpArgs...)
    // airdump doesn't seem to use stdout and seem to use stderr for everything
    output, err = cmd.StderrPipe()
    if err != nil {
        glog.Error(err.Error())
        return err
    }

    go signalListen()
    err = cmd.Start()
    if err != nil {
        glog.Error(err.Error())
        return err
    }

    process(nil)

    // Seem to need to flush output buffer before exit...
    err = cmd.Wait()
    if err != nil {
        glog.Error(airodumpCmd + ": " + err.Error())
    }

    return err
}

// Sends SIGINT to airodump-ng. Does not wait for process to end
func StopAirdump() {
    if cmd != nil {
        cmd.Process.Signal(syscall.SIGTERM)
        // TODO: Maybe flush output and wait for exit
    }
}

// Process and parse the output from airodump-ng
// Returns the current visible clients (mac-addresses) on the supplied channel
func process(macChan chan MAC) {
    //var re = regexp.MustCompile(`^([0-9A-F]{2}[:]){5}([0-9A-F]{2})$`)
    reader := bufio.NewReader(output)
    for {
        line, err := reader.ReadString('\n')
        fmt.Println(line)
        if err != nil {
            if err != io.EOF {
                glog.Error(err.Error())
            }
            return
        }
        // FIXME
        //fmt.Print(re.FindString(line))
    }
}

// Turn on monitor mode on the supplied iface
func SetupMonitor(iface string) error {
    // Check if we actually need to enable monitor mode
    mon, _ := net.InterfaceByName("mon0")
    if mon != nil {
        glog.Info("Monitor interace already enabled")
        return nil
    }

    glog.Info("Enabling monitor mode")
    return execAirmon("start", iface)
}

// Disable the monitor mode. Note: should be pased mon0 for example
// instead of wlan0
func TeardownMonitor(iface string) error {
    glog.Info("Disabling monitor mode")
    return execAirmon("stop", iface)
}

// Helper function
func execAirmon(method, iface string) error {
    c := exec.Command("airmon-ng", method, iface)
    return c.Run()
}

func signalListen() {
    // Handle SIGINT and SIGTERM.
    ch := make(chan os.Signal)
    signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
    glog.Warning("Caught signal: ", <-ch)
}

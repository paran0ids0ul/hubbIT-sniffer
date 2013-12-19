// Start and parse tshark output
package main

import (
    "bufio"
    "github.com/golang/glog"
    "io"
    "net"
    "os"
    "os/exec"
    "os/signal"
    "strconv"
    "strings"
    "syscall"
    "time"
)

const (
    tsharkCmd  = "tshark"
    dispFilter = "(wlan.fc.type_subtype eq 4 || wlan.fc.type_subtype eq 2 || wlan.fc.type_subtype eq 0 || (wlan.fc.tods == 1 and wlan.fc.fromds == 0)) and wlan"
    // Capture filter is faster and less CPU-intensive (filters on the kernel level) but may not show all clients as good as the dispFilter
    captureFilter   = "(subtype assocreq or subtype reassocreq or subtype probereq or subtype null)"
    outputSeparator = "|"
)

type Filter int

const (
    DisplayFilter Filter = iota
    CaptureFilter
)

var (
    //tsharkArgs = []string{"-i", "mon0", "-l", "-n", "-T", "fields",
    //    "-e", "wlan.sa", "-e", "frame.time_epoch", "-E", "separator=|"}
    tsharkArgs  = []string{"-r", "/home/eda/moodelsmote.pcap", "-n", "-T", "fields", "-e", "wlan.sa", "-e", "frame.time_epoch", "-E", "separator=|"}
    dispArgs    = []string{"-2", "-R", dispFilter}
    captureArgs = []string{"-f", captureFilter}
    cmd         *exec.Cmd
)

type MAC string
type CapturedFrame struct {
    Mac       MAC
    Timestamp time.Time
}

// Start the airodump-ng process
func StartTshark(filter Filter, capchan chan CapturedFrame) (err error) {
    glog.Info("Starting tshark")

    var args []string
    if filter == DisplayFilter {
        args = append(tsharkArgs, dispArgs...)
    } else {
        args = append(tsharkArgs, captureArgs...)
    }
    cmd = exec.Command(tsharkCmd, args...)
    //stderr, err := cmd.StderrPipe()
    //if err != nil {
    //glog.Error(err.Error())
    //return err
    //}
    output, err := cmd.StdoutPipe()
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

    // TODO: Print stderr when needed...
    //go io.Copy(os.Stderr, stderr)
    process(output, capchan)

    err = cmd.Wait()
    if err != nil {
        glog.Error(tsharkCmd + ": " + err.Error())
    }

    return err
}

// Sends SIGINT to tshark.
func StopTshark() {
    if cmd != nil {
        cmd.Process.Signal(syscall.SIGINT)
    }
}

// Process and parse the output from airodump-ng
// Returns the current visible clients (mac-addresses) on the supplied channel
func process(output io.ReadCloser, capchan chan CapturedFrame) {
    scanner := bufio.NewScanner(output)
    var cur, prev MAC
    var timestamp string

    for scanner.Scan() {
        // Note that we need to include the frame timestamp since we don't know
        // when the message will be read from the channel.
        // (Also helps when reading pcap files)
        cur, timestamp = splitOutput(scanner.Text())
        if cur != prev {
            capchan <- CapturedFrame{cur, parseEpoch(timestamp)}
            prev = cur
        }
    }
    if err := scanner.Err(); err != nil {
        glog.Warning("reading standard input:" + err.Error())
    }
    // We're done with the channel at this point
    close(capchan)
}

// This assumes that the split will produce exactly 2 words
// Convenience function
func splitOutput(line string) (mac MAC, date string) {
    slice := strings.Split(line, outputSeparator)
    mac = MAC(slice[0])
    date = slice[1]
    return
}

// Asumes this format: 1387372947.665215000
// Where seconds.nanoseconds since the epoch
func parseEpoch(epoch string) time.Time {
    slice := strings.Split(epoch, ".")
    sec, err := strconv.ParseInt(slice[0], 10, 64)
    if err != nil {
        glog.Error(err.Error())
    }
    nsec, err := strconv.ParseInt(slice[1], 10, 64)
    if err != nil {
        glog.Error(err.Error())
    }
    // FIXME i'm broken
    return time.Unix(sec, nsec)
}

// Whether or not the monitor interface exist
func InterfaceExists(iface string) bool {
    if mon, _ := net.InterfaceByName(iface); mon != nil {
        return true
    }
    return false
}

func signalListen() {
    // TODO: Handle SIGUSR1, SIGUSR2, reset current and close wireshark respectivly
    // Handle SIGINT and SIGTERM.
    ch := make(chan os.Signal)
    signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
    glog.Warning("Caught signal: ", <-ch)
    StopTshark()
}

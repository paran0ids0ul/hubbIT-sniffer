/* This is where the magic happens. This module is responsible for
 * starting and parsing the output of tshark (the application that
 * does the actual sniffing)
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
    "bufio"
    "bytes"
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
    dispFilter = "(wlan.fc.type_subtype eq 4 || wlan.fc.type_subtype eq 2 || wlan.fc.type_subtype eq 0 || (wlan.fc.type_subtype eq 0x24 and wlan.fc.tods == 1 and wlan.fc.fromds == 0)) and wlan"
    // Capture filter is faster and less CPU-intensive (filters on the kernel level) but may not show all clients as good as the dispFilter
    captureFilter   = "(subtype assocreq or subtype reassocreq or subtype probereq or (subtype null and dir tods))"
    outputSeparator = "|"
)

type FilterType int

// Serves as a enumerator type for choosing the filter type.
const (
    DisplayFilter FilterType = iota
    CaptureFilter
)

var (
    tsharkLiveArgs = []string{"-i", *iface, "-p", "-l", "-n", "-T", "fields",
        "-e", "wlan.sa", "-e", "frame.time_epoch", "-E", "separator=|"}
    tsharkPcapArgs = []string{"-n", "-T", "fields", "-e", "wlan.sa", "-e", "frame.time_epoch", "-E", "separator=|"}
    dispArgs       = []string{"-2", "-R", dispFilter}
    captureArgs    = []string{"-f", captureFilter}
    cmd            *exec.Cmd
)

type MAC string
type CapturedFrame struct {
    Mac       MAC
    Timestamp time.Time
}

// Start the tshark process. The captured mac-addresses will be
// returned on the supplied channel.
func StartTshark(filter FilterType, capchan chan<- *CapturedFrame) (err error) {
    glog.Info("Starting tshark")

    // Decide whether to use live recording or a pcap file (debugging purposes)
    var tsharkArgs []string
    if len(*pcap) != 0 {
        file := []string{"-r", *pcap}
        tsharkArgs = append(file, tsharkPcapArgs...)
    } else {
        tsharkArgs = tsharkLiveArgs
    }

    var args []string
    if filter == DisplayFilter || len(*pcap) != 0 {
        args = append(tsharkArgs, dispArgs...)
    } else {
        args = append(tsharkArgs, captureArgs...)
    }
    cmd = exec.Command(tsharkCmd, args...)

    var stderr bytes.Buffer
    cmd.Stderr = &stderr

    output, err := cmd.StdoutPipe()
    if err != nil {
        glog.Error(err.Error())
        return err
    }

    go listenCloseSignal()
    err = cmd.Start()
    if err != nil {
        glog.Error(err.Error())
        return err
    }

    process(output, capchan)

    err = cmd.Wait()
    if err != nil {
        glog.Error(tsharkCmd + ": " + err.Error())
        if errstr := stderr.String(); errstr != "<nil>" {
            glog.Error(errstr)
        }
    }

    return err
}

// Sends SIGINT to tshark.
func StopTshark() {
    if cmd != nil {
        cmd.Process.Signal(syscall.SIGINT)
    }
}

// Process and parse the output from tshark
// Assumes the format of: macaddress|epoch_timestamp
// E.g. 38:aa:3c:3e:f2:da|1387487630.925985000
// Puts the current visible clients (mac-addresses) onto the supplied channel
func process(output io.ReadCloser, capchan chan<- *CapturedFrame) {
    var (
        scanner   = bufio.NewScanner(output)
        cur       MAC
        timestamp string
    )

    for scanner.Scan() {
        // Note that we need to include the frame timestamp since we don't know
        // when the message will be read from the channel.
        // (Also helps when reading pcap files)
        cur, timestamp = splitOutput(scanner.Text())
        capchan <- &CapturedFrame{cur, parseEpoch(timestamp)}
    }
    if err := scanner.Err(); err != nil {
        glog.Warning("reading standard output:" + err.Error())
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

// Assumes this format: 1387372947.665215000
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
    return time.Unix(sec, nsec)
}

// Whether or not the supplied interface exist. Does not check whether
// or not tshark can actually capture from this interface.
func InterfaceExists(iface string) bool {
    if mon, _ := net.InterfaceByName(iface); mon != nil {
        return true
    }
    return false
}

// Listen and handle SIGINT and SIGTERM.
func listenCloseSignal() {
    ch := make(chan os.Signal)
    signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
    glog.Warning("Caught signal: ", <-ch)
    StopTshark()
}

package main

import (
    "bufio"
    "fmt"
    "github.com/golang/glog"
    "io"
    "os"
    "os/exec"
)

type MAC string

var (
    airodumpArgs = []string{"--update", "1",
        "--berlin", "20", "mon0"}
    airodumpCmd = "airodump-ng"
    cmd         *exec.Cmd
    stdoutPipe  io.ReadCloser
)

func init() {

}

// Start the airodump-ng process
func start() (err error) {
    glog.Info("Starting airodump-ng")

    cmd = exec.Command(airodumpCmd, airodumpArgs...)
    // airdump does not seem to use stdout and uses stderr for everything
    stdoutPipe, err = cmd.StderrPipe()
    if err != nil {
        glog.Error(err.Error())
        return err
    }

    err = cmd.Start()
    if err != nil {
        glog.Error(err.Error())
        return err
    }

    scanner := bufio.NewScanner(stdoutPipe)
    for scanner.Scan() {
        fmt.Println(scanner.Text())
    }

    if err := scanner.Err(); err != nil {
        glog.Error(err.Error())
    }

    err = cmd.Wait()
    if err != nil {
        glog.Error(airodumpCmd + ": " + err.Error())
    }

    return err
}

func stop() error {
    if cmd != nil {
        cmd.Process.Signal(os.Interrupt)
        err := cmd.Wait()
        if err != nil {
            return err
        }
    }
    return nil
}

func readAll(io.ReadCloser) (out string) {
    scanner := bufio.NewScanner(stdoutPipe)
    for scanner.Scan() {
        out += scanner.Text()
    }

    if err := scanner.Err(); err != nil {
        glog.Error(err.Error())
    }
    return out
}

// Process and parse the output from airodump-ng
// Returns the current visible clients (mac-addresses)
func process() []MAC {
    var clients []MAC
    clients = make([]MAC, 0, 100)

    return clients
}

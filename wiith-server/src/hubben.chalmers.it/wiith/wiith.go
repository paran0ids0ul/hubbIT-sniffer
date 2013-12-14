package main

import (
    "flag"
    "github.com/golang/glog"
    "os"
    "os/user"
)

func main() {
    // Logging options
    flag.Set("logtostderr", "false")
    flag.Set("log_dir", "./logs")
    flag.Set("stderrthreshold", "INFO")
    defer glog.Flush()

    if !isRoot() {
        glog.Error("Server not run with root privileges")
        os.Exit(1)
    }

    glog.Info("Starting whoIsInTheHubb server")
    defer glog.Info("Stopping whoIsInTheHubb server")
    start()
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

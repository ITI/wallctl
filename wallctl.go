package main

import (
    "fmt"
    "flag"
    "os"

)

var (
    who, configfile, ntiaddr, serport, ipbcast string
    ntiport int
    pon, poff, debugnti, debugwall bool
)


func main() {
    flag.BoolVar(&pon, "on", false, "Power on the wall panels")
    flag.BoolVar(&poff, "off", false, "Power off the wall panels")
    flag.StringVar(&ipbcast, "b", "192.168.80.255", "Broadcast address for WOL")
    flag.StringVar(&who, "c", "", "Config section to load")
    flag.StringVar(&configfile, "f", "./wall.conf", "Config file to load")
    flag.StringVar(&ntiaddr, "i", "192.168.80.4",
        "IP Address of the NTI switch")
    flag.IntVar(&ntiport, "p", 2005, "TCP Port on the NTI switch")
    flag.StringVar(&serport, "s", "/dev/ttyS0", "Serial port of the videowall")
    flag.BoolVar(&debugnti, "dn", false, "Run video switch in debug mode")
    flag.BoolVar(&debugwall, "dw", false, "Run wall in debug mode")
    flag.Parse()

    config, err := parseConfig(configfile)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Config file broke: %v\n", err)
        os.Exit(2)
    }


    if pon || poff {
        err := doPower(&config.Config)
        if err != nil {
            fmt.Fprintf(os.Stderr, "Power operation failed: %v\n", err)
            os.Exit(6)
        } else {
            os.Exit(0)
        }
    }


    for _, layout := range config.Layouts {
        if layout.Name == who {
            err = configWall(config.Config, layout)
            if err != nil {
                fmt.Fprintf(os.Stderr, "Something went terribly wrong: %v\n", err)
                os.Exit(5)
            }
            os.Exit(0)
        }
    }
    // If we reach here, Error
    fmt.Fprintf(os.Stderr, "Configuration %s not found\n", who)
    os.Exit(4)
}

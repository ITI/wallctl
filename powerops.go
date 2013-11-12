package main

import (
    "errors"

    libwall "github.com/iti/libwall"

    wol "github.com/ghthor/gowol"
)

func doPower(config *Config) (error) {
    if pon && poff {
        return errors.New("Either on or off, not both")
    }

    if poff {
        allpanels := libwall.NewPanel(byte(0xfe), serport, debugwall)
        err := allpanels.Set("power", libwall.OFF)
        if err != nil {
            return err
        }
    }
    if pon {
        for _,p := range config.Outputs {
            if m := p.MacAddr; m != "" {
                err := wol.SendMagicPacket(m, ipbcast)
                if err != nil {
                    return err
                }
            }
        }
    }

    return nil
}

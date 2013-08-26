package main

import (
    "errors"

    wol "github.com/ghthor/gowol"
)

func doPower(config *Config) (error) {
    if pon && poff {
        return errors.New("Either on or off, not both")
    }



    for _,p := range config.Outputs {
        if pon {
            if m := p.MacAddr; m != "" {
                err := wol.SendMagicPacket(m, ipbcast)
                if err != nil {
                    return err
                }
            }
        }
        if poff {
        }
    }


    return nil
}

package main

import (
    "fmt"
    "os"
    "time"
    "errors"

    libnti "github.com/iti/libnti"
    libwall "github.com/iti/libwall"
)


func configWall(config Config, layout Layout) (error) {
    /*  Steps to configure the wall
        1. Setup video switch
        2. Create All Panels
        2b. Set all Panels to DVI/Wall off
        3. Set Sources
        4. Build Walls (using panels from 2)
        5. For each wall, wall on
    */

    vswitch := libnti.Veemux{
        IP: ntiaddr,
        Port: ntiport,
        Debug: debugnti,
    }

    // Step 1: Set up video switch
    for _, input := range layout.VideoSwitch {
        iid, err := getInputID(input.Name, config.Inputs)
        if err != nil {
            fmt.Fprintf(os.Stderr, "Video Switch Error: %v\n", err)
        }

        for _, x := range input.Outputs {
            oid, err := getOutputID(x, config.Outputs)
            if err != nil {
                fmt.Fprintf(os.Stderr, "Video Switch Error: %v\n", err)
            }
            fmt.Printf("ConnectPort %v,%v\n", iid, oid)
            vswitch.SendCommand("ConnectSource",
                fmt.Sprintf("%v,%v", iid, oid))
        }
        fmt.Printf("\n")
    }

    // Step 2: create all panels and set source to DVI

    allpanels := libwall.NewPanel(byte(0xfe), serport, debugwall)
    err := allpanels.Set("wall", libwall.OFF)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error resetting panels: %v", err)
    }
    err = allpanels.Set("source", libwall.Sources["dvi"])
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error resetting panels: %v", err)
    }
    time.Sleep(3 * time.Second)

    panels := make(map[string]*libwall.Panel)
    for _,p := range config.Outputs {
        if p.WallID != 0 {
            panel := libwall.NewPanel(byte(p.WallID), serport, debugwall)
            panels[p.Name] = panel
        }
    }


    // Step 3
    for _,p := range layout.PanelConfigs {
        panel := panels[p.Name]
        if panel != nil {
            panel.Set("source", libwall.Sources[p.Source])
        }
    }

    // Step 4/5
    for _,wall := range layout.VideoWalls {
        w := new(libwall.Wall)
        for _,p := range wall.Panels {
            if ok := panels[p.Name]; ok != nil {
                panels[p.Name].X = byte(wall.Size.X)
                panels[p.Name].Y = byte(wall.Size.Y)
                panels[p.Name].Position = byte(p.Position)
                w.Panels = append(w.Panels, panels[p.Name])
            }
        }
        if err := w.On(); err != nil {
            return err
        }
    }
    return nil
}



// Need interface here??
func getInputID(input string, config []Input) (int, error) {
    for _,i := range config {
        if i.GetName() == input {
            return i.GetPort(), nil
        }
    }
    return 0, errors.New(fmt.Sprintf("Input %v not found", input))
}

func getOutputID(output string, config []Output) (int, error) {
    for _,i := range config {
        if i.GetName() == output {
            return i.GetPort(), nil
        }
    }
    return 0, errors.New(fmt.Sprintf("Input %v not found", output))
}

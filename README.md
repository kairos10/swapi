# golang API for Skywatcher(c) wifi modules or AZGTi mounts
This GO module is usefull to control a wifi-aware *Skywatcher(c)* mount over wifi.

The module implements a couple of SW methods to control the motors (**SW** functions) and a couple of combo methods to control the mount in a more complex manner.

As of now, the API is not aware of the sky coordinates

# usage
To find all available mounts available on you network, use the **FindMounts()** function:
```
package main

import (
        "fmt"
        "github.com/kairos10/swapi/motor/wifi"
)

func main() {
        mounts := wifi.FindMounts()
        if len(mounts) > 0 {
                for i := 0; i < len(mounts); i++ {
                        m := mounts[i]
                        fmt.Printf("found: %s:%d, time[%s]\n", m.UDPAddr.IP, m.UDPAddr.Port, m.DiscoveryTime)
                }
        } else {
                fmt.Println("nothing found")
        }
}
```

The **RetrieveMountParameters()** method could be used to retrieve the mount's parameters:
```
...
                        err := m.RetrieveMountParameters()
                        if err != nil {
                                fmt.Println("RetrieveMountParameters error: ", err)
                        }
                        fmt.Printf("Mount parameters: counts/revolution=%v timerFrequency=%d\n", m.MCParamCPR, m.MCParamFrequency)
...
```

To slew the RA axis CCW, at a medium speed, for one second:
```
...
                        err = m.SetSlewRate(wifi.AXIS_RA_AZ, -wifi.SLEW_SPEED_5, 1000 * time.Millisecond)
                        if err != nil {
                                fmt.Println("SetSlewRate error: ", err)
                        } else {
                                fmt.Println("SetSlewRate: done")
                        }
...
```

# golang API for Skywatcher(c) wifi modules or AZGTi mounts
This GO module is usefull to control a wifi-aware *Skywatcher(c)* mount over wifi.

The module implements a couple of SW methods to control the motors (**SW** functions) and a couple of combo methods to control the mount in a more complex manner.

As of now, the API is not aware of the terrestrial or sky coordinates; therefore, tracking is only possible in EQ mode and the slew rate does not take into consideration the atmospheric diffraction

**Note:** I'm in no way related to Skywatcher(c) and this API is based only on my very limited experimentation with the protocols involved. Please be aware that using this API for any reason may have unexpected consequences

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

GOTO to a specific tick position
```
			err = m.GoToPosition(wifi.AXIS_RA_AZ, 0)
```

Move an axis for a specified number of ticks
```
			err = m.GoToRelativeIncrement(wifi.AXIS_RA_AZ, m.MCParamCPR/360) // 1 angular degree gets transformed to [MCParamCPR/360] ticks
```

# Samples

## app/sw_joy.go
controls a SW mount over wifi, from a gamepad connected to your computer
```
go run app/sw_joy.go
```

## app/equtil.go
issues simple commands to a wifi connected motor controller (meridian flip, reset home position, etc)
```
go run app/equtil.go [command]
```

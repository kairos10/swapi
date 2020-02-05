speed0
[G]SetMotionMode   ALT/DEC    slew/+    10 --> []
[I]SetStepPeriod   ALT/DEC    1329693    --> []

speed1
[I]SetStepPeriod   ALT/DEC    664846     --> []

speed2
[I]SetStepPeriod   ALT/DEC    83106      --> []

speed3
[I]SetStepPeriod   ALT/DEC    41553      --> []

speed4
[I]SetStepPeriod   ALT/DEC    20776      --> []

speed5
[I]SetStepPeriod   ALT/DEC    10388      --> []

speed6
[I]SetStepPeriod   ALT/DEC    5194       --> []

speed7
[G]SetMotionMode   ALT/DEC    SLEW/+     30 --> []
[I]SetStepPeriod   ALT/DEC    1662       --> []

speed8
[I]SetStepPeriod   ALT/DEC    1108       --> []

speed9
[I]SetStepPeriod   ALT/DEC    831        --> []



[j]GetAxisPosition AZ/RA                 --> 1009
[j]GetAxisPosition ALT/DEC               --> 102504


up speed5
[f]GetAxisStatus   ALT/DEC               --> [111]
[K]AxisStop (NI)   ALT/DEC               --> []
[f]GetAxisStatus   ALT/DEC               --> [101]

[G]SetMotionMode   ALT/DEC    slew/+     --> []
[I]SetStepPeriod   ALT/DEC    10388      --> []
[J]StartMotion     ALT/DEC               --> []
[K]AxisStop (NI)   ALT/DEC               --> []

up speed9

righsStatus   ALT/DEC               --> [301]
[G]SetMotionMode   ALT/DEC    SLEW/+     --> []
[I]SetStepPeriod   ALT/DEC    831        --> []
[J]StartMotion     ALT/DEC               --> []
[K]AxisStop (NI)   ALT/DEC               --> []

speed5
[f]GetAxisStatus   AZ/RA                 --> [101]
[G]SetMotionMode   AZ/RA      slew/+     --> []
[I]SetStepPeriod   AZ/RA      10388      --> []
[J]StartMotion     AZ/RA                 --> []
[K]AxisStop (NI)   AZ/RA                 --> []

down speed5
[f]GetAxisStatus   ALT/DEC               --> [101]
[G]SetMotionMode   ALT/DEC    slew/-     --> []
[I]SetStepPeriod   ALT/DEC    10388      --> []
[J]StartMotion     ALT/DEC               --> []
[K]AxisStop (NI)   ALT/DEC               --> []

left speed5
[f]GetAxisStatus   AZ/RA                 --> [101]
[G]SetMotionMode   AZ/RA      slew/-     --> []
[I]SetStepPeriod   AZ/RA      10388      --> []
[J]StartMotion     AZ/RA                 --> []
[K]AxisStop (NI)   AZ/RA                 --> []


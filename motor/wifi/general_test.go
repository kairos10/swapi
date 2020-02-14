package wifi_test

import "github.com/kairos10/swapi/motor/wifi"

func ExampleMount_Resolve() {
	m := new(wifi.Mount)
	m.Resolve("192.168.4.1", -1) // port defaults to 11880 [= SW_UDP_PORT]

}

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"

	"tinygo.org/x/bluetooth"
	"tinygo.org/x/bluetooth/rawterm"
)

type BLE_COMMAND struct {
	Command string `json:command`
}

type concentractor struct {
	Ssid                   string `json:ssid`
	Password               string `json:password`
	Network_server_address string `json:network_server_address`
	Freq_plan              string `json:freq_plan`
}

const WIFI_CONFIG_PATH string = "/etc/wpa_supplicant/wpa_supplicant.conf"
const CONCENTRACTOR_CONFIG_PATH string = "./siliq_lorawan_concentractor_conf.json"
const LORA_GLOBAL_CONFIG_PATH string = "/home/pi/lora/packet_forwarder/lora_pkt_fwd/global_conf.json"

// Reading files requires checking most calls for errors.
// This helper will streamline our error checks below.
func check(e error) {
	if e != nil {
		panic(e)
	}
}

func setFrequencyPlan(freq_plan string) {
	// Perhaps the most basic file reading task is
	// slurping a file's entire contents into memory.
	targetFreqPlanData, err := ioutil.ReadFile("/home/pi/lora/lorasdk/global_conf_" + freq_plan + ".json")
	check(err)
	targetFreqPlanDataSubstrings := strings.Split(string(targetFreqPlanData), `"gateway_conf":`)

	globalConfig, err := ioutil.ReadFile(LORA_GLOBAL_CONFIG_PATH)
	check(err)
	globalConfigSubstrings := strings.Split(string(globalConfig), `"gateway_conf":`)

	newStrings := targetFreqPlanDataSubstrings[0] + `"gateway_conf":` + globalConfigSubstrings[1]
	// fmt.Println(newStrings) // 0
	ioutil.WriteFile(LORA_GLOBAL_CONFIG_PATH, []byte(newStrings), 0644)
	check(err)
	// json.Unmarshal(dat, &jsonObj)
	// fmt.Print(jsonObj["gateway_conf"].(string))

}

func setNetworkServerAddress(network_server_address string) {

	// Perhaps the most basic file reading task is
	// slurping a file's entire contents into memory.
	dat, err := ioutil.ReadFile(LORA_GLOBAL_CONFIG_PATH)
	check(err)
	substrings := strings.Split(string(dat), "server_address")
	indexComma := strings.Index(substrings[1], ",")
	// fmt.Println(indexComma) // 0
	substring := substrings[1][indexComma:len(substrings[1])]
	// fmt.Println(substring) // 0
	substring = substrings[0] + `server_address": "` + network_server_address + `"` + substring
	// fmt.Println(substring) // 0
	ioutil.WriteFile(LORA_GLOBAL_CONFIG_PATH, []byte(substring), 0644)
	check(err)
}

func setWifi(ssid string, password string) {
	// Perhaps the most basic file reading task is
	// slurping a file's entire contents into memory.
	currentData, err := ioutil.ReadFile(WIFI_CONFIG_PATH)
	check(err)
	substrings := strings.Split(string(currentData), "network")
	indexComma := strings.Index(substrings[1], "}")
	// fmt.Println(indexComma) // 0
	substring := substrings[1][indexComma:len(substrings[1])]
	// fmt.Println(substring) // 0
	substring = substrings[0] + `network={
        ssid="` + ssid + `"
        psk="` + password + `"
        key_mgmt=WPA-PSK
` + substring
	// fmt.Println(substring) // 0
	ioutil.WriteFile(WIFI_CONFIG_PATH, []byte(substring), 0644)
	check(err)
}

func bleCommandExecute(bleMessage string) {
	currentData, err := ioutil.ReadFile(CONCENTRACTOR_CONFIG_PATH)
	check(err)
	var concentractor concentractor
	err = json.Unmarshal([]byte(currentData), &concentractor)

	// fmt.Println("concentractorCurrentStatus", concentractor)

	if err != nil {
		fmt.Println("JsonToMapDemo err: ", err)
	}

	var bleMessageMap map[string]interface{}
	err = json.Unmarshal([]byte(bleMessage), &bleMessageMap)
	if err != nil {
		fmt.Println("JsonToMapDemo err: ", err)
	}

	for k, v := range bleMessageMap {
		switch k {
		case "status":
			if fmt.Sprint(v) == "read" {
				statusString, err := json.Marshal(concentractor)
				if err != nil {
					fmt.Println(err)
					return
				}
				fmt.Println(string(statusString))
			}
		case "ssid":
			concentractor.Ssid = fmt.Sprint(v)
			setWifi(concentractor.Ssid, concentractor.Password)
			fmt.Println("Ssid setting ok")
		case "password":
			concentractor.Password = fmt.Sprint(v)
			setWifi(concentractor.Ssid, concentractor.Password)
			fmt.Println("Password setting ok")
		case "network_server_address":
			concentractor.Network_server_address = fmt.Sprint(v)
			setNetworkServerAddress(concentractor.Network_server_address)
			fmt.Println("Network server address setting ok")
		case "freq_plan":
			concentractor.Freq_plan = fmt.Sprint(v)
			setFrequencyPlan(concentractor.Freq_plan)
			fmt.Println("Freq_plan setting ok")
		default:
			fmt.Println("error")
		}
	}
	statusString, err := json.Marshal(concentractor)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(statusString))
	ioutil.WriteFile("./siliq_lorawan_concentractor_conf.json", []byte(statusString), 0644)
	check(err)
}

var (
	serviceUUID = bluetooth.ServiceUUIDNordicUART
	rxUUID      = bluetooth.CharacteristicUUIDUARTRX
	txUUID      = bluetooth.CharacteristicUUIDUARTTX
)

func must(action string, err error) {
	if err != nil {
		panic("failed to " + action + ": " + err.Error())
	}
}

func main() {
	println("starting")
	bleCommandExecute(`
	{
		"ssid": "livingroom",
		"password": "12345678",
		"network_server_address": "192.168.1.104",
		"freq_plan": "EU868"
	}
	`)
	adapter := bluetooth.DefaultAdapter
	must("enable BLE stack", adapter.Enable())
	adv := adapter.DefaultAdvertisement()
	must("config adv", adv.Configure(bluetooth.AdvertisementOptions{
		LocalName:    "NUS", // Nordic UART Service
		ServiceUUIDs: []bluetooth.UUID{serviceUUID},
	}))
	must("start adv", adv.Start())

	var rxChar bluetooth.Characteristic
	var txChar bluetooth.Characteristic
	must("add service", adapter.AddService(&bluetooth.Service{
		UUID: serviceUUID,
		Characteristics: []bluetooth.CharacteristicConfig{
			{
				Handle: &rxChar,
				UUID:   rxUUID,
				Flags:  bluetooth.CharacteristicWritePermission | bluetooth.CharacteristicWriteWithoutResponsePermission,
				WriteEvent: func(client bluetooth.Connection, offset int, value []byte) {
					txChar.Write(value)
					for _, c := range value {
						rawterm.Putchar(c)
					}
				},
			},
			{
				Handle: &txChar,
				UUID:   txUUID,
				Flags:  bluetooth.CharacteristicNotifyPermission | bluetooth.CharacteristicReadPermission,
			},
		},
	}))

	rawterm.Configure()
	defer rawterm.Restore()
	print("NUS console enabled, use Ctrl-X to exit\r\n")
	var line []byte
	for {
		ch := rawterm.Getchar()
		rawterm.Putchar(ch)
		line = append(line, ch)
		// Send the current line to the central.
		if ch == '\x18' {
			// The user pressed Ctrl-X, exit the terminal.
			break
		} else if ch == '}' {
			rxbuf := line // copy buffer
			// Reset the slice while keeping the buffer in place.
			line = line[:0]
			fmt.Println("123" + string(rxbuf))
			// Send the sendbuf after breaking it up in pieces.
			//                      for len(sendbuf) != 0 {
			//                              // Chop off up to 20 bytes from the sendbuf.
			//                              partlen := 20
			//                              if len(sendbuf) < 20 {
			//                                      partlen = len(sendbuf)
			//                              }
			//                              part := sendbuf[:partlen]
			//                              sendbuf = sendbuf[partlen:]
			//                              // This also sends a notification.
			//                              _, err := txChar.Write(part)
			//                              must("send notification", err)
			//                      }
		}
	}
}

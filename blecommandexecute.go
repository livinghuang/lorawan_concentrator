package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"tinygo.org/x/bluetooth"
)

type BLE_COMMAND struct {
	Command string `json:command`
}

type concentractor struct {
	Ssid                   string `json:Ssid`
	Password               string `json:Password`
	Network_server_address string `json:Network_server_address`
	Freq_plan              string `json:Freq_plan`
}

const WIFI_CONFIG_PATH string = "/etc/wpa_supplicant/wpa_supplicant.conf"
const CONCENTRACTOR_CONFIG_PATH string = "./siliq_lorawan_concentractor_conf.json"
const LORA_GLOBAL_CONFIG_PATH string = "/home/pi/lora/packet_forwarder/lora_pkt_fwd/global_conf.json"
const LOCAL_LORA_GLOBAL_CONFIG_PATH string = "/home/pi/lora/packet_forwarder/lora_pkt_fwd/local_conf.json"

var rxChar bluetooth.Characteristic
var txChar bluetooth.Characteristic
var ble_name string

// Reading files requires checking most calls for errors.
// This helper will streamline our error checks below.
func check(e error) string {
	if e != nil {
		ErrorLogger.Println(e)
		return "error"
		// panic(e)
	}
	return "ok"
}

func set_BLE_NAME() error {

	// Perhaps the most basic file reading task is
	// slurping a file's entire contents into memory.
	dat, err := ioutil.ReadFile(LOCAL_LORA_GLOBAL_CONFIG_PATH)
	check(err)
	substrings := strings.Split(string(dat), "gateway_ID")
	var i int
	i = 0
	for ; i < len(substrings); i++ {
		if (substrings[1][i] >= 48) && (substrings[1][i] <= 57) || (substrings[1][i] >= 65) && (substrings[1][i] <= 70) || (substrings[1][i] >= 97) && (substrings[1][i] <= 102) {
			break
		}
	}
	substring := substrings[1][i+10 : i+16]
	fmt.Println(substring)
	ble_name = "SILIQ-" + substring
	check(err)
	return err
}

func setFrequencyPlan(freq_plan string) error {
	// Perhaps the most basic file reading task is
	// slurping a file's entire contents into memory.
	targetFreqPlanData, err := ioutil.ReadFile("/home/pi/lora/lorasdk/global_conf_" + freq_plan + ".json")
	check(err)
	targetFreqPlanDataSubstrings := strings.Split(string(targetFreqPlanData), `"gateway_conf":`)

	globalConfig, err := ioutil.ReadFile(LORA_GLOBAL_CONFIG_PATH)
	check(err)
	globalConfigSubstrings := strings.Split(string(globalConfig), `"gateway_conf":`)

	newStrings := targetFreqPlanDataSubstrings[0] + `"gateway_conf":` + globalConfigSubstrings[1]
	//fmt.Println(newStrings) // 0
	ioutil.WriteFile(LORA_GLOBAL_CONFIG_PATH, []byte(newStrings), 0644)
	check(err)
	return err
}

func setNetworkServerAddress(network_server_address string) error {

	// Perhaps the most basic file reading task is
	// slurping a file's entire contents into memory.
	dat, err := ioutil.ReadFile(LORA_GLOBAL_CONFIG_PATH)
	check(err)
	substrings := strings.Split(string(dat), "server_address")
	indexComma := strings.Index(substrings[1], ",")
	substring := substrings[1][indexComma:len(substrings[1])]
	substring = substrings[0] + `server_address": "` + network_server_address + `"` + substring
	//fmt.Println(substring) // 0
	ioutil.WriteFile(LORA_GLOBAL_CONFIG_PATH, []byte(substring), 0644)
	check(err)
	return err
}

func setWifi(ssid string, password string) error {
	// Perhaps the most basic file reading task is
	// slurping a file's entire contents into memory.
	currentData, err := ioutil.ReadFile(WIFI_CONFIG_PATH)
	check(err)
	substrings := strings.Split(string(currentData), "network")
	indexComma := strings.Index(substrings[1], "}")
	substring := substrings[1][indexComma:len(substrings[1])]
	substring = substrings[0] + `network={
        ssid="` + ssid + `"
        psk="` + password + `"
        key_mgmt=WPA-PSK
` + substring
	ioutil.WriteFile(WIFI_CONFIG_PATH, []byte(substring), 0644)
	check(err)
	return err
}

func bleCommandExecute(bleMessage string) {
	currentData, err := ioutil.ReadFile(CONCENTRACTOR_CONFIG_PATH)
	check(err)
	var concentractor concentractor
	err = json.Unmarshal([]byte(currentData), &concentractor)

	// fmt.Println("concentractorCurrentStatus", concentractor)

	if err != nil {
		ErrorLogger.Println("JsonToMap err: ", err)
		blePrint([]byte(fmt.Sprint("JsonToMap err")))
	}

	var bleMessageMap map[string]interface{}
	err = json.Unmarshal([]byte(bleMessage), &bleMessageMap)
	if err != nil {
		ErrorLogger.Println("JsonToMap err: ", err)
		blePrint([]byte(fmt.Sprint("JsonToMap err")))
	}

	for k, v := range bleMessageMap {
		switch k {
		case "status":
			if fmt.Sprint(v) == "read" {
				statusString, err := json.Marshal(concentractor)
				if err != nil {
					// fmt.Println(err)
					InfoLogger.Println("BLE COMMAND RESULT:" + k + " read fail")
					return
				}
				InfoLogger.Println("BLE COMMAND RESULT:" + k + " read done")
				// InfoLogger.Println("send Lorawan concentrator status to BLE: " + string(statusString))
				blePrint(statusString)
				continue
			} else {
				InfoLogger.Println("BLE COMMAND RESULT:" + k + " got illegal command")
				blePrint([]byte("illegal command"))
			}
		case "Ssid":
			concentractor.Ssid = fmt.Sprint(v)
			err = setWifi(concentractor.Ssid, concentractor.Password)
		case "Password":
			concentractor.Password = fmt.Sprint(v)
			err = setWifi(concentractor.Ssid, concentractor.Password)
		case "Network_server_address":
			concentractor.Network_server_address = fmt.Sprint(v)
			err = setNetworkServerAddress(concentractor.Network_server_address)
		case "Freq_plan":
			concentractor.Freq_plan = fmt.Sprint(v)
			err = setFrequencyPlan(concentractor.Freq_plan)

		default:
			// fmt.Println("error")
			ErrorLogger.Println("BLE COMMAND IllEGAL")
			blePrint([]byte("command illegal"))
		}
		InfoLogger.Println("BLE COMMAND RESULT:" + k + " setting " + check(err))
		blePrint([]byte(k + " setting " + check(err)))
	}
	statusString, err := json.Marshal(concentractor)
	if err != nil {
		ErrorLogger.Println(err)
		blePrint([]byte(fmt.Sprint(err)))
		return
	}
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
		ErrorLogger.Println("failed to " + action + ": " + err.Error())
		panic("failed to " + action + ": " + err.Error())
	}
}

func blePrint(sendbuf []byte) {
	// Send the sendbuf after breaking it up in pieces.
	for len(sendbuf) != 0 {
		// Chop off up to 20 bytes from the sendbuf.
		partlen := 160
		if len(sendbuf) < 160 {
			partlen = len(sendbuf)
		}
		part := sendbuf[:partlen]
		sendbuf = sendbuf[partlen:]
		// This also sends a notification.
		_, err := txChar.Write(part)
		must("send notification", err)
	}
}

var (
	WarningLogger *log.Logger
	InfoLogger    *log.Logger
	ErrorLogger   *log.Logger
)

func init() {
	file, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	InfoLogger = log.New(file, "INFO: ", log.LstdFlags|log.Lshortfile)
	WarningLogger = log.New(file, "WARNING: ", log.LstdFlags|log.Lshortfile)
	ErrorLogger = log.New(file, "ERROR: ", log.LstdFlags|log.Lshortfile)
}

func main() {
	InfoLogger.Println("starting")
	set_BLE_NAME()
	adapter := bluetooth.DefaultAdapter
	must("enable BLE stack", adapter.Enable())
	adv := adapter.DefaultAdvertisement()
	must("config adv", adv.Configure(bluetooth.AdvertisementOptions{
		LocalName:    ble_name,
		ServiceUUIDs: []bluetooth.UUID{serviceUUID},
	}))
	must("start adv", adv.Start())

	must("add service", adapter.AddService(&bluetooth.Service{
		UUID: serviceUUID,
		Characteristics: []bluetooth.CharacteristicConfig{
			{
				Handle: &rxChar,
				UUID:   rxUUID,
				Flags:  bluetooth.CharacteristicWritePermission | bluetooth.CharacteristicWriteWithoutResponsePermission,
				WriteEvent: func(client bluetooth.Connection, offset int, value []byte) {
					txChar.Write(value)
					var line []byte
					for _, c := range value {
						// rawterm.Putchar(c)
						line = append(line, c)
						if c == '}' {
							InfoLogger.Println("Get BLE message")
							// InfoLogger.Println("Get BLE message:" + string(line))
							bleCommandExecute(string(line))
						}
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
	for {
	}
}


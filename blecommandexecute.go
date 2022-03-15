package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
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

const WIFI_CONFIG_PATH string = "./wpa_supplicant.conf"
const CONCENTRACTOR_CONFIG_PATH string = "./siliq_lorawan_concentractor_conf.json"
const LORA_GLOBAL_CONFIG_PATH string = "./global_conf.json"

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
	targetFreqPlanData, err := ioutil.ReadFile("./lora/lorasdk/global_conf_" + freq_plan + ".json")
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

func main() {
	bleCommandExecute(`
	{
		"ssid": "livingroom",
		"password": "12345678",
		"network_server_address": "192.168.1.104",
		"freq_plan": "EU868"
	}
	`)
}

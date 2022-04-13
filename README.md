# Lorawan Concentrator
A Simple LoRaWAN contractor To Connect LoRaWAN Network Sever

## Function

1. RasbainOS 
0. SPI interface connect to Semtech Sx130x gateway chip base concentrator.
0. BLE interface to connect with phone app for setting WiFi SSID and WiFi Password and also could set LoRaWAN network server address, LoRaWAN Frequency Plan.

## Minimum system requested

This project should be install at RPi0,RPi3/4 version with minimum SD card Flash size: 4GB

## BLE Commands

* BLE command could be set by UTF string. Command list as below

1. Read Status -> {"status": "read"}
0. Set SSID and password -> {"Ssid":"your_ssid","Password":"your_password"}
0. Set LoRaWAN network sever address and Frequency Plan-> {"Network_server_address": "192.168.1.105","Freq_plan": "AS923_2"}

* BLE Command feedback

1. {"status": "read"}-> return

> {"Ssid":"your_ssid","Password":"your_password","Network_server_address":"your_server_address","Freq_plan":"AS923_2"}

2. {"Ssid": "your wifi ssid","Password": "your wifi password"}-> return

> "Password setting ok"

3. {"Network_server_address": "192.168.1.105"}-> return

> Network_server_address setting ok

4. {"Freq_plan": "AS923_2"}-> return

> Freq_plan setting ok

## Implement

1. Read BLE strings
0. Parse BLE strings
0. Setting ssid / password 
0. Setting networksever address
0. Setting frequency plan

## Install

1. install go
2. Install BLE
3. Get code from github
4. Make it autorun as service
5. Make wpa_supplicant.conf as any user could read and write
6. make service and check

### Install go

```bash
# 下載最新版的Go軟體並安裝(RPi 4)
sudo apt-get install golang
# 設定go路徑
nano ~/.profile
```

```text
export GOROOT=/usr/lib/go
export GOPATH=$HOME/go
export PATH=$GOPATH/bin:$GOROOT/bin:$PATH
```

更新 profile

```bash
source ~/.profile
```

### Install BLE

* install bluez

```bash
sudo apt update
sudo apt install bluez
```

* Git Clone Go Bluetooth

```bash
git clone https://github.com/tinygo-org/bluetooth.git
```

### get code from github

```bash
# mkdir -p $GOPATH/src/github.com/user
mkdir -p $GOPATH/src/github.com/livinghuang/lorawan_concentrator
cd $GOPATH/src/github.com/livinghuang/lorawan_concentrator
git clone https://github.com/livinghuang/lorawan_concentrator.git
git mod init
git mod tidy
git build blecommandexecute.go
```

### make it autorun as service

```bash
sudo nano /etc/systemd/system/lorawan_concentrator.service
```

```text
[Unit]
Description=BLE Command Execution Server
ConditionPathExists=/home/pi/go/src/github.com/livinghuang/lorawan_concentrator
After=rc-local.service
Requires=bluetooth.service

[Service]
Type=simple

Restart=on-failure
RestartSec=10

WorkingDirectory=/home/pi/go/src/github.com/livinghuang/lorawan_concentrator
ExecStart=/home/pi/go/src/github.com/livinghuang/lorawan_concentrator/blecommandexecute

[Install]
WantedBy=multi-user.target
```

### make wpa_supplicant.conf as any user could read and write

```bash
sudo chmod 755 /etc/systemd/system/lorawan_concentrator.service
```

### make service and check

```bash
sudo systemctl enable lorawan_concentrator.service
sudo systemctl start lorawan_concentrator
sudo systemctl status lorawan_concentrator
sudo journalctl -f -u lorawan_concentrator
```

## check

```bash
nano /etc/wpa_supplicant/wpa_supplicant.conf
nano /home/pi/go/src/github.com/livinghuang/lorawan_concentrator/siliq_lorawan_concentractor_conf.json
nano /home/pi/lora/packet_forwarder/lora_pkt_fwd/global_conf.json
```

* BLE READ SETTING

``` ble command
{
    "status": "read"
}
```

* BLE SET PARAMETER

``` ble command
{
    "Ssid": "your ssid",
    "Password": "your wifi password",
    "Network_server_address": "192.168.1.105",
    "Freq_plan": "AS923_2"
}
```

* BLE SET PARAMETER

``` ble command
{
    "Ssid": "your wifi ssid",
    "Password": "your wifi password",
}
```

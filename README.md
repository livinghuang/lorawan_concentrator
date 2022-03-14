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

1. Read Status -> [ status :read]
0. Set SSID and password -> [ ssid: your_ssid, password:your_password ]
0. Set LoRaWAN network sever address and Frequency Plan-> [ network_server_address: 192.168.1.1, frequency_plan: as923_2]


* BLE Command feedback

1. [status:read]-> return "{Ethernet :ok/ng ,ssid: your_ssid,freqency_plan:as923_2,network_address:192.169.0.1"}
0.  [ ssid: your_ssid, password:your_password ]-> "{settings: done, ssid: your_ssid, password:your_password}"
0. [network_server_address: 192.168.1.1, frequency_plan: as923_2]-> "{setting :done, network_server_address: 192.168.1.1, frequency_plan: as923_2}"

## Implement 

1. Read BLE strings
0. Parse BLE strings
0. Setting ssid / password 
0. Setting networksever address
0. Setting frequency plan

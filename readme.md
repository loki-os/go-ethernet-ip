# go-ethernet-ip

[![standard-readme compliant](https://img.shields.io/badge/readme%20style-standard-brightgreen.svg)](https://github.com/RichardLitt/standard-readme)
[![](https://img.shields.io/github/go-mod/go-version/loki-os/go-ethernet-ip)]()
[![](https://img.shields.io/github/license/loki-os/go-ethernet-ip)]()

A complete golang implementation of Ethernet/ip protocol.

This repository contains:

1. A implementation of ethernet/ip protocol.
2. A lightweight message router.
3. Session management.
4. A lightweight api interface makes you don't need to focus binary steam.
5. Examples of go-ethernet-ip.

## Table of Contents

- [Background](#Background)
- [Install](#Install)
- [Usage](#Usage)
	- [Find all LAN devices](#Find-all-LAN-devices)
	- [Device connection config](#Device-connection-config)
	- [New device](#New-device)
	- [List identity](#List-identity)
	- [List interface](#List-interface)
	- [List services](#List-services)
	- [Send RRData](#List-services)
	- [Send UnitData](#List-services)
	- [Disconnect](#Disconnect)
- [Maintainers](#Maintainers)
- [Contributing](#Contributing)
- [License](#License)

## Background

`go-ethernet-ip` started with the my own project goplc which used for communication with rockwell control/compactLogix PLCs.

I separate common industrial protocol from ethernet/ip.

If your communication protocol is common industrial protocol, you should move to [go-cip](https://github.com/loki-os/go-cip) which base on this repository.

## Install

This project uses [golang](https://golang.org/). Go check them out if you don't have them locally installed.

```sh
$ go get github.com/loki-os/go-ethernet-ip
```

Also go modules is supported.

## Usage

I use some cip cases for demonstration.

Before you use these cases.You should block your main thread.Because all logic run in go routine.

```go
func block(){
	some_case()

	// you'd better find other way to do this. Sleep is not recommended.
	time.Sleep(time.Second * 10)
}
``` 

### Find all LAN devices

Before we start to communication with other device, we need to find them via lan. If you have clear ip skip this step.

```go
devices, err := GetLanDevices(time.Second)
```

### New device

```go
device, err := NewDevice("192.168.0.100")
```

### Device connection config

Device Connect() is optional.You can skip this step to use default config.

Connect() should be call before other function, otherwise it will fail.

```go
cfg := DefaultConfig()
cfg.TCPTimeout = time.Second * 5

device.Connect(cfg)
```

### List identity

```go
identity, err := device.ListIdentity()
```

### List interface

```go
interface, err := device.ListInterface()
```

### List services

```go
services, err := device.ListServices()
```

### Send RRData

```go
device.SendRRData(cpf *commonPacketFormat, timeout typedef.Uint)
```

### Send UnitData

```go
device.SendUnitData(cpf *commonPacketFormat, timeout typedef.Uint)
```

### Disconnect

```go
device.Disconnect()
```

## Maintainers

[@末日上投](https://github.com/MiguelValentine).

## Contributing

Feel free to dive in! [Open an issue](https://github.com/loki-os/go-ethernet-ip/issues/new) or submit PRs.

## License

[MIT](LICENSE) © 末日上投
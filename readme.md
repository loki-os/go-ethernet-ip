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
	- [List Identity](#List-Identity)
	- [List Interface](#List-Interface)
	- [List services](#List-services)
- [Maintainers](#Maintainers)
- [Contributing](#Contributing)
- [License](#License)

## Background

`go-ethernet-ip` started with the my own project goplc which used for communication with rockwell control/compactLogix PLCs.

I separate common industrial protocol from ethernet/ip.

If your communication protocol is common industrial protocol, you should move to [go-cip](https://github.com/loki-os/go-cip) which base on this repository.
s
## Install

This project uses [golang](https://golang.org/). Go check them out if you don't have them locally installed.

```sh
$ go get github.com/loki-os/go-ethernet-ip
```

Also go modules is supported.

## Usage

I use some cip cases for demonstration.

Before you use these cases.You should block your main thread.

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
func ListAllLanDevices() {
	udp, e := NewUDPWithAutoScan(nil)
	if e != nil {
		log.Println(e)
		return
	}

	e1 := udp.Connect()
	if e1 != nil {
		log.Println(e1)
		return
	}
	defer udp.Close()

	udp.ListIdentity()

	// you should sleep for result because udp use broadcast message
	time.Sleep(time.Second)

	b, _ := json.MarshalIndent(udp.Devices, "", "\t")
	log.Println(string(b))
}
```

### List Identity

```go
func ListIdentity() {
	// tcp
	tcp, e := NewTcpWithAddress("192.168.0.100", nil)
	if e != nil {
		log.Println(e)
		return
	}

	e1 := tcp.Connect()
	if e1 != nil {
		log.Println(e1)
		return
	}

	tcp.ListIdentity(func(data interface{}, err error) {
		if err != nil {
			log.Println(err)
			return
		}

		b, _ := json.MarshalIndent(data, "", "\t")
		log.Println(string(b))
	})

	// udp
	udp, e2 := NewUDPWithAddress("192.168.0.100", nil)
	if e2 != nil {
		log.Println(e)
		return
	}

	e3 := udp.Connect()
	if e3 != nil {
		log.Println(e1)
		return
	}

	udp.ListIdentity()
	time.Sleep(time.Second)

	b, _ := json.MarshalIndent(udp.Devices, "", "\t")
	log.Println(string(b))
}
```

### List Interface

```go
func ListInterface() {
	// tcp
	tcp, e := NewTcpWithAddress("10.211.55.7", nil)
	if e != nil {
		log.Println(e)
		return
	}

	e1 := tcp.Connect()
	if e1 != nil {
		log.Println(e1)
		return
	}

	tcp.ListInterface(func(data interface{}, err error) {
		if err != nil {
			log.Println(err)
			return
		}

		log.Printf("%+v\n", data)
	})
	
	// udp
	// we supported use udp to list interface but not recommended.
}
```

### List services

```go
func ListServices() {
	// tcp
	tcp, e := NewTcpWithAddress("10.211.55.7", nil)
	if e != nil {
		log.Println(e)
		return
	}

	e1 := tcp.Connect()
	if e1 != nil {
		log.Println(e1)
		return
	}

	tcp.ListServices(func(data interface{}, err error) {
		if err != nil {
			log.Println(err)
			return
		}

		b, _ := json.MarshalIndent(data, "", "\t")
		log.Println(string(b), string(data.(*ListServices).Items[0].Name))
	})

	// udp
	// we supported use udp to list services but not recommended.
}
```

## Maintainers

[@末日上投](https://github.com/MiguelValentine).

## Contributing

Feel free to dive in! [Open an issue](https://github.com/loki-os/go-ethernet-ip/issues/new) or submit PRs.

## License

[MIT](LICENSE) © 末日上投
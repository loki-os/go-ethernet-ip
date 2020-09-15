# go-ethernet-ip

[![standard-readme compliant](https://img.shields.io/badge/readme%20style-standard-brightgreen.svg)](https://github.com/RichardLitt/standard-readme)
[![](https://img.shields.io/github/go-mod/go-version/loki-os/go-ethernet-ip)]()
[![](https://img.shields.io/github/license/loki-os/go-ethernet-ip)]()

A complete golang implementation of Ethernet/ip protocol.

This repository contains:

1. A implementation of ethernet/ip protocol.
2. A lightweight message router.
3. A lightweight api interface makes you don't need to focus binary steam.
4. Examples of go-ethernet-ip.

## Table of Contents

- [Background](#Background)
- [Install](#Install)
- [Usage](#Usage)
	- [Scan](#Scan)
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

I used some cip cases for demonstration.

### Scan

```go

```

## Maintainers

[@末日上投](https://github.com/MiguelValentine).
[@Skyblue](https://github.com/skyblue).

## Contributing

Feel free to dive in! [Open an issue](https://github.com/loki-os/go-ethernet-ip/issues/new) or submit PRs.

## License

[MIT](LICENSE) © 末日上投

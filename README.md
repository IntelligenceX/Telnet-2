# Telnet-2

Telnet-2 has 2 features which the original telnet does not provide:

1. By default, it connects via Tor (expecting a local proxy to run at 127.0.0.1:9050)
2. You can optionally define the connection timeout (default is 10s).

Both of those features can be important when investigating high-profile targets such as servers in North Korea. Using Tor helps to conceal your identity and helps to bypass IP blocking.

## Use

```
Telnet-2 [Host] [Port] [Optional Timeout]
```

Timeout is in seconds (the default is 10 seconds).

## Compile

Download [Go](https://golang.org/dl/) and then compile the project using the command `go build`. You can compile it on Linux, Mac and Windows.

## Copyright

This is free and unencumbered software released into the public domain.
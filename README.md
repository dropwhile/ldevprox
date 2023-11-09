# ldevprox

## Install

```
go install github.com/dropwhile/ldevprox@latest
```

## Usage

```
Usage: ldevprox

A local ssl proxy for development

Flags:
  -h, --help                       Show context-sensitive help.
  -l, --listen="127.0.0.1:8080"    listen address:port
      --tls-cert="server.crt"      tls cert file path
      --tls-key="server.key"       tls key file path
      --upstream="http://127.0.0.1:8000"
                                   upstream url
  -v, --verbose                    enable verbose logging
  -V, --version                    Print version information and quit
```

## Example

```
ldevprox \
    -l 127.0.0.1:8080 \
    --tls-cert server.crt \
    --tls-key server.key \
    --upstream "http://127.0.0.1:8000"
```

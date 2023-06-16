# export local service to public 

- [export local service to public](#export-local-service-to-public)
- [`install`](#install)
  - [`use go install`](#use-go-install)
  - [`download release`](#download-release)
- [start](#start)
  - [server](#server)
  - [client](#client)

# `install`

## `use go install`

```
go install github.com/0x1042/localtunnel@latest
```

## `download release`

[download](https://github.com/0x1042/localtunnel/releases)

# start
## server 

```shell
localtunnel server --help
NAME:
   localtunnel server - start tunnel server

USAGE:
   localtunnel server [command options] [arguments...]

OPTIONS:
   --port value, -p value    tunnel listen port (default: 7853)
   --secret value, -s value  secret [$LT_TOKEN]
   --help, -h                show help
```

## client

```shell
./localtunnel client --help
NAME:
   localtunnel client - start tunnel client

USAGE:
   localtunnel client [command options] [arguments...]

OPTIONS:
   --tunnel value, -t value  tunnel server addr
   --secret value, -s value  secret [$LT_TOKEN]
   --local value, -l value   local port (default: 0)
   --remote value, -r value  remote port (default: random)
   --help, -h                show help
```
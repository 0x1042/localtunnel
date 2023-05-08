# 本地端口映射

# 运行 

## server 

```shell
lt server --help
NAME:
   lt server - start tunnel server

USAGE:
   lt server [command options] [arguments...]

OPTIONS:
   --port value, -p value    tunnel listen port (default: 7853)
   --secret value, -s value  secret [$LT_TOKEN]
   --help, -h                show help
```

## client

```shell
lt client --help
NAME:
   lt client - start tunnel client

USAGE:
   lt client [command options] [arguments...]

OPTIONS:
   --tunnel value, -t value  tunnel server addr
   --secret value, -s value  secret [$LT_TOKEN]
   --local value, -l value   local port (default: 0)
   --remote value, -r value  remote port (default: random)
   --help, -h                show help
```
# # hikarian
hikarian just make you crossover the wall:)

#tunnel
##install
```bash
make tunnel
```

##client
```bash
bin/tunnel --client 127.0.0.1:1081 --server 127.0.0.1:1082 --mode encrypt --algo rc4
```
##server
```bash
bin/tunnel --client 127.0.0.1:1082 --server 127.0.0.1:1080 --mode decrypt --algo rc4
```

#socks5
```bash
make socks5
```
##run
```bash
bin/socks5
```
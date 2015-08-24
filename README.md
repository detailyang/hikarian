#hikarian
hikarian just make you crossover the wall:)

##install
```bash
make tunnel
make socks5
or
make
```
## run
* tunnel client

```bash
#there is a bug, if we use rc4 to encrypt https, it will cause "Encrypted Alert"
bin/tunnel --client 127.0.0.1:1081 --server 127.0.0.1:1082 --mode encrypt --algo rc4 or aes
```

* tunnel server

```bash
bin/tunnel --client 127.0.0.1:1082 --server 127.0.0.1:1080 --mode decrypt --algo rc4 or aes
```

* socks5 server

```bash
bin/socks5
```


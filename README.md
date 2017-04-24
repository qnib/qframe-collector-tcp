# qframe-collector-tcp
TCP collector for the qframe framework.

## main.go

The example script will instantiate the collector and wait for a message send to it.


```bash
go run main.go
2017/04/21 12:37:29 [II] Dispatch broadcast for Data and Tick
2017/04/21 12:37:29 [  INFO] test >> Listening on 127.0.0.1:11001
```

Once send...

```bash
$ echo "Test-$(date +%s)" | nc -w1  localhost 11001
```

... the message will be displayed and the script exits:

```
#### Received (remote:127.0.0.1:60846): Test-1492771635
```


## Deveolpment

```bash
$ docker run -ti --name qframe-collector-tcp --rm -e SKIP_ENTRYPOINTS=1 \
           -v ${GOPATH}/src/github.com/qnib/qframe-collector-tcp:/usr/local/src/github.com/qnib/qframe-collector-tcp \
           -v ${GOPATH}/src/github.com/qnib/qframe-collector-docker-events:/usr/local/src/github.com/qnib/qframe-collector-docker-events \
           -v ${GOPATH}/src/github.com/qnib/qframe-types:/usr/local/src/github.com/qnib/qframe-types \
           -v /var/run/docker.sock:/var/run/docker.sock \
           -w /usr/local/src/github.com/qnib/qframe-collector-tcp \
            qnib/uplain-golang bash
# govendor update github.com/qnib/qframe-collector-tcp/lib github.com/qnib/qframe-collector-docker-events/lib github.com/qnib/qframe-types
```
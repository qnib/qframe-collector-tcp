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
#### Received: Test-1492771065
```
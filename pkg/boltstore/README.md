# Boltstore

Boltstore implements the [store interface](https://github.com/eclipse/paho.mqtt.golang/blob/master/store.go) of  [paho.mqtt.golang](https://github.com/eclipse/paho.mqtt.golang).

```go
type Store interface {
    Open()
    Put(key string, message packets.ControlPacket)
    Get(key string) packets.ControlPacket
    All() []string
    Del(key string)
    Close()
    Reset()
}
```

---

## Run Test

`go test -v`

## Environment Variables

> Why environment variable? Refer to the [12 factors](http://12factor.net/)

- `BOLTSTORE_CONF`: config file path

Don't worry if you do not provide a config file via environment variable, boltstore will use its [default values](config.go).

## Config File Example

```json
{
    "BoltStore": {
        "BucketName": "DB",
        "OpenTimeout": 2,
        "Logger": {
            "Level": 3,
            "JSON": true,
            "ToFile": true,
            "Filename": "boltstore.test.log"
        }
    }
}
```

### Logger Level

- Panic : 0
- Fatal : 1
- Error : 2
- Warn  : 3
- Info  : 4
- Debug : 5

---

## Dependent libraries

- [bolt](github.com/boltdb/bolt) -  a pure Go key/value store, [MIT License](https://github.com/boltdb/bolt/blob/master/LICENSE).
- [paho.mqtt.golang](github.com/eclipse/paho.mqtt.golang) - Eclipse Paho MQTT Go client, [EPL 1.0 License](https://github.com/eclipse/paho.mqtt.golang/blob/master/LICENSE)
- [logrus](github.com/sirupsen/logrus) - a structured logger for Go, [MIT License](https://github.com/sirupsen/logrus/blob/master/LICENSE)
- [multiconfig](https://github.com/koding/multiconfig) - load configuration from multiple sources in Go, [MIT License](https://github.com/koding/multiconfig/blob/master/LICENSE)

## Licensing

The package is licensed under the Apache License, Version 2.0. See
[LICENSE](https://github.com/taka-wang/apollo-edge/blob/master/LICENSE) for the full
license text.

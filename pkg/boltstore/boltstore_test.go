package boltstore

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/eclipse/paho.mqtt.golang/packets"
	"github.com/takawang/sugar"
)

const (
	inboundPrefix  = "i."
	outboundPrefix = "o."
)

// Return a string of the form "i.[id]"
func inboundKeyFromMID(id uint16) string {
	return fmt.Sprintf("%s%d", inboundPrefix, id)
}

// Return a string of the form "o.[id]"
func outboundKeyFromMID(id uint16) string {
	return fmt.Sprintf("%s%d", outboundPrefix, id)
}

// remove temporary bolt databases
func housekeeping() {
	files, err := filepath.Glob("*.db")
	if err != nil {
		panic(err)
	}

	for _, f := range files {
		if err := os.Remove(f); err != nil {
			panic(err)
		}
	}
	fmt.Println("Finished housekeeping")
}

// new control packet
func newPacket(topic string, mid uint16) *packets.PublishPacket {
	pub := packets.NewControlPacket(packets.Publish).(*packets.PublishPacket)
	pub.Qos = 1
	pub.TopicName = topic
	pub.Payload = []byte{0xBE, 0xEF, 0xED}
	pub.MessageID = mid
	return pub
}

// Test Cases =========================
func TestBoltStore(t *testing.T) {
	s := sugar.New(t)

	s.Assert("Create and open a new store 'with' a given path name", func(logf sugar.Log) bool {
		given := "my.db"
		store := NewBoltStore(given)
		store.Open()
		defer store.db.Close()
		logf("store path is: %s\n", store.path)
		logf("store is opened at: %v\n", store.opened)
		if store.path == given {
			logf("store path is correct!")
			return true
		}
		logf("store path is wrong!")
		return false
	})

	s.Assert("Create and open a new store 'without' a path name", func(logf sugar.Log) bool {
		store := NewBoltStore("")
		store.Open()
		defer store.db.Close()
		logf("store path is: %s\n", store.path)
		logf("store is opened: %v\n", store.opened)
		if store.opened {
			return true
		}
		return false
	})

	s.Assert("Close a bolt store", func(logf sugar.Log) bool {
		store := NewBoltStore("")
		store.Open()
		store.Close()
		logf("store is opened: %v\n", store.opened)
		if !store.opened {
			return true
		}
		return false
	})

	s.Assert("Put two packets with topic 'hello' and 'hello1' into the boltstore", func(logf sugar.Log) bool {
		store := NewBoltStore("my.db")
		store.Open()
		defer store.Close()

		m1 := newPacket("hello", 91)
		key1 := inboundKeyFromMID(m1.MessageID)
		store.Put(key1, m1)
		logf("put: %v\n", m1)

		m2 := newPacket("hello2", 92)
		key2 := outboundKeyFromMID(m2.MessageID)
		store.Put(key2, m2)
		logf("put: %v\n", m2)

		return true
	})

	s.Assert("Reset a bolt store", func(logf sugar.Log) bool {
		store := NewBoltStore("my.db")
		store.Open()
		defer store.Close()

		store.Reset()
		if m := store.All(); len(m) > 0 {
			return false
		}
		return true
	})

	s.Assert("Put two packets with topic 'hello' and 'hello1' into the boltstore again", func(logf sugar.Log) bool {
		store := NewBoltStore("my.db")
		store.Open()
		defer store.Close()

		m1 := newPacket("hello", 91)
		key1 := inboundKeyFromMID(m1.MessageID)
		store.Put(key1, m1)
		logf("put: %v\n", m1)

		m2 := newPacket("hello2", 92)
		key2 := outboundKeyFromMID(m2.MessageID)
		store.Put(key2, m2)
		logf("put: %v\n", m2)

		return true
	})

	s.Assert("Get a packet with key 'i.91' from the boltstore", func(logf sugar.Log) bool {
		store := NewBoltStore("my.db")
		store.Open()
		defer store.Close()

		m := store.Get("i.91")
		logf("get: %v\n", m)
		if m != nil {
			return true
		}
		return false
	})

	s.Assert("Get all keys from the boltstore", func(logf sugar.Log) bool {
		store := NewBoltStore("my.db")
		store.Open()
		defer store.Close()

		m := store.All()
		logf("get: %v\n", m)
		if m != nil {
			return true
		}
		return false
	})

	s.Assert("Delete a packet with key 'o.92' from the boltstore", func(logf sugar.Log) bool {
		store := NewBoltStore("my.db")
		store.Open()
		defer store.Close()

		store.Del("o.92")
		logf("Delete key 'o.92'")
		return true
	})

	s.Assert("Get all keys from the boltstore", func(logf sugar.Log) bool {
		store := NewBoltStore("my.db")
		store.Open()
		defer store.Close()

		m := store.All()
		logf("get: %v\n", m)
		if m != nil {
			return true
		}
		return false
	})

	// remove all temp db
	housekeeping()

	if s.IsFailed() {
		fmt.Println("the tests failed :/")
	}
}

func TestIntegrateWithPaho(t *testing.T) {
	s := sugar.New(t)

	s.Assert("MQTT integration test", func(logf sugar.Log) bool {
		// Enable MQTT logger
		MQTT.DEBUG = log.New(os.Stdout, "", 0)
		MQTT.ERROR = log.New(os.Stdout, "", 0)
		// Enable debug level logger
		defaultConfig.BoltStore.Logger.Level = 5
		setLogger()

		store := NewBoltStore("mqtt.db")
		//store := MQTT.NewMemoryStore()

		msgChan := make(chan MQTT.Message)

		opts := MQTT.NewClientOptions()
		opts.AddBroker("tcp://iot.eclipse.org:1883")
		opts.SetStore(store)

		var callback MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
			logf("Got [TOPIC: %s, MSG: %s]", msg.Topic(), msg.Payload())
			msgChan <- msg
		}

		c := MQTT.NewClient(opts)
		defer c.Disconnect(250)
		if token := c.Connect(); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}

		if token := c.Subscribe("/boltstore/sample", 1, callback); token.Wait() && token.Error() != nil {
			logf("subscribe error: %v", token.Error())
			return false
		}

		go func() {
			for i := 0; i < 5; i++ {
				token := c.Publish("/boltstore/sample", 1, false, fmt.Sprintf("this is msg #%d!", i))
				token.Wait()
			}
		}()

		var recvCount = 0
		for {
			select {
			case <-msgChan:
				recvCount = recvCount + 1
				if recvCount > 4 {
					logf("Got all messages")
					return true
				}
			case <-time.After(10 * time.Second):
				logf("SELECT Timeout")
				return false
			}
		}
	})

	// remove all temp db
	housekeeping()

	if s.IsFailed() {
		fmt.Println("the tests failed :/")
	}
}

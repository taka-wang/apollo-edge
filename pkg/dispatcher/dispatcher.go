package dispatcher

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	log "github.com/sirupsen/logrus"
)

var (
	keepAlive        = 2
	pingTimeout      = 1
	qos1        byte = 1
	qos2        byte = 2
)

type routeType struct {
	client MQTT.Client
	dest   string
}
type clientType struct {
	subTopic   []string
	subHandler []MQTT.MessageHandler
	// pubChan message to be published
	pubChan chan Msg
	client  MQTT.Client
}

// Dispatcher structure
type Dispatcher struct {
	id       string
	routeMap map[string]routeType
	edgeHub  clientType
	cloudHub clientType
}

/*
data :=
`
{
	"routes": {
		"sensorToFilter": "FROM /messages/modules/tempSensor/outputs/temperatureOutput INTO BrokeredEndpoint(\"/modules/filtermodule/inputs/input1\")",
		"filterToIoTHub": "FROM /messages/modules/filtermodule/outputs/output1 INTO $upstream", // stamp
		"xxx": "FROM /devices/hello/messages/events INTO $upstream", // stamp
		"yyy": "FROM /devices/world/messages/events INTO BrokeredEndpoint(\"/modules/filtermodule/inputs/input2\")",
		"zzz": "FROM /devices/world/messages/devicebound ===> world
	}
}
`
*/

// RouterType routing rule JSON type
type RouterType struct {
	Routes map[string]string `json:"routes"`
}

// stamped append field to arbitrary JSON string
func stamped(b []byte) {
	inputJSON := `{"environment": "production", "runbook":"http://url","message":"there is a problem"}`
	out := map[string]interface{}{}
	if err := json.Unmarshal([]byte(inputJSON), &out); err != nil {
		// do something
		return
	}

	out["name"] = "taka"
	out["command"] = "wang"
	out["status"] = 2

	outputJSON, _ := json.Marshal(out)

	fmt.Printf("%s", outputJSON)
}

// loadRouteFile load the routing JSON file
func loadRouteFile(filename string) map[string]routeType {

	// read routing file
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		log.WithField("error", err).Fatal(ErrLoadRouteFile)
	}
	log.WithField("route file", string(file)).Infoln("Read routing file")

	// load the routing rules from the JSON file to the route map.
	routeMap := make(map[string]routeType)

	return routeMap
}

// Routes routing JSON
type Routes struct {
	Route struct {
		SensorToFilter string `json:"sensorToFilter"`
		FilterToIoTHub string `json:"filterToIoTHub"`
	} `json:"routes"`
}

// NewDispatcher returns a pointer to a new instance of Dispatcher
func NewDispatcher(deviceID string, edgeHubConnStr string, cloudHubConnStr string) *Dispatcher {

	dp := &Dispatcher{
		id:       deviceID,
		routeMap: loadRouteFile(""),
		edgeHub: clientType{
			subTopic: []string{
				"/messages/modules/+/outputs/#",
				"/devices/+/messages/events",
			},
			pubChan: make(chan Msg),
			client:  nil,
		},
		cloudHub: clientType{
			subTopic: []string{
				"/devices/" + deviceID + "/messages/devicebound",
				// who are in my group?
			},
			pubChan: make(chan Msg),
			client:  nil,
		},
	}
	return dp
}

/*
// Start start routing loop
func (d *Dispatcher) Start() {
	// start subscribe
	go d.listenEdgeHub(qos1)
	go d.listenCloudHub(qos2)

	for {
		select {
		// dispatch messages to the edge MQTT broker
		case ingress := <-d.msgToEdgeHub:
			go d.publish(d.edgeHubClient, ingress)
		// dispatch messages to the cloud MQTT broker
		case egress := <-d.msgToCloudHub:
			go d.publish(d.cloudHubClient, egress)
		// graceful shutdown
		case <-d.stopSignal:
			fmt.Println("ready to shut down")
			// TODO: do some cleanup here
			return
		}
	}
}

// publish send a message to a MQTT broker
func (d *Dispatcher) publish(client MQTT.Client, msg Msg) {
	token := client.Publish(msg.Topic(), qos1, false, msg.Payload())
	if token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
}

// listenEdgeHub subscribe messages from the edge hub
func (d *Dispatcher) listenEdgeHub(qos byte) {
	var topic = "/msg/module/+/outputs/#"
	if token := d.edgeHubClient.Subscribe(topic, qos, nil); token.Wait() && token.Error() != nil {
		log.WithField("error", token.Error()).Errorln("listenEdgeHub error")
	}
}

// listenCloudHub subscribe messages from the cloud hub
func (d *Dispatcher) listenCloudHub(qos byte) {
	var topic = "/devices/" + d.id + "/messages/devicebound"
	if token := d.cloudHubClient.Subscribe(topic, qos, nil); token.Wait() && token.Error() != nil {
		log.WithField("error", token.Error()).Errorln("listenCloudHub error")
	}
}
*/

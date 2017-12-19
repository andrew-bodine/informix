package wiot

import (
    "encoding/json"
    "fmt"
    "sync"

    MQTT "github.com/eclipse/paho.mqtt.golang"
)

// NOTE: Watson IoT Platform documentation for MQTT clients and gateways here:
// https://console.ng.bluemix.net/docs/services/IoT/gateways/mqtt.html#mqtt

func NewClient(o *Options) *client {
    if o == nil {
        return nil
    }

    return &client{opts:  o}
}

type client struct {
    opts    *Options

    cli     MQTT.Client

    sync.Mutex
}

// Connect creates a new MQTT client and tries to connect to the remote
// MQTT server.
func (c *client) Connect() error {
    c.Lock()
    defer c.Unlock()

    o := MQTT.NewClientOptions()

    o.AddBroker(c.opts.Broker)
    o.SetClientID(c.opts.ClientId)
    o.SetUsername(c.opts.Username)
    o.SetPassword(c.opts.Password)

    o.SetAutoReconnect(false)

    cli := MQTT.NewClient(o)

    if t := cli.Connect(); t.Wait() && t.Error() != nil {
        return t.Error()
    }
    c.cli = cli

    return nil
}

// Publish extracts configuration params from environment variables, and
// publishes the message payload to the appropriate MQTT topic.
func (c *client) Publish(t string, payload map[string]interface{}) error {
    c.Lock()
    if c.cli == nil {
        return MQTT.ErrNotConnected
    }
    c.Unlock()

    data, err := json.Marshal(payload)
    if err != nil {
        return err
    }

    topic := fmt.Sprintf(
        "iot-2/type/%s/id/%s/evt/%s/fmt/json",
        c.opts.DeviceType(),
        c.opts.DeviceId(),
        t,
    )

    if t := c.cli.Publish(topic, 0, false, data); t.Wait() && t.Error() != nil {
        return t.Error()
    }

    return nil
}

package wiot

import (
    "encoding/json"
    "fmt"
    "sync"

    MQTT "github.com/eclipse/paho.mqtt.golang"

    "github.com/andrew-bodine/informix/downstream"
)

// NOTE: Watson IoT Platform documentation for MQTT clients and gateways here:
// https://console.ng.bluemix.net/docs/services/IoT/gateways/mqtt.html#mqtt

func NewClient(o *Options) downstream.Downstreamer {
    if o == nil {
        return nil
    }

    return &Client{opts:  o}
}

// A downstream.Downstreamer implementation.
type Client struct {
    opts    *Options

    cli     MQTT.Client

    sync.Mutex
}

// Exposes the underlying MQTT client.
func (c *Client) MQTTClient() MQTT.Client {
    c.Lock()
    defer c.Unlock()

    return c.cli
}

// Implement the downstream.Downstreamer interface.
func (c *Client) Connect() error {
    c.Lock()
    defer c.Unlock()

    o := MQTT.NewClientOptions()

    o.AddBroker(c.opts.Broker)
    o.SetClientID(c.opts.ClientId)
    o.SetUsername(c.opts.Username)
    o.SetPassword(c.opts.Password)

    o.SetAutoReconnect(false)

    // Create a new MQTT client.
    cli := MQTT.NewClient(o)

    // Try to connect to the remote MQTT server.
    if t := cli.Connect(); t.Wait() && t.Error() != nil {
        return t.Error()
    }
    c.cli = cli

    return nil
}

// Implement the downstream.Downstreamer interface.
func (c *Client) Publish(t string, payload map[string]interface{}) error {
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

    // Publish the message payload to the specified topic.
    if t := c.cli.Publish(topic, 0, false, data); t.Wait() && t.Error() != nil {
        return t.Error()
    }

    return nil
}

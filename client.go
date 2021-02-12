package mqtt

import (
	"context"
	"time"

	paho "github.com/eclipse/paho.mqtt.golang"
	"github.com/loadimpact/k6/lib"
)

// Mqtt is the objet to be used in tests
type Mqtt struct {
}

// Connect create a connection to mqtt
func (*Mqtt) Connect(
	// The list of URL of  MQTT server to connect to
	servers []string,
	// A username to authenticate to the MQTT server
	user,
	// Password to match username
	password string,
	// clean session setting
	cleansess bool,
	// Client id for reader
	clientid string,
	// timeout ms
	timeout uint,

) paho.Client {

	opts := paho.NewClientOptions()
	for i := range servers {
		opts.AddBroker(servers[i])
	}
	opts.SetClientID(clientid)
	opts.SetUsername(user)
	opts.SetPassword(password)
	opts.SetCleanSession(cleansess)
	client := paho.NewClient(opts)
	token := client.Connect()
	if !token.WaitTimeout(time.Duration(timeout) * time.Millisecond) {
		ReportError(token.Error(), "Connect timeout")
		return nil
	}
	if token.Error() != nil {
		ReportError(token.Error(), "Connect failed")
		return nil
	}
	return client
}

// Close the given client
func (*Mqtt) Close(
	ctx context.Context,
	// Mqtt client to be closed
	client paho.Client,
	// timeout ms
	timeout uint,
) error {
	state := lib.GetState(ctx)
	if state == nil {
		ReportError(ErrorNilState, "Subscribe Cannot determine state")
		return ErrorNilState
	}
	client.Disconnect(timeout)
	return nil
}
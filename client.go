package mqtt

import (
	"context"
	"fmt"
	"net/url"

	"github.com/eclipse/paho.golang/autopaho"
	"github.com/eclipse/paho.golang/paho"
	"go.k6.io/k6/js/common"
	"go.k6.io/k6/lib"
)

// Mqtt is the objet to be used in tests
type Mqtt struct {
}

// Connect create a connection to mqtt
func (*Mqtt) Connect(
	ctx context.Context,
	// The list of URL of  MQTT server to connect to
	servers []string,
	// A username to authenticate to the MQTT server
	user,
	// Password to match username
	password string,
	// Topic
	topic string,
	// clean session setting
	cleansess bool,
	// Client id for reader
	clientid string,
	// timeout ms
	timeout uint,
	// path to local cert
	certPath string,

) *autopaho.ConnectionManager {
	state := lib.GetState(ctx)
	if state == nil {
		common.Throw(common.GetRuntime(ctx), ErrorState)
		return nil
	}
	urls := make([]*url.URL, len(servers))
	for i,v := range servers {
		urls[i], _ = url.Parse(v)
	}
	print(urls)
	// var servers_urls := []*url.URL
	opts := autopaho.ClientConfig{
		BrokerUrls: urls,
		KeepAlive: 5,
		ConnectRetryDelay: 5,
		Debug:             paho.NOOPLogger{},
		OnConnectionUp:    func(cm *autopaho.ConnectionManager, connAck *paho.Connack) {
			fmt.Println("mqtt connection up")
			if _, err := cm.Subscribe(context.Background(), &paho.Subscribe{
				Subscriptions: map[string]paho.SubscribeOptions{
					topic: {QoS: 0},
				},
			}); err != nil {
				common.Throw(common.GetRuntime(ctx), ErrorSubscribe)
				fmt.Printf("failed to subscribe (%s). This is likely to mean no messages will be received.", err)
				return
			}
			fmt.Println("mqtt subscription made")
		},
		OnConnectError:    func(err error) {
			common.Throw(common.GetRuntime(ctx), ErrorClient)
			fmt.Printf("error whilst attempting connection: %s\n", err)
		},

		ClientConfig: paho.ClientConfig{
			ClientID:      clientid,
			OnServerDisconnect: func(d *paho.Disconnect) {
				if d.Properties != nil {
					fmt.Printf("server requested disconnect: %s\n", d.Properties.ReasonString)
				} else {
					fmt.Printf("server requested disconnect; reason code: %d\n", d.ReasonCode)
				}
			},
		},
	}
	opts.SetUsernamePassword(user, []byte(password))
	opts.SetWillMessage(fmt.Sprintf("lwt/%s", user),[]byte("{ \"status\": \"offline\", \"exit\": \"unclean\" }"),0,false)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cm, err := autopaho.NewConnection(ctx, opts)
	if err != nil {
		panic(err)
	}
	err = cm.AwaitConnection(ctx)

	if err != nil {
		fmt.Printf("Failed to connect. (%s)\n", err)
		panic(err)
	}
	return cm
}

// Close the given client
func (*Mqtt) Close(
	ctx context.Context,
	// Mqtt client to be closed
	cm *autopaho.ConnectionManager,
	// timeout ms
	timeout uint,
) {
	state := lib.GetState(ctx)
	if state == nil {
		common.Throw(common.GetRuntime(ctx), ErrorState)
		return
	}
	cm.Disconnect(ctx)
	return
}

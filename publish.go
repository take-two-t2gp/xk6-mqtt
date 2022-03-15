package mqtt

import (
	"context"

	"github.com/eclipse/paho.golang/autopaho"
	"github.com/eclipse/paho.golang/paho"
	"go.k6.io/k6/js/common"
	"go.k6.io/k6/lib"
)

// Publish allow to publish one message
func (*Mqtt) Publish(
	ctx context.Context,
	cm *autopaho.ConnectionManager,
	topic string,
	qos int,
	message string,
	retain bool,
	timeout int,
) {
	state := lib.GetState(ctx)
	if state == nil {
		common.Throw(common.GetRuntime(ctx), ErrorState)
		return
	}
	if cm == nil {
		common.Throw(common.GetRuntime(ctx), ErrorClient)
		return
	}

	pr, err := cm.Publish(ctx, &paho.Publish{
		QoS: 0,
		Topic: topic,
		Payload: []byte(message),
	})

	if err != nil {
		common.Throw(common.GetRuntime(ctx), ErrorPublish)
		return
	} else if pr.ReasonCode != 0 && pr.ReasonCode != 16 {
		common.Throw(common.GetRuntime(ctx), ErrorPublish)
	}
	return
}

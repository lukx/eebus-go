package ship

import (
	"fmt"

	"github.com/enbility/eebus-go/ship/model"
)

// Handshake initialization covers the states cmiState...

// CMI_STATE_INIT_START
func (c *ShipConnection) handshakeInit_cmiStateInitStart() {
	switch c.role {
	case ShipRoleClient:
		// CMI_STATE_CLIENT_SEND
		c.setState(cmiStateClientSend)
		if err := c.DataHandler.WriteMessageToDataConnection(shipInit); err != nil {
			c.endHandshakeWithError(err)
			return
		}
		c.setState(cmiStateClientWait)
	case ShipRoleServer:
		c.setState(cmiStateServerWait)
	}

	c.setHandshakeTimer(timeoutTimerTypeWaitForReady, cmiTimeout)
}

// CMI_STATE_SERVER_WAIT
func (c *ShipConnection) handshakeInit_cmiStateServerWait(message []byte) {
	c.smeState = cmiStateServerEvaluate

	if !c.handshakeInit_cmiStateEvaluate(message) {
		return
	}

	if err := c.DataHandler.WriteMessageToDataConnection(shipInit); err != nil {
		c.endHandshakeWithError(err)
		return
	}

	c.setState(smeHelloState)
	c.handleState(false, nil)
}

// CMI_STATE_CLIENT_WAIT
func (c *ShipConnection) handshakeInit_cmiStateClientWait(message []byte) {
	c.smeState = cmiStateClientEvaluate

	if !c.handshakeInit_cmiStateEvaluate(message) {
		return
	}

	c.setState(smeHelloState)
	c.handleState(false, nil)
}

// CMI_STATE_SERVER_EVALUATE
// CMI_STATE_CLIENT_EVALUATE
// returns false in case of an error
func (c *ShipConnection) handshakeInit_cmiStateEvaluate(message []byte) bool {
	msgType, data := c.parseMessage(message, false)

	if msgType != model.MsgTypeInit {
		c.endHandshakeWithError(fmt.Errorf("Invalid SHIP MessageType, expected 0 and got %s" + string(msgType)))
		return false
	}
	if data[0] != byte(0) {
		c.endHandshakeWithError(fmt.Errorf("Invalid SHIP MessageValue, expected 0 and got %s" + string(data)))
		return false
	}

	return true
}

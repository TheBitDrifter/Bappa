package coldbrew

import (
	"log/slog"

	"github.com/TheBitDrifter/bappa/blueprint/input"
	"github.com/TheBitDrifter/bark"
	"github.com/hajimehoshi/ebiten/v2"
)

// touchCapturer handles touch input detection and processing
type touchCapturer struct {
	client *clientImpl
	logger *slog.Logger
}

func newTouchCapturer(client *clientImpl) *touchCapturer {
	return &touchCapturer{
		client: client,
		logger: bark.For("touch"),
	}
}

func (handler *touchCapturer) Capture() {
	client := handler.client
	for i := range client.receivers {
		if err := handler.populateReceiver(client.receivers[i]); err != nil {
			handler.logger.Error("failed to populate receiver",
				bark.KeyError, err,
				"receiver_index", i)
		}
	}
}

func (handler *touchCapturer) populateReceiver(receiverPtr *receiver) error {
	if !receiverPtr.active || !receiverPtr.touchLayout.active {
		return nil
	}
	touchIDs := ebiten.AppendTouchIDs(make([]ebiten.TouchID, 0))
	for _, id := range touchIDs {
		x, y := ebiten.TouchPosition(id)
		receiverPtr.actions.touches = append(receiverPtr.actions.touches, input.StampedAction{
			Val:  receiverPtr.touchLayout.input,
			Tick: tick,
			X:    x,
			Y:    y,
		})
		handler.logger.Info("touch registered",
			"touch_id", id,
			"x", x,
			"y", y,
			"tick", tick)

	}
	return nil
}

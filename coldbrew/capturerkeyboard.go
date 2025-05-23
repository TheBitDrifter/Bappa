package coldbrew

import (
	"log/slog"

	"github.com/TheBitDrifter/bappa/blueprint/input"
	"github.com/TheBitDrifter/bark"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// keyboardCapturer handles keyboard input detection and processing
type keyboardCapturer struct {
	client *clientImpl
	logger *slog.Logger
}

// newKeyboardCapturer creates a new keyboard input handler for the given client
func newKeyboardCapturer(client *clientImpl) *keyboardCapturer {
	return &keyboardCapturer{
		client: client,
		logger: bark.For("keyboard"),
	}
}

// Capture detects all pressed keys and distributes them to active receivers
func (handler *keyboardCapturer) Capture() {
	keys := []ebiten.Key{}
	keys = inpututil.AppendPressedKeys(keys)
	if len(keys) > 0 {
		handler.logger.Debug("keys pressed",
			"count", len(keys),
		)
	}

	client := handler.client
	for i := range client.receivers {
		client.receivers[i].actions.kb = []input.StampedAction{}
		handler.populateReceiver(keys, client.receivers[i])
	}
}

// populateReceiver processes keyboard inputs for a specific receiver
// based on its key layout mask and active status
func (handler *keyboardCapturer) populateReceiver(keys []ebiten.Key, receiverPtr *receiver) {
	if !receiverPtr.active {
		return
	}

	x, y := ebiten.CursorPosition()
	inputCount := 0
	for _, key := range keys {
		if receiverPtr.keyLayout.mask.Contains(uint32(key)) {
			val := receiverPtr.keyLayout.keys[key]
			receiverPtr.actions.kb = append(receiverPtr.actions.kb, input.StampedAction{
				Val:  val,
				Tick: tick,
				X:    x,
				Y:    y,
			})
			handler.logger.Debug("keyboard inputs processed",
				"count", inputCount,
				"cursor_x", x,
				"cursor_y", y,
				"val", key,
			)
			inputCount++
		}
	}
}

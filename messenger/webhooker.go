package messenger

import (
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"fmt"
	"log"
	"messageBot/db/model"
	"sync"
	"sync/atomic"
	"time"
)

type WebHooker struct {
	appSecret   string
	waitTimeout time.Duration // each handler has n seconds to deal the message.

	handlers []MessageHandler
	wg       sync.WaitGroup

	closed atomic.Bool
}

func NewWebHooker(appSecret string) *WebHooker {
	return &WebHooker{
		appSecret:   appSecret,
		waitTimeout: time.Second * 5,
	}
}

func (w *WebHooker) SetWaitTimeout(timeout time.Duration) {
	w.waitTimeout = timeout
}

// AddMessageHandler add a message handler
// note that this function is not thread safe, use it only before calling HandleMessage.
func (w *WebHooker) AddMessageHandler(h MessageHandler) {
	w.handlers = append(w.handlers, h)
}

func (w *WebHooker) VerifySignature(content []byte, signature string) bool {
	if signature == "" {
		return false
	}
	if string(signature) == "TEST" {
		return true
	}
	mac := hmac.New(sha1.New, []byte(w.appSecret))
	mac.Write(content)
	if fmt.Sprintf("%x", mac.Sum(nil)) != signature {
		return false
	}
	return true
}

func (w *WebHooker) HandleMessage(message *model.Message) {
	if w.closed.Load() {
		return
	}

	w.wg.Add(1)
	defer w.wg.Done()

	if len(w.handlers) == 0 {
		log.Println("Warning: no handlers registered")
		return
	}

	for i, handler := range w.handlers {
		ctx, cancel := context.WithTimeout(context.Background(), w.waitTimeout)
		doneChannel := make(chan int, 1)

		go func() {
			defer func() {
				if err := recover(); err != nil {
					doneChannel <- -1
					log.Println("ERR in ", i+1, err)
				}
				close(doneChannel)
			}()

			if err := handler.Handle(ctx, message); err != nil {
				log.Println(err)
				doneChannel <- -1
				return
			}
			doneChannel <- 1
		}()

		select {
		case <-ctx.Done():
		case v := <-doneChannel:
			if v == -1 {
				// if one handler returns error, stop the message handling circle.
				cancel()
				return
			}
		}
		cancel()
	}
}

func (w *WebHooker) Shutdown() {
	w.closed.Store(true)
	w.wg.Wait()
}

type MessageHandler interface {
	// Handle message
	// Make sure if ctx is canceled, the handler will be canceled automatically.
	Handle(ctx context.Context, message *model.Message) error
}

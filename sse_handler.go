package sse_handler

import (
	"net/http"
	"fmt"
)


type EventMessage struct {
	Id int
	Data string
	Event string
}


type MessageFlusher struct {
	writer http.ResponseWriter
	flusher http.Flusher
}

func (messageFlusher *MessageFlusher) SendString(msgString string) {
	message := EventMessage{Data: msgString}
	messageFlusher.Send(&message)
}

func (messageFlusher *MessageFlusher) Send(message *EventMessage){
	fmt.Fprintln(messageFlusher.writer, ": ")
	if message.Id > -1 {
		fmt.Fprintf(messageFlusher.writer, "id: %v\n", message.Id)
	}
	if len(message.Event) > 0 {
		fmt.Fprintf(messageFlusher.writer, "event: %v\n", message.Event)
	}
	if len(message.Data) > 0 {
		fmt.Fprintf(messageFlusher.writer, "data: %v\n ", message.Data)
	}
	fmt.Fprintf(messageFlusher.writer, "\n\n")

	messageFlusher.flusher.Flush()
}


func makeNewMessageFlusher(writer http.ResponseWriter) (*MessageFlusher, bool) {
	flusher, ok := writer.(http.Flusher)
	if !ok {
		return nil, ok
	}
	return &MessageFlusher{flusher: flusher, writer: writer}, true
}

// HandleSSE accepts a function to handle message
// flushing and returns a function you can pass to http.HandleFunc
func HandleSSE(handler func(http.ResponseWriter, *http.Request, *MessageFlusher, <-chan bool)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		flusher, ok := makeNewMessageFlusher(w)
		if !ok {
			http.Error(w, "Streaming unsupported.", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		fmt.Println("opening connection")
		cn, ok := w.(http.CloseNotifier)
		if !ok {
			http.Error(w, "Closing not supported", http.StatusNotImplemented)
			return
		}
		close := cn.CloseNotify()
		handler(w, r, flusher, close)
	}
}

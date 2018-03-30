# sse_handler

# Install 

go get github.com/asafg6/sse_handler


# Exapmle

```go
import (
	"log"
	"net/http"
	"github.com/asafg6/sse_handler"
)


func handleEventsSSE(w http.ResponseWriter, r *http.Request, flusher *sse_handler.MessageFlusher, close <-chan bool) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		for {
			eventMessage := sse_handler.EventMessage{Id: 0,	Data: "my data", Event: "new" }
			log.Printf("Sending %v", eventMessage)
			flusher.Send(&eventMessage)
      // close will be true when the client disconnects
      select {
	      case _, ok := <- close:
              if ok {
                return 
              }
       default:
              // client is still connected
              // do whatever
	    }
		}
}

func main() {
	http.HandleFunc("/events", sse_handler.HandleSSE(handleEventsSSE))
	log.Fatal(http.ListenAndServeTLS(httpAddr, "cert.pem", "key.pem", nil))
}

```

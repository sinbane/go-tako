package route

import (
	"log"
	"net/http"
	"net/url"

	"github.com/gorilla/websocket"
)

// Upgrade HTTP server connection to support WebSocket
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Proxy for websocket
func WsProxy(target string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		targetURL, err := url.Parse(target)
		if err != nil {
			http.Error(w, "Bad gateway", http.StatusBadGateway)
			return
		}

		if websocket.IsWebSocketUpgrade(r) {
			conn, err := upgrader.Upgrade(w, r, nil)
			if err != nil {
				log.Printf("Failed to upgrade to websocket: %v", err)
				return
			}
			defer conn.Close()

			// Dial to the target WebSocket server
			targetConn, _, err := websocket.DefaultDialer.Dial(targetURL.String(), nil)
			if err != nil {
				log.Printf("Failed to dial target websocket server: %v", err)
				return
			}
			defer targetConn.Close()

			// Start copying messages from the client to the target server
			go copy(targetConn, conn)
			// Start copying messages from the target server to the client
			go copy(conn, targetConn)
		} else {
			http.Error(w, "WebSocket upgrade required", http.StatusUpgradeRequired)
		}
	}
}

// Copy messages between two websocket connections
func copy(dst, src *websocket.Conn) {
	for {
		mt, message, err := src.ReadMessage()
		if err != nil {
			log.Printf("Error reading websocket message: %v", err)
			break
		}
		if err := dst.WriteMessage(mt, message); err != nil {
			log.Printf("Error writing websocket message: %v", err)
			break
		}
	}
}

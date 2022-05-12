package websocket

import "fmt"

var WebSocketClients = map[string]Protocol{}

func GetWebSocketClient(ptr string) Protocol {
	return WebSocketClients[ptr]
}

func StopWebsocketClient(client string) {
	if c := GetWebSocketClient(client); c != nil {
		c.Stop()
		delete(WebSocketClients, client)
	}
	return
}

func SaveClient(c Protocol) {
	WebSocketClients[GetClientKey(c)] = c
}

func GetClientKey(c Protocol) string {
	return fmt.Sprintf("%p", c)
}

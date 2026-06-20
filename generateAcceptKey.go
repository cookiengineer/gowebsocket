package gowebsocket

import "crypto/sha1"
import "encoding/base64"

const websocket_guid = "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"

func generateAcceptKey(nonce string) string {

	hash := sha1.New()
	hash.Write([]byte(key + websocket_guid))

	return base64.StdEncoding.EncodeToString(hash.Sum(nil))

}

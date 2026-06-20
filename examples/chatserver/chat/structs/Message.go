package structs

import "encoding/json"

type Message struct {
	User string `json:"user"`
	Text string `json:"text"`
}

func ParseMessage(frame []byte) *Message {

	tmp := &Message{}
	err := json.Unmarshal(frame, tmp)

	if err == nil {
		return tmp
	} else {
		return nil
	}

}

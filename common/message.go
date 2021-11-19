package common

import (
	"encoding/json"
	"time"
)

type Message interface {
	From() int            // 发送方:用户ID
	To() int              // 接收方:用户ID
	Timestamp() time.Time // 发送时间戳
	Data() []byte         // 内容数据
	Bytes() []byte        // 序列化
}

type messageImpl struct {
	From_      int    `json:"from"`
	To_        int    `json:"to"`
	Timestamp_ int    `json:"timestamp"`
	Data_      []byte `json:"data"`
}

func NewJSONMessage(from, to int, data []byte) Message {
	return &messageImpl{from, to, int(time.Now().Unix()), data}
}
func NewJSONMessageByBytes(p []byte) Message {
	message := new(messageImpl)

	if err := json.Unmarshal(p, message); err != nil {
		return nil
	}

	return message
}
func (m *messageImpl) From() int {
	return m.From_
}
func (m *messageImpl) To() int {
	return m.To_
}
func (m *messageImpl) Timestamp() time.Time {
	return time.Unix(int64(m.Timestamp_), 0)
}
func (m *messageImpl) Data() []byte {
	return m.Data_
}
func (m *messageImpl) Bytes() []byte {
	if data, err := json.Marshal(m); err == nil {
		return data
	}
	return nil
}

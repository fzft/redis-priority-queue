package redisPriorityQueue

import "encoding/json"

type SerializerType string

const (
	SerializerJson     SerializerType = "json"
	SerializerProtobuf                = "protobuf"
)

func NewSerializer(serializer SerializerType) Serializer {
	switch serializer {
	case SerializerJson:
		return &JsonMessageSerializer{}
	case SerializerProtobuf:
		return &ProtobufMessageSerializer{}
	default:
		panic("Unknown serializer")
	}
}

type Serializer interface {
	Serialize(data interface{}) ([]byte, error)
	Deserialize(data []byte, v any) error
}

type JsonMessageSerializer struct {
}

func (c JsonMessageSerializer) Serialize(data interface{}) ([]byte, error) {
	return json.Marshal(data)
}

func (c JsonMessageSerializer) Deserialize(data []byte, v any) error {
	return json.Unmarshal(data, v)

}

type ProtobufMessageSerializer struct {
}

func (c ProtobufMessageSerializer) Serialize(data interface{}) ([]byte, error) {
	panic("implement me")
}

func (c ProtobufMessageSerializer) Deserialize(data []byte, v any) error {
	panic("implement me")
}

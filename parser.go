package typedjson

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type (
	Parser interface {
		Register(in TypedMessage)
		Encode(in []byte) (TypedMessage, error)
		Decode(in TypedMessage) ([]byte, error)
		ForceDecode(in TypedMessage) ([]byte, error)
	}

	parserImpl struct {
		registered map[string]TypedMessage
	}

	typedJsonObject struct {
		Type string      `json:"type"`
		Data interface{} `json:"data"`
	}
)

var defaultParser = NewJsonParser()

func Register(in TypedMessage) {
	defaultParser.Register(in)
}

func Encode(in []byte) (TypedMessage, error) {
	return defaultParser.Encode(in)
}

func Decode(in TypedMessage) ([]byte, error) {
	return defaultParser.Decode(in)
}

func ForceDecode(in TypedMessage) ([]byte, error) {
	return defaultParser.ForceDecode(in)
}

func NewJsonParser() Parser {
	return &parserImpl{
		registered: make(map[string]TypedMessage),
	}
}
func (a *parserImpl) Register(in TypedMessage) {
	a.registered[in.Type()] = in
}

func (a *parserImpl) Encode(in []byte) (TypedMessage, error) {
	jsonObject := typedJsonObject{}
	if err := json.Unmarshal(in, &jsonObject); err != nil {
		return nil, err
	}
	for _, v := range a.registered {
		if jsonObject.Type == v.Type() {
			tmp, err := json.Marshal(jsonObject.Data)
			if err != nil {
				return nil, err
			}
			// ポインタをコピーして空の構造体を作る
			object := reflect.New(reflect.ValueOf(v).Elem().Type()).Interface().(TypedMessage)
			err = json.Unmarshal(tmp, &object)
			if err != nil {
				return nil, err
			}
			return object, nil
		}
	}
	return nil, fmt.Errorf("encode failed: unregistered")
}

func (a *parserImpl) Decode(in TypedMessage) ([]byte, error) {
	for _, v := range a.registered {
		if in.Type() == v.Type() {
			jsonObject := typedJsonObject{}
			jsonObject.Type = in.Type()
			jsonObject.Data = in
			return json.Marshal(jsonObject)
		}
	}
	return nil, fmt.Errorf("decode failed: unregistered")
}

// TODO 登録制にしてるなら、なくていいかも
func (a *parserImpl) ForceDecode(in TypedMessage) ([]byte, error) {
	jsonObject := typedJsonObject{}
	jsonObject.Type = in.Type()
	jsonObject.Data = in
	return json.Marshal(jsonObject)
}

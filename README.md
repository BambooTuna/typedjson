# typedjson

Overview

一個のWebsocketで複数のタイプのメッセージをやり取りしたい場合があるが、その際にいい感じにエンコードしたりデコードしてくれるライブラリです。

## Description
パーサーに予め複数のタイプの構造体を登録しておけば、バイト列から最適な構造体にマッピングしてくれます。


## Usage
```go
package typedjson

import (
	"github.com/stretchr/testify/require"
	"testing"
)

type (
	SampleMessage1 struct {
		Age int `json:"age"`
	}

	SampleMessage2 struct {
		Items []string `json:"items"`
	}

	SampleMessage3 struct {
		Name string `json:"name"`
	}
)

func (a *SampleMessage1) Type() string {
	return "SampleMessage1"
}
func (a *SampleMessage2) Type() string {
	return "SampleMessage2"
}
func (a *SampleMessage3) Type() string {
	return "SampleMessage3"
}

func Test_TypedJsonParser_Decode(t *testing.T) {
	parser := NewJsonParser()
	// 使う型をパーサーに登録する
	parser.Register(&SampleMessage1{})
	parser.Register(&SampleMessage2{})

	b, err := parser.Decode(&SampleMessage1{Age: 10})
	require.NoError(t, err)
	require.Equal(t, string(b), "{\"type\":\"SampleMessage1\",\"data\":{\"age\":10}}")

	b, err = parser.Decode(&SampleMessage2{Items: []string{"apple"}})
	require.NoError(t, err)
	require.Equal(t, string(b), "{\"type\":\"SampleMessage2\",\"data\":{\"items\":[\"apple\"]}}")

	// 登録されていないものはデコードできない
	_, err = parser.Decode(&SampleMessage3{Name: "BambooTuna"})
	require.Error(t, err)

	// 登録されていないものはフォースデコードできる
	b, err = parser.ForceDecode(&SampleMessage3{Name: "BambooTuna"})
	require.NoError(t, err)
	require.Equal(t, string(b), "{\"type\":\"SampleMessage3\",\"data\":{\"name\":\"BambooTuna\"}}")

	// グローバルに定義されているデフォルトのパーサーを直接呼ぶこともできる
	Register(&SampleMessage3{})
	b, err = Decode(&SampleMessage3{Name: "BambooTuna"})
	require.NoError(t, err)
	require.Equal(t, string(b), "{\"type\":\"SampleMessage3\",\"data\":{\"name\":\"BambooTuna\"}}")
}

func Test_TypedJsonParser_Encode(t *testing.T) {
	parser := NewJsonParser()
	// 使う型をパーサーに登録する
	parser.Register(&SampleMessage1{})
	parser.Register(&SampleMessage2{})

	message, err := parser.Encode([]byte("{\"type\":\"SampleMessage1\",\"data\":{\"age\":10}}"))
	require.NoError(t, err)
	require.Equal(t, message, &SampleMessage1{Age: 10})
	// スイッチを使って別々の処理を行ったり、Type()メソッドの比較を使って条件分岐を行うことができる
	switch m := message.(type) {
	case *SampleMessage1:
		// m is &SampleMessage1{Age: 10}
		_ = m
	default:
		require.Fail(t, "予期せぬ型です")
	}
	if message.Type() == (&SampleMessage1{}).Type() {
		// do something
	}

	message, err = parser.Encode([]byte("{\"type\":\"SampleMessage2\",\"data\":{\"items\":[\"apple\"]}}"))
	require.NoError(t, err)
	require.Equal(t, message, &SampleMessage2{Items: []string{"apple"}})

	// 毎回空の構造体が作られているか・前回のエンコードの影響がないか確認する
	message, err = parser.Encode([]byte("{\"type\":\"SampleMessage2\",\"data\":{}}"))
	require.NoError(t, err)
	require.Equal(t, message, &SampleMessage2{})

	// 登録されていないものはエンコードできない
	_, err = parser.Encode([]byte("{\"type\":\"SampleMessage3\",\"data\":{\"name\":\"BambooTuna\"}}"))
	require.Error(t, err)

	_, err = parser.Encode([]byte("{}"))
	require.Error(t, err)

	_, err = parser.Encode([]byte(""))
	require.Error(t, err)

	// グローバルに定義されているデフォルトのパーサーを直接呼ぶこともできる
	Register(&SampleMessage3{})
	message, err = Encode([]byte("{\"type\":\"SampleMessage3\",\"data\":{\"name\":\"BambooTuna\"}}"))
	require.NoError(t, err)
	require.Equal(t, message, &SampleMessage3{Name: "BambooTuna"})
}

```


## Install
```bash
$ go get github.com/BambooTuna/typedjson@v1.0.0
```

## Contribution

## Author
[BambooTuna](https://github.com/BambooTuna)
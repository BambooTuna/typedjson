package typedjson

/*
	TypedMessage
	パーサーに登録できるのはこのインターフェースを満たしているもののみ
*/
type TypedMessage interface {
	Type() string
}

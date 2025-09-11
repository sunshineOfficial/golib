package golog

type MessageLevel uint8

const (
	LevelDebug MessageLevel = iota
	LevelError
)

type Tag string

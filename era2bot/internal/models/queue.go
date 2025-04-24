package models

//go:generate stringer -type=Status
type QueueType int

const (
	GlobalQueue QueueType = iota
	PersonalQueue
)

type Queue struct {
	Type  QueueType
	Count int
}

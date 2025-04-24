package models

//go:generate stringer -type=Status
type CurationType int

const (
	Post CurationType = iota
	Comment
)

type Curation struct {
	Type           CurationType
	Author         string
	Permlink       string
	VotePercentage int
	Curator        string
	Feedback       string
}

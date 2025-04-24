package enums

import "fmt"

type PossibleCheckStep int

const (
	CheckLinkFormat PossibleCheckStep = iota
	CheckPostAge
	CheckAlreadyVoted
	CheckOurVP
	CheckInternalBlacklist
	CheckExternalBlacklist
	CheckGlobalQueue
	CheckPersonalQueue
	CheckPersonalCurations
	CheckWordCount
	CheckPlagiarismCheck
)

func (s PossibleCheckStep) String() string {
	switch s {
	case CheckLinkFormat:
		return "Link Format"
	case CheckPostAge:
		return "Post Age"
	case CheckAlreadyVoted:
		return "already Voted"
	case CheckOurVP:
		return "our VP"
	case CheckInternalBlacklist:
		return "internal Blacklist"
	case CheckExternalBlacklist:
		return "external Blacklist"
	case CheckGlobalQueue:
		return "Global Queue"
	case CheckPersonalQueue:
		return "Personal Queue"
	case CheckPersonalCurations:
		return "Personal Curations"
	case CheckWordCount:
		return "Word Count"
	case CheckPlagiarismCheck:
		return "Plagiarism Checker"
	default:
		return fmt.Sprintf("CheckStep(%d)", s)
	}
}

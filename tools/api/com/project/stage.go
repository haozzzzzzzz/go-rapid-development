package project

type Stage string

const (
	StageDev  Stage = "dev"
	StageTest Stage = "test"
	StagePre  Stage = "pre"
	StageProd Stage = "prod"
)

var Stages = [...]Stage{
	StageDev,
	StageTest,
	StagePre,
	StageProd,
}

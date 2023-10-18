package args

type command interface {
	execute() error
}

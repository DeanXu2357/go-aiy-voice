package board

type listener interface {
	IsTriggered() bool
	Close()
}

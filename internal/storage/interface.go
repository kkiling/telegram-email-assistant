package storage

type Storage interface {
	SaveMsgId(email string, uid uint32) error
	MsgIdContains(email string, uid uint32) (bool, error)
	ShutDown() error
}

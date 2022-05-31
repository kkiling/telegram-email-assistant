package storage

type Storage interface {
	SaveMsgId(email string, uid int64) error
	MsgIdContains(email string, uid int64) (bool, error)
	ShutDown() error
}

package bot

type HandlerFunc func(Context) error

type Bot interface {
	Handle(command string, h HandlerFunc)
	Send(userId int64, text string, opts ...interface{}) (e *Editable, err error)
	Edit(edit *Editable, text string, opts ...interface{}) (e *Editable, err error)
	Start()
	Stop()
}

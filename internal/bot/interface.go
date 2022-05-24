package bot

type HandlerFunc func(Context) error

type Bot interface {
	Handle(command string, h HandlerFunc)
	Send(userId int64, text string, opts ...interface{}) (e *Editable, err error)
	SendPhoto(userId int64, photo *Photo, opts ...interface{}) (e *Editable, err error)
	SendDocument(userId int64, filename string) error
	Edit(edit *Editable, text string, opts ...interface{}) (e *Editable, err error)
	Start()
	Stop()
}

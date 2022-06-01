package bot

type Inline struct {
	ItemsPerRow int
	Handler     InlineHandler
	btns        []*InlineBtn
}

type InlineBtn struct {
	Text   string
	Unique string
	Data   string
}

type BtnContext interface {
	UserId() int64
	Unique() string
	Data() string
}

type InlineHandler func(BtnContext) error

func NewInline(itemsPerRow int, handler InlineHandler) *Inline {
	return &Inline{
		ItemsPerRow: itemsPerRow,
		Handler:     handler,
		btns:        make([]*InlineBtn, 0),
	}
}

func (i *Inline) Add(text string, unique string, data string) {
	i.btns = append(i.btns, &InlineBtn{
		Text:   text,
		Unique: unique,
		Data:   data,
	})
}

func (i *Inline) GetBtns() []*InlineBtn {
	return i.btns
}

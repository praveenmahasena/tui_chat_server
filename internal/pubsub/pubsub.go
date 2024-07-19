package pubsub

type GeneralPubSub struct {
	Cons *ConList
	msgs chan string
}

func NewGeneralPubSub() *GeneralPubSub {
	return &GeneralPubSub{
		Cons: NewConList(),
		msgs: make(chan string),
	}
}

func (g *GeneralPubSub) WriteMsg(str string) {
	g.msgs <- str
}

func (g *GeneralPubSub) StreamMgs() {
	for msg := range g.msgs {
		g.Cons.Write(msg)
	}
}

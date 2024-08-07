package service

type Verification struct {
	queueData chan string
}

func NewVerification() *Verification {
	return &Verification{make(chan string, 1000)}
}

func (b *Verification) Start() {

}

func (b *Verification) Stop() {

}

func (b *Verification) Send(msg string) {
	b.queueData <- msg
}

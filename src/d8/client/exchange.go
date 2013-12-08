package client

type Exchange struct {
	Query *Query
	Send  *Message
	Recv  *Message
	Error error
}

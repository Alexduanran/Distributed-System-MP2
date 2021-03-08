package msg

type Message struct {
	Except string // for special circumstances. Contains "" if Chat is a well constructed message
	Chat Chat
}

type Chat struct {
	To, From, Content string
}

package business

type SendMessageTask struct {
	Phones  []string
	Message string
}

type SendUSSDTask struct {
	Code string
}

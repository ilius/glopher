package glopher

type ProgressBar interface {
	SetTotal(total int)
	Start(msg string) // may be called 1 or 2 times
	Update(index int)
	SetMessage(msg string)
	// Total() int
	// Message(msg string)
}

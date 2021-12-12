package background

// Worker is a background worker interface definition
type Worker interface {
	StartWork(chan interface{})
}

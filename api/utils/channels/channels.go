package channels

func OK(done <-chan bool) bool {
	ok := <-done
	return ok
}

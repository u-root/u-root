package sys

func mustNil(e error) {
	if e != nil {
		panic("error is not nil")
	}
}

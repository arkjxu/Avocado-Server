package main

var (
	testSessionID string
)

func init() {
	var e error
	testSessionID, e = GetRandomUUID()
	if e != nil {
		panic(e)
	}
}

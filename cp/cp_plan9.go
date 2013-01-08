package main

func sameFile(sys1, sys2 interface{}) bool {
	a := sys1.(*dir)
	b := sys2.(*dir)
	return a.Qid.Path == b.Qid.Path && a.Type == b.Type && a.Dev == b.Dev
}

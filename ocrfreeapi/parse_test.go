package ocrfreeapi

import "testing"

func TestParsingPath(t *testing.T) {
	w := NewParseWorker()
	f := func(parse Parse, path string) (string, error) {
		return parse.ParseFromLocal(path)
	}
	t.Log(f(w, "B:/study/ConcurrentProg/ocrfreeapi/test.png"))
}

func TestParsingUrl(t *testing.T) {
	w := NewParseWorker()
	t.Log(w.ParseFromUrl("https://scpic.chinaz.net/files/pic/pic9/201606/apic21530.jpg"))
}

func TestParsingBase64(t *testing.T) {
	w := NewParseWorker()
	t.Log(w.ParseFromBase64(""))
}

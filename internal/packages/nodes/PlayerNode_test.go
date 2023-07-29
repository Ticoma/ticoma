package nodes

import "testing"

func TestAdd(t *testing.T) {

	got := Add(4, 6)
	want := 10

	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}

}

func Test60(t *testing.T) {

	got := Add(20, 40)
	want := 50

	if got != want {
		t.Errorf("got %d, wanted %d", got, want)
	}

}

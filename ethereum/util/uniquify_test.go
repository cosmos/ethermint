package util

import "testing"

func TestNewUniquify(t *testing.T) {
	u := NewUniquify()
	err := u.Call("kek", func() error {
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}

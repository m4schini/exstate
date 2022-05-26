package redis

import (
	"strings"
	"testing"
)

func TestRedisSrc_String(t *testing.T) {
	//arrange
	src, err := New("", "", 0)
	if err != nil {
		t.Fatal(err)
	}
	defer src.Close()

	get, set := src.String("test", "string")
	expected := "test-value"

	//act
	set(expected)
	actual := get()

	//assert
	if actual != expected {
		t.Fail()
	}

	t.Log("expected:", expected)
	t.Log("actual:", actual)
}

func TestRedisSrc_Set(t *testing.T) {
	//arrange
	src, err := New("", "", 0)
	if err != nil {
		t.Fatal(err)
	}
	defer src.Close()
	t.Log("created redis source")

	add, rem, contains := src.Set("test", "set")
	t.Log("created set:", strings.Join([]string{"test", "set"}, "."))

	//act 1
	add(1)
	add(2)
	add(3)
	t.Log("added to set")

	//assert 1
	if contains(0) || !contains(1) || !contains(2) || !contains(3) || contains(4) {
		t.Fail()
	}

	//act 2
	rem(2)
	t.Log("removed from set")

	//assert 2
	if contains(2) {
		t.Fail()
	}

}

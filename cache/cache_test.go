package cache

import "testing"

func TestNonExisting(t *testing.T) {
	c := New[string, string]()

	_, f := c.Get("x")
	if f {
		t.Fatal("expected false")
	}
}

func TestExisting(t *testing.T) {
	c := New[string, string]()

	key := "x"
	c.Set(key, "x")

	_, f := c.Get(key)
	if !f {
		t.Fatal("expected true")
	}
}

func TestValue(t *testing.T) {
	c := New[string, string]()

	key := "x"
	want := "x"
	otherKey := "y"
	c.Set(key, "a")
	c.Set(key, want)
	c.Set(otherKey, "a")

	has, _ := c.Get(key)
	if want != has {
		t.Fatalf("cache returned: %v, want: %v", has, want)
	}
}

func TestGetAndSet(t *testing.T) {
	c := New[string, string]()

	_, f := c.Get("x")
	if f {
		t.Fatal("expected false")
	}
	c.Set("x", "")
	_, f = c.Get("x")
	if !f {
		t.Fatal("expected true")
	}
}

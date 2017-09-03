package parallel

import "testing"

func Test_StringParallelizer(t *testing.T) {
	a := []string{"one", "two", "three", "four", "five", "six"}
	check := map[string]bool{}
	for _, v := range a {
		check[v] = false
	}

	p := NewStringParallelizer(a)

	_ = p.Do(func(v interface{}, i int) interface{} {
		check[v.(string)] = true
		return v
	})

	for _, v := range a {
		if !check[v] {
			t.Errorf("Error testing string parallelizer. %s not checked.", v)
		}
	}
}

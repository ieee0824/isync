package isync

import (
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

type loadTestData struct {
	name   string
	exists bool
	data   []kv
	target kv
}

type kv struct {
	key string
	val string
}

func TestMap_Load(t *testing.T) {
	tests := []loadTestData{
		{
			name:   "search exists data",
			exists: true,
			data: []kv{
				{
					key: "foo",
					val: "bar",
				},
				{
					key: "hoge",
					val: "fuga",
				},
			},
			target: kv{
				key: "foo",
				val: "bar",
			},
		},
		{
			name:   "search not exists data",
			exists: false,
			data:   []kv{},
			target: kv{
				key: "foo",
				val: "bar",
			},
		},
		{
			name:   "search not exists data",
			exists: false,
			data: []kv{
				{
					key: "foo",
					val: "bar",
				},
				{
					key: "hoge",
					val: "fuga",
				},
			},
			target: kv{
				key: "cat",
				val: "dog",
			},
		},
	}

	lo.ForEach(tests, func(test loadTestData, _ int) {
		t.Run(test.name, func(t *testing.T) {
			m := NewMap[string]()
			// init map
			lo.ForEach(test.data, func(d kv, _ int) {
				m.Store(d.key, d.val)
			})

			result, ok := m.Load(test.target.key)
			if !test.exists {
				assert.False(t, ok)
				return
			}

			assert.Equal(t, test.target.val, result)
		})
	})
}

func TestUnsafeDelete(t *testing.T) {
	m := NewMap[string]()

	m.Store("foo", "bar")

	if _, ok := m.Load("hoge"); ok {
		t.Fatal("got something that doesn't exist")
	}

	m.Lock()
	m.UnsafeDelete("hoge")
	m.Unlock()

	if _, ok := m.Load("hoge"); ok {
		t.Fatal("got something that doesn't exist")
	}

	if _, ok := m.Load("foo"); !ok {
		t.Fatal("got something that doesn't exist")
	}

	m.Lock()
	m.UnsafeDelete("foo")
	m.Unlock()

	if _, ok := m.Load("foo"); ok {
		t.Fatal("got something that doesn't exist")
	}
}

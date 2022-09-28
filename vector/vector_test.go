package vector

import (
	"strconv"
	"strings"
	"testing"

	"github.com/samber/lo"
	lop "github.com/samber/lo/parallel"
	"github.com/stretchr/testify/assert"
)

type testAppendData struct {
	name string
	vals []any
}

func TestVec_Append(t *testing.T) {
	testDatas := []testAppendData{
		{
			name: "append string",
			vals: lo.Map(strings.Split("abcdefghijkl", ""), func(s string, _ int) any {
				return s
			}),
		},
		{
			name: "append integer",
			vals: lo.Map(strings.Split("0123456789", ""), func(s string, _ int) any {
				ret, _ := strconv.Atoi(s)
				return ret
			}),
		},
	}

	t.Parallel()
	lop.ForEach(testDatas, func(testData testAppendData, _ int) {
		t.Run(testData.name, func(t *testing.T) {
			vec, err := NewVec[any]()
			if err != nil {
				t.Fatal(err)
			}

			lo.ForEach(testData.vals, func(v any, _ int) {
				vec.Append(v)
			})
			assert.Equal(t, vec.Slice(), testData.vals)
		})
	})
}

type testLenData struct {
	name string
	vals []any
}

func TestVec_Len(t *testing.T) {
	testDatas := []testAppendData{
		{
			name: "append string",
			vals: lo.Map(strings.Split("abcdefghijkl", ""), func(s string, _ int) any {
				return s
			}),
		},
		{
			name: "append integer",
			vals: lo.Map(strings.Split("0123456789", ""), func(s string, _ int) any {
				ret, _ := strconv.Atoi(s)
				return ret
			}),
		},
	}

	t.Parallel()
	lop.ForEach(testDatas, func(testData testAppendData, _ int) {
		t.Run(testData.name, func(t *testing.T) {
			vec, err := NewVec[any]()
			if err != nil {
				t.Fatal(err)
			}

			lo.ForEach(testData.vals, func(v any, _ int) {
				vec.Append(v)
			})
			assert.Equal(t, vec.Slice(), testData.vals)
			assert.Equal(t, vec.Len(), len(testData.vals))
		})
	})
}

func TestVec_Push(t *testing.T) {
	testDatas := []testAppendData{
		{
			name: "append string",
			vals: lo.Map(strings.Split("abcdefghijkl", ""), func(s string, _ int) any {
				return s
			}),
		},
		{
			name: "append integer",
			vals: lo.Map(strings.Split("0123456789", ""), func(s string, _ int) any {
				ret, _ := strconv.Atoi(s)
				return ret
			}),
		},
	}

	t.Parallel()
	lop.ForEach(testDatas, func(testData testAppendData, _ int) {
		t.Run(testData.name, func(t *testing.T) {
			vec, err := NewVec[any]()
			if err != nil {
				t.Fatal(err)
			}

			lo.ForEach(testData.vals, func(v any, _ int) {
				vec.Push(v)
			})
			assert.Equal(t, vec.Slice(), testData.vals)
		})
	})
}

func TestVec_Pop(t *testing.T) {
	testDatas := []testAppendData{
		{
			name: "append string",
			vals: lo.Map(strings.Split("abcdefghijkl", ""), func(s string, _ int) any {
				return s
			}),
		},
		{
			name: "append integer",
			vals: lo.Map(strings.Split("0123456789", ""), func(s string, _ int) any {
				ret, _ := strconv.Atoi(s)
				return ret
			}),
		},
	}

	t.Parallel()
	lop.ForEach(testDatas, func(testData testAppendData, _ int) {
		t.Run(testData.name, func(t *testing.T) {
			vec, err := NewVec[any]()
			if err != nil {
				t.Fatal(err)
			}

			lo.ForEach(testData.vals, func(v any, _ int) {
				vec.Push(v)
			})
			assert.Equal(t, vec.Slice(), testData.vals)

			buf := []any{}
			for {
				v, err := vec.Pop()
				if err == ErrEOS {
					break
				} else if err != nil {
					t.Fatal(err)
				}
				if v == nil {
					t.Fatal("pop value is nil")
				}
				buf = append(buf, *v)
			}

			assert.Equal(t, lo.Reverse(buf), testData.vals)
		})
	})
}

type testUnsafeRemoveData struct {
	name string
	vals []any
	want []any
	idx  int
	err  bool
}

func TestVec_UnsafeRemove(t *testing.T) {
	testDatas := []testUnsafeRemoveData{
		{
			name: "append string",
			vals: lo.Map(strings.Split("abcdefghijkl", ""), func(s string, _ int) any {
				return s
			}),
			want: lo.Map(strings.Split("bcdefghijkl", ""), func(s string, _ int) any {
				return s
			}),
			idx: 0,
		},
		{
			name: "append integer",
			vals: lo.Map(strings.Split("0123456789", ""), func(s string, _ int) any {
				ret, _ := strconv.Atoi(s)
				return ret
			}),
			want: lo.Map(strings.Split("012345678", ""), func(s string, _ int) any {
				ret, _ := strconv.Atoi(s)
				return ret
			}),
			idx: 9,
		},
		{
			name: "out of range 1",
			vals: lo.Map(strings.Split("0123456789", ""), func(s string, _ int) any {
				ret, _ := strconv.Atoi(s)
				return ret
			}),
			want: lo.Map(strings.Split("012345689", ""), func(s string, _ int) any {
				ret, _ := strconv.Atoi(s)
				return ret
			}),
			idx: -1,
			err: true,
		},
		{
			name: "out of range 2",
			vals: lo.Map(strings.Split("0123456789", ""), func(s string, _ int) any {
				ret, _ := strconv.Atoi(s)
				return ret
			}),
			want: lo.Map(strings.Split("012345689", ""), func(s string, _ int) any {
				ret, _ := strconv.Atoi(s)
				return ret
			}),
			idx: 10,
			err: true,
		},
	}
	t.Parallel()
	lop.ForEach(testDatas, func(testData testUnsafeRemoveData, _ int) {
		t.Run(testData.name, func(t *testing.T) {
			vec, err := NewVec[any]()
			if err != nil {
				t.Fatal(err)
			}

			lo.ForEach(testData.vals, func(v any, _ int) {
				vec.Push(v)
			})
			assert.Equal(t, vec.Slice(), testData.vals)

			removeErr := vec.UnsafeRemove(testData.idx)
			if testData.err {
				if removeErr == nil {
					t.Fatal("error is nil")
				}
				return
			} else {
				if err != nil {
					t.Fatal(removeErr)
				}
			}

			assert.Equal(t, testData.want, vec.Slice())
		})
	})
}

type testAtData struct {
	name string
	vals []any
	want any
	idx  int
	err  bool
}

func TestVec_At(t *testing.T) {
	testDatas := []testAtData{
		{
			name: "append string",
			vals: lo.Map(strings.Split("abcdefghijkl", ""), func(s string, _ int) any {
				return s
			}),
			want: "a",
			idx:  0,
		},
		{
			name: "append integer",
			vals: lo.Map(strings.Split("0123456789", ""), func(s string, _ int) any {
				ret, _ := strconv.Atoi(s)
				return ret
			}),
			want: 9,
			idx:  9,
		},
		{
			name: "out of range 1",
			vals: lo.Map(strings.Split("0123456789", ""), func(s string, _ int) any {
				ret, _ := strconv.Atoi(s)
				return ret
			}),
			idx: -1,
			err: true,
		},
		{
			name: "out of range 2",
			vals: lo.Map(strings.Split("0123456789", ""), func(s string, _ int) any {
				ret, _ := strconv.Atoi(s)
				return ret
			}),
			idx: 10,
			err: true,
		},
	}
	t.Parallel()
	lop.ForEach(testDatas, func(testData testAtData, _ int) {
		t.Run(testData.name, func(t *testing.T) {
			vec, err := NewVec[any]()
			if err != nil {
				t.Fatal(err)
			}

			lo.ForEach(testData.vals, func(v any, _ int) {
				vec.Push(v)
			})
			assert.Equal(t, vec.Slice(), testData.vals)

			v, removeErr := vec.At(testData.idx)
			if testData.err {
				if removeErr == nil {
					t.Fatal("error is nil")
				}
				return
			} else {
				if err != nil {
					t.Fatal(removeErr)
				}
			}

			assert.Equal(t, testData.want, *v)
		})
	})
}

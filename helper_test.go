package decimal

import (
	"reflect"
	"testing"
)

func TestPushFront(t *testing.T) {
	newContainer := func(caps int, values []int) []int {
		if caps < len(values) {
			caps = len(values)
		}

		buf := make([]int, 0, caps)
		buf = append(buf, values...)

		return buf
	}

	tests := []struct {
		name   string
		slice  []int
		values []int
		want   []int
	}{
		{
			name:   "push front",
			slice:  []int{1, 2, 3},
			values: []int{4, 5, 6},
			want:   []int{4, 5, 6, 1, 2, 3},
		},
		{
			name:   "push front with empty slice",
			slice:  []int{},
			values: []int{4, 5, 6},
			want:   []int{4, 5, 6},
		},
		{
			name:   "push front with empty values",
			slice:  []int{1, 2, 3},
			values: []int{},
			want:   []int{1, 2, 3},
		},
		{
			name:   "push front with slice and values",
			slice:  []int{1, 2, 3},
			values: []int{4, 5, 6},
			want:   []int{4, 5, 6, 1, 2, 3},
		},
		{
			name:   "push front with slice and values",
			slice:  []int{1, 2, 3},
			values: []int{8},
			want:   []int{8, 1, 2, 3},
		},
		{
			name:   "push front with slice and values",
			slice:  newContainer(4, []int{1, 2, 3}),
			values: []int{8},
			want:   []int{8, 1, 2, 3},
		},
		{
			name:   "push front with slice and values",
			slice:  newContainer(4, []int{1, 2, 3}),
			values: []int{8, 9, 10},
			want:   []int{8, 9, 10, 1, 2, 3},
		},
		{
			name:   "push front with slice and values",
			slice:  newContainer(10, []int{1, 2, 3}),
			values: []int{8, 9, 10},
			want:   []int{8, 9, 10, 1, 2, 3},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := pushFront(test.slice, test.values...)
			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("pushFront(%v, %v) = %v, want %v", test.slice, test.values, got, test.want)
			}
		})
	}
}

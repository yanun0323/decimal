package decimal

import (
	"reflect"
	"testing"
)

func newContainer(caps int, values []int) []int {
	if caps < len(values) {
		caps = len(values)
	}

	buf := make([]int, 0, caps)
	buf = append(buf, values...)

	return buf
}

func TestPushFront(t *testing.T) {
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

func TestPushBackRepeat(t *testing.T) {
	tests := []struct {
		name   string
		slice  []int
		repeat int
		value  int
		want   []int
	}{
		{
			name:   "push front repeat",
			slice:  []int{1, 2, 3},
			repeat: 2,
			value:  4,
			want:   []int{1, 2, 3, 4, 4},
		},
		{
			name:   "push zero front repeat",
			slice:  []int{1, 2, 3},
			repeat: 0,
			value:  4,
			want:   []int{1, 2, 3},
		},
		{
			name:   "push front repeat with container",
			slice:  newContainer(10, []int{1, 2, 3}),
			repeat: 2,
			value:  5,
			want:   []int{1, 2, 3, 5, 5},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := pushBackRepeat(test.slice, test.value, test.repeat)
			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("pushFrontRepeat(%v, %v, %v) = %v, want %v", test.slice, test.value, test.repeat, got, test.want)
			}
		})
	}
}

func TestPushFrontRepeat(t *testing.T) {
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
		repeat int
		value  int
		want   []int
	}{
		{
			name:   "push front repeat",
			slice:  []int{1, 2, 3},
			repeat: 2,
			value:  4,
			want:   []int{4, 4, 1, 2, 3},
		},
		{
			name:   "push zero front repeat",
			slice:  []int{1, 2, 3},
			repeat: 0,
			value:  4,
			want:   []int{1, 2, 3},
		},
		{
			name:   "push front repeat with container",
			slice:  newContainer(10, []int{1, 2, 3}),
			repeat: 2,
			value:  5,
			want:   []int{5, 5, 1, 2, 3},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := pushFrontRepeat(test.slice, test.value, test.repeat)
			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("pushFrontRepeat(%v, %v, %v) = %v, want %v", test.slice, test.value, test.repeat, got, test.want)
			}
		})
	}
}

func TestTake(t *testing.T) {
	tests := []struct {
		name  string
		slice []int
		idx   int
		want  []int
	}{
		{
			name:  "take from empty slice",
			slice: []int{},
			idx:   0,
			want:  []int{},
		},
		{
			name:  "take second from slice",
			slice: []int{1, 2, 3},
			idx:   1,
			want:  []int{1, 3},
		},
		{
			name:  "take last from slice",
			slice: []int{1, 2, 3},
			idx:   2,
			want:  []int{1, 2},
		},
		{
			name:  "take out of range",
			slice: []int{1, 2, 3},
			idx:   3,
			want:  []int{1, 2, 3},
		},
		{
			name:  "take negative index",
			slice: []int{1, 2, 3},
			idx:   -1,
			want:  []int{1, 2, 3},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := remove(test.slice, test.idx)
			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("take(%v, %d) = %v, want %v", test.slice, test.idx, got, test.want)
			}
		})
	}
}

func TestExtend(t *testing.T) {
	tests := []struct {
		name   string
		slice  []int
		target int
		want   int
	}{
		{
			name:   "extend with empty slice",
			slice:  []int{},
			target: 10,
			want:   10,
		},
		{
			name:   "extend with slice",
			slice:  []int{1, 2, 3},
			target: 10,
			want:   10,
		},
		{
			name:   "extend with slice",
			slice:  []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 0},
			target: 5,
			want:   10,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := extend(test.slice, test.target)
			if cap(got) != test.want {
				t.Errorf("extend(%v, %d) = %v, want %v", test.slice, test.want, got, test.want)
			}
		})
	}
}

func TestInsert(t *testing.T) {
	tests := []struct {
		name  string
		slice []int
		idx   int
		value int
		want  []int
	}{
		{
			name:  "insert into empty slice",
			slice: []int{},
			idx:   0,
			value: 1,
			want:  []int{1},
		},
		{
			name:  "insert into slice",
			slice: []int{1, 2, 3},
			idx:   1,
			value: 4,
			want:  []int{1, 4, 2, 3},
		},
		{
			name:  "insert first into slice",
			slice: []int{1, 2, 3},
			idx:   0,
			value: 6,
			want:  []int{6, 1, 2, 3},
		},
		{
			name:  "insert last into slice",
			slice: []int{1, 2, 3},
			idx:   3,
			value: 7,
			want:  []int{1, 2, 3, 7},
		},
		{
			name:  "insert out of range",
			slice: []int{1, 2, 3},
			idx:   4,
			value: 8,
			want:  []int{1, 2, 3},
		},
		{
			name:  "insert negative index",
			slice: []int{1, 2, 3},
			idx:   -1,
			value: 9,
			want:  []int{1, 2, 3},
		},
		{
			name:  "insert into empty slice with container",
			slice: newContainer(10, []int{}),
			idx:   0,
			value: 1,
			want:  []int{1},
		},
		{
			name:  "insert into slice with container",
			slice: newContainer(10, []int{1, 2, 3}),
			idx:   1,
			value: 4,
			want:  []int{1, 4, 2, 3},
		},
		{
			name:  "insert first into slice with container",
			slice: newContainer(10, []int{1, 2, 3}),
			idx:   0,
			value: 6,
			want:  []int{6, 1, 2, 3},
		},
		{
			name:  "insert last into slice with container",
			slice: newContainer(10, []int{1, 2, 3}),
			idx:   3,
			value: 7,
			want:  []int{1, 2, 3, 7},
		},
		{
			name:  "insert out of range with container",
			slice: newContainer(10, []int{1, 2, 3}),
			idx:   4,
			value: 8,
			want:  []int{1, 2, 3},
		},
		{
			name:  "insert negative index with container",
			slice: newContainer(10, []int{1, 2, 3}),
			idx:   -1,
			value: 9,
			want:  []int{1, 2, 3},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := insert(test.slice, test.idx, test.value)
			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("insert(%v, %d, %d) = %v, want %v", test.slice, test.idx, test.value, got, test.want)
			}
		})
	}
}

func TestRemoveFront(t *testing.T) {
	tests := []struct {
		name      string
		slice     []int
		count     int
		want      []int
		remainCap int
	}{
		{
			name:      "remove front from empty slice",
			slice:     []int{},
			count:     0,
			want:      []int{},
			remainCap: 0,
		},
		{
			name:      "remove front from slice",
			slice:     []int{1, 2, 3},
			count:     1,
			want:      []int{2, 3},
			remainCap: 3,
		},
		{
			name:      "remove front from slice",
			slice:     []int{1, 2, 3},
			count:     2,
			want:      []int{3},
			remainCap: 3,
		},
		{
			name:      "remove front from slice",
			slice:     []int{1, 2, 3},
			count:     3,
			want:      []int{},
			remainCap: 3,
		},
		{
			name:      "remove front from slice",
			slice:     []int{1, 2, 3},
			count:     4,
			want:      []int{},
			remainCap: 3,
		},
		{
			name:      "remove front with negative count",
			slice:     []int{1, 2, 3},
			count:     -1,
			want:      []int{1, 2, 3},
			remainCap: 3,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := trimFront(test.slice, test.count)
			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("removeFront(%v, %d) = %v, want %v", test.slice, test.count, got, test.want)
			}

			if cap(got) != test.remainCap {
				t.Errorf("removeFront(%v, %d) = %v, cap %d, want %d", test.slice, test.count, got, cap(got), test.remainCap)
			}
		})
	}
}

func TestRemoveBack(t *testing.T) {
	tests := []struct {
		name      string
		slice     []int
		count     int
		want      []int
		remainCap int
	}{
		{
			name:      "remove back from empty slice",
			slice:     []int{},
			count:     0,
			want:      []int{},
			remainCap: 0,
		},
		{
			name:      "remove back from slice",
			slice:     []int{1, 2, 3},
			count:     1,
			want:      []int{1, 2},
			remainCap: 3,
		},
		{
			name:      "remove back from slice",
			slice:     []int{1, 2, 3},
			count:     2,
			want:      []int{1},
			remainCap: 3,
		},
		{
			name:      "remove back from slice",
			slice:     []int{1, 2, 3},
			count:     3,
			want:      []int{},
			remainCap: 3,
		},
		{
			name:      "remove back from slice",
			slice:     []int{1, 2, 3},
			count:     4,
			want:      []int{},
			remainCap: 3,
		},
		{
			name:      "remove back with negative count",
			slice:     []int{1, 2, 3},
			count:     -1,
			want:      []int{1, 2, 3},
			remainCap: 3,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := trimBack(test.slice, test.count)
			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("removeBack(%v, %d) = %v, want %v", test.slice, test.count, got, test.want)
			}
		})
	}
}

func TestTrim(t *testing.T) {
	tests := []struct {
		name      string
		slice     []int
		start     int
		end       int
		want      []int
		remainCap int
	}{
		{
			name:      "trim from empty slice",
			slice:     []int{},
			start:     0,
			end:       0,
			want:      []int{},
			remainCap: 0,
		},
		{
			name:      "trim from slice",
			slice:     []int{1, 2, 3},
			start:     0,
			end:       0,
			want:      []int{},
			remainCap: 3,
		},
		{
			name:      "trim from slice",
			slice:     []int{1, 2, 3},
			start:     0,
			end:       2,
			want:      []int{1, 2},
			remainCap: 3,
		},
		{
			name:      "trim from slice",
			slice:     []int{1, 2, 3},
			start:     1,
			end:       3,
			want:      []int{2, 3},
			remainCap: 3,
		},
		{
			name:      "trim from slice",
			slice:     []int{1, 2, 3},
			start:     1,
			end:       2,
			want:      []int{2},
			remainCap: 3,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := trim(test.slice, test.start, test.end)
			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("trim(%v, %d, %d) = %v, want %v", test.slice, test.start, test.end, got, test.want)
			}

			if cap(got) != test.remainCap {
				t.Errorf("trim(%v, %d, %d) = %v, cap %d, want %d", test.slice, test.start, test.end, got, cap(got), test.remainCap)
			}
		})
	}
}

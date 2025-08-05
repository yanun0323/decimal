package decimal

import (
	"reflect"
	"testing"
)

func newContainer(caps int, values []byte) []byte {
	if caps < len(values) {
		caps = len(values)
	}

	buf := make([]byte, 0, caps)
	buf = append(buf, values...)

	return buf
}

func TestPushFront(t *testing.T) {
	tests := []struct {
		name   string
		slice  []byte
		values []byte
		want   []byte
	}{
		{
			name:   "push front",
			slice:  []byte{1, 2, 3},
			values: []byte{4, 5, 6},
			want:   []byte{4, 5, 6, 1, 2, 3},
		},
		{
			name:   "push front with empty slice",
			slice:  []byte{},
			values: []byte{4, 5, 6},
			want:   []byte{4, 5, 6},
		},
		{
			name:   "push front with empty values",
			slice:  []byte{1, 2, 3},
			values: []byte{},
			want:   []byte{1, 2, 3},
		},
		{
			name:   "push front with slice and values",
			slice:  []byte{1, 2, 3},
			values: []byte{4, 5, 6},
			want:   []byte{4, 5, 6, 1, 2, 3},
		},
		{
			name:   "push front with slice and values",
			slice:  []byte{1, 2, 3},
			values: []byte{8},
			want:   []byte{8, 1, 2, 3},
		},
		{
			name:   "push front with slice and values",
			slice:  newContainer(4, []byte{1, 2, 3}),
			values: []byte{8},
			want:   []byte{8, 1, 2, 3},
		},
		{
			name:   "push front with slice and values",
			slice:  newContainer(4, []byte{1, 2, 3}),
			values: []byte{8, 9, 10},
			want:   []byte{8, 9, 10, 1, 2, 3},
		},
		{
			name:   "push front with slice and values",
			slice:  newContainer(10, []byte{1, 2, 3}),
			values: []byte{8, 9, 10},
			want:   []byte{8, 9, 10, 1, 2, 3},
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
		slice  []byte
		repeat int
		value  byte
		want   []byte
	}{
		{
			name:   "push front repeat",
			slice:  []byte{1, 2, 3},
			repeat: 2,
			value:  4,
			want:   []byte{1, 2, 3, 4, 4},
		},
		{
			name:   "push zero front repeat",
			slice:  []byte{1, 2, 3},
			repeat: 0,
			value:  4,
			want:   []byte{1, 2, 3},
		},
		{
			name:   "push front repeat with container",
			slice:  newContainer(10, []byte{1, 2, 3}),
			repeat: 2,
			value:  5,
			want:   []byte{1, 2, 3, 5, 5},
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
	newContainer := func(caps int, values []byte) []byte {
		if caps < len(values) {
			caps = len(values)
		}

		buf := make([]byte, 0, caps)
		buf = append(buf, values...)

		return buf
	}

	tests := []struct {
		name   string
		slice  []byte
		repeat int
		value  byte
		want   []byte
	}{
		{
			name:   "push front repeat",
			slice:  []byte{1, 2, 3},
			repeat: 2,
			value:  4,
			want:   []byte{4, 4, 1, 2, 3},
		},
		{
			name:   "push zero front repeat",
			slice:  []byte{1, 2, 3},
			repeat: 0,
			value:  4,
			want:   []byte{1, 2, 3},
		},
		{
			name:   "push front repeat with container",
			slice:  newContainer(10, []byte{1, 2, 3}),
			repeat: 2,
			value:  5,
			want:   []byte{5, 5, 1, 2, 3},
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
		slice []byte
		idx   int
		want  []byte
	}{
		{
			name:  "take from empty slice",
			slice: []byte{},
			idx:   0,
			want:  []byte{},
		},
		{
			name:  "take second from slice",
			slice: []byte{1, 2, 3},
			idx:   1,
			want:  []byte{1, 3},
		},
		{
			name:  "take last from slice",
			slice: []byte{1, 2, 3},
			idx:   2,
			want:  []byte{1, 2},
		},
		{
			name:  "take out of range",
			slice: []byte{1, 2, 3},
			idx:   3,
			want:  []byte{1, 2, 3},
		},
		{
			name:  "take negative index",
			slice: []byte{1, 2, 3},
			idx:   -1,
			want:  []byte{1, 2, 3},
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
		slice  []byte
		target int
		want   int
	}{
		{
			name:   "extend with empty slice",
			slice:  []byte{},
			target: 10,
			want:   10,
		},
		{
			name:   "extend with slice",
			slice:  []byte{1, 2, 3},
			target: 10,
			want:   10,
		},
		{
			name:   "extend with slice",
			slice:  []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 0},
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
		slice []byte
		idx   int
		value byte
		want  []byte
	}{
		{
			name:  "insert into empty slice",
			slice: []byte{},
			idx:   0,
			value: 1,
			want:  []byte{1},
		},
		{
			name:  "insert into slice",
			slice: []byte{1, 2, 3},
			idx:   1,
			value: 4,
			want:  []byte{1, 4, 2, 3},
		},
		{
			name:  "insert first into slice",
			slice: []byte{1, 2, 3},
			idx:   0,
			value: 6,
			want:  []byte{6, 1, 2, 3},
		},
		{
			name:  "insert last into slice",
			slice: []byte{1, 2, 3},
			idx:   3,
			value: 7,
			want:  []byte{1, 2, 3, 7},
		},
		{
			name:  "insert out of range",
			slice: []byte{1, 2, 3},
			idx:   4,
			value: 8,
			want:  []byte{1, 2, 3},
		},
		{
			name:  "insert negative index",
			slice: []byte{1, 2, 3},
			idx:   -1,
			value: 9,
			want:  []byte{1, 2, 3},
		},
		{
			name:  "insert into empty slice with container",
			slice: newContainer(10, []byte{}),
			idx:   0,
			value: 1,
			want:  []byte{1},
		},
		{
			name:  "insert into slice with container",
			slice: newContainer(10, []byte{1, 2, 3}),
			idx:   1,
			value: 4,
			want:  []byte{1, 4, 2, 3},
		},
		{
			name:  "insert first into slice with container",
			slice: newContainer(10, []byte{1, 2, 3}),
			idx:   0,
			value: 6,
			want:  []byte{6, 1, 2, 3},
		},
		{
			name:  "insert last into slice with container",
			slice: newContainer(10, []byte{1, 2, 3}),
			idx:   3,
			value: 7,
			want:  []byte{1, 2, 3, 7},
		},
		{
			name:  "insert out of range with container",
			slice: newContainer(10, []byte{1, 2, 3}),
			idx:   4,
			value: 8,
			want:  []byte{1, 2, 3},
		},
		{
			name:  "insert negative index with container",
			slice: newContainer(10, []byte{1, 2, 3}),
			idx:   -1,
			value: 9,
			want:  []byte{1, 2, 3},
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
		slice     []byte
		count     int
		want      []byte
		remainCap int
	}{
		{
			name:      "remove front from empty slice",
			slice:     []byte{},
			count:     0,
			want:      []byte{},
			remainCap: 0,
		},
		{
			name:      "remove front from slice",
			slice:     []byte{1, 2, 3},
			count:     1,
			want:      []byte{2, 3},
			remainCap: 3,
		},
		{
			name:      "remove front from slice",
			slice:     []byte{1, 2, 3},
			count:     2,
			want:      []byte{3},
			remainCap: 3,
		},
		{
			name:      "remove front from slice",
			slice:     []byte{1, 2, 3},
			count:     3,
			want:      []byte{},
			remainCap: 3,
		},
		{
			name:      "remove front from slice",
			slice:     []byte{1, 2, 3},
			count:     4,
			want:      []byte{},
			remainCap: 3,
		},
		{
			name:      "remove front with negative count",
			slice:     []byte{1, 2, 3},
			count:     -1,
			want:      []byte{1, 2, 3},
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
		slice     []byte
		count     int
		want      []byte
		remainCap int
	}{
		{
			name:      "remove back from empty slice",
			slice:     []byte{},
			count:     0,
			want:      []byte{},
			remainCap: 0,
		},
		{
			name:      "remove back from slice",
			slice:     []byte{1, 2, 3},
			count:     1,
			want:      []byte{1, 2},
			remainCap: 3,
		},
		{
			name:      "remove back from slice",
			slice:     []byte{1, 2, 3},
			count:     2,
			want:      []byte{1},
			remainCap: 3,
		},
		{
			name:      "remove back from slice",
			slice:     []byte{1, 2, 3},
			count:     3,
			want:      []byte{},
			remainCap: 3,
		},
		{
			name:      "remove back from slice",
			slice:     []byte{1, 2, 3},
			count:     4,
			want:      []byte{},
			remainCap: 3,
		},
		{
			name:      "remove back with negative count",
			slice:     []byte{1, 2, 3},
			count:     -1,
			want:      []byte{1, 2, 3},
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
		slice     []byte
		start     int
		end       int
		want      []byte
		remainCap int
	}{
		{
			name:      "trim from empty slice",
			slice:     []byte{},
			start:     0,
			end:       0,
			want:      []byte{},
			remainCap: 0,
		},
		{
			name:      "trim from slice",
			slice:     []byte{1, 2, 3},
			start:     0,
			end:       0,
			want:      []byte{},
			remainCap: 3,
		},
		{
			name:      "trim from slice",
			slice:     []byte{1, 2, 3},
			start:     0,
			end:       2,
			want:      []byte{1, 2},
			remainCap: 3,
		},
		{
			name:      "trim from slice",
			slice:     []byte{1, 2, 3},
			start:     1,
			end:       3,
			want:      []byte{2, 3},
			remainCap: 3,
		},
		{
			name:      "trim from slice",
			slice:     []byte{1, 2, 3},
			start:     1,
			end:       2,
			want:      []byte{2},
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

package slices

import (
	"fmt"
	"testing"
)

func expectPanic(t *testing.T, funcname string, f func()) {
	t.Helper()
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("%s did not panic", funcname)
		}
	}()

	f()
}

func equals[T comparable](t *testing.T, have, want []T) bool {
	t.Helper()
	if len(have) != len(want) {
		t.Errorf("lengths differ: %d != %d", len(have), len(want))
		return false
	}
	for i, v := range have {
		if v != want[i] {
			t.Errorf("values differ at index %d: have %v != want %v", i, v, want[i])
			return false
		}
	}
	return true
}

func fuzzyEquals[T comparable](t *testing.T, have, want []T) bool {
	t.Helper()
	if len(have) != len(want) {
		t.Errorf("lengths differ: %d != %d", len(have), len(want))
		return false
	}
	elementsHave := make(map[T]bool)
	for _, v := range have {
		elementsHave[v] = true
	}
	for _, v := range want {
		if _, ok := elementsHave[v]; !ok {
			t.Errorf("value %v (part of want) not found in %v", v, have)
			return false
		}
		delete(elementsHave, v)
	}
	if len(elementsHave) > 0 {
		for v := range elementsHave {
			t.Errorf("value %v (part of have) not found in %v", v, want)
		}
		return false
	}
	return true
}

func TestRemove(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}

	removed := Remove(slice, 1)
	if !fuzzyEquals(t, removed, []int{2, 3, 4, 5}) {
		t.Errorf("Remove(%v, %v) = %v, want %v", slice, 1, removed, []int{2, 3, 4, 5})
	}

	removed = Remove(slice, 5)
	if !fuzzyEquals(t, removed, []int{1, 2, 3, 4}) {
		t.Errorf("Remove(%v, %v) = %v, want %v", slice, 5, removed, []int{1, 2, 3, 4})
	}

	removed = Remove(slice, 6)
	if !equals(t, removed, slice) {
		t.Errorf("Remove(%v, %v) = %v, want %v", slice, 6, removed, slice)
	}
}

func TestRemoveAt(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}

	removed := RemoveAt(slice, 1)
	if !fuzzyEquals(t, removed, []int{1, 3, 4, 5}) {
		t.Errorf("RemoveAt(%v, %v) = %v, want %v", slice, 1, removed, []int{1, 3, 4, 5})
	}

	expectPanic(t, fmt.Sprintf("RemoveAt(%v, %v)", slice, 6), func() {
		RemoveAt(slice, 6)
	})
}

func TestFilter(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}

	filtered := Filter(slice, func(i int) bool { return i%2 == 0 })

	if !fuzzyEquals(t, filtered, []int{2, 4}) {
		t.Errorf("Filter(%v, %v) = %v, want %v", slice, "i%2==0", filtered, []int{2, 4})
	}

	filtered = Filter(slice, func(i int) bool { return i%2 == 1 })
	if !fuzzyEquals(t, filtered, []int{1, 3, 5}) {
		t.Errorf("Filter(%v, %v) = %v, want %v", slice, "i%2==1", filtered, []int{1, 3, 5})
	}

	filtered = Filter(slice, func(i int) bool { return i > 3 })
	if !fuzzyEquals(t, filtered, []int{4, 5}) {
		t.Errorf("Filter(%v, %v) = %v, want %v", slice, "i>3", filtered, []int{4, 5})
	}
}

func TestMap(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}

	mapped := Map(slice, func(i int) int { return i * 2 })
	if !equals(t, mapped, []int{2, 4, 6, 8, 10}) {
		t.Errorf("Map(%v, %v) = %v, want %v", slice, "i*2", mapped, []int{2, 4, 6, 8, 10})
	}

	mapped = Map(slice, func(i int) int { return 0 })
	if !equals(t, mapped, []int{0, 0, 0, 0, 0}) {
		t.Errorf("Map(%v, %v) = %v, want %v", slice, "0", mapped, []int{0, 0, 0, 0, 0})
	}
}

func TestContains(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}

	if !Contains(slice, 1) {
		t.Errorf("Contains(%v, %v) = false, want true", slice, 1)
	}

	if !Contains(slice, 5) {
		t.Errorf("Contains(%v, %v) = false, want true", slice, 5)
	}

	if Contains(slice, 6) {
		t.Errorf("Contains(%v, %v) = true, want false", slice, 6)
	}
}

func TestIndexOf(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}

	if IndexOf(slice, 1) != 0 {
		t.Errorf("IndexOf(%v, %v) = %v, want %v", slice, 1, IndexOf(slice, 1), 0)
	}

	if IndexOf(slice, 5) != 4 {
		t.Errorf("IndexOf(%v, %v) = %v, want %v", slice, 5, IndexOf(slice, 5), 4)
	}

	if IndexOf(slice, 6) != -1 {
		t.Errorf("IndexOf(%v, %v) = %v, want %v", slice, 6, IndexOf(slice, 6), -1)
	}
}

func TestMax(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}

	if Max(slice) != 5 {
		t.Errorf("Max(%v) = %v, want %v", slice, Max(slice), 5)
	}

	slice = []int{5, 4, 3, 2, 1}
	if Max(slice) != 5 {
		t.Errorf("Max(%v) = %v, want %v", slice, Max(slice), 5)
	}
}

func TestMin(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}

	if Min(slice) != 1 {
		t.Errorf("Min(%v) = %v, want %v", slice, Min(slice), 1)
	}

	slice = []int{5, 4, 3, 2, 1}
	if Min(slice) != 1 {
		t.Errorf("Min(%v) = %v, want %v", slice, Min(slice), 1)
	}
}

func TestAverage(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}

	if Average(slice) != 3 {
		t.Errorf("Average(%v) = %v, want %v", slice, Average(slice), 3)
	}

	slice = []int{5, 4, 3, 2, 1}
	if Average(slice) != 3 {
		t.Errorf("Average(%v) = %v, want %v", slice, Average(slice), 3)
	}

	fslice := []float32{1.0, 2.0, 3.0, 4.0, 5.0}
	if Average(fslice) != 3.0 {
		t.Errorf("Average(%v) = %v, want %v", slice, Average(slice), 3.0)
	}

	fslice = []float32{5.0, 4.0, 3.0, 2.0, 1.0}
	if Average(fslice) != 3.0 {
		t.Errorf("Average(%v) = %v, want %v", slice, Average(slice), 3.0)
	}
}

func TestUnique(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5, 1, 2, 3, 4, 5}

	unique := Unique(slice)
	if !equals(t, unique, []int{1, 2, 3, 4, 5}) {
		t.Errorf("Unique(%v) = %v, want %v", slice, unique, []int{1, 2, 3, 4, 5})
	}

	slice = []int{1, 1, 1, 1, 1, 1, 1, 1, 1, 1}
	unique = Unique(slice)
	if !equals(t, unique, []int{1}) {
		t.Errorf("Unique(%v) = %v, want %v", slice, unique, []int{1})
	}

	slice = []int{}
	unique = Unique(slice)
	if !equals(t, unique, []int{}) {
		t.Errorf("Unique(%v) = %v, want %v", slice, unique, []int{})
	}
}

func TestUniqueBy(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5, 1, 2, 3, 4, 5}

	unique := UniqueBy(slice, func(i int) int { return i })
	if !equals(t, unique, []int{1, 2, 3, 4, 5}) {
		t.Errorf("UniqueBy(%v, %v) = %v, want %v", slice, "i", unique, []int{1, 2, 3, 4, 5})
	}

	type user struct {
		id int
	}

	slice2 := []user{{1}, {2}, {3}, {4}, {5}, {1}, {2}, {3}, {4}, {5}}
	unique2 := UniqueBy(slice2, func(i user) int { return i.id })
	for i, u := range unique2 {
		if u.id != i+1 {
			t.Errorf("UniqueBy(%v, %v) = %v, want %v", slice2, "i.id", unique2, []user{{1}, {2}, {3}, {4}, {5}})
		}
	}
}

func TestReduce(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}

	reduced := Reduce(slice, 0, func(a, b int) int { return a + b })
	if reduced != 15 {
		t.Errorf("Reduce(%v, %v) = %v, want %v", slice, "a+b", reduced, 15)
	}

	reduced = Reduce(slice, 1, func(a, b int) int { return a * b })
	if reduced != 120 {
		t.Errorf("Reduce(%v, %v) = %v, want %v", slice, "a*b", reduced, 120)
	}

	reduced = Reduce([]int{}, 0, func(a, b int) int { return a + b })
	if reduced != 0 {
		t.Errorf("Reduce(%v, %v) = %v, want %v", slice, "a+b", reduced, 0)
	}
}

func TestSome(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}

	if !Some(slice, func(i int) bool { return i > 3 }) {
		t.Errorf("Some(%v, %v) = %v, want %v", slice, "i > 3", false, true)
	}

	if Some(slice, func(i int) bool { return i > 6 }) {
		t.Errorf("Some(%v, %v) = %v, want %v", slice, "i > 6", true, false)
	}

	if Some([]int{}, func(i int) bool { return i > 6 }) {
		t.Errorf("Some(%v, %v) = %v, want %v", slice, "i > 6", true, false)
	}
}

func TestEvery(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}

	if !Every(slice, func(i int) bool { return i > 0 }) {
		t.Errorf("Every(%v, %v) = %v, want %v", slice, "i > 0", false, true)
	}

	if Every(slice, func(i int) bool { return i > 3 }) {
		t.Errorf("Every(%v, %v) = %v, want %v", slice, "i > 3", true, false)
	}

	if !Every([]int{}, func(i int) bool { return i > 6 }) {
		t.Errorf("Every(%v, %v) = %v, want %v", slice, "i > 6", false, true)
	}
}

func TestAssociate(t *testing.T) {
	k := []string{"a", "b", "c"}
	v := []int{1, 2, 3}

	assoc := Associate(k, v)

	if val, ok := assoc["a"]; !ok || val != 1 {
		t.Errorf("Associate(%v, %v)[%v] = %v, want %v", k, v, "a", val, 1)
	}

	assoc = Associate(k, []int{})
	if len(assoc) != 0 {
		t.Errorf("Associate(%v, %v) should be empty, got %v", k, v, assoc)
	}

	assoc = Associate([]string{}, []int{1, 2, 3})
	if len(assoc) != 0 {
		t.Errorf("Associate(%v, %v) should be empty, got %v", k, v, assoc)
	}
}

func TestAssociateBy(t *testing.T) {
	k := []string{"a", "b", "c"}

	assoc := AssociateBy(k, func(k string) int { return int(k[0]) })

	if val, ok := assoc["a"]; !ok || val != int('a') {
		t.Errorf("AssociateBy(%v, %v)[%v] = %v, want %v", k, "int(k[0])", "a", val, int('a'))
	}

	assoc = AssociateBy([]string{}, func(k string) int { return int(k[0]) })
	if len(assoc) != 0 {
		t.Errorf("AssociateBy(%v, %v) should be empty, got %v", k, "int(k[0])", assoc)
	}
}

func TestAssociateReverseBy(t *testing.T) {
	v := []string{"a", "b", "c"}

	assoc := AssociateReverseBy(v, func(k string) int { return int(k[0]) })

	if val, ok := assoc[int('a')]; !ok || val != "a" {
		t.Errorf("AssociateReverseBy(%v, %v)[%v] = %v, want %v", v, "int(k[0])", "a", val, "a")
	}

	assoc = AssociateReverseBy([]string{}, func(k string) int { return int(k[0]) })
	if len(assoc) != 0 {
		t.Errorf("AssociateReverseBy(%v, %v) should be empty, got %v", v, "int(k[0])", assoc)
	}
}

func TestIntersperse(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}

	interspersed := Intersperse(slice, 0)
	if !equals(t, interspersed, []int{1, 0, 2, 0, 3, 0, 4, 0, 5}) {
		t.Errorf("Intersperse(%v, %v) = %v, want %v", slice, 0, interspersed, []int{1, 0, 2, 0, 3, 0, 4, 0, 5})
	}

	interspersed = Intersperse([]int{}, 0)
	if !equals(t, interspersed, []int{}) {
		t.Errorf("Intersperse(%v, %v) = %v, want %v", slice, 0, interspersed, []int{})
	}
}

func TestIntersperseBy(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}

	interspersed := IntersperseBy(slice, func(i int) int { return i })
	if !equals(t, interspersed, []int{1, 1, 2, 2, 3, 3, 4, 4, 5}) {
		t.Errorf("IntersperseBy(%v, %v) = %v, want %v", slice, "i", interspersed, []int{1, 1, 2, 2, 3, 3, 4, 4, 5})
	}

	interspersed = IntersperseByIndex([]int{}, func(i int) int { return i })
	if !equals(t, interspersed, []int{}) {
		t.Errorf("IntersperseBy(%v, %v) = %v, want %v", slice, "i", interspersed, []int{})
	}
}

func TestIntersperseByIndex(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}

	interspersed := IntersperseByIndex(slice, func(i int) int { return i })
	if !equals(t, interspersed, []int{1, 0, 2, 1, 3, 2, 4, 3, 5}) {
		t.Errorf("IntersperseBy(%v, %v) = %v, want %v", slice, "i", interspersed, []int{1, 0, 2, 1, 3, 2, 4, 3, 5})
	}

	interspersed = IntersperseByIndex([]int{}, func(i int) int { return i })
	if !equals(t, interspersed, []int{}) {
		t.Errorf("IntersperseBy(%v, %v) = %v, want %v", slice, "i", interspersed, []int{})
	}
}

func TestShuffle(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}

	shuffled := copy(slice)
	Shuffle(shuffled)

	if !fuzzyEquals(t, shuffled, slice) {
		t.Errorf("Shuffle(%v) = %v, want %v", slice, shuffled, slice)
	}
}

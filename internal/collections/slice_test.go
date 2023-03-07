package collections

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilter(t *testing.T) {
	expected := []int{1, 2, 3}
	actual := Filter([]int{1, 2, 3, 4, 5, 6}, func(x int) bool { return x <= 3 })
	assert.ElementsMatch(t, expected, actual)
	expected1 := []string{"test"}
	actual1 := Filter([]string{"asdasd", "123123", "test", "test123123"}, func(x string) bool { return x == "test" })
	assert.ElementsMatch(t, expected1, actual1)
}

func TestFirst(t *testing.T) {
	expected := 1
	actual := FirstMatch([]int{-3, -2, -1, 0, 1, 2, 3}, func(x int) bool { return x > 0 })
	assert.Equal(t, expected, actual)
	expected = 0
	actual = FirstMatch([]int{1, 2, 3}, func(x int) bool { return x < 0 })
	assert.Equal(t, expected, actual)
}

func TestLast(t *testing.T) {
	expected := 3
	actual := LastMatch([]int{-3, -2, -1, 0, 1, 2, 3}, func(x int) bool { return x > 0 })
	assert.Equal(t, expected, actual)
	expected = 0
	actual = LastMatch([]int{1, 2, 3}, func(x int) bool { return x < 0 })
	assert.Equal(t, expected, actual)
}

func TestAny(t *testing.T) {
	assert.True(t, AnyMatch([]int{1, 1, 2, 3, 5, 8, 13}, func(x int) bool { return x > 10 }))
	assert.False(t, AnyMatch([]int{1, 1, 2, 3, 5, 8, 13}, func(x int) bool { return x < 0 }))
}

func TestAll(t *testing.T) {
	assert.True(t, AllMatch([]int{1, 2, 3}, func(x int) bool { return x > 0 }))
	assert.False(t, AllMatch([]int{-1, 0, 1, 2, 3}, func(x int) bool { return x > 0 }))
}

func TestMap(t *testing.T) {
	assert.ElementsMatch(t,
		[]int{2, 4, 6, 8},
		Map([]int{1, 2, 3, 4}, func(x int) int { return 2 * x }),
	)
	assert.ElementsMatch(t,
		[]int{4, 1, 6},
		Map([]string{"test", "a", "golang"}, func(x string) int { return len(x) }),
	)
}

func TestReduce(t *testing.T) {
	assert.Equal(t, true, Reduce([]bool{true, true, true, true}, func(acc bool, elem bool) bool { return acc && elem }, true))
	assert.Equal(t, 24, Reduce([]int{1, 2, 3, 4}, func(acc int, elem int) int { return acc * elem }, 1))
}

func TestIndexOf(t *testing.T) {
	assert.Equal(t, 2, IndexOf([]int{1, 2, 3, 4, 5}, 3))
	assert.Equal(t, -1, IndexOf([]int{1, 2, 3, 4, 5}, 6))
}

func TestContains(t *testing.T) {
	assert.True(t, Contains([]int{1, 2, 3, 4}, 3))
	assert.False(t, Contains([]string{"a", "b", "test"}, "wednesday"))
}

type person struct {
	ID        int
	Name      string
	ManagerID int
}

func TestGroupBy(t *testing.T) {
	alice := person{1, "Alice", 0}
	bob := person{2, "Bob", 1}
	carol := person{9, "Carol", 1}
	david := person{11, "David", 9}
	eugene := person{13, "Eugene", 2}
	francene := person{12, "Francene", 9}
	employees := []person{alice, bob, carol, david, eugene, francene}

	assert.Equal(t,
		map[int][]person{
			0: {alice},
			1: {bob, carol},
			2: {eugene},
			9: {david, francene},
		},
		GroupBy(employees, func(p person) int { return p.ManagerID }, func(p person) person { return p }),
	)
}

func TestToSet(t *testing.T) {
	people := []person{
		{
			ID:   1,
			Name: "Alice",
		},
		{
			ID:   2,
			Name: "Bob",
		},
		{
			ID:   9,
			Name: "Carol",
		},
	}

	assert.Equal(t,
		map[int]struct{}{1: {}, 2: {}, 9: {}},
		ToSet(people, func(p person) int { return p.ID }),
	)

	assert.Equal(t,
		map[string]struct{}{
			"Alice": {},
			"Bob":   {},
			"Carol": {},
		},
		ToSet(people, func(p person) string { return p.Name }),
	)
}

func TestToMap(t *testing.T) {
	alice := person{1, "Alice", 0}
	bob := person{2, "Bob", 1}
	carol := person{9, "Carol", 1}
	people := []person{alice, bob, carol}

	assert.Equal(t,
		map[int]person{
			1: alice,
			2: bob,
			9: carol,
		},
		ToMap(people, func(p person) int { return p.ID }, func(p person) person { return p }),
	)

	assert.Equal(t,
		map[string]person{
			"Alice": alice,
			"Bob":   bob,
			"Carol": carol,
		},
		ToMap(people, func(p person) string { return p.Name }, func(p person) person { return p }),
	)
}

func TestUniq(t *testing.T) {
	assert.ElementsMatch(t, []int{22, 31, 12}, Uniq([]int{22, 22, 31, 12}, func(num int) int { return num }))
	assert.ElementsMatch(t, []string{"a", "b", "c"}, Uniq([]string{"a", "b", "c", "a"}, func(str string) string { return str }))
}

// BenchmarkFilter
// goos: linux
// goarch: amd64
// pkg: git.woa.com/tencent_cloud_mobile_tools/QAPM_CLOUD/go-between/library/collections/slices
// cpu: AMD EPYC 7K62 48-Core Processor
// BenchmarkFilter-16    	 9866012	       122.2 ns/op	      56 B/op	       3 allocs/op
// BenchmarkFilter
//
//	@param b *testing.B
//	@author: Kevineluo 2022-11-22 07:12:14
func BenchmarkFilter(b *testing.B) {
	nums := []int{1, 2, 3, 4, 5, 6}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Filter(nums, func(x int) bool { return x <= 3 })
	}
}

// BenchmarkFirst
// goos: linux
// goarch: amd64
// pkg: git.woa.com/tencent_cloud_mobile_tools/QAPM_CLOUD/go-between/library/collections/slices
// cpu: AMD EPYC 7K62 48-Core Processor
// BenchmarkFirst
// BenchmarkFirst-16    	100000000	        11.12 ns/op	       0 B/op	       0 allocs/op
//
//	@param b *testing.B
//	@author: Kevineluo 2022-11-22 07:14:26
func BenchmarkFirst(b *testing.B) {
	nums := []int{-3, -2, -1, 0, 1, 2, 3}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		FirstMatch(nums, func(x int) bool { return x > 0 })
	}
}

// BenchmarkLastMatch
// goos: linux
// goarch: amd64
// pkg: git.woa.com/tencent_cloud_mobile_tools/QAPM_CLOUD/go-between/library/collections/slices
// cpu: AMD EPYC 7K62 48-Core Processor
// BenchmarkLastMatch
// BenchmarkLastMatch-16    	128857670	         9.219 ns/op	       0 B/op	       0 allocs/op
//
//	@param b *testing.B
//	@author: Kevineluo 2022-11-22 07:16:03
func BenchmarkLastMatch(b *testing.B) {
	nums := []int{-3, -2, -1, -1, 1, 2, 3}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		LastMatch(nums, func(x int) bool { return x < 0 })
	}
}

// BenchmarkAnyMatch
// goos: linux
// goarch: amd64
// pkg: git.woa.com/tencent_cloud_mobile_tools/QAPM_CLOUD/go-between/library/collections/slices
// cpu: AMD EPYC 7K62 48-Core Processor
// BenchmarkAnyMatch
// BenchmarkAnyMatch-16    	77031542	        15.50 ns/op	       0 B/op	       0 allocs/op
//
//	@param b *testing.B
//	@author: Kevineluo 2022-11-23 10:22:01
func BenchmarkAnyMatch(b *testing.B) {
	nums := []int{1, 1, 2, 3, 5, 8, 13}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		AnyMatch(nums, func(x int) bool { return x > 10 })
	}
}

// BenchmarkAllMatch
// goos: linux
// goarch: amd64
// pkg: git.woa.com/tencent_cloud_mobile_tools/QAPM_CLOUD/go-between/library/collections/slices
// cpu: AMD EPYC 7K62 48-Core Processor
// BenchmarkAllMatch
// BenchmarkAllMatch-16    	352560135	         3.417 ns/op	       0 B/op	       0 allocs/op
//
//	@param b *testing.B
//	@author: Kevineluo 2022-11-23 11:01:39
func BenchmarkAllMatch(b *testing.B) {
	nums := []int{1, 1, 2, 3, 5, 8, 13}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		AllMatch(nums, func(x int) bool { return x > 5 })
	}
}

// BenchmarkMap
// goos: linux
// goarch: amd64
// pkg: git.woa.com/tencent_cloud_mobile_tools/QAPM_CLOUD/go-between/library/collections/slices
// cpu: AMD EPYC 7K62 48-Core Processor
// BenchmarkMap
// BenchmarkMap-16    	25056819	        48.83 ns/op	      64 B/op	       1 allocs/op
//
//	@param b *testing.B
//	@author: Kevineluo 2022-11-23 11:04:52
func BenchmarkMap(b *testing.B) {
	nums := []int{1, 1, 2, 3, 5, 8, 13}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Map(nums, func(x int) int { return x ^ 5 })
	}
}

// BenchmarkReduce
// goos: linux
// goarch: amd64
// pkg: git.woa.com/tencent_cloud_mobile_tools/QAPM_CLOUD/go-between/library/collections/slices
// cpu: AMD EPYC 7K62 48-Core Processor
// BenchmarkReduce
// BenchmarkReduce-16    	71148554	        15.80 ns/op	       0 B/op	       0 allocs/op
//
//	@param b *testing.B
//	@author: Kevineluo 2022-11-23 11:07:49
func BenchmarkReduce(b *testing.B) {
	nums := []int{1, 1, 2, 3, 5, 8, 13}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Reduce(nums, func(x, y int) int { return x * y }, 1)
	}
}

// BenchmarkIndexOf
// goos: linux
// goarch: amd64
// pkg: git.woa.com/tencent_cloud_mobile_tools/QAPM_CLOUD/go-between/library/collections/slices
// cpu: AMD EPYC 7K62 48-Core Processor
// BenchmarkIndexOf
// BenchmarkIndexOf-16    	255968445	         4.511 ns/op	       0 B/op	       0 allocs/op
//
//	@param b *testing.B
//	@author: Kevineluo 2022-11-23 11:15:00
func BenchmarkIndexOf(b *testing.B) {
	nums := []int{1, 1, 2, 3, 5, 8, 13}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		IndexOf(nums, 8)
	}
}

// BenchmarkContains
// goos: linux
// goarch: amd64
// pkg: git.woa.com/tencent_cloud_mobile_tools/QAPM_CLOUD/go-between/library/collections/slices
// cpu: AMD EPYC 7K62 48-Core Processor
// BenchmarkContains
// BenchmarkContains-16    	260143281	         5.038 ns/op	       0 B/op	       0 allocs/op
//
//	@param b *testing.B
//	@author: Kevineluo 2022-11-23 11:15:54
func BenchmarkContains(b *testing.B) {
	nums := []int{1, 1, 2, 3, 5, 8, 13}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Contains(nums, 8)
	}
}

// BenchmarkGroupBy
// goos: linux
// goarch: amd64
// pkg: git.woa.com/tencent_cloud_mobile_tools/QAPM_CLOUD/go-between/library/collections/slices
// cpu: AMD EPYC 7K62 48-Core Processor
// BenchmarkGroupBy
// BenchmarkGroupBy-16    	 1845642	       641.8 ns/op	     592 B/op	       8 allocs/op
//
//	@param b *testing.B
//	@author: Kevineluo 2022-11-23 11:17:45
func BenchmarkGroupBy(b *testing.B) {
	alice := person{1, "Alice", 0}
	bob := person{2, "Bob", 1}
	carol := person{9, "Carol", 1}
	david := person{11, "David", 9}
	eugene := person{13, "Eugene", 2}
	francene := person{12, "Francene", 9}
	employees := []person{alice, bob, carol, david, eugene, francene}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		GroupBy(employees, func(p person) int { return p.ManagerID }, func(p person) person { return p })
	}
}

// BenchmarkToSet
// goos: linux
// goarch: amd64
// pkg: git.woa.com/tencent_cloud_mobile_tools/QAPM_CLOUD/go-between/library/collections/slices
// cpu: AMD EPYC 7K62 48-Core Processor
// BenchmarkToSet
// BenchmarkToSet-16    	 8366940	       141.5 ns/op	     128 B/op	       2 allocs/op
//
//	@param b *testing.B
//	@author: Kevineluo 2022-11-23 11:19:28
func BenchmarkToSet(b *testing.B) {
	people := []person{
		{
			ID:   1,
			Name: "Alice",
		},
		{
			ID:   2,
			Name: "Bob",
		},
		{
			ID:   9,
			Name: "Carol",
		},
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ToSet(people, func(p person) int { return p.ID })
	}
}

// BenchmarkToMap
// goos: linux
// goarch: amd64
// pkg: git.woa.com/tencent_cloud_mobile_tools/QAPM_CLOUD/go-between/library/collections/slices
// cpu: AMD EPYC 7K62 48-Core Processor
// BenchmarkToMap
// BenchmarkToMap-16    	 5085471	       231.7 ns/op	     464 B/op	       2 allocs/op
//
//	@param b *testing.B
//	@author: Kevineluo 2022-11-23 11:20:22
func BenchmarkToMap(b *testing.B) {
	alice := person{1, "Alice", 0}
	bob := person{2, "Bob", 1}
	carol := person{9, "Carol", 1}
	people := []person{alice, bob, carol}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ToMap(people, func(p person) string { return p.Name }, func(p person) person { return p })
	}
}

// BenchmarkUniq
// goos: linux
// goarch: amd64
// pkg: git.woa.com/tencent_cloud_mobile_tools/QAPM_CLOUD/go-between/library/collections/slices
// cpu: AMD EPYC 7K62 48-Core Processor
// BenchmarkUniq
// BenchmarkUniq-16    	 3330757	       358.9 ns/op	     120 B/op	       4 allocs/op
//
//	@param b *testing.B
//	@author: Kevineluo 2022-11-23 11:21:35
func BenchmarkUniq(b *testing.B) {
	nums := []int{1, 1, 2, 3, 5, 8, 13}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Uniq(nums, func(num int) int { return num })
	}
}

package collections

import (
	"fmt"
	"unsafe"
)

// header returns pointer, length, and capacity of a slice.
func header[T any](s []T) (ptr uintptr, ln, cp int) {
	data := unsafe.SliceData(s)
	return uintptr(unsafe.Pointer(data)), len(s), cap(s)
}

func addr0[T any](s []T) uintptr {
	if len(s) == 0 {
		return 0
	}
	return uintptr(unsafe.Pointer(&s[0]))
}

func logSlice[T any](label string, s []T) {
	ptr, ln, cp := header(s)
	fmt.Printf("%-28s len=%-3d cap=%-3d ptr=0x%X addr0=0x%X data=%v\n", label, ln, cp, ptr, addr0(s), s)

}

func Slice() {
	// init
	fmt.Println("===What a slice really is (pointer, len, cap)===")
	var a []int
	logSlice("nil slice", a)

	a = make([]int, 0)
	logSlice("make([int, 0])", a)

	if a != nil {
		fmt.Println("it's not nil")
	}

	a = make([]int, 3, 5)
	a[0], a[1], a[2] = 10, 20, 30

	logSlice("make([]int, 3, 5)", a)

	// append
	fmt.Println("\n===Append and reallocation of backing array===")
	s := make([]int, 0, 2)
	prevPtr, _, _ := header(s)
	logSlice("start", s)
	for i := 1; i <= 8; i++ {
		s = append(s, i)
		p, ln, cp := header(s)
		if p != prevPtr {
			fmt.Printf("  -> reallocated: len=%d cap=%d ptr 0x%X -> 0x%X\n", ln, cp, prevPtr, p)
			prevPtr = p
		}
	}
	logSlice("final after appends", s)

	// Preallocation
	fmt.Println("\n===Preallocation avoids reallocations===")
	s2 := make([]int, 0, 8) // Preallocated enough capacity
	prevPtr, _, _ = header(s2)
	logSlice("start prealloc", s2)
	for i := 1; i <= 8; i++ {
		s2 = append(s2, i)
		p, _, _ := header(s2)
		if p != prevPtr {
			fmt.Printf(" -> (surpirse) reallocated! ptr 0x%X -> 0x%X\n", prevPtr, p)
			prevPtr = p
		}
	}
	logSlice("final no reallocation", s2)

	// backing array
	fmt.Println("\n===Sub-slices SHARE the same backing array===")
	base := []int{1, 2, 3, 4, 5}
	logSlice("base", base)
	sub := base[1:4]
	logSlice("sub := base[1:4]", sub)

	sub[1] = 99 //change base as well
	logSlice("sub after sub[1] = 99", sub)
	logSlice("base reflected", base)

	fmt.Println("\nAppend to sub (might affect base if cap shared)")
	logSlice("before append sub", sub)
	sub = append(sub, 777)
	logSlice("after append sub", sub)
	logSlice("base after append sub", base)

	fmt.Println("\nForce reallocation (cap exhausted)")
	t := base[:len(base):len(base)] // full slice trick: cap = len
	logSlice("t := base[:len:len]", t)
	t = append(t, 888) // new allocation
	logSlice("t after append 888", t)
	logSlice("base unchanged", base)

	fmt.Println("\n===Prepend (insert at start) is O(n) ===")
	u := []int{10, 20, 30, 40}
	logSlice("u original", u)
	u = prepend(u, 0)
	logSlice("u after prepend(0)", u)

	fmt.Println("\n===Insert in the middle is also O(n) ===")
	v := []int{1, 2, 3, 4, 5}
	logSlice("v original", v)
	v = insertAt(v, 999, 2)
	logSlice("v after insertAt 2", v)

	fmt.Println("\n===Show raw bytes of the backing array ===")
	ints := []int32{1, 2, 3, 4}
	logSlice("ints (int32)", ints)
	printBackingBytes(ints)
}

func prepend[T any](s []T, x T) []T {

	s = append(s, x) // this append is only to doble the cap
	copy(s[1:], s[0:len(s)-1])
	s[0] = x
	return s
}

func insertAt[T any](s []T, x T, i int) []T {
	if i < 0 && i > len(s) {
		panic("index out of range")
	}

	s = append(s, x)
	copy(s[i+1:], s[i:len(s)-1])
	s[i] = x

	return s
}

// printBackingBytes shows raw memory of []int32 backing array.
func printBackingBytes(s []int32) {
	if len(s) == 0 {
		fmt.Println("(empty slice)")
		return
	}
	p, ln, _ := header(s)
	total := ln * int(unsafe.Sizeof(s[0]))
	fmt.Printf("Backing array bytes (len=%d * %dB): ", ln, unsafe.Sizeof(s[0]))
	for i := 0; i < total; i++ {
		b := *(*byte)(unsafe.Pointer(p + uintptr(i)))
		fmt.Printf("%02X ", b)
	}
	fmt.Println()
}

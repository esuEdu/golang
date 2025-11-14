package collections

import (
	"fmt"
	"unsafe"
)

func Array() {
	fmt.Println("=== 1. Array allocation and initialization ===")
	var arr [4]int32 // zero-initialized
	fmt.Println("Before initialization:", arr)

	// Initialize array
	arr = [4]int32{1, 2, 3, 4}
	fmt.Println("After initialization:", arr)

	fmt.Printf("Total size of array: %d bytes\n", unsafe.Sizeof(arr))
	fmt.Printf("Each element size: %d bytes\n\n", unsafe.Sizeof(arr[0]))

	// Show memory addresses
	fmt.Println("Memory layout (addresses are contiguous):")
	for i := range arr {
		fmt.Printf("  arr[%d] = %d, address = %p\n", i, arr[i], &arr[i])
	}

	// === 2. Offsets and address math ===
	fmt.Println("\n=== 2. Pointer arithmetic (offset calculation) ===")
	base := uintptr(unsafe.Pointer(&arr[0]))
	for i := 0; i < len(arr); i++ {
		offset := uintptr(i) * unsafe.Sizeof(arr[0])
		addr := base + offset
		fmt.Printf("arr[%d] offset: %d bytes, computed address: 0x%X\n", i, offset, addr)
	}

	// === 3. Access and modify in O(1) ===
	fmt.Println("\n=== 3. Constant-time access and modification ===")
	fmt.Printf("arr[2] before = %d\n", arr[2])
	arr[2] = 99
	fmt.Printf("arr[2] after  = %d\n", arr[2])

	// === 4. Attempting to resize (impossible for arrays) ===
	fmt.Println("\n=== 4. Arrays are fixed-size ===")

	fmt.Println("To 'resize', we must create a new array and copy data:")

	newArr := [5]int32{}
	copy(newArr[:], arr[:])
	newArr[4] = 5
	fmt.Println("New array:", newArr)

	// === 5. Raw bytes in memory ===
	fmt.Println("\n===5. Raw memory bytes (hex view) ===")
	ptr := uintptr(unsafe.Pointer(&arr[0]))
	totalBytes := len(arr) * int(unsafe.Sizeof(arr[0]))

	for i := 0; i < totalBytes; i++ {
		b := *(*byte)(unsafe.Pointer(ptr + uintptr(i)))
		fmt.Printf("%02X ", b)
	}
	fmt.Println()

	fmt.Println("\n=== End of array demonstration ===")
}

// Package lookup3 implements the lookup3 hash functions designed by Bob Jenkins.
// Reference: http://burtleburtle.net/bob/c/lookup3.c
package lookup3

import (
	"encoding/binary"
	"fmt"
)

func rot(n uint32, k uint8) uint32 {
	return (n << k) | (n >> (32 - k))
}

// Mix mixes three 32-bit values (a, b, c) reversibly.
func Mix(pa, pb, pc *uint32) {
	a, b, c := *pa, *pb, *pc
	a -= c
	a ^= rot(c, 4)
	c += b
	b -= a
	b ^= rot(a, 6)
	a += c
	c -= b
	c ^= rot(b, 8)
	b += a
	a -= c
	a ^= rot(c, 16)
	c += b
	b -= a
	b ^= rot(a, 19)
	a += c
	c -= b
	c ^= rot(b, 4)
	b += a
	*pa, *pb, *pc = a, b, c
}

// Final implements the final mixing of three 32-bit values (a, b, c) into c.
func Final(pa, pb, pc *uint32) {
	a, b, c := *pa, *pb, *pc
	c ^= b
	c -= rot(b, 14)
	a ^= c
	a -= rot(c, 11)
	b ^= a
	b -= rot(a, 25)
	c ^= b
	c -= rot(b, 16)
	a ^= c
	a -= rot(c, 4)
	b ^= a
	b -= rot(a, 14)
	c ^= b
	c -= rot(b, 24)
	*pa, *pb, *pc = a, b, c
}

func HashWord(key []uint32, length, initValue uint32) uint32 {
	var a, b, c uint32
	a = 0xdeadbeef + (length << 2) + initValue
	b = a
	c = a

	k := 0

	for length > 3 {
		a += key[k+0]
		b += key[k+1]
		c += key[k+2]
		Mix(&a, &b, &c)
		length -= 3
		k += 3
	}

	switch length {
	case 3:
		c += key[k+2]
		fallthrough
	case 2:
		b += key[k+1]
		fallthrough
	case 1:
		a += key[k+0]
		Final(&a, &b, &c)
	case 0:
	}

	return c
}

func HashWord2(key []uint32, length uint32, pc, pb *uint32) {
	var a, b, c uint32
	a = 0xdeadbeef + (length << 2) + *pc
	b = a
	c = a + *pb

	i := 0

	for length > 3 {
		a += key[i+0]
		b += key[i+1]
		c += key[i+2]
		Mix(&a, &b, &c)
		length -= 3
		i += 3
	}

	switch length {
	case 3:
		c += key[i+2]
		fallthrough
	case 2:
		b += key[i+1]
		fallthrough
	case 1:
		a += key[i+0]
		Final(&a, &b, &c)
	case 0:
	}

	*pc = c
	*pb = b
}

func HashLittle(key []uint8, length, initValue uint32) uint32 {
	var a, b, c uint32
	a = 0xdeadbeef + length + initValue
	b = a
	c = a

	ui := binary.LittleEndian.Uint32(key[0:4])
	i := 0

	k := func(j int) uint32 {
		return binary.LittleEndian.Uint32(key[i+(j*4) : i+(j*4)+4])
	}

	k16 := func(j int) uint32 {
		return uint32(binary.LittleEndian.Uint16(key[i+(j*2) : i+(j*2)+2]))
	}

	k8 := func(j int) uint8 {
		return key[i+j]
	}

	if ui&0x3 == 0 {
		for length > 12 {
			a += k(0)
			b += k(1)
			c += k(2)
			Mix(&a, &b, &c)
			length -= 12
			i += 3 * 4
		}
		switch length {
		case 12:
			c += k(2)
			b += k(1)
			a += k(0)
		case 11:
			c += uint32(k8(10)) << 16
			fallthrough
		case 10:
			c += uint32(k8(9)) << 8
			fallthrough
		case 9:
			c += uint32(k8(8))
			fallthrough
		case 8:
			b += k(1)
			a += k(0)
		case 7:
			b += uint32(k8(6)) << 16
			fallthrough
		case 6:
			b += uint32(k8(5)) << 8
			fallthrough
		case 5:
			b += uint32(k8(4))
			fallthrough
		case 4:
			a += k(0)
		case 3:
			a += uint32(k8(2)) << 16
			fallthrough
		case 2:
			a += uint32(k8(1)) << 8
			fallthrough
		case 1:
			a += uint32(k8(0))
		case 0:
			return c
		}
	} else if ui%0x1 == 0 {
		for length > 12 {
			a += k16(0) + (k16(1) << 16)
			b += k16(2) + (k16(3) << 16)
			c += k16(4) + (k16(5) << 16)
			Mix(&a, &b, &c)
			length -= 12
			i += 6 * 2
		}
		switch length {
		case 12:
			c += k16(4) + (k16(5) << 16)
			b += k16(2) + (k16(3) << 16)
			a += k16(0) + (k16(1) << 16)
		case 11:
			c += uint32(k8(10)) << 16
			fallthrough
		case 10:
			c += k16(4)
			b += k16(2) + (k16(3) << 16)
			a += k16(0) + (k16(1) << 16)
		case 9:
			c += uint32(k8(8))
			fallthrough
		case 8:
			b += k16(2) + (k16(3) << 16)
			a += k16(0) + (k16(1) << 16)
		case 7:
			b += uint32(k8(6)) << 16
			fallthrough
		case 6:
			b += k16(2)
			a += k16(0) + (k16(1) << 16)
		case 5:
			b += uint32(k8(4))
			fallthrough
		case 4:
			a += k16(0) + (k16(1) << 16)
		case 3:
			a += uint32(k8(2)) << 16
			fallthrough
		case 2:
			a += k16(0)
		case 1:
			a += uint32(k8(0))
		case 0:
			return c
		}
	} else {
		for length > 12 {
			a += uint32(k8(0))
			a += uint32(k8(1)) << 8
			a += uint32(k8(2)) << 16
			a += uint32(k8(3)) << 24
			b += uint32(k8(4))
			b += uint32(k8(5)) << 8
			b += uint32(k8(6)) << 16
			b += uint32(k8(7)) << 24
			c += uint32(k8(8))
			c += uint32(k8(9)) << 8
			c += uint32(k8(10)) << 16
			c += uint32(k8(11)) << 24
			Mix(&a, &b, &c)
			length -= 12
			i += 12
		}
		switch length {
		case 12:
			c += uint32(k(11)) << 24
			fallthrough
		case 11:
			c += uint32(k8(10)) << 16
			fallthrough
		case 10:
			c += uint32(k8(9)) << 8
			fallthrough
		case 9:
			c += uint32(k8(8))
			fallthrough
		case 8:
			b += uint32(k(7)) << 24
			fallthrough
		case 7:
			b += uint32(k8(6)) << 16
			fallthrough
		case 6:
			b += uint32(k8(5)) << 8
			fallthrough
		case 5:
			b += uint32(k8(4))
			fallthrough
		case 4:
			a += uint32(k(3)) << 24
			fallthrough
		case 3:
			a += uint32(k8(2)) << 16
			fallthrough
		case 2:
			a += uint32(k8(1)) << 8
			fallthrough
		case 1:
			a += uint32(k8(0))
		case 0:
			return c
		}
	}

	Final(&a, &b, &c)
	return c
}

func HashLittle2(key []uint8, length uint32, pc, pb *uint32) {
	var a, b, c uint32
	a = 0xdeadbeef + length + *pc
	b = a
	c = a + *pb

	ui := binary.LittleEndian.Uint32(key[0:4])
	i := 0

	k := func(j int) uint32 {
		return binary.LittleEndian.Uint32(key[i+(j*4) : i+(j*4)+4])
	}

	k16 := func(j int) uint32 {
		return uint32(binary.LittleEndian.Uint16(key[i+(j*2) : i+(j*2)+2]))
	}

	k8 := func(j int) uint8 {
		return key[i+j]
	}

	if ui&0x3 == 0 {
		for length > 12 {
			a += k(0)
			b += k(1)
			c += k(2)
			Mix(&a, &b, &c)
			length -= 12
			i += 3 * 4
		}
		switch length {
		case 12:
			c += k(2)
			b += k(1)
			a += k(0)
		case 11:
			c += uint32(k8(10)) << 16
			fallthrough
		case 10:
			c += uint32(k8(9)) << 8
			fallthrough
		case 9:
			c += uint32(k8(8))
			fallthrough
		case 8:
			b += k(1)
			a += k(0)
		case 7:
			b += uint32(k8(6)) << 16
			fallthrough
		case 6:
			b += uint32(k8(5)) << 8
			fallthrough
		case 5:
			b += uint32(k8(4))
			fallthrough
		case 4:
			a += k(0)
		case 3:
			a += uint32(k8(2)) << 16
			fallthrough
		case 2:
			a += uint32(k8(1)) << 8
			fallthrough
		case 1:
			a += uint32(k8(0))
		case 0:
			*pc = c
			*pb = b
			return
		}
	} else if ui%0x1 == 0 {
		for length > 12 {
			a += k16(0) + (k16(1) << 16)
			b += k16(2) + (k16(3) << 16)
			c += k16(4) + (k16(5) << 16)
			Mix(&a, &b, &c)
			length -= 12
			i += 6 * 2
		}
		switch length {
		case 12:
			c += k16(4) + (k16(5) << 16)
			b += k16(2) + (k16(3) << 16)
			a += k16(0) + (k16(1) << 16)
		case 11:
			c += uint32(k8(10)) << 16
			fallthrough
		case 10:
			c += k16(4)
			b += k16(2) + (k16(3) << 16)
			a += k16(0) + (k16(1) << 16)
		case 9:
			c += uint32(k8(8))
			fallthrough
		case 8:
			b += k16(2) + (k16(3) << 16)
			a += k16(0) + (k16(1) << 16)
		case 7:
			b += uint32(k8(6)) << 16
			fallthrough
		case 6:
			b += k16(2)
			a += k16(0) + (k16(1) << 16)
		case 5:
			b += uint32(k8(4))
			fallthrough
		case 4:
			a += k16(0) + (k16(1) << 16)
		case 3:
			a += uint32(k8(2)) << 16
			fallthrough
		case 2:
			a += k16(0)
		case 1:
			a += uint32(k8(0))
		case 0:
			*pc = c
			*pb = b
			return
		}
	} else {
		fmt.Println("else")
		for length > 12 {
			a += uint32(k8(0))
			a += uint32(k8(1)) << 8
			a += uint32(k8(2)) << 16
			a += uint32(k8(3)) << 24
			b += uint32(k8(4))
			b += uint32(k8(5)) << 8
			b += uint32(k8(6)) << 16
			b += uint32(k8(7)) << 24
			c += uint32(k8(8))
			c += uint32(k8(9)) << 8
			c += uint32(k8(10)) << 16
			c += uint32(k8(11)) << 24
			Mix(&a, &b, &c)
			length -= 12
			i += 12
		}
		switch length {
		case 12:
			c += uint32(k(11)) << 24
			fallthrough
		case 11:
			c += uint32(k8(10)) << 16
			fallthrough
		case 10:
			c += uint32(k8(9)) << 8
			fallthrough
		case 9:
			c += uint32(k8(8))
			fallthrough
		case 8:
			b += uint32(k(7)) << 24
			fallthrough
		case 7:
			b += uint32(k8(6)) << 16
			fallthrough
		case 6:
			b += uint32(k8(5)) << 8
			fallthrough
		case 5:
			b += uint32(k8(4))
			fallthrough
		case 4:
			a += uint32(k(3)) << 24
			fallthrough
		case 3:
			a += uint32(k8(2)) << 16
			fallthrough
		case 2:
			a += uint32(k8(1)) << 8
			fallthrough
		case 1:
			a += uint32(k8(0))
		case 0:
			*pc = c
			*pb = b
			return
		}
	}

	Final(&a, &b, &c)
	*pc = c
	*pb = b
}

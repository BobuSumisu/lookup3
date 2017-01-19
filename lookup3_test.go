package lookup3

import "testing"

// Test values generated (and taken) from lookup3.c.

var rotTests = []struct {
	n uint32
	k uint8
	r uint32
}{
	{1, 1, 2},
	{2, 2, 8},
	{0xdeadbeef, 1, 0xbd5b7ddf},
	{0xdeadbeef, 16, 0xbeefdead},
}

var mixTests = []struct {
	aIn, bIn, cIn    uint32
	aOut, bOut, cOut uint32
}{
	{1, 0, 0, 0xbfbe7f3a, 0x7d633d2b, 0x9810de96},
	{0, 1, 0, 0xff000003, 0xf80ff701, 0x960f65eb},
	{0, 0, 1, 0xeae4d1c9, 0x2dcd1961, 0x9397e6b6},
	{1, 1, 1, 0xf5e6e7d0, 0x9ad19143, 0x15bbdf44},
	{0xdeadbeef, 0xdeadbeef, 0xdeadbeef, 0x108a17ff, 0xc0bfb8ee, 0x8c37cf7c},
}

var finalTests = []struct {
	aIn, bIn, cIn    uint32
	aOut, bOut, cOut uint32
}{
	{1, 0, 0, 0x1e1de1e1, 0x67a59a59, 0x40f3f2be},
	{0, 1, 0, 0x3a989830, 0xa15719ed, 0x02961431},
	{0, 0, 1, 0xbd324300, 0xb60d8bc5, 0xf02a6a49},
	{1, 1, 1, 0xb1881916, 0x462df523, 0xec1761c6},
	{0xdeadbeef, 0xdeadbeef, 0xdeadbeef, 0x1523f639, 0x6d004bb2, 0x31b8a510},
}

var hashWordTests = []struct {
	key       []uint32
	keyLength uint32
	initValue uint32
	hash      uint32
}{
	{[]uint32{1}, 1, 0, 0x72a82a9b},
	{[]uint32{1, 2}, 2, 1, 0x73989811},
	{[]uint32{1, 2, 3}, 3, 0, 0xa46158f5},
	{[]uint32{1, 2, 3, 4}, 4, 1, 0x044ec9ea},
	{[]uint32{1, 2, 3, 4, 5}, 5, 2, 0x39a100d5},
}

var hashWord2Tests = []struct {
	key        []uint32
	keyLength  uint32
	cIn, bIn   uint32
	cOut, bOut uint32
}{
	{[]uint32{1, 2}, 2, 1, 1, 0x301b0127, 0x3ce9fe7e},
	{[]uint32{1, 2, 3}, 3, 0, 0, 0xa46158f5, 0x45915a7e},
	{[]uint32{1, 2, 3, 4}, 4, 1, 0, 0x044ec9ea, 0x729b6663},
	{[]uint32{1, 2, 3, 4, 5}, 5, 0, 1, 0x7489c25b, 0x898e47dd},
}

var hashLittleTests = []struct {
	key       string
	keyLength uint32
	initValue uint32
	hash      uint32
}{
	{"Four score and seven years ago", 30, 0, 0x17770551},
	{"Four score and seven years ago", 30, 1, 0xcd628161},
	{"hello world", 12, 0, 0x14973b58},
	{"hello world", 12, 1, 0x4746e99f},
}

var hashLittle2Tests = []struct {
	key        string
	keyLength  uint32
	initValue1 uint32
	initValue2 uint32
	hash1      uint32
	hash2      uint32
}{
	{"", 0, 0, 0, 0xdeadbeef, 0xdeadbeef},
	{"", 0, 0, 0xdeadbeef, 0xbd5b7dde, 0xdeadbeef},
	{"", 0, 0xdeadbeef, 0xdeadbeef, 0x9c093ccd, 0xbd5b7dde},
	{"Four score and seven years ago", 30, 0, 0, 0x17770551, 0xce7226e6},
	{"Four score and seven years ago", 30, 1, 0, 0xcd628161, 0x6cbea4b3},
	{"Four score and seven years ago", 30, 0, 1, 0xe3607cae, 0xbd371de4},
	{"hello world", 12, 0, 0, 0x14973b58, 0xd2d6bcf5},
	{"hello world", 12, 1, 0, 0x4746e99f, 0xcf623d47},
	{"hello world", 12, 0, 1, 0xe1c4d473, 0x884cad01},
}

func TestRot(t *testing.T) {
	for _, tt := range rotTests {
		r := rot(tt.n, tt.k)
		if r != tt.r {
			t.Errorf("rot(0x%08x, %d) => 0x%08x, want 0x%08x", tt.n, tt.k, r, tt.r)
		}
	}
}

func TestMix(t *testing.T) {
	for _, tt := range mixTests {
		a, b, c := tt.aIn, tt.bIn, tt.cIn
		Mix(&a, &b, &c)
		if a != tt.aOut || b != tt.bOut || c != tt.cOut {
			t.Errorf("Mix(%d, %d, %d) => (0x%08x, 0x%08x, 0x%08x), want (0x%08x, 0x%08x, 0x%08x)\n",
				tt.aIn, tt.bIn, tt.cIn, a, b, c, tt.aOut, tt.bOut, tt.cOut)
		}
	}
}

func TestFinal(t *testing.T) {
	for _, tt := range finalTests {
		a, b, c := tt.aIn, tt.bIn, tt.cIn
		Final(&a, &b, &c)
		if a != tt.aOut || b != tt.bOut || c != tt.cOut {
			t.Errorf("Final(%d, %d, %d) => (0x%08x, 0x%08x, 0x%08x), want (0x%08x, 0x%08x, 0x%08x)\n",
				tt.aIn, tt.bIn, tt.cIn, a, b, c, tt.aOut, tt.bOut, tt.cOut)
		}
	}
}

func TestHashWord(t *testing.T) {
	for _, tt := range hashWordTests {
		h := HashWord(tt.key, tt.keyLength, tt.initValue)
		if h != tt.hash {
			t.Errorf("HashWord(%q, %d, %d) => 0x%08x, want 0x%08x\n", tt.key, tt.keyLength,
				tt.initValue, h, tt.hash)
		}
	}
}

func TestHashWord2(t *testing.T) {
	for _, tt := range hashWord2Tests {
		c, b := tt.cIn, tt.bIn
		HashWord2(tt.key, tt.keyLength, &c, &b)
		if c != tt.cOut || b != tt.bOut {
			t.Errorf("HashWord2(%q, %d, %d, %d) => (0x%08x, 0x%08x), want (0x%08x, 0x%08x)\n",
				tt.key, tt.keyLength, tt.cIn, tt.bIn, c, b, tt.cOut, tt.bOut)
		}
	}
}

func TestHashLittle(t *testing.T) {
	for _, tt := range hashLittleTests {
		h := HashLittle([]uint8(tt.key), tt.keyLength, tt.initValue)
		if h != tt.hash {
			t.Errorf("HashLittle(%q, %d, %d) => 0x%08x, want 0x%08x", tt.key, tt.keyLength,
				tt.initValue, h, tt.hash)
		}
	}
}

func TestHashLittle2(t *testing.T) {
	for _, tt := range hashLittle2Tests {
		h1 := tt.initValue1
		h2 := tt.initValue2
		HashLittle2([]uint8(tt.key), tt.keyLength, &h1, &h2)
		if h1 != tt.hash1 || h2 != tt.hash2 {
			t.Errorf("HashLittle(%q, %d, %d, %d) => (0x%08x, 0x%08x), want (0x%08x, 0x%08x)",
				tt.key, tt.keyLength, tt.initValue1, tt.initValue2, h1, h2, tt.hash1, tt.hash2)
		}
	}
}

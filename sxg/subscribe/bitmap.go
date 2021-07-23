package subscribe

import "encoding/binary"

var LimitMaxValue uint64 = 0x100000000 //2^32

type BitMap interface {
	Existed(value uint64) bool
	Put(value uint64) error
	Remove(value uint64)
	Resize(value uint64) error
	Range() (min, max uint64)
	Size() int
	Count() int
}

func indexAndMask(value uint64) (index uint64, mask byte) {
	index = value >> 3
	mod := (byte)(value & 0x07)
	if mod == 0 {
		index -= 1
		mask = 0x01 << 7
	} else {
		mask = 0x01 << (mod - 1)
	}
	return
}

func normalizedWithLimit(value, offset, limit uint64) (uint64, bool) {
	value = value - offset
	return value, (value > 0 && value <= limit)
}

//value - offset must range from 1 to LimitMaxValue
func normalized(value, offset uint64) (uint64, bool) {
	return normalizedWithLimit(value, offset, LimitMaxValue)
}

//caculate the sum of 1(from redis)
func caculateBitCount(p []byte) int {
	count := 0
	size := len(p)
	left := size
	for left >= 28 {
		var v1, v2, v3, v4, v5, v6, v7 uint32

		v1 = binary.BigEndian.Uint32(p[size-left : size-left+4])
		v2 = binary.BigEndian.Uint32(p[size-left+4 : size-left+8])
		v3 = binary.BigEndian.Uint32(p[size-left+8 : size-left+12])
		v4 = binary.BigEndian.Uint32(p[size-left+12 : size-left+16])
		v5 = binary.BigEndian.Uint32(p[size-left+16 : size-left+20])
		v6 = binary.BigEndian.Uint32(p[size-left+20 : size-left+24])
		v7 = binary.BigEndian.Uint32(p[size-left+24 : size-left+28])

		v1 = v1 - ((v1 >> 1) & 0x55555555)
		v1 = (v1 & 0x33333333) + ((v1 >> 2) & 0x33333333)
		v2 = v2 - ((v2 >> 1) & 0x55555555)
		v2 = (v2 & 0x33333333) + ((v2 >> 2) & 0x33333333)
		v3 = v3 - ((v3 >> 1) & 0x55555555)
		v3 = (v3 & 0x33333333) + ((v3 >> 2) & 0x33333333)
		v4 = v4 - ((v4 >> 1) & 0x55555555)
		v4 = (v4 & 0x33333333) + ((v4 >> 2) & 0x33333333)
		v5 = v5 - ((v5 >> 1) & 0x55555555)
		v5 = (v5 & 0x33333333) + ((v5 >> 2) & 0x33333333)
		v6 = v6 - ((v6 >> 1) & 0x55555555)
		v6 = (v6 & 0x33333333) + ((v6 >> 2) & 0x33333333)
		v7 = v7 - ((v7 >> 1) & 0x55555555)
		v7 = (v7 & 0x33333333) + ((v7 >> 2) & 0x33333333)

		bitCount := ((((v1 + (v1 >> 4)) & 0x0F0F0F0F) +
			((v2 + (v2 >> 4)) & 0x0F0F0F0F) +
			((v3 + (v3 >> 4)) & 0x0F0F0F0F) +
			((v4 + (v4 >> 4)) & 0x0F0F0F0F) +
			((v5 + (v5 >> 4)) & 0x0F0F0F0F) +
			((v6 + (v6 >> 4)) & 0x0F0F0F0F) +
			((v7 + (v7 >> 4)) & 0x0F0F0F0F)) * 0x01010101) >> 24
		count += int(bitCount)

		left -= 28
	}
	if left > 0 {
		for _, b := range p[size-left : size] {
			num := byte(0)
			mask := byte(0x80)
			for k := 7; k >= 0; k-- {
				num += (b & mask) >> k
				mask = mask >> 1
			}
			count += int(num)
		}
	}
	return count
}

package oracleparser

import (
	"encoding/binary"
	"bytes"
)

const UseBigClrChunks = true

func GetInt64(size int, compress bool, bigEndian bool, buffer []byte, index int) (int64, int) {
	var ret int64
	negFlag := false
	if compress {
		rb := buffer[index]
		index += 1
		size = int(rb)
		if size&0x80 > 0 {
			negFlag = true
			size = size & 0x7F
		}
		bigEndian = true
	}
	rb := buffer[index:index+size]
	index += size
	temp := make([]byte, 8)
	if bigEndian {
		copy(temp[8-size:], rb)
		ret = int64(binary.BigEndian.Uint64(temp))
	} else {
		copy(temp[:size], rb)
		//temp = append(pck.buffer[pck.index: pck.index + size], temp...)
		ret = int64(binary.LittleEndian.Uint64(temp))
	}
	if negFlag {
		ret = ret * -1
	}
	return ret, index
}

func GetInt(size int, compress bool, bigEndian bool, buffer []byte, index int) (int, int)  {
	temp, returnindex := GetInt64(size, compress, bigEndian, buffer, index)
	return int(temp), returnindex
}

func GetKeyVal(buffer []byte, index int) (key []byte, val []byte, num int, returnIndex int) {
	key, index = GetDlc(buffer, index)
	val,index = GetDlc(buffer, index)
	num, returnIndex = GetInt(4, true, true, buffer, index)
	return
}

func GetDlc(buffer []byte, index int) (output []byte, returnIndex int) {
	var length int
	length, returnIndex = GetInt(4, true, true, buffer, index)
	if length > 0 {
		output, returnIndex = GetClr(buffer, index)
		if len(output) > length {
			output = output[:length]
		}
	}
	return
}

func GetClr(buffer []byte, index int) (output []byte, returnIndex int) {
	var nb byte
	nb, index = GetByte(buffer, index)
	if nb == 0 || nb == 0xFF {
		output = nil
		return
	}
	chunkSize := int(nb)
	var chunk []byte
	var tempBuffer bytes.Buffer
	if chunkSize == 0xFE {
		for chunkSize > 0 {
			if UseBigClrChunks {
				chunkSize, index = GetInt(4, true, true, buffer, index)
			} else {
				nb, index = GetByte(buffer, index)
				chunkSize = int(nb)
			}
			chunk, index = GetBytes(chunkSize, buffer, index)
			tempBuffer.Write(chunk)
		}
	} else {
		chunk, index = GetBytes(chunkSize, buffer, index)
		tempBuffer.Write(chunk)
	}
	output = tempBuffer.Bytes()
	returnIndex = index
	return
}

func GetByte(buffer []byte, index int) (uint8, int) {
	rb := buffer[index]
	return rb, index+1
}

func GetBytes(length int, buffer []byte, index int) ([]byte, int) {
	rb := buffer[index: index+length]
	index += length
	return rb , index
}





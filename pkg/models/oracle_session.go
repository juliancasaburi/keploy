package models

import (
	"bytes"
	"encoding/binary"
)

func (session *PacketSession) GetInt64(size int, compress bool, bigEndian bool, buffer []byte, index int) (int64, int) {
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
	rb := buffer[index : index+size]
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

func (session *PacketSession) GetInt(size int, compress bool, bigEndian bool, buffer []byte, index int) (int, int) {
	temp, returnindex := session.GetInt64(size, compress, bigEndian, buffer, index)
	return int(temp), returnindex
}

func (session *PacketSession) GetKeyVal(buffer []byte, index int) (key []byte, val []byte, num int, returnIndex int) {
	key, index = session.GetDlc(buffer, index)
	val, index = session.GetDlc(buffer, index)
	num, returnIndex = session.GetInt(4, true, true, buffer, index)
	return
}

func (session *PacketSession) GetDlc(buffer []byte, index int) (output []byte, returnIndex int) {
	var length int
	length, returnIndex = session.GetInt(4, true, true, buffer, index)
	if length > 0 {
		output, returnIndex = session.GetClr(buffer, index)
		if len(output) > length {
			output = output[:length]
		}
	}
	return
}

func (session *PacketSession) GetClr(buffer []byte, index int) (output []byte, returnIndex int) {
	var nb byte
	nb, index = session.GetByte(buffer, index)
	if nb == 0 || nb == 0xFF {
		output = nil
		return
	}
	chunkSize := int(nb)
	var chunk []byte
	var tempBuffer bytes.Buffer
	if chunkSize == 0xFE {
		for chunkSize > 0 {
			if session.UseBigClrChunks {
				chunkSize, index = session.GetInt(4, true, true, buffer, index)
			} else {
				nb, index = session.GetByte(buffer, index)
				chunkSize = int(nb)
			}
			chunk, index = session.GetBytes(chunkSize, buffer, index)
			tempBuffer.Write(chunk)
		}
	} else {
		chunk, index = session.GetBytes(chunkSize, buffer, index)
		tempBuffer.Write(chunk)
	}
	returnIndex = index
	output = tempBuffer.Bytes()
	return
}

func (session *PacketSession) GetByte(buffer []byte, index int) (uint8, int) {
	rb := buffer[index]
	return rb, index + 1
}

func (session *PacketSession) GetBytes(length int, buffer []byte, index int) ([]byte, int) {
	rb := buffer[index : index+length]
	index += length
	return rb, index
}

func (session *PacketSession) GetNullTerminatedString(buffer []byte, index int) (string, int) {
	startingIndex := index
	for index < len(buffer) && buffer[index] != 0 {
		index++
	}
	driverName := string(buffer[startingIndex:index])
	index += 1
	return driverName, index
}

func (session *PacketSession) GetNullTerminatedArray(buffer []byte, index int) ([]byte, int) {
	var arrayTerminator []byte
	for index < len(buffer) && buffer[index] != 0 {
		arrayTerminator = append(arrayTerminator, buffer[index])
		index += 1
	}
	index += 1
	return arrayTerminator, index
}
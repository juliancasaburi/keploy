package oracleparser

import (
	"encoding/binary"
	"strings"

	"go.keploy.io/server/pkg/models"
)

func DecodeRedirectPacket(Packets [][]byte, dataPacketType models.DataPacketType) (models.OracleHeader, interface{}, bool, models.DataPacketType, interface{}, error) {
	var packetData []byte
	var packetLength int
	var nextDataPacketType models.DataPacketType
	var redirectAddress string
	var redirectData string
	for _, slice := range Packets {
		packetData = append(packetData, slice...)
	}
	if session.Context.Version >= 315 {
		packetLength = int(binary.BigEndian.Uint32(Packets[0][0:]))
	} else {
		packetLength = int(binary.BigEndian.Uint16(Packets[0][0:]))
	}
	dataOffset := uint16(10)
	dataLen := binary.BigEndian.Uint16(packetData[8:])
	requestHeader := models.OracleHeader{
		PacketLength: packetLength,
		PacketType:   models.PacketTypeFromUint8(packetData[4]),
		PacketFlag:   Packets[0][5],
		Session:      session,
	}
	if requestHeader.PacketLength <= int(dataOffset) {
		nextDataPacketType = models.OracleRedirectDataMessageType
	} else {
		data := string(packetData[10 : 10+dataLen])
		length := strings.Index(data, "\x00")
		if requestHeader.PacketFlag&2 != 0 && length > 0 {
			redirectAddress = data[:length]
			redirectData = data[length:]
		} else {
			redirectAddress = data
		}
	}
	requestMessage := models.OracleRedirectMessage{
		DataOffset:      10,
		DataLength:      dataLen,
		RedirectAddress: redirectAddress,
		RedirectData:    redirectData,
	}
	return requestHeader, requestMessage, false, nextDataPacketType, nil, nil
}

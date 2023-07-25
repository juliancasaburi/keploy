package oracleparser

import (
	"encoding/binary"
	"strings"

	"github.com/sijms/go-ora/v2/network"
	"go.keploy.io/server/pkg/models"
)

func DecodeRedirectPacket(Packets [][]byte, dataPacketType models.DataPacketType) (models.OracleHeader, interface{}, bool, models.DataPacketType, error) {
	var packetData []byte
	var packetLength interface{}
	var nextDataPacketType models.DataPacketType
	var redirectAddress string
	var redirectData string
	for _, slice := range Packets {
		packetData = append(packetData, slice...)
	}
	if session.Context.Version >= 315 {
		packetLength = binary.BigEndian.Uint32(Packets[0][0:])
	} else {
		packetLength = binary.BigEndian.Uint16(Packets[0][0:])
	}
	dataOffset := uint16(10)
	dataLen := binary.BigEndian.Uint16(packetData[8:])
	requestHeader := models.OracleHeader{
		PacketLength: packetLength,
		PacketType:   network.PacketType(Packets[0][4]),
		PacketFlag:   Packets[0][5],
		Session:      session,
	}
	if requestHeader.PacketLength.(uint16) <= dataOffset {
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
		DATA_OFFSET:      10,
		DATA_LENGTH:      dataLen,
		REDIRECT_ADDRESS: redirectAddress,
		REDIRECT_DATA:    redirectData,
	}
	return requestHeader, requestMessage, false, nextDataPacketType, nil
}

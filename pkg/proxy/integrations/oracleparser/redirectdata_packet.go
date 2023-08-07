package oracleparser

import (
	"encoding/binary"
	"strings"

	"go.keploy.io/server/pkg/models"
)

func DecodeRedirectDataMessage(Packets [][]byte) (models.OracleHeader, interface{}, bool, models.DataPacketType, interface{}, error) {
	var packetData []byte
	var packetLength int
	var redirectAddress string
	var redirectData string
	for _, slice := range Packets {
		packetData = append(packetData, slice...)
	}
	if session.Context.Version >= 315 {
		packetLength = int(binary.BigEndian.Uint32(packetData[0:]))
	} else {
		packetLength = int(binary.BigEndian.Uint16(packetData[0:]))
	}
	requestHeader := models.OracleHeader{
		PacketLength: packetLength,
		PacketType:   models.PacketTypeFromUint8(packetData[4]),
		PacketFlag:   packetData[5],
		Session:      session,
	}
	data := string(packetData[10:])
	length := strings.Index(data, "\x00")
	if requestHeader.PacketFlag&2 != 0 && length > 0 {
		redirectAddress = data[:length]
		redirectData = data[length:]
	} else {
		redirectAddress = data
	}
	oracleRedirectDataMessage := models.OracleRedirectDataMessage{
		RedirectAddress: redirectAddress,
		RedirectData:    redirectData,
	}
	requestMessage := models.OracleDataMessage{
		DataOffset:      10,
		DataMessageType: models.OracleRedirectDataMessageType,
		DataMessage:     oracleRedirectDataMessage,
	}
	return requestHeader, requestMessage, false, models.DefaultDataPacket, nil, nil
}

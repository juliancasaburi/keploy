package oracleparser

import (
	"encoding/binary"

	"go.keploy.io/server/pkg/models"
)

func DecodeConnectionDataMessage(Packets [][]byte) (models.OracleHeader, interface{}, bool, models.DataPacketType, interface{}, error) {
	var packetData []byte
	var packetLength int
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
	connectionString := string(packetData[10:])
	oracleConnectionDataMessage := models.OracleConnectionDataMessage{
		ConnectString: connectionString,
	}
	requestMessage := models.OracleDataMessage{
		DataOffset:      10,
		DataMessageType: models.OracleConnectionDataMessageType,
		DataMessage:     oracleConnectionDataMessage,
	}
	return requestHeader, requestMessage, false, models.DefaultDataPacket, nil, nil
}

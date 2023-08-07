package oracleparser

import (
	"encoding/binary"

	"go.keploy.io/server/pkg/models"
)

func DecodeMarkerPacket(Packets [][]byte, dataPacketType models.DataPacketType) (models.OracleHeader, interface{}, bool, models.DataPacketType, interface{}, error) {
	var packetData []byte
	var packetLength int
	for _, slice := range Packets {
		packetData = append(packetData, slice...)
	}
	if session.HandShakeComplete && session.Context.Version >= 315 {
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
	requestMessage := models.OracleMarkerMessage{
		MarkerType: models.MarkerTypeFromInt(int(packetData[8])),
		MarkerData: packetData[10],
	}
	if (requestMessage.MarkerType == models.TNS_MARKER_TYPE_BREAK && requestMessage.MarkerData != 2) || requestMessage.MarkerType == models.TNS_MARKER_TYPE_RESET {
		return requestHeader, requestMessage, true, dataPacketType, nil, nil
	}
	return requestHeader, requestMessage, false, dataPacketType, nil, nil
}

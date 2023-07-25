package oracleparser

import (
	"encoding/binary"

	"github.com/sijms/go-ora/v2/network"
	"go.keploy.io/server/pkg/models"
)

func DecodeMarkerPacket(Packets [][]byte, dataPacketType models.DataPacketType) (models.OracleHeader, interface{}, bool, models.DataPacketType, error) {
	var packetData []byte
	var packetLength interface{}
	for _, slice := range Packets {
		packetData = append(packetData, slice...)
	}
	if session.HandShakeComplete && session.Context.Version >= 315 {
		packetLength = binary.BigEndian.Uint32(packetData[0:])
	} else {
		packetLength = binary.BigEndian.Uint16(packetData[0:])
	}
	requestHeader := models.OracleHeader{
		PacketLength: packetLength,
		PacketType:   network.PacketType(packetData[4]),
		PacketFlag:   packetData[5],
		Session:      session,
	}
	requestMessage := models.OracleMarkerMessage{
		MARKER_TYPE: packetData[8],
		MARKER_DATA: packetData[10],
	}
	if (requestMessage.MARKER_TYPE == 1 && requestMessage.MARKER_DATA != 2) || requestMessage.MARKER_TYPE == 0 {
		return requestHeader, requestMessage, true, dataPacketType, nil
	}
	return requestHeader, requestMessage, false, dataPacketType, nil
}

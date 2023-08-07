package oracleparser

import (
	"encoding/binary"

	"go.keploy.io/server/pkg/models"
)

func DecodeControlPacket(Packets [][]byte, dataPacketType models.DataPacketType) (models.OracleHeader, interface{}, bool, models.DataPacketType, interface{}, error) {
	var packetData []byte
	var packetLength int
	var requestMessage models.OracleControlMessage
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
	controlType := models.ControlTypeFromInt(int(binary.BigEndian.Uint16(packetData[8:])))
	if controlType == "TNS_CONTROL_TYPE_INBAND_NOTIFICATION" {
		requestMessage = models.OracleControlMessage{
			ControlType:  controlType,
			ControlError: models.ControlErrorFromInt(int(binary.BigEndian.Uint32(packetData[14:]))),
		}
	} else {
		requestHeader.Session.SupportOOB = false
		session.SupportOOB = false
		requestMessage = models.OracleControlMessage{
			ControlType: controlType,
		}
	}
	return requestHeader, requestMessage, false, dataPacketType, nil, nil
}

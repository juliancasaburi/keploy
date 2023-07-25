package oracleparser

import (
	"encoding/binary"

	"github.com/sijms/go-ora/v2/network"
	"go.keploy.io/server/pkg/models"
)

func DecodeControlPacket(Packets [][]byte, dataPacketType models.DataPacketType) (models.OracleHeader, interface{}, bool, models.DataPacketType, error) {
	var packetData []byte
	var packetLength interface{}
	var requestMessage models.OracleControlMessage
	for _, slice := range Packets {
		packetData = append(packetData, slice...)
	}
	if session.Context.Version >= 315 {
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
	controlType := models.ControlType(binary.BigEndian.Uint16(packetData[8:]))
	if controlType == 8 {
		requestMessage = models.OracleControlMessage{
			CONTROL_TYPE:  controlType,
			CONTROL_ERROR: models.ControlError(binary.BigEndian.Uint32(packetData[14:])),
		}
	} else {
		requestHeader.Session.SupportOOB = false
		session.SupportOOB = false
		requestMessage = models.OracleControlMessage{
			CONTROL_TYPE: controlType,
		}
	}
	return requestHeader, requestMessage, false, dataPacketType, nil
}

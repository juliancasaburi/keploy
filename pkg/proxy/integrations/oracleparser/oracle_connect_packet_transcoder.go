package oracleparser

import (
	"encoding/binary"

	"github.com/sijms/go-ora/v2/network"
	"go.keploy.io/server/pkg/models"
)

func DecodeConnectPacket(Packets [][]byte, dataPacketType models.DataPacketType) (models.OracleHeader, interface{}, bool, models.DataPacketType, error) {
	var connectionString string
	var nextDataPacketType models.DataPacketType
	PacketData := Packets[0]
	connectionStringLength := binary.BigEndian.Uint16(PacketData[24:])
	if connectionStringLength > TNS_MAX_CONNECT_DATA {
		nextDataPacketType = models.OracleConnectionDataMessageType
	}
	dataOffSet := binary.BigEndian.Uint16(PacketData[26:])
	PacketLength := binary.BigEndian.Uint16(PacketData[0:])
	if nextDataPacketType == models.Default {
		connectionString = string(PacketData[dataOffSet:])
	}
	requestHeader := models.OracleHeader{
		PacketLength: PacketLength,
		PacketType:   network.PacketType(PacketData[4]),
		PacketFlag:   PacketData[5],
	}
	requestMessage := models.OracleConnectMessage{
		TNS_VERSION_DESIRED:          binary.BigEndian.Uint16(PacketData[8:]),
		TNS_VERSION_MINIMUM:          binary.BigEndian.Uint16(PacketData[10:]),
		SERVICE_OPTIONS:              binary.BigEndian.Uint16(PacketData[12:]),
		TNS_SDU_16:                   binary.BigEndian.Uint16(PacketData[14:]),
		TNS_TDU_16:                   binary.BigEndian.Uint16(PacketData[16:]),
		TNS_PROTOCOL_CHARACTERISTICS: binary.BigEndian.Uint16(PacketData[18:]),
		LINE_TURNAROUND:              binary.BigEndian.Uint16(PacketData[20:]),
		OURONE:                       binary.BigEndian.Uint16(PacketData[22:]),
		CONNECTION_STRING_LENGTH:     connectionStringLength,
		OFFSET_OF_CONNECTION_DATA:    dataOffSet,
		ACFL0:                        PacketData[32],
		ACFL1:                        PacketData[33],
		TNS_SDU_32:                   binary.BigEndian.Uint32(PacketData[58:]),
		TNS_TDU_32:                   binary.BigEndian.Uint32(PacketData[62:]),
		CONNECT_STRING:               connectionString,
	}
	if requestMessage.OFFSET_OF_CONNECTION_DATA > 70 {
		requestMessage.CONNECT_FLAG_1 = binary.BigEndian.Uint32(PacketData[66:])
		requestMessage.CONNECT_FLAG_2 = binary.BigEndian.Uint32(PacketData[70:])
	}
	session.Context = &network.SessionContext{
		LoVersion:         requestMessage.TNS_VERSION_MINIMUM,
		Options:           requestMessage.SERVICE_OPTIONS,
		SessionDataUnit:   requestMessage.TNS_SDU_32,
		TransportDataUnit: requestMessage.TNS_TDU_32,
		ACFL0:             requestMessage.ACFL0,
		ACFL1:             requestMessage.ACFL1,
		OurOne:            requestMessage.OURONE,
	}
	if requestMessage.CONNECT_FLAG_2 == 1 {
		session.SupportOOB = true
	}
	requestHeader.Session = session
	if nextDataPacketType == 0 {
		return requestHeader, requestMessage, false, nextDataPacketType, nil
	} else {
		return requestHeader, requestMessage, true, nextDataPacketType, nil

	}
}

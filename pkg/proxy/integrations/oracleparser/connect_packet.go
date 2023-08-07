package oracleparser

import (
	"encoding/binary"

	"github.com/sijms/go-ora/v2/network"
	"go.keploy.io/server/pkg/models"
)

func DecodeConnectPacket(Packets [][]byte, dataPacketType models.DataPacketType) (models.OracleHeader, interface{}, bool, models.DataPacketType, interface{}, error) {
	var connectionString string
	nextDataPacketType := models.DefaultDataPacket
	packetData := Packets[0]
	connectionStringLength := binary.BigEndian.Uint16(packetData[24:])
	if connectionStringLength > TNS_MAX_CONNECT_DATA {
		nextDataPacketType = models.OracleConnectionDataMessageType
	}
	dataOffSet := binary.BigEndian.Uint16(packetData[26:])
	PacketLength := binary.BigEndian.Uint16(packetData[0:])
	if nextDataPacketType == models.DefaultDataPacket {
		connectionString = string(packetData[dataOffSet:])
	}
	requestHeader := models.OracleHeader{
		PacketLength: int(PacketLength),
		PacketType:   models.PacketTypeFromUint8(packetData[4]),
		PacketFlag:   packetData[5],
	}
	requestMessage := models.OracleConnectMessage{
		TnsVersionDesired:          binary.BigEndian.Uint16(packetData[8:]),
		TnsVersionMinimum:          binary.BigEndian.Uint16(packetData[10:]),
		ServiceOptions:             binary.BigEndian.Uint16(packetData[12:]),
		TNS_SDU_16:                 binary.BigEndian.Uint16(packetData[14:]),
		TNS_TDU_16:                 binary.BigEndian.Uint16(packetData[16:]),
		TnsProtocolCharacteristics: binary.BigEndian.Uint16(packetData[18:]),
		LineTurnaround:             binary.BigEndian.Uint16(packetData[20:]),
		OurOne:                     binary.BigEndian.Uint16(packetData[22:]),
		ConnectionStringLength:     connectionStringLength,
		OffsetOfConnectionData:     dataOffSet,
		ACFL0:                      packetData[32],
		ACFL1:                      packetData[33],
		TNS_SDU_32:                 binary.BigEndian.Uint32(packetData[58:]),
		TNS_TDU_32:                 binary.BigEndian.Uint32(packetData[62:]),
		ConnectString:              connectionString,
	}
	if requestMessage.OffsetOfConnectionData > 70 {
		requestMessage.ConnectFlag1 = binary.BigEndian.Uint32(packetData[66:])
		requestMessage.ConnectFlag2 = binary.BigEndian.Uint32(packetData[70:])
	}
	session.Context = &network.SessionContext{
		LoVersion:         requestMessage.TnsVersionMinimum,
		Options:           requestMessage.ServiceOptions,
		SessionDataUnit:   requestMessage.TNS_SDU_32,
		TransportDataUnit: requestMessage.TNS_TDU_32,
		ACFL0:             requestMessage.ACFL0,
		ACFL1:             requestMessage.ACFL1,
		OurOne:            requestMessage.OurOne,
	}
	if requestMessage.ConnectFlag2 == 1 {
		session.SupportOOB = true
	}
	requestHeader.Session = session
	if nextDataPacketType == models.DefaultDataPacket {
		return requestHeader, requestMessage, false, nextDataPacketType, nil, nil
	} else {
		return requestHeader, requestMessage, true, nextDataPacketType, nil, nil

	}
}

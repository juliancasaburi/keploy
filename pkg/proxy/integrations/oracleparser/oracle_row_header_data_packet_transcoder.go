package oracleparser

import (
	"encoding/binary"

	"github.com/sijms/go-ora/v2/network"
	"go.keploy.io/server/pkg/models"
)

func DecodeRowHeaderPacket(Packets [][]byte, dataPacketType models.DataPacketType) (models.OracleHeader, interface{}, bool, models.DataPacketType, error) {
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
	requestMessage := models.OracleDataMessage{
		DATA_OFFSET:       10,
		DATA_MESSAGE_TYPE: models.OracleRowHeaderDataMessage,
	}
	var oracleRowHeaderDataMessage models.OracleRowHeaderTypeDataMessage
	index := 11
	oracleRowHeaderDataMessage.FLAGS = packetData[index]
	index += 1
	oracleRowHeaderDataMessage.COLUMN_COUNT, index = session.GetInt(2, true, true, packetData, index)
	oracleRowHeaderDataMessage.ITERATION_NUM, index = session.GetInt(4, true, true, packetData, index)
	oracleRowHeaderDataMessage.ROW_COUNT, index = session.GetInt(4, true, true, packetData, index)
	oracleRowHeaderDataMessage.BUFFER_LENGTH, index = session.GetInt(2, true, true, packetData, index)
	oracleRowHeaderDataMessage.BIT_VECTOR, index = session.GetDlc(packetData, index)
	_, _ = session.GetDlc(packetData, index)
	requestMessage.DATA_MESSAGE = oracleRowHeaderDataMessage
	return requestHeader, requestMessage, false, models.Default, nil
}


func DecodeRowDataPacket(Packets [][]byte, dataPacketType models.DataPacketType) (models.OracleHeader, interface{}, bool, models.DataPacketType, error) {
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
	requestMessage := models.OracleDataMessage{
		DATA_OFFSET:       10,
		DATA_MESSAGE_TYPE: models.OracleRowHeaderDataMessage,
	}
	var oracleRowHeaderDataMessage models.OracleRowHeaderTypeDataMessage
	index := 11
	oracleRowHeaderDataMessage.FLAGS = packetData[index]
	index += 1
	oracleRowHeaderDataMessage.COLUMN_COUNT, index = session.GetInt(2, true, true, packetData, index)
	oracleRowHeaderDataMessage.ITERATION_NUM, index = session.GetInt(4, true, true, packetData, index)
	oracleRowHeaderDataMessage.ROW_COUNT, index = session.GetInt(4, true, true, packetData, index)
	oracleRowHeaderDataMessage.BUFFER_LENGTH, index = session.GetInt(2, true, true, packetData, index)
	oracleRowHeaderDataMessage.BIT_VECTOR, index = session.GetDlc(packetData, index)
	_, _ = session.GetDlc(packetData, index)
	requestMessage.DATA_MESSAGE = oracleRowHeaderDataMessage
	return requestHeader, requestMessage, false, models.Default, nil
}



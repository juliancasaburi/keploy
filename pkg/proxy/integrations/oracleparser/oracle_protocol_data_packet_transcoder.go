package oracleparser

import (
	"encoding/binary"

	"github.com/sijms/go-ora/v2/network"
	"go.keploy.io/server/pkg/models"
)

func DecodeOracleProtocolDataMessage(Packets [][]byte, isRequest bool) (models.OracleHeader, interface{}, bool, models.DataPacketType, error) {
	var packetData []byte
	var packetLength interface{}
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
	requestMessage := models.OracleDataMessage{
		DATA_OFFSET:       10,
		DATA_MESSAGE_TYPE: models.OracleProtocolDataMessageType,
	}
	if isRequest {
		index := 12
		arrayTerminator, index := session.GetNullTerminatedArray(packetData, index)
		driverName, _ := session.GetNullTerminatedString(packetData, index)
		requestMessage.DATA_MESSAGE = models.OracleProtocolDataMessageRequest{
			PROTOCOL_VERSION:      packetData[11],
			ARRAY_TERMINATOR_LIST: arrayTerminator,
			DRIVER_NAME:           driverName,
		}
		return requestHeader, requestMessage, false, models.Default, nil
	} else {
		oracleProtocolDataMessage := models.OracleProtocolDataMessageResponse{}
		oracleProtocolDataMessage.PROTOCOL_VERSION = packetData[11]
		index := 12
		oracleProtocolDataMessage.ARRAY_TERMINATOR_LIST, index = session.GetNullTerminatedArray(packetData, index)
		oracleProtocolDataMessage.PROTOCOL_SERVER_NAME, index = session.GetNullTerminatedString(packetData, index)
		oracleProtocolDataMessage.SERVER_CHARACTER_SET, index = session.GetInt(2, false, false, packetData, index)
		oracleProtocolDataMessage.SERVER_FLAGS = packetData[index]
		index += 1
		oracleProtocolDataMessage.CHARACTER_SET_ELEMENT, index = session.GetInt(2, false, false, packetData, index)
		if oracleProtocolDataMessage.CHARACTER_SET_ELEMENT > 0 {
			index = index + (oracleProtocolDataMessage.CHARACTER_SET_ELEMENT * 5)
		}
		oracleProtocolDataMessage.ARRAY_LENGTH, index = session.GetInt(2, false, true, packetData, index)
		oracleProtocolDataMessage.NUMBER_ARRAY = packetData[index : index+oracleProtocolDataMessage.ARRAY_LENGTH]
		index += oracleProtocolDataMessage.ARRAY_LENGTH
		oracleProtocolDataMessage.SERVER_COMPILE_TIME_CAPS_LENGHT = packetData[index]
		index += 1
		oracleProtocolDataMessage.SERVER_COMPILE_TIME_CAPS = packetData[index : index+int(oracleProtocolDataMessage.SERVER_COMPILE_TIME_CAPS_LENGHT)]
		index += int(oracleProtocolDataMessage.SERVER_COMPILE_TIME_CAPS_LENGHT)
		session.ServerCompileTimeCaps = oracleProtocolDataMessage.SERVER_COMPILE_TIME_CAPS
		oracleProtocolDataMessage.SERVER_RUN_TIME_CAPS_LENGTH = packetData[index]
		index += 1
		oracleProtocolDataMessage.SERVER_RUN_TIME_CAPS = packetData[index : index+int(oracleProtocolDataMessage.SERVER_RUN_TIME_CAPS_LENGTH)]
		index += int(oracleProtocolDataMessage.SERVER_COMPILE_TIME_CAPS_LENGHT)
		session.ServerRunTimeCaps = oracleProtocolDataMessage.SERVER_RUN_TIME_CAPS
		if oracleProtocolDataMessage.SERVER_COMPILE_TIME_CAPS[15]&1 != 0 {
			session.HasEOSCapability = true
		}
		if oracleProtocolDataMessage.SERVER_COMPILE_TIME_CAPS[16]&1 != 0 {
			session.HasFSAPCapability = true
		}
		if len(oracleProtocolDataMessage.SERVER_COMPILE_TIME_CAPS) > 37 && oracleProtocolDataMessage.SERVER_COMPILE_TIME_CAPS[37]&32 != 0 {
			session.UseBigClrChunks = true
			session.ClrChunkSize = 0x7FFF
		}
		requestHeader.Session = session
		requestMessage.DATA_MESSAGE = oracleProtocolDataMessage
		return requestHeader, requestMessage, false, models.Default, nil
	}
}

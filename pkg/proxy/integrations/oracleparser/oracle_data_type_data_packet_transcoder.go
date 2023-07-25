package oracleparser

import (
	"encoding/binary"

	"github.com/sijms/go-ora/v2/network"
	"go.keploy.io/server/pkg/models"
)

func DecodeOracleDataTypeDataMessage(Packets [][]byte, isRequest bool) (models.OracleHeader, interface{}, bool, models.DataPacketType, error) {
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
		DATA_MESSAGE_TYPE: models.OracleDataTypeDataMessageType,
	}
	if isRequest {
		oracleDataTypeDataMessage := models.OracleDataTypeDataMessageRequest{}
		index := 11
		oracleDataTypeDataMessage.SERVER_CHARACTER_SET, index = session.GetInt(2, false, false, packetData, index)
		oracleDataTypeDataMessage.SERVER_CHARACTER_SET, index = session.GetInt(2, false, false, packetData, index)
		oracleDataTypeDataMessage.SERVER_FLAGS = packetData[index]
		index += 1
		oracleDataTypeDataMessage.COMPILE_TIME_CAPS_LENGHT = packetData[index]
		index += 1
		oracleDataTypeDataMessage.COMPILE_TIME_CAPS = packetData[index : index+int(oracleDataTypeDataMessage.COMPILE_TIME_CAPS_LENGHT)]
		session.CompileTimeCaps = oracleDataTypeDataMessage.COMPILE_TIME_CAPS
		index += int(oracleDataTypeDataMessage.COMPILE_TIME_CAPS_LENGHT)
		oracleDataTypeDataMessage.RUN_TIME_CAPS_LENGTH = packetData[index]
		index += 1
		oracleDataTypeDataMessage.RUN_TIME_CAPS = packetData[index : index+int(oracleDataTypeDataMessage.RUN_TIME_CAPS_LENGTH)]
		session.RunTimeCaps = oracleDataTypeDataMessage.RUN_TIME_CAPS
		index += int(oracleDataTypeDataMessage.RUN_TIME_CAPS_LENGTH)
		if oracleDataTypeDataMessage.RUN_TIME_CAPS[1]&1 == 1 {
			index += 11
			if oracleDataTypeDataMessage.COMPILE_TIME_CAPS[37]&2 == 2 {
				oracleDataTypeDataMessage.CLIENT_TZ_VERSION, index = session.GetInt(4, false, true, packetData, index)
			}
		}
		oracleDataTypeDataMessage.SERVERN_CHARACTER_SET, index = session.GetInt(2, false, false, packetData, index)
		var dataTypeInfoList []models.DataTypeInfo
		if oracleDataTypeDataMessage.COMPILE_TIME_CAPS[27] != 0 {
			for {
				dataTypeInfo := models.DataTypeInfo{}
				dataType, newIndex := session.GetInt(2, false, true, packetData, index)
				index = newIndex
				if dataType == 0 {
					break
				}
				dataTypeInfo.DATA_TYPE = models.DataType(dataType)
				dataTypeConv, newIndex := session.GetInt(2, false, true, packetData, index)
				index = newIndex
				dataTypeInfo.CONV_DATA_TYPE = models.DataType(dataTypeConv)
				if dataTypeInfo.CONV_DATA_TYPE == 0 {
					dataTypeInfoList = append(dataTypeInfoList, dataTypeInfo)
					continue
				}
				representation, newIndex := session.GetInt(2, false, true, packetData, index)
				index = newIndex
				dataTypeInfo.REPRESENTATION = models.DataType(representation)
				_, newIndex = session.GetInt(2, false, true, packetData, index)
				index = newIndex
				dataTypeInfoList = append(dataTypeInfoList, dataTypeInfo)
			}
		} else {
			for {
				dataTypeInfo := models.DataTypeInfo{}
				dataTypeInfo.DATA_TYPE = models.DataType(uint8(packetData[index]))
				index += 1
				if dataTypeInfo.DATA_TYPE == 0 {
					break
				}
				dataTypeInfo.CONV_DATA_TYPE = models.DataType(uint8(packetData[index]))
				index += 1
				if dataTypeInfo.CONV_DATA_TYPE == 0 {
					dataTypeInfoList = append(dataTypeInfoList, dataTypeInfo)
					continue
				}
				dataTypeInfo.REPRESENTATION = models.DataType(uint8(packetData[index]))
				index += 1
				_ = uint8(packetData[index])
				index += 1
				dataTypeInfoList = append(dataTypeInfoList, dataTypeInfo)
			}
		}
		oracleDataTypeDataMessage.DATA_TYPE = dataTypeInfoList
		requestHeader.Session = session
		requestMessage.DATA_MESSAGE = oracleDataTypeDataMessage
		return requestHeader, requestMessage, false, models.Default, nil
	} else {
		oracleDataTypeDataMessage := models.OracleDataTypeDataMessageResponse{}
		index := 11
		if session.RunTimeCaps[1] == 1 {
			index += 11
			if session.CompileTimeCaps[37]&2 == 2 {
				oracleDataTypeDataMessage.SERVER_TZ_VERSION, index = session.GetInt(4, false, true, packetData, index)
			}
		}
		var dataTypeInfoList []models.DataTypeInfo
		if session.CompileTimeCaps[27] != 0 {
			for {
				dataTypeInfo := models.DataTypeInfo{}
				dataType, newIndex := session.GetInt(2, false, true, packetData, index)
				index = newIndex
				if dataType == 0 {
					break
				}
				dataTypeInfo.DATA_TYPE = models.DataType(dataType)
				convoDataType, newIndex := session.GetInt(2, false, true, packetData, index)
				index = newIndex
				dataTypeInfo.CONV_DATA_TYPE = models.DataType(convoDataType)
				if dataTypeInfo.CONV_DATA_TYPE == 0 {
					dataTypeInfoList = append(dataTypeInfoList, dataTypeInfo)
					continue
				}
				representation, newIndex := session.GetInt(2, false, true, packetData, index)
				index = newIndex
				dataTypeInfo.REPRESENTATION = models.DataType(representation)
				index += 2
				dataTypeInfoList = append(dataTypeInfoList, dataTypeInfo)
			}
		} else {
			for {
				dataTypeInfo := models.DataTypeInfo{}
				dataType, newIndex := session.GetInt(1, false, false, packetData, index)
				index = newIndex
				if dataType == 0 {
					break
				}
				dataTypeInfo.DATA_TYPE = models.DataType(dataType)
				convoDataType, newIndex := session.GetInt(1, false, false, packetData, index)
				index = newIndex
				dataTypeInfo.CONV_DATA_TYPE = models.DataType(convoDataType)
				if dataTypeInfo.CONV_DATA_TYPE == 0 {
					dataTypeInfoList = append(dataTypeInfoList, dataTypeInfo)
					continue
				}
				representation, newIndex := session.GetInt(1, false, false, packetData, index)
				index = newIndex
				dataTypeInfo.REPRESENTATION = models.DataType(representation)
				index += 1
				dataTypeInfoList = append(dataTypeInfoList, dataTypeInfo)
			}
		}
		oracleDataTypeDataMessage.DATA_TYPE = dataTypeInfoList
		requestHeader.Session = session
		requestMessage.DATA_MESSAGE = oracleDataTypeDataMessage
	}
	session.TTCVersion = session.CompileTimeCaps[7]
	session.UseBigScn = session.ServerCompileTimeCaps[7] >= 8
	if session.ServerCompileTimeCaps[7] < session.TTCVersion {
		session.TTCVersion = session.ServerCompileTimeCaps[7]
	}
	return requestHeader, requestMessage, false, models.Default, nil

}

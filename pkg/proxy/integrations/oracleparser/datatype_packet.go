package oracleparser

import (
	"encoding/binary"

	"go.keploy.io/server/pkg/models"
)

func DecodeOracleDataTypeDataMessage(Packets [][]byte, isRequest bool) (models.OracleHeader, interface{}, bool, models.DataPacketType, interface{}, error) {
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
	requestMessage := models.OracleDataMessage{
		DataOffset:      10,
		DataMessageType: models.OracleDataTypeDataMessageType,
	}
	if isRequest {
		oracleDataTypeDataMessage := models.OracleDataTypeDataMessageRequest{}
		index := 11
		oracleDataTypeDataMessage.ServerCharacterSet, index = session.GetInt(2, false, false, packetData, index)
		oracleDataTypeDataMessage.ServerCharacterSet, index = session.GetInt(2, false, false, packetData, index)
		oracleDataTypeDataMessage.ServerFlags = packetData[index]
		index += 1
		oracleDataTypeDataMessage.CompileTimeCapsLength = packetData[index]
		index += 1
		oracleDataTypeDataMessage.CompileTimeCaps = packetData[index : index+int(oracleDataTypeDataMessage.CompileTimeCapsLength)]
		session.CompileTimeCaps = oracleDataTypeDataMessage.CompileTimeCaps
		index += int(oracleDataTypeDataMessage.CompileTimeCapsLength)
		oracleDataTypeDataMessage.RunTimeCapsLength = packetData[index]
		index += 1
		oracleDataTypeDataMessage.RunTimeCaps = packetData[index : index+int(oracleDataTypeDataMessage.RunTimeCapsLength)]
		session.RunTimeCaps = oracleDataTypeDataMessage.RunTimeCaps
		index += int(oracleDataTypeDataMessage.RunTimeCapsLength)
		if oracleDataTypeDataMessage.RunTimeCaps[1]&1 == 1 {
			index += 11
			if oracleDataTypeDataMessage.CompileTimeCaps[37]&2 == 2 {
				oracleDataTypeDataMessage.ClientTZVersion, index = session.GetInt(4, false, true, packetData, index)
			}
		}
		oracleDataTypeDataMessage.ServernCharacterSet, index = session.GetInt(2, false, false, packetData, index)
		var dataTypeInfoList []models.DataTypeInfo
		if oracleDataTypeDataMessage.CompileTimeCaps[27] != 0 {
			for {
				dataTypeInfo := models.DataTypeInfo{}
				dataType, newIndex := session.GetInt(2, false, true, packetData, index)
				index = newIndex
				if dataType == 0 {
					break
				}
				dataTypeInfo.DataType = models.DataTypeFromInt(dataType)
				dataTypeConv, newIndex := session.GetInt(2, false, true, packetData, index)
				index = newIndex
				dataTypeInfo.ConvDataType = models.DataTypeFromInt(dataTypeConv)
				if dataTypeInfo.ConvDataType == models.DEFAULT_DATA_TYPE {
					dataTypeInfoList = append(dataTypeInfoList, dataTypeInfo)
					continue
				}
				representation, newIndex := session.GetInt(2, false, true, packetData, index)
				index = newIndex
				dataTypeInfo.Representation = models.DataTypeFromInt(representation)
				_, newIndex = session.GetInt(2, false, true, packetData, index)
				index = newIndex
				dataTypeInfoList = append(dataTypeInfoList, dataTypeInfo)
			}
		} else {
			for {
				dataTypeInfo := models.DataTypeInfo{}
				dataTypeInfo.DataType = models.DataTypeFromInt(int(packetData[index]))
				index += 1
				if dataTypeInfo.DataType == models.DEFAULT_DATA_TYPE {
					break
				}
				dataTypeInfo.ConvDataType = models.DataTypeFromInt(int(packetData[index]))
				index += 1
				if dataTypeInfo.ConvDataType == models.DEFAULT_DATA_TYPE {
					dataTypeInfoList = append(dataTypeInfoList, dataTypeInfo)
					continue
				}
				dataTypeInfo.Representation = models.DataTypeFromInt(int(packetData[index]))
				index += 1
				_ = uint8(packetData[index])
				index += 1
				dataTypeInfoList = append(dataTypeInfoList, dataTypeInfo)
			}
		}
		oracleDataTypeDataMessage.DataType = dataTypeInfoList
		requestHeader.Session = session
		requestMessage.DataMessage = oracleDataTypeDataMessage
		return requestHeader, requestMessage, false, models.DefaultDataPacket, nil, nil
	} else {
		oracleDataTypeDataMessage := models.OracleDataTypeDataMessageResponse{}
		index := 11
		if session.RunTimeCaps[1] == 1 {
			index += 11
			if session.CompileTimeCaps[37]&2 == 2 {
				oracleDataTypeDataMessage.ServerTZVersion, index = session.GetInt(4, false, true, packetData, index)
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
				dataTypeInfo.DataType = models.DataTypeFromInt(dataType)
				convoDataType, newIndex := session.GetInt(2, false, true, packetData, index)
				index = newIndex
				dataTypeInfo.ConvDataType = models.DataTypeFromInt(convoDataType)
				if dataTypeInfo.ConvDataType == models.DEFAULT_DATA_TYPE {
					dataTypeInfoList = append(dataTypeInfoList, dataTypeInfo)
					continue
				}
				representation, newIndex := session.GetInt(2, false, true, packetData, index)
				index = newIndex
				dataTypeInfo.Representation = models.DataTypeFromInt(representation)
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
				dataTypeInfo.DataType = models.DataTypeFromInt(dataType)
				convoDataType, newIndex := session.GetInt(1, false, false, packetData, index)
				index = newIndex
				dataTypeInfo.ConvDataType = models.DataTypeFromInt(convoDataType)
				if dataTypeInfo.ConvDataType == models.DEFAULT_DATA_TYPE {
					dataTypeInfoList = append(dataTypeInfoList, dataTypeInfo)
					continue
				}
				representation, newIndex := session.GetInt(1, false, false, packetData, index)
				index = newIndex
				dataTypeInfo.Representation = models.DataTypeFromInt(representation)
				index += 1
				dataTypeInfoList = append(dataTypeInfoList, dataTypeInfo)
			}
		}
		oracleDataTypeDataMessage.DataType = dataTypeInfoList
		requestHeader.Session = session
		requestMessage.DataMessage = oracleDataTypeDataMessage
	}
	session.TTCVersion = session.CompileTimeCaps[7]
	session.UseBigScn = session.ServerCompileTimeCaps[7] >= 8
	if session.ServerCompileTimeCaps[7] < session.TTCVersion {
		session.TTCVersion = session.ServerCompileTimeCaps[7]
	}
	return requestHeader, requestMessage, false, models.DefaultDataPacket, nil, nil

}

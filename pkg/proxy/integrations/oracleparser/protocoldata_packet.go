package oracleparser

import (
	"encoding/binary"

	"github.com/sijms/go-ora/v2/converters"
	"go.keploy.io/server/pkg/models"
)

func DecodeOracleProtocolDataMessage(Packets [][]byte, isRequest bool) (models.OracleHeader, interface{}, bool, models.DataPacketType, interface{}, error) {
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
		DataMessageType: models.OracleProtocolDataMessageType,
	}
	if isRequest {
		index := 12
		arrayTerminator, index := session.GetNullTerminatedArray(packetData, index)
		driverName, _ := session.GetNullTerminatedString(packetData, index)
		requestMessage.DataMessage = models.OracleProtocolDataMessageRequest{
			ProtocolVersion:     packetData[11],
			ArrayTerminatorList: arrayTerminator,
			DriveName:           driverName,
		}
		return requestHeader, requestMessage, false, models.DefaultDataPacket, nil, nil
	} else {
		oracleProtocolDataMessage := models.OracleProtocolDataMessageResponse{}
		oracleProtocolDataMessage.ProtocolVersion = packetData[11]
		index := 12
		oracleProtocolDataMessage.ArrayTerminatorList, index = session.GetNullTerminatedArray(packetData, index)
		oracleProtocolDataMessage.ProtocolServerName, index = session.GetNullTerminatedString(packetData, index)
		oracleProtocolDataMessage.ServerCharacterSet, index = session.GetInt(2, false, false, packetData, index)
		oracleProtocolDataMessage.ServerFlags = packetData[index]
		index += 1
		oracleProtocolDataMessage.CharacterSetElement, index = session.GetInt(2, false, false, packetData, index)
		if oracleProtocolDataMessage.CharacterSetElement > 0 {
			index = index + (oracleProtocolDataMessage.CharacterSetElement * 5)
		}
		oracleProtocolDataMessage.ArrayLength, index = session.GetInt(2, false, true, packetData, index)
		oracleProtocolDataMessage.NumberArray = packetData[index : index+oracleProtocolDataMessage.ArrayLength]
		index += oracleProtocolDataMessage.ArrayLength
		oracleProtocolDataMessage.ServerCompileTimeCapsLength = packetData[index]
		index += 1
		oracleProtocolDataMessage.ServerCompileTimeCaps = packetData[index : index+int(oracleProtocolDataMessage.ServerCompileTimeCapsLength)]
		index += int(oracleProtocolDataMessage.ServerCompileTimeCapsLength)
		session.ServerCompileTimeCaps = oracleProtocolDataMessage.ServerCompileTimeCaps
		oracleProtocolDataMessage.ServerRunTimeCapsLength = packetData[index]
		index += 1
		oracleProtocolDataMessage.ServerRunTimeCaps = packetData[index : index+int(oracleProtocolDataMessage.ServerRunTimeCapsLength)]
		index += int(oracleProtocolDataMessage.ServerCompileTimeCapsLength)
		session.ServerRunTimeCaps = oracleProtocolDataMessage.ServerRunTimeCaps
		if oracleProtocolDataMessage.ServerCompileTimeCaps[15]&1 != 0 {
			session.HasEOSCapability = true
		}
		if oracleProtocolDataMessage.ServerCompileTimeCaps[16]&1 != 0 {
			session.HasFSAPCapability = true
		}
		if len(oracleProtocolDataMessage.ServerCompileTimeCaps) > 37 && oracleProtocolDataMessage.ServerCompileTimeCaps[37]&32 != 0 {
			session.UseBigClrChunks = true
			session.ClrChunkSize = 0x7FFF
		}
		session.SStrConv = converters.NewStringConverter(oracleProtocolDataMessage.ServerCharacterSet)
		session.StrConv = session.SStrConv
		session.ServerCharacterSet = oracleProtocolDataMessage.ServerCharacterSet
		num := int(6 + (oracleProtocolDataMessage.NumberArray[5]) + (oracleProtocolDataMessage.NumberArray[6]))
		ServernCharset := int(binary.BigEndian.Uint16(oracleProtocolDataMessage.NumberArray[(num + 3):(num + 5)]))
		session.ServernCharacterSet = ServernCharset
		session.NStrConv = converters.NewStringConverter(ServernCharset)
		requestHeader.Session = session
		requestMessage.DataMessage = oracleProtocolDataMessage
		return requestHeader, requestMessage, false, models.DefaultDataPacket, nil, nil
	}
}

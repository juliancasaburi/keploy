package oracleparser

import (
	"encoding/binary"

	"github.com/sijms/go-ora/v2/network"
	"go.keploy.io/server/pkg/models"
)

func DecodeAcceptPacket(Packets [][]byte, dataPacketType models.DataPacketType) (models.OracleHeader, interface{}, bool, models.DataPacketType, interface{}, error) {
	var packetData []byte
	for _, slice := range Packets {
		packetData = append(packetData, slice...)
	}
	requestHeader := models.OracleHeader{
		PacketLength: int(binary.BigEndian.Uint16(packetData[0:])),
		PacketType:   models.PacketTypeFromUint8(packetData[4]),
		PacketFlag:   Packets[0][5],
	}
	requestMessage := models.OracleAcceptMessage{
		TnsVersion:      binary.BigEndian.Uint16(packetData[8:]),
		ProtocolOptions: binary.BigEndian.Uint16(packetData[10:]),
		SDU_16:          binary.BigEndian.Uint16(packetData[12:]),
		TDU_16:          binary.BigEndian.Uint16(packetData[14:]),
		HistOne:         binary.BigEndian.Uint16(packetData[16:]),
		DataLength:      binary.BigEndian.Uint16(packetData[18:]),
		DataOffset:      binary.BigEndian.Uint16(packetData[20:]),
		ACFL0:           packetData[22],
		ACFL1:           packetData[23],
		ReconAddStart:   binary.BigEndian.Uint16(packetData[28:]),
		ReconAddLength:  binary.BigEndian.Uint16(packetData[30:]),
	}
	requestMessage.Buffer = packetData[int(requestMessage.DataOffset):]
	SDU := uint32(requestMessage.SDU_16)
	TDU := uint32(requestMessage.TDU_16)
	if requestMessage.TnsVersion >= 315 {
		requestMessage.SDU_32 = binary.BigEndian.Uint32(packetData[32:])
		requestMessage.TDU_32 = binary.BigEndian.Uint32(packetData[36:])
		SDU = requestMessage.SDU_32
		TDU = requestMessage.TDU_32
	}
	reconAdd := ""
	if requestMessage.ReconAddStart != 0 && requestMessage.ReconAddLength != 0 && uint16(len(packetData)) > (requestMessage.ReconAddStart+requestMessage.ReconAddLength) {
		reconAdd = string(packetData[requestMessage.ReconAddStart:(requestMessage.ReconAddStart + requestMessage.ReconAddLength)])
	}
	requestMessage.ReconAdd = reconAdd
	if (requestHeader.PacketFlag & 1) > 0 {
		requestMessage.Sid = packetData[requestHeader.PacketLength:]
	}
	session.Context = &network.SessionContext{
		Version:           requestMessage.TnsVersion,
		LoVersion:         0,
		Options:           0,
		NegotiatedOptions: requestMessage.ProtocolOptions,
		SessionDataUnit:   SDU,
		TransportDataUnit: TDU,
		ACFL0:             requestMessage.ACFL0,
		ACFL1:             requestMessage.ACFL1,
		SID:               requestMessage.Sid,
		Histone:           requestMessage.HistOne,
		ReconAddr:         requestMessage.ReconAdd,
	}
	session.HandShakeComplete = true
	requestHeader.Session = session
	return requestHeader, requestMessage, false, dataPacketType, nil, nil
}

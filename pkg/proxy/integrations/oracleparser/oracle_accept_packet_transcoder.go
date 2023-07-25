package oracleparser

import (
	"encoding/binary"

	"github.com/sijms/go-ora/v2/network"
	"go.keploy.io/server/pkg/models"
)

func DecodeAcceptPacket(Packets [][]byte, dataPacketType models.DataPacketType) (models.OracleHeader, interface{}, bool, models.DataPacketType, error) {

	var packetData []byte
	for _, slice := range Packets {
		packetData = append(packetData, slice...)
	}
	requestHeader := models.OracleHeader{
		PacketLength: binary.BigEndian.Uint16(Packets[0][0:]),
		PacketType:   network.PacketType(Packets[0][4]),
		PacketFlag:   Packets[0][5],
	}
	requestMessage := models.OracleAcceptMessage{
		TNS_VERSION:      binary.BigEndian.Uint16(packetData[8:]),
		PROTOCOL_OPTIONS: binary.BigEndian.Uint16(packetData[10:]),
		SDU_16:           binary.BigEndian.Uint16(packetData[12:]),
		TDU_16:           binary.BigEndian.Uint16(packetData[14:]),
		HISTONE:          binary.BigEndian.Uint16(packetData[16:]),
		DATA_LENGTH:      binary.BigEndian.Uint16(packetData[18:]),
		DATA_OFFSET:      binary.BigEndian.Uint16(packetData[20:]),
		ACFL0:            packetData[22],
		ACFL1:            packetData[23],
		RECON_ADD_START:  binary.BigEndian.Uint16(packetData[28:]),
		RECON_ADD_LENGTH: binary.BigEndian.Uint16(packetData[30:]),
	}
	requestMessage.BUFFER = packetData[int(requestMessage.DATA_OFFSET):]
	SDU := uint32(requestMessage.SDU_16)
	TDU := uint32(requestMessage.TDU_16)
	if requestMessage.TNS_VERSION >= 315 {
		requestMessage.SDU_32 = binary.BigEndian.Uint32(packetData[32:])
		requestMessage.TDU_32 = binary.BigEndian.Uint32(packetData[36:])
		SDU = requestMessage.SDU_32
		TDU = requestMessage.TDU_32
	}
	reconAdd := ""
	if requestMessage.RECON_ADD_START != 0 && requestMessage.RECON_ADD_LENGTH != 0 && uint16(len(packetData)) > (requestMessage.RECON_ADD_START+requestMessage.RECON_ADD_LENGTH) {
		reconAdd = string(packetData[requestMessage.RECON_ADD_START:(requestMessage.RECON_ADD_START + requestMessage.RECON_ADD_LENGTH)])
	}
	requestMessage.RECON_ADD = reconAdd
	if (requestHeader.PacketFlag & 1) > 0 {
		requestMessage.SID = packetData[requestHeader.PacketLength.(int):]
	}
	session.Context = &network.SessionContext{
		Version:           requestMessage.TNS_VERSION,
		LoVersion:         0,
		Options:           0,
		NegotiatedOptions: requestMessage.PROTOCOL_OPTIONS,
		SessionDataUnit:   SDU,
		TransportDataUnit: TDU,
		ACFL0:             requestMessage.ACFL0,
		ACFL1:             requestMessage.ACFL1,
		SID:               requestMessage.SID,
		Histone:           requestMessage.HISTONE,
		ReconAddr:         requestMessage.RECON_ADD,
	}
	session.HandShakeComplete = true
	requestHeader.Session = session
	return requestHeader, requestMessage, false, dataPacketType, nil
}

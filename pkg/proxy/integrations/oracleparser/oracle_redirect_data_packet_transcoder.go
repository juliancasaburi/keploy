package oracleparser

import (
	"encoding/binary"
	"strings"

	"github.com/sijms/go-ora/v2/network"
	"go.keploy.io/server/pkg/models"
)

func DecodeRedirectDataMessage(Packets [][]byte) (models.OracleHeader, interface{}, bool, models.DataPacketType, error) {
	var packetData []byte
	var packetLength interface{}
	var redirectAddress string
	var redirectData string
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
	data := string(packetData[10:])
	length := strings.Index(data, "\x00")
	if requestHeader.PacketFlag&2 != 0 && length > 0 {
		redirectAddress = data[:length]
		redirectData = data[length:]
	} else {
		redirectAddress = data
	}
	oracleRedirectDataMessage := models.OracleRedirectDataMessage{
		REDIRECT_ADDRESS: redirectAddress,
		REDIRECT_DATA:    redirectData,
	}
	requestMessage := models.OracleDataMessage{
		DATA_OFFSET:       10,
		DATA_MESSAGE_TYPE: models.OracleRedirectDataMessageType,
		DATA_MESSAGE:      oracleRedirectDataMessage,
	}
	return requestHeader, requestMessage, false, models.Default, nil
}

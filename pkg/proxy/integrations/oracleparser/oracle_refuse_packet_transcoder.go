package oracleparser

import (
	"encoding/binary"
	"regexp"
	"strconv"
	"strings"

	"github.com/sijms/go-ora/v2/network"
	"go.keploy.io/server/pkg/models"
)

func DecodeRefusePacket(Packets [][]byte, dataPacketType models.DataPacketType) (models.OracleHeader, interface{}, bool, models.DataPacketType, error) {
	var packetData []byte
	var packetLength interface{}
	var message string
	for _, slice := range Packets {
		packetData = append(packetData, slice...)
	}
	if session.Context.Version >= 315 {
		packetLength = binary.BigEndian.Uint32(Packets[0][0:])
	} else {
		packetLength = binary.BigEndian.Uint16(Packets[0][0:])
	}
	requestHeader := models.OracleHeader{
		PacketLength: packetLength,
		PacketType:   network.PacketType(Packets[0][4]),
		PacketFlag:   Packets[0][5],
		Session:      session,
	}
	dataLength := binary.BigEndian.Uint16(packetData[10:])
	if uint16(len(packetData)) >= 12+dataLength {
		message = string(packetData[12 : 12+dataLength])
	}
	requestMessage := models.OracleRefuseMessage{
		DATA_OFFSET:   12,
		SYSTEM_REASON: packetData[9],
		USER_REASON:   packetData[8],
		DATA_LENGTH:   dataLength,
		DATA:          message,
	}
	requestMessage.ORACLE_ERROR.ErrCode, requestMessage.ORACLE_ERROR.ErrMsg = extractCode(requestMessage.DATA)
	return requestHeader, requestMessage, false, dataPacketType, nil
}

func extractCode(data string) (errCode int, errMsg string) {

	errCode = 12564
	errMsg = "ORA-12564: TNS connection refused"
	if len(data) == 0 {
		return
	}
	r, err := regexp.Compile(`\(\s*ERR\s*=\s*([0-9]+)\s*\)`)
	if err != nil {
		return
	}
	msg := strings.ToUpper(data)
	matches := r.FindStringSubmatch(msg)
	if len(matches) != 2 {
		return
	}
	strErrCode := matches[1]
	ErrCode, err := strconv.ParseInt(strErrCode, 10, 32)
	if err == nil {
		errCode = int(ErrCode)
		oracleError := models.OracleError{ErrCode: errCode}
		errMsg = oracleError.Error()
	}
	r, err = regexp.Compile(`\(\s*ERROR\s*=([A-Z0-9=\(\)]+)`)
	if err != nil {
		return
	}
	matches = r.FindStringSubmatch(msg)
	if len(matches) != 2 {
		return
	}
	codeStr := matches[1]
	r, err = regexp.Compile(`CODE\s*=\s*([0-9]+)`)
	if err != nil {
		return
	}
	matches = r.FindStringSubmatch(codeStr)
	if len(matches) != 2 {
		return
	}
	strErrCode = matches[1]
	ErrCode, err = strconv.ParseInt(strErrCode, 10, 32)
	if err == nil {
		errCode = int(ErrCode)
		oracleError := models.OracleError{ErrCode: errCode}
		errMsg = oracleError.Error()
	}
	return
}

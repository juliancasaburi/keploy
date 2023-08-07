package oracleparser

import (
	"encoding/binary"

	"fmt"

	"go.keploy.io/server/pkg/models"
)

func DecodeOraclePiggyBackDataMessage(Packets [][]byte) (models.OracleHeader, interface{}, bool, models.DataPacketType, interface{}, error) {
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
		DataMessageType: models.OracleFunctionDataMesssageType,
	}
	oraclePiggyBackMsg := models.PiggyBackMsg{
		PiggyBackCode:  models.FunctionTypeFromInt(int(packetData[11])),
		SequenceNumber: packetData[12],
	}
	index := 13
	var piggyBackMsg interface{}
	switch oraclePiggyBackMsg.PiggyBackCode {
	case models.TNS_FUNC_SET_SCHEMA:
		fmt.Println("SchemaPiggyBackMsg")
		piggyBackMsg, index = DecodePiggyBackSetSchema(packetData, index)
	case models.TNS_FUNC_CLOSE_CURSORS:
		fmt.Println("CloseCursorPiggyBackMsg")
		piggyBackMsg, index = DecodePiggyCloseCursor(packetData, index)
	case models.TNS_FUNC_LOB_OP:
		fmt.Println("LOBPiggyBackMsg")
		piggyBackMsg, index = DecodePiggyBackLOB(packetData, index)
	case models.TNS_FUNC_SET_END_TO_END_ATTR:
		fmt.Println("EndToEndPiggyBackMsg")
		piggyBackMsg, index = DecodePiggyBackEndToEnd(packetData, index)
	}
	oraclePiggyBackMsg.PiggyBackData = piggyBackMsg
	oracleFunctionTypeDataMessage := models.OracleFuntionTypeDataMessage{
		FunctionCode:   models.FunctionTypeFromInt(int(packetData[index+1])),
		SequenceNumber: packetData[index+2],
	}
	index += 3
	var funtionData interface{}
	var stmt interface{}
	nextDataPacketType := models.DefaultDataPacket
	switch oracleFunctionTypeDataMessage.FunctionCode {
	case models.TNS_FUNC_PING:
		fmt.Println("FUNC_PING")
		nextDataPacketType = models.OracleMessageWithDataMessageType
	case models.TNS_FUNC_COMMIT:
		fmt.Println("FUNC_COMMIT")
		nextDataPacketType = models.OracleMessageWithDataMessageType
	case models.TNS_FUNC_FETCH:
		fmt.Println("FUNC_FETCH")
		funtionData, stmt = DecodeFetchFunctionMessage(packetData)
		nextDataPacketType = models.OracleMessageWithDataMessageType
	case models.TNS_FUNC_ROLLBACK:
		fmt.Println("FUNC_ROLLBACK")
		nextDataPacketType = models.OracleMessageWithDataMessageType
	case models.TNS_FUNC_LOGOFF:
		fmt.Println("FUNC_LOGOFF")
		nextDataPacketType = models.OracleMessageWithDataMessageType
	case models.TNS_FUNC_LOB_OP:
		fmt.Println("FUNC_LOB_OP")
		DecodeLOBOPFunctionMessage(packetData)
		nextDataPacketType = models.OracleMessageWithDataMessageType
	case models.TNS_FUNC_AUTH_PHASE_ONE:
		fmt.Println("TNS_FUNC_AUTH_PHASE_ONE")
		funtionData = DecodeAuthMessage(packetData)
		nextDataPacketType = models.OracleAuthPhaseOneDataMessageType
	case models.TNS_FUNC_AUTH_PHASE_TWO:
		fmt.Println("TNS_FUNC_AUTH_PHASE_TWO")
		funtionData = DecodeAuthMessage(packetData)
		nextDataPacketType = models.OracleAuthPhaseTwoDataMessageType
	case models.TNS_FUNC_EXECUTE:
		fmt.Println("FUNC_EXECUTE")
		funtionData, stmt = DecodeExecuteFunctionMessage(packetData, 13)
		nextDataPacketType = models.OracleMessageWithDataMessageType
	case models.TNS_FUNC_REEXECUTE:
		fmt.Println("FUNC_REEXECUTE")
		funtionData, stmt = DecodeReExecuteFunctionMessage(packetData, 13)
		nextDataPacketType = models.OracleMessageWithDataMessageType
	case models.TNS_FUNC_REEXECUTE_AND_FETCH:
		fmt.Println("FUNC_REEXECUTE_AND_FETCH")
		funtionData, stmt = DecodeReExecuteAndFetchFunctionMessage(packetData, 13)
		nextDataPacketType = models.OracleMessageWithDataMessageType
	case models.TNS_FUNC_GET_DB_VERSION:
		fmt.Println("FUNC_GET_DB_VERSION")
		nextDataPacketType = models.OracleGetDBVersionDataMessageType
	}
	var piggyAndFunctionData []interface{}
	piggyAndFunctionData = append(piggyAndFunctionData, oraclePiggyBackMsg)
	piggyAndFunctionData = append(piggyAndFunctionData, funtionData)
	oracleFunctionTypeDataMessage.FunctionData = piggyAndFunctionData
	requestMessage.DataMessage = oracleFunctionTypeDataMessage
	return requestHeader, requestMessage, false, nextDataPacketType, stmt, nil
}

func DecodePiggyBackSetSchema(packetData []byte, index int) (models.SchemaPiggyBackMsg, int) {
	piggyBackMsg := models.SchemaPiggyBackMsg{}
	_, index = session.GetByte(packetData, index)
	_, index = session.GetInt(4, true, true, packetData, index)
	len := packetData[index]
	index += 1
	piggyBackMsg.SchemaBytes = packetData[index : index+int(len)]
	index += int(len)
	return piggyBackMsg, index
}

func DecodePiggyCloseCursor(packetData []byte, index int) (models.CloseCursorPiggyBackMsg, int) {
	piggyBackMsg := models.CloseCursorPiggyBackMsg{}
	_, index = session.GetByte(packetData, index)
	var num int
	num, index = session.GetInt(4, true, true, packetData, index)
	var cursorIdList []int
	var cursorId int
	for i := 0; i < num; i++ {
		cursorId, index = session.GetInt(4, true, true, packetData, index)
		cursorIdList = append(cursorIdList, cursorId)
		delete(stmtMap, cursorId)
	}
	piggyBackMsg.CursorIds = cursorIdList
	return piggyBackMsg, index
}

func DecodePiggyBackLOB(packetData []byte, index int) (models.CloseTempLobsPiggyBackMsg, int) {
	piggyBackMsg := models.CloseTempLobsPiggyBackMsg{}
	index += 1
	piggyBackMsg.LobSize, index = session.GetInt(4, true, true, packetData, index)
	index += 1
	_, index = session.GetInt(4, true, true, packetData, index)
	_, index = session.GetInt(4, true, true, packetData, index)
	_, index = session.GetInt(4, true, true, packetData, index)
	index += 3
	piggyBackMsg.OpCode, index = session.GetInt(4, true, true, packetData, index)
	index += 3
	_, index = session.GetInt(4, true, true, packetData, index)
	_, index = session.GetInt(8, true, true, packetData, index)
	_, index = session.GetInt(8, true, true, packetData, index)
	index += 2
	_, index = session.GetInt(4, true, true, packetData, index)
	index += 1
	_, index = session.GetInt(4, true, true, packetData, index)
	index += 1
	_, index = session.GetInt(4, true, true, packetData, index)
	var LobList [][]byte
	// TODO: Captue LOB_TO_CLOSE but we don't know the length of each LOB_TO_CLOSE
	piggyBackMsg.LobToClose = LobList
	return piggyBackMsg, index
}

func DecodePiggyBackEndToEnd(packetData []byte, index int) (models.EndToEndPiggyBackMsg, int) {
	piggyBackMsg := models.EndToEndPiggyBackMsg{}
	index += 2
	piggyBackMsg.Flags, index = session.GetInt(4, true, true, packetData, index)

	if packetData[index] == 1 {
		piggyBackMsg.ClientIdentifierModified = true
		index += 1
		var num int
		num, index = session.GetInt(4, true, true, packetData, index)
		if num != 0 {
			piggyBackMsg.ClientIdentifier = true
		}
		index += 1
	} else {
		index += 1
		_, index = session.GetInt(4, true, true, packetData, index)
	}

	if packetData[index] == 1 {
		piggyBackMsg.ModuleModified = true
		index += 1
		var num int
		num, index = session.GetInt(4, true, true, packetData, index)
		if num != 0 {
			piggyBackMsg.Module = true
		}
		index += 1
	} else {
		index += 1
		_, index = session.GetInt(4, true, true, packetData, index)
	}

	if packetData[index] == 1 {
		piggyBackMsg.ActionModified = true
		index += 1
		var num int
		num, index = session.GetInt(4, true, true, packetData, index)
		if num != 0 {
			piggyBackMsg.Action = true
		}
		index += 1
	} else {
		index += 1
		_, index = session.GetInt(4, true, true, packetData, index)
	}

	index += 1
	_, index = session.GetInt(4, true, true, packetData, index)
	index += 1
	_, index = session.GetInt(4, true, true, packetData, index)

	if packetData[index] == 1 {
		piggyBackMsg.ClientInfoModified = true
		index += 1
		var num int
		num, index = session.GetInt(4, true, true, packetData, index)
		if num != 0 {
			piggyBackMsg.ClientInfo = true
		}
		index += 1
	} else {
		index += 1
		_, index = session.GetInt(4, true, true, packetData, index)
	}

	index += 1
	_, index = session.GetInt(4, true, true, packetData, index)
	index += 1
	_, index = session.GetInt(4, true, true, packetData, index)

	if packetData[index] == 1 {
		piggyBackMsg.DbopModified = true
		index += 1
		var num int
		num, index = session.GetInt(4, true, true, packetData, index)
		if num != 0 {
			piggyBackMsg.Dbop = true
		}
		index += 1
	} else {
		index += 1
		_, index = session.GetInt(4, true, true, packetData, index)
	}

	if piggyBackMsg.ClientIdentifierModified && piggyBackMsg.ClientIdentifier {
		len := packetData[index]
		index += 1
		piggyBackMsg.ClientIdentifierBytes = packetData[index : index+int(len)]
		index += int(len)
	}
	if piggyBackMsg.ModuleModified && piggyBackMsg.Module {
		len := packetData[index]
		index += 1
		piggyBackMsg.ModuleBytes = packetData[index : index+int(len)]
		index += int(len)
	}
	if piggyBackMsg.ActionModified && piggyBackMsg.Action {
		len := packetData[index]
		index += 1
		piggyBackMsg.ActionBytes = packetData[index : index+int(len)]
		index += int(len)
	}
	if piggyBackMsg.ClientInfoModified && piggyBackMsg.ClientInfo {
		len := packetData[index]
		index += 1
		piggyBackMsg.ClientInfoBytes = packetData[index : index+int(len)]
		index += int(len)
	}
	if piggyBackMsg.DbopModified && piggyBackMsg.Dbop {
		len := packetData[index]
		index += 1
		piggyBackMsg.DbopBytes = packetData[index : index+int(len)]
		index += int(len)
	}
	return piggyBackMsg, index
}

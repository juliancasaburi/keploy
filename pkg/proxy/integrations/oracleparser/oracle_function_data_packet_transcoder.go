package oracleparser

import (
	"encoding/binary"
	"fmt"

	"github.com/sijms/go-ora/v2/network"
	"go.keploy.io/server/pkg/models"
)

func DecodeOracleFunctionDataMessage(Packets [][]byte) (models.OracleHeader, interface{}, bool, models.DataPacketType, error) {
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
		DATA_MESSAGE_TYPE: models.OracleFunctionDataMesssageType,
	}
	oracleFunctionTypeDataMessage := models.OracleFuntionTypeDataMessage{
		FUNCTION_CODE:   models.FunctionType(packetData[11]),
		SEQUENCE_NUMBER: packetData[12],
	}
	var funtionData interface{}
	switch models.FunctionType(Packets[0][11]) {
	case models.TNS_FUNC_EXECUTE:
		fmt.Println("FUNC_EXECUTE")
		funtionData = DecodeExecuteFunctionMessage(packetData)
	case models.TNS_FUNC_PING:
		fmt.Println("FUNC_PING")
		funtionData = nil
	case models.TNS_FUNC_COMMIT:
		fmt.Println("FUNC_COMMIT")
		funtionData = nil
	case models.TNS_FUNC_FETCH:
		fmt.Println("FUNC_FETCH")
		funtionData = DecodeFetchFunctionMessage(packetData)
	case models.TNS_FUNC_ROLLBACK:
		fmt.Println("FUNC_ROLLBACK")
	case models.TNS_FUNC_LOGOFF:
		fmt.Println("FUNC_ROLLBACK")
	case models.TNS_FUNC_LOB_OP:
		fmt.Println("FUNC_LOB_OP")
		DecodeLOBOPFunctionMessage(packetData)
	case models.TNS_FUNC_REEXECUTE:
		fmt.Println("FUNC_REEXECUTE")
	case models.TNS_FUNC_REEXECUTE_AND_FETCH:
		fmt.Println("FUNC_REEXECUTE_FETCH")
	}
	oracleFunctionTypeDataMessage.FUNCTION_DATA = funtionData
	requestMessage.DATA_MESSAGE = oracleFunctionTypeDataMessage
	return requestHeader, requestMessage, false, models.Default, nil
}

func DecodeFetchFunctionMessage(packetData []byte) models.OracleFetchFunctionTypeDataMessage {
	oracleFetchFunction := models.OracleFetchFunctionTypeDataMessage{}
	index := 13
	oracleFetchFunction.CURSOR_ID, index = session.GetInt(2, true, true, packetData, index)
	oracleFetchFunction.ROWS_TO_FETCH, _ = session.GetInt(2, true, true, packetData, index)
	return oracleFetchFunction
}

func DecodeLOBOPFunctionMessage(packetData []byte) models.OracleLOBFunctionTypeDataMessage {
	oracleLOB := models.OracleLOBFunctionTypeDataMessage{}
	index := 13
	index += 1
	oracleLOB.SOURCE_LENGTH, index = session.GetInt(4, true, true, packetData, index)
	index += 1
	oracleLOB.DEST_LENGTH, index = session.GetInt(4, true, true, packetData, index)

	if session.TTCVersion < 3 {
		oracleLOB.SOURCE_OFF_SET, index = session.GetInt(4, true, true, packetData, index)
		oracleLOB.DEST_OFF_SET, index = session.GetInt(4, true, true, packetData, index)
	} else {
		index += 2
	}
	oracleLOB.IS_CHARSET_ID = packetData[index]
	index += 1
	if session.TTCVersion < 3 {
		if packetData[index] == 1 {
			oracleLOB.SEND_SIZE = true
		}
	}
	index += 1
	if packetData[index] == 1 {
		oracleLOB.BNullO2U = true
	}
	index += 1
	oracleLOB.OPERATION_ID, index = session.GetInt(4, true, true, packetData, index)
	index += 1
	oracleLOB.SCN_LENGTH, index = session.GetInt(4, true, true, packetData, index)
	if session.TTCVersion >= 3 {
		oracleLOB.SOURCE_OFF_SET, index = session.GetInt(8, true, true, packetData, index)
		oracleLOB.DEST_OFF_SET, index = session.GetInt(8, true, true, packetData, index)
		if packetData[index] == 1 {
			oracleLOB.SEND_SIZE = true
		}
		index += 1
	}
	if session.TTCVersion >= 4 {
		index += 1
	}

	if oracleLOB.SOURCE_LENGTH > 0 {
		oracleLOB.SOURCE_LOCATOR = packetData[index : index+oracleLOB.SOURCE_LENGTH]
		index += oracleLOB.SOURCE_LENGTH
	}

	if oracleLOB.DEST_LENGTH > 0 {
		oracleLOB.DEST_LOCATOR = packetData[index : index+oracleLOB.DEST_LENGTH]
		index += oracleLOB.DEST_LENGTH
	}

	if oracleLOB.IS_CHARSET_ID == 1 {
		oracleLOB.CHARSET_ID, index = session.GetInt(2, true, true, packetData, index)
	}
	if session.TTCVersion < 3 && oracleLOB.SEND_SIZE {
		oracleLOB.SIZE, index = session.GetInt(4, true, true, packetData, index)
	}
	for x := 0; x < oracleLOB.SCN_LENGTH; x++ {
		var scn int
		scn, index = session.GetInt(4, true, true, packetData, index)
		oracleLOB.SCN = append(oracleLOB.SCN, scn)
	}
	if session.TTCVersion >= 3 && oracleLOB.SEND_SIZE {
		oracleLOB.SIZE, _ = session.GetInt(8, true, true, packetData, index)
	}
	return oracleLOB
}

func DecodeExecuteFunctionMessage(packetData []byte) models.OracleQuery {

	oracleQuery := models.OracleQuery{}
	index := 13
	oracleQuery.OPTIONS, index = session.GetInt(4, true, true, packetData, index)
	oracleQuery.CURSOR_ID, index = session.GetInt(2, true, true, packetData, index)
	if oracleQuery.CURSOR_ID != 0 && packetData[index] == 1 {
		oracleQuery.IS_DDL = 1
	}
	index += 1
	oracleQuery.SQL_LENGTH, index = session.GetInt(4, true, true, packetData, index)
	index += 1
	oracleQuery.ARRAY_LENGTH, index = session.GetInt(2, true, true, packetData, index)
	index += 2
	index += 1
	oracleQuery.ROWS_TO_FETCH, index = session.GetInt(4, true, true, packetData, index)
	oracleQuery.LOB, index = session.GetInt(4, true, true, packetData, index)
	index += 1
	oracleQuery.NUM_PARAMS, index = session.GetInt(2, true, true, packetData, index)
	index += 5
	index += 1
	oracleQuery.NUM_DEFINES, index = session.GetInt(2, true, true, packetData, index)
	if session.TTCVersion >= 4 {
		index += 3
	}
	if session.TTCVersion >= 5 {
		index += 5
	}
	if session.TTCVersion >= 7 {
		index += 1
		oracleQuery.NUM_EXEC, index = session.GetInt(4, true, true, packetData, index)
		index += 1
	}
	if session.TTCVersion >= 8 {
		index += 5
	}
	if session.TTCVersion >= 9 {
		index += 2
	}
	if oracleQuery.CURSOR_ID == 0 || oracleQuery.IS_DDL == 1 {
		oracleQuery.SQL_BYTES, index = GetClr(packetData, index)
	}
	fmt.Println(oracleQuery.SQL_BYTES)
	fmt.Println(string(oracleQuery.SQL_BYTES))
	var al8i4 int
	for x := 0; x < 13; x++ {
		fmt.Println(packetData[index])
		al8i4, index = session.GetInt(4, true, true, packetData, index)
		fmt.Println(al8i4)
		oracleQuery.Al8i4 = append(oracleQuery.Al8i4, al8i4)
		fmt.Println(oracleQuery.Al8i4)
	}
	if oracleQuery.NUM_DEFINES > 0 {
		var defineData []models.DefineData
		for x := 0; x < oracleQuery.NUM_DEFINES; x++ {
			data, newIndex := ReadDefineData(packetData, index)
			index = newIndex
			defineData = append(defineData, data)
		}
		oracleQuery.DefineDataList = defineData
	} else {
		if oracleQuery.NUM_PARAMS > 0 {
			var defineData []models.DefineData
			for x := 0; x < oracleQuery.NUM_PARAMS; x++ {
				data, newIndex := ReadDefineData(packetData, index)
				index = newIndex
				defineData = append(defineData, data)
			}
			oracleQuery.DefineDataList = defineData
		}
	}
	return oracleQuery
}

func ReadDefineData(buffer []byte, index int) (models.DefineData, int) {
	var defineData models.DefineData
	defineData.DataType, index = session.GetByte(buffer, index)
	defineData.Flag, index = session.GetByte(buffer, index)
	defineData.Precision, index = session.GetByte(buffer, index)
	defineData.Scale, index = session.GetByte(buffer, index)
	defineData.MaxLen, index = session.GetInt(4, true, true, buffer, index)
	defineData.MaxNoOfArrayElements, index = session.GetInt(4, true, true, buffer, index)
	if session.TTCVersion >= 10 {
		defineData.ContFlag, index = session.GetInt(8, true, true, buffer, index)
	} else {
		defineData.ContFlag, index = session.GetInt(4, true, true, buffer, index)
	}
	if buffer[index] != 0 {
		_, index = session.GetInt(4, true, true, buffer, index)
		defineData.ToID, index = session.GetClr(buffer, index)
	} else {
		index += 1
	}
	defineData.Version, index = session.GetInt(2, true, true, buffer, index)
	defineData.CharsetID, index = session.GetInt(2, true, true, buffer, index)
	defineData.CharsetForm, index = session.GetByte(buffer, index)
	defineData.MaxCharLen, index = session.GetInt(4, true, true, buffer, index)
	if session.TTCVersion >= 8 {
		defineData.Oaccollid, index = session.GetInt(4, true, true, buffer, index)
	}
	return defineData, index
}

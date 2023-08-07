package oracleparser

import (
	"encoding/binary"
	"fmt"

	"go.keploy.io/server/pkg/models"
)

var stmtMap map[int]models.Stmt

func DecodeOracleFunctionDataMessage(Packets [][]byte) (models.OracleHeader, interface{}, bool, models.DataPacketType, interface{}, error) {
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
	oracleFunctionTypeDataMessage := models.OracleFuntionTypeDataMessage{
		FunctionCode:   models.FunctionTypeFromInt(int(packetData[11])),
		SequenceNumber: packetData[12],
	}
	var funtionData interface{}
	var stmt interface{}
	nextDataPacketType := models.DefaultDataPacket
	switch models.FunctionTypeFromInt(int(packetData[11])) {
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
	oracleFunctionTypeDataMessage.FunctionData = funtionData
	requestMessage.DataMessage = oracleFunctionTypeDataMessage
	return requestHeader, requestMessage, false, nextDataPacketType, stmt, nil
}

func DecodeFetchFunctionMessage(packetData []byte) (models.OracleFetchFunctionTypeDataMessage, models.Stmt) {
	oracleFetchFunction := models.OracleFetchFunctionTypeDataMessage{}
	index := 13
	oracleFetchFunction.CursorId, index = session.GetInt(2, true, true, packetData, index)
	oracleFetchFunction.RowsToFetch, _ = session.GetInt(2, true, true, packetData, index)
	oracleFetchFunction.Stmt = stmtMap[oracleFetchFunction.CursorId]
	return oracleFetchFunction, oracleFetchFunction.Stmt
}

func DecodeLOBOPFunctionMessage(packetData []byte) models.OracleLOBFunctionTypeDataMessage {
	oracleLOB := models.OracleLOBFunctionTypeDataMessage{}
	index := 13
	index += 1
	oracleLOB.SourceLength, index = session.GetInt(4, true, true, packetData, index)
	index += 1
	oracleLOB.DestLength, index = session.GetInt(4, true, true, packetData, index)

	if session.TTCVersion < 3 {
		oracleLOB.SourceOffSet, index = session.GetInt(4, true, true, packetData, index)
		oracleLOB.DestOffSet, index = session.GetInt(4, true, true, packetData, index)
	} else {
		index += 2
	}
	oracleLOB.IsCharsetId = packetData[index]
	index += 1
	if session.TTCVersion < 3 {
		if packetData[index] == 1 {
			oracleLOB.SendSize = true
		}
	}
	index += 1
	if packetData[index] == 1 {
		oracleLOB.BNullO2U = true
	}
	index += 1
	oracleLOB.OperationId, index = session.GetInt(4, true, true, packetData, index)
	index += 1
	oracleLOB.SCNLength, index = session.GetInt(4, true, true, packetData, index)
	if session.TTCVersion >= 3 {
		oracleLOB.SourceOffSet, index = session.GetInt(8, true, true, packetData, index)
		oracleLOB.DestOffSet, index = session.GetInt(8, true, true, packetData, index)
		if packetData[index] == 1 {
			oracleLOB.SendSize = true
		}
		index += 1
	}
	if session.TTCVersion >= 4 {
		index += 1
	}
	if oracleLOB.SourceLength > 0 {
		oracleLOB.SourceLocator = packetData[index : index+oracleLOB.SourceLength]
		index += oracleLOB.SourceLength
	}
	if oracleLOB.DestLength > 0 {
		oracleLOB.DestLocator = packetData[index : index+oracleLOB.DestLength]
		index += oracleLOB.DestLength
	}
	if oracleLOB.IsCharsetId == 1 {
		oracleLOB.CharsetId, index = session.GetInt(2, true, true, packetData, index)
	}
	if session.TTCVersion < 3 && oracleLOB.SendSize {
		oracleLOB.Size, index = session.GetInt(4, true, true, packetData, index)
	}
	for x := 0; x < oracleLOB.SCNLength; x++ {
		var scn int
		scn, index = session.GetInt(4, true, true, packetData, index)
		oracleLOB.SCN = append(oracleLOB.SCN, scn)
	}
	if session.TTCVersion >= 3 && oracleLOB.SendSize {
		oracleLOB.Size, _ = session.GetInt(8, true, true, packetData, index)
	}
	return oracleLOB
}

func DecodeExecuteFunctionMessage(packetData []byte, index int) (models.OracleExecuteMessage, models.Stmt) {
	oracleExecuteMessage := models.OracleExecuteMessage{}
	oracleExecuteMessage.ExeOp, index = session.GetInt(4, true, true, packetData, index)
	oracleExecuteMessage.Stmt.CursorID, index = session.GetInt(2, true, true, packetData, index)
	index += 1
	if packetData[index] != 0 {
		oracleExecuteMessage.Parse = true
	}
	if oracleExecuteMessage.Parse {
		_, index = session.GetInt(4, true, true, packetData, index)
	} else {
		index += 1
	}
	index += 1
	_, index = session.GetInt(2, true, true, packetData, index)
	index += 2
	index += 1
	oracleExecuteMessage.Stmt.NoOfRowsToFetch, index = session.GetInt(4, true, true, packetData, index)
	oracleExecuteMessage.LOB, index = session.GetInt(4, true, true, packetData, index)
	index += 1
	oracleExecuteMessage.NumPars, index = session.GetInt(2, true, true, packetData, index)
	index += 5
	if packetData[index] != 0 {
		oracleExecuteMessage.Define = true
	}
	index += 1
	oracleExecuteMessage.NumCols, index = session.GetInt(2, true, true, packetData, index)
	if session.TTCVersion >= 4 {
		index += 3
	}
	if session.TTCVersion >= 5 {
		index += 5
	}
	if session.TTCVersion >= 7 {
		index += 1
		oracleExecuteMessage.Stmt.ArrayBindCount, index = session.GetInt(4, true, true, packetData, index)
		index += 1
	}
	if session.TTCVersion >= 8 {
		index += 5
	}
	if session.TTCVersion >= 9 {
		index += 2
	}
	if oracleExecuteMessage.Parse {
		var textBytes []byte
		textBytes, index = GetClr(packetData, index)
		oracleExecuteMessage.Stmt.Text = string(textBytes)
	}
	var al8i4 int
	for x := 0; x < 13; x++ {
		al8i4, index = session.GetInt(4, true, true, packetData, index)
		oracleExecuteMessage.Al8i4 = append(oracleExecuteMessage.Al8i4, al8i4)
	}
	if oracleExecuteMessage.Define {
		var colData []models.ParameterInfo
		for x := 0; x < oracleExecuteMessage.NumCols; x++ {
			data, newIndex := ReadParamData(packetData, index)
			index = newIndex
			colData = append(colData, data)
		}
		oracleExecuteMessage.Stmt.Columns = colData
	} else {
		var paramData []models.ParameterInfo
		for x := 0; x < oracleExecuteMessage.NumPars; x++ {
			data, newIndex := ReadParamData(packetData, index)
			index = newIndex
			paramData = append(paramData, data)
		}
		oracleExecuteMessage.Stmt.Pars = paramData
	}
	stmt := models.NewStmt(oracleExecuteMessage.Stmt.Text)
	stmt.CursorID = oracleExecuteMessage.Stmt.CursorID
	stmt.NoOfRowsToFetch = oracleExecuteMessage.Stmt.NoOfRowsToFetch
	stmt.ArrayBindCount = oracleExecuteMessage.Stmt.ArrayBindCount
	stmt.Columns = oracleExecuteMessage.Stmt.Columns
	stmt.Pars = oracleExecuteMessage.Stmt.Pars
	oracleExecuteMessage.Stmt = *stmt
	if oracleExecuteMessage.Stmt.ArrayBindCount > 0 {
		oracleExecuteMessage.Stmt.Pars, _ = GetValuesIntoParsInBulk(oracleExecuteMessage.Stmt.Pars, packetData, index, oracleExecuteMessage.Stmt)
	} else {
		oracleExecuteMessage.Stmt.Pars, _ = GetValuesIntoPars(oracleExecuteMessage.Stmt.Pars, packetData, index, oracleExecuteMessage.Stmt)
	}
	fmt.Println(oracleExecuteMessage.Stmt.Pars)
	for _, par := range oracleExecuteMessage.Stmt.Pars {
		fmt.Println(par.BValue)
	}
	return oracleExecuteMessage, oracleExecuteMessage.Stmt
}

func DecodeReExecuteFunctionMessage(packetData []byte, index int) (models.OracleReExecuteMessage, models.Stmt) {
	oracleReExecuteMessage := models.OracleReExecuteMessage{}
	var cursorId int
	cursorId, index = session.GetInt(2, true, true, packetData, index)
	oracleReExecuteMessage.Stmt = stmtMap[cursorId]
	oracleReExecuteMessage.Count, index = session.GetInt(2, true, true, packetData, index)
	oracleReExecuteMessage.ExeOp, index = session.GetInt(2, true, true, packetData, index)
	oracleReExecuteMessage.ExecFlag, index = session.GetInt(2, true, true, packetData, index)

	if oracleReExecuteMessage.Stmt.ArrayBindCount > 0 {
		oracleReExecuteMessage.Stmt.Pars, _ = GetValuesIntoParsInBulk(oracleReExecuteMessage.Stmt.Pars, packetData, index, oracleReExecuteMessage.Stmt)
	} else {
		oracleReExecuteMessage.Stmt.Pars, _ = GetValuesIntoPars(oracleReExecuteMessage.Stmt.Pars, packetData, index, oracleReExecuteMessage.Stmt)
	}
	fmt.Println(oracleReExecuteMessage.Stmt.Pars)
	for _, par := range oracleReExecuteMessage.Stmt.Pars {
		fmt.Println(par.BValue)
	}
	return oracleReExecuteMessage, oracleReExecuteMessage.Stmt
}

func DecodeReExecuteAndFetchFunctionMessage(packetData []byte, index int) (models.OracleReExecuteAndFetchMessage, models.Stmt) {
	oracleReExecuteAndFetchMessage := models.OracleReExecuteAndFetchMessage{}
	var cursorId int
	cursorId, index = session.GetInt(2, true, true, packetData, index)
	oracleReExecuteAndFetchMessage.Stmt = stmtMap[cursorId]
	oracleReExecuteAndFetchMessage.Count, index = session.GetInt(2, true, true, packetData, index)
	oracleReExecuteAndFetchMessage.ExeOp, index = session.GetInt(2, true, true, packetData, index)
	oracleReExecuteAndFetchMessage.ExecFlag, index = session.GetInt(2, true, true, packetData, index)

	if oracleReExecuteAndFetchMessage.Stmt.ArrayBindCount > 0 {
		oracleReExecuteAndFetchMessage.Stmt.Pars, _ = GetValuesIntoParsInBulk(oracleReExecuteAndFetchMessage.Stmt.Pars, packetData, index, oracleReExecuteAndFetchMessage.Stmt)
	} else {
		oracleReExecuteAndFetchMessage.Stmt.Pars, _ = GetValuesIntoPars(oracleReExecuteAndFetchMessage.Stmt.Pars, packetData, index, oracleReExecuteAndFetchMessage.Stmt)
	}
	fmt.Println(oracleReExecuteAndFetchMessage.Stmt.Pars)
	for _, par := range oracleReExecuteAndFetchMessage.Stmt.Pars {
		fmt.Println(par.BValue)
	}
	return oracleReExecuteAndFetchMessage, oracleReExecuteAndFetchMessage.Stmt
}

func GetValuesIntoParsInBulk(pars []models.ParameterInfo, buffer []byte, index int, stmt models.Stmt) ([]models.ParameterInfo, int) {
	var returnPars []models.ParameterInfo
	for valueIndex := 0; valueIndex < stmt.ArrayBindCount; valueIndex++ {
		returnPars, index = GetValuesIntoPars(pars, buffer, index, stmt)
		pars = returnPars
	}
	return pars, index
}

func GetValuesIntoPars(pars []models.ParameterInfo, buffer []byte, index int, stmt models.Stmt) ([]models.ParameterInfo, int) {
	if len(pars) == 0 {
		return pars, index
	} else {
		var unknown uint8
		unknown, index = session.GetByte(buffer, index)
		fmt.Println(unknown)
		for i, par := range pars {
			if pars[i].Flag == 0x80 {
				continue
			}
			if !stmt.Parse && par.Direction == models.Output && stmt.StmtType != models.PLSQL {
				continue
			}
			if par.DataType == models.REFCURSOR {
				_, index = session.GetBytes(2, buffer, index)
			} else if par.Direction == models.Input &&
				(par.DataType == models.OCIClobLocator || par.DataType == models.OCIBlobLocator || par.DataType == models.OCIFileLocator) {
				_, index = session.GetInt(2, true, true, buffer, index)
				pars[i].BValue, index = session.GetClr(buffer, index)
			} else {
				if par.CusType != nil {
					_, index = session.GetBytes(4, buffer, index)
					_, index = session.GetInt(4, true, true, buffer, index)
					_, index = session.GetBytes(2, buffer, index)
					pars[i].BValue, index = session.GetClr(buffer, index)
				} else {
					fmt.Println(par.MaxNoOfArrayElements)
					if par.MaxNoOfArrayElements > 0 {
						par.BValue = nil
					} else {
						pars[i].BValue, index = session.GetClr(buffer, index)
					}
				}
			}
		}

	}
	return pars, index
}

func ReadParamData(buffer []byte, index int) (models.ParameterInfo, int) {
	var parameterInfo models.ParameterInfo
	var dataType uint8
	dataType, index = session.GetByte(buffer, index)
	parameterInfo.DataType = models.TNSTypeFromInt(int(dataType))
	parameterInfo.Flag, index = session.GetByte(buffer, index)
	parameterInfo.Precision, index = session.GetByte(buffer, index)
	parameterInfo.Scale, index = session.GetByte(buffer, index)
	parameterInfo.MaxLen, index = session.GetInt(4, true, true, buffer, index)
	parameterInfo.MaxNoOfArrayElements, index = session.GetInt(4, true, true, buffer, index)
	if session.TTCVersion >= 10 {
		parameterInfo.ContFlag, index = session.GetInt(8, true, true, buffer, index)
	} else {
		parameterInfo.ContFlag, index = session.GetInt(4, true, true, buffer, index)
	}
	if buffer[index] != 0 {
		_, index = session.GetInt(4, true, true, buffer, index)
		parameterInfo.ToID, index = session.GetClr(buffer, index)
	} else {
		index += 1
	}
	parameterInfo.Version, index = session.GetInt(2, true, true, buffer, index)
	parameterInfo.CharsetID, index = session.GetInt(2, true, true, buffer, index)
	var charsetForm uint8
	charsetForm, index = session.GetByte(buffer, index)
	parameterInfo.CharsetForm = int(charsetForm)
	parameterInfo.MaxCharLen, index = session.GetInt(4, true, true, buffer, index)
	if session.TTCVersion >= 8 {
		parameterInfo.Oaccollid, index = session.GetInt(4, true, true, buffer, index)
	}
	return parameterInfo, index
}

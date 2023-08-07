package oracleparser

import (
	"encoding/binary"
	"reflect"

	"go.keploy.io/server/pkg/models"
)

func DecodeOracleMessageWithData(Packets [][]byte, stmt interface{}) (models.OracleHeader, interface{}, bool, models.DataPacketType, interface{}, error) {
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
		DataMessageType: models.OracleMessageWithDataMessageType,
	}

	var messageList []models.OracleResponseDataMessage
	var message models.OracleResponseDataMessage
	stmtObj, _ := stmt.(models.Stmt)
	oldStmt, _ := stmt.(models.Stmt)
	dataSet := oldStmt.DataSet

	index := 11
	loop := true
	for loop {
		switch models.DataPacketTypeFromInt(int(packetData[index-1])) {
		case models.TNS_MSG_TYPE_ROW_HEADER:
			var messageData models.OracleResponseRowHeaderMessage
			dataSet, index, stmtObj, messageData = DecodeRowHeader(packetData, index, dataSet, stmtObj)
			stmtObj.DataSet = dataSet
			message.MessageType = models.TNS_MSG_TYPE_ROW_HEADER
			message.MessageData = messageData
			messageList = append(messageList, message)
		case models.TNS_MSG_TYPE_ROW_DATA:
			var messageData models.OracleResponseRowMessage
			dataSet, index, stmtObj, messageData = DecodeRowData(packetData, index, dataSet, stmtObj)
			stmtObj.DataSet = dataSet
			message.MessageType = models.TNS_MSG_TYPE_LOB_DATA
			message.MessageData = messageData
			messageList = append(messageList, message)
		case models.TNS_MSG_TYPE_FLUSH_OUT_BINDS:
			var messageData models.OracleResponseFlushOutBindsMesage
			dataSet, index, stmtObj, messageData = DecodeFlushOutBinds(packetData, index, dataSet, stmtObj)
			stmtObj.DataSet = dataSet
			message.MessageType = models.TNS_MSG_TYPE_FLUSH_OUT_BINDS
			message.MessageData = messageData
			messageList = append(messageList, message)
		case models.TNS_MSG_TYPE_DESCRIBE_INFO:
			var messageData models.OracleResponseDescribeInfoMessage
			dataSet, index, stmtObj, messageData = DecodeDescribeInfo(packetData, index, dataSet, stmtObj)
			stmtObj.DataSet = dataSet
			message.MessageType = models.TNS_MSG_TYPE_DESCRIBE_INFO
			message.MessageData = messageData
			messageList = append(messageList, message)
		case models.TNS_MSG_TYPE_BIT_VECTOR:
			var messageData models.OracleResponseBitVectorMessage
			dataSet, index, stmtObj, messageData = DecodeBitVector(packetData, index, dataSet, stmtObj)
			stmtObj.DataSet = dataSet
			message.MessageType = models.TNS_MSG_TYPE_BIT_VECTOR
			message.MessageData = messageData
			messageList = append(messageList, message)
		case models.TNS_MSG_TYPE_IO_VECTOR:
			var messageData models.OracleResponseIOVectorMessage
			dataSet, index, stmtObj, messageData = DecodeIOVector(packetData, index, dataSet, stmtObj)
			stmtObj.DataSet = dataSet
			message.MessageType = models.TNS_MSG_TYPE_IO_VECTOR
			message.MessageData = messageData
			messageList = append(messageList, message)
		case models.OracleMsgTypeError:
			var messageData models.OracleResponseErrorMessage
			index, stmtObj, messageData = DecodeError(packetData, index, stmtObj)
			message.MessageType = models.OracleMsgTypeError
			message.MessageData = messageData
			messageList = append(messageList, message)
			if messageData.Summary.RetCode != 0 {
				if messageData.Summary.RetCode == 1403 {
					stmtObj.HasMoreRows = false
					messageData.Summary = nil
				}
			}
			requestHeader.Session.Summary = messageData.Summary
			session.Summary = requestHeader.Session.Summary
		case models.OracleMsgTypeWarning:
			var messageData models.OracleDataMessageWarning
			messageData, index = DecodeWarningInfoDataMessage(packetData, index)
			message.MessageType = models.OracleMsgTypeWarning
			message.MessageData = messageData
			messageList = append(messageList, message)
		case models.OracleMsgTypeServerSidePiggyback:
			var messageData models.OracleDataMessageServerSidePiggyback
			messageData, index = DecodePiggyBackMessage(packetData, index)
			message.MessageType = models.OracleMsgTypeServerSidePiggyback
			message.MessageData = messageData
			messageList = append(messageList, message)
		case models.OracleMsgTypeStatus:
			var messageData models.OracleDataMessageStatus
			messageData, index = DecodeStatusDataMessage(packetData, index)
			message.MessageType = models.OracleMsgTypeStatus
			message.MessageData = messageData
			messageList = append(messageList, message)
			requestHeader.Session = session
			loop = false
		case models.OracleMsgTypeParameter:
			var messageData models.OracleResponseParameterMessage
			index, stmtObj, messageData = DecodeParameter(packetData, index, stmtObj)
			message.MessageType = models.OracleMsgTypeParameter
			message.MessageData = messageData
			messageList = append(messageList, message)
			requestHeader.Session = session
		}
		index += 1
		if len(packetData) <= index {
			if session.Summary != nil {
				stmtObj.CursorID = session.Summary.CursorID
			}
			loop = false
		}
	}

	if stmtObj.StmtType != models.DEFAULT {
		if oldStmt.StmtType != models.DEFAULT {
			if _, ok := stmtMap[oldStmt.CursorID]; ok {
				delete(stmtMap, oldStmt.CursorID)
				stmtMap[stmtObj.CursorID] = stmtObj
			} else {
				if stmtMap == nil {
					stmtMap = make(map[int]models.Stmt)
				}
				stmtMap[stmtObj.CursorID] = stmtObj
			}
		} else {
			if stmtMap == nil {
				stmtMap = make(map[int]models.Stmt)
			}
			stmtMap[stmtObj.CursorID] = stmtObj
		}
	}

	requestMessage.DataMessage = messageList
	return requestHeader, requestMessage, false, models.DefaultDataPacket, stmtObj, nil
}

func DecodeRowHeader(packetData []byte, index int, dataSet models.DataSet, stmt models.Stmt) (models.DataSet, int, models.Stmt, models.OracleResponseRowHeaderMessage) {
	var columnCount int
	var num int
	var bitVector []byte
	_, index = session.GetByte(packetData, index)
	columnCount, index = session.GetInt(2, true, true, packetData, index)
	num, index = session.GetInt(4, true, true, packetData, index)
	columnCount += num * 0x100
	if columnCount > dataSet.ColumnCount {
		dataSet.ColumnCount = columnCount
	}
	if len(dataSet.CurrentRow) != dataSet.ColumnCount {
		dataSet.CurrentRow = make(models.Row, dataSet.ColumnCount)
	}
	dataSet.RowCount, index = session.GetInt(4, true, true, packetData, index)
	dataSet.UACBufferLength, index = session.GetInt(2, true, true, packetData, index)
	bitVector, index = session.GetDlc(packetData, index)
	dataSet.SetBitVector(bitVector)
	_, index = session.GetDlc(packetData, index)
	messaage := models.OracleResponseRowHeaderMessage{
		ColumnCount:     columnCount,
		Num:             num,
		RowCount:        dataSet.RowCount,
		UacBufferLength: dataSet.UACBufferLength,
		BitVector:       bitVector,
	}
	return dataSet, index, stmt, messaage
}

func DecodeFlushOutBinds(packetData []byte, index int, dataSet models.DataSet, stmt models.Stmt) (models.DataSet, int, models.Stmt, models.OracleResponseFlushOutBindsMesage) {
	return dataSet, index, stmt, models.OracleResponseFlushOutBindsMesage{}
}

func DecodeBitVector(packetData []byte, index int, dataSet models.DataSet, stmt models.Stmt) (models.DataSet, int, models.Stmt, models.OracleResponseBitVectorMessage) {
	var noOfColumnSent int
	noOfColumnSent, index = session.GetInt(2, true, true, packetData, index)
	bitVectorLen := dataSet.ColumnCount / 8
	if dataSet.ColumnCount%8 > 0 {
		bitVectorLen++
	}
	bitVector := make([]byte, bitVectorLen)
	for x := 0; x < bitVectorLen; x++ {
		bitVector[x], index = session.GetByte(packetData, index)
	}
	dataSet.SetBitVector(bitVector)
	message := models.OracleResponseBitVectorMessage{
		NumCloumnsSent: noOfColumnSent,
		BitVector:      bitVector,
	}
	return dataSet, index, stmt, message
}

func DecodeIOVector(packetData []byte, index int, dataSet models.DataSet, stmt models.Stmt) (models.DataSet, int, models.Stmt, models.OracleResponseIOVectorMessage) {
	var columnCount int
	var num int
	var bitVector []byte
	var direction uint8
	var directionArray []models.ParameterDirection
	_, index = session.GetByte(packetData, index)
	columnCount, index = session.GetInt(2, true, true, packetData, index)
	num, index = session.GetInt(4, true, true, packetData, index)
	columnCount += num * 0x100
	if columnCount > dataSet.ColumnCount {
		dataSet.ColumnCount = columnCount
	}
	if len(dataSet.CurrentRow) != dataSet.ColumnCount {
		dataSet.CurrentRow = make(models.Row, dataSet.ColumnCount)
	}
	dataSet.RowCount, index = session.GetInt(4, true, true, packetData, index)
	dataSet.UACBufferLength, index = session.GetInt(2, true, true, packetData, index)
	bitVector, index = session.GetDlc(packetData, index)
	dataSet.SetBitVector(bitVector)
	_, index = session.GetDlc(packetData, index)
	for x := 0; x < dataSet.ColumnCount; x++ {
		direction, index = session.GetByte(packetData, index)
		switch direction {
		case 32:
			directionArray = append(directionArray, models.Input)
			stmt.Pars[x].Direction = models.Input
		case 16:
			directionArray = append(directionArray, models.Output)
			stmt.Pars[x].Direction = models.Output
			stmt.ContainOutputPars = true
		case 48:
			directionArray = append(directionArray, models.InOut)
			stmt.Pars[x].Direction = models.InOut
			stmt.ContainOutputPars = true
		}
	}
	message := models.OracleResponseIOVectorMessage{
		ColumnCount:            columnCount,
		Num:                    num,
		RowCount:               dataSet.RowCount,
		UacBufferLength:        dataSet.UACBufferLength,
		BitVector:              bitVector,
		ParameterDirectorArray: directionArray,
	}
	return dataSet, index, stmt, message
}

func DecodeDescribeInfo(packetData []byte, index int, dataSet models.DataSet, stmt models.Stmt) (models.DataSet, int, models.Stmt, models.OracleResponseDescribeInfoMessage) {
	var size uint8
	size, index = session.GetByte(packetData, index)
	_, index = session.GetBytes(int(size), packetData, index)
	dataSet.MaxRowSize, index = session.GetInt(4, true, true, packetData, index)
	dataSet.ColumnCount, index = session.GetInt(4, true, true, packetData, index)
	if dataSet.ColumnCount > 0 {
		_, index = session.GetByte(packetData, index)
	}
	dataSet.Cols = make([]models.ParameterInfo, dataSet.ColumnCount)
	var paramInfoList []models.ParameterInfo
	var paramInfo models.ParameterInfo
	for x := 0; x < dataSet.ColumnCount; x++ {
		paramInfo, index = dataSet.Cols[x].Load(packetData, index, session)
		paramInfoList = append(paramInfoList, paramInfo)
		if dataSet.Cols[x].DataType == models.LONG || dataSet.Cols[x].DataType == models.LongRaw {
			stmt.HasLONG = true
		}
		if dataSet.Cols[x].DataType == models.OCIClobLocator || dataSet.Cols[x].DataType == models.OCIBlobLocator {
			stmt.HasBLOB = true
		}
	}
	stmt.Columns = make([]models.ParameterInfo, dataSet.ColumnCount)
	copy(stmt.Columns, dataSet.Cols)
	_, index = session.GetDlc(packetData, index)
	if session.TTCVersion >= 3 {
		_, index = session.GetInt(4, true, true, packetData, index)
		_, index = session.GetInt(4, true, true, packetData, index)
	}
	if session.TTCVersion >= 4 {
		_, index = session.GetInt(4, true, true, packetData, index)
		_, index = session.GetInt(4, true, true, packetData, index)
	}
	if session.TTCVersion >= 5 {
		_, index = session.GetDlc(packetData, index)
	}
	message := models.OracleResponseDescribeInfoMessage{
		ColumnCount:   dataSet.ColumnCount,
		MaxRowSize:    dataSet.RowCount,
		Size:          size,
		ParamInfoList: paramInfoList,
	}
	return dataSet, index, stmt, message
}

func DecodeRowData(packetData []byte, index int, dataSet models.DataSet, stmt models.Stmt) (models.DataSet, int, models.Stmt, models.OracleResponseRowMessage) {
	var numList []int
	var num int
	var oraclePrimeValue models.OraclePrimeValue
	var oracleColumnValue models.CalculateColumnValue
	var CalculateParameterValueList []models.OraclePrimeValue
	var CalculateColumnValueList []models.CalculateColumnValue
	var CursorArray []models.RefCursor
	if stmt.HasReturnClause && stmt.ContainOutputPars {
		for x := 0; x < len(stmt.Pars); x++ {
			if stmt.Pars[x].Direction == models.Output {
				num, index = session.GetInt(4, true, true, packetData, index)
				numList = append(numList, num)
				if num == 0 {
					stmt.Pars[x].BValue = nil
					stmt.Pars[x].Value = nil
				} else {
					oraclePrimeValue, index = stmt.CalculateParameterValue(&stmt.Pars[x], session, packetData, index)
					CalculateParameterValueList = append(CalculateParameterValueList, oraclePrimeValue)
					_, index = session.GetInt(2, true, true, packetData, index)
				}
			}
		}
	} else {
		var cursor *models.RefCursor
		if stmt.ContainOutputPars {
			for x := 0; x < len(stmt.Pars); x++ {
				if stmt.Pars[x].DataType == models.REFCURSOR {
					typ := reflect.TypeOf(stmt.Pars[x].Value)
					if typ.Kind() == reflect.Ptr {
						cursor, _ = stmt.Pars[x].Value.(*models.RefCursor)

						cursor.Parent = &stmt
						cursor.AutoClose = true
						cursor, index = cursor.Load(session, packetData, index)
						CursorArray = append(CursorArray, *cursor)
						if stmt.StmtType == models.PLSQL {
							_, index = session.GetInt(2, true, true, packetData, index)
						}
					}
				} else {
					if stmt.Pars[x].Direction != models.Input {
						oraclePrimeValue, index = stmt.CalculateParameterValue(&stmt.Pars[x], session, packetData, index)
						CalculateParameterValueList = append(CalculateParameterValueList, oraclePrimeValue)
						_, index = session.GetInt(2, true, true, packetData, index)
					}
				}
			}
		} else {
			// see if it is re-executed
			if len(dataSet.Cols) == 0 && len(stmt.Columns) > 0 {
				dataSet.Cols = make([]models.ParameterInfo, len(stmt.Columns))
				copy(dataSet.Cols, stmt.Columns)
			}
			for x := 0; x < len(dataSet.Cols); x++ {
				if dataSet.Cols[x].GetDataFromServer {
					oracleColumnValue, index = stmt.CalculateColumnValue(&dataSet.Cols[x], false, session, packetData, index)
					CalculateColumnValueList = append(CalculateColumnValueList, oracleColumnValue)
					if dataSet.Cols[x].DataType == models.LONG || dataSet.Cols[x].DataType == models.LongRaw {
						_, index = session.GetInt(4, true, true, packetData, index)
						_, index = session.GetInt(4, true, true, packetData, index)
					}
				}
			}
			newRow := make(models.Row, dataSet.ColumnCount)
			for x := 0; x < len(dataSet.Cols); x++ {
				newRow[x] = dataSet.Cols[x].OPrimValue
			}
			//copy(newRow, dataSet.currentRow)
			dataSet.Rows = append(dataSet.Rows, newRow)
		}
	}
	return dataSet, index, stmt, models.OracleResponseRowMessage{
		NumList:                     numList,
		CalculateParameterValueList: CalculateParameterValueList,
		CalculateColumnValueList:    CalculateColumnValueList,
		CursorArray:                 CursorArray,
	}
}

func DecodeError(packetData []byte, index int, stmt models.Stmt) (int, models.Stmt, models.OracleResponseErrorMessage) {
	summary := new(models.SummaryObject)
	summary, index = NewSummary(packetData, index)
	stmt.CursorID = summary.CursorID
	stmt.DisableCompression = summary.Flags&0x20 != 0
	oracleErrorMsg := models.OracleResponseErrorMessage{
		Summary: summary,
	}
	return index, stmt, oracleErrorMsg
}

func DecodeParameter(packetData []byte, index int, stmt models.Stmt) (int, models.Stmt, models.OracleResponseParameterMessage) {
	var size1 int
	var size2 int
	var size3 int
	var key []byte
	var val []byte
	var num int
	var keyList [][]byte
	var valList [][]byte
	var numList []int
	var bty []byte
	var length int
	size1, index = session.GetInt(2, true, true, packetData, index)
	var ScnForSnapshotList []int
	for x := 0; x < 2; x++ {
		stmt.ScnForSnapshot[x], index = session.GetInt(4, true, true, packetData, index)
		ScnForSnapshotList = append(ScnForSnapshotList, stmt.ScnForSnapshot[x])
	}
	for x := 2; x < size1; x++ {
		_, index = session.GetInt(4, true, true, packetData, index)
	}
	_, index = session.GetInt(2, true, true, packetData, index)
	size2, index = session.GetInt(2, true, true, packetData, index)
	for x := 0; x < size2; x++ {
		key, val, num, index = session.GetKeyVal(packetData, index)
		keyList = append(keyList, key)
		valList = append(valList, val)
		numList = append(numList, num)
		if num == 163 {
			session.TimeZone = val
		}
	}
	if session.TTCVersion >= 4 {
		// get queryID
		size3, index = session.GetInt(4, true, true, packetData, index)
		if size3 > 0 {
			bty, index = session.GetBytes(size3, packetData, index)
			if len(bty) >= 8 {
				stmt.QueryID = binary.LittleEndian.Uint64(bty[size3-8:])
			}
		}
	}
	if session.TTCVersion >= 7 && stmt.StmtType == models.DML && stmt.ArrayBindCount > 0 {
		length, index = session.GetInt(4, true, true, packetData, index)
		for i := 0; i < length; i++ {
			_, index = session.GetInt(8, true, true, packetData, index)
		}
	}
	oracleRepsonseParameterMessage := models.OracleResponseParameterMessage{
		Size1:              size1,
		Size2:              size2,
		Size3:              size3,
		KeyList:            keyList,
		ValList:            valList,
		NumList:            numList,
		Bty:                bty,
		Length:             length,
		ScnForSnapshotList: ScnForSnapshotList,
	}
	return index, stmt, oracleRepsonseParameterMessage
}

func DecodeOracleGetDBVersion(Packets [][]byte) (models.OracleHeader, interface{}, bool, models.DataPacketType, interface{}, error) {
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
		DataMessageType: models.OracleGetDBVersionDataMessageType,
	}
	index := 11
	var oracleDBVersionMessage models.OracleDBVersionMessage
	oracleDBVersionMessage.Length, index = session.GetInt(2, true, true, packetData, index)
	oracleDBVersionMessage.Info, index = session.GetString(int(oracleDBVersionMessage.Length), packetData, index)
	oracleDBVersionMessage.Number, _ = session.GetInt(4, true, true, packetData, index)
	requestMessage.DataMessage = oracleDBVersionMessage
	return requestHeader, requestMessage, false, models.DefaultDataPacket, nil, nil
}

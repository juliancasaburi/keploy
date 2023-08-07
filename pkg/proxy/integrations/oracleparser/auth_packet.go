package oracleparser

import (
	"encoding/binary"
	"fmt"

	"go.keploy.io/server/pkg/models"
)

func DecodeAuthMessage(packetData []byte) models.OracleAuthDataMessageRequest {

	oracleAuthDataMessagePhaseOne := models.OracleAuthDataMessageRequest{}
	oracleAuthDataMessagePhaseOne.HasUser = packetData[13]
	index := 14
	var userByteLength int
	userByteLength, index = session.GetInt(4, true, true, packetData, index)
	oracleAuthDataMessagePhaseOne.UserByteLength = uint32(userByteLength)
	var authMode int
	authMode, index = session.GetInt(4, true, true, packetData, index)
	oracleAuthDataMessagePhaseOne.AuthMode = models.AuthModeFromInt(authMode)
	index += 1
	var numPair int
	numPair, index = session.GetInt(4, true, true, packetData, index)
	oracleAuthDataMessagePhaseOne.NumPair = uint32(numPair)
	index += 2
	if oracleAuthDataMessagePhaseOne.HasUser == 1 {
		userByteLength := packetData[index]
		index += 1
		bytes := packetData[index : index+int(userByteLength)]
		index = index + int(userByteLength)
		oracleAuthDataMessagePhaseOne.UserBytes = string(bytes)
	}
	var authOracleKeyValueList []models.OracleKeyValue
	for i := 0; i < int(numPair); i++ {
		keyBytes, valueBytes, _, retIndex := session.GetKeyVal(packetData, index)
		index = retIndex
		authOracleKeyValueList = append(authOracleKeyValueList, models.OracleKeyValue{
			Key:   string(keyBytes),
			Value: string(valueBytes),
		})
	}
	oracleAuthDataMessagePhaseOne.AuthKeyValue = authOracleKeyValueList
	return oracleAuthDataMessagePhaseOne
}

func DecodeOracleAuthPhaseOneResponse(Packets [][]byte) (models.OracleHeader, interface{}, bool, models.DataPacketType, interface{}, error) {
	var packetData []byte
	var dataPacketType = models.DefaultDataPacket
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
	index := 10
	responseMessages := []models.OracleResponseDataMessage{}
	responseMessage := models.OracleDataMessage{
		DataMessageType: models.OracleAuthPhaseOneDataMessageType,
	}
	loop := true
	for loop {
		var messageCode uint8
		messageCode, index = session.GetByte(packetData, index)
		switch models.DataPacketTypeFromInt(int(messageCode)) {
		case models.OracleMsgTypeStatus:
			fmt.Println("TNS_MSG_TYPE_STATUS")
			var statusMessage models.OracleDataMessageStatus
			statusMessage, index = DecodeStatusDataMessage(packetData, index)
			responseMessages = append(responseMessages, models.OracleResponseDataMessage{
				MessageType: models.OracleMsgTypeStatus,
				MessageData: statusMessage,
			})
		case models.OracleMsgTypeWarning:
			fmt.Println("TNS_MSG_TYPE_WARNING")
			var warningInfo models.OracleDataMessageWarning
			warningInfo, index = DecodeWarningInfoDataMessage(packetData, index)
			responseMessages = append(responseMessages, models.OracleResponseDataMessage{
				MessageType: models.OracleMsgTypeWarning,
				MessageData: warningInfo,
			})

		case models.OracleMsgTypeServerSidePiggyback:
			fmt.Println("TNS_MSG_TYPE_SERVER_SIDE_PIGGYBACK")
			var serverInfo models.OracleDataMessageServerSidePiggyback
			serverInfo, index = DecodePiggyBackMessage(packetData, index)
			responseMessages = append(responseMessages, models.OracleResponseDataMessage{
				MessageType: models.OracleMsgTypeServerSidePiggyback,
				MessageData: serverInfo,
			})

		case models.OracleMsgTypeParameter:
			var dictLen int
			dictLen, index = session.GetInt(4, true, true, packetData, index)
			parameterMessage := models.OracleAuthDataMessageParameter{
				NumParams:    uint32(dictLen),
				AuthKeyValue: []models.OracleKeyValue{},
			}
			for x := 0; x < dictLen; x++ {
				keyBytes, valueBytes, code, retIndex := session.GetKeyVal(packetData, index)
				index = retIndex
				parameterMessage.AuthKeyValue = append(parameterMessage.AuthKeyValue, models.OracleKeyValue{
					Key:   string(keyBytes),
					Value: string(valueBytes),
					Code:  code,
				})
			}
			responseMessages = append(responseMessages, models.OracleResponseDataMessage{
				MessageType: models.OracleMsgTypeParameter,
				MessageData: parameterMessage,
			})
		case models.OracleMsgTypeError:
			requestHeader.Session.Summary, index = NewSummary(packetData, index)
			responseMessages = append(responseMessages, models.OracleResponseDataMessage{
				MessageType: models.OracleMsgTypeError,
				MessageData: *requestHeader.Session.Summary,
			})
			session.Summary = requestHeader.Session.Summary
			loop = false
		default:

		}
	}
	responseMessage.DataMessage = responseMessages
	return requestHeader, responseMessage, false, dataPacketType, nil, nil
}

func DecodeOracleAuthPhaseTwoResponse(Packets [][]byte) (models.OracleHeader, interface{}, bool, models.DataPacketType, interface{}, error) {
	var packetData []byte
	var dataPacketType = models.DefaultDataPacket
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
	responseMessages := []models.OracleResponseDataMessage{}
	responseMessage := models.OracleDataMessage{
		DataMessageType: models.OracleAuthPhaseTwoDataMessageType,
	}
	index := 10
	stop := false
	for !stop {
		msgCode := packetData[index]
		index += 1
		switch models.DataPacketTypeFromInt(int(msgCode)) {
		case models.OracleMsgTypeError:
			fmt.Println("TNS_MSG_TYPE_ERROR")
			requestHeader.Session.Summary, index = NewSummary(packetData, index)
			responseMessages = append(responseMessages, models.OracleResponseDataMessage{
				MessageType: models.OracleMsgTypeError,
				MessageData: *requestHeader.Session.Summary,
			})
			session.Summary = requestHeader.Session.Summary
			stop = true
		case models.OracleMsgTypeParameter:
			fmt.Println("TNS_MSG_TYPE_PARAMETER")
			var dictLen int
			dictLen, index = session.GetInt(2, true, true, packetData, index)
			sessionProperties := []models.OracleKeyValue{}
			for x := 0; x < dictLen; x++ {
				var (
					keyBytes, valueBytes []byte
					code                 int
				)
				keyBytes, valueBytes, code, index = session.GetKeyVal(packetData, index)
				sessionProperties = append(sessionProperties, models.OracleKeyValue{
					Key:   string(keyBytes),
					Value: string(valueBytes),
					Code:  code,
				})
			}
			responseMessages = append(responseMessages, models.OracleResponseDataMessage{
				MessageType: models.OracleMsgTypeParameter,
				MessageData: models.OracleAuthDataMessageParameter{
					// MessageCode:    int(msgCode),
					NumParams:    uint32(dictLen),
					AuthKeyValue: sessionProperties,
				}})
		case models.OracleMsgTypeWarning:
			fmt.Println("TNS_MSG_TYPE_WARNING")
			var warningInfo models.OracleDataMessageWarning
			warningInfo, index = DecodeWarningInfoDataMessage(packetData, index)
			responseMessages = append(responseMessages, models.OracleResponseDataMessage{
				MessageType: models.OracleMsgTypeWarning,
				MessageData: warningInfo,
			})
		case models.OracleMsgTypeServerSidePiggyback:
			fmt.Println("TNS_MSG_TYPE_SERVER_SIDE_PIGGYBACK")
			var serverInfo models.OracleDataMessageServerSidePiggyback
			serverInfo, index = DecodePiggyBackMessage(packetData, index)
			responseMessages = append(responseMessages, models.OracleResponseDataMessage{
				MessageType: models.OracleMsgTypeServerSidePiggyback,
				MessageData: serverInfo,
			})
		default:
			return requestHeader, nil, false, dataPacketType, nil, fmt.Errorf("message code error: received code %d", msgCode)
		}
	}
	responseMessage.DataMessage = responseMessages
	return requestHeader, responseMessage, false, dataPacketType, nil, nil
}

func NewSummary(packetData []byte, index int) (*models.SummaryObject, int) {
	result := new(models.SummaryObject)
	// var err error
	if session.HasEOSCapability {
		result.EndOfCallStatus, index = session.GetInt(4, true, true, packetData, index)
	}
	if session.TTCVersion >= 3 {
		if session.HasFSAPCapability {
			result.EndToEndECIDSequence, index = session.GetInt(2, true, true, packetData, index)
		}
	}
	result.CurRowNumber, index = session.GetInt(4, true, true, packetData, index)
	result.RetCode, index = session.GetInt(2, true, true, packetData, index)
	result.ArrayElmWError, index = session.GetInt(2, true, true, packetData, index)
	result.ArrayElmErrno, index = session.GetInt(2, true, true, packetData, index)
	// if err != nil {
	// 	return nil, err
	// }
	result.CursorID, index = session.GetInt(2, true, true, packetData, index)
	// if err != nil {
	// 	return nil, err
	// }
	result.ErrorPos, index = session.GetInt(2, true, true, packetData, index)
	// if err != nil {
	// 	return nil, err
	// }
	result.SqlType, index = session.GetByte(packetData, index)
	// if err != nil {
	// 	return nil, err
	// }
	result.OerFatal, index = session.GetByte(packetData, index)
	// if err != nil {
	// 	return nil, err
	// }
	if session.TTCVersion >= 6 {
		result.Flags, index = session.GetInt(2, true, true, packetData, index)
		// if err != nil {
		// 	return nil, err
		// }
		result.UserCursorOPT, index = session.GetInt(2, true, true, packetData, index)
		// if err != nil {
		// 	return nil, err
		// }
	} else {
		var temp uint8
		temp, index = session.GetByte(packetData, index)
		// if err != nil {
		// 	return nil, err
		// }
		result.Flags = int(temp)
		temp, index = session.GetByte(packetData, index)
		// if err != nil {
		// 	return nil, err
		// }
		result.UserCursorOPT = int(temp)
	}
	result.UpiParam, index = session.GetByte(packetData, index)
	// if err != nil {
	// 	return nil, err
	// }
	result.WarningFlag, index = session.GetByte(packetData, index)
	// if err != nil {
	// 	return nil, err
	// }
	result.Rba, index = session.GetInt(4, true, true, packetData, index)
	// if err != nil {
	// 	return nil, err
	// }
	result.PartitionID, index = session.GetInt(2, true, true, packetData, index)
	// if err != nil {
	// 	return nil, err
	// }
	result.TableID, index = session.GetByte(packetData, index)
	// if err != nil {
	// 	return nil, err
	// }
	result.BlockNumber, index = session.GetInt(4, true, true, packetData, index)
	// if err != nil {
	// 	return nil, err
	// }
	result.SlotNumber, index = session.GetInt(2, true, true, packetData, index)
	// if err != nil {
	// 	return nil, err
	// }
	result.OsError, index = session.GetInt(4, true, true, packetData, index)
	// if err != nil {
	// 	return nil, err
	// }
	result.StmtNumber, index = session.GetByte(packetData, index)
	// if err != nil {
	// 	return nil, err
	// }
	result.CallNumber, index = session.GetByte(packetData, index)
	// if err != nil {
	// 	return nil, err
	// }
	result.Pad1, index = session.GetInt(2, true, true, packetData, index)
	// if err != nil {
	// 	return nil, err
	// }
	result.SuccessIter, index = session.GetInt(4, true, true, packetData, index)
	// if err != nil {
	// 	return nil, err
	// }
	_, index = session.GetDlc(packetData, index)
	if session.TTCVersion < 7 {
		_, index = session.GetDlc(packetData, index)
		_, index = session.GetDlc(packetData, index)
		_, index = session.GetDlc(packetData, index)
	} else {
		var length int
		length, index = session.GetInt(2, true, true, packetData, index)
		// if err != nil {
		// 	return nil, err
		// }
		if length > 0 {
			result.BindErrors = make([]models.BindError, length)
			var num byte
			num, index = session.GetByte(packetData, index)
			// if err != nil {
			// 	return nil, err
			// }
			flag := num == 0xFE
			for x := 0; x < length; x++ {
				if flag {
					if session.UseBigClrChunks {
						_, index = session.GetInt(4, true, true, packetData, index)
					} else {
						_, index = session.GetByte(packetData, index)
					}
				}
				result.BindErrors[x].ErrorCode, index = session.GetInt(2, true, true, packetData, index)
				// if err != nil {
				// 	return nil, err
				// }
			}
			if flag {
				_, index = session.GetByte(packetData, index)
			}
		}
		length, index = session.GetInt(4, true, true, packetData, index)
		// if err != nil {
		// 	return nil, err
		// }
		if length > 0 {
			if len(result.BindErrors) == 0 {
				result.BindErrors = make([]models.BindError, length)
			}
			var num byte
			num, index = session.GetByte(packetData, index)
			// if err != nil {
			// 	return nil, err
			// }
			flag := num == 0xFE
			for x := 0; x < length; x++ {
				if flag {
					if session.UseBigClrChunks {
						_, index = session.GetInt(4, true, true, packetData, index)
					} else {
						_, index = session.GetByte(packetData, index)
					}
				}
				result.BindErrors[x].RowOffset, index = session.GetInt(4, true, true, packetData, index)
				// if err != nil {
				// 	return nil, err
				// }
			}
			if flag {
				_, index = session.GetByte(packetData, index)
			}
		}
		length, index = session.GetInt(2, true, true, packetData, index)
		// if err != nil {
		// 	return nil, err
		// }
		if length > 0 {
			if len(result.BindErrors) == 0 {
				result.BindErrors = make([]models.BindError, length)
			}
			_, index = session.GetByte(packetData, index)
			for x := 0; x < length; x++ {
				_, index = session.GetInt(2, true, true, packetData, index)
				// if err != nil {
				// 	return nil, err
				// }
				result.BindErrors[x].ErrorMsg, index = session.GetClr(packetData, index)
				// if err != nil {
				// 	return nil, err
				// }
				_, index = session.GetByte(packetData, index)
				_, index = session.GetByte(packetData, index)
			}
		}
		if session.TTCVersion >= 7 {
			result.RetCode, index = session.GetInt(4, true, true, packetData, index)
			// if err != nil {
			// 	return nil, err
			// }
			result.CurRowNumber, index = session.GetInt(8, true, true, packetData, index)
			// if err != nil {
			// 	return nil, err
			// }
		}
	}
	if result.RetCode != 0 {
		var errMsg []byte
		errMsg, index = session.GetClr(packetData, index)
		result.ErrorMessage = string(errMsg)
	}
	if len(result.BindErrors) > 0 && result.RetCode == 24381 {
		result.RetCode = result.BindErrors[0].ErrorCode
		result.ErrorMessage = string(result.BindErrors[0].ErrorMsg)
	}
	return result, index
}

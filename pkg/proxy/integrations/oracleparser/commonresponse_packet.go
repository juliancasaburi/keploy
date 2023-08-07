package oracleparser

import (
	"fmt"

	"go.keploy.io/server/pkg/models"
)

func DecodeStatusDataMessage(packetData []byte, index int) (models.OracleDataMessageStatus, int) {
	statusPacket := models.OracleDataMessageStatus{}
	if session.HasEOSCapability {
		statusPacket.EndOfCallStatus, index = session.GetInt(4, true, true, packetData, index)
		if session.Summary != nil {
			session.Summary.EndOfCallStatus = statusPacket.EndOfCallStatus
		}
	}
	if session.HasFSAPCapability {
		if session.Summary == nil {
			session.Summary = new(models.SummaryObject)
		}
		statusPacket.EndToEndECIDSequence, index = session.GetInt(2, true, true, packetData, index)
		session.Summary.EndToEndECIDSequence = statusPacket.EndToEndECIDSequence
	}
	return statusPacket, index
}

func DecodeWarningInfoDataMessage(packetData []byte, index int) (models.OracleDataMessageWarning, int) {
	var result models.OracleDataMessageWarning
	result.ErrCode, index = session.GetInt(2, true, true, packetData, index)

	result.ErrLength, index = session.GetInt(2, true, true, packetData, index)

	result.Flag, index = session.GetInt(2, true, true, packetData, index)

	if result.ErrCode != 0 && result.ErrLength > 0 {
		var msgBytes []byte
		msgBytes, index = session.GetClr(packetData, index)

		result.ErrMsg = string(msgBytes)
	}
	return result, index
}

func DecodePiggyBackMessage(packetData []byte, index int) (models.OracleDataMessageServerSidePiggyback, int) {
	var result models.OracleDataMessageServerSidePiggyback
	var code byte
	code, index = session.GetByte(packetData, index)
	fmt.Println("the opcode in the piggybackServerMessage: ", code)
	switch models.PiggyBackTypeFromInt(int(code)) {
	case models.TNS_SERVER_PIGGYBACK_EXT_SYNC:
		extSyncPacket := models.OracleServerPiggybackExtSync{}
		extSyncPacket.DTYCount, index = session.GetInt(2, true, true, packetData, index)

		extSyncPacket.DTYLength, index = session.GetByte(packetData, index)
		result.OpCode = models.TNS_SERVER_PIGGYBACK_EXT_SYNC
		result.ServerSideInfo = extSyncPacket
	case models.TNS_SERVER_PIGGYBACK_OS_PID_MTS:
		osPidPacket := models.OracleServerPiggybackOsPid{}
		osPidPacket.Length, index = session.GetInt(2, true, true, packetData, index)

		osPidPacket.Byte, index = session.GetByte(packetData, index)

		osPidPacket.Pid, index = session.GetBytes(osPidPacket.Length, packetData, index)
		result.OpCode = models.TNS_SERVER_PIGGYBACK_OS_PID_MTS
		result.ServerSideInfo = osPidPacket
	case models.TNS_SERVER_PIGGYBACK_AC_REPLAY_CONTEXT:
		replayCtxPacket := models.OracleServerPiggybackReplayCtx{}
		replayCtxPacket.DTYCount, index = session.GetInt(2, true, true, packetData, index)

		replayCtxPacket.DTYLength, index = session.GetByte(packetData, index)

		replayCtxPacket.Flags, index = session.GetInt(4, true, true, packetData, index)

		replayCtxPacket.ErrorCode, index = session.GetInt(4, true, true, packetData, index)

		replayCtxPacket.Queue, index = session.GetByte(packetData, index)

		replayCtxPacket.Length, index = session.GetInt(4, true, true, packetData, index)

		replayCtxPacket.ReplayContext, index = session.GetClr(packetData, index)
		result.OpCode = models.TNS_SERVER_PIGGYBACK_AC_REPLAY_CONTEXT
		result.ServerSideInfo = replayCtxPacket
	case models.TNS_SERVER_PIGGYBACK_LTXID:
		ltxIdPacket := models.OracleServerPiggybackLtxId{}
		ltxIdPacket.Length, index = session.GetInt(4, true, true, packetData, index)

		ltxIdPacket.TransactionID, index = session.GetClr(packetData, index)
		result.ServerSideInfo = ltxIdPacket
		result.OpCode = models.TNS_SERVER_PIGGYBACK_LTXID
	case models.TNS_SERVER_PIGGYBACK_SESS_RET:
		sessionReturnPacket := models.OracleServerPiggybackSessionReturn{}
		sessionReturnPacket.SkipInt, index = session.GetInt(2, true, true, packetData, index)

		sessionReturnPacket.SkipByte, index = session.GetByte(packetData, index)

		sessionReturnPacket.Length, index = session.GetInt(2, true, true, packetData, index)
		// get nls data
		for i := 0; i < sessionReturnPacket.Length; i++ {
			var (
				nlsKey  []byte
				nlsVal  []byte
				nlsCode int
			)
			nlsKey, nlsVal, nlsCode, index = session.GetKeyVal(packetData, index)

			sessionReturnPacket.SessionProps = append(sessionReturnPacket.SessionProps, models.OracleKeyValue{
				Key:   string(nlsKey),
				Value: string(nlsVal),
				Code:  nlsCode,
			})
			// conn.NLSData.SaveNLSValue(string(nlsKey), string(nlsVal), nlsCode)
		}
		sessionReturnPacket.Flag, index = session.GetInt(4, true, true, packetData, index)

		sessionReturnPacket.SessionId, index = session.GetInt(4, true, true, packetData, index)

		sessionReturnPacket.SerialId, index = session.GetInt(2, true, true, packetData, index)
		result.OpCode = models.TNS_SERVER_PIGGYBACK_SESS_RET
		result.ServerSideInfo = sessionReturnPacket
	case models.TNS_SERVER_PIGGYBACK_SYNC:
		syncPacket := models.OracleServerPiggybackSync{}
		syncPacket.DTYCount, index = session.GetInt(2, true, true, packetData, index)

		syncPacket.DTYLength, index = session.GetByte(packetData, index)

		syncPacket.NUM_PAIRS, index = session.GetInt(4, true, true, packetData, index)

		syncPacket.Length, index = session.GetByte(packetData, index)

		for i := 0; i < syncPacket.NUM_PAIRS; i++ {
			var (
				nlsKey  []byte
				nlsVal  []byte
				nlsCode int
			)
			nlsKey, nlsVal, nlsCode, index = session.GetKeyVal(packetData, index)

			syncPacket.SessionProps = append(syncPacket.SessionProps, models.OracleKeyValue{
				Key:   string(nlsKey),
				Value: string(nlsVal),
				Code:  nlsCode,
			})
		}
		syncPacket.Flags, index = session.GetInt(4, true, true, packetData, index)
		result.ServerSideInfo = syncPacket
		result.OpCode = models.TNS_SERVER_PIGGYBACK_SYNC

	}
	return result, index
}

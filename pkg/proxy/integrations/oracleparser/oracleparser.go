package oracleparser

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"time"

	"go.keploy.io/server/pkg/hooks"
	"go.keploy.io/server/pkg/models"
	"go.keploy.io/server/pkg/proxy/util"
	"go.uber.org/zap"
)

var Emoji = "\U0001F430" + " Keploy:"
var TNS_MAX_CONNECT_DATA uint16 = 230
var session models.PacketSession

type RequestResponse struct {
	Requests  [][]byte
	Responses [][]byte
}

// Determines whether the outgoing packets belong to oracle
func IsOutgoingOracle(buffer []byte) bool {
	messageLength := binary.BigEndian.Uint16(buffer[0:2])
	if int(messageLength) == len(buffer) {
		return true
	} else if int(messageLength) < len(buffer) {
		var sum = messageLength
		sum += binary.BigEndian.Uint16(buffer[messageLength : messageLength+2])
		return int(sum) == len(buffer)
	}
	return false
}

// Processes the Oracle packets
func ProcessOraclePackets(reuqestBuffer []byte, clientConn, destConn net.Conn, h *hooks.Hook, logger *zap.Logger, FilterPid *bool, port uint32, kernalPid uint32) {
	switch models.GetMode() {
	case models.MODE_RECORD:
		err := ReadWriteProtocol(reuqestBuffer, clientConn, destConn, logger, port, FilterPid, kernalPid, h)
		if err != nil {
			logger.Error(Emoji+"failed to call next", zap.Error(err))
			clientConn.Close()
			return
		}
	case models.MODE_TEST:
		fmt.Println("not yet ready")
	default:
		fmt.Println("not yet ready")
		clientConn.Close()
	}

}

// This function facilitates bidirectional communication by reading packets from both the client and the server,
// and subsequently writing the appropriate responses back to both entities. Meanwhile this fuction also creates
// mocks in record mode and use mocks in test mode.
func ReadWriteProtocol(firstbuffer []byte, clientConn, destConn net.Conn, logger *zap.Logger, port uint32, FilterPid *bool, kernalPid uint32, h *hooks.Hook) error {

	fmt.Println(Emoji, "trying to forward requests to target: ", destConn.RemoteAddr().String())
	defer destConn.Close()

	var nextRequest [][]byte
	var nextResponse [][]byte
	dataPacketType := models.DefaultDataPacket
	nextDataPacketType := models.DefaultDataPacket
	var stmt interface{}
	for {
		fmt.Println("inside connection request")
		var (
			oracleRequests  = []models.OracleRequest{}
			oracleResponses = []models.OracleResponse{}
		)
		for {
			started := time.Now()
			var (
				requestHeader     models.OracleHeader
				requestMessage    interface{}
				continueStatement bool
				err               error
				buffer            [][]byte
				clientBuffer      []byte
			)
			for {
				if len(firstbuffer) > 0 {
					clientBuffer = firstbuffer
					firstbuffer = []byte{}
				} else {
					if len(nextRequest) == 0 {
						clientBuffer, err = util.ReadBytes(clientConn)
						if err != nil {
							logger.Error(Emoji+"failed to read the request message in proxy", zap.Error(err), zap.Any("proxy port", port))
							return err
						}
					} else {
						clientBuffer = nextRequest[0]
						nextRequest = nextRequest[:0]
					}
				}
				buffer = append(buffer, clientBuffer)
				breakStatement, packetNumber, index := recievedAllPackets(buffer)
				if breakStatement {
					if packetNumber != 0 || index != 0 {
						buffer, nextRequest, err = cut2DSlice(buffer, packetNumber, index)
						if err != nil {
							return err
						}
					}
					break
				}
			}
			fmt.Println("Request")
			fmt.Println(buffer)
			requestHeader, requestMessage, continueStatement, nextDataPacketType, stmt, err = Decode(buffer, dataPacketType, true, nil)
			if err != nil {
				logger.Error(Emoji+"failed to read the request message in proxy", zap.Error(err), zap.Any("proxy port", port))
				return err
			}
			readRequestDelay := time.Since(started)
			for _, packet := range buffer {
				_, err = destConn.Write(packet)
				if err != nil {
					logger.Error(Emoji+"failed to write request message to the destination server", zap.Error(err))
					return err
				}
			}
			oracleRequests = append(oracleRequests, models.OracleRequest{
				Header:    requestHeader,
				Message:   requestMessage,
				ReadDelay: int64(readRequestDelay),
			})
			dataPacketType = nextDataPacketType
			if continueStatement {
				continue
			}
			fmt.Println(oracleRequests)
			break
		}

		for {
			started := time.Now()
			var (
				responseHeader    models.OracleHeader
				responseMessage   interface{}
				continueStatement bool
				err               error
				buffer            [][]byte
				serverBuffer      []byte
			)
			for {
				if len(nextResponse) == 0 {
					serverBuffer, err = util.ReadBytes(destConn)
					if err != nil {
						logger.Error(Emoji+"failed to read the response message in proxy", zap.Error(err), zap.Any("proxy port", port))
						return err
					}
				} else {
					serverBuffer = nextResponse[0]
					nextResponse = nextResponse[:0]
				}
				buffer = append(buffer, serverBuffer)
				breakStatement, packetNumber, index := recievedAllPackets(buffer)
				if breakStatement {
					if packetNumber != 0 || index != 0 {
						buffer, nextResponse, err = cut2DSlice(buffer, packetNumber, index)
						if err != nil {
							return err
						}
					}
					break
				}
			}
			fmt.Println("response")
			fmt.Println(buffer)
			responseHeader, responseMessage, continueStatement, nextDataPacketType, stmt, err = Decode(buffer, dataPacketType, false, stmt)
			if err != nil {
				logger.Error(Emoji+"failed to read the response message in proxy", zap.Error(err), zap.Any("proxy port", port))
				return err
			}
			readRequestDelay := time.Since(started)
			for _, packet := range buffer {
				_, err = clientConn.Write(packet)
				if err != nil {
					logger.Error(Emoji+"failed to write response to the client", zap.Error(err))
					return err
				}
			}
			oracleResponses = append(oracleResponses, models.OracleResponse{
				Header:    responseHeader,
				Message:   responseMessage,
				ReadDelay: int64(readRequestDelay),
			})
			dataPacketType = nextDataPacketType

			if continueStatement {
				continue
			}
			fmt.Println(oracleResponses)
			break
		}
	}
}

func Decode(Packets [][]byte, dataPacketType models.DataPacketType, isRequest bool, stmt interface{}) (models.OracleHeader, interface{}, bool, models.DataPacketType, interface{}, error) {
	switch models.PacketTypeFromUint8(Packets[0][4]) {
	case models.CONNECT:
		fmt.Println("CONNECT")
		return DecodeConnectPacket(Packets, dataPacketType)
	case models.ACCEPT:
		fmt.Println("ACCEPT")
		return DecodeAcceptPacket(Packets, dataPacketType)
	case models.REFUSE:
		fmt.Println("REFUSE")
		return DecodeRefusePacket(Packets, dataPacketType)
	case models.REDIRECT:
		fmt.Println("REDIRECT")
		return DecodeRedirectPacket(Packets, dataPacketType)
	case models.DATA:
		fmt.Println("DATA")
		switch DecodeDataPacketType(dataPacketType, Packets) {
		case models.OracleConnectionDataMessageType:
			fmt.Println("CONNECTION_DATA")
			return DecodeConnectionDataMessage(Packets)
		case models.OracleRedirectDataMessageType:
			fmt.Println("REDIRECT_DATA")
			return DecodeRedirectDataMessage(Packets)
		case models.OracleProtocolDataMessageType:
			fmt.Println("PROTOCOL_DATA")
			return DecodeOracleProtocolDataMessage(Packets, isRequest)
		case models.OracleDataTypeDataMessageType:
			fmt.Println("DATA_TYPE_DATA")
			return DecodeOracleDataTypeDataMessage(Packets, isRequest)
		case models.OracleFunctionDataMesssageType:
			fmt.Println("FUNCTION_TYPE_DATA")
			return DecodeOracleFunctionDataMessage(Packets)
		case models.OracleAuthPhaseOneDataMessageType:
			fmt.Println("RESP_AUTH_PHASE_ONE")
			return DecodeOracleAuthPhaseOneResponse(Packets)
		case models.OracleAuthPhaseTwoDataMessageType:
			fmt.Println("RESP_AUTH_PHASE_TWO")
			return DecodeOracleAuthPhaseTwoResponse(Packets)
		case models.OraclePiggyBackDataMesssageType:
			fmt.Println("FUNCTION_PIGGY_BACK")
			return DecodeOraclePiggyBackDataMessage(Packets)
		case models.OracleMessageWithDataMessageType:
			fmt.Println("MESSAGE_WITH_DATA")
			return DecodeOracleMessageWithData(Packets, stmt)
		case models.OracleGetDBVersionDataMessageType:
			fmt.Println("GET_DB_VERSION_RESPONSE")
			return DecodeOracleGetDBVersion(Packets)
		default:
			isAdvNego := checkforAdvanceNego(Packets)
			if isAdvNego {
				fmt.Println("ADV_NEGO")
				return DecodeAdvNegoDataMesssage(Packets, isRequest)
			} else {
				return models.OracleHeader{}, nil, false, models.DefaultDataPacket, nil, errors.New("unsupported Message type")
			}
		}
	case models.RESEND:
		fmt.Println("RESEND")
		return DecodeResendPacket(Packets, dataPacketType)
	case models.MARKER:
		fmt.Println("MARKER")
		return DecodeMarkerPacket(Packets, dataPacketType)
	case models.CTRL:
		fmt.Println("CTRL")
		return DecodeControlPacket(Packets, dataPacketType)
	case models.ATTN:
		return models.OracleHeader{}, nil, false, models.DefaultDataPacket, nil, errors.New("unsupported Message type ATTN")
	case models.HIGHEST:
		return models.OracleHeader{}, nil, false, models.DefaultDataPacket, nil, errors.New("unsupported Message type HIGHEST")
	case models.ACK:
		return models.OracleHeader{}, nil, false, models.DefaultDataPacket, nil, errors.New("unsupported Message type ACK")
	case models.NULL:
		return models.OracleHeader{}, nil, false, models.DefaultDataPacket, nil, errors.New("unsupported Message type NULL")
	case models.ABORT:
		return models.OracleHeader{}, nil, false, models.DefaultDataPacket, nil, errors.New("unsupported Message type ABORT")
	default:
		return models.OracleHeader{}, nil, false, models.DefaultDataPacket, nil, errors.New("unsupported Message type ABORT")
	}
}

func checkforAdvanceNego(Packets [][]byte) bool {
	var packetData []byte
	for _, slice := range Packets {
		packetData = append(packetData, slice...)
	}
	num, _ := session.GetInt64(4, false, true, packetData, 10)
	return num == 0xDEADBEEF
}

func recievedAllPackets(Packets [][]byte) (bool, int, int) {
	var lengthOfPacketsProvided int
	if session.Context != nil && session.Context.Version >= 315 {
		lengthOfPacketsProvided = int(binary.BigEndian.Uint32(Packets[0][0:]))
	} else {
		lengthOfPacketsProvided = int(binary.BigEndian.Uint16(Packets[0][0:]))
	}
	lengthOfPacketsReceived := 0
	for i, packet := range Packets {
		lengthOfPacketsReceived += len(packet)
		if lengthOfPacketsReceived > int(lengthOfPacketsProvided) {
			return true, i, len(packet) - (lengthOfPacketsReceived - lengthOfPacketsProvided)
		}
	}
	if lengthOfPacketsReceived == int(lengthOfPacketsProvided) {
		return true, 0, 0
	} else {
		return false, 0, 0
	}
}

func DecodeDataPacketType(dataPacketType models.DataPacketType, Packets [][]byte) models.DataPacketType {
	if dataPacketType != models.DefaultDataPacket {
		return dataPacketType
	} else {
		var packetData []byte
		for _, slice := range Packets {
			packetData = append(packetData, slice...)
		}
		return models.DataPacketTypeFromInt(int(packetData[10]))
	}
}

func cut2DSlice(slice [][]byte, row int, col int) ([][]byte, [][]byte, error) {
	if row < 0 || row >= len(slice) || col < 0 || col >= len(slice[0]) {
		return nil, nil, fmt.Errorf("invalid row or column index")
	}
	// Copy the rows above and below the cut point
	above := make([][]byte, row)
	below := make([][]byte, len(slice)-row-1)
	copy(above, slice[:row])
	copy(below, slice[row+1:])
	// Split the row at the cut point
	left := make([]byte, col)
	right := make([]byte, len(slice[0])-col)
	for i := range slice {
		if i == row {
			copy(left, slice[i][:col])
			copy(right, slice[i][col:])
			break
		}
	}
	above = append(above, left)
	below = append(below, right)
	return above, below, nil
}

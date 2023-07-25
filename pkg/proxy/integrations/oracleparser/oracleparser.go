package oracleparser

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/sijms/go-ora/v2/network"

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

// var m map[int]RequestResponse

// Determines whether the outgoing packets belong to oracle
func IsOutgoingOracle(buffer []byte) bool {
	messageLength := uint32(binary.BigEndian.Uint16(buffer[0:2]))
	return int(messageLength) == len(buffer)
}

// Processes the Oracle packets
func ProcessOraclePackets(reuqestBuffer []byte, clientConn, destConn net.Conn, h *hooks.Hook, logger *zap.Logger, FilterPid *bool, port uint32, kernalPid uint32) {

	switch models.GetMode() {
	case models.MODE_RECORD:
		fmt.Println("record mode")
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

// // This function facilitates bidirectional communication by reading packets from both the client and the server,
// // and subsequently writing the appropriate responses back to both entities. Meanwhile this fuction also creates
// // mocks in record mode and use mocks in test mode.
// func ReadWriteProtocol(firstbuffer []byte, clientConn, destConn net.Conn, logger *zap.Logger, port uint32, FilterPid *bool, kernalPid uint32, h *hooks.Hook) error {

// 	fmt.Println(Emoji, "trying to forward requests to target: ", destConn.RemoteAddr().String())

// 	defer destConn.Close()

// 	// Create channels
// 	destinationWriteChannel := make(chan []byte)
// 	clientWriteChannel := make(chan []byte)

// 	m = make(map[int]RequestResponse)

// 	// fmt.Println("writing buffer to destination", firstbuffer)

// 	_, err := destConn.Write(firstbuffer)
// 	if err != nil {
// 		logger.Error(Emoji+"failed to write request message to the destination server", zap.Error(err))
// 		return err
// 	}

// 	go func() {
// 		for {
// 			if *FilterPid {
// 				err, pid := h.GetApplicationPID()
// 				if err != nil {
// 					logger.Error(Emoji+"failed to get application pid after filtering")
// 					return
// 				}
// 				if (kernalPid == uint32(pid)) {
// 						FillUncapturedMocks()
// 						break
// 					} else {
// 						break
// 					}
// 			}
// 		}
// 	}()

// 	Requests := [][]byte{firstbuffer}
// 	Responses := make([][]byte, 0)
// 	isResponseDone := false
// 	RequestNum := 0

// 	for {

// 		fmt.Println("inside connection")

// 		// go routine to read from client
// 		go func() {
// 			buffer, err := util.ReadBytes(clientConn)
// 			if err != nil {
// 				logger.Error(Emoji+"failed to read the request message in proxy", zap.Error(err), zap.Any("proxy port", port))
// 				return
// 			}

// 			if (!isResponseDone) {
// 				Requests = append(Requests, buffer)
// 			} else {
// 				m[RequestNum] = RequestResponse{
// 					Requests:  Requests,
// 					Responses: Responses,
// 				}
// 				RequestNum += 1
// 				isResponseDone = false
// 				Requests = [][]byte{}
// 				Responses = [][]byte{}
// 				Requests = append(Requests, buffer)
// 			}

// 			// fmt.Println("buffer from client connection")
// 			// fmt.Println(buffer)
// 			// fmt.Println(string(buffer))
// 			destinationWriteChannel <- buffer

// 		}()

// 		// go routine to read from destination
// 		go func() {
// 			buffer, err := util.ReadBytes(destConn)
// 			if err != nil {
// 				logger.Error(Emoji+"failed to read the request message in proxy", zap.Error(err), zap.Any("proxy port", port))
// 				return
// 			}
// 			isResponseDone = true
// 			Responses = append(Responses, buffer)
// 			// fmt.Println("buffer from destination connection")
// 			// fmt.Println(buffer)
// 			// fmt.Println(string(buffer))
// 			clientWriteChannel <- buffer
// 		}()

// 		select {
// 		case requestBuffer := <-destinationWriteChannel:
// 			// Write the request message to the actual destination server
// 			// fmt.Println("writing buffer to destination", requestBuffer)
// 			_, err := destConn.Write(requestBuffer)
// 			if err != nil {
// 				logger.Error(Emoji+"failed to write request message to the destination server", zap.Error(err))
// 				return err
// 			}

// 		case responseBuffer := <-clientWriteChannel:
// 			// Write the response message to the client
// 			// fmt.Println("writing buffer to client", responseBuffer)
// 			_, err := clientConn.Write(responseBuffer)
// 			if err != nil {
// 				logger.Error(Emoji+"failed to write response to the client", zap.Error(err))
// 				return err
// 			}
// 			// fmt.Println(Emoji, "Successfully wrote response to the user client ", destConn.RemoteAddr().String())

// 		}
// 	}

// }

// This function facilitates bidirectional communication by reading packets from both the client and the server,
// and subsequently writing the appropriate responses back to both entities. Meanwhile this fuction also creates
// mocks in record mode and use mocks in test mode.
func ReadWriteProtocol(firstbuffer []byte, clientConn, destConn net.Conn, logger *zap.Logger, port uint32, FilterPid *bool, kernalPid uint32, h *hooks.Hook) error {

	fmt.Println(Emoji, "trying to forward requests to target: ", destConn.RemoteAddr().String())

	defer destConn.Close()

	// _, err := destConn.Write(firstbuffer)
	// if err != nil {
	// 	logger.Error(Emoji+"failed to write request message to the destination server", zap.Error(err))
	// 	return err
	// }

	// go func() {
	// 	for {
	// 		if *FilterPid {
	// 			err, pid := h.GetApplicationPID()
	// 			if err != nil {
	// 				logger.Error(Emoji+"failed to get application pid after filtering")
	// 				return
	// 			}
	// 			if (kernalPid == uint32(pid)) {
	// 					FillUncapturedMocks()
	// 					break
	// 				} else {
	// 					break
	// 				}
	// 		}
	// 	}
	// }()

	// Requests := [][]byte{firstbuffer}
	// Responses := make([][]byte, 0)
	// isResponseDone := false
	// RequestNum := 0

	for {
		fmt.Println("inside connection request")
		var (
			oracleRequests     = []models.OracleRequest{}
			oracleResponses    = []models.OracleResponse{}
			dataPacketType     models.DataPacketType
			nextDataPacketType models.DataPacketType
			nextRequest        [][]byte
			nextResponse       [][]byte
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
					}
				}
				buffer = append(buffer, clientBuffer)
				breakStatement, packetNumber, index := recievedAllPackets(buffer)
				if breakStatement {
					if packetNumber != 0 && index != 0 {
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
			requestHeader, requestMessage, continueStatement, nextDataPacketType, err = Decode(buffer, dataPacketType, true)
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
			// if (!isResponseDone) {
			// 	Requests = append(Requests, buffer)
			// } else {
			// 	m[RequestNum] = RequestResponse{
			// 		Requests:  Requests,
			// 		Responses: Responses,
			// 	}
			// 	RequestNum += 1
			// 	isResponseDone = false
			// 	Requests = [][]byte{}
			// 	Responses = [][]byte{}
			// 	Requests = append(Requests, buffer)
			// }

			// fmt.Println("buffer from client connection")
			// fmt.Println(buffer)
			// fmt.Println(string(buffer))
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
				}
				fmt.Println("this")
				fmt.Println(serverBuffer)
				buffer = append(buffer, serverBuffer)
				breakStatement, packetNumber, index := recievedAllPackets(buffer)
				fmt.Println(breakStatement, packetNumber, index)
				if breakStatement {
					if packetNumber != 0 || index != 0 {
						buffer, nextResponse, err = cut2DSlice(buffer, packetNumber, index)
						fmt.Println(buffer, nextResponse)
						if err != nil {
							return err
						}
					}
					break
				}
			}
			fmt.Println("response")
			fmt.Println(buffer)
			responseHeader, responseMessage, continueStatement, nextDataPacketType, err = Decode(buffer, dataPacketType, false)
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

			// isResponseDone = true
			// Responses = append(Responses, buffer)
			// // fmt.Println("buffer from destination connection")
			// // fmt.Println(buffer)
			// // fmt.Println(string(buffer))
			// clientWriteChannel <- buffer
		}

		// select {
		// case requestBuffer := <-destinationWriteChannel:
		// 	// Write the request message to the actual destination server
		// 	// fmt.Println("writing buffer to destination", requestBuffer)
		// 	_, err := destConn.Write(requestBuffer)
		// 	if err != nil {
		// 		logger.Error(Emoji+"failed to write request message to the destination server", zap.Error(err))
		// 		return err
		// 	}

		// case responseBuffer := <-clientWriteChannel:
		// 	// Write the response message to the client
		// 	// fmt.Println("writing buffer to client", responseBuffer)
		// 	_, err := clientConn.Write(responseBuffer)
		// 	if err != nil {
		// 		logger.Error(Emoji+"failed to write response to the client", zap.Error(err))
		// 		return err
		// 	}
		// 	// fmt.Println(Emoji, "Successfully wrote response to the user client ", destConn.RemoteAddr().String())

		// }
	}

}

// func FillUncapturedMocks() {

// 	for key, rr := range m {
// 		fmt.Printf("Key: %d\n", key)
// 		fmt.Println("Requests:")
// 		for i, req := range rr.Requests {
// 			fmt.Printf("\tRequest %d: %s\n", i, string(req))
// 			fmt.Println(req)

// 		}
// 		// fmt.Println(Decode(rr.Requests))
// 		fmt.Println("Responses:")
// 		for i, res := range rr.Responses {
// 			fmt.Printf("\tResponse %d: %s\n", i, string(res))
// 			fmt.Println(res)
// 		}
// 		// Decode(rr.Responses)
// 		fmt.Println()
// 	}
// }

func Decode(Packets [][]byte, dataPacketType models.DataPacketType, isRequest bool) (models.OracleHeader, interface{}, bool, models.DataPacketType, error) {
	switch network.PacketType(Packets[0][4]) {
	case network.CONNECT:
		fmt.Println("CONNECT")
		return DecodeConnectPacket(Packets, dataPacketType)
	case network.ACCEPT:
		fmt.Println("ACCEPT")
		return DecodeAcceptPacket(Packets, dataPacketType)
	case network.REFUSE:
		fmt.Println("REFUSE")
		return DecodeRefusePacket(Packets, dataPacketType)
	case network.REDIRECT:
		fmt.Println("REDIRECT")
		return DecodeRedirectPacket(Packets, dataPacketType)
	case network.DATA:
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
			return DecodeOracleFunctionDataMessage(Packets, isRequest)
		}
	case network.RESEND:
		fmt.Println("RESEND")
		return DecodeResendPacket(Packets, dataPacketType)
	case network.MARKER:
		fmt.Println("MARKER")
		return DecodeMarkerPacket(Packets, dataPacketType)
	case network.CTRL:
		fmt.Println("CTRL")
		return DecodeControlPacket(Packets, dataPacketType)
	case network.ATTN:
		return models.OracleHeader{}, nil, false, models.Default, errors.New("unsupported Message type ATTN")
	case network.HIGHEST:
		return models.OracleHeader{}, nil, false, models.Default, errors.New("unsupported Message type HIGHEST")
	case network.ACK:
		return models.OracleHeader{}, nil, false, models.Default, errors.New("unsupported Message type ACK")
	case network.NULL:
		return models.OracleHeader{}, nil, false, models.Default, errors.New("unsupported Message type NULL")
	case network.ABORT:
		return models.OracleHeader{}, nil, false, models.Default, errors.New("unsupported Message type ABORT")
	default:
		return models.OracleHeader{}, nil, false, models.Default, nil
	}
	return models.OracleHeader{}, nil, false, models.Default, nil
}

func recievedAllPackets(Packets [][]byte) (bool, int, int) {
	var lengthOfPacketsProvided int
	if session.Context != nil && session.Context.Version >= 315 {
		lengthOfPacketsProvided = int(binary.BigEndian.Uint32(Packets[0][0:]))
	} else {
		lengthOfPacketsProvided = int(binary.BigEndian.Uint16(Packets[0][0:]))
	}
	lengthOfPacketsReceived := 0
	fmt.Println(lengthOfPacketsProvided)
	for i, packet := range Packets {
		lengthOfPacketsReceived += len(packet)
		if lengthOfPacketsReceived > int(lengthOfPacketsProvided) {
			fmt.Println(len(packet) - (lengthOfPacketsReceived - lengthOfPacketsProvided))
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
	if dataPacketType != models.Default {
		return dataPacketType
	} else {
		var packetData []byte
		for _, slice := range Packets {
			packetData = append(packetData, slice...)
		}
		switch models.DataPacketType(packetData[10]) {
		case models.OracleProtocolDataMessageType:
			return models.OracleProtocolDataMessageType
		case models.OracleDataTypeDataMessageType:
			return models.OracleDataTypeDataMessageType
		case models.OracleFunctionDataMesssageType:
			return models.OracleFunctionDataMesssageType
		}
	}
	return models.Default

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

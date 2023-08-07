package oracleparser

import (
	"encoding/binary"

	"go.keploy.io/server/pkg/models"
)

func DecodeAdvNegoDataMesssage(Packets [][]byte, isRequest bool) (models.OracleHeader, interface{}, bool, models.DataPacketType, interface{}, error) {
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
		DataMessageType: models.OracleAdvNegoMessagetype,
	}
	index := 10
	if isRequest {
		var advNegoHeaderMessage models.AdvNegoHeaderMessage
		_, index = session.GetInt64(4, false, true, packetData, index)
		advNegoHeaderMessage.Length, index = session.GetInt(2, false, true, packetData, index)
		advNegoHeaderMessage.Version, index = session.GetInt(4, false, true, packetData, index)
		advNegoHeaderMessage.ServCount, index = session.GetInt(2, false, true, packetData, index)
		advNegoHeaderMessage.ErrFlags, index = session.GetInt(1, false, true, packetData, index)
		var Service_Data_List []models.ServiceData
		for i := 0; i < advNegoHeaderMessage.ServCount; i++ {
			var Service_Data models.ServiceData
			var advNegoServiceHeader models.AdvNegoServiceHeader
			var serviceType int
			serviceType, index = session.GetInt(2, false, true, packetData, index)
			advNegoServiceHeader.ServiceType = models.ServiceTypeFromInt(serviceType)
			advNegoServiceHeader.ServiceSubPackets, index = session.GetInt(2, false, true, packetData, index)
			_, index = session.GetInt(4, false, true, packetData, index)
			switch advNegoServiceHeader.ServiceType {
			case models.AUTH_SERVICE:
				var authServiceData models.AuthServiceData
				var ServiceIds []models.ServicePacket
				var ServiceNames []models.ServicePacket
				var ServiceId models.ServicePacket
				var ServiceName models.ServicePacket
				authServiceData.Version, index = readVersion(packetData, index)
				authServiceData.Unknown, index = readUB2(packetData, index)
				authServiceData.Status, index = readStatus(packetData, index)
				for j := 0; j < advNegoServiceHeader.ServiceSubPackets-3; j++ {
					ServiceId, index = readUB1(packetData, index)
					ServiceName, index = readString(packetData, index)
					ServiceIds = append(ServiceIds, ServiceId)
					ServiceNames = append(ServiceNames, ServiceName)
				}
				authServiceData.ServiceIds = ServiceIds
				authServiceData.ServiceNames = ServiceNames
				Service_Data.AdvNegoServiceHeader = advNegoServiceHeader
				Service_Data.ServiceData = authServiceData
			case models.ENCRYPT_SERVICE:
				var encrytionServiceData models.EncrytionServiceData
				encrytionServiceData.Version, index = readVersion(packetData, index)
				encrytionServiceData.ServiceIds, index = readBytes(packetData, index)
				_, index = readUB1(packetData, index)
				Service_Data.AdvNegoServiceHeader = advNegoServiceHeader
				Service_Data.ServiceData = encrytionServiceData
			case models.DATA_INTEGRITY_SERVICE:
				var dataServiceData models.DataServiceData
				dataServiceData.Version, index = readVersion(packetData, index)
				dataServiceData.ServiceIds, index = readBytes(packetData, index)
				Service_Data.AdvNegoServiceHeader = advNegoServiceHeader
				Service_Data.ServiceData = dataServiceData
			case models.SUPERVISOR_SERVICE:
				var supervisorServiceData models.SupervisorServiceData
				supervisorServiceData.Version, index = readVersion(packetData, index)
				supervisorServiceData.ServiceIds, index = readBytes(packetData, index)
				supervisorServiceData.ServiceArray, index = readUB2Array(packetData, index)
				Service_Data.AdvNegoServiceHeader = advNegoServiceHeader
				Service_Data.ServiceData = supervisorServiceData
			}
			Service_Data_List = append(Service_Data_List, Service_Data)
		}
		requestMessage.DataMessage = models.OracleAdvNegoMessage{
			AdvNegoHeaderMessage: advNegoHeaderMessage,
			ServiceDataList:      Service_Data_List,
		}
	} else {
		var advNegoHeaderMessage models.AdvNegoHeaderMessage
		_, index = session.GetInt64(4, false, true, packetData, index)
		advNegoHeaderMessage.Length, index = session.GetInt(2, false, true, packetData, index)
		advNegoHeaderMessage.Version, index = session.GetInt(4, false, true, packetData, index)
		advNegoHeaderMessage.ServCount, index = session.GetInt(2, false, true, packetData, index)
		advNegoHeaderMessage.ErrFlags, index = session.GetInt(1, false, true, packetData, index)
		var Service_Data_List []models.ServiceData
		for i := 0; i < advNegoHeaderMessage.ServCount; i++ {
			var Service_Data models.ServiceData
			var advNegoServiceHeader models.AdvNegoServiceHeader
			var serviceType int
			serviceType, index = session.GetInt(2, false, true, packetData, index)
			advNegoServiceHeader.ServiceType = models.ServiceTypeFromInt(serviceType)
			advNegoServiceHeader.ServiceSubPackets, index = session.GetInt(2, false, true, packetData, index)
			_, index = session.GetInt(4, false, true, packetData, index)
			switch advNegoServiceHeader.ServiceType {
			case models.AUTH_SERVICE:
				var authServiceData models.AuthServiceDataResponse
				authServiceData.Version, index = readVersion(packetData, index)
				authServiceData.Status, index = readStatus(packetData, index)
				if authServiceData.Status.Value == 0xFAFF && advNegoServiceHeader.ServiceSubPackets > 2 {
					// get 1 byte with header
					_, index = readUB1(packetData, index)
					authServiceData.ServiceName, index = readString(packetData, index)
					if advNegoServiceHeader.ServiceSubPackets > 4 {
						_, index = readVersion(packetData, index)
						_, index = readUB4(packetData, index)
						_, index = readUB4(packetData, index)
					}
				}
				Service_Data.AdvNegoServiceHeader = advNegoServiceHeader
				Service_Data.ServiceData = authServiceData
			case models.ENCRYPT_SERVICE:
				var encrytionServiceData models.EncrytionServiceDataResponse
				encrytionServiceData.Version, index = readVersion(packetData, index)
				encrytionServiceData.ServiceId, index = readUB1(packetData, index)
				Service_Data.AdvNegoServiceHeader = advNegoServiceHeader
				Service_Data.ServiceData = encrytionServiceData
			case models.DATA_INTEGRITY_SERVICE:
				var dataServiceData models.DataServiceDataResponse
				dataServiceData.Version, index = readVersion(packetData, index)
				dataServiceData.ServiceId, index = readUB1(packetData, index)
				if advNegoServiceHeader.ServiceSubPackets == 8 {
					dataServiceData.DhGenLen, index = readUB2(packetData, index)
					dataServiceData.DPrimLen, index = readUB2(packetData, index)
					dataServiceData.GenBytes, index = readBytes(packetData, index)
					dataServiceData.PrimeBytes, index = readBytes(packetData, index)
					dataServiceData.ServerPublicKeyBytes, index = readBytes(packetData, index)
					dataServiceData.IV, index = readBytes(packetData, index)
				}
				Service_Data.AdvNegoServiceHeader = advNegoServiceHeader
				Service_Data.ServiceData = dataServiceData
			case models.SUPERVISOR_SERVICE:
				var supervisorServiceData models.SupervisorServiceDataResponse
				supervisorServiceData.Version, index = readVersion(packetData, index)
				supervisorServiceData.Status, index = readStatus(packetData, index)
				supervisorServiceData.ServiceArray, index = readUB2Array(packetData, index)
				Service_Data.AdvNegoServiceHeader = advNegoServiceHeader
				Service_Data.ServiceData = supervisorServiceData
			}
			Service_Data_List = append(Service_Data_List, Service_Data)
		}
		requestMessage.DataMessage = models.OracleAdvNegoMessage{
			AdvNegoHeaderMessage: advNegoHeaderMessage,
			ServiceDataList:      Service_Data_List,
		}
	}
	return requestHeader, requestMessage, false, models.DefaultDataPacket, nil, nil
}

func readUB2Array(packetData []byte, index int) (models.ServicePacket, int) {
	var servicePacket models.ServicePacket
	var leng int
	var value []int

	length, receivedType, index := ReadPacketHeader(packetData, index)
	servicePacket.Length = length
	servicePacket.Type = receivedType
	_, index = session.GetInt(4, false, true, packetData, index)
	_, index = session.GetInt(2, false, true, packetData, index)
	leng, index = session.GetInt(4, false, true, packetData, index)
	for i := 0; i < leng; i++ {
		tempValue := 0
		tempValue, index = session.GetInt(2, false, true, packetData, index)
		value = append(value, tempValue)
	}
	servicePacket.Value = value
	return servicePacket, index
}

func readUB4(packetData []byte, index int) (models.ServicePacket, int) {
	var servicePacket models.ServicePacket
	length, receivedType, index := ReadPacketHeader(packetData, index)
	servicePacket.Length = length
	servicePacket.Type = receivedType
	servicePacket.Value, index = session.GetInt(4, false, true, packetData, index)
	return servicePacket, index
}

func readVersion(packetData []byte, index int) (models.ServicePacket, int) {
	var servicePacket models.ServicePacket
	length, receivedType, index := ReadPacketHeader(packetData, index)
	servicePacket.Length = length
	servicePacket.Type = receivedType
	servicePacket.Value, index = session.GetInt(4, false, true, packetData, index)
	return servicePacket, index
}

func readStatus(packetData []byte, index int) (models.ServicePacket, int) {
	var servicePacket models.ServicePacket
	length, receivedType, index := ReadPacketHeader(packetData, index)
	servicePacket.Length = length
	servicePacket.Type = receivedType
	servicePacket.Value, index = session.GetInt(2, false, true, packetData, index)
	return servicePacket, index
}

func readUB2(packetData []byte, index int) (models.ServicePacket, int) {
	var servicePacket models.ServicePacket
	length, receivedType, index := ReadPacketHeader(packetData, index)
	servicePacket.Length = length
	servicePacket.Type = receivedType
	servicePacket.Value, index = session.GetInt(2, false, true, packetData, index)
	return servicePacket, index
}

func readBytes(packetData []byte, index int) (models.ServicePacket, int) {
	var servicePacket models.ServicePacket
	length, receivedType, index := ReadPacketHeader(packetData, index)
	servicePacket.Length = length
	servicePacket.Type = receivedType
	servicePacket.Value, index = session.GetBytes(length, packetData, index)
	return servicePacket, index
}

func ReadPacketHeader(packetData []byte, index int) (length int, receivedType int, retIndex int) {
	length, retIndex = session.GetInt(2, false, true, packetData, index)
	receivedType, retIndex = session.GetInt(2, false, true, packetData, retIndex)
	return length, receivedType, retIndex
}

func readUB1(packetData []byte, index int) (models.ServicePacket, int) {
	var servicePacket models.ServicePacket
	length, receivedType, index := ReadPacketHeader(packetData, index)
	servicePacket.Length = length
	servicePacket.Type = receivedType
	servicePacket.Value, index = session.GetByte(packetData, index)
	return servicePacket, index
}

func readString(packetData []byte, index int) (models.ServicePacket, int) {
	var servicePacket models.ServicePacket
	length, receivedType, index := ReadPacketHeader(packetData, index)
	servicePacket.Length = length
	servicePacket.Type = receivedType
	servicePacket.Value, index = session.GetBytes(length, packetData, index)
	return servicePacket, index
}

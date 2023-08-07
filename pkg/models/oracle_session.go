package models

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"

	"github.com/sijms/go-ora/v2/converters"
)

func PacketTypeFromUint8(val uint8) PacketType {
	return PacketTypeValues[val]
}

func PacketTypeFromString(val string) uint8 {
	return ReversePacketTypeValues[PacketType(val)]
}

func ControlTypeFromInt(val int) ControlType {
	return ControlTypeValues[val]
}

func ControlTypeFromString(val string) int {
	return ReverseControlTypeValues[ControlType(val)]
}

func ControlErrorFromInt(val int) ControlError {
	return ControlErrorValues[val]
}

func ControlErrorFromString(val string) int {
	return ReverseControlErrorValues[ControlError(val)]
}

func MarkerTypeFromInt(val int) MarkerType {
	return MarkerTypeValues[val]
}

func MarkerTypeFromString(val string) int {
	return ReverseMarkerTypeValues[MarkerType(val)]
}

func StmtTypeFromInt(val int) StmtType {
	return StmtTypeValues[val]
}

func StmtTypeFromString(val string) int {
	return ReverseStmtTypeValues[StmtType(val)]
}

func DirectionFromInt(val int) ParameterDirection {
	return DirectionValues[val]
}

func DirectionFromString(val string) int {
	return ReverseDirectionValues[ParameterDirection(val)]
}

func TNSTypeFromInt(val int) TNSType {
	return TnsTypeValues[val]
}

func TNSTypeFromString(val string) int {
	return ReverseTNSTypeValues[TNSType(val)]
}

func DataPacketTypeFromInt(val int) DataPacketType {
	return DataPacketTypeValues[val]
}

func DataPacketTypeFromString(val string) int {
	return ReverseDataPacketTypeValues[DataPacketType(val)]
}

func FunctionTypeFromInt(val int) FunctionType {
	return FunctionTypeValues[val]
}

func FunctionTypeFromString(val string) int {
	return ReverseFunctionTypeValues[FunctionType(val)]
}

func AuthModeFromInt(val int) AuthMode {
	return AuthModeValues[val]
}

func AuthModeFromString(val string) int {
	return ReverseAuthModeValues[AuthMode(val)]
}

func PiggyBackTypeFromInt(val int) PiggyBackType {
	return PiggyBackTypeValues[val]
}

func PiggyBackTypeFromString(val string) int {
	return ReversePiggyBackTypeValues[PiggyBackType(val)]
}

func DataTypeFromInt(val int) DataType {
	return DataTypeValues[val]
}

func DataTypeFromString(val string) int {
	return ReverseDataTypeValues[DataType(val)]
}

func ServiceTypeFromInt(val int) ServiceType {
	return ServiceTypeValues[val]
}

func ServiceTypeFromString(val string) int {
	return ReverseServiceTypeValues[ServiceType(val)]
}

func (err *OracleError) Error() string {
	if len(err.ErrMsg) == 0 {
		err.translate()
	}
	return err.ErrMsg
}

func (err *OracleError) translate() {
	switch err.ErrCode {
	case 1:
		err.ErrMsg = "ORA-00001: Unique constraint violation"
	case 900:
		err.ErrMsg = "ORA-00900: Invalid SQL statement"
	case 901:
		err.ErrMsg = "ORA-00901: Invalid CREATE command"
	case 902:
		err.ErrMsg = "ORA-00902: Invalid data type"
	case 903:
		err.ErrMsg = "ORA-00903: Invalid table name"
	case 904:
		err.ErrMsg = "ORA-00904: Invalid identifier"
	case 905:
		err.ErrMsg = "ORA-00905: Misspelled keyword"
	case 906:
		err.ErrMsg = "ORA-00906: Missing left parenthesis"
	case 907:
		err.ErrMsg = "ORA-00907: Missing right parenthesis"
	case 12631:
		err.ErrMsg = "ORA-12631: Username retrieval failed"
	case 12564:
		err.ErrMsg = "ORA-12564: TNS connection refused"
	case 12506:
		err.ErrMsg = "ORA-12506: TNS:listener rejected connection based on service ACL filtering"
	case 12514:
		err.ErrMsg = "ORA-12514: TNS:listener does not currently know of service requested in connect descriptor"
	case 3135:
		err.ErrMsg = "ORA-03135: connection lost contact"
	default:
		err.ErrMsg = "ORA-" + strconv.Itoa(err.ErrCode)
	}
}

func (session *PacketSession) GetInt64(size int, compress bool, bigEndian bool, buffer []byte, index int) (int64, int) {
	var ret int64
	negFlag := false
	if compress {
		rb := buffer[index]
		index += 1
		size = int(rb)
		if size&0x80 > 0 {
			negFlag = true
			size = size & 0x7F
		}
		bigEndian = true
	}
	rb := buffer[index : index+size]
	index += size
	temp := make([]byte, 8)
	if bigEndian {
		copy(temp[8-size:], rb)
		ret = int64(binary.BigEndian.Uint64(temp))
	} else {
		copy(temp[:size], rb)
		//temp = append(pck.buffer[pck.index: pck.index + size], temp...)
		ret = int64(binary.LittleEndian.Uint64(temp))
	}
	if negFlag {
		ret = ret * -1
	}
	return ret, index
}

func (session *PacketSession) GetInt(size int, compress bool, bigEndian bool, buffer []byte, index int) (int, int) {
	temp, returnindex := session.GetInt64(size, compress, bigEndian, buffer, index)
	return int(temp), returnindex
}

func (session *PacketSession) GetKeyVal(buffer []byte, index int) (key []byte, val []byte, num int, returnIndex int) {
	key, index = session.GetDlc(buffer, index)
	val, index = session.GetDlc(buffer, index)
	num, returnIndex = session.GetInt(4, true, true, buffer, index)
	return
}

func (session *PacketSession) GetDlc(buffer []byte, index int) (output []byte, returnIndex int) {
	var length int
	length, returnIndex = session.GetInt(4, true, true, buffer, index)
	if length > 0 {
		output, returnIndex = session.GetClr(buffer, returnIndex)
		if len(output) > length {
			output = output[:length]
		}
	}
	return
}

func (session *PacketSession) GetClr(buffer []byte, index int) (output []byte, returnIndex int) {
	var nb byte
	nb, index = session.GetByte(buffer, index)
	if nb == 0 || nb == 0xFF {
		output = nil
		returnIndex = index
		return
	}
	chunkSize := int(nb)
	var chunk []byte
	var tempBuffer bytes.Buffer
	if chunkSize == 0xFE {
		for chunkSize > 0 {
			if session.UseBigClrChunks {
				chunkSize, index = session.GetInt(4, true, true, buffer, index)
			} else {
				nb, index = session.GetByte(buffer, index)
				chunkSize = int(nb)
			}
			chunk, index = session.GetBytes(chunkSize, buffer, index)
			tempBuffer.Write(chunk)
		}
	} else {
		chunk, index = session.GetBytes(chunkSize, buffer, index)
		tempBuffer.Write(chunk)
	}
	returnIndex = index
	output = tempBuffer.Bytes()
	return
}

func (session *PacketSession) GetByte(buffer []byte, index int) (uint8, int) {
	rb := buffer[index]
	return rb, index + 1
}

func (session *PacketSession) GetBytes(length int, buffer []byte, index int) ([]byte, int) {
	rb := buffer[index : index+length]
	index += length
	return rb, index
}

func (session *PacketSession) GetNullTerminatedString(buffer []byte, index int) (string, int) {
	startingIndex := index
	for index < len(buffer) && buffer[index] != 0 {
		index++
	}
	driverName := string(buffer[startingIndex:index])
	index += 1
	return driverName, index
}

func (session *PacketSession) GetNullTerminatedArray(buffer []byte, index int) ([]byte, int) {
	var arrayTerminator []byte
	for index < len(buffer) && buffer[index] != 0 {
		arrayTerminator = append(arrayTerminator, buffer[index])
		index += 1
	}
	index += 1
	return arrayTerminator, index
}

func (session *PacketSession) GetString(length int, buffer []byte, index int) (string, int) {
	ret, index := session.GetClr(buffer, index)
	return string(ret[:length]), index
}

func NewStmt(text string) *Stmt {
	ret := &Stmt{
		ReSendParDef: false,
		Parse:        true,
		Execute:      true,
		Define:       false,
	}
	ret.Text = text
	ret.HasBLOB = false
	ret.HasLONG = false
	ret.DisableCompression = false
	ret.ArrayBindCount = 0
	ret.ScnForSnapshot = make([]int, 2)
	// get stmt type
	uCmdText := strings.ToUpper(text)
	for {
		uCmdText = strings.TrimSpace(uCmdText) // trim leading white-space
		if strings.HasPrefix(uCmdText, "--") {
			i := strings.Index(uCmdText, "\n")
			if i <= 0 {
				break
			}
			uCmdText = uCmdText[i+1:]
		} else if strings.HasPrefix(uCmdText, "/*") {
			i := strings.Index(uCmdText, "*/")
			if i <= 0 {
				break
			}
			uCmdText = uCmdText[i+2:]
		} else {
			break
		}
	}
	if strings.HasPrefix(uCmdText, "(") {
		uCmdText = uCmdText[1:]
	}
	if strings.HasPrefix(uCmdText, "SELECT") || strings.HasPrefix(uCmdText, "WITH") {
		ret.StmtType = SELECT
	} else if strings.HasPrefix(uCmdText, "INSERT") ||
		strings.HasPrefix(uCmdText, "MERGE") {
		ret.StmtType = DML
		ret.BulkExec = true
	} else if strings.HasPrefix(uCmdText, "UPDATE") ||
		strings.HasPrefix(uCmdText, "DELETE") {
		ret.StmtType = DML
	} else if strings.HasPrefix(uCmdText, "DECLARE") || strings.HasPrefix(uCmdText, "BEGIN") {
		ret.StmtType = PLSQL
	} else {
		ret.StmtType = OTHERS
	}
	// returning clause
	var err error
	if ret.StmtType != PLSQL {
		ret.HasReturnClause, err = regexp.MatchString(`\bRETURNING\b\s+(\w+\s*,\s*)*\s*\w+\s+\bINTO\b`, uCmdText)
		if err != nil {
			ret.HasReturnClause = false
		}
	}
	return ret
}

func (dataSet *DataSet) SetBitVector(bitVector []byte) {
	index := dataSet.ColumnCount / 8
	if dataSet.ColumnCount%8 > 0 {
		index++
	}
	if len(bitVector) > 0 {
		for x := 0; x < len(bitVector); x++ {
			for i := 0; i < 8; i++ {
				if (x*8)+i < dataSet.ColumnCount {
					dataSet.Cols[(x*8)+i].GetDataFromServer = bitVector[x]&(1<<i) > 0
				}
			}
		}
	} else {
		for x := 0; x < len(dataSet.Cols); x++ {
			dataSet.Cols[x].GetDataFromServer = true
		}
	}
}

func (par *ParameterInfo) Load(packetData []byte, index int, session PacketSession) (ParameterInfo, int) {
	var dataType byte
	var scale int
	var num1 int
	var bName []byte
	par.GetDataFromServer = true
	dataType, index = session.GetByte(packetData, index)
	par.DataType = TNSTypeFromInt(int(dataType))
	par.Flag, index = session.GetByte(packetData, index)
	par.Precision, index = session.GetByte(packetData, index)
	switch par.DataType {
	case NUMBER:
		fallthrough
	case TimeStampDTY:
		fallthrough
	case TimeStampTZ_DTY:
		fallthrough
	case IntervalDS_DTY:
		fallthrough
	case TIMESTAMP:
		fallthrough
	case TIMESTAMPTZ:
		fallthrough
	case IntervalDS:
		fallthrough
	case TimeStampLTZ_DTY:
		fallthrough
	case TimeStampeLTZ:
		scale, index = session.GetInt(2, true, true, packetData, index)
		if scale == -127 {
			par.Precision = uint8(math.Ceil(float64(par.Precision) * 0.30103))
			par.Scale = 0xFF
		} else {
			par.Scale = uint8(scale)
		}
	default:
		par.Scale, index = session.GetByte(packetData, index)
	}
	if par.DataType == NUMBER && par.Precision == 0 && (par.Scale == 0 || par.Scale == 0xFF) {
		par.Precision = 38
		par.Scale = 0xFF
	}
	par.MaxLen, index = session.GetInt(4, true, true, packetData, index)
	switch par.DataType {
	case ROWID:
		par.MaxLen = 128
	case DATE:
		par.MaxLen = 7
	case IBFloat:
		par.MaxLen = 4
	case IBDouble:
		par.MaxLen = 8
	case TimeStampTZ_DTY:
		par.MaxLen = 13
	case IntervalYM_DTY:
		fallthrough
	case IntervalDS_DTY:
		fallthrough
	case IntervalYM:
		fallthrough
	case IntervalDS:
		par.MaxLen = 11
	}
	par.MaxNoOfArrayElements, index = session.GetInt(4, true, true, packetData, index)
	if session.TTCVersion >= 10 {
		par.ContFlag, index = session.GetInt(8, true, true, packetData, index)
	} else {
		par.ContFlag, index = session.GetInt(4, true, true, packetData, index)
	}
	par.ToID, index = session.GetDlc(packetData, index)
	par.Version, index = session.GetInt(2, true, true, packetData, index)
	par.CharsetID, index = session.GetInt(2, true, true, packetData, index)
	par.CharsetForm, index = session.GetInt(1, false, false, packetData, index)
	par.MaxCharLen, index = session.GetInt(4, true, true, packetData, index)
	if session.TTCVersion >= 8 {
		par.Oaccollid, index = session.GetInt(4, true, true, packetData, index)
	}
	num1, index = session.GetInt(1, false, false, packetData, index)
	par.AllowNull = num1 > 0
	_, index = session.GetByte(packetData, index)
	bName, index = session.GetDlc(packetData, index)
	par.Name = session.StrConv.Decode(bName)
	_, index = session.GetDlc(packetData, index)
	bName, index = session.GetDlc(packetData, index)
	par.TypeName = strings.ToUpper(session.StrConv.Decode(bName))
	if par.TypeName == "XMLTYPE" {
		par.DataType = XMLType
		par.IsXmlType = true
	}
	parameterInfo := ParameterInfo{
		DataType:             par.DataType,
		Flag:                 par.Flag,
		Precision:            par.Precision,
		Scale:                par.Scale,
		MaxLen:               par.MaxLen,
		MaxNoOfArrayElements: par.MaxNoOfArrayElements,
		ContFlag:             par.ContFlag,
		ToID:                 par.ToID,
		Version:              par.Version,
		CharsetID:            par.CharsetID,
		CharsetForm:          par.CharsetForm,
		MaxCharLen:           par.MaxCharLen,
		Oaccollid:            par.Oaccollid,
		AllowNull:            par.AllowNull,
		Name:                 par.Name,
		TypeName:             par.TypeName,
	}
	if session.TTCVersion < 3 {
		return parameterInfo, index
	}
	_, index = session.GetInt(2, true, true, packetData, index)
	if session.TTCVersion < 6 {
		return parameterInfo, index
	}
	_, index = session.GetInt(4, true, true, packetData, index)
	return parameterInfo, index
}

func (stmt *Stmt) CalculateParameterValue(param *ParameterInfo, session PacketSession, packetData []byte, index int) (OraclePrimeValue, int) {
	if param.DataType == OCIBlobLocator || param.DataType == OCIClobLocator {
		stmt.HasBLOB = true
	}
	return param.DecodeParameterValue(session, packetData, index)
}

func (par *ParameterInfo) DecodeParameterValue(session PacketSession, packetData []byte, index int) (OraclePrimeValue, int) {
	return par.DecodePrimValue(session, packetData, index, false)
}

func (par *ParameterInfo) DecodePrimValue(session PacketSession, packetData []byte, index int, udt bool) (OraclePrimeValue, int) {
	var err error
	var size int
	var parameterValueList []OraclePrimeValue
	var parameterValue OraclePrimeValue
	var oraclePrimeValue OraclePrimeValue
	var bValue []byte
	par.OPrimValue = nil
	par.BValue = nil
	var rowid *Rowid
	var urowid *Urowid
	var locator []byte
	var decodeObj DecodeObject
	if par.MaxNoOfArrayElements > 0 {
		size, index = session.GetInt(4, true, true, packetData, index)
		par.MaxNoOfArrayElements = size
		if size > 0 {
			pars := make([]ParameterInfo, 0, size)
			for x := 0; x < size; x++ {
				tempPar := par.Clone()
				parameterValue, index = tempPar.DecodeParameterValue(session, packetData, index)
				parameterValueList = append(parameterValueList, parameterValue)
				if x < size-1 {
					_, index = session.GetInt(2, true, true, packetData, index)
				}
				pars = append(pars, tempPar)
			}
			par.OPrimValue = pars
		}
		oraclePrimeValue = OraclePrimeValue{
			Size:               size,
			ParameterValueList: parameterValueList,
		}
		return oraclePrimeValue, index
	}
	if par.DataType == XMLType {
		_, index = session.GetDlc(packetData, index)      // contain toid and some 0s
		_, index = session.GetBytes(3, packetData, index) // 3 0s
		_, index = session.GetInt(4, true, true, packetData, index)
		_, index = session.GetByte(packetData, index)
		_, index = session.GetByte(packetData, index)
	}
	if par.DataType == ROWID {
		rowid, index = NewRowID(session, packetData, index)
		if rowid != nil {
			par.OPrimValue = string(rowid.getBytes())
		}
		oraclePrimeValue = OraclePrimeValue{
			Size:               size,
			ParameterValueList: parameterValueList,
			RowId:              *rowid,
		}
		return oraclePrimeValue, index
	}
	if par.DataType == UROWID {
		urowid, index = NewURowID(session, packetData, index)
		if rowid != nil {
			par.OPrimValue = string(urowid.getBytes())
		}
		oraclePrimeValue = OraclePrimeValue{
			Size:               size,
			ParameterValueList: parameterValueList,
			UrowId:             *urowid,
		}
		return oraclePrimeValue, index
	}
	if (par.DataType == NCHAR || par.DataType == CHAR) && par.MaxCharLen == 0 {
		return oraclePrimeValue, index
	}
	if par.DataType == RAW && par.MaxLen == 0 {
		return oraclePrimeValue, index
	}
	bValue, index = session.GetClr(packetData, index)
	par.BValue = bValue
	oraclePrimeValue = OraclePrimeValue{
		Bvalue: bValue,
	}
	if par.BValue == nil {
		return oraclePrimeValue, index
	}
	switch par.DataType {
	case NCHAR, CHAR, LONG:
		strConv := GetStrConv(session, par.CharsetID)
		par.OPrimValue = strConv.Decode(par.BValue)
	case RAW:
		par.OPrimValue = par.BValue
	case NUMBER:
		if par.Scale == 0 && par.Precision <= 18 {
			par.OPrimValue, err = converters.NumberToInt64(par.BValue)
		} else if par.Scale == 0 && (converters.CompareBytes(par.BValue, converters.Int64MaxByte) > 0 &&
			converters.CompareBytes(par.BValue, converters.Uint64MaxByte) < 0) {
			par.OPrimValue, err = converters.NumberToUInt64(par.BValue)
		} else if par.Scale > 0 {
			par.OPrimValue, err = converters.NumberToString(par.BValue)
		} else {
			par.OPrimValue = converters.DecodeNumber(par.BValue)
		}
	case TimeStampDTY, TimeStampeLTZ, TimeStampLTZ_DTY, TIMESTAMPTZ, TimeStampTZ_DTY:
		fallthrough
	case TIMESTAMP, DATE:
		par.OPrimValue, err = converters.DecodeDate(par.BValue)
	case OCIClobLocator, OCIBlobLocator:
		if !udt {
			locator, index = session.GetClr(packetData, index)
		} else {
			locator = par.BValue

		}
		par.OPrimValue = Lob{
			sourceLocator: locator,
			sourceLen:     len(locator),
			charsetID:     par.CharsetID,
		}
	case OCIFileLocator:
		locator, index = session.GetClr(packetData, index)
		par.OPrimValue = BFile{
			isOpened: false,
			lob: Lob{
				sourceLocator: locator,
				sourceLen:     len(locator),
				charsetID:     par.CharsetID,
			},
		}
	case IBFloat:
		par.OPrimValue = float64(converters.ConvertBinaryFloat(par.BValue))
	case IBDouble:
		par.OPrimValue = converters.ConvertBinaryDouble(par.BValue)
	case IntervalYM_DTY:
		par.OPrimValue = converters.ConvertIntervalYM_DTY(par.BValue)
	case IntervalDS_DTY:
		par.OPrimValue = converters.ConvertIntervalDS_DTY(par.BValue)
	case XMLType:
		decodeObj, index = decodeObject(session, packetData, index, par)
	default:
		return oraclePrimeValue, index
	}
	oraclePrimeValue.DecodeObj = decodeObj
	if err != nil {
		fmt.Println("error")
	}
	return oraclePrimeValue, index
}

func (par *ParameterInfo) Clone() ParameterInfo {
	tempPar := ParameterInfo{}
	tempPar.DataType = par.DataType
	tempPar.CusType = par.CusType
	tempPar.TypeName = par.TypeName
	tempPar.MaxLen = par.MaxLen
	tempPar.MaxCharLen = par.MaxCharLen
	tempPar.CharsetID = par.CharsetID
	tempPar.CharsetForm = par.CharsetForm
	tempPar.Scale = par.Scale
	tempPar.Precision = par.Precision
	return tempPar
}

func NewRowID(session PacketSession, packetData []byte, index int) (*Rowid, int) {
	var temp byte
	var num byte
	temp, index = session.GetByte(packetData, index)
	if temp > 0 {
		ret := new(Rowid)
		ret.RBA, index = session.GetInt64(4, true, true, packetData, index)
		ret.PartitionId, index = session.GetInt64(2, true, true, packetData, index)
		num, index = session.GetByte(packetData, index)
		ret.BlockNumber, index = session.GetInt64(4, true, true, packetData, index)
		ret.SlotNumber, index = session.GetInt64(2, true, true, packetData, index)
		if ret.RBA == 0 && ret.PartitionId == 0 && num == 0 && ret.BlockNumber == 0 && ret.SlotNumber == 0 {
			return nil, index
		}
		return ret, index
	}
	return nil, index
}

func (id *Rowid) getBytes() []byte {
	output := make([]byte, 0, 18)
	output = append(output, convertRowIDToByte(id.RBA, 6)...)
	output = append(output, convertRowIDToByte(id.PartitionId, 3)...)
	output = append(output, convertRowIDToByte(id.BlockNumber, 6)...)
	output = append(output, convertRowIDToByte(id.SlotNumber, 3)...)
	return output
}

func convertRowIDToByte(number int64, size int) []byte {
	var buffer = []byte{
		65, 66, 67, 68, 69, 70, 71, 72,
		73, 74, 75, 76, 77, 78, 79, 80,
		81, 82, 83, 84, 85, 86, 87, 88,
		89, 90, 97, 98, 99, 100, 101, 102,
		103, 104, 105, 106, 107, 108, 109, 110,
		111, 112, 113, 114, 115, 116, 117, 118,
		119, 120, 121, 122, 48, 49, 50, 51,
		52, 53, 54, 55, 56, 57, 43, 47,
	}
	output := make([]byte, size)
	for x := size; x > 0; x-- {
		output[x-1] = buffer[number&0x3F]
		if number >= 0 {
			number = number >> 6
		} else {
			number = (number >> 6) + (2 << (32 + ^6))
		}
	}
	return output
}

func NewURowID(session PacketSession, packetData []byte, index int) (*Urowid, int) {
	var length int
	length, index = session.GetInt(4, true, true, packetData, index)
	ret := new(Urowid)
	if length > 0 {
		ret.Data, index = session.GetClr(packetData, index)
		return ret, index
	}
	return nil, index
}

func (id *Urowid) getBytes() []byte {
	if id.Data[0] == 1 {
		return id.physicalRawIDToByteArray()
	} else {
		return id.logicalRawIDToByteArray()
	}
}

func (id *Urowid) physicalRawIDToByteArray() []byte {
	// physical
	temp32 := binary.BigEndian.Uint32(id.Data[1:5])
	id.RBA = int64(temp32)
	temp16 := binary.BigEndian.Uint16(id.Data[5:7])
	id.PartitionId = int64(temp16)
	temp32 = binary.BigEndian.Uint32(id.Data[7:11])
	id.BlockNumber = int64(temp32)
	temp16 = binary.BigEndian.Uint16(id.Data[11:13])
	id.SlotNumber = int64(temp16)
	if id.RBA == 0 {
		return []byte(fmt.Sprintf("%08X.%04X.%04X", id.BlockNumber, id.SlotNumber, id.PartitionId))
	} else {
		return id.Rowid.getBytes()
	}
}
func (id *Urowid) logicalRawIDToByteArray() []byte {
	length1 := len(id.Data)
	num1 := length1 / 3
	num2 := length1 % 3
	num3 := num1 * 4
	num4 := 0
	if num2 > 1 {
		num4 = 3
	} else {
		num4 = num2
	}
	length2 := num3 + num4
	var output []byte = nil
	if length2 > 0 {
		KGRD_INDBYTE_CHAR := []byte{65, 42, 45, 40, 41}
		var buffer = []byte{
			65, 66, 67, 68, 69, 70, 71, 72,
			73, 74, 75, 76, 77, 78, 79, 80,
			81, 82, 83, 84, 85, 86, 87, 88,
			89, 90, 97, 98, 99, 100, 101, 102,
			103, 104, 105, 106, 107, 108, 109, 110,
			111, 112, 113, 114, 115, 116, 117, 118,
			119, 120, 121, 122, 48, 49, 50, 51,
			52, 53, 54, 55, 56, 57, 43, 47,
		}
		output = make([]byte, length2)
		srcIndex := 0
		dstIndex := 1
		output[dstIndex] = KGRD_INDBYTE_CHAR[id.Data[srcIndex]-1]
		length1 -= 1
		srcIndex++
		dstIndex++
		for length1 > 0 {
			output[dstIndex] = buffer[id.Data[srcIndex]>>2]
			if length1 == 1 {
				output[dstIndex+1] = buffer[(id.Data[srcIndex]&3)<<4]
				break
			}
			output[dstIndex+1] = buffer[(id.Data[srcIndex]&3)<<4|(id.Data[srcIndex+1]&0xF0)>>4]
			if length1 == 2 {
				output[dstIndex+2] = buffer[(id.Data[srcIndex+1]&0xF)<<2]
				break
			}
			output[dstIndex+2] = buffer[(id.Data[srcIndex+1]&0xF)<<2|(id.Data[srcIndex+2]&0xC0)>>6]
			output[dstIndex+3] = buffer[id.Data[srcIndex+2]&63]
			length1 -= 3
			srcIndex += 3
			dstIndex += 3
		}
	}
	return output
}

func GetStrConv(session PacketSession, charsetID int) converters.IStringConverter {
	switch charsetID {
	case session.SStrConv.GetLangID():
		if session.CStrConv != nil {
			return session.CStrConv
		}
		return session.StrConv
	case session.NStrConv.GetLangID():
		return session.NStrConv
	default:
		temp := converters.NewStringConverter(charsetID)
		if temp == nil {
			return temp
		}
		return temp
	}
}

func decodeObject(session PacketSession, packetData []byte, index int, parent *ParameterInfo) (DecodeObject, int) {
	var objectType byte
	var ctl int
	var itemsLen int
	var bvalueArray [][]byte
	var decodeObjArray []DecodeObject
	var decodePrimevalueArray []OraclePrimeValue
	objectType, index = session.GetByte(packetData, index)
	ctl, index = session.GetInt(4, true, true, packetData, index)
	if ctl == 0xFE {
		_, index = session.GetInt(4, false, true, packetData, index)
	}
	switch objectType {
	case 0x88:
		_, index = session.GetInt(2, true, true, packetData, index)
		itemsLen, index = session.GetInt(2, false, true, packetData, index)
		pars := make([]ParameterInfo, 0, itemsLen)
		var decodeObj DecodeObject
		for x := 0; x < itemsLen; x++ {
			tempPar := parent.Clone()
			tempPar.Direction = parent.Direction
			tempPar.BValue, index = session.GetClr(packetData, index)
			bvalueArray = append(bvalueArray, tempPar.BValue)
			decodeObj, index = decodeObject(session, packetData, index, &tempPar)
			decodeObjArray = append(decodeObjArray, decodeObj)
			pars = append(pars, tempPar)
		}
		parent.OPrimValue = pars
	case 0x84:
		pars := make([]ParameterInfo, 0, len(parent.CusType.Attribs))
		var decodePrimeVal OraclePrimeValue
		for _, attrib := range parent.CusType.Attribs {
			tempPar := attrib
			tempPar.Direction = parent.Direction
			decodePrimeVal, index = tempPar.DecodePrimValue(session, packetData, index, true)
			decodePrimevalueArray = append(decodePrimevalueArray, decodePrimeVal)
			pars = append(pars, tempPar)
		}
		parent.OPrimValue = pars
	}
	oracleDecodeObject := DecodeObject{
		ObjType:               objectType,
		Ctl:                   ctl,
		ItemLen:               itemsLen,
		BValueArray:           bvalueArray,
		DecodeObjArray:        decodeObjArray,
		DecodePrimeValueArray: decodePrimevalueArray,
	}
	return oracleDecodeObject, index
}

func (cursor *RefCursor) Load(session PacketSession, packetData []byte, index int) (*RefCursor, int) {
	var columnCount int
	cursor.Text = ""
	cursor.HasLONG = false
	cursor.HasBLOB = false
	cursor.HasReturnClause = false
	cursor.DisableCompression = false
	cursor.ArrayBindCount = 1
	cursor.ScnForSnapshot = make([]int, 2)
	cursor.StmtType = SELECT
	cursor.Len, index = session.GetByte(packetData, index)
	cursor.MaxRowSize, index = session.GetInt(4, true, true, packetData, index)
	columnCount, index = session.GetInt(4, true, true, packetData, index)
	var paramInfoList []ParameterInfo
	if columnCount > 0 {
		cursor.Columns = make([]ParameterInfo, columnCount)
		_, index = session.GetByte(packetData, index)
		var paramInfo ParameterInfo
		for x := 0; x < len(cursor.Columns); x++ {
			paramInfo, index = cursor.Columns[x].Load(packetData, index, session)
			paramInfoList = append(paramInfoList, paramInfo)
			if cursor.Columns[x].DataType == OCIClobLocator || cursor.Columns[x].DataType == OCIBlobLocator {
				cursor.HasBLOB = true
			}
			if cursor.Columns[x].DataType == LONG || cursor.Columns[x].DataType == LongRaw {
				cursor.HasLONG = true
			}
		}
	}
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
	cursor.CursorID, index = session.GetInt(4, true, true, packetData, index)
	var refcursor RefCursor
	refcursor.Text = cursor.Text
	refcursor.ArrayBindCount = cursor.ArrayBindCount
	refcursor.HasBLOB = cursor.HasBLOB
	refcursor.HasLONG = cursor.HasLONG
	refcursor.HasReturnClause = cursor.HasReturnClause
	refcursor.DisableCompression = cursor.DisableCompression
	refcursor.ScnForSnapshot = cursor.ScnForSnapshot
	refcursor.StmtType = cursor.StmtType
	refcursor.Len = cursor.Len
	refcursor.MaxRowSize = cursor.MaxRowSize
	refcursor.ColumnCount = columnCount
	refcursor.ParamInfoList = paramInfoList
	refcursor.CursorID = cursor.CursorID
	return &refcursor, index
}

func (stmt *Stmt) CalculateColumnValue(col *ParameterInfo, udt bool, session PacketSession, packetData []byte, index int) (CalculateColumnValue, int) {
	var refcursor *RefCursor
	if col.DataType == REFCURSOR {
		var cursor = new(RefCursor)
		cursor.Parent = stmt
		cursor.AutoClose = true
		refcursor, index = cursor.Load(session, packetData, index)
		if stmt.StmtType == PLSQL {
			_, index = session.GetInt(2, true, true, packetData, index)
		}
		col.Value = cursor
		return CalculateColumnValue{
			RefCursor: *refcursor,
		}, index
	}
	var decodeColumnValue DecodeColumnValue
	decodeColumnValue, index = col.DecodeColumnValue(session, udt, packetData, index)
	if refcursor != nil {
		return CalculateColumnValue{
			RefCursor:         *refcursor,
			DecodeColumnValue: decodeColumnValue,
		}, index
	} else {
		return CalculateColumnValue{
			DecodeColumnValue: decodeColumnValue,
		}, index
	}
}

func (par *ParameterInfo) DecodeColumnValue(session PacketSession, udt bool, packetData []byte, index int) (DecodeColumnValue, int) {
	var maxSize int
	var flag byte
	var tempByte byte
	var Bvalue []byte
	if !udt && (par.DataType == OCIBlobLocator || par.DataType == OCIClobLocator) {
		maxSize, index = session.GetInt(4, true, true, packetData, index)
		if maxSize > 0 {
			_, index = session.GetInt(8, true, true, packetData, index)
			_, index = session.GetInt(4, true, true, packetData, index)
			if par.DataType == OCIClobLocator {
				flag, index = session.GetByte(packetData, index)
				par.CharsetID = 0
				if flag == 1 {
					par.CharsetID, index = session.GetInt(2, true, true, packetData, index)
				}
				tempByte, index = session.GetByte(packetData, index)
				par.CharsetForm = int(tempByte)
				if par.CharsetID == 0 {
					if par.CharsetForm == 1 {
						par.CharsetID = session.ServerCharacterSet
					} else {
						par.CharsetID = session.ServernCharacterSet
					}
				}
			}
			par.BValue, index = session.GetClr(packetData, index)
			Bvalue = par.BValue
			if par.DataType == OCIClobLocator {
				strConv := GetStrConv(session, par.CharsetID)
				par.OPrimValue = strConv.Decode(par.BValue)
			} else {
				par.OPrimValue = par.BValue
			}
			_, index = session.GetClr(packetData, index)
		} else {
			par.OPrimValue = nil
		}
	}
	var oraclePrimeValue OraclePrimeValue
	oraclePrimeValue, index = par.DecodePrimValue(session, packetData, index, udt)
	return DecodeColumnValue{MaxSize: maxSize, Flag: flag, TempByte: tempByte, BVlaue: Bvalue, OraclePrimeValue: oraclePrimeValue}, index
}

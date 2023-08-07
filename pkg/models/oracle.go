package models

import (
	"database/sql/driver"

	"bytes"
	"reflect"

	"github.com/sijms/go-ora/v2/converters"
	"github.com/sijms/go-ora/v2/network"
)

type PacketType string

const (
	CONNECT  PacketType = "CONNECT"
	ACCEPT   PacketType = "ACCEPT"
	ACK      PacketType = "ACK"
	REFUSE   PacketType = "REFUSE"
	REDIRECT PacketType = "REDIRECT"
	DATA     PacketType = "DATA"
	NULL     PacketType = "NULL"
	ABORT    PacketType = "ABORT"
	RESEND   PacketType = "RESEND"
	MARKER   PacketType = "MARKER"
	ATTN     PacketType = "ATTN"
	CTRL     PacketType = "CTRL"
	HIGHEST  PacketType = "HIGHEST"
)

var PacketTypeValues = map[uint8]PacketType{
	1:  CONNECT,
	2:  ACCEPT,
	3:  ACK,
	4:  REFUSE,
	5:  REDIRECT,
	6:  DATA,
	7:  NULL,
	9:  ABORT,
	11: RESEND,
	12: MARKER,
	13: ATTN,
	14: CTRL,
	19: HIGHEST,
}

var ReversePacketTypeValues = map[PacketType]uint8{
	CONNECT:  1,
	ACCEPT:   2,
	ACK:      3,
	REFUSE:   4,
	REDIRECT: 5,
	DATA:     6,
	NULL:     7,
	ABORT:    9,
	RESEND:   11,
	MARKER:   12,
	ATTN:     13,
	CTRL:     14,
	HIGHEST:  19,
}

type ControlType string

const (
	TNS_CONTROL_TYPE_RESET_OOB           ControlType = "TNS_CONTROL_TYPE_RESET_OOB"
	TNS_CONTROL_TYPE_INBAND_NOTIFICATION ControlType = "TNS_CONTROL_TYPE_INBAND_NOTIFICATION"
)

var ControlTypeValues = map[int]ControlType{
	9: TNS_CONTROL_TYPE_RESET_OOB,
	8: TNS_CONTROL_TYPE_INBAND_NOTIFICATION,
}

var ReverseControlTypeValues = map[ControlType]int{
	TNS_CONTROL_TYPE_RESET_OOB:           9,
	TNS_CONTROL_TYPE_INBAND_NOTIFICATION: 8,
}

type ControlError string

const (
	TNS_ERR_SESSION_SHUTDOWN ControlError = "TNS_ERR_SESSION_SHUTDOWN"
	TNS_ERR_INBAND_MESSAGE   ControlError = "TNS_ERR_INBAND_MESSAGE"
)

var ControlErrorValues = map[int]ControlError{
	12572: TNS_ERR_SESSION_SHUTDOWN,
	12573: TNS_ERR_INBAND_MESSAGE,
}

var ReverseControlErrorValues = map[ControlError]int{
	TNS_ERR_SESSION_SHUTDOWN: 12572,
	TNS_ERR_INBAND_MESSAGE:   12573,
}

type MarkerType string

const (
	TNS_MARKER_TYPE_BREAK     MarkerType = "TNS_MARKER_TYPE_BREAK"
	TNS_MARKER_TYPE_RESET     MarkerType = "TNS_MARKER_TYPE_RESET"
	TNS_MARKER_TYPE_INTERRUPT MarkerType = "TNS_MARKER_TYPE_INTERRUPT"
)

var MarkerTypeValues = map[int]MarkerType{
	1: TNS_MARKER_TYPE_BREAK,
	2: TNS_MARKER_TYPE_RESET,
	3: TNS_MARKER_TYPE_INTERRUPT,
}

// Reverse map for looking up int from string.
var ReverseMarkerTypeValues = map[MarkerType]int{
	TNS_MARKER_TYPE_BREAK:     1,
	TNS_MARKER_TYPE_RESET:     2,
	TNS_MARKER_TYPE_INTERRUPT: 3,
}

type StmtType string

const (
	DEFAULT StmtType = ""
	SELECT  StmtType = "SELECT"
	DML     StmtType = "DML"
	PLSQL   StmtType = "PLSQL"
	OTHERS  StmtType = "OTHERS"
)

var StmtTypeValues = map[int]StmtType{
	0: DEFAULT,
	1: SELECT,
	2: DML,
	3: PLSQL,
	4: OTHERS,
}

// Reverse map for looking up int from string.
var ReverseStmtTypeValues = map[StmtType]int{
	DEFAULT: 0,
	SELECT:  1,
	DML:     2,
	PLSQL:   3,
	OTHERS:  4,
}

type ParameterDirection string

const (
	Input  ParameterDirection = "Input"
	Output ParameterDirection = "Output"
	InOut  ParameterDirection = "InOut"
)

var DirectionValues = map[int]ParameterDirection{
	1: Input,
	2: Output,
	3: InOut,
}

// Reverse map for looking up int from string.
var ReverseDirectionValues = map[ParameterDirection]int{
	Input:  1,
	Output: 2,
	InOut:  3,
}

type TNSType string

const (
	NCHAR            TNSType = "NCHAR"
	NUMBER           TNSType = "NUMBER"
	BInteger         TNSType = "BInteger"
	FLOAT            TNSType = "FLOAT"
	NullStr          TNSType = "NullStr"
	VarNum           TNSType = "VarNum"
	LONG             TNSType = "LONG"
	VARCHAR          TNSType = "VARCHAR"
	ROWID            TNSType = "ROWID"
	DATE             TNSType = "DATE"
	VarRaw           TNSType = "VarRaw"
	BFloat           TNSType = "BFloat"
	BDouble          TNSType = "BDouble"
	RAW              TNSType = "RAW"
	LongRaw          TNSType = "LongRaw"
	UINT             TNSType = "UINT"
	LongVarChar      TNSType = "LongVarChar"
	LongVarRaw       TNSType = "LongVarRaw"
	CHAR             TNSType = "CHAR"
	CHARZ            TNSType = "CHARZ"
	IBFloat          TNSType = "IBFloat"
	IBDouble         TNSType = "IBDouble"
	REFCURSOR        TNSType = "REFCURSOR"
	OCIXMLType       TNSType = "OCIXMLType"
	XMLType          TNSType = "XMLType"
	OCIRef           TNSType = "OCIRef"
	OCIClobLocator   TNSType = "OCIClobLocator"
	OCIBlobLocator   TNSType = "OCIBlobLocator"
	OCIFileLocator   TNSType = "OCIFileLocator"
	ResultSet        TNSType = "ResultSet"
	OCIString        TNSType = "OCIString"
	OCIDate          TNSType = "OCIDate"
	TimeStampDTY     TNSType = "TimeStampDTY"
	TimeStampTZ_DTY  TNSType = "TimeStampTZ_DTY"
	IntervalYM_DTY   TNSType = "IntervalYM_DTY"
	IntervalDS_DTY   TNSType = "IntervalDS_DTY"
	TimeTZ           TNSType = "TimeTZ"
	TIMESTAMP        TNSType = "TIMESTAMP"
	TIMESTAMPTZ      TNSType = "TIMESTAMPTZ"
	IntervalYM       TNSType = "IntervalYM"
	IntervalDS       TNSType = "IntervalDS"
	UROWID           TNSType = "UROWID"
	TimeStampLTZ_DTY TNSType = "TimeStampLTZ_DTY"
	TimeStampeLTZ    TNSType = "TimeStampeLTZ"
)

var TnsTypeValues = map[int]TNSType{
	1:   NCHAR,
	2:   NUMBER,
	3:   BInteger,
	4:   FLOAT,
	5:   NullStr,
	6:   VarNum,
	8:   LONG,
	9:   VARCHAR,
	11:  ROWID,
	12:  DATE,
	15:  VarRaw,
	21:  BFloat,
	22:  BDouble,
	23:  RAW,
	24:  LongRaw,
	68:  UINT,
	94:  LongVarChar,
	95:  LongVarRaw,
	96:  CHAR,
	97:  CHARZ,
	100: IBFloat,
	101: IBDouble,
	102: REFCURSOR,
	108: OCIXMLType,
	109: XMLType,
	110: OCIRef,
	112: OCIClobLocator,
	113: OCIBlobLocator,
	114: OCIFileLocator,
	116: ResultSet,
	155: OCIString,
	156: OCIDate,
	180: TimeStampDTY,
	181: TimeStampTZ_DTY,
	182: IntervalYM_DTY,
	183: IntervalDS_DTY,
	186: TimeTZ,
	187: TIMESTAMP,
	188: TIMESTAMPTZ,
	189: IntervalYM,
	190: IntervalDS,
	208: UROWID,
	231: TimeStampLTZ_DTY,
	232: TimeStampeLTZ,
}

var ReverseTNSTypeValues = map[TNSType]int{
	NCHAR:            1,
	NUMBER:           2,
	BInteger:         3,
	FLOAT:            4,
	NullStr:          5,
	VarNum:           6,
	LONG:             8,
	VARCHAR:          9,
	ROWID:            11,
	DATE:             12,
	VarRaw:           15,
	BFloat:           21,
	BDouble:          22,
	RAW:              23,
	LongRaw:          24,
	UINT:             68,
	LongVarChar:      94,
	LongVarRaw:       95,
	CHAR:             96,
	CHARZ:            97,
	IBFloat:          100,
	IBDouble:         101,
	REFCURSOR:        102,
	OCIXMLType:       108,
	XMLType:          109,
	OCIRef:           110,
	OCIClobLocator:   112,
	OCIBlobLocator:   113,
	OCIFileLocator:   114,
	ResultSet:        116,
	OCIString:        155,
	OCIDate:          156,
	TimeStampDTY:     180,
	TimeStampTZ_DTY:  181,
	IntervalYM_DTY:   182,
	IntervalDS_DTY:   183,
	TimeTZ:           186,
	TIMESTAMP:        187,
	TIMESTAMPTZ:      188,
	IntervalYM:       189,
	IntervalDS:       190,
	UROWID:           208,
	TimeStampLTZ_DTY: 231,
	TimeStampeLTZ:    232,
}

type DataPacketType string

const (
	DefaultDataPacket                 DataPacketType = "Default"
	OracleProtocolDataMessageType     DataPacketType = "OracleProtocolDataMessageType"
	OracleDataTypeDataMessageType     DataPacketType = "OracleDataTypeDataMessageType"
	OracleFunctionDataMesssageType    DataPacketType = "OracleFunctionDataMesssageType"
	OracleAdvNegoMessagetype          DataPacketType = "OracleAdvNegoMessagetype"
	OracleMsgTypeError                DataPacketType = "OracleMsgTypeError"
	TNS_MSG_TYPE_ROW_HEADER           DataPacketType = "TNS_MSG_TYPE_ROW_HEADER"
	TNS_MSG_TYPE_ROW_DATA             DataPacketType = "TNS_MSG_TYPE_ROW_DATA"
	OracleMsgTypeParameter            DataPacketType = "OracleMsgTypeParameter"
	OracleMsgTypeStatus               DataPacketType = "OracleMsgTypeStatus"
	TNS_MSG_TYPE_IO_VECTOR            DataPacketType = "TNS_MSG_TYPE_IO_VECTOR"
	OracleMsgTypeWarning              DataPacketType = "OracleMsgTypeWarning"
	TNS_MSG_TYPE_DESCRIBE_INFO        DataPacketType = "TNS_MSG_TYPE_DESCRIBE_INFO"
	TNS_MSG_TYPE_PIGGYBACK            DataPacketType = "TNS_MSG_TYPE_PIGGYBACK"
	TNS_MSG_TYPE_FLUSH_OUT_BINDS      DataPacketType = "TNS_MSG_TYPE_FLUSH_OUT_BINDS"
	TNS_MSG_TYPE_BIT_VECTOR           DataPacketType = "TNS_MSG_TYPE_BIT_VECTOR"
	OracleMsgTypeServerSidePiggyback  DataPacketType = "OracleMsgTypeServerSidePiggyback"
	OraclePiggyBackDataMesssageType   DataPacketType = "OraclePiggyBackDataMesssageType"
	OracleConnectionDataMessageType   DataPacketType = "OracleConnectionDataMessageType"
	OracleRedirectDataMessageType     DataPacketType = "OracleRedirectDataMessageType"
	OracleMessageWithDataMessageType  DataPacketType = "OracleMessageWithDataMessageType"
	OracleGetDBVersionDataMessageType DataPacketType = "OracleGetDBVersionDataMessageType"
	TNS_MSG_TYPE_LOB_DATA             DataPacketType = "TNS_MSG_TYPE_LOB_DATA"
	OracleAuthPhaseOneDataMessageType DataPacketType = "OracleAuthPhaseOneDataMessageType"
	OracleAuthPhaseTwoDataMessageType DataPacketType = "OracleAuthPhaseTwoDataMessageType"
	TNS_MSG_TYPE_ONEWAY_FN            DataPacketType = "TNS_MSG_TYPE_ONEWAY_FN"
	TNS_MSG_TYPE_IMPLICIT_RESULTSET   DataPacketType = "TNS_MSG_TYPE_IMPLICIT_RESULTSET"
	TNS_MSG_TYPE_RENEGOTIATE          DataPacketType = "TNS_MSG_TYPE_RENEGOTIATE"
	TNS_MSG_TYPE_COOKIE               DataPacketType = "TNS_MSG_TYPE_COOKIE"
)

var DataPacketTypeValues = map[int]DataPacketType{
	0:          DefaultDataPacket,
	1:          OracleProtocolDataMessageType,
	2:          OracleDataTypeDataMessageType,
	3:          OracleFunctionDataMesssageType,
	0xDEADBEEF: OracleAdvNegoMessagetype,
	4:          OracleMsgTypeError,
	6:          TNS_MSG_TYPE_ROW_HEADER,
	7:          TNS_MSG_TYPE_ROW_DATA,
	8:          OracleMsgTypeParameter,
	9:          OracleMsgTypeStatus,
	11:         TNS_MSG_TYPE_IO_VECTOR,
	15:         OracleMsgTypeWarning,
	16:         TNS_MSG_TYPE_DESCRIBE_INFO,
	19:         TNS_MSG_TYPE_FLUSH_OUT_BINDS,
	21:         TNS_MSG_TYPE_BIT_VECTOR,
	23:         OracleMsgTypeServerSidePiggyback,
	17:         OraclePiggyBackDataMesssageType,
	100:        OracleConnectionDataMessageType,
	101:        OracleRedirectDataMessageType,
	103:        OracleMessageWithDataMessageType,
	104:        OracleGetDBVersionDataMessageType,
	14:         TNS_MSG_TYPE_LOB_DATA,
	110:        OracleAuthPhaseOneDataMessageType,
	111:        OracleAuthPhaseTwoDataMessageType,
	26:         TNS_MSG_TYPE_ONEWAY_FN,
	27:         TNS_MSG_TYPE_IMPLICIT_RESULTSET,
	28:         TNS_MSG_TYPE_RENEGOTIATE,
	30:         TNS_MSG_TYPE_COOKIE,
}

var ReverseDataPacketTypeValues = map[DataPacketType]int{
	DefaultDataPacket:                 0,
	OracleProtocolDataMessageType:     1,
	OracleDataTypeDataMessageType:     2,
	OracleFunctionDataMesssageType:    3,
	OracleAdvNegoMessagetype:          0xDEADBEEF,
	OracleMsgTypeError:                4,
	TNS_MSG_TYPE_ROW_HEADER:           6,
	TNS_MSG_TYPE_ROW_DATA:             7,
	OracleMsgTypeParameter:            8,
	OracleMsgTypeStatus:               9,
	TNS_MSG_TYPE_IO_VECTOR:            11,
	OracleMsgTypeWarning:              15,
	TNS_MSG_TYPE_DESCRIBE_INFO:        16,
	OraclePiggyBackDataMesssageType:   17,
	TNS_MSG_TYPE_FLUSH_OUT_BINDS:      19,
	TNS_MSG_TYPE_BIT_VECTOR:           21,
	OracleMsgTypeServerSidePiggyback:  23,
	OracleConnectionDataMessageType:   100,
	OracleRedirectDataMessageType:     101,
	OracleMessageWithDataMessageType:  103,
	OracleGetDBVersionDataMessageType: 104,
	TNS_MSG_TYPE_LOB_DATA:             14,
	OracleAuthPhaseOneDataMessageType: 110,
	OracleAuthPhaseTwoDataMessageType: 111,
	TNS_MSG_TYPE_ONEWAY_FN:            26,
	TNS_MSG_TYPE_IMPLICIT_RESULTSET:   27,
	TNS_MSG_TYPE_RENEGOTIATE:          28,
	TNS_MSG_TYPE_COOKIE:               30,
}

type FunctionType string

const (
	TNS_FUNC_AUTH_PHASE_ONE      FunctionType = "TNS_FUNC_AUTH_PHASE_ONE"
	TNS_FUNC_AUTH_PHASE_TWO      FunctionType = "TNS_FUNC_AUTH_PHASE_TWO"
	TNS_FUNC_CLOSE_CURSORS       FunctionType = "TNS_FUNC_CLOSE_CURSORS"
	TNS_FUNC_COMMIT              FunctionType = "TNS_FUNC_COMMIT"
	TNS_FUNC_EXECUTE             FunctionType = "TNS_FUNC_EXECUTE"
	TNS_FUNC_FETCH               FunctionType = "TNS_FUNC_FETCH"
	TNS_FUNC_LOB_OP              FunctionType = "TNS_FUNC_LOB_OP"
	TNS_FUNC_LOGOFF              FunctionType = "TNS_FUNC_LOGOFF"
	TNS_FUNC_PING                FunctionType = "TNS_FUNC_PING"
	TNS_FUNC_ROLLBACK            FunctionType = "TNS_FUNC_ROLLBACK"
	TNS_FUNC_SET_END_TO_END_ATTR FunctionType = "TNS_FUNC_SET_END_TO_END_ATTR"
	TNS_FUNC_REEXECUTE           FunctionType = "TNS_FUNC_REEXECUTE"
	TNS_FUNC_REEXECUTE_AND_FETCH FunctionType = "TNS_FUNC_REEXECUTE_AND_FETCH"
	TNS_FUNC_SESSION_GET         FunctionType = "TNS_FUNC_SESSION_GET"
	TNS_FUNC_SESSION_RELEASE     FunctionType = "TNS_FUNC_SESSION_RELEASE"
	TNS_FUNC_SET_SCHEMA          FunctionType = "TNS_FUNC_SET_SCHEMA"
	TNS_FUNC_GET_DB_VERSION      FunctionType = "TNS_FUNC_GET_DB_VERSION"
)

var FunctionTypeValues = map[int]FunctionType{
	118: TNS_FUNC_AUTH_PHASE_ONE,
	115: TNS_FUNC_AUTH_PHASE_TWO,
	105: TNS_FUNC_CLOSE_CURSORS,
	14:  TNS_FUNC_COMMIT,
	94:  TNS_FUNC_EXECUTE,
	5:   TNS_FUNC_FETCH,
	96:  TNS_FUNC_LOB_OP,
	9:   TNS_FUNC_LOGOFF,
	147: TNS_FUNC_PING,
	15:  TNS_FUNC_ROLLBACK,
	135: TNS_FUNC_SET_END_TO_END_ATTR,
	4:   TNS_FUNC_REEXECUTE,
	78:  TNS_FUNC_REEXECUTE_AND_FETCH,
	162: TNS_FUNC_SESSION_GET,
	163: TNS_FUNC_SESSION_RELEASE,
	152: TNS_FUNC_SET_SCHEMA,
	59:  TNS_FUNC_GET_DB_VERSION,
}

var ReverseFunctionTypeValues = map[FunctionType]int{
	TNS_FUNC_AUTH_PHASE_ONE:      118,
	TNS_FUNC_AUTH_PHASE_TWO:      115,
	TNS_FUNC_CLOSE_CURSORS:       105,
	TNS_FUNC_COMMIT:              14,
	TNS_FUNC_EXECUTE:             94,
	TNS_FUNC_FETCH:               5,
	TNS_FUNC_LOB_OP:              96,
	TNS_FUNC_LOGOFF:              9,
	TNS_FUNC_PING:                147,
	TNS_FUNC_ROLLBACK:            15,
	TNS_FUNC_SET_END_TO_END_ATTR: 135,
	TNS_FUNC_REEXECUTE:           4,
	TNS_FUNC_REEXECUTE_AND_FETCH: 78,
	TNS_FUNC_SESSION_GET:         162,
	TNS_FUNC_SESSION_RELEASE:     163,
	TNS_FUNC_SET_SCHEMA:          152,
	TNS_FUNC_GET_DB_VERSION:      59,
}

type AuthMode string

const (
	TNS_AUTH_MODE_LOGON           AuthMode = "TNS_AUTH_MODE_LOGON"
	TNS_AUTH_MODE_CHANGE_PASSWORD AuthMode = "TNS_AUTH_MODE_CHANGE_PASSWORD"
	TNS_AUTH_MODE_SYSDBA          AuthMode = "TNS_AUTH_MODE_SYSDBA"
	TNS_AUTH_MODE_SYSOPER         AuthMode = "TNS_AUTH_MODE_SYSOPER"
	TNS_AUTH_MODE_PRELIM          AuthMode = "TNS_AUTH_MODE_PRELIM"
	TNS_AUTH_MODE_WITH_PASSWORD   AuthMode = "TNS_AUTH_MODE_WITH_PASSWORD"
	TNS_AUTH_MODE_SYSASM          AuthMode = "TNS_AUTH_MODE_SYSASM"
	TNS_AUTH_MODE_SYSBKP          AuthMode = "TNS_AUTH_MODE_SYSBKP"
	TNS_AUTH_MODE_SYSDGD          AuthMode = "TNS_AUTH_MODE_SYSDGD"
	TNS_AUTH_MODE_SYSKMT          AuthMode = "TNS_AUTH_MODE_SYSKMT"
	TNS_AUTH_MODE_SYSRAC          AuthMode = "TNS_AUTH_MODE_SYSRAC"
	TNS_AUTH_MODE_IAM_TOKEN       AuthMode = "TNS_AUTH_MODE_IAM_TOKEN"
)

var AuthModeValues = map[int]AuthMode{
	0x00000001: TNS_AUTH_MODE_LOGON,
	0x00000002: TNS_AUTH_MODE_CHANGE_PASSWORD,
	0x00000020: TNS_AUTH_MODE_SYSDBA,
	0x00000040: TNS_AUTH_MODE_SYSOPER,
	0x00000080: TNS_AUTH_MODE_PRELIM,
	0x00000100: TNS_AUTH_MODE_WITH_PASSWORD,
	0x00400000: TNS_AUTH_MODE_SYSASM,
	0x01000000: TNS_AUTH_MODE_SYSBKP,
	0x02000000: TNS_AUTH_MODE_SYSDGD,
	0x04000000: TNS_AUTH_MODE_SYSKMT,
	0x08000000: TNS_AUTH_MODE_SYSRAC,
	0x20000000: TNS_AUTH_MODE_IAM_TOKEN,
}

var ReverseAuthModeValues = map[AuthMode]int{
	TNS_AUTH_MODE_LOGON:           0x00000001,
	TNS_AUTH_MODE_CHANGE_PASSWORD: 0x00000002,
	TNS_AUTH_MODE_SYSDBA:          0x00000020,
	TNS_AUTH_MODE_SYSOPER:         0x00000040,
	TNS_AUTH_MODE_PRELIM:          0x00000080,
	TNS_AUTH_MODE_WITH_PASSWORD:   0x00000100,
	TNS_AUTH_MODE_SYSASM:          0x00400000,
	TNS_AUTH_MODE_SYSBKP:          0x01000000,
	TNS_AUTH_MODE_SYSDGD:          0x02000000,
	TNS_AUTH_MODE_SYSKMT:          0x04000000,
	TNS_AUTH_MODE_SYSRAC:          0x08000000,
	TNS_AUTH_MODE_IAM_TOKEN:       0x20000000,
}

type ServiceType string

const (
	AUTH_SERVICE           ServiceType = "AUTH_SERVICE"
	ENCRYPT_SERVICE        ServiceType = "ENCRYPT_SERVICE"
	DATA_INTEGRITY_SERVICE ServiceType = "DATA_INTEGRITY_SERVICE"
	SUPERVISOR_SERVICE     ServiceType = "SUPERVISOR_SERVICE"
)

var ServiceTypeValues = map[int]ServiceType{
	1: AUTH_SERVICE,
	2: ENCRYPT_SERVICE,
	3: DATA_INTEGRITY_SERVICE,
	4: SUPERVISOR_SERVICE,
}

var ReverseServiceTypeValues = map[ServiceType]int{
	AUTH_SERVICE:           1,
	ENCRYPT_SERVICE:        2,
	DATA_INTEGRITY_SERVICE: 3,
	SUPERVISOR_SERVICE:     4,
}

type PiggyBackType string

const (
	TNS_SERVER_PIGGYBACK_QUERY_CACHE_INVALIDATION PiggyBackType = "TNS_SERVER_PIGGYBACK_QUERY_CACHE_INVALIDATION"
	TNS_SERVER_PIGGYBACK_OS_PID_MTS               PiggyBackType = "TNS_SERVER_PIGGYBACK_OS_PID_MTS"
	TNS_SERVER_PIGGYBACK_TRACE_EVENT              PiggyBackType = "TNS_SERVER_PIGGYBACK_TRACE_EVENT"
	TNS_SERVER_PIGGYBACK_SESS_RET                 PiggyBackType = "TNS_SERVER_PIGGYBACK_SESS_RET"
	TNS_SERVER_PIGGYBACK_SYNC                     PiggyBackType = "TNS_SERVER_PIGGYBACK_SYNC"
	TNS_SERVER_PIGGYBACK_LTXID                    PiggyBackType = "TNS_SERVER_PIGGYBACK_LTXID"
	TNS_SERVER_PIGGYBACK_AC_REPLAY_CONTEXT        PiggyBackType = "TNS_SERVER_PIGGYBACK_AC_REPLAY_CONTEXT"
	TNS_SERVER_PIGGYBACK_EXT_SYNC                 PiggyBackType = "TNS_SERVER_PIGGYBACK_EXT_SYNC"
)

var PiggyBackTypeValues = map[int]PiggyBackType{
	1: TNS_SERVER_PIGGYBACK_QUERY_CACHE_INVALIDATION,
	2: TNS_SERVER_PIGGYBACK_OS_PID_MTS,
	3: TNS_SERVER_PIGGYBACK_TRACE_EVENT,
	4: TNS_SERVER_PIGGYBACK_SESS_RET,
	5: TNS_SERVER_PIGGYBACK_SYNC,
	7: TNS_SERVER_PIGGYBACK_LTXID,
	8: TNS_SERVER_PIGGYBACK_AC_REPLAY_CONTEXT,
	9: TNS_SERVER_PIGGYBACK_EXT_SYNC,
}

var ReversePiggyBackTypeValues = map[PiggyBackType]int{
	TNS_SERVER_PIGGYBACK_QUERY_CACHE_INVALIDATION: 1,
	TNS_SERVER_PIGGYBACK_OS_PID_MTS:               2,
	TNS_SERVER_PIGGYBACK_TRACE_EVENT:              3,
	TNS_SERVER_PIGGYBACK_SESS_RET:                 4,
	TNS_SERVER_PIGGYBACK_SYNC:                     5,
	TNS_SERVER_PIGGYBACK_LTXID:                    7,
	TNS_SERVER_PIGGYBACK_AC_REPLAY_CONTEXT:        8,
	TNS_SERVER_PIGGYBACK_EXT_SYNC:                 9,
}

type DataType string

const (
	DEFAULT_DATA_TYPE            DataType = "Default"
	TNS_DATA_TYPE_VARCHAR        DataType = "TNS_DATA_TYPE_VARCHAR"
	TNS_DATA_TYPE_NUMBER         DataType = "TNS_DATA_TYPE_NUMBER"
	TNS_DATA_TYPE_BINARY_INTEGER DataType = "TNS_DATA_TYPE_BINARY_INTEGER"
	TNS_DATA_TYPE_FLOAT          DataType = "TNS_DATA_TYPE_FLOAT"
	TNS_DATA_TYPE_STR            DataType = "TNS_DATA_TYPE_STR"
	TNS_DATA_TYPE_VNU            DataType = "TNS_DATA_TYPE_VNU"
	TNS_DATA_TYPE_PDN            DataType = "TNS_DATA_TYPE_PDN"
	TNS_DATA_TYPE_LONG           DataType = "TNS_DATA_TYPE_LONG"
	TNS_DATA_TYPE_VCS            DataType = "TNS_DATA_TYPE_VCS"
	TNS_DATA_TYPE_TID            DataType = "TNS_DATA_TYPE_TID"
	TNS_DATA_TYPE_ROWID          DataType = "TNS_DATA_TYPE_ROWID"
	TNS_DATA_TYPE_DATE           DataType = "TNS_DATA_TYPE_DATE"
	TNS_DATA_TYPE_VBI            DataType = "TNS_DATA_TYPE_VBI"
	TNS_DATA_TYPE_RAW            DataType = "TNS_DATA_TYPE_RAW"
	TNS_DATA_TYPE_LONG_RAW       DataType = "TNS_DATA_TYPE_LONG_RAW"
	TNS_DATA_TYPE_UB2            DataType = "TNS_DATA_TYPE_UB2"
	TNS_DATA_TYPE_UB4            DataType = "TNS_DATA_TYPE_UB4"
	TNS_DATA_TYPE_SB1            DataType = "TNS_DATA_TYPE_SB1"
	TNS_DATA_TYPE_SB2            DataType = "TNS_DATA_TYPE_SB2"
	TNS_DATA_TYPE_SB4            DataType = "TNS_DATA_TYPE_SB4"
	TNS_DATA_TYPE_SWORD          DataType = "TNS_DATA_TYPE_SWORD"
	TNS_DATA_TYPE_UWORD          DataType = "TNS_DATA_TYPE_UWORD"
	TNS_DATA_TYPE_PTRB           DataType = "TNS_DATA_TYPE_PTRB"
	TNS_DATA_TYPE_PTRW           DataType = "TNS_DATA_TYPE_PTRW"
	TNS_DATA_TYPE_OER8           DataType = "TNS_DATA_TYPE_OER8"
	TNS_DATA_TYPE_FUN            DataType = "TNS_DATA_TYPE_FUN"
	TNS_DATA_TYPE_AUA            DataType = "TNS_DATA_TYPE_AUA"
	TNS_DATA_TYPE_RXH7           DataType = "TNS_DATA_TYPE_RXH7"
	TNS_DATA_TYPE_NA6            DataType = "TNS_DATA_TYPE_NA6"
	TNS_DATA_TYPE_OAC            DataType = "TNS_DATA_TYPE_OAC"
	TNS_DATA_TYPE_AMS            DataType = "TNS_DATA_TYPE_AMS"
	TNS_DATA_TYPE_BRN            DataType = "TNS_DATA_TYPE_BRN"
	TNS_DATA_TYPE_BRP            DataType = "TNS_DATA_TYPE_BRP"
	TNS_DATA_TYPE_BRV            DataType = "TNS_DATA_TYPE_BRV"
	TNS_DATA_TYPE_KVA            DataType = "TNS_DATA_TYPE_KVA"
	TNS_DATA_TYPE_CLS            DataType = "TNS_DATA_TYPE_CLS"
	TNS_DATA_TYPE_CUI            DataType = "TNS_DATA_TYPE_CUI"
	TNS_DATA_TYPE_DFN            DataType = "TNS_DATA_TYPE_DFN"
	TNS_DATA_TYPE_DQR            DataType = "TNS_DATA_TYPE_DQR"
	TNS_DATA_TYPE_DSC            DataType = "TNS_DATA_TYPE_DSC"
	TNS_DATA_TYPE_EXE            DataType = "TNS_DATA_TYPE_EXE"
	TNS_DATA_TYPE_FCH            DataType = "TNS_DATA_TYPE_FCH"
	TNS_DATA_TYPE_GBV            DataType = "TNS_DATA_TYPE_GBV"
	TNS_DATA_TYPE_GEM            DataType = "TNS_DATA_TYPE_GEM"
	TNS_DATA_TYPE_GIV            DataType = "TNS_DATA_TYPE_GIV"
	TNS_DATA_TYPE_OKG            DataType = "TNS_DATA_TYPE_OKG"
	TNS_DATA_TYPE_HMI            DataType = "TNS_DATA_TYPE_HMI"
	TNS_DATA_TYPE_INO            DataType = "TNS_DATA_TYPE_INO"
	TNS_DATA_TYPE_LNF            DataType = "TNS_DATA_TYPE_LNF"
	TNS_DATA_TYPE_ONT            DataType = "TNS_DATA_TYPE_ONT"
	TNS_DATA_TYPE_OPE            DataType = "TNS_DATA_TYPE_OPE"
	TNS_DATA_TYPE_OSQ            DataType = "TNS_DATA_TYPE_OSQ"
	TNS_DATA_TYPE_SFE            DataType = "TNS_DATA_TYPE_SFE"
	TNS_DATA_TYPE_SPF            DataType = "TNS_DATA_TYPE_SPF"
	TNS_DATA_TYPE_VSN            DataType = "TNS_DATA_TYPE_VSN"
	TNS_DATA_TYPE_UD7            DataType = "TNS_DATA_TYPE_UD7"
	TNS_DATA_TYPE_DSA            DataType = "TNS_DATA_TYPE_DSA"
	TNS_DATA_TYPE_UIN            DataType = "TNS_DATA_TYPE_UIN"
	TNS_DATA_TYPE_PIN            DataType = "TNS_DATA_TYPE_PIN"
	TNS_DATA_TYPE_PFN            DataType = "TNS_DATA_TYPE_PFN"
	TNS_DATA_TYPE_PPT            DataType = "TNS_DATA_TYPE_PPT"
	TNS_DATA_TYPE_STO            DataType = "TNS_DATA_TYPE_STO"
	TNS_DATA_TYPE_ARC            DataType = "TNS_DATA_TYPE_ARC"
	TNS_DATA_TYPE_MRS            DataType = "TNS_DATA_TYPE_MRS"
	TNS_DATA_TYPE_MRT            DataType = "TNS_DATA_TYPE_MRT"
	TNS_DATA_TYPE_MRG            DataType = "TNS_DATA_TYPE_MRG"
	TNS_DATA_TYPE_MRR            DataType = "TNS_DATA_TYPE_MRR"
	TNS_DATA_TYPE_MRC            DataType = "TNS_DATA_TYPE_MRC"
	TNS_DATA_TYPE_VER            DataType = "TNS_DATA_TYPE_VER"
	TNS_DATA_TYPE_LON2           DataType = "TNS_DATA_TYPE_LON2"
	TNS_DATA_TYPE_INO2           DataType = "TNS_DATA_TYPE_INO2"
	TNS_DATA_TYPE_ALL            DataType = "TNS_DATA_TYPE_ALL"
	TNS_DATA_TYPE_UDB            DataType = "TNS_DATA_TYPE_UDB"
	TNS_DATA_TYPE_AQI            DataType = "TNS_DATA_TYPE_AQI"
	TNS_DATA_TYPE_ULB            DataType = "TNS_DATA_TYPE_ULB"
	TNS_DATA_TYPE_ULD            DataType = "TNS_DATA_TYPE_ULD"
	TNS_DATA_TYPE_SLS            DataType = "TNS_DATA_TYPE_SLS"
	TNS_DATA_TYPE_SID            DataType = "TNS_DATA_TYPE_SID"
	TNS_DATA_TYPE_NA7            DataType = "TNS_DATA_TYPE_NA7"
	TNS_DATA_TYPE_LVC            DataType = "TNS_DATA_TYPE_LVC"
	TNS_DATA_TYPE_LVB            DataType = "TNS_DATA_TYPE_LVB"
	TNS_DATA_TYPE_CHAR           DataType = "TNS_DATA_TYPE_CHAR"
	TNS_DATA_TYPE_AVC            DataType = "TNS_DATA_TYPE_AVC"
	TNS_DATA_TYPE_AL7            DataType = "TNS_DATA_TYPE_AL7"
	TNS_DATA_TYPE_K2RPC          DataType = "TNS_DATA_TYPE_K2RPC"
	TNS_DATA_TYPE_BINARY_FLOAT   DataType = "TNS_DATA_TYPE_BINARY_FLOAT"
	TNS_DATA_TYPE_BINARY_DOUBLE  DataType = "TNS_DATA_TYPE_BINARY_DOUBLE"
	TNS_DATA_TYPE_CURSOR         DataType = "TNS_DATA_TYPE_CURSOR"
	TNS_DATA_TYPE_RDD            DataType = "TNS_DATA_TYPE_RDD"
	TNS_DATA_TYPE_XDP            DataType = "TNS_DATA_TYPE_XDP"
	TNS_DATA_TYPE_OSL            DataType = "TNS_DATA_TYPE_OSL"
	TNS_DATA_TYPE_OKO8           DataType = "TNS_DATA_TYPE_OKO8"
	TNS_DATA_TYPE_EXT_NAMED      DataType = "TNS_DATA_TYPE_EXT_NAMED"
	TNS_DATA_TYPE_INT_NAMED      DataType = "TNS_DATA_TYPE_INT_NAMED"
	TNS_DATA_TYPE_EXT_REF        DataType = "TNS_DATA_TYPE_EXT_REF"
	TNS_DATA_TYPE_INT_REF        DataType = "TNS_DATA_TYPE_INT_REF"
	TNS_DATA_TYPE_CLOB           DataType = "TNS_DATA_TYPE_CLOB"
	TNS_DATA_TYPE_BLOB           DataType = "TNS_DATA_TYPE_BLOB"
	TNS_DATA_TYPE_BFILE          DataType = "TNS_DATA_TYPE_BFILE"
	TNS_DATA_TYPE_CFILE          DataType = "TNS_DATA_TYPE_CFILE"
	TNS_DATA_TYPE_RSET           DataType = "TNS_DATA_TYPE_RSET"
	TNS_DATA_TYPE_CWD            DataType = "TNS_DATA_TYPE_CWD"
	TNS_DATA_TYPE_JSON           DataType = "TNS_DATA_TYPE_JSON"
	TNS_DATA_TYPE_NEW_OAC        DataType = "TNS_DATA_TYPE_NEW_OAC"
	TNS_DATA_TYPE_UD12           DataType = "TNS_DATA_TYPE_UD12"
	TNS_DATA_TYPE_AL8            DataType = "TNS_DATA_TYPE_AL8"
	TNS_DATA_TYPE_LFOP           DataType = "TNS_DATA_TYPE_LFOP"
	TNS_DATA_TYPE_FCRT           DataType = "TNS_DATA_TYPE_FCRT"
	TNS_DATA_TYPE_DNY            DataType = "TNS_DATA_TYPE_DNY"
	TNS_DATA_TYPE_OPR            DataType = "TNS_DATA_TYPE_OPR"
	TNS_DATA_TYPE_PLS            DataType = "TNS_DATA_TYPE_PLS"
	TNS_DATA_TYPE_XID            DataType = "TNS_DATA_TYPE_XID"
	TNS_DATA_TYPE_TXN            DataType = "TNS_DATA_TYPE_TXN"
	TNS_DATA_TYPE_DCB            DataType = "TNS_DATA_TYPE_DCB"
	TNS_DATA_TYPE_CCA            DataType = "TNS_DATA_TYPE_CCA"
	TNS_DATA_TYPE_WRN            DataType = "TNS_DATA_TYPE_WRN"
	TNS_DATA_TYPE_TLH            DataType = "TNS_DATA_TYPE_TLH"
	TNS_DATA_TYPE_TOH            DataType = "TNS_DATA_TYPE_TOH"
	TNS_DATA_TYPE_FOI            DataType = "TNS_DATA_TYPE_FOI"
	TNS_DATA_TYPE_SID2           DataType = "TNS_DATA_TYPE_SID2"
	TNS_DATA_TYPE_TCH            DataType = "TNS_DATA_TYPE_TCH"
	TNS_DATA_TYPE_PII            DataType = "TNS_DATA_TYPE_PII"
	TNS_DATA_TYPE_PFI            DataType = "TNS_DATA_TYPE_PFI"
	TNS_DATA_TYPE_PPU            DataType = "TNS_DATA_TYPE_PPU"
	TNS_DATA_TYPE_PTE            DataType = "TNS_DATA_TYPE_PTE"
	TNS_DATA_TYPE_CLV            DataType = "TNS_DATA_TYPE_CLV"
	TNS_DATA_TYPE_RXH8           DataType = "TNS_DATA_TYPE_RXH8"
	TNS_DATA_TYPE_N12            DataType = "TNS_DATA_TYPE_N12"
	TNS_DATA_TYPE_AUTH           DataType = "TNS_DATA_TYPE_AUTH"
	TNS_DATA_TYPE_KVAL           DataType = "TNS_DATA_TYPE_KVAL"
	TNS_DATA_TYPE_DTR            DataType = "TNS_DATA_TYPE_DTR"
	TNS_DATA_TYPE_DUN            DataType = "TNS_DATA_TYPE_DUN"
	TNS_DATA_TYPE_DOP            DataType = "TNS_DATA_TYPE_DOP"
	TNS_DATA_TYPE_VST            DataType = "TNS_DATA_TYPE_VST"
	TNS_DATA_TYPE_ODT            DataType = "TNS_DATA_TYPE_ODT"
	TNS_DATA_TYPE_FGI            DataType = "TNS_DATA_TYPE_FGI"
	TNS_DATA_TYPE_DSY            DataType = "TNS_DATA_TYPE_DSY"
	TNS_DATA_TYPE_DSYR8          DataType = "TNS_DATA_TYPE_DSYR8"
	TNS_DATA_TYPE_DSYH8          DataType = "TNS_DATA_TYPE_DSYH8"
	TNS_DATA_TYPE_DSYL           DataType = "TNS_DATA_TYPE_DSYL"
	TNS_DATA_TYPE_DSYT8          DataType = "TNS_DATA_TYPE_DSYT8"
	TNS_DATA_TYPE_DSYV8          DataType = "TNS_DATA_TYPE_DSYV8"
	TNS_DATA_TYPE_DSYP           DataType = "TNS_DATA_TYPE_DSYP"
	TNS_DATA_TYPE_DSYF           DataType = "TNS_DATA_TYPE_DSYF"
	TNS_DATA_TYPE_DSYK           DataType = "TNS_DATA_TYPE_DSYK"
	TNS_DATA_TYPE_DSYY           DataType = "TNS_DATA_TYPE_DSYY"
	TNS_DATA_TYPE_DSYQ           DataType = "TNS_DATA_TYPE_DSYQ"
	TNS_DATA_TYPE_DSYC           DataType = "TNS_DATA_TYPE_DSYC"
	TNS_DATA_TYPE_DSYA           DataType = "TNS_DATA_TYPE_DSYA"
	TNS_DATA_TYPE_OT8            DataType = "TNS_DATA_TYPE_OT8"
	TNS_DATA_TYPE_DOL            DataType = "TNS_DATA_TYPE_DOL"
	TNS_DATA_TYPE_DSYTY          DataType = "TNS_DATA_TYPE_DSYTY"
	TNS_DATA_TYPE_AQE            DataType = "TNS_DATA_TYPE_AQE"
	TNS_DATA_TYPE_KV             DataType = "TNS_DATA_TYPE_KV"
	TNS_DATA_TYPE_AQD            DataType = "TNS_DATA_TYPE_AQD"
	TNS_DATA_TYPE_AQ8            DataType = "TNS_DATA_TYPE_AQ8"
	TNS_DATA_TYPE_TIME           DataType = "TNS_DATA_TYPE_TIME"
	TNS_DATA_TYPE_TIME_TZ        DataType = "TNS_DATA_TYPE_TIME_TZ"
	TNS_DATA_TYPE_TIMESTAMP      DataType = "TNS_DATA_TYPE_TIMESTAMP"
	TNS_DATA_TYPE_TIMESTAMP_TZ   DataType = "TNS_DATA_TYPE_TIMESTAMP_TZ"
	TNS_DATA_TYPE_INTERVAL_YM    DataType = "TNS_DATA_TYPE_INTERVAL_YM"
	TNS_DATA_TYPE_INTERVAL_DS    DataType = "TNS_DATA_TYPE_INTERVAL_DS"
	TNS_DATA_TYPE_EDATE          DataType = "TNS_DATA_TYPE_EDATE"
	TNS_DATA_TYPE_ETIME          DataType = "TNS_DATA_TYPE_ETIME"
	TNS_DATA_TYPE_ETTZ           DataType = "TNS_DATA_TYPE_ETTZ"
	TNS_DATA_TYPE_ESTAMP         DataType = "TNS_DATA_TYPE_ESTAMP"
	TNS_DATA_TYPE_ESTZ           DataType = "TNS_DATA_TYPE_ESTZ"
	TNS_DATA_TYPE_EIYM           DataType = "TNS_DATA_TYPE_EIYM"
	TNS_DATA_TYPE_EIDS           DataType = "TNS_DATA_TYPE_EIDS"
	TNS_DATA_TYPE_RFS            DataType = "TNS_DATA_TYPE_RFS"
	TNS_DATA_TYPE_RXH10          DataType = "TNS_DATA_TYPE_RXH10"
	TNS_DATA_TYPE_DCLOB          DataType = "TNS_DATA_TYPE_DCLOB"
	TNS_DATA_TYPE_DBLOB          DataType = "TNS_DATA_TYPE_DBLOB"
	TNS_DATA_TYPE_DBFILE         DataType = "TNS_DATA_TYPE_DBFILE"
	TNS_DATA_TYPE_DJSON          DataType = "TNS_DATA_TYPE_DJSON"
	TNS_DATA_TYPE_KPN            DataType = "TNS_DATA_TYPE_KPN"
	TNS_DATA_TYPE_KPDNR          DataType = "TNS_DATA_TYPE_KPDNR"
	TNS_DATA_TYPE_DSYD           DataType = "TNS_DATA_TYPE_DSYD"
	TNS_DATA_TYPE_DSYS           DataType = "TNS_DATA_TYPE_DSYS"
	TNS_DATA_TYPE_DSYR           DataType = "TNS_DATA_TYPE_DSYR"
	TNS_DATA_TYPE_DSYH           DataType = "TNS_DATA_TYPE_DSYH"
	TNS_DATA_TYPE_DSYT           DataType = "TNS_DATA_TYPE_DSYT"
	TNS_DATA_TYPE_DSYV           DataType = "TNS_DATA_TYPE_DSYV"
	TNS_DATA_TYPE_AQM            DataType = "TNS_DATA_TYPE_AQM"
	TNS_DATA_TYPE_OER11          DataType = "TNS_DATA_TYPE_OER11"
	TNS_DATA_TYPE_UROWID         DataType = "TNS_DATA_TYPE_UROWID"
	TNS_DATA_TYPE_AQL            DataType = "TNS_DATA_TYPE_AQL"
	TNS_DATA_TYPE_OTC            DataType = "TNS_DATA_TYPE_OTC"
	TNS_DATA_TYPE_KFNO           DataType = "TNS_DATA_TYPE_KFNO"
	TNS_DATA_TYPE_KFNP           DataType = "TNS_DATA_TYPE_KFNP"
	TNS_DATA_TYPE_KGT8           DataType = "TNS_DATA_TYPE_KGT8"
	TNS_DATA_TYPE_RASB4          DataType = "TNS_DATA_TYPE_RASB4"
	TNS_DATA_TYPE_RAUB2          DataType = "TNS_DATA_TYPE_RAUB2"
	TNS_DATA_TYPE_RAUB1          DataType = "TNS_DATA_TYPE_RAUB1"
	TNS_DATA_TYPE_RATXT          DataType = "TNS_DATA_TYPE_RATXT"
	TNS_DATA_TYPE_RSSB4          DataType = "TNS_DATA_TYPE_RSSB4"
	TNS_DATA_TYPE_RSUB2          DataType = "TNS_DATA_TYPE_RSUB2"
	TNS_DATA_TYPE_RSUB1          DataType = "TNS_DATA_TYPE_RSUB1"
	TNS_DATA_TYPE_RSTXT          DataType = "TNS_DATA_TYPE_RSTXT"
	TNS_DATA_TYPE_RIDL           DataType = "TNS_DATA_TYPE_RIDL"
	TNS_DATA_TYPE_GLRDD          DataType = "TNS_DATA_TYPE_GLRDD"
	TNS_DATA_TYPE_GLRDG          DataType = "TNS_DATA_TYPE_GLRDG"
	TNS_DATA_TYPE_GLRDC          DataType = "TNS_DATA_TYPE_GLRDC"
	TNS_DATA_TYPE_OKO            DataType = "TNS_DATA_TYPE_OKO"
	TNS_DATA_TYPE_DPP            DataType = "TNS_DATA_TYPE_DPP"
	TNS_DATA_TYPE_DPLS           DataType = "TNS_DATA_TYPE_DPLS"
	TNS_DATA_TYPE_DPMOP          DataType = "TNS_DATA_TYPE_DPMOP"
	TNS_DATA_TYPE_TIMESTAMP_LTZ  DataType = "TNS_DATA_TYPE_TIMESTAMP_LTZ"
	TNS_DATA_TYPE_ESITZ          DataType = "TNS_DATA_TYPE_ESITZ"
	TNS_DATA_TYPE_UB8            DataType = "TNS_DATA_TYPE_UB8"
	TNS_DATA_TYPE_STAT           DataType = "TNS_DATA_TYPE_STAT"
	TNS_DATA_TYPE_RFX            DataType = "TNS_DATA_TYPE_RFX"
	TNS_DATA_TYPE_FAL            DataType = "TNS_DATA_TYPE_FAL"
	TNS_DATA_TYPE_CKV            DataType = "TNS_DATA_TYPE_CKV"
	TNS_DATA_TYPE_DRCX           DataType = "TNS_DATA_TYPE_DRCX"
	TNS_DATA_TYPE_KGH            DataType = "TNS_DATA_TYPE_KGH"
	TNS_DATA_TYPE_AQO            DataType = "TNS_DATA_TYPE_AQO"
	TNS_DATA_TYPE_PNTY           DataType = "TNS_DATA_TYPE_PNTY"
	TNS_DATA_TYPE_OKGT           DataType = "TNS_DATA_TYPE_OKGT"
	TNS_DATA_TYPE_KPFC           DataType = "TNS_DATA_TYPE_KPFC"
	TNS_DATA_TYPE_FE2            DataType = "TNS_DATA_TYPE_FE2"
	TNS_DATA_TYPE_SPFP           DataType = "TNS_DATA_TYPE_SPFP"
	TNS_DATA_TYPE_DPULS          DataType = "TNS_DATA_TYPE_DPULS"
	TNS_DATA_TYPE_BOOLEAN        DataType = "TNS_DATA_TYPE_BOOLEAN"
	TNS_DATA_TYPE_AQA            DataType = "TNS_DATA_TYPE_AQA"
	TNS_DATA_TYPE_KPBF           DataType = "TNS_DATA_TYPE_KPBF"
	TNS_DATA_TYPE_TSM            DataType = "TNS_DATA_TYPE_TSM"
	TNS_DATA_TYPE_MSS            DataType = "TNS_DATA_TYPE_MSS"
	TNS_DATA_TYPE_KPC            DataType = "TNS_DATA_TYPE_KPC"
	TNS_DATA_TYPE_CRS            DataType = "TNS_DATA_TYPE_CRS"
	TNS_DATA_TYPE_KKS            DataType = "TNS_DATA_TYPE_KKS"
	TNS_DATA_TYPE_KSP            DataType = "TNS_DATA_TYPE_KSP"
	TNS_DATA_TYPE_KSPTOP         DataType = "TNS_DATA_TYPE_KSPTOP"
	TNS_DATA_TYPE_KSPVAL         DataType = "TNS_DATA_TYPE_KSPVAL"
	TNS_DATA_TYPE_PSS            DataType = "TNS_DATA_TYPE_PSS"
	TNS_DATA_TYPE_NLS            DataType = "TNS_DATA_TYPE_NLS"
	TNS_DATA_TYPE_ALS            DataType = "TNS_DATA_TYPE_ALS"
	TNS_DATA_TYPE_KSDEVTVAL      DataType = "TNS_DATA_TYPE_KSDEVTVAL"
	TNS_DATA_TYPE_KSDEVTTOP      DataType = "TNS_DATA_TYPE_KSDEVTTOP"
	TNS_DATA_TYPE_KPSPP          DataType = "TNS_DATA_TYPE_KPSPP"
	TNS_DATA_TYPE_KOL            DataType = "TNS_DATA_TYPE_KOL"
	TNS_DATA_TYPE_LST            DataType = "TNS_DATA_TYPE_LST"
	TNS_DATA_TYPE_ACX            DataType = "TNS_DATA_TYPE_ACX"
	TNS_DATA_TYPE_SCS            DataType = "TNS_DATA_TYPE_SCS"
	TNS_DATA_TYPE_RXH            DataType = "TNS_DATA_TYPE_RXH"
	TNS_DATA_TYPE_KPDNS          DataType = "TNS_DATA_TYPE_KPDNS"
	TNS_DATA_TYPE_KPDCN          DataType = "TNS_DATA_TYPE_KPDCN"
	TNS_DATA_TYPE_KPNNS          DataType = "TNS_DATA_TYPE_KPNNS"
	TNS_DATA_TYPE_KPNCN          DataType = "TNS_DATA_TYPE_KPNCN"
	TNS_DATA_TYPE_KPS            DataType = "TNS_DATA_TYPE_KPS"
	TNS_DATA_TYPE_APINF          DataType = "TNS_DATA_TYPE_APINF"
	TNS_DATA_TYPE_TEN            DataType = "TNS_DATA_TYPE_TEN"
	TNS_DATA_TYPE_XSSCS          DataType = "TNS_DATA_TYPE_XSSCS"
	TNS_DATA_TYPE_XSSSO          DataType = "TNS_DATA_TYPE_XSSSO"
	TNS_DATA_TYPE_XSSAO          DataType = "TNS_DATA_TYPE_XSSAO"
	TNS_DATA_TYPE_KSRPC          DataType = "TNS_DATA_TYPE_KSRPC"
	TNS_DATA_TYPE_KVL            DataType = "TNS_DATA_TYPE_KVL"
	TNS_DATA_TYPE_SESSGET        DataType = "TNS_DATA_TYPE_SESSGET"
	TNS_DATA_TYPE_SESSREL        DataType = "TNS_DATA_TYPE_SESSREL"
	TNS_DATA_TYPE_XSS            DataType = "TNS_DATA_TYPE_XSS"
	TNS_DATA_TYPE_PDQCINV        DataType = "TNS_DATA_TYPE_PDQCINV"
	TNS_DATA_TYPE_PDQIDC         DataType = "TNS_DATA_TYPE_PDQIDC"
	TNS_DATA_TYPE_KPDQCSTA       DataType = "TNS_DATA_TYPE_KPDQCSTA"
	TNS_DATA_TYPE_KPRS           DataType = "TNS_DATA_TYPE_KPRS"
	TNS_DATA_TYPE_KPDQIDC        DataType = "TNS_DATA_TYPE_KPDQIDC"
	TNS_DATA_TYPE_RTSTRM         DataType = "TNS_DATA_TYPE_RTSTRM"
	TNS_DATA_TYPE_SESSRET        DataType = "TNS_DATA_TYPE_SESSRET"
	TNS_DATA_TYPE_SCN6           DataType = "TNS_DATA_TYPE_SCN6"
	TNS_DATA_TYPE_KECPA          DataType = "TNS_DATA_TYPE_KECPA"
	TNS_DATA_TYPE_KECPP          DataType = "TNS_DATA_TYPE_KECPP"
	TNS_DATA_TYPE_SXA            DataType = "TNS_DATA_TYPE_SXA"
	TNS_DATA_TYPE_KVARR          DataType = "TNS_DATA_TYPE_KVARR"
	TNS_DATA_TYPE_KPNGN          DataType = "TNS_DATA_TYPE_KPNGN"
	TNS_DATA_TYPE_XSNSOP         DataType = "TNS_DATA_TYPE_XSNSOP"
	TNS_DATA_TYPE_XSATTR         DataType = "TNS_DATA_TYPE_XSATTR"
	TNS_DATA_TYPE_XSNS           DataType = "TNS_DATA_TYPE_XSNS"
	TNS_DATA_TYPE_TXT            DataType = "TNS_DATA_TYPE_TXT"
	TNS_DATA_TYPE_XSSESSNS       DataType = "TNS_DATA_TYPE_XSSESSNS"
	TNS_DATA_TYPE_XSATTOP        DataType = "TNS_DATA_TYPE_XSATTOP"
	TNS_DATA_TYPE_XSCREOP        DataType = "TNS_DATA_TYPE_XSCREOP"
	TNS_DATA_TYPE_XSDETOP        DataType = "TNS_DATA_TYPE_XSDETOP"
	TNS_DATA_TYPE_XSDESOP        DataType = "TNS_DATA_TYPE_XSDESOP"
	TNS_DATA_TYPE_XSSETSP        DataType = "TNS_DATA_TYPE_XSSETSP"
	TNS_DATA_TYPE_XSSIDP         DataType = "TNS_DATA_TYPE_XSSIDP"
	TNS_DATA_TYPE_XSPRIN         DataType = "TNS_DATA_TYPE_XSPRIN"
	TNS_DATA_TYPE_XSKVL          DataType = "TNS_DATA_TYPE_XSKVL"
	TNS_DATA_TYPE_XSSS2          DataType = "TNS_DATA_TYPE_XSSS2"
	TNS_DATA_TYPE_XSNSOP2        DataType = "TNS_DATA_TYPE_XSNSOP2"
	TNS_DATA_TYPE_XSNS2          DataType = "TNS_DATA_TYPE_XSNS2"
	TNS_DATA_TYPE_IMPLRES        DataType = "TNS_DATA_TYPE_IMPLRES"
	TNS_DATA_TYPE_OER            DataType = "TNS_DATA_TYPE_OER"
	TNS_DATA_TYPE_UB1ARRAY       DataType = "TNS_DATA_TYPE_UB1ARRAY"
	TNS_DATA_TYPE_SESSSTATE      DataType = "TNS_DATA_TYPE_SESSSTATE"
	TNS_DATA_TYPE_AC_REPLAY      DataType = "TNS_DATA_TYPE_AC_REPLAY"
	TNS_DATA_TYPE_AC_CONT        DataType = "TNS_DATA_TYPE_AC_CONT"
	TNS_DATA_TYPE_KPDNREQ        DataType = "TNS_DATA_TYPE_KPDNREQ"
	TNS_DATA_TYPE_KPDNRNF        DataType = "TNS_DATA_TYPE_KPDNRNF"
	TNS_DATA_TYPE_KPNGNC         DataType = "TNS_DATA_TYPE_KPNGNC"
	TNS_DATA_TYPE_KPNRI          DataType = "TNS_DATA_TYPE_KPNRI"
	TNS_DATA_TYPE_AQENQ          DataType = "TNS_DATA_TYPE_AQENQ"
	TNS_DATA_TYPE_AQDEQ          DataType = "TNS_DATA_TYPE_AQDEQ"
	TNS_DATA_TYPE_AQJMS          DataType = "TNS_DATA_TYPE_AQJMS"
	TNS_DATA_TYPE_KPDNRPAY       DataType = "TNS_DATA_TYPE_KPDNRPAY"
	TNS_DATA_TYPE_KPDNRACK       DataType = "TNS_DATA_TYPE_KPDNRACK"
	TNS_DATA_TYPE_KPDNRMP        DataType = "TNS_DATA_TYPE_KPDNRMP"
	TNS_DATA_TYPE_KPDNRDQ        DataType = "TNS_DATA_TYPE_KPDNRDQ"
	TNS_DATA_TYPE_CHUNKINFO      DataType = "TNS_DATA_TYPE_CHUNKINFO"
	TNS_DATA_TYPE_SCN            DataType = "TNS_DATA_TYPE_SCN"
	TNS_DATA_TYPE_SCN8           DataType = "TNS_DATA_TYPE_SCN8"
	TNS_DATA_TYPE_UDS            DataType = "TNS_DATA_TYPE_UDS"
	TNS_DATA_TYPE_TNP            DataType = "TNS_DATA_TYPE_TNP"
)

var DataTypeValues = map[int]DataType{
	0:   DEFAULT_DATA_TYPE,
	1:   TNS_DATA_TYPE_VARCHAR,
	2:   TNS_DATA_TYPE_NUMBER,
	3:   TNS_DATA_TYPE_BINARY_INTEGER,
	4:   TNS_DATA_TYPE_FLOAT,
	5:   TNS_DATA_TYPE_STR,
	6:   TNS_DATA_TYPE_VNU,
	7:   TNS_DATA_TYPE_PDN,
	8:   TNS_DATA_TYPE_LONG,
	9:   TNS_DATA_TYPE_VCS,
	10:  TNS_DATA_TYPE_TID,
	11:  TNS_DATA_TYPE_ROWID,
	12:  TNS_DATA_TYPE_DATE,
	15:  TNS_DATA_TYPE_VBI,
	23:  TNS_DATA_TYPE_RAW,
	24:  TNS_DATA_TYPE_LONG_RAW,
	25:  TNS_DATA_TYPE_UB2,
	26:  TNS_DATA_TYPE_UB4,
	27:  TNS_DATA_TYPE_SB1,
	28:  TNS_DATA_TYPE_SB2,
	29:  TNS_DATA_TYPE_SB4,
	30:  TNS_DATA_TYPE_SWORD,
	31:  TNS_DATA_TYPE_UWORD,
	32:  TNS_DATA_TYPE_PTRB,
	33:  TNS_DATA_TYPE_PTRW,
	290: TNS_DATA_TYPE_OER8,
	291: TNS_DATA_TYPE_FUN,
	292: TNS_DATA_TYPE_AUA,
	293: TNS_DATA_TYPE_RXH7,
	294: TNS_DATA_TYPE_NA6,
	39:  TNS_DATA_TYPE_OAC,
	40:  TNS_DATA_TYPE_AMS,
	41:  TNS_DATA_TYPE_BRN,
	298: TNS_DATA_TYPE_BRP,
	299: TNS_DATA_TYPE_BRV,
	300: TNS_DATA_TYPE_KVA,
	301: TNS_DATA_TYPE_CLS,
	302: TNS_DATA_TYPE_CUI,
	303: TNS_DATA_TYPE_DFN,
	304: TNS_DATA_TYPE_DQR,
	305: TNS_DATA_TYPE_DSC,
	306: TNS_DATA_TYPE_EXE,
	307: TNS_DATA_TYPE_FCH,
	308: TNS_DATA_TYPE_GBV,
	309: TNS_DATA_TYPE_GEM,
	310: TNS_DATA_TYPE_GIV,
	311: TNS_DATA_TYPE_OKG,
	312: TNS_DATA_TYPE_HMI,
	313: TNS_DATA_TYPE_INO,
	315: TNS_DATA_TYPE_LNF,
	316: TNS_DATA_TYPE_ONT,
	317: TNS_DATA_TYPE_OPE,
	318: TNS_DATA_TYPE_OSQ,
	319: TNS_DATA_TYPE_SFE,
	320: TNS_DATA_TYPE_SPF,
	321: TNS_DATA_TYPE_VSN,
	322: TNS_DATA_TYPE_UD7,
	323: TNS_DATA_TYPE_DSA,
	68:  TNS_DATA_TYPE_UIN,
	325: TNS_DATA_TYPE_PIN,
	326: TNS_DATA_TYPE_PFN,
	327: TNS_DATA_TYPE_PPT,
	329: TNS_DATA_TYPE_STO,
	331: TNS_DATA_TYPE_ARC,
	332: TNS_DATA_TYPE_MRS,
	333: TNS_DATA_TYPE_MRT,
	334: TNS_DATA_TYPE_MRG,
	335: TNS_DATA_TYPE_MRR,
	336: TNS_DATA_TYPE_MRC,
	337: TNS_DATA_TYPE_VER,
	338: TNS_DATA_TYPE_LON2,
	339: TNS_DATA_TYPE_INO2,
	340: TNS_DATA_TYPE_ALL,
	341: TNS_DATA_TYPE_UDB,
	342: TNS_DATA_TYPE_AQI,
	343: TNS_DATA_TYPE_ULB,
	344: TNS_DATA_TYPE_ULD,
	91:  TNS_DATA_TYPE_SLS,
	346: TNS_DATA_TYPE_SID,
	347: TNS_DATA_TYPE_NA7,
	94:  TNS_DATA_TYPE_LVC,
	95:  TNS_DATA_TYPE_LVB,
	96:  TNS_DATA_TYPE_CHAR,
	97:  TNS_DATA_TYPE_AVC,
	354: TNS_DATA_TYPE_AL7,
	355: TNS_DATA_TYPE_K2RPC,
	100: TNS_DATA_TYPE_BINARY_FLOAT,
	101: TNS_DATA_TYPE_BINARY_DOUBLE,
	102: TNS_DATA_TYPE_CURSOR,
	104: TNS_DATA_TYPE_RDD,
	359: TNS_DATA_TYPE_XDP,
	106: TNS_DATA_TYPE_OSL,
	360: TNS_DATA_TYPE_OKO8,
	108: TNS_DATA_TYPE_EXT_NAMED,
	109: TNS_DATA_TYPE_INT_NAMED,
	110: TNS_DATA_TYPE_EXT_REF,
	111: TNS_DATA_TYPE_INT_REF,
	112: TNS_DATA_TYPE_CLOB,
	113: TNS_DATA_TYPE_BLOB,
	114: TNS_DATA_TYPE_BFILE,
	115: TNS_DATA_TYPE_CFILE,
	116: TNS_DATA_TYPE_RSET,
	117: TNS_DATA_TYPE_CWD,
	119: TNS_DATA_TYPE_JSON,
	120: TNS_DATA_TYPE_NEW_OAC,
	380: TNS_DATA_TYPE_UD12,
	381: TNS_DATA_TYPE_AL8,
	382: TNS_DATA_TYPE_LFOP,
	383: TNS_DATA_TYPE_FCRT,
	384: TNS_DATA_TYPE_DNY,
	385: TNS_DATA_TYPE_OPR,
	386: TNS_DATA_TYPE_PLS,
	387: TNS_DATA_TYPE_XID,
	388: TNS_DATA_TYPE_TXN,
	389: TNS_DATA_TYPE_DCB,
	390: TNS_DATA_TYPE_CCA,
	391: TNS_DATA_TYPE_WRN,
	393: TNS_DATA_TYPE_TLH,
	394: TNS_DATA_TYPE_TOH,
	395: TNS_DATA_TYPE_FOI,
	396: TNS_DATA_TYPE_SID2,
	397: TNS_DATA_TYPE_TCH,
	398: TNS_DATA_TYPE_PII,
	399: TNS_DATA_TYPE_PFI,
	400: TNS_DATA_TYPE_PPU,
	401: TNS_DATA_TYPE_PTE,
	146: TNS_DATA_TYPE_CLV,
	404: TNS_DATA_TYPE_RXH8,
	405: TNS_DATA_TYPE_N12,
	406: TNS_DATA_TYPE_AUTH,
	407: TNS_DATA_TYPE_KVAL,
	152: TNS_DATA_TYPE_DTR,
	153: TNS_DATA_TYPE_DUN,
	154: TNS_DATA_TYPE_DOP,
	155: TNS_DATA_TYPE_VST,
	156: TNS_DATA_TYPE_ODT,
	413: TNS_DATA_TYPE_FGI,
	414: TNS_DATA_TYPE_DSY,
	415: TNS_DATA_TYPE_DSYR8,
	416: TNS_DATA_TYPE_DSYH8,
	417: TNS_DATA_TYPE_DSYL,
	418: TNS_DATA_TYPE_DSYT8,
	419: TNS_DATA_TYPE_DSYV8,
	420: TNS_DATA_TYPE_DSYP,
	421: TNS_DATA_TYPE_DSYF,
	422: TNS_DATA_TYPE_DSYK,
	423: TNS_DATA_TYPE_DSYY,
	424: TNS_DATA_TYPE_DSYQ,
	425: TNS_DATA_TYPE_DSYC,
	426: TNS_DATA_TYPE_DSYA,
	427: TNS_DATA_TYPE_OT8,
	428: TNS_DATA_TYPE_DOL,
	429: TNS_DATA_TYPE_DSYTY,
	430: TNS_DATA_TYPE_AQE,
	431: TNS_DATA_TYPE_KV,
	432: TNS_DATA_TYPE_AQD,
	433: TNS_DATA_TYPE_AQ8,
	178: TNS_DATA_TYPE_TIME,
	179: TNS_DATA_TYPE_TIME_TZ,
	180: TNS_DATA_TYPE_TIMESTAMP,
	181: TNS_DATA_TYPE_TIMESTAMP_TZ,
	182: TNS_DATA_TYPE_INTERVAL_YM,
	183: TNS_DATA_TYPE_INTERVAL_DS,
	184: TNS_DATA_TYPE_EDATE,
	185: TNS_DATA_TYPE_ETIME,
	186: TNS_DATA_TYPE_ETTZ,
	187: TNS_DATA_TYPE_ESTAMP,
	188: TNS_DATA_TYPE_ESTZ,
	189: TNS_DATA_TYPE_EIYM,
	190: TNS_DATA_TYPE_EIDS,
	449: TNS_DATA_TYPE_RFS,
	450: TNS_DATA_TYPE_RXH10,
	195: TNS_DATA_TYPE_DCLOB,
	196: TNS_DATA_TYPE_DBLOB,
	197: TNS_DATA_TYPE_DBFILE,
	198: TNS_DATA_TYPE_DJSON,
	454: TNS_DATA_TYPE_KPN,
	455: TNS_DATA_TYPE_KPDNR,
	456: TNS_DATA_TYPE_DSYD,
	457: TNS_DATA_TYPE_DSYS,
	458: TNS_DATA_TYPE_DSYR,
	459: TNS_DATA_TYPE_DSYH,
	460: TNS_DATA_TYPE_DSYT,
	461: TNS_DATA_TYPE_DSYV,
	462: TNS_DATA_TYPE_AQM,
	463: TNS_DATA_TYPE_OER11,
	208: TNS_DATA_TYPE_UROWID,
	469: TNS_DATA_TYPE_AQL,
	470: TNS_DATA_TYPE_OTC,
	471: TNS_DATA_TYPE_KFNO,
	472: TNS_DATA_TYPE_KFNP,
	473: TNS_DATA_TYPE_KGT8,
	474: TNS_DATA_TYPE_RASB4,
	475: TNS_DATA_TYPE_RAUB2,
	476: TNS_DATA_TYPE_RAUB1,
	477: TNS_DATA_TYPE_RATXT,
	478: TNS_DATA_TYPE_RSSB4,
	479: TNS_DATA_TYPE_RSUB2,
	480: TNS_DATA_TYPE_RSUB1,
	481: TNS_DATA_TYPE_RSTXT,
	482: TNS_DATA_TYPE_RIDL,
	483: TNS_DATA_TYPE_GLRDD,
	484: TNS_DATA_TYPE_GLRDG,
	485: TNS_DATA_TYPE_GLRDC,
	486: TNS_DATA_TYPE_OKO,
	487: TNS_DATA_TYPE_DPP,
	488: TNS_DATA_TYPE_DPLS,
	489: TNS_DATA_TYPE_DPMOP,
	231: TNS_DATA_TYPE_TIMESTAMP_LTZ,
	232: TNS_DATA_TYPE_ESITZ,
	233: TNS_DATA_TYPE_UB8,
	490: TNS_DATA_TYPE_STAT,
	491: TNS_DATA_TYPE_RFX,
	492: TNS_DATA_TYPE_FAL,
	493: TNS_DATA_TYPE_CKV,
	494: TNS_DATA_TYPE_DRCX,
	495: TNS_DATA_TYPE_KGH,
	496: TNS_DATA_TYPE_AQO,
	241: TNS_DATA_TYPE_PNTY,
	498: TNS_DATA_TYPE_OKGT,
	499: TNS_DATA_TYPE_KPFC,
	500: TNS_DATA_TYPE_FE2,
	501: TNS_DATA_TYPE_SPFP,
	502: TNS_DATA_TYPE_DPULS,
	252: TNS_DATA_TYPE_BOOLEAN,
	507: TNS_DATA_TYPE_AQA,
	508: TNS_DATA_TYPE_KPBF,
	513: TNS_DATA_TYPE_TSM,
	514: TNS_DATA_TYPE_MSS,
	516: TNS_DATA_TYPE_KPC,
	517: TNS_DATA_TYPE_CRS,
	518: TNS_DATA_TYPE_KKS,
	519: TNS_DATA_TYPE_KSP,
	520: TNS_DATA_TYPE_KSPTOP,
	521: TNS_DATA_TYPE_KSPVAL,
	522: TNS_DATA_TYPE_PSS,
	523: TNS_DATA_TYPE_NLS,
	524: TNS_DATA_TYPE_ALS,
	525: TNS_DATA_TYPE_KSDEVTVAL,
	526: TNS_DATA_TYPE_KSDEVTTOP,
	527: TNS_DATA_TYPE_KPSPP,
	528: TNS_DATA_TYPE_KOL,
	529: TNS_DATA_TYPE_LST,
	530: TNS_DATA_TYPE_ACX,
	531: TNS_DATA_TYPE_SCS,
	532: TNS_DATA_TYPE_RXH,
	533: TNS_DATA_TYPE_KPDNS,
	534: TNS_DATA_TYPE_KPDCN,
	535: TNS_DATA_TYPE_KPNNS,
	536: TNS_DATA_TYPE_KPNCN,
	537: TNS_DATA_TYPE_KPS,
	538: TNS_DATA_TYPE_APINF,
	539: TNS_DATA_TYPE_TEN,
	540: TNS_DATA_TYPE_XSSCS,
	541: TNS_DATA_TYPE_XSSSO,
	542: TNS_DATA_TYPE_XSSAO,
	543: TNS_DATA_TYPE_KSRPC,
	560: TNS_DATA_TYPE_KVL,
	563: TNS_DATA_TYPE_SESSGET,
	564: TNS_DATA_TYPE_SESSREL,
	565: TNS_DATA_TYPE_XSS,
	572: TNS_DATA_TYPE_PDQCINV,
	573: TNS_DATA_TYPE_PDQIDC,
	574: TNS_DATA_TYPE_KPDQCSTA,
	575: TNS_DATA_TYPE_KPRS,
	576: TNS_DATA_TYPE_KPDQIDC,
	578: TNS_DATA_TYPE_RTSTRM,
	579: TNS_DATA_TYPE_SESSRET,
	580: TNS_DATA_TYPE_SCN6,
	581: TNS_DATA_TYPE_KECPA,
	582: TNS_DATA_TYPE_KECPP,
	583: TNS_DATA_TYPE_SXA,
	584: TNS_DATA_TYPE_KVARR,
	585: TNS_DATA_TYPE_KPNGN,
	590: TNS_DATA_TYPE_XSNSOP,
	591: TNS_DATA_TYPE_XSATTR,
	592: TNS_DATA_TYPE_XSNS,
	593: TNS_DATA_TYPE_TXT,
	594: TNS_DATA_TYPE_XSSESSNS,
	595: TNS_DATA_TYPE_XSATTOP,
	596: TNS_DATA_TYPE_XSCREOP,
	597: TNS_DATA_TYPE_XSDETOP,
	598: TNS_DATA_TYPE_XSDESOP,
	599: TNS_DATA_TYPE_XSSETSP,
	600: TNS_DATA_TYPE_XSSIDP,
	601: TNS_DATA_TYPE_XSPRIN,
	602: TNS_DATA_TYPE_XSKVL,
	603: TNS_DATA_TYPE_XSSS2,
	604: TNS_DATA_TYPE_XSNSOP2,
	605: TNS_DATA_TYPE_XSNS2,
	611: TNS_DATA_TYPE_IMPLRES,
	612: TNS_DATA_TYPE_OER,
	613: TNS_DATA_TYPE_UB1ARRAY,
	614: TNS_DATA_TYPE_SESSSTATE,
	615: TNS_DATA_TYPE_AC_REPLAY,
	616: TNS_DATA_TYPE_AC_CONT,
	622: TNS_DATA_TYPE_KPDNREQ,
	623: TNS_DATA_TYPE_KPDNRNF,
	624: TNS_DATA_TYPE_KPNGNC,
	625: TNS_DATA_TYPE_KPNRI,
	626: TNS_DATA_TYPE_AQENQ,
	627: TNS_DATA_TYPE_AQDEQ,
	628: TNS_DATA_TYPE_AQJMS,
	629: TNS_DATA_TYPE_KPDNRPAY,
	630: TNS_DATA_TYPE_KPDNRACK,
	631: TNS_DATA_TYPE_KPDNRMP,
	632: TNS_DATA_TYPE_KPDNRDQ,
	636: TNS_DATA_TYPE_CHUNKINFO,
	637: TNS_DATA_TYPE_SCN,
	638: TNS_DATA_TYPE_SCN8,
	639: TNS_DATA_TYPE_UDS,
	640: TNS_DATA_TYPE_TNP,
}

var ReverseDataTypeValues = map[DataType]int{
	TNS_DATA_TYPE_VARCHAR:        1,
	TNS_DATA_TYPE_NUMBER:         2,
	TNS_DATA_TYPE_BINARY_INTEGER: 3,
	TNS_DATA_TYPE_FLOAT:          4,
	TNS_DATA_TYPE_STR:            5,
	TNS_DATA_TYPE_VNU:            6,
	TNS_DATA_TYPE_PDN:            7,
	TNS_DATA_TYPE_LONG:           8,
	TNS_DATA_TYPE_VCS:            9,
	TNS_DATA_TYPE_TID:            10,
	TNS_DATA_TYPE_ROWID:          11,
	TNS_DATA_TYPE_DATE:           12,
	TNS_DATA_TYPE_VBI:            15,
	TNS_DATA_TYPE_RAW:            23,
	TNS_DATA_TYPE_LONG_RAW:       24,
	TNS_DATA_TYPE_UB2:            25,
	TNS_DATA_TYPE_UB4:            26,
	TNS_DATA_TYPE_SB1:            27,
	TNS_DATA_TYPE_SB2:            28,
	TNS_DATA_TYPE_SB4:            29,
	TNS_DATA_TYPE_SWORD:          30,
	TNS_DATA_TYPE_UWORD:          31,
	TNS_DATA_TYPE_PTRB:           32,
	TNS_DATA_TYPE_PTRW:           33,
	TNS_DATA_TYPE_OER8:           290,
	TNS_DATA_TYPE_FUN:            291,
	TNS_DATA_TYPE_AUA:            292,
	TNS_DATA_TYPE_RXH7:           293,
	TNS_DATA_TYPE_NA6:            294,
	TNS_DATA_TYPE_OAC:            39,
	TNS_DATA_TYPE_AMS:            40,
	TNS_DATA_TYPE_BRN:            41,
	TNS_DATA_TYPE_BRP:            298,
	TNS_DATA_TYPE_BRV:            299,
	TNS_DATA_TYPE_KVA:            300,
	TNS_DATA_TYPE_CLS:            301,
	TNS_DATA_TYPE_CUI:            302,
	TNS_DATA_TYPE_DFN:            303,
	TNS_DATA_TYPE_DQR:            304,
	TNS_DATA_TYPE_DSC:            305,
	TNS_DATA_TYPE_EXE:            306,
	TNS_DATA_TYPE_FCH:            307,
	TNS_DATA_TYPE_GBV:            308,
	TNS_DATA_TYPE_GEM:            309,
	TNS_DATA_TYPE_GIV:            310,
	TNS_DATA_TYPE_OKG:            311,
	TNS_DATA_TYPE_HMI:            312,
	TNS_DATA_TYPE_INO:            313,
	TNS_DATA_TYPE_LNF:            315,
	TNS_DATA_TYPE_ONT:            316,
	TNS_DATA_TYPE_OPE:            317,
	TNS_DATA_TYPE_OSQ:            318,
	TNS_DATA_TYPE_SFE:            319,
	TNS_DATA_TYPE_SPF:            320,
	TNS_DATA_TYPE_VSN:            321,
	TNS_DATA_TYPE_UD7:            322,
	TNS_DATA_TYPE_DSA:            323,
	TNS_DATA_TYPE_UIN:            68,
	TNS_DATA_TYPE_PIN:            325,
	TNS_DATA_TYPE_PFN:            326,
	TNS_DATA_TYPE_PPT:            327,
	TNS_DATA_TYPE_STO:            329,
	TNS_DATA_TYPE_ARC:            331,
	TNS_DATA_TYPE_MRS:            332,
	TNS_DATA_TYPE_MRT:            333,
	TNS_DATA_TYPE_MRG:            334,
	TNS_DATA_TYPE_MRR:            335,
	TNS_DATA_TYPE_MRC:            336,
	TNS_DATA_TYPE_VER:            337,
	TNS_DATA_TYPE_LON2:           338,
	TNS_DATA_TYPE_INO2:           339,
	TNS_DATA_TYPE_ALL:            340,
	TNS_DATA_TYPE_UDB:            341,
	TNS_DATA_TYPE_AQI:            342,
	TNS_DATA_TYPE_ULB:            343,
	TNS_DATA_TYPE_ULD:            344,
	TNS_DATA_TYPE_SLS:            91,
	TNS_DATA_TYPE_SID:            346,
	TNS_DATA_TYPE_NA7:            347,
	TNS_DATA_TYPE_LVC:            94,
	TNS_DATA_TYPE_LVB:            95,
	TNS_DATA_TYPE_CHAR:           96,
	TNS_DATA_TYPE_AVC:            97,
	TNS_DATA_TYPE_AL7:            354,
	TNS_DATA_TYPE_K2RPC:          355,
	TNS_DATA_TYPE_BINARY_FLOAT:   100,
	TNS_DATA_TYPE_BINARY_DOUBLE:  101,
	TNS_DATA_TYPE_CURSOR:         102,
	TNS_DATA_TYPE_RDD:            104,
	TNS_DATA_TYPE_XDP:            359,
	TNS_DATA_TYPE_OSL:            106,
	TNS_DATA_TYPE_OKO8:           360,
	TNS_DATA_TYPE_EXT_NAMED:      108,
	TNS_DATA_TYPE_INT_NAMED:      109,
	TNS_DATA_TYPE_EXT_REF:        110,
	TNS_DATA_TYPE_INT_REF:        111,
	TNS_DATA_TYPE_CLOB:           112,
	TNS_DATA_TYPE_BLOB:           113,
	TNS_DATA_TYPE_BFILE:          114,
	TNS_DATA_TYPE_CFILE:          115,
	TNS_DATA_TYPE_RSET:           116,
	TNS_DATA_TYPE_CWD:            117,
	TNS_DATA_TYPE_JSON:           119,
	TNS_DATA_TYPE_NEW_OAC:        120,
	TNS_DATA_TYPE_UD12:           380,
	TNS_DATA_TYPE_AL8:            381,
	TNS_DATA_TYPE_LFOP:           382,
	TNS_DATA_TYPE_FCRT:           383,
	TNS_DATA_TYPE_DNY:            384,
	TNS_DATA_TYPE_OPR:            385,
	TNS_DATA_TYPE_PLS:            386,
	TNS_DATA_TYPE_XID:            387,
	TNS_DATA_TYPE_TXN:            388,
	TNS_DATA_TYPE_DCB:            389,
	TNS_DATA_TYPE_CCA:            390,
	TNS_DATA_TYPE_WRN:            391,
	TNS_DATA_TYPE_TLH:            393,
	TNS_DATA_TYPE_TOH:            394,
	TNS_DATA_TYPE_FOI:            395,
	TNS_DATA_TYPE_SID2:           396,
	TNS_DATA_TYPE_TCH:            397,
	TNS_DATA_TYPE_PII:            398,
	TNS_DATA_TYPE_PFI:            399,
	TNS_DATA_TYPE_PPU:            400,
	TNS_DATA_TYPE_PTE:            401,
	TNS_DATA_TYPE_CLV:            146,
	TNS_DATA_TYPE_RXH8:           404,
	TNS_DATA_TYPE_N12:            405,
	TNS_DATA_TYPE_AUTH:           406,
	TNS_DATA_TYPE_KVAL:           407,
	TNS_DATA_TYPE_DTR:            152,
	TNS_DATA_TYPE_DUN:            153,
	TNS_DATA_TYPE_DOP:            154,
	TNS_DATA_TYPE_VST:            155,
	TNS_DATA_TYPE_ODT:            156,
	TNS_DATA_TYPE_FGI:            413,
	TNS_DATA_TYPE_DSY:            414,
	TNS_DATA_TYPE_DSYR8:          415,
	TNS_DATA_TYPE_DSYH8:          416,
	TNS_DATA_TYPE_DSYL:           417,
	TNS_DATA_TYPE_DSYT8:          418,
	TNS_DATA_TYPE_DSYV8:          419,
	TNS_DATA_TYPE_DSYP:           420,
	TNS_DATA_TYPE_DSYF:           421,
	TNS_DATA_TYPE_DSYK:           422,
	TNS_DATA_TYPE_DSYY:           423,
	TNS_DATA_TYPE_DSYQ:           424,
	TNS_DATA_TYPE_DSYC:           425,
	TNS_DATA_TYPE_DSYA:           426,
	TNS_DATA_TYPE_OT8:            427,
	TNS_DATA_TYPE_DOL:            428,
	TNS_DATA_TYPE_DSYTY:          429,
	TNS_DATA_TYPE_AQE:            430,
	TNS_DATA_TYPE_KV:             431,
	TNS_DATA_TYPE_AQD:            432,
	TNS_DATA_TYPE_AQ8:            433,
	TNS_DATA_TYPE_TIME:           178,
	TNS_DATA_TYPE_TIME_TZ:        179,
	TNS_DATA_TYPE_TIMESTAMP:      180,
	TNS_DATA_TYPE_TIMESTAMP_TZ:   181,
	TNS_DATA_TYPE_INTERVAL_YM:    182,
	TNS_DATA_TYPE_INTERVAL_DS:    183,
	TNS_DATA_TYPE_EDATE:          184,
	TNS_DATA_TYPE_ETIME:          185,
	TNS_DATA_TYPE_ETTZ:           186,
	TNS_DATA_TYPE_ESTAMP:         187,
	TNS_DATA_TYPE_ESTZ:           188,
	TNS_DATA_TYPE_EIYM:           189,
	TNS_DATA_TYPE_EIDS:           190,
	TNS_DATA_TYPE_RFS:            449,
	TNS_DATA_TYPE_RXH10:          450,
	TNS_DATA_TYPE_DCLOB:          195,
	TNS_DATA_TYPE_DBLOB:          196,
	TNS_DATA_TYPE_DBFILE:         197,
	TNS_DATA_TYPE_DJSON:          198,
	TNS_DATA_TYPE_KPN:            454,
	TNS_DATA_TYPE_KPDNR:          455,
	TNS_DATA_TYPE_DSYD:           456,
	TNS_DATA_TYPE_DSYS:           457,
	TNS_DATA_TYPE_DSYR:           458,
	TNS_DATA_TYPE_DSYH:           459,
	TNS_DATA_TYPE_DSYT:           460,
	TNS_DATA_TYPE_DSYV:           461,
	TNS_DATA_TYPE_AQM:            462,
	TNS_DATA_TYPE_OER11:          463,
	TNS_DATA_TYPE_UROWID:         208,
	TNS_DATA_TYPE_AQL:            469,
	TNS_DATA_TYPE_OTC:            470,
	TNS_DATA_TYPE_KFNO:           471,
	TNS_DATA_TYPE_KFNP:           472,
	TNS_DATA_TYPE_KGT8:           473,
	TNS_DATA_TYPE_RASB4:          474,
	TNS_DATA_TYPE_RAUB2:          475,
	TNS_DATA_TYPE_RAUB1:          476,
	TNS_DATA_TYPE_RATXT:          477,
	TNS_DATA_TYPE_RSSB4:          478,
	TNS_DATA_TYPE_RSUB2:          479,
	TNS_DATA_TYPE_RSUB1:          480,
	TNS_DATA_TYPE_RSTXT:          481,
	TNS_DATA_TYPE_RIDL:           482,
	TNS_DATA_TYPE_GLRDD:          483,
	TNS_DATA_TYPE_GLRDG:          484,
	TNS_DATA_TYPE_GLRDC:          485,
	TNS_DATA_TYPE_OKO:            486,
	TNS_DATA_TYPE_DPP:            487,
	TNS_DATA_TYPE_DPLS:           488,
	TNS_DATA_TYPE_DPMOP:          489,
	TNS_DATA_TYPE_TIMESTAMP_LTZ:  231,
	TNS_DATA_TYPE_ESITZ:          232,
	TNS_DATA_TYPE_UB8:            233,
	TNS_DATA_TYPE_STAT:           490,
	TNS_DATA_TYPE_RFX:            491,
	TNS_DATA_TYPE_FAL:            492,
	TNS_DATA_TYPE_CKV:            493,
	TNS_DATA_TYPE_DRCX:           494,
	TNS_DATA_TYPE_KGH:            495,
	TNS_DATA_TYPE_AQO:            496,
	TNS_DATA_TYPE_PNTY:           241,
	TNS_DATA_TYPE_OKGT:           498,
	TNS_DATA_TYPE_KPFC:           499,
	TNS_DATA_TYPE_FE2:            500,
	TNS_DATA_TYPE_SPFP:           501,
	TNS_DATA_TYPE_DPULS:          502,
	TNS_DATA_TYPE_BOOLEAN:        252,
	TNS_DATA_TYPE_AQA:            507,
	TNS_DATA_TYPE_KPBF:           508,
	TNS_DATA_TYPE_TSM:            513,
	TNS_DATA_TYPE_MSS:            514,
	TNS_DATA_TYPE_KPC:            516,
	TNS_DATA_TYPE_CRS:            517,
	TNS_DATA_TYPE_KKS:            518,
	TNS_DATA_TYPE_KSP:            519,
	TNS_DATA_TYPE_KSPTOP:         520,
	TNS_DATA_TYPE_KSPVAL:         521,
	TNS_DATA_TYPE_PSS:            522,
	TNS_DATA_TYPE_NLS:            523,
	TNS_DATA_TYPE_ALS:            524,
	TNS_DATA_TYPE_KSDEVTVAL:      525,
	TNS_DATA_TYPE_KSDEVTTOP:      526,
	TNS_DATA_TYPE_KPSPP:          527,
	TNS_DATA_TYPE_KOL:            528,
	TNS_DATA_TYPE_LST:            529,
	TNS_DATA_TYPE_ACX:            530,
	TNS_DATA_TYPE_SCS:            531,
	TNS_DATA_TYPE_RXH:            532,
	TNS_DATA_TYPE_KPDNS:          533,
	TNS_DATA_TYPE_KPDCN:          534,
	TNS_DATA_TYPE_KPNNS:          535,
	TNS_DATA_TYPE_KPNCN:          536,
	TNS_DATA_TYPE_KPS:            537,
	TNS_DATA_TYPE_APINF:          538,
	TNS_DATA_TYPE_TEN:            539,
	TNS_DATA_TYPE_XSSCS:          540,
	TNS_DATA_TYPE_XSSSO:          541,
	TNS_DATA_TYPE_XSSAO:          542,
	TNS_DATA_TYPE_KSRPC:          543,
	TNS_DATA_TYPE_KVL:            560,
	TNS_DATA_TYPE_SESSGET:        563,
	TNS_DATA_TYPE_SESSREL:        564,
	TNS_DATA_TYPE_XSS:            565,
	TNS_DATA_TYPE_PDQCINV:        572,
	TNS_DATA_TYPE_PDQIDC:         573,
	TNS_DATA_TYPE_KPDQCSTA:       574,
	TNS_DATA_TYPE_KPRS:           575,
	TNS_DATA_TYPE_KPDQIDC:        576,
	TNS_DATA_TYPE_RTSTRM:         578,
	TNS_DATA_TYPE_SESSRET:        579,
	TNS_DATA_TYPE_SCN6:           580,
	TNS_DATA_TYPE_KECPA:          581,
	TNS_DATA_TYPE_KECPP:          582,
	TNS_DATA_TYPE_SXA:            583,
	TNS_DATA_TYPE_KVARR:          584,
	TNS_DATA_TYPE_KPNGN:          585,
	TNS_DATA_TYPE_XSNSOP:         590,
	TNS_DATA_TYPE_XSATTR:         591,
	TNS_DATA_TYPE_XSNS:           592,
	TNS_DATA_TYPE_TXT:            593,
	TNS_DATA_TYPE_XSSESSNS:       594,
	TNS_DATA_TYPE_XSATTOP:        595,
	TNS_DATA_TYPE_XSCREOP:        596,
	TNS_DATA_TYPE_XSDETOP:        597,
	TNS_DATA_TYPE_XSDESOP:        598,
	TNS_DATA_TYPE_XSSETSP:        599,
	TNS_DATA_TYPE_XSSIDP:         600,
	TNS_DATA_TYPE_XSPRIN:         601,
	TNS_DATA_TYPE_XSKVL:          602,
	TNS_DATA_TYPE_XSSS2:          603,
	TNS_DATA_TYPE_XSNSOP2:        604,
	TNS_DATA_TYPE_XSNS2:          605,
	TNS_DATA_TYPE_IMPLRES:        611,
	TNS_DATA_TYPE_OER:            612,
	TNS_DATA_TYPE_UB1ARRAY:       613,
	TNS_DATA_TYPE_SESSSTATE:      614,
	TNS_DATA_TYPE_AC_REPLAY:      615,
	TNS_DATA_TYPE_AC_CONT:        616,
	TNS_DATA_TYPE_KPDNREQ:        622,
	TNS_DATA_TYPE_KPDNRNF:        623,
	TNS_DATA_TYPE_KPNGNC:         624,
	TNS_DATA_TYPE_KPNRI:          625,
	TNS_DATA_TYPE_AQENQ:          626,
	TNS_DATA_TYPE_AQDEQ:          627,
	TNS_DATA_TYPE_AQJMS:          628,
	TNS_DATA_TYPE_KPDNRPAY:       629,
	TNS_DATA_TYPE_KPDNRACK:       630,
	TNS_DATA_TYPE_KPDNRMP:        631,
	TNS_DATA_TYPE_KPDNRDQ:        632,
	TNS_DATA_TYPE_CHUNKINFO:      636,
	TNS_DATA_TYPE_SCN:            637,
	TNS_DATA_TYPE_SCN8:           638,
	TNS_DATA_TYPE_UDS:            639,
	TNS_DATA_TYPE_TNP:            640,
}

type OracleRequest struct {
	Header    OracleHeader
	Message   interface{}
	ReadDelay int64
}

type OracleResponse struct {
	Header    OracleHeader
	Message   interface{}
	ReadDelay int64
}

type OracleHeader struct {
	PacketLength int
	PacketType   PacketType
	PacketFlag   uint8
	Session      PacketSession
}

type PacketSession struct {
	Context               *network.SessionContext
	TimeZone              []byte
	TTCVersion            uint8
	HasEOSCapability      bool
	HasFSAPCapability     bool
	Summary               *SummaryObject
	States                []network.SessionState
	StrConv               converters.IStringConverter
	SStrConv              converters.IStringConverter
	CStrConv              converters.IStringConverter
	NStrConv              converters.IStringConverter
	ServerCharacterSet    int
	ServernCharacterSet   int
	UseBigClrChunks       bool
	UseBigScn             bool
	ClrChunkSize          int
	SupportOOB            bool
	HandShakeComplete     bool
	CompileTimeCaps       []byte
	RunTimeCaps           []byte
	ServerCompileTimeCaps []byte
	ServerRunTimeCaps     []byte
}

type DecodeColumnValue struct {
	MaxSize          int
	Flag             byte
	TempByte         byte
	BVlaue           []byte
	OraclePrimeValue OraclePrimeValue
}

type CalculateColumnValue struct {
	RefCursor         RefCursor
	DecodeColumnValue DecodeColumnValue
}

type Lob struct {
	sourceLocator []byte
	destLocator   []byte
	scn           []byte
	sourceOffset  int64
	destOffset    int64
	sourceLen     int
	destLen       int
	charsetID     int
	size          int64
	data          bytes.Buffer
	bNullO2U      bool
	isNull        bool
	sendSize      bool
}

type BindError struct {
	ErrorCode int
	RowOffset int
	ErrorMsg  []byte
}

type SummaryObject struct {
	EndOfCallStatus      int // uint32
	EndToEndECIDSequence int // uint16
	CurRowNumber         int // uint32
	RetCode              int // uint16
	ArrayElmWError       int // uint16
	ArrayElmErrno        int //uint16
	CursorID             int // uint16
	ErrorPos             int // uint16
	SqlType              uint8
	OerFatal             uint8
	Flags                int // uint16
	UserCursorOPT        int // uint16
	UpiParam             uint8
	WarningFlag          uint8
	Rba                  int // uint32
	PartitionID          int // uint16
	TableID              uint8
	BlockNumber          int // uint32
	SlotNumber           int // uint16
	OsError              int // uint32
	StmtNumber           uint8
	CallNumber           uint8
	Pad1                 int // uint16
	SuccessIter          int // uint16
	ErrorMessage         string
	BindErrors           []BindError
}

type OracleConnectMessage struct {
	TnsVersionDesired          uint16
	TnsVersionMinimum          uint16
	ServiceOptions             uint16
	TNS_SDU_16                 uint16
	TNS_TDU_16                 uint16
	TNS_SDU_32                 uint32
	TNS_TDU_32                 uint32
	TnsProtocolCharacteristics uint16
	LineTurnaround             uint16
	OurOne                     uint16
	ConnectionStringLength     uint16
	OffsetOfConnectionData     uint16
	ACFL0                      uint8
	ACFL1                      uint8
	ConnectFlag1               uint32
	ConnectFlag2               uint32
	ConnectString              string
}

type OracleAcceptMessage struct {
	TnsVersion      uint16
	Sid             []byte
	ProtocolOptions uint16
	HistOne         uint16
	ACFL0           uint8
	ACFL1           uint8
	SDU_16          uint16
	TDU_16          uint16
	SDU_32          uint32
	TDU_32          uint32
	DataOffset      uint16
	DataLength      uint16
	ReconAddStart   uint16
	ReconAddLength  uint16
	ReconAdd        string
	Buffer          []byte
}

type OracleRefuseMessage struct {
	DataOffset   uint16
	SystemReason uint8
	UserReason   uint8
	DataLength   uint16
	Data         string
	OracleError  OracleError
}

type OracleRedirectMessage struct {
	DataOffset      uint16
	DataLength      uint16
	RedirectAddress string
	RedirectData    string
}

type OracleMarkerMessage struct {
	MarkerType MarkerType
	MarkerData uint8
}

type OracleControlMessage struct {
	ControlType  ControlType
	ControlError ControlError
}

type OracleDataMessage struct {
	DataOffset      uint16
	DataMessageType DataPacketType
	DataMessage     interface{}
}

type OracleAdvNegoMessage struct {
	AdvNegoHeaderMessage AdvNegoHeaderMessage
	ServiceDataList      []ServiceData
}

type ServiceData struct {
	AdvNegoServiceHeader AdvNegoServiceHeader
	ServiceData          interface{}
}

type AuthServiceData struct {
	Version      ServicePacket
	Unknown      ServicePacket
	Status       ServicePacket
	ServiceIds   []ServicePacket
	ServiceNames []ServicePacket
}

type AuthServiceDataResponse struct {
	Version     ServicePacket
	Status      ServicePacket
	ServiceName ServicePacket
}

type ServicePacket struct {
	Length int
	Type   int
	Value  interface{}
}

type EncrytionServiceData struct {
	Version    ServicePacket
	ServiceIds ServicePacket
}

type EncrytionServiceDataResponse struct {
	Version   ServicePacket
	ServiceId ServicePacket
}

type DataServiceData struct {
	Version    ServicePacket
	ServiceIds ServicePacket
}

type DataServiceDataResponse struct {
	Version              ServicePacket
	ServiceId            ServicePacket
	DhGenLen             ServicePacket
	DPrimLen             ServicePacket
	GenBytes             ServicePacket
	PrimeBytes           ServicePacket
	ServerPublicKeyBytes ServicePacket
	IV                   ServicePacket
}

type SupervisorServiceData struct {
	Version      ServicePacket
	ServiceIds   ServicePacket
	ServiceArray ServicePacket
}

type SupervisorServiceDataResponse struct {
	Version      ServicePacket
	Status       ServicePacket
	ServiceArray ServicePacket
}

type AdvNegoHeaderMessage struct {
	Length    int
	Version   int
	ServCount int
	ErrFlags  int
}

type AdvNegoServiceHeader struct {
	ServiceType       ServiceType
	ServiceSubPackets int
}

type OracleResponseDataMessage struct {
	MessageType DataPacketType
	MessageData interface{}
}

type OracleResponseRowHeaderMessage struct {
	ColumnCount     int
	Num             int
	RowCount        int
	UacBufferLength int
	BitVector       []byte
}

type OracleResponseBitVectorMessage struct {
	NumCloumnsSent int
	BitVector      []byte
}

type OracleResponseFlushOutBindsMesage struct {
}

type OracleResponseIOVectorMessage struct {
	ColumnCount            int
	Num                    int
	RowCount               int
	UacBufferLength        int
	BitVector              []byte
	ParameterDirectorArray []ParameterDirection
}

type OracleResponseDescribeInfoMessage struct {
	Size          uint8
	MaxRowSize    int
	ColumnCount   int
	ParamInfoList []ParameterInfo
}

type OracleResponseErrorMessage struct {
	Summary *SummaryObject
}

type OracleResponseParameterMessage struct {
	Size1              int
	Size2              int
	Size3              int
	ScnForSnapshotList []int
	KeyList            [][]byte
	ValList            [][]byte
	NumList            []int
	Bty                []byte
	Length             int
}

type OracleResponseRowMessage struct {
	NumList                     []int
	CalculateParameterValueList []OraclePrimeValue
	CalculateColumnValueList    []CalculateColumnValue
	CursorArray                 []RefCursor
}
type OraclePrimeValue struct {
	Size               int
	ParameterValueList interface{}
	RowId              Rowid
	UrowId             Urowid
	Bvalue             []byte
	DecodeObj          DecodeObject
}

type RefCursor struct {
	Stmt
	Len           uint8
	MaxRowSize    int
	Parent        *Stmt
	ColumnCount   int
	ParamInfoList []ParameterInfo
}

type OracleConnectionDataMessage struct {
	ConnectString string
}

type OracleRedirectDataMessage struct {
	RedirectAddress string
	RedirectData    string
}

type OracleProtocolDataMessageRequest struct {
	ProtocolVersion     uint8
	ArrayTerminatorList []uint8
	DriveName           string
}

type OracleProtocolDataMessageResponse struct {
	ProtocolVersion             uint8
	ArrayTerminatorList         []uint8
	ProtocolServerName          string
	ServerCharacterSet          int
	ServerFlags                 uint8
	CharacterSetElement         int
	ArrayLength                 int
	NumberArray                 []byte
	ServerCompileTimeCapsLength uint8
	ServerRunTimeCapsLength     uint8
	ServerCompileTimeCaps       []byte
	ServerRunTimeCaps           []byte
}

type OracleDataTypeDataMessageRequest struct {
	ServerCharacterSet    int
	ServernCharacterSet   int
	ServerFlags           uint8
	CompileTimeCapsLength uint8
	RunTimeCapsLength     uint8
	CompileTimeCaps       []byte
	RunTimeCaps           []byte
	ClientTZVersion       int
	DataType              []DataTypeInfo
}

type OracleDataTypeDataMessageResponse struct {
	ServerTZVersion int
	DataType        []DataTypeInfo
}

type OracleFuntionTypeDataMessage struct {
	FunctionCode   FunctionType
	FunctionData   interface{}
	SequenceNumber uint8
}

type OracleFetchFunctionTypeDataMessage struct {
	CursorId    int
	RowsToFetch int
	Stmt        Stmt
}

type OracleLOBFunctionTypeDataMessage struct {
	DestLength    int
	SourceLength  int
	SourceOffSet  int
	DestOffSet    int
	SendSize      bool
	BNullO2U      bool
	OperationId   int
	SCNLength     int
	SourceLocator []byte
	DestLocator   []byte
	CharsetId     int
	SCN           []int
	Size          int
	IsCharsetId   uint8
}

type OracleRowHeaderTypeDataMessage struct {
	Flags        uint8
	ColumnCount  int
	IterationNum int
	RowCount     int
	BufferLength int
	BitVector    []byte
}

type OracleAuthDataMessageRequest struct {
	HasUser        uint8
	UserByteLength uint32
	AuthMode       AuthMode
	NumPair        uint32
	UserBytes      string
	AuthKeyValue   []OracleKeyValue
}

type OracleAuthDataMessageParameter struct {
	NumParams    uint32
	AuthKeyValue []OracleKeyValue
}

type OracleDataMessageWarning struct {
	ErrCode   int
	ErrLength int
	Flag      int
	ErrMsg    string
}

type OracleDataMessageServerSidePiggyback struct {
	OpCode         PiggyBackType
	ServerSideInfo interface{}
}

type OracleServerPiggybackSync struct {
	DTYCount     int
	DTYLength    uint8
	NUM_PAIRS    int
	Length       uint8
	SessionProps []OracleKeyValue
	Flags        int
}

type OracleServerPiggybackLtxId struct {
	Length        int
	TransactionID []byte
}

type OracleServerPiggybackSessionReturn struct {
	SkipInt      int
	SkipByte     byte
	Length       int
	SessionProps []OracleKeyValue
	Flag         int
	SessionId    int
	SerialId     int
}

type OracleServerPiggybackReplayCtx struct {
	DTYCount      int
	DTYLength     byte
	Flags         int
	ErrorCode     int
	Queue         byte
	Length        int
	ReplayContext []byte
}

type OracleServerPiggybackExtSync struct {
	DTYCount  int
	DTYLength byte
}

type OracleServerPiggybackOsPid struct {
	Length int
	Byte   byte
	Pid    []byte
}

type OracleKeyValue struct {
	Key   string
	Value string
	Code  int
}

type DataTypeInfo struct {
	DataType       DataType
	ConvDataType   DataType
	Representation DataType
}

type Stmt struct {
	ReSendParDef       bool
	Parse              bool
	Execute            bool
	Define             bool
	BulkExec           bool
	Text               string
	DisableCompression bool
	HasLONG            bool
	HasBLOB            bool
	HasMoreRows        bool
	HasReturnClause    bool
	NoOfRowsToFetch    int
	StmtType           StmtType
	CursorID           int
	QueryID            uint64
	Pars               []ParameterInfo
	Columns            []ParameterInfo
	ScnForSnapshot     []int
	ArrayBindCount     int
	ContainOutputPars  bool
	AutoClose          bool
	DataSet            DataSet
}

type DataSet struct {
	ColumnCount     int
	RowCount        int
	UACBufferLength int
	MaxRowSize      int
	Cols            []ParameterInfo
	Rows            []Row
	CurrentRow      Row
	Lasterr         error
	Index           int
}

type Row []driver.Value

type ParameterInfo struct {
	Name                 string
	TypeName             string
	Direction            ParameterDirection
	IsNull               bool
	AllowNull            bool
	ColAlias             string
	DataType             TNSType
	IsXmlType            bool
	Flag                 uint8
	Precision            uint8
	Scale                uint8
	MaxLen               int
	MaxCharLen           int
	MaxNoOfArrayElements int
	ContFlag             int
	ToID                 []byte
	Version              int
	CharsetID            int
	CharsetForm          int
	BValue               []byte
	Value                driver.Value
	IPrimValue           driver.Value
	OPrimValue           driver.Value
	OutputVarPtr         interface{}
	GetDataFromServer    bool
	Oaccollid            int
	CusType              *CustomType
}

type Rowid struct {
	RBA         int64
	PartitionId int64
	Filter      byte
	BlockNumber int64
	SlotNumber  int64
}

type Urowid struct {
	Data []byte
	Rowid
}

type BFile struct {
	isOpened bool
	lob      Lob
}

type DecodeObject struct {
	ObjType               byte
	Ctl                   int
	ItemLen               int
	BValueArray           [][]byte
	DecodeObjArray        interface{}
	DecodePrimeValueArray []OraclePrimeValue
}

type CustomType struct {
	Owner         string
	Name          string
	ArrayTypeName string
	Attribs       []ParameterInfo
	Typ           reflect.Type
	Toid          []byte // type oid
	ArrayTOID     []byte
	FieldMap      map[string]int
}

type OracleExecuteMessage struct {
	PiggyBack PiggyBackMsg
	Stmt      Stmt
	ExeOp     int
	LOB       int
	Al8i4     []int
	Parse     bool
	Define    bool
	NumPars   int
	NumCols   int
}

type OracleReExecuteMessage struct {
	PiggyBack PiggyBackMsg
	Stmt      Stmt
	Count     int
	ExeOp     int
	ExecFlag  int
}

type OracleReExecuteAndFetchMessage struct {
	PiggyBack PiggyBackMsg
	Stmt      Stmt
	Count     int
	ExeOp     int
	ExecFlag  int
}

type PiggyBackMsg struct {
	PiggyBackCode  FunctionType
	PiggyBackData  interface{}
	SequenceNumber uint8
}

type SchemaPiggyBackMsg struct {
	SchemaBytes []byte
}

type CloseCursorPiggyBackMsg struct {
	CursorIds []int
}

type OracleDataMessageStatus struct {
	EndOfCallStatus      int
	EndToEndECIDSequence int
}

type CloseTempLobsPiggyBackMsg struct {
	LobSize    int
	OpCode     int
	LobToClose [][]byte
}

type EndToEndPiggyBackMsg struct {
	Flags                    int
	ClientIdentifierModified bool
	ClientIdentifier         bool
	ModuleModified           bool
	Module                   bool
	ActionModified           bool
	Action                   bool
	ClientInfoModified       bool
	ClientInfo               bool
	DbopModified             bool
	Dbop                     bool
	ClientIdentifierBytes    []byte
	ModuleBytes              []byte
	ActionBytes              []byte
	ClientInfoBytes          []byte
	DbopBytes                []byte
}

type DefineData struct {
	DataType             TNSType
	Flag                 uint8
	Precision            uint8
	Scale                uint8
	MaxLen               int
	MaxNoOfArrayElements int
	ContFlag             int
	ToID                 []byte
	Version              int
	CharsetID            int
	CharsetForm          int
	MaxCharLen           int
	Oaccollid            int
}

type OracleError struct {
	ErrCode int
	ErrMsg  string
}

type OracleDBVersionMessage struct {
	Length int
	Info   string
	Number int
}

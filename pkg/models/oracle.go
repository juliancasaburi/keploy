package models

import (
	"strconv"

	"github.com/sijms/go-ora/v2/converters"
	"github.com/sijms/go-ora/v2/network"
)

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
	PacketLength interface{}
	PacketType   network.PacketType
	PacketFlag   uint8
	Session      PacketSession
}

type PacketSession struct {
	Context               *network.SessionContext
	TimeZone              []byte
	TTCVersion            uint8
	HasEOSCapability      bool
	HasFSAPCapability     bool
	Summary               *network.SummaryObject
	states                []network.SessionState
	StrConv               converters.IStringConverter
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

type OracleConnectMessage struct {
	TNS_VERSION_DESIRED          uint16
	TNS_VERSION_MINIMUM          uint16
	SERVICE_OPTIONS              uint16
	TNS_SDU_16                   uint16
	TNS_TDU_16                   uint16
	TNS_SDU_32                   uint32
	TNS_TDU_32                   uint32
	TNS_PROTOCOL_CHARACTERISTICS uint16
	LINE_TURNAROUND              uint16
	OURONE                       uint16
	CONNECTION_STRING_LENGTH     uint16
	OFFSET_OF_CONNECTION_DATA    uint16
	ACFL0                        uint8
	ACFL1                        uint8
	CONNECT_FLAG_1               uint32
	CONNECT_FLAG_2               uint32
	CONNECT_STRING               string
}

type OracleAcceptMessage struct {
	TNS_VERSION      uint16
	SID              []byte
	PROTOCOL_OPTIONS uint16
	HISTONE          uint16
	ACFL0            uint8
	ACFL1            uint8
	SDU_16           uint16
	TDU_16           uint16
	SDU_32           uint32
	TDU_32           uint32
	DATA_OFFSET      uint16
	DATA_LENGTH      uint16
	RECON_ADD_START  uint16
	RECON_ADD_LENGTH uint16
	RECON_ADD        string
	BUFFER           []byte
}

type OracleRefuseMessage struct {
	DATA_OFFSET   uint16
	SYSTEM_REASON uint8
	USER_REASON   uint8
	DATA_LENGTH   uint16
	DATA          string
	ORACLE_ERROR  OracleError
}

type OracleRedirectMessage struct {
	DATA_OFFSET      uint16
	DATA_LENGTH      uint16
	REDIRECT_ADDRESS string
	REDIRECT_DATA    string
}

type OracleMarkerMessage struct {
	MARKER_TYPE uint8
	MARKER_DATA uint8
}

type OracleControlMessage struct {
	CONTROL_TYPE  ControlType
	CONTROL_ERROR ControlError
}

type OracleDataMessage struct {
	DATA_OFFSET       uint16
	DATA_MESSAGE_TYPE DataPacketType
	DATA_MESSAGE      interface{}
}

type OracleConnectionDataMessage struct {
	CONNECT_STRING string
}

type OracleRedirectDataMessage struct {
	REDIRECT_ADDRESS string
	REDIRECT_DATA    string
}

type OracleProtocolDataMessageRequest struct {
	PROTOCOL_VERSION      uint8
	ARRAY_TERMINATOR_LIST []uint8
	DRIVER_NAME           string
}

type OracleProtocolDataMessageResponse struct {
	PROTOCOL_VERSION                uint8
	ARRAY_TERMINATOR_LIST           []uint8
	PROTOCOL_SERVER_NAME            string
	SERVER_CHARACTER_SET            int
	SERVER_FLAGS                    uint8
	CHARACTER_SET_ELEMENT           int
	ARRAY_LENGTH                    int
	NUMBER_ARRAY                    []byte
	SERVER_COMPILE_TIME_CAPS_LENGHT uint8
	SERVER_RUN_TIME_CAPS_LENGTH     uint8
	SERVER_COMPILE_TIME_CAPS        []byte
	SERVER_RUN_TIME_CAPS            []byte
}

type OracleDataTypeDataMessageRequest struct {
	SERVER_CHARACTER_SET     int
	SERVERN_CHARACTER_SET    int
	SERVER_FLAGS             uint8
	COMPILE_TIME_CAPS_LENGHT uint8
	RUN_TIME_CAPS_LENGTH     uint8
	COMPILE_TIME_CAPS        []byte
	RUN_TIME_CAPS            []byte
	CLIENT_TZ_VERSION        int
	DATA_TYPE                []DataTypeInfo
}

type OracleDataTypeDataMessageResponse struct {
	SERVER_TZ_VERSION int
	DATA_TYPE         []DataTypeInfo
}

type OracleFuntionTypeDataMessage struct {
	FUNCTION_CODE   FunctionType
	FUNCTION_DATA   interface{}
	SEQUENCE_NUMBER uint8
}

type OracleRowHeaderTypeDataMessage struct {
	FLAGS         uint8
	COLUMN_COUNT  int
	ITERATION_NUM int
	ROW_COUNT     int
	BUFFER_LENGTH int
	BIT_VECTOR []byte
}

type OracleAuthDataMessagePhaseOne struct {
	HAS_USER         uint8
	USER_BYTE_LENGTH uint32
	AUTH_MODE        AuthMode
	NUM_PAIR         uint32
	USER_BYTES       string
	AUTH_KEY_VALUE   []AuthKeyValue
}

type OracleParamterTypeAuthMessage struct {
}

type AuthKeyValue struct {
	KEY   string
	VALUE string
}

type DataTypeInfo struct {
	DATA_TYPE      DataType
	CONV_DATA_TYPE DataType
	REPRESENTATION DataType
}

type OracleQuery struct {
	OPTIONS              int
	CURSOR_ID            int
	IS_DDL               int
	SQL_LENGTH           int
	ROWS_TO_FETCH        int
	LOB                  int
	NUM_PARAMS           int
	NUM_DEFINES          int
	ARRAY_LENGTH         int
	PREFETCH_BUFFER_SIZE uint32
	REGISTRATION_ID      uint32
	NUM_EXEC             int
	SQL_BYTES            []byte
	Al8i4                []int
	IS_QUERY             int
	DML_OPTIONS          int
	DefineDataList       []DefineData
}

type DefineData struct {
	DataType             uint8
	Flag                 uint8
	Precision            uint8
	Scale                uint8
	MaxLen               int
	MaxNoOfArrayElements int
	ContFlag             int
	ToID                 []byte
	Version              int
	CharsetID            int
	CharsetForm          uint8
	MaxCharLen           int
	Oaccollid            int
}

type PiggyBackMsg struct {
}

type OracleQueryResponse struct {
	Query string
}

type OracleError struct {
	ErrCode int
	ErrMsg  string
}

type ControlType int
type ControlError int

const (
	TNS_CONTROL_TYPE_RESET_OOB           = ControlType(9)
	TNS_CONTROL_TYPE_INBAND_NOTIFICATION = ControlType(8)
)

const (
	TNS_ERR_SESSION_SHUTDOWN = ControlError(12572)
	TNS_ERR_INBAND_MESSAGE   = ControlError(12573)
)

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

type DataPacketType uint8

const (
	Default                           = DataPacketType(0)
	OracleProtocolDataMessageType     = DataPacketType(1)
	OracleDataTypeDataMessageType     = DataPacketType(2)
	OracleFunctionDataMesssageType    = DataPacketType(3)
	OracleRowHeaderDataMessage        = DataPacketType(6)
	OracleParamterTypeAuthMessageType = DataPacketType(102)
	OracleConnectionDataMessageType   = DataPacketType(100)
	OracleRedirectDataMessageType     = DataPacketType(101)
)

type DataType uint16

const (
	TNS_DATA_TYPE_VARCHAR        = DataType(1)
	TNS_DATA_TYPE_NUMBER         = DataType(2)
	TNS_DATA_TYPE_BINARY_INTEGER = DataType(3)
	TNS_DATA_TYPE_FLOAT          = DataType(4)
	TNS_DATA_TYPE_STR            = DataType(5)
	TNS_DATA_TYPE_VNU            = DataType(6)
	TNS_DATA_TYPE_PDN            = DataType(7)
	TNS_DATA_TYPE_LONG           = DataType(8)
	TNS_DATA_TYPE_VCS            = DataType(9)
	TNS_DATA_TYPE_TID            = DataType(10)
	TNS_DATA_TYPE_ROWID          = DataType(11)
	TNS_DATA_TYPE_DATE           = DataType(12)
	TNS_DATA_TYPE_VBI            = DataType(15)
	TNS_DATA_TYPE_RAW            = DataType(23)
	TNS_DATA_TYPE_LONG_RAW       = DataType(24)
	TNS_DATA_TYPE_UB2            = DataType(25)
	TNS_DATA_TYPE_UB4            = DataType(26)
	TNS_DATA_TYPE_SB1            = DataType(27)
	TNS_DATA_TYPE_SB2            = DataType(28)
	TNS_DATA_TYPE_SB4            = DataType(29)
	TNS_DATA_TYPE_SWORD          = DataType(30)
	TNS_DATA_TYPE_UWORD          = DataType(31)
	TNS_DATA_TYPE_PTRB           = DataType(32)
	TNS_DATA_TYPE_PTRW           = DataType(33)
	TNS_DATA_TYPE_OER8           = DataType(34 + 256)
	TNS_DATA_TYPE_FUN            = DataType(35 + 256)
	TNS_DATA_TYPE_AUA            = DataType(36 + 256)
	TNS_DATA_TYPE_RXH7           = DataType(37 + 256)
	TNS_DATA_TYPE_NA6            = DataType(38 + 256)
	TNS_DATA_TYPE_OAC            = DataType(39)
	TNS_DATA_TYPE_AMS            = DataType(40)
	TNS_DATA_TYPE_BRN            = DataType(41)
	TNS_DATA_TYPE_BRP            = DataType(42 + 256)
	TNS_DATA_TYPE_BRV            = DataType(43 + 256)
	TNS_DATA_TYPE_KVA            = DataType(44 + 256)
	TNS_DATA_TYPE_CLS            = DataType(45 + 256)
	TNS_DATA_TYPE_CUI            = DataType(46 + 256)
	TNS_DATA_TYPE_DFN            = DataType(47 + 256)
	TNS_DATA_TYPE_DQR            = DataType(48 + 256)
	TNS_DATA_TYPE_DSC            = DataType(49 + 256)
	TNS_DATA_TYPE_EXE            = DataType(50 + 256)
	TNS_DATA_TYPE_FCH            = DataType(51 + 256)
	TNS_DATA_TYPE_GBV            = DataType(52 + 256)
	TNS_DATA_TYPE_GEM            = DataType(53 + 256)
	TNS_DATA_TYPE_GIV            = DataType(54 + 256)
	TNS_DATA_TYPE_OKG            = DataType(55 + 256)
	TNS_DATA_TYPE_HMI            = DataType(56 + 256)
	TNS_DATA_TYPE_INO            = DataType(57 + 256)
	TNS_DATA_TYPE_LNF            = DataType(59 + 256)
	TNS_DATA_TYPE_ONT            = DataType(60 + 256)
	TNS_DATA_TYPE_OPE            = DataType(61 + 256)
	TNS_DATA_TYPE_OSQ            = DataType(62 + 256)
	TNS_DATA_TYPE_SFE            = DataType(63 + 256)
	TNS_DATA_TYPE_SPF            = DataType(64 + 256)
	TNS_DATA_TYPE_VSN            = DataType(65 + 256)
	TNS_DATA_TYPE_UD7            = DataType(66 + 256)
	TNS_DATA_TYPE_DSA            = DataType(67 + 256)
	TNS_DATA_TYPE_UIN            = DataType(68)
	TNS_DATA_TYPE_PIN            = DataType(71 + 256)
	TNS_DATA_TYPE_PFN            = DataType(72 + 256)
	TNS_DATA_TYPE_PPT            = DataType(73 + 256)
	TNS_DATA_TYPE_STO            = DataType(75 + 256)
	TNS_DATA_TYPE_ARC            = DataType(77 + 256)
	TNS_DATA_TYPE_MRS            = DataType(78 + 256)
	TNS_DATA_TYPE_MRT            = DataType(79 + 256)
	TNS_DATA_TYPE_MRG            = DataType(80 + 256)
	TNS_DATA_TYPE_MRR            = DataType(81 + 256)
	TNS_DATA_TYPE_MRC            = DataType(82 + 256)
	TNS_DATA_TYPE_VER            = DataType(83 + 256)
	TNS_DATA_TYPE_LON2           = DataType(84 + 256)
	TNS_DATA_TYPE_INO2           = DataType(85 + 256)
	TNS_DATA_TYPE_ALL            = DataType(86 + 256)
	TNS_DATA_TYPE_UDB            = DataType(87 + 256)
	TNS_DATA_TYPE_AQI            = DataType(88 + 256)
	TNS_DATA_TYPE_ULB            = DataType(89 + 256)
	TNS_DATA_TYPE_ULD            = DataType(90 + 256)
	TNS_DATA_TYPE_SLS            = DataType(91)
	TNS_DATA_TYPE_SID            = DataType(92 + 256)
	TNS_DATA_TYPE_NA7            = DataType(93 + 256)
	TNS_DATA_TYPE_LVC            = DataType(94)
	TNS_DATA_TYPE_LVB            = DataType(95)
	TNS_DATA_TYPE_CHAR           = DataType(96)
	TNS_DATA_TYPE_AVC            = DataType(97)
	TNS_DATA_TYPE_AL7            = DataType(98 + 256)
	TNS_DATA_TYPE_K2RPC          = DataType(99 + 256)
	TNS_DATA_TYPE_BINARY_FLOAT   = DataType(100)
	TNS_DATA_TYPE_BINARY_DOUBLE  = DataType(101)
	TNS_DATA_TYPE_CURSOR         = DataType(102)
	TNS_DATA_TYPE_RDD            = DataType(104)
	TNS_DATA_TYPE_XDP            = DataType(103 + 256)
	TNS_DATA_TYPE_OSL            = DataType(106)
	TNS_DATA_TYPE_OKO8           = DataType(107 + 256)
	TNS_DATA_TYPE_EXT_NAMED      = DataType(108)
	TNS_DATA_TYPE_INT_NAMED      = DataType(109)
	TNS_DATA_TYPE_EXT_REF        = DataType(110)
	TNS_DATA_TYPE_INT_REF        = DataType(111)
	TNS_DATA_TYPE_CLOB           = DataType(112)
	TNS_DATA_TYPE_BLOB           = DataType(113)
	TNS_DATA_TYPE_BFILE          = DataType(114)
	TNS_DATA_TYPE_CFILE          = DataType(115)
	TNS_DATA_TYPE_RSET           = DataType(116)
	TNS_DATA_TYPE_CWD            = DataType(117)
	TNS_DATA_TYPE_JSON           = DataType(119)
	TNS_DATA_TYPE_NEW_OAC        = DataType(120)
	TNS_DATA_TYPE_UD12           = DataType(124 + 256)
	TNS_DATA_TYPE_AL8            = DataType(125 + 256)
	TNS_DATA_TYPE_LFOP           = DataType(126 + 256)
	TNS_DATA_TYPE_FCRT           = DataType(127 + 256)
	TNS_DATA_TYPE_DNY            = DataType(128 + 256)
	TNS_DATA_TYPE_OPR            = DataType(129 + 256)
	TNS_DATA_TYPE_PLS            = DataType(130 + 256)
	TNS_DATA_TYPE_XID            = DataType(131 + 256)
	TNS_DATA_TYPE_TXN            = DataType(132 + 256)
	TNS_DATA_TYPE_DCB            = DataType(133 + 256)
	TNS_DATA_TYPE_CCA            = DataType(134 + 256)
	TNS_DATA_TYPE_WRN            = DataType(135 + 256)
	TNS_DATA_TYPE_TLH            = DataType(137 + 256)
	TNS_DATA_TYPE_TOH            = DataType(138 + 256)
	TNS_DATA_TYPE_FOI            = DataType(139 + 256)
	TNS_DATA_TYPE_SID2           = DataType(140 + 256)
	TNS_DATA_TYPE_TCH            = DataType(141 + 256)
	TNS_DATA_TYPE_PII            = DataType(142 + 256)
	TNS_DATA_TYPE_PFI            = DataType(143 + 256)
	TNS_DATA_TYPE_PPU            = DataType(144 + 256)
	TNS_DATA_TYPE_PTE            = DataType(145 + 256)
	TNS_DATA_TYPE_CLV            = DataType(146)
	TNS_DATA_TYPE_RXH8           = DataType(148 + 256)
	TNS_DATA_TYPE_N12            = DataType(149 + 256)
	TNS_DATA_TYPE_AUTH           = DataType(150 + 256)
	TNS_DATA_TYPE_KVAL           = DataType(151 + 256)
	TNS_DATA_TYPE_DTR            = DataType(152)
	TNS_DATA_TYPE_DUN            = DataType(153)
	TNS_DATA_TYPE_DOP            = DataType(154)
	TNS_DATA_TYPE_VST            = DataType(155)
	TNS_DATA_TYPE_ODT            = DataType(156)
	TNS_DATA_TYPE_FGI            = DataType(157 + 256)
	TNS_DATA_TYPE_DSY            = DataType(158 + 256)
	TNS_DATA_TYPE_DSYR8          = DataType(159 + 256)
	TNS_DATA_TYPE_DSYH8          = DataType(160 + 256)
	TNS_DATA_TYPE_DSYL           = DataType(161 + 256)
	TNS_DATA_TYPE_DSYT8          = DataType(162 + 256)
	TNS_DATA_TYPE_DSYV8          = DataType(163 + 256)
	TNS_DATA_TYPE_DSYP           = DataType(164 + 256)
	TNS_DATA_TYPE_DSYF           = DataType(165 + 256)
	TNS_DATA_TYPE_DSYK           = DataType(166 + 256)
	TNS_DATA_TYPE_DSYY           = DataType(167 + 256)
	TNS_DATA_TYPE_DSYQ           = DataType(168 + 256)
	TNS_DATA_TYPE_DSYC           = DataType(169 + 256)
	TNS_DATA_TYPE_DSYA           = DataType(170 + 256)
	TNS_DATA_TYPE_OT8            = DataType(171 + 256)
	TNS_DATA_TYPE_DOL            = DataType(172)
	TNS_DATA_TYPE_DSYTY          = DataType(173 + 256)
	TNS_DATA_TYPE_AQE            = DataType(174 + 256)
	TNS_DATA_TYPE_KV             = DataType(175 + 256)
	TNS_DATA_TYPE_AQD            = DataType(176 + 256)
	TNS_DATA_TYPE_AQ8            = DataType(177 + 256)
	TNS_DATA_TYPE_TIME           = DataType(178)
	TNS_DATA_TYPE_TIME_TZ        = DataType(179)
	TNS_DATA_TYPE_TIMESTAMP      = DataType(180)
	TNS_DATA_TYPE_TIMESTAMP_TZ   = DataType(181)
	TNS_DATA_TYPE_INTERVAL_YM    = DataType(182)
	TNS_DATA_TYPE_INTERVAL_DS    = DataType(183)
	TNS_DATA_TYPE_EDATE          = DataType(184)
	TNS_DATA_TYPE_ETIME          = DataType(185)
	TNS_DATA_TYPE_ETTZ           = DataType(186)
	TNS_DATA_TYPE_ESTAMP         = DataType(187)
	TNS_DATA_TYPE_ESTZ           = DataType(188)
	TNS_DATA_TYPE_EIYM           = DataType(189)
	TNS_DATA_TYPE_EIDS           = DataType(190)
	TNS_DATA_TYPE_RFS            = DataType(193 + 256)
	TNS_DATA_TYPE_RXH10          = DataType(194 + 256)
	TNS_DATA_TYPE_DCLOB          = DataType(195)
	TNS_DATA_TYPE_DBLOB          = DataType(196)
	TNS_DATA_TYPE_DBFILE         = DataType(197)
	TNS_DATA_TYPE_DJSON          = DataType(198)
	TNS_DATA_TYPE_KPN            = DataType(198 + 256)
	TNS_DATA_TYPE_KPDNR          = DataType(199 + 256)
	TNS_DATA_TYPE_DSYD           = DataType(200 + 256)
	TNS_DATA_TYPE_DSYS           = DataType(201 + 256)
	TNS_DATA_TYPE_DSYR           = DataType(202 + 256)
	TNS_DATA_TYPE_DSYH           = DataType(203 + 256)
	TNS_DATA_TYPE_DSYT           = DataType(204 + 256)
	TNS_DATA_TYPE_DSYV           = DataType(205 + 256)
	TNS_DATA_TYPE_AQM            = DataType(206 + 256)
	TNS_DATA_TYPE_OER11          = DataType(207 + 256)
	TNS_DATA_TYPE_UROWID         = DataType(208)
	TNS_DATA_TYPE_AQL            = DataType(210 + 256)
	TNS_DATA_TYPE_OTC            = DataType(211 + 256)
	TNS_DATA_TYPE_KFNO           = DataType(212 + 256)
	TNS_DATA_TYPE_KFNP           = DataType(213 + 256)
	TNS_DATA_TYPE_KGT8           = DataType(214 + 256)
	TNS_DATA_TYPE_RASB4          = DataType(215 + 256)
	TNS_DATA_TYPE_RAUB2          = DataType(216 + 256)
	TNS_DATA_TYPE_RAUB1          = DataType(217 + 256)
	TNS_DATA_TYPE_RATXT          = DataType(218 + 256)
	TNS_DATA_TYPE_RSSB4          = DataType(219 + 256)
	TNS_DATA_TYPE_RSUB2          = DataType(220 + 256)
	TNS_DATA_TYPE_RSUB1          = DataType(221 + 256)
	TNS_DATA_TYPE_RSTXT          = DataType(222 + 256)
	TNS_DATA_TYPE_RIDL           = DataType(223 + 256)
	TNS_DATA_TYPE_GLRDD          = DataType(224 + 256)
	TNS_DATA_TYPE_GLRDG          = DataType(225 + 256)
	TNS_DATA_TYPE_GLRDC          = DataType(226 + 256)
	TNS_DATA_TYPE_OKO            = DataType(227 + 256)
	TNS_DATA_TYPE_DPP            = DataType(228 + 256)
	TNS_DATA_TYPE_DPLS           = DataType(229 + 256)
	TNS_DATA_TYPE_DPMOP          = DataType(230 + 256)
	TNS_DATA_TYPE_TIMESTAMP_LTZ  = DataType(231)
	TNS_DATA_TYPE_ESITZ          = DataType(232)
	TNS_DATA_TYPE_UB8            = DataType(233)
	TNS_DATA_TYPE_STAT           = DataType(234 + 256)
	TNS_DATA_TYPE_RFX            = DataType(235 + 256)
	TNS_DATA_TYPE_FAL            = DataType(236 + 256)
	TNS_DATA_TYPE_CKV            = DataType(237 + 256)
	TNS_DATA_TYPE_DRCX           = DataType(238 + 256)
	TNS_DATA_TYPE_KGH            = DataType(239 + 256)
	TNS_DATA_TYPE_AQO            = DataType(240 + 256)
	TNS_DATA_TYPE_PNTY           = DataType(241)
	TNS_DATA_TYPE_OKGT           = DataType(242 + 256)
	TNS_DATA_TYPE_KPFC           = DataType(243 + 256)
	TNS_DATA_TYPE_FE2            = DataType(244 + 256)
	TNS_DATA_TYPE_SPFP           = DataType(245 + 256)
	TNS_DATA_TYPE_DPULS          = DataType(246 + 256)
	TNS_DATA_TYPE_BOOLEAN        = DataType(252)
	TNS_DATA_TYPE_AQA            = DataType(253 + 256)
	TNS_DATA_TYPE_KPBF           = DataType(254 + 256)
	TNS_DATA_TYPE_TSM            = DataType(513)
	TNS_DATA_TYPE_MSS            = DataType(514)
	TNS_DATA_TYPE_KPC            = DataType(516)
	TNS_DATA_TYPE_CRS            = DataType(517)
	TNS_DATA_TYPE_KKS            = DataType(518)
	TNS_DATA_TYPE_KSP            = DataType(519)
	TNS_DATA_TYPE_KSPTOP         = DataType(520)
	TNS_DATA_TYPE_KSPVAL         = DataType(521)
	TNS_DATA_TYPE_PSS            = DataType(522)
	TNS_DATA_TYPE_NLS            = DataType(523)
	TNS_DATA_TYPE_ALS            = DataType(524)
	TNS_DATA_TYPE_KSDEVTVAL      = DataType(525)
	TNS_DATA_TYPE_KSDEVTTOP      = DataType(526)
	TNS_DATA_TYPE_KPSPP          = DataType(527)
	TNS_DATA_TYPE_KOL            = DataType(528)
	TNS_DATA_TYPE_LST            = DataType(529)
	TNS_DATA_TYPE_ACX            = DataType(530)
	TNS_DATA_TYPE_SCS            = DataType(531)
	TNS_DATA_TYPE_RXH            = DataType(532)
	TNS_DATA_TYPE_KPDNS          = DataType(533)
	TNS_DATA_TYPE_KPDCN          = DataType(534)
	TNS_DATA_TYPE_KPNNS          = DataType(535)
	TNS_DATA_TYPE_KPNCN          = DataType(536)
	TNS_DATA_TYPE_KPS            = DataType(537)
	TNS_DATA_TYPE_APINF          = DataType(538)
	TNS_DATA_TYPE_TEN            = DataType(539)
	TNS_DATA_TYPE_XSSCS          = DataType(540)
	TNS_DATA_TYPE_XSSSO          = DataType(541)
	TNS_DATA_TYPE_XSSAO          = DataType(542)
	TNS_DATA_TYPE_KSRPC          = DataType(543)
	TNS_DATA_TYPE_KVL            = DataType(560)
	TNS_DATA_TYPE_SESSGET        = DataType(563)
	TNS_DATA_TYPE_SESSREL        = DataType(564)
	TNS_DATA_TYPE_XSS            = DataType(565)
	TNS_DATA_TYPE_PDQCINV        = DataType(572)
	TNS_DATA_TYPE_PDQIDC         = DataType(573)
	TNS_DATA_TYPE_KPDQCSTA       = DataType(574)
	TNS_DATA_TYPE_KPRS           = DataType(575)
	TNS_DATA_TYPE_KPDQIDC        = DataType(576)
	TNS_DATA_TYPE_RTSTRM         = DataType(578)
	TNS_DATA_TYPE_SESSRET        = DataType(579)
	TNS_DATA_TYPE_SCN6           = DataType(580)
	TNS_DATA_TYPE_KECPA          = DataType(581)
	TNS_DATA_TYPE_KECPP          = DataType(582)
	TNS_DATA_TYPE_SXA            = DataType(583)
	TNS_DATA_TYPE_KVARR          = DataType(584)
	TNS_DATA_TYPE_KPNGN          = DataType(585)
	TNS_DATA_TYPE_XSNSOP         = DataType(590)
	TNS_DATA_TYPE_XSATTR         = DataType(591)
	TNS_DATA_TYPE_XSNS           = DataType(592)
	TNS_DATA_TYPE_TXT            = DataType(593)
	TNS_DATA_TYPE_XSSESSNS       = DataType(594)
	TNS_DATA_TYPE_XSATTOP        = DataType(595)
	TNS_DATA_TYPE_XSCREOP        = DataType(596)
	TNS_DATA_TYPE_XSDETOP        = DataType(597)
	TNS_DATA_TYPE_XSDESOP        = DataType(598)
	TNS_DATA_TYPE_XSSETSP        = DataType(599)
	TNS_DATA_TYPE_XSSIDP         = DataType(600)
	TNS_DATA_TYPE_XSPRIN         = DataType(601)
	TNS_DATA_TYPE_XSKVL          = DataType(602)
	TNS_DATA_TYPE_XSSS2          = DataType(603)
	TNS_DATA_TYPE_XSNSOP2        = DataType(604)
	TNS_DATA_TYPE_XSNS2          = DataType(605)
	TNS_DATA_TYPE_IMPLRES        = DataType(611)
	TNS_DATA_TYPE_OER            = DataType(612)
	TNS_DATA_TYPE_UB1ARRAY       = DataType(613)
	TNS_DATA_TYPE_SESSSTATE      = DataType(614)
	TNS_DATA_TYPE_AC_REPLAY      = DataType(615)
	TNS_DATA_TYPE_AC_CONT        = DataType(616)
	TNS_DATA_TYPE_KPDNREQ        = DataType(622)
	TNS_DATA_TYPE_KPDNRNF        = DataType(623)
	TNS_DATA_TYPE_KPNGNC         = DataType(624)
	TNS_DATA_TYPE_KPNRI          = DataType(625)
	TNS_DATA_TYPE_AQENQ          = DataType(626)
	TNS_DATA_TYPE_AQDEQ          = DataType(627)
	TNS_DATA_TYPE_AQJMS          = DataType(628)
	TNS_DATA_TYPE_KPDNRPAY       = DataType(629)
	TNS_DATA_TYPE_KPDNRACK       = DataType(630)
	TNS_DATA_TYPE_KPDNRMP        = DataType(631)
	TNS_DATA_TYPE_KPDNRDQ        = DataType(632)
	TNS_DATA_TYPE_CHUNKINFO      = DataType(636)
	TNS_DATA_TYPE_SCN            = DataType(637)
	TNS_DATA_TYPE_SCN8           = DataType(638)
	TNS_DATA_TYPE_UDS            = DataType(639)
	TNS_DATA_TYPE_TNP            = DataType(640)
)

type FunctionType uint8

const (
	TNS_FUNC_AUTH_PHASE_ONE      = FunctionType(118)
	TNS_FUNC_AUTH_PHASE_TWO      = FunctionType(115)
	TNS_FUNC_CLOSE_CURSORS       = FunctionType(105)
	TNS_FUNC_COMMIT              = FunctionType(14)
	TNS_FUNC_EXECUTE             = FunctionType(94)
	TNS_FUNC_FETCH               = FunctionType(5)
	TNS_FUNC_LOB_OP              = FunctionType(96)
	TNS_FUNC_LOGOFF              = FunctionType(9)
	TNS_FUNC_PING                = FunctionType(147)
	TNS_FUNC_ROLLBACK            = FunctionType(15)
	TNS_FUNC_SET_END_TO_END_ATTR = FunctionType(135)
	TNS_FUNC_REEXECUTE           = FunctionType(4)
	TNS_FUNC_REEXECUTE_AND_FETCH = FunctionType(78)
	TNS_FUNC_SESSION_GET         = FunctionType(162)
	TNS_FUNC_SESSION_RELEASE     = FunctionType(163)
	TNS_FUNC_SET_SCHEMA          = FunctionType(152)
)

type AuthMode uint32

const (
	TNS_AUTH_MODE_LOGON           = AuthMode(0x00000001)
	TNS_AUTH_MODE_CHANGE_PASSWORD = AuthMode(0x00000002)
	TNS_AUTH_MODE_SYSDBA          = AuthMode(0x00000020)
	TNS_AUTH_MODE_SYSOPER         = AuthMode(0x00000040)
	TNS_AUTH_MODE_PRELIM          = AuthMode(0x00000080)
	TNS_AUTH_MODE_WITH_PASSWORD   = AuthMode(0x00000100)
	TNS_AUTH_MODE_SYSASM          = AuthMode(0x00400000)
	TNS_AUTH_MODE_SYSBKP          = AuthMode(0x01000000)
	TNS_AUTH_MODE_SYSDGD          = AuthMode(0x02000000)
	TNS_AUTH_MODE_SYSKMT          = AuthMode(0x04000000)
	TNS_AUTH_MODE_SYSRAC          = AuthMode(0x08000000)
	TNS_AUTH_MODE_IAM_TOKEN       = AuthMode(0x20000000)
)

package util

const MaxUint = 0xffffffff
const MinUint = 0
const MaxInt = MaxUint >> 1
const MinInt = -MaxInt - 1

const SecondsOfDay = 24 * 60 * 60

const (
	BlockHashLength = 64
)

// error code
const (
	ErrorCodeUnknown           = -1
	ErrorCodeSuccess           = 0
	ErrorCodeInvalidParams     = 1
	ErrorCodeInvalidParamRange = 2
	ErrorCodeNotExist          = 3
	ErrorCodeInternalError     = 100
)

//// error message
//var ErrorMessages = map[int]string{
//	ErrorCodeUnknown:           "unknown",
//	ErrorCodeSuccess:           "success",
//	ErrorCodeInvalidParams:     "invalid parameters",
//	ErrorCodeInvalidParamRange: "invalid parameter range",
//	ErrorCodeNotExist:          "not exist",
//	ErrorCodeInternalError:     "internal error",
//}

// websocket message type
const (
	MessageTypeInfo  = 1 // node info
	MessageTypeBlock = 2 // block
	MessageTypeTx    = 3 // tx
)

// default quantity
const (
	DefaultLatestBlocks = 10 // default latest blocks
	DefaultLatestTxs    = 10 // default latest txs

	DefaultMaxLatestTxCounts      = 90 // 3 months
	DefaultMaxLatestAccountCounts = 90 // 3 months

	DefaultAccountPageSize = 50
	DefaultMaxAccountPage  = 300
	DefaultBlockPageSize   = 50
	DefaultMaxBlockPage    = 300
	DefaultTxPageSize      = 50
	DefaultMaxTxPage       = 300

	DefaultDisableWebSocketBlockHeight = 10
)

// default duration for cache, unit: second
const (
	// Global Information
	DurationInfo = 60

	// Account
	DurationAccount = 5

	// Block
	DurationLatestBlocks = 5
	DurationBlock        = 30
	DurationBlockCount   = 5

	// Tx
	DurationLatestTxs = 5
	DurationTx        = 30
	DurationTxCount   = 5

	// TxCount
	DurationLatestTxCounts      = 600
	DurationLatestAccountCounts = 600

	// TxType
	DurationTxTypes = 600
)

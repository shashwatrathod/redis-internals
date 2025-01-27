package resp

const (
	RespSimpleStringIdentifier byte = '+'
	RespSimpleErrorIdentifier  byte = '-'
	RespIntegerIdentifier      byte = ':'
	RespBulkStringIdentifier   byte = '$'
	RespArrayIdentifier        byte = '*'
)

type RespDataTypes int

const (
	SimpleString RespDataTypes = iota
	BulkString
	SimpleError
	RespArray
	RespInteger
)

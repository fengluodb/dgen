package parser

import "strconv"

type tokenType int

const (
	T_Builtin       tokenType = iota // 内置类型，如int32, uint32等
	T_Enum                           // enum关键字
	T_Message                        // message关键字
	T_Service                        // service关键字
	T_Seq                            // seq关键字
	T_Optional                       // option关键字
	T_Return                         // return关键字
	T_Identifier                     // 标识符
	T_Num                            // 数字
	T_Assign                         // =
	T_LSmallBracket                  // (
	T_RSmallBracket                  // )
	T_LBracket                       // [
	T_RBracket                       // ]
	T_LCurlyBracket                  // {
	T_RCurlyBracket                  // }
	T_Comma                          // ,
	T_Semicolon                      // ;
)

type token struct {
	typ    tokenType
	val    string
	row    int
	column int
}

var tokenTypeMap map[string]tokenType
var ignoreCharMap map[byte]struct{}

func init() {
	tokenTypeMap = make(map[string]tokenType)
	ignoreCharMap = make(map[byte]struct{})

	tokenTypeMap["uint8"] = T_Builtin
	tokenTypeMap["uint16"] = T_Builtin
	tokenTypeMap["uint32"] = T_Builtin
	tokenTypeMap["uint64"] = T_Builtin
	tokenTypeMap["int8"] = T_Builtin
	tokenTypeMap["int16"] = T_Builtin
	tokenTypeMap["int32"] = T_Builtin
	tokenTypeMap["int64"] = T_Builtin
	tokenTypeMap["float32"] = T_Builtin
	tokenTypeMap["float64"] = T_Builtin
	tokenTypeMap["string"] = T_Builtin
	tokenTypeMap["list"] = T_Builtin
	tokenTypeMap["map"] = T_Builtin
	tokenTypeMap["enum"] = T_Enum
	tokenTypeMap["message"] = T_Message
	tokenTypeMap["service"] = T_Service
	tokenTypeMap["seq"] = T_Seq
	tokenTypeMap["optional"] = T_Optional
	tokenTypeMap["return"] = T_Return
	tokenTypeMap["="] = T_Assign
	tokenTypeMap["("] = T_LSmallBracket
	tokenTypeMap[")"] = T_RSmallBracket
	tokenTypeMap["["] = T_LBracket
	tokenTypeMap["]"] = T_RBracket
	tokenTypeMap["{"] = T_LCurlyBracket
	tokenTypeMap["}"] = T_RCurlyBracket
	tokenTypeMap[","] = T_Comma
	tokenTypeMap[";"] = T_Semicolon

	ignoreCharMap[' '] = struct{}{}
	ignoreCharMap['\t'] = struct{}{}
	ignoreCharMap['\n'] = struct{}{}
	ignoreCharMap['#'] = struct{}{}
}

func isNum(s string) bool {
	_, err := strconv.Atoi(s)

	return err == nil
}

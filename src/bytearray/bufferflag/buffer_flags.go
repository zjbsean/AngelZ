package bufferflag

import "sync"

//ByteArrayFlag : data in ByteArray Flag
//  FlagType:
//      Raw / HadEncrypt / HadCompress
type BufferFlag byte

const (
	_          BufferFlag = 0x00
	HadEncrypt            = 0x01
	HadCompres            = 0x02

	CARRIER  = 0x10
	PROTOCOL = 0x20
)

type BufferFlagHelper struct{}

var helpInstance *BufferFlagHelper
var once sync.Once

func GetInstance() *BufferFlagHelper {
	once.Do(func() {
		helpInstance = &BufferFlagHelper{}
	})
	return helpInstance
}

func (*BufferFlagHelper) TestFlag(beTestFlag, flagBit BufferFlag) bool {
	return (beTestFlag & flagBit) != 0
}

func (*BufferFlagHelper) SetFlag(theFlag, flagBit BufferFlag) BufferFlag {
	return theFlag | flagBit
}

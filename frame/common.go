package frame

import (
	"errors"
	"github.com/nicle-lin/golang-csp/stream"
	"log"
)

/*
 * 格式:[2byte][2byte][1byte][data]
 *      前两2个是协议开头标志　　后面2个是标识数据(data)长度(仅包括data的长度)　接着的1个是操作类型　　data为真实数据
 */

type OpType byte //操作类型

const (
	//客户端用于向服务端查询 增加 删除信息的操作类型
	GET    OpType = 1 //0000 0001
	ADD    OpType = 2 //0000 0010
	DELETE OpType = 4 //0000 0100

	//服务端用于返回操作是否成功
	SUCCESS OpType = 8  //0000 1000
	FAIL    OpType = 16 //0001 0000

	FrameHeaderLen = 5
	MaxFrameLen    = 1024   //随便定一个长度,这长度包括FrameHeaderLen及data
	FrameFlag      = 0xf3db //协议开头标志
)

type Frame interface {
	//根据操作类型返回可用于发送的[]byte类型的数据
	BuildFrame() (frame []byte, err error)
	SetFrameData(data []byte)
	SetFrameOpType(optype OpType)
}

type frame struct {
	data   []byte
	optype OpType
}

func NewFrame(data []byte, optype OpType) Frame {
	return &frame{
		data:   data,
		optype: optype,
	}
}

func (f *frame) SetFrameData(data []byte) {
	f.data = data
}

func (f *frame) SetFrameOpType(optype OpType) {
	f.optype = optype
}

/*
 * 格式:[2byte][2byte][1byte][data]
 *      前两2个是协议开头标志　　后面2个是标识数据(data)长度　接着的1个是操作类型　　data为真实数据
 */
func (f *frame) BuildFrame() (frame []byte, err error) {
	//
	buf := make([]byte, FrameHeaderLen+len(f.data))

	st := stream.NewLEStream(buf)
	err = st.WriteUint16(FrameFlag)
	checkError(err)
	err = st.WriteUint16(uint16(len(f.data)))
	checkError(err)
	err = st.WriteByte(byte(f.optype))
	checkError(err)
	err = st.WriteBuff(f.data)
	checkError(err)
	frame = st.Data()
	return
}

//解析协议头部是否合法,如果合法返回接下来要接收的数据长度,不合法就返回0
//传入参数必须是头部5个字节的长度
func ParseFrame(frameHeader []byte) (length uint16, err error) {
	if len(frameHeader) != FrameHeaderLen {
		err = errors.New("frameHeader isn't equal FrameHeaderLen(5)")
		return
	}
	st := stream.NewLEStream(frameHeader)
	flag, _ := st.ReadUint16()
	if flag != FrameFlag {
		err = errors.New("flag isn't equal 0xf3db")
		return
	}

	length, _ = st.ReadUint16()
	if length+FrameHeaderLen > MaxFrameLen {
		err = errors.New("frame length over MaxFrameLen(1024byte)")
		return
	}

	op, _ := st.ReadByte()
	switch OpType(op) {
	case ADD:
	case GET:
	case DELETE:
	case SUCCESS:
	case FAIL:
	case ADD | SUCCESS:
	case DELETE | SUCCESS:
	case GET | SUCCESS:
	case ADD | FAIL:
	case DELETE | FAIL:
	case GET | FAIL:
	default: //所有都没匹配中说明不是合法的操作
		log.Printf("operation:%d", op)
		err = errors.New("operation is illegal")
	}
	return
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)

	}
}

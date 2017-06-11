package session

import (
	"github.com/nicle-lin/golang-csp/frame"
	"github.com/nicle-lin/golang-csp/stream"
	"io"
	"net"
)

type Session interface {
	//从网络中读取数据,读回来是一个完整协议的,不完整就返回失败
	//连接被对端关闭了就返回一个io.ErrUnexpectedEOF
	Read() (st stream.Stream, err error)
	//向网络中写入一段完整协议的数据
	//连接被对端关闭了就返回一个io.EOF
	Write(f frame.Frame) (n int, err error)
	//关闭conn
	Close()
}

type session struct {
	conn net.Conn
}

func NewSession(conn net.Conn) Session {
	return &session{conn: conn}
}

func (s *session) Read() (st stream.Stream, err error) {
	//分配frame头部长度的大小的空间
	frameHeader := make([]byte, frame.FrameHeaderLen)
	if _, err = io.ReadFull(s.conn, frameHeader); err != nil {
		return
	}
	var dataLength uint16
	//校验收到的头部是否正确
	dataLength, err = frame.ParseFrame(frameHeader)
	if err != nil {
		return
	}

	//读取data
	data := make([]byte, dataLength)
	if _, err = io.ReadFull(s.conn, data); err != nil {
		return
	}

	//TODO: maybe there anther more effective,such as bytes.join
	newData := make([]byte, frame.FrameHeaderLen+dataLength)
	copy(newData, frameHeader)
	copy(newData[frame.FrameHeaderLen:], data)
	st = stream.NewLEStream(newData)
	return
}

func (s *session) Write(f frame.Frame) (n int, err error) {
	frame, err1 := f.BuildFrame()
	if err1 != nil {
		err = err1
		return
	}
	return s.conn.Write(frame)
}

func (s *session) Close() {
	s.conn.Close()
}

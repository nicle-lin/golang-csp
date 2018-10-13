package main

import (
	"github.com/go-ini/ini"
	"github.com/golang/protobuf/proto"
	"github.com/nicle-lin/golang-csp/frame"
	"github.com/nicle-lin/golang-csp/session"
	"io"
	"log"
	"math/rand"
	"net"
	"os"
	"time"
)

var conf = "../../conf/conf.ini"
var network, address string = "tcp", "127.0.0.1:8989"
var optypes = []frame.OpType{
	frame.GET,
	frame.ADD,
	frame.DELETE,
}

func init() {
	//一开始就读文件,没有回锁也不会有竞争,因为压根没人来竞争
	conf := "../../conf/conf.ini"
	cfg, err := ini.Load(conf)
	if err != nil {
		log.Printf("can't get the network address from %s,\n error:%s", conf, err)
		log.Printf("set the network address:%s", address)
	} else {
		address = cfg.Section("Host").Key("host").String()
	}

}
func RandomNumber(number int32) int32 {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Int31n(number)
}

func RandomString(length int) string {
	bytes := []byte{'a','b','c','d','e','f','g','h','i',
					'j','k','l','m','n','o','p','q','r',
					's','t','u','v','w','x','y','z'}
	var result []byte
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

func write(optype frame.OpType, data []byte, se session.Session) {
	f := frame.NewFrame(data, optype)
	_, err := se.Write(f)
	if err != nil {
		log.Println("wirte error:", err)
	}
}

func start() {
	log.Println("trying to connect")
	conn, err := net.Dial(network, address)
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("connecting successfully(%s)", conn.RemoteAddr())
	defer conn.Close()

	//for Get
	newRequest := &frame.Request{
		Id:      proto.Int32(RandomNumber(10)),
		Message: proto.String("Show me Person"),
	}
	dataGet, err := proto.Marshal(newRequest)
	if err != nil {
		log.Println(err)
	}

	//for Delete
	newRequest.Id = proto.Int32(RandomNumber(10))
	newRequest.Message = proto.String("delete Person")
	dataDelete, err := proto.Marshal(newRequest)
	if err != nil {
		log.Println(err)
	}

	//for Add
	newPerson := &frame.Person{
		Id:   proto.Int32(RandomNumber(10)),
		Age:  proto.Int32(RandomNumber(100)),
		Name: proto.String(RandomString(5)),
		City: proto.String(RandomString(5)),
	}
	dataAdd, err := proto.Marshal(newPerson)
	if err != nil {
		log.Println(err)
	}
	//用于判断Write Read是否完成了
	done := make(chan struct{}, 2)
	se := session.NewSession(conn)

	//Write
	go func() {
		for _, op := range optypes {
			switch op {
			case frame.ADD:
				write(op, dataAdd, se)
			case frame.DELETE:
				write(op, dataDelete, se)
			default: //默认为GET
				write(op, dataGet, se)
			}
		}
		done <- struct{}{}
	}()
	//Read
	go func() {
		//因为只写了3次,所有只读3即可,然后关闭conn
		for i := 0; i < 3; i++ {
			st, err := se.Read()
			if err == io.ErrUnexpectedEOF {
				log.Println("connection has been closed by peer")
				break
			} else if err != nil {
				log.Println(err)
				break
			}

			st.ReadUint16()
			st.ReadUint16()
			optype, _ := st.ReadByte()
			op := frame.OpType(optype)
			PrintResult(op, st.DataSelect(st.Pos(), st.Size()))

		}
		done <- struct{}{}
		se.Close()
	}()
	<-done
	<-done
}

func PrintResult(op frame.OpType, data []byte) {
	var operate string = ""

	switch { //此处 switch op{}是不正确的,因为判断应交给下面的case语句
	case (op & frame.ADD) == frame.ADD:
		operate = "Add"
	case op&frame.DELETE == frame.DELETE:
		operate = "Delete"
	case op&frame.GET == frame.GET:
		operate = "Get"
		newPerson := &frame.Person{}
		err := proto.Unmarshal(data, newPerson)
		if err != nil {
			log.Println(err)
		}
		log.Printf("Get the Person:\nID:%d\nName:%s\nAge:%d\nCity:%s\n",
			newPerson.GetId(), newPerson.GetName(), newPerson.GetAge(), newPerson.GetCity())

	default:
		log.Printf("unknow operate:%+v\n", op)
	}
	switch {
	case op&frame.SUCCESS == frame.SUCCESS:
		log.Printf("execute operate %s success\n", operate)
	case op&frame.FAIL == frame.FAIL:
		log.Printf("execute operate %s fail\n", operate)
	default:
		log.Printf("unknow operate:%+v\n", op)
	}

}

func main() {

	go func() {
		select {
		case <-time.NewTimer(time.Second * 8).C:
			log.Println("time out, exit")
			os.Exit(1)
		}
	}()

	log.Println("start to connect")
	start()
	log.Println("All operation done")

}

package common

import (
	"errors"
	"github.com/go-ini/ini"
	"github.com/golang/protobuf/proto"
	"github.com/nicle-lin/golang-csp/frame"
	"github.com/nicle-lin/golang-csp/stream"
	"log"
	"strconv"
	"sync"
)

var conf = "../../conf/conf.ini"

type Result struct {
	person  frame.Person //执行失败就为空,当optype为ADD时,解析结果存在此
	request frame.Request
	optype  frame.OpType //命令执行成功与否已经包含于此
}

func ParseProtobuf(optype frame.OpType, data []byte) (newFrame Result, err error) {
	if optype == frame.ADD {
		err = proto.Unmarshal(data, &newFrame.person)
	} else if optype == frame.GET || optype == frame.DELETE {
		err = proto.Unmarshal(data, &newFrame.request)
	} else {
		log.Printf("opType is illegal:%d", optype)
		err = errors.New("opType is illegal")
		return
	}
	newFrame.optype = optype
	return
}

func ParseCmd(st stream.Stream) Result {
	st.ReadUint16()
	st.ReadUint16()
	op, _ := st.ReadByte()
	newFrame, err := ParseProtobuf(frame.OpType(op), st.DataSelect(st.Pos(), st.Size()))
	if err != nil {
		log.Printf("Parse Request fail:%s\n", err)
		return Result{}
	}
	return newFrame
}

func ExecuteCmd(result Result, confLock *sync.RWMutex) ([]byte, frame.OpType) {
	var re Result
	switch result.optype {
	case frame.GET:
		re = getResultFromConf(result, confLock)
		data, err := proto.Marshal(&re.person)
		if err != nil {
			log.Println(err)

			return nil, re.optype
		}
		return data, re.optype
	case frame.ADD:
		re = addResultToConf(result, confLock)
	case frame.DELETE:
		re = deleteResultFromConf(result, confLock)

	default:
		return nil, re.optype
	}
	data, err := proto.Marshal(&re.request)
	if err != nil {
		log.Println(err)
		return nil, re.optype
	}
	return data, re.optype

}

func addResultToConf(result Result, confLock *sync.RWMutex) Result {
	confLock.Lock()
	defer confLock.Unlock()
	cfg, err := ini.Load(conf)
	if err != nil {
		log.Println(err)
		result.optype |= frame.FAIL
		return result
	}
	section := "Person" + strconv.Itoa(int(result.person.GetId()))
	_, err = cfg.Section(section).NewKey("Id", strconv.Itoa(int(result.person.GetId())))
	if err != nil {
		log.Println(err)
		result.optype |= frame.FAIL
		return result
	}

	_, err = cfg.Section(section).NewKey("Name", result.person.GetName())
	if err != nil {
		log.Println(err)
		result.optype |= frame.FAIL
		return result
	}

	_, err = cfg.Section(section).NewKey("Age", strconv.Itoa(int(result.person.GetAge())))
	if err != nil {
		log.Println(err)
		result.optype |= frame.FAIL
		return result
	}

	_, err = cfg.Section(section).NewKey("City", result.person.GetCity())
	if err != nil {
		log.Println(err)
		result.optype |= frame.FAIL
		return result
	}

	if err := cfg.SaveTo(conf); err != nil {
		log.Println(err)
		result.optype |= frame.FAIL
		return result
	}
	result.optype |= frame.SUCCESS
	log.Printf("trying to Add Person whose id is %d success\n", result.person.GetId())
	return result
}

func deleteResultFromConf(result Result, confLock *sync.RWMutex) Result {
	confLock.Lock()
	defer confLock.Unlock()
	cfg, err := ini.Load(conf)
	if err != nil {
		log.Println(err)
		result.optype |= frame.FAIL
		return result
	}
	section := "Person" + strconv.Itoa(int(result.request.GetId()))
	cfg.DeleteSection(section)
	if err := cfg.SaveTo(conf); err != nil {
		log.Println(err)
		result.optype |= frame.FAIL
		return result
	}
	log.Printf("trying to delete Person whose id is %d success\n", result.request.GetId())
	result.optype |= frame.SUCCESS
	return result
}

func getResultFromConf(result Result, confLock *sync.RWMutex) Result {
	confLock.RLock()
	defer confLock.RUnlock()
	cfg, err := ini.Load(conf)
	if err != nil {
		log.Println(err)
		result.optype |= frame.FAIL
		return result
	}
	section := "Person" + strconv.Itoa(int(result.request.GetId()))
	if _, err := cfg.GetSection(section); err != nil {
		 person  := frame.Person{
			Id:   proto.Int32(0),
			Name: proto.String("doesn't exist"),
			Age:  proto.Int32(0),
			City: proto.String("doesn't exist"),
		}
		result.person = person
		result.optype |= frame.FAIL
		log.Printf("trying to get Person whose id is %d fail,doesn't exist\n", result.request.GetId())
		return result
	}

	id, _ := strconv.Atoi(cfg.Section(section).Key("Id").String())
	age, _ := strconv.Atoi(cfg.Section(section).Key("Age").String())
	person := frame.Person{
		Id:   proto.Int32(int32(id)),
		Name: proto.String(cfg.Section(section).Key("Name").String()),
		Age:  proto.Int32(int32(age)),
		City: proto.String(cfg.Section(section).Key("City").String()),
	}
	result.person = person

	result.optype |= frame.SUCCESS
	log.Printf("trying to get Person whose id is %d success\n", result.request.GetId())
	return result
}

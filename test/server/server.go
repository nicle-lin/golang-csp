package main

import (
	"context"
	"github.com/go-ini/ini"
	"github.com/nicle-lin/golang-csp/common"
	"github.com/nicle-lin/golang-csp/frame"
	"github.com/nicle-lin/golang-csp/session"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type Server interface {
	//开始监听等一系列后续动作
	Start()
	//通知所有的goroutines退出
	Stop()
}

type server struct {
	confLock  *sync.RWMutex      //读写配置用到的锁
	exitCh    chan struct{}      //收到Signal后通知所有goroutines退出
	waitGroup *sync.WaitGroup    //等待所有的goroutines退出
	context   context.Context    //用于通知所有子goroutines cancel
	cancel    context.CancelFunc //用于通知所有子goroutines cancel
}

func NewServer() Server {
	//用于通知所有子goroutines cancel
	ctx, cancel := context.WithCancel(context.Background())
	return &server{
		confLock:  new(sync.RWMutex),
		exitCh:    make(chan struct{}),
		waitGroup: &sync.WaitGroup{},
		context:   ctx,
		cancel:    cancel,
	}
}

//从socket里把数据读出来,并解析出命令,然后执行命令,把命令结果通过channel发送给
//writeLoop()函数去发送数据
func (srv *server) readLoop(se session.Session, resultCh chan<- interface{}, done chan<- struct{}) {
	srv.waitGroup.Add(1)
	defer func() {
		close(resultCh) //只有发送方才可以关闭通道,要不然会引起panic
		done <- struct{}{}
		srv.waitGroup.Done()
	}()
	for {
		select {
		case <-srv.context.Done():
			break
		default:

		}
		//从socket读取
		st, err := se.Read()
		if err == io.ErrUnexpectedEOF || err == io.EOF {
			log.Println("connection has been closed by peer")
			break
		} else if err != nil {
			log.Println("read error:", err)
			break
		}
		//分析socket读出来数据,返回一个Result结构体
		result := common.ParseCmd(st)
		//执行命令
		data, optype := common.ExecuteCmd(result, srv.confLock)
		//通过channel发送结果给writeLoop()
		resultCh <- data
		resultCh <- optype

	}

	return

}

//从channel里读数据,所有把数据发送到socket里
func (srv *server) writeLoop(se session.Session, resultCh <-chan interface{}, done chan<- struct{}) {
	srv.waitGroup.Add(1)
	defer func() {
		done <- struct{}{}
		se.Close()
		srv.waitGroup.Done()
	}()
	f := frame.NewFrame(nil, frame.SUCCESS)
	for {
		select {
		case <-srv.context.Done():
			break
		default:

		}
		for i := 0; i < 2; i++ {
			elem, isClosed := <-resultCh
			if !isClosed {
				return
			}
			switch elem.(type) {
			case []byte:
				data, ok := interface{}(elem).([]byte)
				if ok {
					f.SetFrameData(data)
				} else {
					log.Println("change interface{} to []byte fail")
				}

			case frame.OpType:
				optype, ok := interface{}(elem).(frame.OpType)
				if ok {

					f.SetFrameOpType(optype)

				} else {
					log.Println("change interface{} to OpType fail")
				}

			}
		}

		//发送从channel中读取的数据,发送到socket
		_, err := se.Write(f)
		if err == io.EOF || err == io.ErrClosedPipe {
			log.Println("connection has been closed by peer")
			break
		} else if err != nil {
			log.Println("write error:", err)
			break
		}

	}
	return

}

func (srv *server) newConn(conn net.Conn) {
	srv.waitGroup.Add(1)

	//如果3秒都没有读到或写完说明网络太差了
	//conn.SetReadDeadline(time.Now().Add(time.Second * 3))
	se := session.NewSession(conn)

	done := make(chan struct{}, 2)
	defer func() {
		close(done)
		se.Close()
		srv.waitGroup.Done()
	}()
	//用于从readLoop的结果发送到writeLoop
	resultCh := make(chan interface{}, 2)
	go srv.readLoop(se, resultCh, done)
	go srv.writeLoop(se, resultCh, done)
	<-done
	<-done
	log.Println("close connection....")
	return

}

func (srv *server) Start() {
	go func() {
		<-srv.exitCh //收到退出信号
		srv.cancel() //取消所有goroutines
	}()

	l, err := net.Listen(network, address)
	if err != nil {
		log.Println(err)
	}
	defer l.Close()
	log.Printf("Listening on %s\n", address)

	for {
		select {
		case <-srv.exitCh: //收到退出信号
			srv.cancel() //取消所有goroutines
			return
		default:

		}
		log.Println("waiting for a new connection.....")
		conn, err := l.Accept()
		if err != nil {
			log.Println(err)
		}
		//新的连接到来开始处理
		log.Printf("a new connection coming(%s)....", conn.RemoteAddr())
		go srv.newConn(conn)
	}

}

//收到退出信号调用此函数
func (srv *server) Stop() {
	close(srv.exitCh)
	srv.waitGroup.Wait() //等待所有goroutines优雅退出
}

var network, address string = "tcp", "127.0.0.1:8989"

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

func main() {

	srv := NewServer()
	go srv.Start()

	SigRecv := make(chan os.Signal, 1)
	//catch ctr+c signal
	signal.Notify(SigRecv, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	sig := <-SigRecv
	log.Printf("catch signal(%s) to exit\n", sig.String())
	srv.Stop() //开始通知退出

}

package geerpc

import (
	"encoding/json"
	"errors"
	"fmt"
	"geerpc/codec"
	"io"
	"log"
	"net"
	"reflect"
	"strings"
	"sync"
	"time"
)

const MagicNumber = 0x3bef5c

type Option struct {
	MagicNumber    int
	CodecType      codec.Type
	ConnectTimeout time.Duration // 连接超时
	HandleTimeout  time.Duration
}

var DefaultOption = &Option{
	MagicNumber:    MagicNumber,
	CodecType:      codec.GobType,
	ConnectTimeout: time.Second * 10, // 设置超时时长 10s
}

type Server struct {
	serviceMap sync.Map
}

func NewServer() *Server {
	return &Server{}
}

var DefaultServer = NewServer()

func (server *Server) Accept(lis net.Listener) {
	// for 循环等待 socket 连接建立
	for {
		conn, err := lis.Accept()
		if err != nil {
			log.Println("rpc server: accept error: ", err)
			return
		}
		go server.ServerConn(conn) // 开启子协程处理
	}
}

func Accept(lis net.Listener) {
	DefaultServer.Accept(lis)
}

// ServerConn 处理 socket 连接，反序列化 Option 实例
func (server *Server) ServerConn(conn io.ReadWriteCloser) {
	defer func() { _ = conn.Close() }()
	var opt Option
	if err := json.NewDecoder(conn).Decode(&opt); err != nil {
		log.Println("rpc server: options error: ", err)
		return
	}

	if opt.MagicNumber != MagicNumber {
		log.Printf("rpc server: invalid magic number: %x ", opt.MagicNumber)
		return
	}

	f := codec.NewCodecFuncMap[opt.CodecType]
	if f == nil {
		log.Printf("rpc server: invalid codec type %s", opt.CodecType)
	}
	server.serverCodec(f(conn), &opt)
}

var invalidRequest = struct{}{}

// serverCodec
//func (server *Server) serverCodec_v1(cc codec.Codec) {
//	sending := new(sync.Mutex)
//	wg := new(sync.WaitGroup)
//	for {
//		req, err := server.readRequest(cc)
//		if err != nil {
//			if req == nil {
//				break
//			}
//			req.h.Error = err.Error()
//			server.sendResponse(cc, req.h, invalidRequest, sending)
//			continue
//		}
//		wg.Add(1)
//		go server.handleRequest(cc, req, sending, wg) // handleRequest 使用协程并发执行请求
//
//	}
//}

// serverCodec 添加超时
func (server *Server) serverCodec(cc codec.Codec, opt *Option) {
	sending := new(sync.Mutex)
	wg := new(sync.WaitGroup)
	for {
		req, err := server.readRequest(cc)
		if err != nil {

			if req == nil {
				break
			}
			req.h.Error = err.Error()
			server.sendResponse(cc, req.h, invalidRequest, sending)
			continue
		}
		wg.Add(1)
		go server.handleRequest(cc, req, sending, wg, opt.HandleTimeout)
	}
	wg.Wait()
	_ = cc.Close()
}

// 处理请求时并发的，但回复请求的报文必须逐个发送，并发回复，客户端无法解析

type request struct {
	h            *codec.Header
	argv, replyv reflect.Value
	mtype        *methodType
	svc          *service
}

func (server *Server) readRequestHander(cc codec.Codec) (*codec.Header, error) {
	var h codec.Header
	if err := cc.ReadHeader(&h); err != nil {
		if err != io.EOF && err != io.ErrUnexpectedEOF {
			log.Println("rpc server: read header error: ", err)
		}
		return nil, err
	}
	return &h, nil
}

func (server *Server) readRequest_v1(cc codec.Codec) (*request, error) {
	h, err := server.readRequestHander(cc)
	if err != nil {
		return nil, err
	}

	req := &request{h: h}

	req.argv = reflect.New(reflect.TypeOf(""))
	if err = cc.ReadBody(req.argv.Interface()); err != nil {
		log.Println("rpc server: read argv err: ", err)
	}
	return req, err
}

func (server *Server) readRequest(cc codec.Codec) (*request, error) {
	h, err := server.readRequestHander(cc)
	if err != nil {
		return nil, err
	}
	req := &request{h: h}

	req.svc, req.mtype, err = server.findService(h.ServiceMethod) // 查找服务
	if err != nil {
		return req, err
	}
	req.argv = req.mtype.newArgv() // 创建入参实例
	req.replyv = req.mtype.newReplyv()

	argvi := req.argv.Interface() // argvi 确保是一个指针
	if req.argv.Type().Kind() != reflect.Ptr {
		argvi = req.argv.Addr().Interface()
	}
	if err = cc.ReadBody(argvi); err != nil {
		return req, err
	}
	return req, err

}

func (server *Server) sendResponse(cc codec.Codec, h *codec.Header, body interface{}, sending *sync.Mutex) {
	sending.Lock()
	defer sending.Unlock()
	if err := cc.Write(h, body); err != nil {
		log.Println("rpc server: write response error: ", err)
	}
}

// handleRequest 使用协程并发执行请求
func (server *Server) handleRequest_v1(cc codec.Codec, req *request, sending *sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done()
	log.Println(req.h, req.argv.Elem())
	req.replyv = reflect.ValueOf(fmt.Sprintf("geerpc resp %d", req.h.Seq))
	server.sendResponse(cc, req.h, req.replyv.Interface(), sending)
}

func (server *Server) handleRequest_V2(cc codec.Codec, req *request, sending *sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done()
	err := req.svc.call(req.mtype, req.argv, req.replyv)
	if err != nil {
		req.h.Error = err.Error()
		server.sendResponse(cc, req.h, invalidRequest, sending)
		return
	}
	server.sendResponse(cc, req.h, req.replyv.Interface(), sending)
}

// handleRequest 服务端处理超时
func (server *Server) handleRequest(cc codec.Codec, req *request, sending *sync.Mutex, wg *sync.WaitGroup, timeout time.Duration) {
	defer wg.Done()
	called := make(chan struct{})
	sent := make(chan struct{})
	go func() {
		err := req.svc.call(req.mtype, req.argv, req.replyv)
		called <- struct{}{}
		if err != nil {
			req.h.Error = err.Error()
			server.sendResponse(cc, req.h, invalidRequest, sending) // 无效请求
			sent <- struct{}{}
			return
		}
		server.sendResponse(cc, req.h, req.replyv.Interface(), sending)
		sent <- struct{}{}
	}()

	if timeout == 0 {
		<-called
		<-sent
		return
	}

	select {
	case <-time.After(timeout):
		req.h.Error = fmt.Sprintf("rpc server: request handle timeout: expect within %s", timeout)
		server.sendResponse(cc, req.h, invalidRequest, sending)
	case <-called:
		<-sent
	}
}

func (server *Server) Register(rcvr interface{}) error {
	s := newService(rcvr)
	if _, dup := server.serviceMap.LoadOrStore(s.name, s); dup {
		return errors.New("rpc: service already defined: " + s.name)
	}
	return nil
}

func Register(rcvr interface{}) error {
	return DefaultServer.Register(rcvr)
}

func (server *Server) findService(serviceMethod string) (svc *service, mtype *methodType, err error) {
	dot := strings.LastIndex(serviceMethod, ".")
	if dot < 0 {
		err = errors.New("rpc server: service/method request ill-formed: " + serviceMethod)
		return
	}

	// Service 服务名，第二个方法名
	serviceName, methodName := serviceMethod[:dot], serviceMethod[dot+1:]
	svci, ok := server.serviceMap.Load(serviceName) // 找到对应的 service 实例
	if !ok {
		err = errors.New("rpc server: can't find service " + serviceName)
		return
	}
	svc = svci.(*service)
	mtype = svc.method[methodName] // 再从 service 实例的 method 中找到对应的 methodType
	if mtype == nil {
		err = errors.New("rpc server: can't find method " + methodName)
	}
	return
}

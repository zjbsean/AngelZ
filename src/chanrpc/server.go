package chanrpc

import (
	"conf"
	"errors"
	"fmt"
	"logger"
	"runtime"
)

//Server :
//  manage functions that can be call, and deal the call information that from channel
type Server struct {
	//can be call functions : id -> functions
	//
	//functions:
	//  func(args []interface{})
	//  func(args []interface{}) interface{}
	//  func(args []interface{}) []interface{}
	functions map[interface{}]interface{}

	//call infomation chan
	ChanCall chan *CallInfo
}

//CallInfo : the call information that be deal by Server object.
type CallInfo struct {
	//need execute function
	//  func(args []interface{})
	//  func(args []interface{}) interface{}
	//  func(args []interface{}) []interface{}
	f interface{}

	//the parameters of function f
	args []interface{}

	//RetInfo
	//	the channel of return infomation
	chanRet chan *RetInfo

	//the function that be execute by Client object when call back:
	//  func(err error)
	//  func(ret interface{}, err error)
	//  func(ret []interface{}, err error)
	cb interface{}
}

//RetInfo :
//  the return infomation, that Server object execute Client object call infomation.
type RetInfo struct {
	// the return data that Server object execute Client object call infomation
	ret interface{}

	// err != nil : when Server object execute Client object call infomation have an error
	err error

	// the callback function
	//  func(err error)
	//  func(ret interface{}, err error)
	//  func(ret []interface{}, err error)
	cb interface{}
}

//NewServer :
// create a Server Object
//Params:
//	l : size of ChanCall chan
func NewServer(l int) *Server {
	s := new(Server)
	s.functions = make(map[interface{}]interface{})
	s.ChanCall = make(chan *CallInfo, l)
	return s
}

//interfaceToSliceInterface :
//	make interface{} obj to []interface{}
//
//Params:
//	i : nil or []interface{} interface
func interfaceToSliceInterface(i interface{}) []interface{} {
	if i == nil {
		return nil
	}

	return i.([]interface{})
}

//Open :
//	open Client object
//
//Params:
//	l : size of client asynchronous callback chan
func (s *Server) Open(l int) *Client {
	c := NewClient(l)
	c.Attach(s)
	return c
}

//Register :
//	register functions that can be call by Client object and execute with Server object
func (s *Server) Register(id interface{}, f interface{}) {
	switch f.(type) {
	case func([]interface{}):
	case func([]interface{}) interface{}:
	case func([]interface{}) []interface{}:
	default:
		panic(fmt.Sprintf("function id %v: definition of function is invalid", id))
	}
	if _, ok := s.functions[id]; ok {
		panic(fmt.Sprintf("function id %v: already registered", id))
	}

	s.functions[id] = f
}

//ret :
//	set callback function to return data and push to return chan
//
//Params:
// ci : the call information that Client object request
// ri : the return infomation after execute call information
func (s *Server) ret(ci *CallInfo, ri *RetInfo) (err error) {
	if ci.chanRet == nil {
		return
	}

	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()

	ri.cb = ci.cb
	ci.chanRet <- ri
	return
}

//exec:
//	exectue call infomation and deal return data
//
//Params:
//	ci : call infomation
func (s *Server) exec(ci *CallInfo) (err error) {
	defer func() {
		if r := recover(); r != nil {
			if conf.LenStackBuf > 0 {
				buf := make([]byte, conf.LenStackBuf)
				l := runtime.Stack(buf, false)
				err = fmt.Errorf("%v %s", r, buf[:l])
			} else {
				err = fmt.Errorf("%v", r)
			}

			s.ret(ci, &RetInfo{err: fmt.Errorf("%v", r)})
		}
	}()

	switch ci.f.(type) {
	case func([]interface{}):
		ci.f.(func([]interface{}))(ci.args)
		return s.ret(ci, &RetInfo{})
	case func([]interface{}) interface{}:
		rd := ci.f.(func([]interface{}) interface{})(ci.args)
		return s.ret(ci, &RetInfo{ret: rd})
	case func([]interface{}) []interface{}:
		rd := ci.f.(func([]interface{}) []interface{})(ci.args)
		return s.ret(ci, &RetInfo{ret: rd})
	}

	panic("the call function type is unkonw !")
}

//Exec :
//	exectue call infomation and log error
func (s *Server) Exec(ci *CallInfo) {
	err := s.exec(ci)
	if err != nil {
		logger.Error("%v", err)
	}
}

//Go : (goroutine safe)
//	execute command asynchronous
func (s *Server) Go(id interface{}, args ...interface{}) {
	f := s.functions[id]
	if f == nil {
		return
	}

	defer func() {
		recover()
	}()

	s.ChanCall <- &CallInfo{
		f:    f,
		args: args,
	}
}

//Call0 :
//	call command synchronization with no return data
//
//Params:
//	id : command id
//	args : command params
func (s *Server) Call0(id interface{}, args ...interface{}) error {
	return s.Open(0).Call0(id, args...)
}

//Call1 :
//	call command synchronization with one return data
//
//Params:
//	id : command id
//	args : command params
func (s *Server) Call1(id interface{}, args ...interface{}) (interface{}, error) {
	return s.Open(0).Call1(id, args...)
}

//CallN :
//	call command synchronization with many return data
//
//Params:
//	id : command id
//	args : command params
func (s *Server) CallN(id interface{}, args ...interface{}) ([]interface{}, error) {
	return s.Open(0).CallN(id, args...)
}

//Close :
//	close command server
func (s *Server) Close() {
	close(s.ChanCall)

	for ci := range s.ChanCall {
		s.ret(ci, &RetInfo{
			err: errors.New("chanrps server closed"),
		})
	}
}

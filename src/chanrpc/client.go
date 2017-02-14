package chanrpc

import (
	"conf"
	"errors"
	"fmt"
	"logger"
	"runtime"
)

//Client :
//	call the function that register in Server and deal the callback.
//	deal callback support synchronization and asynchronous mode.
type Client struct {
	//the object deal the call
	s *Server

	//callback informaction channel with synchronization mode
	chanSyncRet chan *RetInfo

	//callback informaction channel with asynchronous mode
	ChanAsynRet chan *RetInfo

	//pending asynchronous call count
	pendingAsynCall int
}

//NewClient :
//	new a Client object
//
//Params:
//	l : size of asynchronous callback chan
func NewClient(l int) *Client {
	c := new(Client)
	c.chanSyncRet = make(chan *RetInfo, 1)
	c.ChanAsynRet = make(chan *RetInfo, l)
	return c
}

//Attach :
//	attach a Server object
func (c *Client) Attach(s *Server) {
	c.s = s
}

//call :
//	call call infomaction by synchronization or asynchronous mode
//
//Params:
//	ci : call informaction
//	block : synchronization -> true; asynchronous -> false
func (c *Client) call(ci *CallInfo, block bool) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()

	if block {
		c.s.ChanCall <- ci
	} else {
		select {
		case c.s.ChanCall <- ci:
		default:
			err = errors.New("chanrpc channel full")
		}
	}
	return
}

//f :
//	function get from registed in server
//
//Params:
//	id : function call command
//	cbParamMode : return data mode 0/1/2
func (c *Client) f(id interface{}, cbParamMode int) (f interface{}, err error) {
	if c.s == nil {
		err = errors.New("server not attached")
		return
	}

	f = c.s.functions[id]
	if f == nil {
		err = fmt.Errorf("function id %v: function not registed", id)
		return
	}

	var ok bool
	switch cbParamMode {
	case 0:
		_, ok = f.(func([]interface{}))
	case 1:
		_, ok = f.(func([]interface{}) interface{})
	case 2:
		_, ok = f.(func([]interface{}) []interface{})
	default:
		panic("call function return data count only 0/1/2 three mode")
	}

	if !ok {
		err = fmt.Errorf("function id %v: return type mismatch, want return data count %d", id, cbParamMode)
	}
	return
}

//Call0 :
//	call command synchronous, command no return date
//
//Params:
//	id : command id
//	args : command params
func (c *Client) Call0(id interface{}, args ...interface{}) (err error) {
	f, err := c.f(id, 0)
	if err != nil {
		return err
	}

	err = c.call(&CallInfo{
		f:       f,
		args:    args,
		chanRet: c.chanSyncRet,
	}, true)
	if err != nil {
		return err
	}
	ri := <-c.chanSyncRet
	return ri.err
}

//Call1 :
//	call command synchronous, command had noe return date
//
//Params:
//	id : command id
//	args : command params
func (c *Client) Call1(id interface{}, args ...interface{}) (ret interface{}, err error) {
	f, err := c.f(id, 1)
	if err != nil {
		return nil, err
	}

	err = c.call(&CallInfo{
		f:       f,
		args:    args,
		chanRet: c.chanSyncRet,
	}, true)
	if err != nil {
		return nil, err
	}
	ri := <-c.chanSyncRet
	return ri.ret, ri.err
}

//CallN :
//	call command synchronous, command had many return date
//
//Params:
//	id : command id
//	args : command params
func (c *Client) CallN(id interface{}, args ...interface{}) (ret []interface{}, err error) {
	f, err := c.f(id, 2)
	if err != nil {
		return nil, err
	}

	err = c.call(&CallInfo{
		f:       f,
		args:    args,
		chanRet: c.chanSyncRet,
	}, true)
	if err != nil {
		return nil, err
	}
	ri := <-c.chanSyncRet
	return interfaceToSliceInterface(ri.ret), err
}

//asynCall :
//	call command asynchronous
//
//Params:
//	id : call command
//	args : call command params
// 	cb : asynchronous call callback
//  cbParamMode : return data mode 0/1/2
func (c *Client) asynCall(id interface{}, args []interface{}, cb interface{}, cbParamMode int) {
	f, err := c.f(id, cbParamMode)
	if err != nil {
		c.ChanAsynRet <- &RetInfo{err: err, cb: cb}
		return
	}
	err = c.call(&CallInfo{
		f:       f,
		args:    args,
		chanRet: c.ChanAsynRet,
		cb:      cb,
	}, false)
	if err != nil {
		c.ChanAsynRet <- &RetInfo{err: err, cb: cb}
		return
	}
}

//AsynCall :
//	call command asynchronous
//
//Params:
//	id : call command
//	_args : call command params + callback function
func (c *Client) AsynCall(id interface{}, _args ...interface{}) {
	if len(_args) < 1 {
		panic("callback func not found")
	}

	args := _args[:len(_args)-1]
	cb := _args[len(_args)-1]

	var n int
	switch cb.(type) {
	case func(error):
		n = 0
	case func(interface{}, error):
		n = 1
	case func([]interface{}, error):
		n = 2
	default:
		panic("definition of callback function is invalid")
	}

	if c.pendingAsynCall >= cap(c.ChanAsynRet) {
		execCb(&RetInfo{err: errors.New("too many calls"), cb: cb})
		return
	}

	c.asynCall(id, args, cb, n)
	c.pendingAsynCall++
}

//execCb :
//	execute command callback function
//
//Params:
//	ri : command return data
func execCb(ri *RetInfo) {
	defer func() {
		if r := recover(); r != nil {
			if conf.LenStackBuf > 0 {
				buf := make([]byte, conf.LenStackBuf)
				l := runtime.Stack(buf, false)
				logger.Error("%v: %s", r, buf[:l])
			} else {
				logger.Error("%v", r)
			}
		}
	}()

	switch ri.cb.(type) {
	case func(error):
		ri.cb.(func(error))(ri.err)
	case func(interface{}, error):
		ri.cb.(func(interface{}, error))(ri.ret, ri.err)
	case func([]interface{}, error):
		ri.cb.(func([]interface{}, error))(interfaceToSliceInterface(ri.ret), ri.err)
	default:
		panic("definition of callback function is invalid")
	}
	return
}

//Cb :
//	execute callback function
//
//Params:
//	ri : command return data
func (c *Client) Cb(ri *RetInfo) {
	c.pendingAsynCall--
	execCb(ri)
}

//Close :
//	close client
func (c *Client) Close() {
	for c.pendingAsynCall > 0 {
		c.Cb(<-c.ChanAsynRet)
	}
}

//Idle :
//	client is idle
func (c *Client) Idle() bool {
	return c.pendingAsynCall == 0
}

func (c *Client) PendingAsynCllCount() int {
	return c.pendingAsynCall
}

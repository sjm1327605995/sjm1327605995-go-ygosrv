package protocol

import (
	"bytes"
	"github.com/panjf2000/gnet/v2"
	"github.com/panjf2000/gnet/v2/pkg/logging"
	"go-ygosrv/server/protocol/tcp"
	"go-ygosrv/server/protocol/websocket"
	"time"
)

type Server struct {
	gnet.BuiltinEventEngine

	addr      string
	multicore bool
	eng       gnet.Engine
}

func (wss *Server) OnBoot(eng gnet.Engine) gnet.Action {
	wss.eng = eng
	logging.Infof("echo server with multi-core=%t is listening on %s", wss.multicore, wss.addr)
	return gnet.None
}

func (wss *Server) OnOpen(c gnet.Conn) ([]byte, gnet.Action) {
	c.SetContext(new(Context))

	return nil, gnet.None
}

func (wss *Server) OnClose(c gnet.Conn, err error) (action gnet.Action) {
	if err != nil {
		logging.Warnf("error occurred on connection=%s, %v\n", c.RemoteAddr().String(), err)
	}

	logging.Infof("conn[%v] disconnected", c.RemoteAddr().String())
	return gnet.None
}

func (wss *Server) OnTraffic(c gnet.Conn) (action gnet.Action) {
	ws := c.Context().(*Context)
	//先把数据读取到buff中
	if ws.readBufferBytes(c) == gnet.Close {
		return gnet.Close
	}

	switch ws.protocol {
	case 0: //等待解析
		//不包含websocket当做tcp处理
		if !bytes.Contains(ws.buf.Bytes(), []byte("Upgrade: websocket")) {
			ws.protocol = 1
			ws.Decoder = &tcp.TCPDecoder{}
		} else {
			ws.protocol = 2
			ws.Decoder = &websocket.WsDecoder{}
		}
	case 1, 2:

	default:
		return gnet.Close
	}

	return ws.Decoder.Decode(c, &ws.buf)
}

func (wss *Server) OnTick() (delay time.Duration, action gnet.Action) {
	return 3 * time.Second, gnet.None
}

type Decoder interface {
	Decode(c gnet.Conn, buffer *bytes.Buffer) gnet.Action
}
type Context struct {
	protocol uint8
	buf      bytes.Buffer // 从实际socket中读取到的数据缓存
	Decoder  Decoder
}

func (w *Context) readBufferBytes(c gnet.Conn) gnet.Action {
	size := c.InboundBuffered()
	buf := make([]byte, size, size)
	read, err := c.Read(buf)
	if err != nil {
		logging.Infof("read err! %w", err)
		return gnet.Close
	}
	if read < size {
		logging.Infof("read bytes len err! size: %d read: %d", size, read)
		return gnet.Close
	}
	w.buf.Write(buf)
	return gnet.None
}

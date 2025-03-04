package httpconn

import (
	"errors"
	"io"
	"slices"
)

// PrivateProtocolHead HTTP客户端私有协议头
type PrivateProtocolHead struct {
	RequestHead  []byte // 请求前的私有协议头
	ResponseHead []byte // 响应前的私有协议头
	Strict       bool   // 是否严格模式，严格模式下，私有协议头内容必须完全一致，否则只检查长度一致
}

// Clone 克隆
func (p *PrivateProtocolHead) Clone() *PrivateProtocolHead {
	// 准备空间
	tmp := &PrivateProtocolHead{
		RequestHead:  make([]byte, len(p.RequestHead)),
		ResponseHead: make([]byte, len(p.ResponseHead)),
		Strict:       p.Strict,
	}
	// 拷贝协议
	copy(tmp.RequestHead, p.RequestHead)
	copy(tmp.ResponseHead, p.ResponseHead)
	// OK
	return tmp
}

// WriteRequestHead 写入请求前的私有协议头
func (p *PrivateProtocolHead) WriteRequestHead(writer io.Writer) error {
	// 检查私有协议头长度
	priLen := len(p.RequestHead)
	if priLen > 0 {
		// 发送私有协议头
		n, err := writer.Write(p.RequestHead)
		if err != nil {
			return err
		}
		if n != priLen {
			return errors.New("write request private protocol length error")
		}
	}
	// OK
	return nil
}

// ReadRequestHead 读取请求前的私有协议头
func (p *PrivateProtocolHead) ReadRequestHead(reader io.Reader) error {
	// 检查私有协议头长度
	priLen := len(p.RequestHead)
	if priLen > 0 {
		// 读取私有协议头
		privateProtocol := make([]byte, priLen)
		n, err := reader.Read(privateProtocol)
		if err != nil {
			return err
		}
		// 检查长度是否一致
		if n != priLen {
			return errors.New("read request private protocol length error")
		}
		// 严格模式需要比对内容是否一致
		if p.Strict {
			if !slices.Equal(privateProtocol, p.RequestHead) {
				return errors.New("read request private protocol content error")
			}
		}
	}
	// OK
	return nil
}

// WriteResponseHead 写入响应前的私有协议头
func (p *PrivateProtocolHead) WriteResponseHead(writer io.Writer) error {
	// 检查私有协议头长度
	priLen := len(p.ResponseHead)
	if priLen > 0 {
		// 发送私有协议头
		n, err := writer.Write(p.ResponseHead)
		if err != nil {
			return err
		}
		if n != priLen {
			return errors.New("write response private protocol length error")
		}
	}
	// OK
	return nil
}

// ReadResponseHead 读取响应前的私有协议头
func (p *PrivateProtocolHead) ReadResponseHead(reader io.Reader) error {
	// 检查私有协议头长度
	priLen := len(p.ResponseHead)
	if priLen > 0 {
		// 读取私有协议头
		buf := make([]byte, priLen)
		n, err := reader.Read(buf)
		if err != nil {
			return err
		}
		// 检查长度是否一致
		if n != priLen {
			return errors.New("read response private protocol length error")
		}
		// 严格模式需要比对内容是否一致
		if p.Strict {
			if !slices.Equal(buf, p.ResponseHead) {
				return errors.New("read response private protocol content error")
			}
		}
	}
	// OK
	return nil
}

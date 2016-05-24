package fdfs_client

import (
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net"
	"os"
	"time"
)

var ErrClosed = errors.New("pool is closed")

type pConn struct {
	net.Conn
	pool *ConnectionPool
}

func (c pConn) Close() error {
	return c.pool.put(c.Conn)
}

type ConnectionPool struct {
	hosts     []string
	ports     []int
	minConns  int
	maxConns  int
	busyConns []bool
	conns     chan net.Conn
}

func minInt(a int, b int) int {
	if b-a > 0 {
		return a
	} else {
		return b
	}
}

func NewConnectionPool(hosts []string, ports []int, minConns int, maxConns int) (*ConnectionPool, error) {
	if minConns < 0 || maxConns <= 0 || minConns > maxConns {
		err := errors.New("invalid conns settings")
		logger.Error(err.Error())
		return nil, err
	}
	cp := &ConnectionPool{
		hosts:     hosts,
		ports:     ports,
		minConns:  minConns,
		maxConns:  maxConns,
		conns:     make(chan net.Conn, maxConns),
		busyConns: make([]bool, len(hosts)),
	}
	//logger.Debug("cp made")
	for i := 0; i < minInt(MINCONN, len(hosts)); i++ {
		conn, err := cp.makeConn()
		if err != nil {
			cp.Close()
			logger.Error("make connection error" + err.Error())
			return nil, err
		}
		cp.conns <- conn
	}
	return cp, nil
}

func (this *ConnectionPool) Get() (net.Conn, error) {
	conns := this.getConns()
	if conns == nil {
		return nil, ErrClosed
	}

	for {
		select {
		case conn := <-conns:
			if conn == nil {
				break
				//return nil, ErrClosed
			}
			if err := this.activeConn(conn); err != nil {
				break
			}
			return this.wrapConn(conn), nil
		default:
			if this.Len() >= this.maxConns {
				errmsg := fmt.Sprintf("Too many connctions %d", this.Len())
				return nil, errors.New(errmsg)
			}
			conn, err := this.makeConn()
			if err != nil {
				return nil, err
			}

			this.conns <- conn
			//put connection to pool and go next `for` loop
			//return this.wrapConn(conn), nil
		}
	}

}

func (this *ConnectionPool) Close() {
	conns := this.conns
	this.conns = nil

	if conns == nil {
		return
	}

	close(conns)
	logger.Debugf("%d", len(conns))
	for conn := range conns {
		conn.Close()
	}
}

func (this *ConnectionPool) Len() int {
	return len(this.getConns())
}

func (this *ConnectionPool) makeConn() (net.Conn, error) {
	var n int
	for {
		n = rand.Intn(len(this.hosts))
		if !this.busyConns[n] {
			this.busyConns[n] = true
			break
		}
	}
	host := this.hosts[n]
	addr := fmt.Sprintf("%s:%d", host, this.ports[n])

	return net.DialTimeout("tcp", addr, time.Minute)
}

func (this *ConnectionPool) getConns() chan net.Conn {
	conns := this.conns
	return conns
}

func (this *ConnectionPool) put(conn net.Conn) error {
	if conn == nil {
		return errors.New("connection is nil")
	}
	if this.conns == nil {
		return conn.Close()
	}

	select {
	case this.conns <- conn:
		return nil
	default:
		return conn.Close()
	}
}

func (this *ConnectionPool) wrapConn(conn net.Conn) net.Conn {
	c := pConn{pool: this}
	c.Conn = conn
	return c
}

func (this *ConnectionPool) activeConn(conn net.Conn) error {
	th := &trackerHeader{}
	th.cmd = FDFS_PROTO_CMD_ACTIVE_TEST
	th.sendHeader(conn)
	th.recvHeader(conn)
	if th.cmd == 100 && th.status == 0 {
		return nil
	}
	return errors.New("Conn unaliviable")
}

func TcpSendData(conn net.Conn, bytesStream []byte) error {
	if _, err := conn.Write(bytesStream); err != nil {
		return err
	}
	return nil
}

func TcpSendFile(conn net.Conn, filename string) error {
	file, err := os.Open(filename)
	defer file.Close()
	if err != nil {
		return err
	}

	var fileSize int64 = 0
	if fileInfo, err := file.Stat(); err == nil {
		fileSize = fileInfo.Size()
	}

	if fileSize == 0 {
		errmsg := fmt.Sprintf("file size is zeor [%s]", filename)
		return errors.New(errmsg)
	}

	fileBuffer := make([]byte, fileSize)

	_, err = file.Read(fileBuffer)
	if err != nil {
		return err
	}

	return TcpSendData(conn, fileBuffer)
}

func TcpRecvResponse(conn net.Conn, bufferSize int64) ([]byte, int64, error) {
	recvBuff := make([]byte, 0, bufferSize)
	tmp := make([]byte, 256)
	var total int64
	for {
		n, err := conn.Read(tmp)
		total += int64(n)
		recvBuff = append(recvBuff, tmp[:n]...)
		if err != nil {
			if err != io.EOF {
				return nil, 0, err
			}
			break
		}
		if total == bufferSize {
			break
		}
	}
	return recvBuff, total, nil
}

func TcpRecvFile(conn net.Conn, localFilename string, bufferSize int64) (int64, error) {
	file, err := os.Create(localFilename)
	defer file.Close()
	if err != nil {
		return 0, err
	}

	recvBuff, total, err := TcpRecvResponse(conn, bufferSize)
	if _, err := file.Write(recvBuff); err != nil {
		return 0, err
	}
	return total, nil
}

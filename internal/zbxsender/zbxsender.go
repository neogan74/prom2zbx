package zbxsender

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"time"
)

// Metric ...
type Metric struct {
	Host  string `json:"host"`
	Key   string `json:"key"`
	Value string `json:"value"`
	Clock int64  `json:"clock"`
}

//NewMetric func
func NewMetric(host, key, value string, clock ...int64) *Metric {
	m := &Metric{
		Host:  host,
		Key:   key,
		Value: value,
	}
	if m.Clock = time.Now().Unix(); len(clock) > 0 {
		m.Clock = int64(clock[0])
	}
	return m
}

//ZbxPacket struct
type ZbxPacket struct {
	Request string    `json:"request"`
	Data    []*Metric `json:"data"`
	Clock   int64     `json:"clock"`
}

//NewPacket ...
func NewPacket(data []*Metric, clock ...int64) *ZbxPacket {
	p := &ZbxPacket{
		Request: `sender data`,
		Data:    data,
	}
	if p.Clock = time.Now().Unix(); len(clock) > 0 {
		p.Clock = int64(clock[0])
	}
	return p
}

//DataLen ...
func (p *ZbxPacket) DataLen() []byte {
	datalen := make([]byte, 8)
	JSONData, _ := json.Marshal(p)
	binary.LittleEndian.PutUint32(datalen, uint32(len(JSONData)))
	return datalen
}

//Zsender ...
type Zsender struct {
	Host string
	Port int
}

//NewSender func
func NewSender(host string, port int) *Zsender {
	s := &Zsender{
		Host: host,
		Port: port}
	return s
}

func (zs *Zsender) getHeader() []byte {
	return []byte("ZBXD\\x01")
}

func (zs *Zsender) getTCPAddr() (ipaddr *net.TCPAddr, err error) {
	addr := fmt.Sprintf("%s:%d", zs.Host, zs.Port)

	ipaddr, err = net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		err = fmt.Errorf("Connection failed: %s", err.Error())
		return
	}

	return
}

func (zs *Zsender) read(conn *net.TCPConn) (res []byte, err error) {
	res = make([]byte, 1024)
	res, err = ioutil.ReadAll(conn)
	if err != nil {
		err = fmt.Errorf("Error while receiving data: %s", err.Error())
		return
	}
	return
}

func (zs *Zsender) connect() (conn *net.TCPConn, err error) {
	type DialResp struct {
		Conn  *net.TCPConn
		Error error
	}

	ipaddr, err := zs.getTCPAddr()
	if err != nil {
		return
	}
	ch := make(chan DialResp)

	go func() {
		conn, err = net.DialTCP("tcp", nil, ipaddr)
		ch <- DialResp{Conn: conn, Error: err}
	}()
	select {
	case <-time.After(5 * time.Second):
		err = fmt.Errorf("Connection timeout")
	case resp := <-ch:
		if resp.Error != nil {
			err = resp.Error
			break
		}
		conn = resp.Conn
	}
	return
}

//Send method
func (zs *Zsender) Send(packet *ZbxPacket) (res []byte, err error) {
	conn, err := zs.connect()
	if err != nil {
		return
	}
	defer conn.Close()

	dataPacket, _ := json.Marshal(packet)

	buffer := append(zs.getHeader(), packet.DataLen()...)
	buffer = append(buffer, dataPacket...)

	_, err = conn.Write(buffer)
	if err != nil {
		err = fmt.Errorf("Error while sending the data: %s", err.Error())
		return
	}

	res, err = zs.read(conn)

	return

}

package fdfs_client

import (
	"fmt"
	"testing"
)

func getConn(pool *ConnectionPool) {
	conn, err := pool.Get()
	defer func() {
		if conn != nil {
			conn.Close()
		}
	}()
	if err != nil {
		fmt.Printf("get conn error:%s\n", err)
	}
}

func TestGetConnection(t *testing.T) {
	Config, err := getConf("client.conf")
	if err != nil {
		t.Error(err)
	}
	hosts := Config.TrackerIp
	ports := Config.TrackerPort
	minConns := Config.MinConn
	maxConns := Config.MaxConn
	pool, err := NewConnectionPool(hosts, ports, minConns, maxConns)
	if err != nil {
		t.Error(err)
		return
	}
	for i := 0; i < 100; i++ {
		go getConn(pool)
	}
}
func TestConnetionPoolClose(t *testing.T) {
	Config, err := getConf("client.conf")
	if err != nil {
		t.Error(err)
	}
	hosts := Config.TrackerIp
	ports := Config.TrackerPort
	minConns := Config.MinConn
	maxConns := Config.MaxConn
	pool, err := NewConnectionPool(hosts, ports, minConns, maxConns)
	if err != nil {
		t.Error(err)
		return
	}
	pool.Close()
}
func BenchmarkGetConnection(b *testing.B) {
	Config, err := getConf("client.conf")
	if err != nil {
		b.Error(err)
	}
	hosts := Config.TrackerIp
	ports := Config.TrackerPort
	minConns := Config.MinConn
	maxConns := Config.MaxConn
	pool, err := NewConnectionPool(hosts, ports, minConns, maxConns)
	if err != nil {
		b.Error(err)
		return
	}
	b.StopTimer()
	b.StartTimer()
	for i := 0; i < 10000; i++ {
		go getConn(pool)
	}
}

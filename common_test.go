package fdfs_client

import (
	//"fmt"
	"testing"
)

func TestGetConf(t *testing.T) {
	Config := &Config{}
	Config, err := getConf("client.conf")
	if err != nil {
		t.Error(err)
	}
	//t.Log("\nTrackerIp   : %s\nTrackerPort : %d\nMaxConn     : %d\nNet_Timeout : %d\nCon_Timeout : %d", Config.TrackerIp[0], Config.TrackerPort[0], Config.MaxConn, Config.Net_Timeout, Config.Con_Timeout)
	//fmt.Print("Hello")
	t.Log("Hello")
	t.Log(Config.TrackerIp)
	t.Log(Config.TrackerPort)
	t.Log(Config.MaxConn)
	t.Log(Config.Net_Timeout)
	t.Log(Config.Con_Timeout)
}

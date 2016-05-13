package fdfs_client

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var (
	MAXCONN int
	MINCONN int
)

type Config struct {
	TrackerIp   []string
	TrackerPort []int
	MaxConn     int
	MinConn     int
	Net_Timeout int
	Con_Timeout int
}

func getConf(ConfPath string) (*Config, error) {
	//logger.Debug("getConf begin")
	fc := &FdfsConfigParser{}
	cf, err := fc.Read(ConfPath)
	if err != nil {
		logger.Errorf("Read conf error :%s", err)
		return nil, err
	}

	trackerListString, _ := cf.RawString("DEFAULT", "tracker_server")
	trackerList := strings.Split(trackerListString, ",")

	max_conn, _ := cf.RawString("DEFAULT", "max_conn")

	network_timeout, _ := cf.RawString("DEFAULT", "network_timeout")

	connect_timeout, _ := cf.RawString("DEFAULT", "connect_timeout")

	min_conn, _ := cf.RawString("DEFAULT", "min_conn")

	var (
		trackerIpList   []string
		trackerPortList []int
		maxc            int
		minc            int
		nt              int
		ct              int
	)

	maxc, err = conv2int(max_conn, 0)
	if err != nil {
		return nil, err
	}

	nt, err = conv2int(network_timeout, 1)
	if err != nil {
		return nil, err
	}

	ct, err = conv2int(connect_timeout, 2)
	if err != nil {
		return nil, err
	}

	minc, err = conv2int(min_conn, 3)
	if err != nil {
		return nil, err
	}
	for _, tr := range trackerList {
		var trackerIp string
		var trackerPort int
		tr = strings.TrimSpace(tr)
		parts := strings.Split(tr, ":")
		if len(parts) != 2 {
			return nil, errors.New("Wrong format with section 'tracker_server' of config file")
		}
		trackerIp = parts[0]

		//if len(parts) == 2 {
		trackerPort, err := strconv.Atoi(parts[1]) //append(trackerPort, parts[1])
		if err != nil {
			return nil, errors.New("Wrong format with section 'ip port' of config file")
		}

		/*} else if len(parts > 2) {

																																																																																																																}*/
		//if trackerIp != "" {
		trackerIpList = append(trackerIpList, trackerIp)
		//}

		//if trackerPort !=
		trackerPortList = append(trackerPortList, trackerPort)
	}

	Config := &Config{
		TrackerIp:   trackerIpList,
		TrackerPort: trackerPortList,
		MaxConn:     maxc,
		MinConn:     minc,
		Net_Timeout: nt,
		Con_Timeout: ct,
	}
	MAXCONN = maxc
	MINCONN = minc
	return Config, nil

}

func conv2int(stringname string, id int) (int, error) {
	i, err := strconv.Atoi(stringname)
	var name string
	switch id {
	case 0:
		name = "max_conn"
	case 1:
		name = "network_timeout"
	case 2:
		name = "connect_timeout"
	case 3:
		name = "min_conn"
	default:
		return -1, fmt.Errorf("wrong id number %d", id)
	}
	if err != nil {
		return -1, fmt.Errorf("Wrong format with section '%s' of config file", name)
	}
	return i, nil
}

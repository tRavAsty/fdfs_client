package fdfs_client

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"time"

	"github.com/weilaihui/goconfig/config"
)

const (
	base64Table = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_"
)

var coder = base64.NewEncoding(base64Table)

type Errno struct {
	status int
}

type fileInfo struct {
	createTimeStamp int32
	crc32           int32
	sourceId        int
	fileSize        int64
	sourceIpAddress string
}

func (e Errno) Error() string {
	errmsg := fmt.Sprintf("errno [%d] ", e.status)
	switch e.status {
	case 17:
		errmsg += "File Exist"
	case 22:
		errmsg += "Argument Invlid"
	}
	return errmsg
}

type FdfsConfigParser struct{}

var (
	ConfigFile *config.Config
)

func (this *FdfsConfigParser) Read(filename string) (*config.Config, error) {
	return config.ReadDefault(filename)
}

func fdfsCheckFile(filename string) error {
	if _, err := os.Stat(filename); err != nil {
		return err
	}
	return nil
}

func readCstr(buff io.Reader, length int) (string, error) {
	str := make([]byte, length)
	n, err := buff.Read(str)
	if err != nil || n != len(str) {
		return "", Errno{255}
	}

	for i, v := range str {
		if v == 0 {
			str = str[0:i]
			break
		}
	}
	return string(str), nil
}
func getFileExt(filename string) string {
	parts := strings.Split(filename, ".")
	if len(parts) >= 2 {
		return parts[len(parts)-1]
	}
	return ""
}

func splitRemoteFileId(remoteFileId string) ([]string, error) {
	parts := strings.SplitN(remoteFileId, "/", 2)
	if len(parts) < 2 {
		return nil, errors.New("error remoteFileId")
	}
	return parts, nil
}
func inet_ntoa(bytes []byte) (string, error) {
	if len(bytes) != 4 {
		return "", errors.New("error ip address")
	}

	return net.IPv4(bytes[0], bytes[1], bytes[2], bytes[3]).String(), nil
}

func (this *FdfsClient) getFileInfo(remotFileId string) (*fileInfo, error) {
	parts, err := splitRemoteFileId(remotFileId)
	if err != nil {
		return nil, err
	}
	fileLen := len(parts[1])
	if fileLen < FDFS_NORMAL_LOGIC_FILENAME_LENGTH {
		return nil, errors.New("error remoteFileName")
	}
	fileInfo := &fileInfo{}
	var buffer bytes.Buffer
	buffer.WriteString(parts[1][FDFS_LOGIC_FILE_PATH_LEN : FDFS_LOGIC_FILE_PATH_LEN+FDFS_FILENAME_BASE64_LENGTH])
	buffer.WriteString("=")
	logger.Info(buffer.String())
	decode, err := coder.DecodeString(buffer.String())
	if err != nil {
		return nil, err
	}
	ip, err := inet_ntoa(decode[:4])
	if err != nil {
		return nil, err
	}
	fileInfo.sourceIpAddress = ip

	b_buf := bytes.NewBuffer(decode[4:])

	var (
		createTimeStamp int32
		crc32           int32
		//sourceId        int
		fileSize int64
	)

	if err = binary.Read(b_buf, binary.BigEndian, &createTimeStamp); err != nil {
		logger.Error("1")
		return nil, err
	}
	if err = binary.Read(b_buf, binary.BigEndian, &fileSize); err != nil {
		logger.Error("2")
		return nil, err
	}
	if err = binary.Read(b_buf, binary.BigEndian, &crc32); err != nil {
		logger.Error("3")
		return nil, err
	}
	//logger.Infof("filesize:%ld", fileSize)
	if (fileSize >> 63) == 0 {
		fileSize &= 0x00000000FFFFFFFF
	}
	//logger.Infof("filesize:%ld", fileSize)
	if (fileSize >> 57) != 0 {
		logger.Info("appender file")
		return this.QueryFileInfo(parts[0], parts[1])

	} else {
		logger.Info("not appender file")
	}
	fileInfo.createTimeStamp = createTimeStamp
	fileInfo.crc32 = crc32
	return fileInfo, nil

}
func (fileInfo *fileInfo) Print() {

	logger.Info("createtime:" + time.Unix(int64(fileInfo.createTimeStamp), 0).String())
	logger.Infof("crc:%d", fileInfo.crc32)
	logger.Info("source ip:" + fileInfo.sourceIpAddress)
	logger.Infof("filesize:%d", fileInfo.fileSize)
}

package fdfs_client

import (
	"errors"
	"fmt"
	"github.com/Sirupsen/logrus"
	"os"
	//"strings"
	"testing"
	"time"
)

var (
	uploadResponse *UploadFileResponse
	deleteResponse *DeleteFileResponse
)

func TestParserFdfsConfig(t *testing.T) {
	fc := &FdfsConfigParser{}
	c, err := fc.Read("client.conf")
	if err != nil {
		t.Error(err)
		return
	}
	v, _ := c.String("DEFAULT", "base_path")
	t.Log(v)
}
func TestNewFdfsClientByTracker(t *testing.T) {

	tracker, err := getTrackerConf("client.conf")
	if err != nil {
		t.Error(err)
	}
	_, err = NewFdfsClientByTracker(tracker)
	if err != nil {
		t.Error(err)
	}
}

func TestUploadByFilename(t *testing.T) {
	logger.WithFields(logrus.Fields{
		"file":     "client_test.go",
		"function": "TestUploadByFilename",
	}).Info("Begin to upload by filename")
	//logger.Debug("upload by file name")
	fdfsClient, err := NewFdfsClient("client.conf")
	if err != nil {
		t.Errorf("New FdfsClient error %s", err.Error())
		return
	}

	uploadResponse, err = fdfsClient.UploadByFilename("x0")
	if err != nil {
		t.Errorf("UploadByfilename error %s", err.Error())
	}
	t.Log(uploadResponse.GroupName)
	t.Log(uploadResponse.RemoteFileId)
	fdfsClient.DeleteFile(uploadResponse.RemoteFileId)
}

func TestUploadByBuffer(t *testing.T) {
	fdfsClient, err := NewFdfsClient("client.conf")
	if err != nil {
		t.Errorf("New FdfsClient error %s", err.Error())
		return
	}

	file, err := os.Open("testfile") // For read access.
	if err != nil {
		t.Fatal(err)
	}

	var fileSize int64 = 0
	if fileInfo, err := file.Stat(); err == nil {
		fileSize = fileInfo.Size()
	}
	fileBuffer := make([]byte, fileSize)
	_, err = file.Read(fileBuffer)
	if err != nil {
		t.Fatal(err)
	}

	uploadResponse, err = fdfsClient.UploadByBuffer(fileBuffer, "txt")
	if err != nil {
		t.Errorf("TestUploadByBuffer error %s", err.Error())
	}

	t.Log(uploadResponse.GroupName)
	t.Log(uploadResponse.RemoteFileId)
	fdfsClient.DeleteFile(uploadResponse.RemoteFileId)
}

func TestUploadSlaveByFilename(t *testing.T) {
	fdfsClient, err := NewFdfsClient("client.conf")
	if err != nil {
		t.Errorf("New FdfsClient error %s", err.Error())
		return
	}

	uploadResponse, err = fdfsClient.UploadByFilename("client.conf")
	if err != nil {
		t.Errorf("UploadByfilename error %s", err.Error())
	}
	t.Log(uploadResponse.GroupName)
	t.Log(uploadResponse.RemoteFileId)

	masterFileId := uploadResponse.RemoteFileId
	uploadResponse, err = fdfsClient.UploadSlaveByFilename("testfile", masterFileId, "_test")
	if err != nil {
		t.Errorf("UploadByfilename error %s", err.Error())
	}
	t.Log(uploadResponse.GroupName)
	t.Log(uploadResponse.RemoteFileId)

	fdfsClient.DeleteFile(masterFileId)
	fdfsClient.DeleteFile(uploadResponse.RemoteFileId)
}

func TestDownloadToFile(t *testing.T) {
	fdfsClient, err := NewFdfsClient("client.conf")
	if err != nil {
		t.Errorf("New FdfsClient error %s", err.Error())
		return
	}

	uploadResponse, err = fdfsClient.UploadByFilename("/usr/include/stdlib.h")
	defer fdfsClient.DeleteFile(uploadResponse.RemoteFileId)
	if err != nil {
		t.Errorf("UploadByfilename error %s", err.Error())
	}
	t.Log(uploadResponse.GroupName)
	t.Log(uploadResponse.RemoteFileId)

	var (
		downloadResponse *DownloadFileResponse
		localFilename    string = "download.txt"
	)
	downloadResponse, err = fdfsClient.DownloadToFile(localFilename, uploadResponse.RemoteFileId, 0, 0)
	if err != nil {
		t.Errorf("DownloadToFile error %s", err.Error())
	}
	t.Log(downloadResponse.DownloadSize)
	t.Log(downloadResponse.RemoteFileId)
}

func TestDownloadToBuffer(t *testing.T) {
	fdfsClient, err := NewFdfsClient("client.conf")
	if err != nil {
		t.Errorf("New FdfsClient error %s", err.Error())
		return
	}

	uploadResponse, err = fdfsClient.UploadByFilename("client.conf")
	defer fdfsClient.DeleteFile(uploadResponse.RemoteFileId)
	if err != nil {
		t.Errorf("UploadByfilename error %s", err.Error())
	}
	t.Log(uploadResponse.GroupName)
	t.Log(uploadResponse.RemoteFileId)

	var (
		downloadResponse *DownloadFileResponse
	)
	downloadResponse, err = fdfsClient.DownloadToBuffer(uploadResponse.RemoteFileId, 0, 0)
	if err != nil {
		t.Errorf("DownloadToBuffer error %s", err.Error())
	}
	t.Log(downloadResponse.DownloadSize)
	t.Log(downloadResponse.RemoteFileId)
}

func BenchmarkUploadByBuffer(b *testing.B) {
	fdfsClient, err := NewFdfsClient("client.conf")
	if err != nil {
		fmt.Errorf("New FdfsClient error %s", err.Error())
		return
	}
	file, err := os.Open("testfile") // For read access.
	if err != nil {
		fmt.Errorf("%s", err.Error())
	}

	var fileSize int64 = 0
	if fileInfo, err := file.Stat(); err == nil {
		fileSize = fileInfo.Size()
	}
	fileBuffer := make([]byte, fileSize)
	_, err = file.Read(fileBuffer)
	if err != nil {
		fmt.Errorf("%s", err.Error())
	}

	b.StopTimer()
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		uploadResponse, err = fdfsClient.UploadByBuffer(fileBuffer, "txt")
		if err != nil {
			fmt.Errorf("TestUploadByBuffer error %s", err.Error())
		}

		fdfsClient.DeleteFile(uploadResponse.RemoteFileId)
	}
}

func BenchmarkUploadByFilename(b *testing.B) {
	fdfsClient, err := NewFdfsClient("client.conf")
	if err != nil {
		fmt.Errorf("New FdfsClient error %s", err.Error())
		return
	}

	b.StopTimer()
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		uploadResponse, err = fdfsClient.UploadByFilename("client.conf")
		if err != nil {
			fmt.Errorf("UploadByfilename error %s", err.Error())
		}
		_, err = fdfsClient.DeleteFile(uploadResponse.RemoteFileId)
		if err != nil {
			fmt.Errorf("DeleteFile error %s", err.Error())
		}
	}
}

func BenchmarkDownloadToFile(b *testing.B) {
	fdfsClient, err := NewFdfsClient("client.conf")
	if err != nil {
		fmt.Errorf("New FdfsClient error %s", err.Error())
		return
	}

	uploadResponse, err = fdfsClient.UploadByFilename("client.conf")
	defer fdfsClient.DeleteFile(uploadResponse.RemoteFileId)
	if err != nil {
		fmt.Errorf("UploadByfilename error %s", err.Error())
	}
	b.StopTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		var (
			localFilename string = "download.txt"
		)
		_, err = fdfsClient.DownloadToFile(localFilename, uploadResponse.RemoteFileId, 0, 0)
		if err != nil {
			fmt.Errorf("DownloadToFile error %s", err.Error())
		}

		// fmt.Println(downloadResponse.RemoteFileId)
	}
}

func TestUploadAppenderByFilename(t *testing.T) {
	/*logger.WithFields(logrus.Fields{
		"animal": "walrus",
		"size":   20,
	}).Info("A group of walrus emerges from the ocean")*/
	//logger.Debug("upload by file name")
	fdfsClient, err := NewFdfsClient("client.conf")
	if err != nil {
		t.Errorf("New FdfsClient error %s", err.Error())
		return
	}

	uploadResponse, err = fdfsClient.UploadAppenderByFilename("x1")
	//uploadResponse, err = fdfsClient.UploadByFilename("x1")
	if err != nil {
		t.Errorf("AppendByfilename error %s", err.Error())
	}

	t.Log(uploadResponse.GroupName)
	t.Log(uploadResponse.RemoteFileId)
	parts, err := splitRemoteFileId(uploadResponse.RemoteFileId)
	if err != nil {
		t.Error(err)
	}
	groupName := parts[0]
	remoteFileName := parts[1]

	fileInfo, err := fdfsClient.getFileInfo(uploadResponse.RemoteFileId)
	if err != nil {
		t.Error("get file info error" + err.Error())
	}

	fileInfo.Print()
	fileSize := fileInfo.fileSize
	//group1/M00/00/03/wKj_glc-fQiEISCUAAAAAChSBpE4174280
	if deleteResponse, err = fdfsClient.TruncAppenderByFilename(uploadResponse.RemoteFileId, fileSize/2); err != nil {
		t.Errorf("Truncate Appender File error %s", err.Error())
	}
	t.Log(deleteResponse.groupName)
	t.Log(deleteResponse.remoteFilename)
	fileInfo, err = fdfsClient.getFileInfo(uploadResponse.RemoteFileId)
	if err != nil {
		t.Error("get file info error" + err.Error())
	}

	fileInfo.Print()
	if fileInfo.fileSize != fileSize/2 {
		t.Errorf("filesize:%d != %d", fileInfo.fileSize, fileSize/2)
	}

	if err = fdfsClient.AppendByFileName("x1", groupName, remoteFileName); err != nil {
		t.Error("can't append file")
	}
	fileInfo, err = fdfsClient.getFileInfo(uploadResponse.RemoteFileId)
	if err != nil {
		t.Error("get file info error" + err.Error())
	}
	fileInfo.Print()

	offset := fileInfo.fileSize
	if err = fdfsClient.ModifyByFileName("x1", offset, groupName, remoteFileName); err != nil {
		t.Error("can't modify file")
	}

	fileInfo, err = fdfsClient.getFileInfo(uploadResponse.RemoteFileId)
	if err != nil {
		t.Error("get file info error" + err.Error())
	}
	fileInfo.Print()
}
func TestDeleteFile(t *testing.T) {
	fdfsClient, err := NewFdfsClient("client.conf")
	if err != nil {
		t.Errorf("New FdfsClient error %s", err.Error())
		return
	}

	uploadResponse, err = fdfsClient.UploadByFilename("x3")
	if err != nil {
		t.Errorf("UploadByfilename error %s", err.Error())
	}
	logger.Info(uploadResponse.GroupName)
	logger.Info(uploadResponse.RemoteFileId)
	deleteResponse, err = fdfsClient.DeleteFile(uploadResponse.RemoteFileId)
	if err != nil {
		t.Error(err)
	}
	t.Log(deleteResponse.groupName)
	t.Log(deleteResponse.remoteFilename)

}

func TestCrackTracker(t *testing.T) {
	fdfsClient, err := NewFdfsClient("client.conf")
	if err != nil {
		t.Errorf("New FdfsClient error %s", err.Error())
		return
	}

	uploadResponse, err = fdfsClient.UploadByFilename("/usr/include/stdlib.h")
	defer fdfsClient.DeleteFile(uploadResponse.RemoteFileId)
	if err != nil {
		t.Errorf("UploadByfilename error %s", err.Error())
	}
	t.Log(uploadResponse.GroupName)
	t.Log(uploadResponse.RemoteFileId)

	time.Sleep(5 * time.Minute)

	var (
		downloadResponse *DownloadFileResponse
		localFilename    string = "download.txt"
	)
	downloadResponse, err = fdfsClient.DownloadToFile(localFilename, uploadResponse.RemoteFileId, 0, 0)
	if err != nil {
		t.Errorf("DownloadToFile error %s", err.Error())
	}
	t.Log(downloadResponse.DownloadSize)
	t.Log(downloadResponse.RemoteFileId)
}
func TestCrackStorage(t *testing.T) {
	fdfsClient, err := NewFdfsClient("client.conf")
	if err != nil {
		t.Errorf("New FdfsClient error %s", err.Error())
		return
	}

	uploadResponse, err = fdfsClient.UploadByFilename("/usr/include/stdlib.h")
	defer fdfsClient.DeleteFile(uploadResponse.RemoteFileId)
	if err != nil {
		t.Errorf("UploadByfilename error %s", err.Error())
	}
	t.Log(uploadResponse.GroupName)
	t.Log(uploadResponse.RemoteFileId)

	time.Sleep(5*time.Minute + 3*time.Second)
	t.Error(errors.New("can't download file"))
	/*var (
		downloadResponse *DownloadFileResponse
		localFilename    string = "download.txt"
	)
	downloadResponse, err = fdfsClient.DownloadToFile(localFilename, uploadResponse.RemoteFileId, 0, 0)
	if err != nil {
		t.Errorf("DownloadToFile error %s", err.Error())
	}
	t.Log(downloadResponse.DownloadSize)
	t.Log(downloadResponse.RemoteFileId)*/
}

var ch chan int = make(chan int)

func Test10Upload(t *testing.T) {
	fdfsClient, err := NewFdfsClient("client.conf")
	if err != nil {
		t.Errorf("New FdfsClient error %s", err.Error())
		return
	}
	for i := 0; i < 10; i++ {
		go fdfsClient.conUpload("x0", ch)
	}
	for i := 0; i < 10; i++ {
		t.Log("err=", <-ch)
	}
}

func (fdfsclient *FdfsClient) conUpload(filename string, ch chan int) {
	uploadResponse, err := fdfsclient.UploadByFilename(filename)
	if err != nil {
		ch <- -1
	} else {
		ch <- 0
	}
	logger.Info(uploadResponse.GroupName)
	logger.Info(uploadResponse.RemoteFileId)
}

/*func (fdfsclient *FdfsClient) conDownload(localFilename string, ch chan int) {
	var (
		downloadResponse *DownloadFileResponse
	)
	downloadResponse, err := fdfsclient.DownloadToFile(localFilename, "group1/M00/00/01/wKiWhFc8dFWAAQD3AAAoAChSBpE9021000", 0, 0)
	if err != nil {
		ch <- -1
	} else {
		ch <- 0
	}
	logger.Info(downloadResponse.DownloadSize)
	logger.Info(downloadResponse.RemoteFileId)
}*/

func Test10Download(t *testing.T) {
	fdfsClient, err := NewFdfsClient("client.conf")
	if err != nil {
		t.Errorf("New FdfsClient error %s", err.Error())
		return
	}

	var (
		downloadResponse *DownloadFileResponse
		localFilename    string = "x11"
	)
	for i := 0; i < 10; i++ {
		downloadResponse, err = fdfsClient.DownloadToFile(localFilename, "group1/M00/00/01/wKiWhFc8dFWAAQD3AAAoAChSBpE9021000", 0, 0)
		if err != nil {
			t.Error()
		}
		logger.Info(downloadResponse.DownloadSize)
		logger.Info(downloadResponse.RemoteFileId)

	}

}

func TestGetFileInfo(t *testing.T) {
	fdfsClient, err := NewFdfsClient("client.conf")
	if err != nil {
		t.Errorf("New FdfsClient error %s", err.Error())
		return
	}

	fileInfo, err := fdfsClient.getFileInfo("group1/M00/00/03/wKj_glc-fQiEISCUAAAAAChSBpE4174280")
	if err != nil {
		t.Error("get file info error" + err.Error())
	}
	logger.Info("createtime:" + time.Unix(int64(fileInfo.createTimeStamp), 0).String())
	logger.Infof("crc:%d", fileInfo.crc32)
	logger.Info("source ip:" + fileInfo.sourceIpAddress)
	logger.Infof("filesize:%d", fileInfo.fileSize)
}

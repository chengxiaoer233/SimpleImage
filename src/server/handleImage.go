package server

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

// request struct
type ReqImageAnalyze struct {
	Url      string
	FilePath string
}

// resp struct
type RespImageStruct struct {
	Normal               bool  `json:"normal"`
	ReWrite              bool  `json:"reWrite"`
	TotalContentLength   int64 `json:"totalContentLength"`
	ReWriteContentLength int64 `json:"reWriteContentLength"`
}

// image analyze,can check file from local or remote
func HandleAnalyzeImage(ctx context.Context, req *ReqImageAnalyze) (resp RespImageStruct, err error) {

	// get file tmpBuf and totalBuf
	var contentTypeBuf, totalBuf []byte
	if req.Url != "" { // from local
		contentTypeBuf, totalBuf, err = getFileFromRemote(req.Url)

	} else { // from remote
		contentTypeBuf, totalBuf, err = getFileFromLocal(req.FilePath)
	}

	// get file content-type
	contentType := http.DetectContentType(contentTypeBuf)

	// analyze
	resp, err = analyze(contentType, totalBuf)
	data, err := json.Marshal(resp)

	fmt.Println("HandleImageCheck resp=", string(data), ",err=", err)

	return
}

// load file from local
func getFileFromLocal(filePath string) (tmpBytes, totalBytes []byte, err error) {

	// open file and close
	file, err := os.Open(filePath)
	if err != nil {
		panic(fmt.Sprintf("os.Open error,err%s", err))
	}
	defer file.Close()

	// read tmpBuf and totalBuf
	// tmpBuf is used to build file content-type
	tmpBytes = make([]byte, 521)
	_, err = file.Read(tmpBytes)
	if err != nil {
		return tmpBytes, totalBytes, err
	}

	// read totalBuf
	totalBytes = append(totalBytes, tmpBytes...)
	for {

		tmp := make([]byte, 4096)

		n, err := file.Read(tmp)
		if err != nil && err != io.EOF {
			return tmpBytes, totalBytes, err
		} else if err == io.EOF {
			totalBytes = append(totalBytes, tmp[:n]...)
			break
		}

		totalBytes = append(totalBytes, tmp[:n]...)
	}

	return
}

// load file from remote
func getFileFromRemote(url string) (tmpBytes, totalBytes []byte, err error) {

	// download file from remote
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	reader := bufio.NewReader(resp.Body)

	// read tmpBuf
	tmpBytes = make([]byte, 512)
	_, err = reader.Read(tmpBytes)
	if err != nil {
		return tmpBytes, totalBytes, err
	}

	// read totalBuf
	totalBytes = append(totalBytes, tmpBytes...)
	for {
		tmp := make([]byte, 4096)
		n, err := reader.Read(tmp)

		if err != nil && err != io.EOF { // other error
			return tmpBytes, totalBytes, err
		} else if err == io.EOF { // EOF
			totalBytes = append(totalBytes, tmp[:n]...)
			break
		}

		totalBytes = append(totalBytes, tmp[:n]...)
	}

	// ioutil.ReadAll
	/*tmp,err := ioutil.ReadAll(resp.Body)
	if err != nil{
		return tmpBytes, totalBytes, err
	}
	totalBytes = append(totalBytes,tmp...)*/

	return
}

// analyze detail of file
func analyze(contentType string, totalBuf []byte) (resp RespImageStruct, err error) {

	var image image.Image

	bytesReader := bytes.NewReader(totalBuf)
	resp.TotalContentLength = int64(bytesReader.Len())

	switch contentType {
	case "image/png":
		image, err = png.Decode(bytesReader)
	case "image/jpeg":
		image, err = jpeg.Decode(bytesReader)
	case "image/gif":
		image, err = gif.Decode(bytesReader)
	default:
		return
	}

	// decode error,return
	if err != nil {
		fmt.Println("decode error,err=", err)
		resp.Normal = false
		return
	}

	_ = image

	// read left data
	left, err := ioutil.ReadAll(bytesReader)
	if err != nil {
		resp.Normal = false
		resp.TotalContentLength = int64(bytesReader.Len())
		return
	}

	// resp the result
	if len(left) == 0 {
		resp.Normal = true
		resp.ReWrite = false
		resp.ReWriteContentLength = 0

		return
	}

	// have been reWrite
	resp.Normal = false
	resp.ReWrite = true
	resp.ReWriteContentLength = int64(len(left))

	// build a ts file
	tsDir := "./build.ts"
	tsFile, err := os.OpenFile("./build.ts", os.O_TRUNC|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		panic(fmt.Sprintf("os.OpenFile err=%s", err))
	}
	defer tsFile.Close()

	// write
	_, err = tsFile.Write(left)
	if err != nil {
		fmt.Println("tsFile.Write failed,err=",err)
	}

	fmt.Println("tsFile.Write success,dir=",tsDir,",len=",len(left))
	return
}

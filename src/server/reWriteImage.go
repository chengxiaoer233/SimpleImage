package server

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

func RewriteImage() {

	reWritePng()

//	reWriteJpeg()

//	reWriteGif()
}

// rewrite ts into png
func reWritePng() {

	// build tmp file
	src := "./etc/data/base/base_367_bytes.png"
	ts := "./etc/data/base/1.ts"
	dst := "./etc/data/build/png-ts.png"

	reWrite(dst, src, ts)
}

// rewrite ts into jpeg
func reWriteJpeg() {

	// build tmp file
	src := "./etc/data/base/base_367_bytes.jpeg"
	ts := "./etc/data/base/1.ts"
	dst := "./etc/data/build/png-ts.jpeg"

	reWrite(dst, src, ts)
}

// rewrite ts into gif
func reWriteGif() {

	// build tmp file
	src := "./etc/data/base/base_367_bytes.gif"
	ts := "./etc/data/base/1.ts"
	dst := "./etc/data/build/png-ts.gif"

	reWrite(dst, src, ts)
}

// rewrite
func reWrite(dst, src, ts string) (err error) {

	// file copy,base image to tmp image
	copyFile(dst, src)

	// read a ts file
	tsFile, err := os.Open(ts)
	if err != nil {
		panic(fmt.Sprintf("os.Open error =%s", err))
	}
	defer tsFile.Close()

	// readAll
	tsBuf, err := ioutil.ReadAll(tsFile)
	if err != nil {
		fmt.Println("ioutil.ReadAll error,err=", err)
		return err
	}

	// reWrite：insert a ts file into tmp png
	file, err := os.OpenFile(dst, os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		panic(fmt.Sprintf("os.OpenFile error =%s", err))
	}
	defer file.Close()

	// reWrite
	_, err = file.Write(tsBuf)
	if err != nil {
		fmt.Println("file.Write error,err=", err)
		return err
	}

	fmt.Println("reWrite success,dst=", dst, ",src=", src, ",len(ts)=", len(tsBuf))
	return err
}

// file copy
func copyFile(dst string, src string) {

	srcFile, err := os.Open(src)
	if err != nil {
		panic(fmt.Sprintf(" open file error,err=%s", err))
	}
	defer srcFile.Close()

	dstFile, err := os.OpenFile(dst, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0666)
	if err != nil {
		panic(fmt.Sprintf(" os.OpenFile error,err=%s", err))
	}
	defer dstFile.Close()

	io.Copy(dstFile, srcFile)
}

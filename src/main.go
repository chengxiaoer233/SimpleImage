package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"simpleImage/server"
	"syscall"
)

func main() {

	var req server.ReqImageAnalyze
	//req.Url = "https://static.yximgs.com/udata/pkg/test-wuxiaolong/test-new/1Group_1312317608111111.png"

	 req.FilePath = "./etc/data/build/png-ts.png"
	// req.FilePath = "./etc/data/build/png-ts.png"
	// req.FilePath = "./etc/data/build/png-ts.png"
	// req.FilePath = "./etc/data/build/png-ts.png"

	// analyze
	 server.HandleAnalyzeImage(context.Background(), &req)

	// reWrite
	 // server.RewriteImage()

	// wait for return
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGALRM)
	<-c
	fmt.Println("bye")
}

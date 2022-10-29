<div align="center">

# 一个简易go语言图片防盗链检测器
</div>


<div align="center">

![](https://img.shields.io/github/languages/code-size/chengxiaoer233/simpleImage?label=CodeSize)
![](https://img.shields.io/github/stars/chengxiaoer233/simpleImage?label=GitHub)
![](https://img.shields.io/github/watchers/chengxiaoer233/simpleImage?label=Watch)
[![Go Report Card](https://goreportcard.com/badge/github.com/chengxiaoer233/simpleImage)](https://goreportcard.com/report/github.com/chengxiaoer233/simpleImage)
[![LICENSE](https://img.shields.io/badge/license-MIT-green)](https://mit-license.org/)
</div>


<div align="center">

<img  src="https://my-source666.obs.cn-south-1.myhuaweicloud.com/myBlog/golang-jixiangwu-image.png" width="600" height="350"/>

</div>


#### 实现的功能
- 1：检测图片类型

- 2：检查图片是否被修改（图片中添加了其他数据，主要用于实现盗链，用别人的存储当做图床）

- 3：修改原始图片内容，插入ts文件，模拟黑产上传至**网（上传过程忽略，主要是用于学习）

- 4：解析被修改后的图片，从中提取出原始插入的ts文件，并进行转码，ts文件转MP4

#### 函数分析  

##### **RewriteImage：** 

重写图片，在原始的图片后面加入自己的内容，可以伪造成图片，上传到对方的服务器，把对方的服务器当成图床,
下载的时候，再解码出自己写入的文件。  

黑产常见的手段：在图片中插入视频ts片段，服务器一般不会校验数据是否被篡改，上传成功后，黑产拿到对应的url，
下载改图片解码出对应的ts文件，用自己的播放器播放即可。  

1. 函数讲解：
    - 1: ./etc/data/base下为原始图片，未篡改内容
    
    - 2: ./etc/data/ts下面的ts文件为需要我们插入的1.ts文件，该ts文件对应的视频文件为1.mp4
   
    - 3: 先基于原始base图片，生成一个待写入的tmp ts文件
    
    - 4: 读取ts文件内容，将ts文件插入到tmp ts文件，并保存
    
    - 5: 上传改ts文件至别人的服务器，例如：**的图片评论系统，获取图片地址（这里就不演示了）
    
    - 6: 通过该url下载资源，解析出写入的原始ts内容，将ts内容进行转码，转成mp4文件，
    就可以播放了（参照HandleAnalyzeImage函数）

##### **HandleAnalyzeImage：**  

1. 函数讲解：
    - 1: 根据携带的参数读取本地或远端的url资源
    
    - 2: 分别读取文件的前512个字节和总字节数
   
    - 3: 根据前面的512个字节数，判断该文件的content-type
    
    - 4: 根据不同的content-type调用不同的decode函数
    
    - 5: 读取decode后的数据，判断是否有被修改
    
    - 6: 将多余的ts数据写入到临时tmp文件，进行ts转mp4格式转换

#### 举例说明 :
* reWrite：伪造修改原始png图片
```text

  （1）调用reWritePng()方法,
  （2）将./etc/data/ts/1.ts中的ts写入到./etc/data/ts/png-ts.png
    
    结果输出：
    reWrite success,
        dst= ./etc/data/ts/png-ts.png ,
        src= ./etc/data/base/base1_367bytes.png ,
        len(ts)= 1474674
    
    可以看到ts文件已经被插入到临时的png-ts.png
```


* analyze：下载解析修改后的png图片,并还原出插入的ts文件，进行格式转换为mp4文件
```text

  （1）req.FilePath = "./etc/data/ts/png-ts.png"
  （2）调用server.HandleAnalyzeImage
    
   结果输出：
     HandleImageCheck resp= 
     {
       "normal": false,               // 解析异常
       "reWrite": true,               // 图片被改写了
       "totalContentLength": 1475041, // 图片总大小
       "reWriteContentLength": 1474674 // 插入数据(ts)的大小 
     } ,
    err= <nil>

   可以看到插入的ts文件被正常的解析和还原下来，转码后就可以播放了

``` 
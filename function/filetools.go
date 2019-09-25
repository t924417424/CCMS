package function

import (
	"os"
	"fmt"
	"bytes"
	"strings"
	"io/ioutil"
	"crypto/des"
	"crypto/md5"
	"encoding/hex"
	"crypto/cipher"
)

func Padding(src []byte,blocksize int) []byte {
	padnum:=blocksize-len(src)%blocksize
	pad:=bytes.Repeat([]byte{byte(padnum)},padnum)
	return append(src,pad...)
}

func Unpadding(src []byte) []byte {
	n:=len(src)
	unpadnum:=int(src[n-1])
	return src[:n-unpadnum]
}

func Encrypt3DES(src []byte,key []byte) []byte {
	block,_:=des.NewTripleDESCipher(key)
	src=Padding(src,block.BlockSize())
	blockmode:=cipher.NewCBCEncrypter(block,key[:block.BlockSize()])
	blockmode.CryptBlocks(src,src)
	return src
}

func Decrypt3DES(src []byte,key []byte) []byte {
	block,_:=des.NewTripleDESCipher(key)
	blockmode:=cipher.NewCBCDecrypter(block,key[:block.BlockSize()])
	blockmode.CryptBlocks(src,src)
	src=Unpadding(src)
	return src
}

func ReadAllIntoMemory(filename string) (content []byte,err error) {

	fp, err := os.Open(filename) // 获取文件指针
	if err != nil {
		return nil, err
	}
	defer fp.Close()

	fileInfo, err := fp.Stat()
	if err != nil {
		return nil, err
	}
	buffer := make([]byte, fileInfo.Size())
	_, err = fp.Read(buffer) // 文件内容读取到buffer中
	if err != nil {
		return nil, err
	}
	return buffer, nil
}
func WriteWithFileWrite(name,content string){
	fileObj,err := os.OpenFile(name,os.O_RDWR|os.O_CREATE|os.O_TRUNC,0644)
	if err != nil {
		//fmt.Println("Failed to open the file",err.Error())
		//os.Exit(2)
	}
	defer fileObj.Close()
	if _,err := fileObj.WriteString(content);err == nil {
		//fmt.Println("Successful writing to the file with os.OpenFile and *File.WriteString method.")
	}
	contents := []byte(content)
	if _,err := fileObj.Write(contents);err == nil {
		//fmt.Println("Successful writing to thr file with os.OpenFile and *File.Write method.")
	}
}

func WriteWithIoutil(name,content string) {
	data :=  []byte(content)
	if ioutil.WriteFile(name,data,0644) == nil {
		fmt.Println("写入文件成功")
	}
}

func GetBetweenStr(str, start, end string) string {
	n := strings.Index(str, start)
	if n == -1 {
		n = 0
	} else {
		n = n + len(start)  // 增加了else，不加的会把start带上
	}
	str = string([]byte(str)[n:])
	m := strings.Index(str, end)
	if m == -1 {
		m = len(str)
	}
	str = string([]byte(str)[:m])
	return str
}

func IsExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
	// 或者
	//return err == nil || !os.IsNotExist(err)
	// 或者
	//return !os.IsNotExist(err)
}
func Md5V(str string) string  {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}
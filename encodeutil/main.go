package main

import (
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"unicode/utf16"
)

func main() {
	port := 8853
	Args := os.Args
	if len(Args) >= 2 {
		arg1 := Args[1]
		port, e := strconv.Atoi(arg1)
		if e != nil || port <= 0 || port > 65535 {
			log.Println("端口错误，默认使用8853端口")
			port = 8853
		}
	}
	fmt.Printf("打开浏览器输入localhost:%d访问\n", port)
	fmt.Println("在此期间请勿关闭此窗口")
	fmt.Println("可以通过参数形式指定端口号")
	dir, _ := os.Getwd()
	http.Handle("/", http.FileServer(http.Dir(dir)))
	http.HandleFunc("/content", func(writer http.ResponseWriter, request *http.Request) {
		m := make(map[string]struct{})
		var build strings.Builder
		build.WriteString("{")

		value := request.FormValue("txt")
		newLine := 0
		for _, v := range value {
			if _, ok := m[string(v)]; ok {
				continue
			}
			m[string(v)] = struct{}{}

			//gbk编码
			ss, _ := simplifiedchinese.GBK.NewEncoder().Bytes([]byte(string(v)))
			gbk := "0x"
			for i := len(ss) - 1; i >= 0; i-- {
				x := strconv.FormatInt(int64(ss[i]), 16)
				gbk = gbk + x
			}

			encode := utf16.Encode([]rune(string(v)))
			utf := "0x" + strconv.FormatInt(int64(encode[0]), 16)
			build.WriteString("{" + gbk + "," + utf + "},")
			newLine++
			if newLine == 4 {
				build.WriteByte('\n')
				newLine = 0
			}
		}
		r := build.String()
		r = r[0:len(r)-1] + "}"
		writer.Write([]byte(r))
	})
	err := http.ListenAndServe("localhost:"+strconv.Itoa(port), nil)
	if err != nil {
		fmt.Print(err)
	}
}

type Value struct {
}

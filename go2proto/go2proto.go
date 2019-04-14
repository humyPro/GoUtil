package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

var writer *bufio.Writer
var reader *bufio.Reader
var (
	startOfStruct = regexp.MustCompile(`^type +(\w+) +struct {$`)
	pack          = regexp.MustCompile(`^package +(\w+)$`)
	field1        = regexp.MustCompile(`^(\w+) +(\w+) *(//.+)?$`) //单独写的
	field2        = regexp.MustCompile(`^(?P<head>( *\w+, *)+) *(?P<last>\w+) +(?P<type>\w+) *(//.+)?$`)
)
//GO类型和proto类型
var typeMap = map[string]string{
	"float64": "double",
	"float32": "float",
	"int8":    "int32",
	"int16":   "int32",
	"int32":   "int32",
	"int":     "int32",
	"int64":   "int64",
	"uint8":   "uint32",
	"uint16":  "uint32",
	"uint32":  "uint32",
	"unit64":  "unit64",
	"bool":    "bool",
	"string":  "string",
	//todo  other type map
}

var packageFlag = true
var startEndFlag = false
var lineCount = 0
var fieldNumber = 1

//const (
//	sourceFileName = "./oth/model.go"
//	targetFileName = "./proto/model.proto"
//)

var goFilePath = flag.String("g", "", "指定需要编译的go文件")
var protoFilePath = flag.String("p", "./proto", "指定存放proto文件的位置，默认在本路径下的proto目录下")

// 读取一行
func nextLine() (string, bool) {
	s, e := reader.ReadString('\n')
	s = strings.TrimSpace(s)
	lineCount++
	if e != nil && e.Error() == "EOF" {
		return s, false
	}
	if len(s) == 0 || strings.HasPrefix(s, "/") || strings.HasPrefix(s, "*") {
		return nextLine()
	}
	return s, true
}

func genProtoFile() {
	for {
		s, b := nextLine()
		// do something
		if packageFlag {
			handPackage(s)
			continue
		}
		handStruct(s)
		if !b {
			break
		}
	}
	log.Println("生成proto文件成功")
}

func handStruct(s string) {
	if s == "" {
		return
	}
	strArr := startOfStruct.FindStringSubmatch(s)
	// 匹配输出message name {
	if len(strArr) == 2 {
		if startEndFlag {
			err(s)
		}
		startEndFlag = true
		writer.WriteString("message " + strArr[1] + " {\n")
		writer.Flush()
		return
	}
	if !startEndFlag {
		err(s)
	}
	strArr = field1.FindStringSubmatch(s)
	if len(strArr) == 4 {
		f := strArr[1]
		t := strArr[2]
		desc := strArr[3]
		writer.WriteString("\t" + typeMap[t] + " " + f + " = " + strconv.Itoa(fieldNumber) + ";  " + desc + "\n")
		fieldNumber++
		writer.Flush()
		return
	}

	strArr = field2.FindStringSubmatch(s)
	if len(strArr) == 6 {
		fields := strings.Split(strArr[1], ",")
		fields = fields[:len(fields)-1]
		last := strArr[3]
		t := strArr[4]
		desc := strArr[5]
		for _, f := range fields {
			writer.WriteString("\t" + typeMap[t] + " " + strings.TrimSpace(f) + " = " + strconv.Itoa(fieldNumber) + ";  " + desc + "\n")
			fieldNumber++
		}
		writer.WriteString("\t" + typeMap[t] + " " + strings.TrimSpace(last) + " = " + strconv.Itoa(fieldNumber) + ";  " + desc + "\n")
		fieldNumber++
		writer.Flush()
		return
	}

	if strings.TrimSpace(s) == "}" {
		writer.WriteString(s + "\n\n")
		writer.Flush()
		fieldNumber = 1
		startEndFlag = false
		return
	}

	if s == "gorm.Model" {
		writer.WriteString("\tint32 id = " + strconv.Itoa(fieldNumber) + ";\n")
		fieldNumber++
		writer.WriteString("\tint64 createdAt = " + strconv.Itoa(fieldNumber) + ";\n")
		fieldNumber++
		writer.WriteString("\tint64 updatedAt = " + strconv.Itoa(fieldNumber) + ";\n")
		fieldNumber++
		writer.WriteString("\tint64 deletedAt = " + strconv.Itoa(fieldNumber) + ";\n")
		fieldNumber++
	}
	err(s)

}
func err(s string) {
	s = fmt.Sprintf("%d:语法错误,%s\n", lineCount, s)
	writer.WriteString(s)
	writer.Flush()
	log.Fatal(s)
}

func handPackage(s string) {
	if packageFlag {
		ss := pack.FindStringSubmatch(s)
		if len(ss) == 2 {
			writer.WriteString("package " + ss[1] + ";\n")
			writer.Flush()
		}
		packageFlag = false
	} else {
		if strings.HasPrefix(s, "package") {
			fmt.Errorf("%d行:语法错误，出现多个package定义: %s", lineCount, s)
			os.Exit(2)
		}
	}
}

func main() {
	writer.WriteString(`syntax = "proto3";` + "\n")
	writer.Flush()
	genProtoFile()
}

//初始化reader writer
func init() {

	flag.Parse()
	if *goFilePath == "" {
		log.Fatal("必须指定go文件")
	}

	if !exists(*goFilePath) {
		log.Fatal(*goFilePath, ": 目标go文件不存在")
	}

	file, e := os.Open(*goFilePath)
	if e != nil {
		fmt.Println(e)
		os.Exit(1)
	}
	reader = bufio.NewReader(file)

	f := strings.TrimSpace(*protoFilePath) + "/" + strings.Split(filepath.Base(*goFilePath), ".")[0] + ".proto"
	if exists(f) {
		log.Fatal(f, ": 目标proto文件已存在")
	}

	if !exists(*protoFilePath) {
		e := os.Mkdir(*protoFilePath, 777)
		if e != nil {
			log.Fatal(e, ":创建文件夹失败")
		}
	}
	open, e := os.OpenFile(f, os.O_RDWR|os.O_CREATE, 0766)
	if e != nil {
		log.Fatal(e)
	}
	writer = bufio.NewWriter(open)
}

// 判断所给路径文件/文件夹是否存在
func exists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

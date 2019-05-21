package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
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
	field1        = regexp.MustCompile(`^(\w+) +(\w+) *(.+)?$`) //单独写的
	field2        = regexp.MustCompile(`^(?P<head>( *\w+, *)+) *(?P<last>\w+) +(?P<type>\w+) *(.+)?$`)
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
	"uint":    "uint32",
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

var filename = ""

var goFilePath = flag.String("g", "", "指定需要编译的go文件")
var protoFilePath = flag.String("p", "./proto", "指定存放proto文件的位置，默认在本路径下的proto目录下")
var goDir = flag.String("d", "", "指定存放go文件的位置，会扫描这个路径下的所有go文件")

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
	writer.WriteString(`syntax = "proto3";` + "\n")
	writer.Flush()
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
	log.Printf("生成%s.proto文件成功\n", filename)
	packageFlag = true
	startEndFlag = false
	lineCount = 0
	fieldNumber = 1
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

	strArr = field1.FindStringSubmatch(s)
	if len(strArr) == 4 {
		f := strArr[1]
		t := strArr[2]
		desc := strArr[3]
		if strings.TrimSpace(desc) != "" && !strings.HasPrefix(desc, "//") {
			desc = "//" + desc
		}
		writer.WriteString("\t" + getType(t) + " " + f + " = " + strconv.Itoa(fieldNumber) + ";  " + desc + "\n")
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
		if strings.TrimSpace(desc) != "" && !strings.HasPrefix(desc, "//") {
			desc = "//" + desc
		}
		for _, f := range fields {
			writer.WriteString("\t" + getType(t) + " " + strings.TrimSpace(f) + " = " + strconv.Itoa(fieldNumber) + ";  " + desc + "\n")
			fieldNumber++
		}
		writer.WriteString("\t" + getType(t) + " " + strings.TrimSpace(last) + " = " + strconv.Itoa(fieldNumber) + ";  " + desc + "\n")
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

	if strings.HasPrefix(s, "gorm.Model") {
		writer.WriteString("\tuint32 ID = " + strconv.Itoa(fieldNumber) + ";\n")
		fieldNumber++
		writer.WriteString("\tuint64 CreatedAt = " + strconv.Itoa(fieldNumber) + ";\n")
		fieldNumber++
		writer.WriteString("\tuint64 UpdatedAt = " + strconv.Itoa(fieldNumber) + ";\n")
		fieldNumber++
		writer.WriteString("\tuint64 DeletedAt = " + strconv.Itoa(fieldNumber) + ";\n")
		fieldNumber++
		return
	}
	if startEndFlag {
		err(s)
	}

}

func getType(s string) string {
	t, ok := typeMap[s]
	if !ok {
		log.Fatal(filename, "-", lineCount, ":无法转换成的go类型->", s)
	}
	return t
}

func err(s string) {
	s = fmt.Sprintf(filename, "-", "%d:语法错误,%s\n", lineCount, s)
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
			fmt.Errorf("文件%s-%d行:语法错误，出现多个package定义: %s", filename, lineCount, s)
			os.Exit(2)
		}
	}
}

func main() {
	flag.Parse()
	if *goDir != "" {
		files, e := ioutil.ReadDir(*goDir)
		if e != nil {
			log.Fatalf("读取%s文件夹错误->%s\n", *goDir, e)
		}

		for _, f := range files {
			name := f.Name()
			if !f.IsDir() && strings.HasSuffix(name, ".go") {
				wr(filepath.Join(*goDir, name))
				genProtoFile()
			}
		}
	}
	if *goFilePath != "" {
		if !exists(*goFilePath) {
			log.Fatal(*goFilePath, ": 目标go文件不存在")
		}
		wr(*goFilePath)
		genProtoFile()
	}

	if *goDir == "" && *goFilePath == "" {
		log.Fatalf("什么也没有发生\n")
	}
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

func wr(path string) {
	filename = path
	file, e := os.Open(path)
	if e != nil {
		fmt.Println(e)
		os.Exit(1)
	}
	reader = bufio.NewReader(file)

	f := strings.TrimSpace(*protoFilePath) + "/" + strings.Split(filepath.Base(path), ".")[0] + ".proto"
	if exists(f) {
		log.Println(f, ": 目标proto文件已存在")
	}

	if !exists(*protoFilePath) {
		e := os.Mkdir(*protoFilePath, 777)
		if e != nil {
			log.Println(e, ":创建文件夹失败")
		}
	}
	open, e := os.OpenFile(f, os.O_RDWR|os.O_CREATE, 0766)
	if e != nil {
		log.Println(e)
	}
	writer = bufio.NewWriter(open)
}

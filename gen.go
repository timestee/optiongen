package optiongen

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"unicode"

	"myitcv.io/gogenerate"
)

type FieldType int

const (
	FieldTypeFunc FieldType = iota
	FieldTypeVar
)

type fileOptionGen struct {
	FilePath    string
	NameDeclare string
	FileName    string
	PkgName     string
	ImportPath  []string

	Comments          []string
	ClassName         string
	ClassOptionFields []optionField
}

type optionField struct {
	FieldType       FieldType
	Name            string
	Type            string
	Body            string
	LastRowComments []string
	SameRowComment  string
	MethodComments  []string
}

type templateData struct {
	ClassOptionInfo     []optionInfo
	ClassComments       []string
	ClassName           string
	ClassOptionTypeName string
	ClassNewFuncName    string
	XConf               bool
}

type optionInfo struct {
	Index           int
	FieldType       FieldType
	Name            string
	NameAsParameter string
	OptionFuncName  string
	VisitFuncName   string
	GenOptionFunc   bool
	Slice           bool
	SliceElemType   template.HTML
	Type            template.HTML
	Body            template.HTML
	LastRowComments []string
	SameRowComment  string
	MethodComments  []string
	Tags            []string
	TagString       string
}

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func SnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}
func LcFirst(str string) string {
	for i, v := range str {
		return string(unicode.ToLower(v)) + str[i+1:]
	}
	return ""
}

// DefaultTrimChars are the characters which are stripped by Trim* functions in default.
var DefaultTrimChars = string([]byte{
	'\t', // Tab.
	'\v', // Vertical tab.
	'\n', // New line (line feed).
	'\r', // Carriage return.
	'\f', // New page.
	' ',  // Ordinary space.
	0x00, // NUL-byte.
	0x85, // Delete.
	0xA0, // Non-breaking space.
})

func StringTrim(str string, characterMask ...string) string {
	trimChars := DefaultTrimChars
	if len(characterMask) > 0 {
		trimChars += characterMask[0]
	}
	return strings.Trim(str, trimChars)
}

// escapeStringBackslash is similar to escapeBytesBackslash but for string.
func escapeStringBackslash(v string) string {
	buf := make([]byte, len(v)*2)
	pos := 0
	for i := 0; i < len(v); i++ {
		c := v[i]
		switch c {
		case '\x00':
			buf[pos] = '\\'
			buf[pos+1] = '0'
			pos += 2
		case '\n':
			buf[pos] = '\\'
			buf[pos+1] = 'n'
			pos += 2
		case '\r':
			buf[pos] = '\\'
			buf[pos+1] = 'r'
			pos += 2
		case '\x1a':
			buf[pos] = '\\'
			buf[pos+1] = 'Z'
			pos += 2
		case '\'':
			buf[pos] = '\\'
			buf[pos+1] = '\''
			pos += 2
		case '"':
			buf[pos] = '\\'
			buf[pos+1] = '"'
			pos += 2
		case '\\':
			buf[pos] = '\\'
			buf[pos+1] = '\\'
			pos += 2
		default:
			buf[pos] = c
			pos++
		}
	}

	return string(buf[:pos])
}
func cleanAsTag(s ...string) string {
	var tmp []string
	for _, v := range s {
		tmp = append(tmp, StringTrim(v, "//"))
	}
	return escapeStringBackslash(strings.Join(tmp, "  "))
}
func (g fileOptionGen) fatal(location string, err error, info ...string) {
	var infos []string
	infos = append(infos, "----------------------------------------- >>>>>>>>> optiongen got fatal")
	infos = append(infos, fmt.Sprintf("file: %s", g.FilePath))
	infos = append(infos, fmt.Sprintf("option: %s", g.NameDeclare))
	infos = append(infos, fmt.Sprintf("location: %s", location))
	infos = append(infos, fmt.Sprintf("error: %s", err.Error()))
	infos = append(infos, info...)
	infos = append(infos, "----------------------------------------- <<<<<<<<<")
	fmt.Println(strings.Join(infos, "\n"))
	os.Exit(1)
}
func (g fileOptionGen) gen() {
	buf := BufWrite{
		buf: bytes.NewBuffer(nil),
	}

	buf.wln(fmt.Sprintf("// Code generated by %s. DO NOT EDIT.", OptionGen))
	buf.wln(fmt.Sprintf("// %s: %s", OptionGen, "github.com/timestee/optiongen"))
	buf.wln()
	buf.wf("package %v\n", g.PkgName)

	for _, importPath := range g.ImportPath {
		buf.wf("import %v\n", importPath)
	}

	tmp := templateData{XConf: AtomicConfig().GetXConf()}

	className := g.ClassName
	indexGot := make(map[int]string)
	for _, val := range g.ClassOptionFields {
		name := strings.Trim(val.Name, "\"")
		funcName := "With"
		if AtomicConfig().GetOptionWithStructName() {
			funcName = funcName + strings.Title(className)
		}
		if strings.HasSuffix(funcName, "Options") {
			funcName = funcName[:len(funcName)-1]
		}
		if strings.HasSuffix(funcName, "Opts") {
			funcName = funcName[:len(funcName)-1]
		}

		if strings.HasPrefix(val.Type, "(") && strings.HasSuffix(val.Type, ")") {
			val.Type = val.Type[1 : len(val.Type)-1]
		}
		ss := strings.Split(name, "@")
		name = ss[0]
		genOptionFunc := !strings.HasSuffix(name, "_") && !strings.HasSuffix(name, "Inner")
		index := 0
		funcName += strings.Title(name)
		xconfTag := ""
		if len(ss) > 1 {
			for i, v := range ss {
				if i == 0 {
					// 跳过字段名
					continue
				}
				v = strings.TrimSpace(v)
				if strings.HasPrefix(v, "#") {
					numStr := strings.TrimPrefix(v, "#")
					i, err := strconv.Atoi(numStr)
					if err != nil {
						g.fatal("parse annotation #", fmt.Errorf("got error:%s when run Atoi", err.Error()))
					}
					index = i
					if got, ok := indexGot[index]; ok {
						g.fatal("parse annotation #", fmt.Errorf("got same index,%s and %s ", got, val.Name))
					}
					indexGot[index] = val.Name
					genOptionFunc = false
				}
				if strings.EqualFold(v, "inner") {
					genOptionFunc = false
				}
				if strings.EqualFold(v, "protected") {
					genOptionFunc = false
				}
				if strings.HasPrefix(v, "xconf#") {
					xconfTag = strings.TrimPrefix(v, "xconf#")
				}
			}
		}

		info := optionInfo{
			Index:           index,
			FieldType:       val.FieldType,
			Name:            name,
			NameAsParameter: LcFirst(name),
			GenOptionFunc:   genOptionFunc,
			OptionFuncName:  funcName,
			VisitFuncName:   "Get" + name,
			Slice:           strings.HasPrefix(val.Type, "[]"),
			SliceElemType:   template.HTML(strings.Replace(val.Type, "[]", "", 1)),
			Type:            template.HTML(val.Type),
			Body:            template.HTML(val.Body),
			LastRowComments: val.LastRowComments,
			SameRowComment:  val.SameRowComment,
			MethodComments:  val.MethodComments,
		}
		if AtomicConfig().GetXConf() {
			if xconfTag == "" {
				xconfTag = SnakeCase(info.Name)
			}
			info.Tags = append(info.Tags, fmt.Sprintf(`xconf:"%s"`, xconfTag))
		}
		if AtomicConfig().GetUsageTagName() != "" {
			s := cleanAsTag(val.MethodComments...)
			if s != "" {
				info.Tags = append(info.Tags, fmt.Sprintf(`%s:"%s"`, AtomicConfig().GetUsageTagName(), s))
			}
		}
		if len(info.Tags) > 0 {
			tag := strings.Join(info.Tags, " ")
			info.TagString = fmt.Sprintf("`%s`", tag)
			if err := validateStructTag(tag); err != nil {
				g.fatal("tag", err, "tag: "+tag)
			}
		}

		// []byte不作为数组类型处理
		if strings.TrimSpace(strings.TrimLeft(val.Type, "[]")) == "byte" {
			info.Slice = false
		}
		tmp.ClassOptionInfo = append(tmp.ClassOptionInfo, info)
	}
	optionTypeName := className + "Option"
	if strings.HasSuffix(className, "Options") {
		optionTypeName = className[:len(className)-1]
	}
	if strings.HasSuffix(className, "Opts") {
		optionTypeName = className[:len(className)-1]
	}
	tmp.ClassOptionTypeName = optionTypeName
	tmp.ClassComments = g.Comments
	tmp.ClassName = g.ClassName
	newFuncName := AtomicConfig().GetNewFunc()
	if newFuncName == "" {
		newFuncName = fmt.Sprintf("New%s", className)
	}
	sort.Slice(tmp.ClassOptionInfo, func(i, j int) bool {
		return tmp.ClassOptionInfo[i].Index < tmp.ClassOptionInfo[j].Index
	})
	var pameters []string
	for _, v := range tmp.ClassOptionInfo {
		if v.Index == 0 {
			continue
		}
		pameters = append(pameters, fmt.Sprintf("%s %s", v.NameAsParameter, v.Type))
	}
	if len(pameters) == 0 {
		tmp.ClassNewFuncName = fmt.Sprintf("%s(opts... %s)", newFuncName, optionTypeName)
	} else {
		tmp.ClassNewFuncName = fmt.Sprintf("%s(%s, opts... %s)", newFuncName, strings.Join(pameters, ","), optionTypeName)
	}

	funcMap := template.FuncMap{
		"unescaped": unescaped,
	}

	t := template.Must(template.New("tmp").Funcs(funcMap).Parse(templateTextWithPreviousSupport))

	err := t.Execute(buf.buf, tmp)
	if err != nil {
		g.fatal("tempalt_execute", err)
	}

	if strings.HasPrefix(g.FileName, "gen_") {
		g.FileName = strings.TrimLeft(g.FileName, "gen_")
	}

	genName := gogenerate.NameFile(g.FileName, OptionGen)
	source := g.goimportsBuf(buf.buf)
	if err := ioutil.WriteFile(genName, source.Bytes(), 0644); err != nil {
		g.fatal("write_file", err)
	}
	if AtomicConfig().GetDebug() {
		log.Println(fmt.Sprintf("%s/%s", g.PkgName, genName))
	}
}

func (g fileOptionGen) goimportsBuf(buf *bytes.Buffer) *bytes.Buffer {
	out := bytes.NewBuffer(nil)
	cmd := exec.Command("goimports")
	data := buf.Bytes()
	cmd.Stdin = bytes.NewReader(data)
	cmd.Stdout = out
	err := cmd.Run()
	if err != nil {
		g.fatal("goimports", err, "==========> INVALID SOURCE CODE <==========", string(data))
	}
	return out
}
func unescaped(str string) template.HTML { return template.HTML(str) }

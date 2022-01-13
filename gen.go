package optiongen

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os/exec"
	"sort"
	"strings"

	"github.com/timestee/optiongen/annotation"
	"github.com/timestee/optiongen/xutil"
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

	Comments  []string
	ClassName string

	ClassOptionFields []optionField
	Annotations       []annotation.Annotation
}

func (g *fileOptionGen) ParseAnnotations() (err error) {
	var allComments []string
	for _, v := range g.ClassOptionFields {
		allComments = append(allComments, v.LastRowComments...)
		allComments = append(allComments, v.SameRowComment)
		allComments = append(allComments, v.MethodComments...)
	}
	allComments = append(allComments, g.Comments...)
	g.Annotations, err = annotation.NewRegistry().ResolveAnnotationsErrorDuplicate(allComments)
	if AtomicConfig().GetDebug() {
		fmt.Printf("\n===>>> ParseAnnotations all comments ===>>> \n %s \n", strings.Join(allComments, "\n"))
		fmt.Printf("\n===>>> ParseAnnotations annotations  ===>>> \n %v \n", g.Annotations)
	}
	return
}
func (g *fileOptionGen) GetAnnotation(name string) annotation.Annotation {
	for _, v := range g.Annotations {
		if strings.EqualFold(v.Name, name) {
			return v
		}
	}
	return annotation.Annotation{}
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
	ClassOptionInfo       []optionInfo
	ClassComments         []string
	ClassName             string
	ClassNameTitle        string
	ClassOptionTypeName   string
	ClassNewFuncSignature string
	ClassNewFuncName      string
	XConf                 bool
	OptionReturnPrevious  bool

	VisitorName   string
	InterfaceName string
}

type optionInfo struct {
	ArgIndex            int
	FieldType           FieldType
	Name                string
	NameAsParameter     string
	OptionFuncName      string
	VisitFuncName       string
	VisitFuncReturnType template.HTML
	GenOptionFunc       bool
	Slice               bool
	SliceElemType       template.HTML
	Type                template.HTML
	Body                template.HTML
	LastRowComments     []string
	SameRowComment      string

	OptionComment    string
	VisitFuncComment string
	Tags             []string
	TagString        string
}

func cleanAsTag(s ...string) string {
	var tmp []string
	for _, v := range s {
		tmp = append(tmp, xutil.StringTrim(v, "//"))
	}
	return xutil.EscapeStringBackslash(strings.Join(tmp, " , "))
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
	panic(strings.Join(infos, "\n"))
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

	tmp := templateData{XConf: AtomicConfig().GetXConf(), OptionReturnPrevious: AtomicConfig().GetOptionReturnPrevious()}

	className := g.ClassName
	indexGot := make(map[int]string)
	for _, val := range g.ClassOptionFields {
		name := strings.Trim(val.Name, "\"")
		funcName := "With"
		if AtomicConfig().GetOptionPrefix() != "" {
			funcName = AtomicConfig().GetOptionPrefix()
		} else {
			if AtomicConfig().GetOptionWithStructName() {
				funcName = funcName + strings.Title(className)
			}
			if strings.HasSuffix(funcName, "Options") {
				funcName = funcName[:len(funcName)-1]
			}
			if strings.HasSuffix(funcName, "Opts") {
				funcName = funcName[:len(funcName)-1]
			}

		}
		if strings.HasPrefix(val.Type, "(") && strings.HasSuffix(val.Type, ")") {
			val.Type = val.Type[1 : len(val.Type)-1]
		}

		name = strings.Split(name, "@")[0]
		nameSnakeCase := xutil.SnakeCase(name)

		funcName += strings.Title(name)

		an := g.GetAnnotation(name)
		private := an.GetBool(AnnotationKeyPrivate, strings.HasSuffix(name, "_") || strings.HasSuffix(name, "Inner"))
		if AtomicConfig().GetDebug() {
			fmt.Printf("===>>> Field Annotation name: %s attributes: %v\n", name, an.Attributes)
		}
		xconfTag := an.GetString(AnnotationKeyXConfTag, nameSnakeCase)
		argIndex := an.GetInt(AnnotationKeyArg)
		getterType := an.GetString(AnnotationKeyGetter, val.Type)
		optionFuncName := an.GetString(AnnotationKeyOption, funcName)
		comment := an.GetString(AnnotationKeyComment)
		deprecated := an.GetString(AnnotationKeyDeprecated)

		if argIndex != 0 {
			if got, ok := indexGot[argIndex]; ok {
				g.fatal("parse annotation "+AnnotationKeyArg, fmt.Errorf("got same index,%s and %s ", got, val.Name))
			}
			indexGot[argIndex] = name
		}
		if AtomicConfig().GetXConf() {
			name = strings.Title(name)
		}
		info := optionInfo{
			ArgIndex:            argIndex,
			FieldType:           val.FieldType,
			Name:                name,
			NameAsParameter:     xutil.LcFirst(name),
			GenOptionFunc:       !private && argIndex == 0,
			OptionFuncName:      optionFuncName,
			VisitFuncName:       "Get" + strings.Title(name),
			VisitFuncReturnType: template.HTML(getterType),
			Slice:               strings.HasPrefix(val.Type, "[]"),
			SliceElemType:       template.HTML(strings.Replace(val.Type, "[]", "", 1)),
			Type:                template.HTML(val.Type),
			Body:                template.HTML(val.Body),
			LastRowComments:     val.LastRowComments,
			SameRowComment:      val.SameRowComment,
		}
		methodComments := val.MethodComments
		if comment != "" {
			if len(methodComments) == 0 {
				methodComments = append(methodComments, comment)
			}
		}
		for index, v := range methodComments {
			methodComments[index] = xutil.StringTrim(v, "//", ",", ".")
		}
		if len(methodComments) == 0 {
			info.OptionComment = fmt.Sprintf("%s option func for filed %s", info.OptionFuncName, info.Name)
		} else {
			info.OptionComment = xutil.WrapString(fmt.Sprintf("%s %s", info.OptionFuncName, xutil.StringTrim(strings.Join(methodComments, ","))), 200)
		}
		info.OptionComment = xutil.CleanAsComment(info.OptionComment)
		infoForUsage := methodComments
		if deprecated != "" {
			info.OptionComment += "\n//"
			info.OptionComment += "\n// Deprecated: " + deprecated

			info.VisitFuncComment += fmt.Sprintf("\n// %s visitor func for filed %s", info.VisitFuncName, info.Name)
			info.VisitFuncComment += "\n//"
			info.VisitFuncComment += "\n// Deprecated: " + deprecated

			infoForUsage = append(infoForUsage, "Deprecated: "+deprecated)
		}
		if AtomicConfig().GetXConf() {
			if xconfTag == "" {
				xconfTag = xutil.SnakeCase(info.Name)
			}
			if AtomicConfig().GetXConfTrimPrefix() != "" {
				xconfTag = strings.TrimPrefix(xconfTag, AtomicConfig().GetXConfTrimPrefix())
			}
			if deprecated != "" && !strings.Contains(xconfTag, "deprecated") {
				xconfTag += ",deprecated"
			}
			info.Tags = append(info.Tags, fmt.Sprintf(`xconf:"%s"`, xconfTag))
		}
		if AtomicConfig().GetUsageTagName() != "" {
			s := cleanAsTag(infoForUsage...)
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
	optionTypeName = strings.Title(optionTypeName)
	tmp.ClassOptionTypeName = optionTypeName
	tmp.ClassComments = g.Comments
	tmp.ClassName = g.ClassName
	tmp.ClassNameTitle = strings.Title(tmp.ClassName)

	tmp.VisitorName = fmt.Sprintf("%sVisitor", tmp.ClassNameTitle)
	tmp.InterfaceName = fmt.Sprintf("%sInterface", tmp.ClassNameTitle)

	newFuncReturn := "*" + className

	if AtomicConfig().GetNewFuncReturn() == NewFuncReturnVisitor {
		newFuncReturn = tmp.VisitorName
	} else if AtomicConfig().GetNewFuncReturn() == NewFuncReturnInterface {
		newFuncReturn = tmp.InterfaceName
	}

	newFuncName := AtomicConfig().GetNewFunc()
	if newFuncName == "" {
		newFuncName = fmt.Sprintf("New%s", strings.Title(className))
	}
	sort.Slice(tmp.ClassOptionInfo, func(i, j int) bool {
		return tmp.ClassOptionInfo[i].ArgIndex < tmp.ClassOptionInfo[j].ArgIndex
	})
	var pameters []string
	for _, v := range tmp.ClassOptionInfo {
		if v.ArgIndex == 0 {
			continue
		}
		pameters = append(pameters, fmt.Sprintf("%s %s", v.NameAsParameter, v.Type))
	}
	if len(pameters) == 0 {
		tmp.ClassNewFuncSignature = fmt.Sprintf("func %s(opts... %s) %s", newFuncName, optionTypeName, newFuncReturn)
	} else {
		tmp.ClassNewFuncSignature = fmt.Sprintf("func %s(%s, opts... %s) %s", newFuncName, strings.Join(pameters, ","), optionTypeName, newFuncReturn)
	}

	tmp.ClassNewFuncName = newFuncName

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

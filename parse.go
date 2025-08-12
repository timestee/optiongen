package optiongen

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"runtime"

	"myitcv.io/gogenerate"
)

func paniwWith(format string, v ...interface{}) {
	info := fmt.Sprintf(format, v...)
	panic(info)
}

var fset = token.NewFileSet()

func inspectDir(wd string) (envFile string, lineNo int) {
	envFile, ok := os.LookupEnv(gogenerate.GOFILE)
	if !ok {
		paniwWith("env not correct; missing %v", gogenerate.GOFILE)
	}

	lineStr, ok := os.LookupEnv(gogenerate.GOLINE)
	if !ok {
		paniwWith("env not correct; missing %v", gogenerate.GOLINE)
	}
	lineNo, err := strconv.Atoi(lineStr)
	if err != nil {
		paniwWith("env not correct; env %v convert to int failed. %s", gogenerate.GOLINE, err)
	}

	tags := make(map[string]bool)

	goos := os.Getenv("GOOS")
	if goos == "" {
		goos = runtime.GOOS
	}
	tags[goos] = true

	goarch := os.Getenv("GOARCH")
	if goarch == "" {
		goarch = runtime.GOARCH
	}
	tags[goarch] = true

	dirFiles, err := gogenerate.FilesContainingCmd(wd, OptionGen, tags)
	if err != nil {
		paniwWith("could not determine if we are the first file: %v", err)
	}

	if dirFiles == nil {
		paniwWith("cannot find any files containing the %v directive", OptionGen)
	}
	return
}

var gofmtBuf bytes.Buffer

// gofmt returns the gofmt-formatted string for an AST node.
func gofmt(n interface{}) string {
	gofmtBuf.Reset()
	err := printer.Fprint(&gofmtBuf, fset, n)
	if err != nil {
		return "<" + err.Error() + ">"
	}
	return gofmtBuf.String()
}

func printLog(format string, a ...interface{}) {
	if AtomicConfig().GetDebug() {
		fmt.Println(fmt.Sprintf(format, a...))
	}
}

var commentReg = regexp.MustCompile(`@MethodComment\(([^)]+)\)`)

func parseComment(comment string) (string, []string) {
	if len(comment) == 0 {
		return "", nil
	}
	out := commentReg.FindAllStringSubmatch(comment, -1)
	if len(out) == 0 {
		return comment, nil
	}
	var mc []string
	for _, v := range out {
		comment = strings.ReplaceAll(comment, v[0], "")
		if len(v[1]) > 0 {
			mc = append(mc, fmt.Sprintf("// %s", v[1]))
		}
	}
	if strings.TrimSpace(comment) == "//" {
		return "", mc
	}
	return comment, mc
}

func ParseDir(dir string) {
	fileName, lineNo := inspectDir(dir)

	DstName := ""
	declareName := ""

	file, err := parser.ParseFile(fset, fileName, nil, parser.ParseComments)
	if err != nil {
		paniwWith("unable to parse %v: %v", fileName, err)
	}
	var importPath []string
	for _, imp := range file.Imports {
		importPath = append(importPath, imp.Path.Value)
	}

	comments := file.Comments
	var classComments []string
	var classOptionFields []optionField
	var typeParams []TypeParam

	var lastMaxBodyPos token.Pos
	for _, d := range file.Decls {
		switch d := d.(type) {
		default:
			if lastMaxBodyPos < d.End() {
				lastMaxBodyPos = d.End()
			}
		case *ast.FuncDecl:
			if d.Recv != nil {
				continue
			}

			// name check valid
			if !strings.HasSuffix(d.Name.Name, optionDeclarationSuffix) {
				if lastMaxBodyPos < d.End() {
					lastMaxBodyPos = d.End()
				}
				continue
			}
			// command line check valid
			p := fset.Position(d.Pos())
			if p.Line != lineNo+1 {
				if lastMaxBodyPos < d.End() {
					lastMaxBodyPos = d.End()
				}
				continue
			}
			// Only allow return expr in class option declaration function
			if len(d.Body.List) != 1 {
				continue
			}
			// Only allow return one value
			stmt, ok := d.Body.List[0].(*ast.ReturnStmt)
			if !ok {
				continue
			}
			if len(stmt.Results) != 1 {
				continue
			}
			var pos token.Pos
			if d.Doc != nil {
				pos = d.Doc.Pos()
			}
			declareName = d.Name.Name
			declarationClassName := strings.TrimPrefix(strings.TrimSuffix(d.Name.Name, optionDeclarationSuffix), "_")
			for k, v := range comments {
				if v.Pos() == pos {
					for _, v1 := range comments[:k] {
						if v1.Pos() < lastMaxBodyPos {
							continue
						}
						for _, v2 := range v1.List {
							classComments = append(classComments, v2.Text)
						}
					}
					comments = comments[k+1:]
					break
				}
			}

			if d.Type != nil && d.Type.TypeParams != nil {
				typeParams = make([]TypeParam, len(d.Type.TypeParams.List))
				for i, f := range d.Type.TypeParams.List {
					if len(f.Names) == 0 {
						paniwWith("optionGen %s type param %d has no name", declareName, i)
					}
					typeParams[i] = TypeParam{
						Name: gofmt(f.Names[0]),
						Type: gofmt(f.Type),
					}
				}
			}

			result := stmt.Results[0].(*ast.CompositeLit)
			optionFields := make([]optionField, len(result.Elts))
			for i, elt := range result.Elts {
				switch elt := elt.(type) {
				case *ast.KeyValueExpr:
					// Option Field Name
					key := elt.Key.(*ast.BasicLit)
					optionFields[i].Name = key.Value
					for _, v := range comments {
						if v.Pos() <= elt.Pos() {
							comments = comments[1:]
							for _, vv := range v.List {
								fc, mc := parseComment(vv.Text)
								if len(fc) > 0 {
									optionFields[i].LastRowComments = append(optionFields[i].LastRowComments, fc)
								}
								if len(mc) > 0 {
									optionFields[i].MethodComments = append(optionFields[i].MethodComments, mc...)
								}
							}
							continue
						}
						eltP := fset.Position(elt.Pos())
						vP := fset.Position(v.Pos())
						if eltP.Line == vP.Line {
							comments = comments[1:]
							for _, vv := range v.List {
								fc, mc := parseComment(vv.Text)
								if len(fc) > 0 {
									optionFields[i].SameRowComment = fc
								}
								if len(mc) > 0 {
									optionFields[i].MethodComments = append(optionFields[i].MethodComments, mc...)
								}
								break
							}
							continue
						}
						break
					}
					switch value := elt.Value.(type) {
					case *ast.FuncLit:
						printLog("%s type: %s", optionFields[i].Name, "ast.FuncLit")
						optionFields[i].FieldType = FieldTypeFunc
						buf := bytes.NewBufferString("")
						// Option func Type
						_ = printer.Fprint(buf, fset, value.Type)
						optionFields[i].Type = buf.String()

						// Option func Body
						buf.Reset()
						_ = printer.Fprint(buf, fset, value.Body)
						optionFields[i].Body = buf.String()
					case *ast.CallExpr:
						printLog("%s type: %s", optionFields[i].Name, "ast.CallExpr")
						optionFields[i].FieldType = FieldTypeVar
						buf := bytes.NewBufferString("")

						// Option Variable Type
						_ = printer.Fprint(buf, fset, value.Fun)
						optionFields[i].Type = buf.String()
						// Option Variable Value
						buf.Reset()
						_ = printer.Fprint(buf, fset, value.Args[0])
						optionFields[i].Body = buf.String()
					case *ast.BasicLit: // token.INT, token.FLOAT, token.IMAG, token.CHAR, or token.STRING
						printLog("%s type: %s", optionFields[i].Name, "ast.BasicLit")
						optionFields[i].FieldType = FieldTypeVar
						switch value.Kind {
						case token.INT:
							optionFields[i].Type = "int"
							optionFields[i].Body = value.Value
						case token.FLOAT:
							optionFields[i].Type = "float"
							optionFields[i].Body = value.Value
						case token.CHAR:
							optionFields[i].Type = "byte"
							optionFields[i].Body = value.Value
						case token.STRING:
							optionFields[i].Type = "string"
							optionFields[i].Body = value.Value
						}
					case *ast.CompositeLit:
						printLog("%s type: %s", optionFields[i].Name, "ast.CompositeLit")
						optionFields[i].FieldType = FieldTypeVar
						optionFields[i].Type = gofmt(value.Type)
						buf := bytes.NewBufferString("")
						buf.Reset()
						var data []string
						for _, p := range value.Elts {
							switch t := p.(type) {
							case *ast.BasicLit:
								data = append(data, t.Value)
							case *ast.Ident:
								if t.Name == "false" || t.Name == "true" {
									data = append(data, t.Name)
								}
							case *ast.KeyValueExpr:
								blKey, okKey := t.Key.(*ast.BasicLit)
								blVal, okVal := t.Value.(*ast.BasicLit)
								if !okKey || !okVal {
									paniwWith("optionGen %s got type %s support basic types only", optionFields[i].Name, optionFields[i].Type)
								}
								data = append(data, fmt.Sprintf("%s:%s", blKey.Value, blVal.Value))
							default:
								paniwWith("optionGen %s got type %s support basic types only", optionFields[i].Name, optionFields[i].Type)
							}
						}
						val := "make(" + optionFields[i].Type + ",0)"
						if AtomicConfig().GetEmptyCompositeNil() {
							val = "nil"
						}
						if len(data) > 0 {
							val = fmt.Sprintf("%s{%s}", optionFields[i].Type, strings.Join(data, ","))
						}
						optionFields[i].Body = val
					case *ast.Ident:
						if value.Obj == nil {
							printLog("%s type: %s", optionFields[i].Name, "ast.Ident")
						} else {
							printLog("%s type: %s Object:%v Type:%v", optionFields[i].Name, "ast.Ident", value.Obj, value.Obj.Type)
						}
						optionFields[i].FieldType = FieldTypeVar
						if value.Name == "false" || value.Name == "true" {
							optionFields[i].Type = "bool"
							optionFields[i].Body = value.Name
							break
						}
						if value.Name == "nil" {
							optionFields[i].Type = "interface{}"
							optionFields[i].Body = value.Name
							break
						}
						paniwWith("optionGen %s got type ast.Ident, please give some type hint like: string(DefaultValue) ", optionFields[i].Name)
					default:
						paniwWith("optionGen %s got type %s support basic types only", optionFields[i].Name, optionFields[i].Type)
					}
				}
			}

			classOptionFields = optionFields
			DstName = declarationClassName
		}
		// find dst file line specify name
		if DstName != "" {
			break
		}
	}
	if DstName == "" {
		paniwWith("specify file %s line %d cannot find generate declare", fileName, lineNo+1)
	}
	info := fmt.Sprintf("%s:%d [%s]", filepath.Join(dir, fileName), lineNo, declareName)
	fmt.Printf("🚀  %s running => %s ...\n", OptionGen, info)
	defer func() {
		if reason := recover(); reason != nil {
			err = fmt.Errorf("%v", reason)
		}
		if err != nil {
			fmt.Printf("💀  got error when process:%s , err:%s\n", info, err.Error())
		}
	}()
	pkgName := file.Name.Name
	g := fileOptionGen{
		FilePath:          filepath.Join(dir, fileName),
		NameDeclare:       declareName,
		FileName:          strings.ToLower(DstName),
		PkgName:           pkgName,
		ImportPath:        importPath,
		ClassName:         DstName,
		ClassOptionFields: classOptionFields,
		Comments:          classComments,
		TypeParams:        typeParams,
	}
	err = g.ParseAnnotations()
	if err != nil {
		panic(err)
	}
	g.gen()
}

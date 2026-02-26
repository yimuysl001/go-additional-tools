package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"
	"unicode"
)

func main() {
	pkgName := "gf/gutil"
	sourceDir := "D:\\GO\\gopath\\pkg\\mod\\github.com\\gogf\\gf\\v2@v2.10.0\\util\\gutil"
	pname := "gutil"
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, sourceDir, nil, parser.ParseComments)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Parse error: %v\n", err)
		os.Exit(1)
	}

	for _, pkg := range pkgs {
		if strings.HasSuffix(pkg.Name, "_test") {
			continue
		}
		s := generateRegisterCode(pkgName, pkg, pname, fset)

		fmt.Println(s)
	}
}

func generateRegisterCode(importPkg string, pkg *ast.Package, pname string, fset *token.FileSet) string {
	registerMap := make(map[string][]string)
	comments := make(map[string]string)

	for _, file := range pkg.Files {
		// 收集注释
		for _, decl := range file.Decls {
			switch d := decl.(type) {
			case *ast.FuncDecl:
				name := d.Name.Name
				if d.Doc != nil {
					comments[name] = d.Doc.Text()
				}
				registerMap["funcs"] = append(registerMap["funcs"], name)

			case *ast.GenDecl:
				if d.Tok == token.TYPE {
					for _, spec := range d.Specs {
						if ts, ok := spec.(*ast.TypeSpec); ok {
							name := ts.Name.Name
							if d.Doc != nil {
								comments[name] = d.Doc.Text()
							}

							switch ts.Type.(type) {
							case *ast.StructType:
								registerMap["structs"] = append(registerMap["structs"], name)
							case *ast.InterfaceType:
								registerMap["interfaces"] = append(registerMap["interfaces"], name)
							default:
								registerMap["types"] = append(registerMap["types"], name)
							}
						}
					}
				}
			}
		}
	}
	var sb bytes.Buffer
	// 生成注册代码
	fmt.Fprintf(&sb, "\n// Register imports for package %s\n", importPkg)
	fmt.Fprintf(&sb, "RegisterImport(%q, map[string]any{\n", importPkg)

	// 按类别组织
	categories := []string{"funcs", "structs"}
	for _, cat := range categories {
		items := registerMap[cat]
		if len(items) > 0 {
			fmt.Printf("\t// %s\n", strings.Title(cat))
			for _, item := range items {
				if !IsFirstLetterUpper(item) {
					continue
				}

				if cat == "funcs" && (strings.HasPrefix(item, "Benchmark_") || strings.HasPrefix(item, "Test_")) {
					continue
				}

				if comment, ok := comments[item]; ok {
					// 格式化注释
					comment = strings.TrimSpace(comment)
					comment = strings.ReplaceAll(comment, "\n", "\n\t// ")
					fmt.Fprintf(&sb, "\t// %s\n", comment)
				}
				if cat == "structs" {
					fmt.Fprintf(&sb, "\t%q: %s.%s,\n", item, pname, item+"{}")

				} else {
					fmt.Fprintf(&sb, "\t%q: %s.%s,\n", item, pname, item)
				}

			}
		}
	}

	fmt.Fprintf(&sb, "})\n")
	return sb.String()
}
func IsFirstLetterUpper(s string) bool {
	if len(s) == 0 {
		return false
	}
	r := []rune(s)[0] // 处理 Unicode 字符
	return unicode.IsUpper(r)
}

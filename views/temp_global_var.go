package views

import (
	"embed"
	"html/template"
)

// main 函数中定义的全局变量 在其他包中无法调用 ，重新定义新包实现全局变量
// go:embed 不支持相对路径，只能获取当前目录下的目录或文件
var (
	//go:embed *.html
	embedTmpl embed.FS

	// 自定义的函数必须在调用ParseFiles() ParseFS()之前创建。
	funcMap = template.FuncMap{}
	GoTpl = template.Must(
		template.New("").
			Funcs(funcMap).
			ParseFS(embedTmpl, "*.html"))
)

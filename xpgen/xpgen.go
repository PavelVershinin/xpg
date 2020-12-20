package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"regexp"
	"strings"
	"text/template"

	"github.com/PavelVershinin/xpg/underscore"
)

var (
	regexpTag = regexp.MustCompile(`(?i)xpg:"([^"]+)"`)
	path      *string
	connect   *string
)

func main() {
	path = flag.String("path", "", "-path=/absolute/path/to/the/model.go")
	connect = flag.String("connect", "main", "-connect=main")
	flag.Parse()

	if len(*path) == 0 {
		fmt.Println("The path flag was not passed")
		os.Exit(0)
	}

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, *path, nil, 0)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, decl := range f.Decls {
		switch decl := decl.(type) {
		case *ast.GenDecl:
			if decl.Tok == token.TYPE {
				for _, spec := range decl.Specs {
					if tspec := spec.(*ast.TypeSpec); tspec.Type != nil && tspec.Name.IsExported() {
						if ident, ok := tspec.Type.(*ast.StructType); ok {
							if code, err := codeGen(tspec.Name.String(), ident.Fields); err != nil {
								fmt.Println(err)
							} else {
								fmt.Println(code)
							}
						}
					}
				}
			}
		}
	}
}

func codeGen(modelName string, fields *ast.FieldList) (string, error) {
	var result = &bytes.Buffer{}
	var validModel bool
	var dbName = underscore.Underscore(modelName)

	for _, f := range fields.List {
		if fmt.Sprintf("%s", f.Type) == "&{xpg Model}" {
			validModel = true
			break
		}
	}

	if !validModel {
		return "", fmt.Errorf("the %s structure is not the correct xpg.Model", modelName)
	}

	var selectColumns, modelColumns, dbColumns []string

	for _, f := range fields.List {
		if len(f.Names) > 0 {
			columnModelName := f.Names[0].Name
			columnDbName := dbColumnName(columnModelName, f.Tag)
			selectColumns = append(selectColumns, `"`+dbName+`"."`+columnDbName+`"`)
			modelColumns = append(modelColumns, columnModelName)
			dbColumns = append(dbColumns, columnDbName)
		}
	}

	t := template.Must(template.New(modelName).Parse(tpl))
	err := t.Execute(result, map[string]interface{}{
		"model_name":        modelName,
		"model_letter":      string([]rune(strings.ToLower(modelName))[0]),
		"model_columns":     modelColumns,
		"db_name":           dbName,
		"db_columns":        dbColumns,
		"db_select_columns": "`\n\t\t" + strings.Join(selectColumns, ",\n\t\t") + "\n\t`",
		"connection_name":   *connect,
	})
	return result.String(), err
}

func dbColumnName(name string, tag *ast.BasicLit) string {
	if tag != nil {
		value := regexpTag.FindStringSubmatch(tag.Value)
		if len(value) < 2 {
			return strings.Fields(value[1])[0]
		}
	}
	return underscore.Underscore(name)
}

var tpl = `
{{$modelColumns := .model_columns}}
{{$modelLetter := .model_letter}}
// Table Возвращает название таблицы в базе данных
func ({{.model_name}}) Table() string {
	return "{{.db_name}}"
}

// Columns Список полей, которые необходимо получать запросом SELECT
func ({{.model_name}}) Columns() string {
	return {{.db_select_columns}}
}

// Connection Возвращает название подключения к БД
func ({{.model_name}}) Connection() (name string) {
	return "{{.connection_name}}"
}

// ScanRow Реализация чтения строки из результата запроса
func ({{.model_name}}) ScanRow(rows pgx.Rows) (xpg.Tabler, error) {
	row := &{{.model_name}}{}
	err := rows.Scan(
		&row.ID,
		{{range .model_columns}}&row.{{.}},
		{{end}}&row.CreatedAt,
		&row.UpdatedAt,
	)

	return row, err
}

// Save Сохранение новой/измененной структуры в БД
func ({{$modelLetter}} *{{.model_name}}) Save() (err error) {
	data := map[string]interface{}{
		"id": {{$modelLetter}}.ID, 
		{{range $i, $k := .db_columns}}"{{$k}}": {{$modelLetter}}.{{index $modelColumns $i}},
		{{end}}
	}
	{{$modelLetter}}.ID, err = xpg.New({{$modelLetter}}).Write(data)
	return err
}

// Delete Удаление записи из БД
func ({{$modelLetter}} *{{.model_name}}) Delete() error {
	return xpg.New({{$modelLetter}}).Where("id", "=", {{$modelLetter}}.ID).Delete()
}

// DbTake Получение записи из БД
func ({{$modelLetter}} *{{.model_name}}) DbTake(force ...bool) error {
	if {{$modelLetter}}.ID > 0 && (!{{$modelLetter}}.Valid || (len(force) > 0 && force[0])) {
		row, err := xpg.New(&{{.model_name}}{}).Where("id", "=", {{$modelLetter}}.ID).First()
		if err != nil {
			return err
		}
		*{{$modelLetter}} = *row.(*{{.model_name}})
		{{$modelLetter}}.Valid = true
	}
	return nil
}
`

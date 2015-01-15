package gen

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"text/template"
	"unicode"

	"github.com/gedex/inflector"
)

func writeToFile(file string, structs structs) error {
	f, err := os.Create(file)
	if err != nil {
		return err
	}

	w := bufio.NewWriter(f)
	defer w.Flush()

	return write(w, structs)
}

func write(w io.Writer, structs structs) error {

	const tplText = `package {{.Package}}

import (
	"database/sql"

	"github.com/monochromegane/goar"
)

var db *sql.DB

func Use(DB *sql.DB) {
	db = DB
}
{{range .}}
var {{.Name}}Scopes = map[string]{{.Name}}Scope{}

type {{.Name}}Scope func(r *{{.Name}}Relation, args ...interface{}) *{{.Name}}Relation

func (m {{.Name}}) Scope(name string, scope {{.Name}}Scope) {
        {{.Name}}Scopes[name] = scope
}

func (m {{.Name}}) Find(id {{.PrimaryKeyType}}) (*{{.Name}}, error) {
	r := m.newRelation()
        q, b := r.Select.Where("{{.PrimaryKeyColumn}}", id).Build()
        row := &{{.Name}}{}
        if err := db.QueryRow(q, b...).Scan({{.FieldNames "&row."| joinField}}); err != nil {
                return nil, err
        }
        return row, nil
}

func (m {{.Name}}) Where(cond string, args ...interface{}) *{{.Name}}Relation {
	r := m.newRelation()
        return r.Where(cond, args...)
}

func (m *{{.Name}}) IsNewRecord() bool {
	return ar.IsZero(m.{{.PrimaryKeyField}})
}

func (m *{{.Name}}) IsPersistent() bool {
	return !m.IsNewRecord()
}

type {{.Name}}Params {{.Name}}

func (m {{.Name}}) Create(p {{.Name}}Params) (*{{.Name}}, error) {
	n := &{{.Name}}{
	{{range .FieldNames ""}}{{.}}: p.{{.}},
	{{end}}
	}
	err := n.Save()
	return n, err
}

func (m *{{.Name}}) Save() error {
	if m.IsNewRecord() {
		ins := &ar.Insert{}
		q, b := ins.Table("{{.Name}}").Params(ar.Params{ {{$f := .FieldNames "m."}}
		{{range $index, $column := .ColumnNames}}"{{$column}}": {{index $f $index}},
		{{end}}
		}).Build()

		if _, err := db.Exec(q, b...); err != nil {
			return err
		}
		return nil
	}else{
		upd := &ar.Update{}
		q, b := upd.Table("{{.Name}}").Params(ar.Params{ {{$f := .FieldNames "m."}}
		{{range $index, $column := .ColumnNames}}"{{$column}}": {{index $f $index}},
		{{end}}
		}).Where("{{.PrimaryKeyColumn}}", m.{{.PrimaryKeyField}}).Build()

		if _, err := db.Exec(q, b...); err != nil {
			return err
		}
		return nil
	}
	return nil
}

func (m *{{.Name}}) newRelation() *{{.Name}}Relation {
	sel := &ar.Select{}
	sel.Table("{{.Name}}").Columns({{.ColumnNames | joinColumn}})
	return &{{.Name}}Relation{sel}
}

type {{.Name}}Relation struct {
	*ar.Select
}

func (r *{{.Name}}Relation) Query() ([]*{{.Name}}, error) {
        q, b := r.Build()
        rows, err := db.Query(q, b...)
        if err != nil {
                return nil, err
        }
        defer rows.Close()

	results := []*{{.Name}}{}
        for rows.Next() {
                row := &{{.Name}}{}
		if err := rows.Scan({{.FieldNames "&row."| joinField}}); err != nil {
                        return nil, err
                }
                results = append(results, row)
        }
        return results, nil
}

func (r *{{.Name}}Relation) First() (*{{.Name}}, error) {
	q, b := r.OrderBy("{{.PrimaryKeyColumn}}", ar.ASC).Limit(1).Build()
        row := &{{.Name}}{}
        if err := db.QueryRow(q, b...).Scan({{.FieldNames "&row."| joinField}}); err != nil {
                return nil, err
        }
        return row, nil
}

func (r *{{.Name}}Relation) Last() (*{{.Name}}, error) {
	q, b := r.OrderBy("{{.PrimaryKeyColumn}}", ar.DESC).Limit(1).Build()
        row := &{{.Name}}{}
        if err := db.QueryRow(q, b...).Scan({{.FieldNames "&row."| joinField}}); err != nil {
                return nil, err
        }
        return row, nil
}

func (r *{{.Name}}Relation) Where(cond string, args ...interface{}) *{{.Name}}Relation {
        r.Select.Where(cond, args...)
        return r
}

func (r *{{.Name}}Relation) And(cond string, args ...interface{}) *{{.Name}}Relation {
        r.Select.And(cond, args...)
        return r
}
{{$model := .}}
{{range .Anotations}}
{{if .BelongsTo}}
func (m *{{$model.Name}}) {{.Arg | capitalize}}() (*{{.Arg | capitalize}}, error) {
	return {{.Arg | capitalize}}{}.Where("{{$model.PrimaryKeyColumn}}", m.{{.Arg | capitalize}}ID).First()
}
{{else if .HasOne}}
func (m *{{$model.Name}}) {{.Arg | capitalize}}() ([]*{{.Arg | capitalize}}, error) {
	return {{.Arg | capitalize | singularize}}{}.Where("{{$model.TableName}}_id", m.{{$model.PrimaryKeyField}}).First()
}
{{else if .HasMany}}
func (m *{{$model.Name}}) {{.Arg | capitalize}}() ([]*{{.Arg | capitalize | singularize}}, error) {
	return {{.Arg | capitalize | singularize}}{}.Where("{{$model.TableName}}_id", m.{{$model.PrimaryKeyField}}).Query()
}
{{else if .Scope}}
func (m {{$model.Name}}) {{.Arg}}(args ...interface{}) *{{$model.Name}}Relation {
        return {{$model.Name}}Scopes["{{.Arg}}"](m.newRelation(), args...)
}

func (r *{{$model.Name}}Relation) {{.Arg}}(args ...interface{}) *{{$model.Name}}Relation {
        return {{$model.Name}}Scopes["{{.Arg}}"](r, args...)
}
{{end}}
{{end}}
{{end}}
`
	t := template.New("t")
	t.Funcs(template.FuncMap{
		"capitalize":  capitalize,
		"joinColumn":  joinColumn,
		"joinField":   joinField,
		"singularize": inflector.Singularize,
		"pluralize":   inflector.Pluralize,
	})
	tpl := template.Must(t.Parse(tplText))
	if err := tpl.Execute(w, structs); err != nil {
		return err
	}
	return nil
}

func capitalize(s string) string {
	c := []rune(s)
	c[0] = unicode.ToUpper(c[0])
	return string(c)
}

func joinColumn(s []string) string {
	return fmt.Sprintf("\"%s\"", strings.Join(s, "\", \""))
}

func joinField(s []string) string {
	return strings.Join(s, ", ")
}

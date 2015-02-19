package gen

import "fmt"

type Templates []*Template

func (ts Templates) ToString() string {
	var template string
	for _, t := range ts {
		template = template + t.toDefine()
	}
	return template
}

type Template struct {
	Name string
	Text string
}

func (t Template) toDefine() string {
	return fmt.Sprintf("{{define \"%s\"}}%s{{end}}\n", t.Name, t.Text)
}

var structTemplates = Templates{
	fieldByName,
	create,
	save,
	sel,
	find,
	findBy,
	relation,
	query,
	queryRow,
	where,
	and,
	first,
	last,
	order,
	limit,
	offset,
	group,
	having,
	validation,
	hasMany,
	hasOne,
	belongsTo,
	joins,
	scope,
	explain,
	delete,
	update,
}

var structDb = `package {{.Package}}

import "database/sql"

var db *sql.DB

func Use(DB *sql.DB) {
	db = DB
}
`

var structTemplate = `package {{.Package}}

import (
	"fmt"

	"github.com/monochromegane/argen"
	"github.com/monochromegane/goban"
)

{{range .}}
{{template "Relation" .}}
{{template "Select" .}}
{{template "Find" .}}
{{template "FindBy" .}}
{{template "First" .}}
{{template "Last" .}}
{{template "Where" .}}
{{template "And" .}}
{{template "Order" .}}
{{template "Limit" .}}
{{template "Offset" .}}
{{template "Group" .}}
{{template "Having" .}}
{{template "Explain" .}}
{{template "Validation" .}}
{{range .Scope}}
{{template "Scope" .}}
{{end}}
{{range .HasMany}}
{{template "HasMany" .}}
{{end}}
{{range .HasOne}}
{{template "HasOne" .}}
{{end}}
{{range .BelongsTo}}
{{template "BelongsTo" .}}
{{end}}
{{range .Joins}}
{{template "Joins" .}}
{{end}}
{{template "Create" .}}
{{template "Save" .}}
{{template "Update" .}}
{{template "Delete" .}}
{{template "Query" .}}
{{template "QueryRow" .}}
{{template "FieldByName" .}}
{{end}}
` + structTemplates.ToString()

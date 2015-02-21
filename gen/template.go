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
	build,
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
	exists,
	validation,
	hasMany,
	hasOne,
	belongsTo,
	joinsHasAny,
	joinsBelongsTo,
	buildHasAny,
	scope,
	delete,
	destroy,
	update,
}

var structDb = `package {{.Package}}

import "database/sql"

var db *sql.DB

func Use(DB *sql.DB) {
	db = DB
}
`

var structLogger = `package {{.Package}}

import (
	"log"
	"os"

	"github.com/monochromegane/argen"
)

var logger = &ar.Logger{Logger: log.New(os.Stdout, "", 0)}

func LogMode(mode bool) {
	logger.LogMode = mode
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
{{template "Validation" .}}
{{range .Scope}}
{{template "Scope" .}}
{{end}}
{{range .HasMany}}
{{template "HasMany" .}}
{{template "JoinsHasAny" .}}
{{template "BuildHasAny" .}}
{{end}}
{{range .HasOne}}
{{template "HasOne" .}}
{{template "JoinsHasAny" .}}
{{template "BuildHasAny" .}}
{{end}}
{{range .BelongsTo}}
{{template "BelongsTo" .}}
{{template "JoinsBelongsTo" .}}
{{end}}
{{template "Build" .}}
{{template "Create" .}}
{{template "Save" .}}
{{template "Update" .}}
{{template "Destroy" .}}
{{template "Delete" .}}
{{template "Query" .}}
{{template "QueryRow" .}}
{{template "Exists" .}}
{{template "FieldByName" .}}
{{end}}
` + structTemplates.ToString()

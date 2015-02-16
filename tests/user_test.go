package tests

import (
	"database/sql"
	"log"
	"os"
	"reflect"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/monochromegane/argen"
)

func TestMain(m *testing.M) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	Use(db)
	sqlStmt := `
	create table users (id integer not null primary key, name text);
	create table posts (id integer not null primary key, user_id integer, name text);
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Fatal(err, sqlStmt)
	}

	os.Exit(m.Run())
}

func TestSelect(t *testing.T) {
	u := &User{Name: "test"}
	u.Save()
	defer User{}.DeleteAll()

	u, err := User{}.Select("id").First()
	assertError(t, err)

	if !ar.IsZero(u.Name) {
		t.Errorf("column value should be empty, but %s", u.Name)
	}
}

func TestFind(t *testing.T) {
	expect := &User{Name: "test"}
	expect.Save()
	defer User{}.DeleteAll()

	u, err := User{}.Find(1)
	assertError(t, err)
	assertEqualStruct(t, expect, u)
}

func TestFindBy(t *testing.T) {
	expect := &User{Name: "test"}
	expect.Save()
	defer User{}.DeleteAll()

	u, err := User{}.FindBy("name", "test")
	assertError(t, err)
	assertEqualStruct(t, expect, u)
}

func TestFirst(t *testing.T) {
	for _, name := range []string{"test1", "test2"} {
		u := &User{Name: name}
		u.Save()
	}
	defer User{}.DeleteAll()

	expect, _ := User{}.Where("name", "test1").QueryRow()

	u, err := User{}.First()
	assertError(t, err)
	assertEqualStruct(t, expect, u)
}

func TestLast(t *testing.T) {
	for _, name := range []string{"test1", "test2"} {
		u := &User{Name: name}
		u.Save()
	}
	defer User{}.DeleteAll()

	expect, _ := User{}.Where("name", "test2").QueryRow()

	u, err := User{}.Last()
	assertError(t, err)
	assertEqualStruct(t, expect, u)
}

func TestWhere(t *testing.T) {
	expect := &User{Name: "test"}
	expect.Save()
	defer User{}.DeleteAll()

	u, err := User{}.Where("name", "test").And("id", 1).QueryRow()

	assertError(t, err)
	assertEqualStruct(t, expect, u)
}

func TestOrder(t *testing.T) {
	expects := []string{"test1", "test2"}
	for _, name := range expects {
		u := &User{Name: name}
		u.Save()
	}
	defer User{}.DeleteAll()

	users, err := User{}.Order("name", "ASC").Query()

	assertError(t, err)
	for i, u := range users {
		if u.Name != expects[i] {
			t.Errorf("column value should be %v, but %v", expects[i], u.Name)
		}
	}
}

func TestLimitAndOffset(t *testing.T) {
	for _, name := range []string{"test1", "test2", "test3"} {
		u := &User{Name: name}
		u.Save()
	}
	defer User{}.DeleteAll()

	users, err := User{}.Limit(2).Offset(1).Order("name", "ASC").Query()

	assertError(t, err)
	expects := []string{"test2", "test3"}
	for i, u := range users {
		if u.Name != expects[i] {
			t.Errorf("column value should be %v, but %v", expects[i], u.Name)
		}
	}
}

func TestGroupByAndHaving(t *testing.T) {
	for _, name := range []string{"testA", "testB", "testB"} {
		u := &User{Name: name}
		u.Save()
	}
	defer User{}.DeleteAll()

	users, err := User{}.Group("name").Having("count(name)", 2).Query()

	assertError(t, err)
	expects := []string{"testB"}
	for i, u := range users {
		if u.Name != expects[i] {
			t.Errorf("column value should be %v, but %v", expects[i], u.Name)
		}
	}
}

func TestExplain(t *testing.T) {
	err := User{}.Where("name", "test").Explain()
	assertError(t, err)
}

func TestIsValid(t *testing.T) {
	p := &Post{Name: "abc"}
	_, errs := p.IsValid()

	if len(errs.Messages["name"]) != 2 {
		t.Errorf("errors count should be 2, but %d", len(errs.Messages["name"]))
	}
}

func TestCreate(t *testing.T) {
	u, errs := User{}.Create(UserParams{
		Name: "TestCreate",
	})
	assertErrors(t, errs)

	expect, _ := User{}.FindBy("name", "TestCreate")
	assertEqualStruct(t, expect, u)
}

func assertEqualStruct(t *testing.T, expect, actual interface{}) {
	if !reflect.DeepEqual(expect, actual) {
		t.Errorf("struct should be equal to %v, but %v", expect, actual)
	}
}

func assertErrors(t *testing.T, errs *ar.Errors) {
	if errs != nil {
		t.Errorf("errors should be nil, but %v", errs)
	}
}

func assertError(t *testing.T, err error) {
	if err != nil {
		t.Errorf("error should be nil, but %v", err)
	}
}

package main

import (
	"fmt"
	"strings"

	"github.com/pingcap/parser"
	"github.com/pingcap/parser/ast"
	"github.com/pingcap/parser/format"
	_ "github.com/pingcap/parser/test_driver"
)

var (
	query = `CREATE FUNCTION say_hello(s text) RETURNS text DETERMINISTIC RETURN CONCAT('hello ', s);`
	q2    = `CREATE FUNCTION say_hello(s text) RETURNS text DETERMINISTIC`

	//q3 = `CREATE FUNCTION sf1 (p1 BIGINT) RETURNS BIGINT DETERMINISTIC MODIFIES SQL DATA;`
	q3 = `CREATE FUNCTION sf1 (p1 BIGINT) RETURNS BIGINT NO SQL;`
	//q3 = `CREATE FUNCTION sf1 (p1 BIGINT) RETURNS BIGINT CONTAINS SQL;`

	f1 = `CREATE FUNCTION sf1 (p1 BIGINT) RETURNS BIGINT BEGIN DECLARE ret INT DEFAULT 0\; SELECT c1*2 INTO ret FROM t1 WHERE c1 = p1\; RETURN ret\; END`
)

func parse(sql string) (*ast.StmtNode, error) {
	p := parser.New()

	node, err := p.ParseOneStmt(sql, "", "")
	if err != nil {
		return nil, err
	}

	return &node, nil
}

func main() {

	parser2 := parser.New()
	var err error
	var createNode ast.StmtNode
	createNode, err = parser2.ParseOneStmt(f1, "", "")
	if err != nil {
		fmt.Println(err)
		return
	}
	createStmt := createNode.(*ast.CreateFunctionStmt)
	fmt.Println(createStmt.Text())

	var newCreateSQLBuilder strings.Builder
	restoreCtx := format.NewRestoreCtx(format.DefaultRestoreFlags, &newCreateSQLBuilder)
	if err = createStmt.Restore(restoreCtx); err != nil {
		fmt.Println("restore err: ", err)
		return
	}
	newCreateSQL := newCreateSQLBuilder.String()
	fmt.Println("parsered sql: ", newCreateSQL)
}

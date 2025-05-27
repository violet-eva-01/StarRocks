// Package sql @author: Violet-Eva @date  : 2025/5/8 @notes :
package sql

import (
	"fmt"
	"github.com/violet-eva-01/StarRocks/sql/ast"
	"testing"
)

func TestAlterRenameParser(t *testing.T) {
	query := `
alter table test1 rename test2;
`
	parser := NewParser(query, "dc", "ddb")
	for priv, tbs := range parser.GetActionTables() {
		for tb, tbCount := range tbs {
			fmt.Printf("priv[%s] tb[%+v] count[%d]\n", priv, tb.String(), tbCount)
		}
	}
}

func TestWithCreateParser(t *testing.T) {
	query1 := `
create table aa(
id int
)
`
	ast.SetCaseSensitive()
	parser1 := NewParser(query1, "dc", "ddb")
	for priv, tbs := range parser1.GetActionTables() {
		for tb, tbCount := range tbs {
			fmt.Printf("priv[%s] tb[%s] count[%d]\n", priv, tb.String(), tbCount)
		}
	}
	query2 := `
create table aa as 
with tmpA as (
select
* from violet.eva.test1 t1 left join tmpA on tmpA.id = t1.id and tmpA.name = t1.name
)
select * from tmpA;
`
	ast.SetCaseSensitive()
	parser2 := NewParser(query2, "dc", "ddb")
	for priv, tbs := range parser2.GetActionTables() {
		for tb, tbCount := range tbs {
			fmt.Printf("priv[%s] tb[%s] count[%d]\n", priv, tb.String(), tbCount)
		}
	}
	query3 := `
create table aa like b
`
	ast.SetCaseSensitive()
	parser3 := NewParser(query3, "dc", "ddb")
	for priv, tbs := range parser3.GetActionTables() {
		for tb, tbCount := range tbs {
			fmt.Printf("priv[%s] tb[%s] count[%d]\n", priv, tb.String(), tbCount)
		}
	}
}

func TestShowCreateTable(t *testing.T) {
	query := `
show create table aa;
`
	parser := NewParser(query, "dc", "ddb")
	for priv, tbs := range parser.GetActionTables() {
		for tb, tbCount := range tbs {
			fmt.Printf("priv[%s] tb[%+v] count[%d]\n", priv, tb.String(), tbCount)
		}
	}
}

func TestWithDeleteParser(t *testing.T) {
	query := `
WITH foo_producers as (
    SELECT * from producers
    where producers.name = 'foo'
)
DELETE FROM films USING foo_producers
WHERE producer_id = foo_producers.id;
`
	ast.SetCaseSensitive()
	parser := NewParser(query, "dc", "ddb")
	for priv, tbs := range parser.GetActionTables() {
		for tb, tbCount := range tbs {
			fmt.Printf("priv[%s] tb[%s] count[%d]\n", priv, tb.String(), tbCount)
		}
	}
}

func TestTruncateParser(t *testing.T) {
	query := `
truncate table aa;
`
	ast.SetCaseSensitive()
	parser := NewParser(query, "dc", "ddb")
	for priv, tbs := range parser.GetActionTables() {
		for tb, tbCount := range tbs {
			fmt.Printf("priv[%s] tb[%s] count[%d]\n", priv, tb.String(), tbCount)
		}
	}
}

func TestDropParser(t *testing.T) {
	query := `
drop table aa force;
`
	ast.SetCaseSensitive()
	parser := NewParser(query, "dc", "ddb")
	for priv, tbs := range parser.GetActionTables() {
		for tb, tbCount := range tbs {
			fmt.Printf("priv[%s] tb[%s] count[%d]\n", priv, tb.String(), tbCount)
		}
	}
}

func TestExport(t *testing.T) {
	query := `
EXPORT TABLE V.E.testTbl 
TO "hdfs://<hdfs_host>:<hdfs_port>/a/b/c/" 
WITH BROKER
(
    "username"="xxx",
    "password"="yyy"
);
`
	ast.SetCaseSensitive()
	parser := NewParser(query, "violet", "eva")
	for priv, tbs := range parser.GetActionTables() {
		for tb, tbCount := range tbs {
			fmt.Printf("priv[%s] tb[%+v] count[%d]\n", priv, tb.String(), tbCount)
		}
	}
}

func TestWithInsertParser(t *testing.T) {
	query := `
insert into violet
with tmpA as (
select
*
from violet.eva.test1
), tmpB as (
select
*
from violet.eva.tmpA left join tmpA on tmpA.id = bbo.test
)
select * from violet.eva.test1 t1 left join tmpB on tmpB.id = t1.id;
`
	parser := NewParser(query, "dc", "ddb")
	for priv, tbs := range parser.GetActionTables() {
		for tb, tbCount := range tbs {
			fmt.Printf("priv[%s] tb[%+v] count[%d]\n", priv, tb.String(), tbCount)
		}
	}
}

func TestLoad(t *testing.T) {
	query1 := `
LOAD LABEL test_db.label10
(
    DATA INFILE("hdfs://<hdfs_host>:<hdfs_port>/user/starrocks/data/input/example10.csv")
    INTO TABLE table10
    COLUMNS TERMINATED BY ","
    (id, temp1, temp2)
    SET
    (
        col1 = hll_hash(temp1),
        col2 = hll_hash(temp2),
        col3 = empty_hll()
    )
)
WITH BROKER
(
    "username" = "<hdfs_username>",
    "password" = "<hdfs_password>"
);
`
	parser1 := NewParser(query1, "", "violet")
	for priv, tbs := range parser1.GetActionTables() {
		for tb, tbCount := range tbs {
			fmt.Printf("priv[%s] tb[%+v] count[%d]\n", priv, tb.String(), tbCount)
		}
	}
	query2 := `
LOAD LABEL db1.label1
(
DATA FROM TABLE hive_t1
INTO TABLE tbl1
SET
(
uuid=bitmap_dict(uuid)
)
)
WITH RESOURCE 'my_spark';
`
	parser2 := NewParser(query2, "violet", "eva")
	for priv, tbs := range parser2.GetActionTables() {
		for tb, tbCount := range tbs {
			fmt.Printf("priv[%s] tb[%+v] count[%d]\n", priv, tb.String(), tbCount)
		}
	}
}

func TestWithSelectParser(t *testing.T) {
	query := `
with tmpA as (
select
*
from violet.eva.test1
), tmpB as (
select
*
from violet.eva.tmpA left join tmpA on tmpA.id = bbo.test
)
select * from violet.eva.test1 t1 left join tmpB on tmpB.id = t1.id;
`
	parser := NewParser(query, "dc", "ddb")
	for priv, tbs := range parser.GetActionTables() {
		for tb, tbCount := range tbs {
			fmt.Printf("priv[%s] tb[%+v] count[%d]\n", priv, tb.String(), tbCount)
		}
	}
}

func TestWithUpdateParser(t *testing.T) {
	query := `
WITH AvgSalary AS (
    SELECT AVG(Salary) AS AverageSalary
    FROM Employees
)
UPDATE Employees
SET Salary = Salary * 1.1  -- 将薪水增加10%
FROM AvgSalary
WHERE Employees.Salary < AvgSalary.AverageSalary;
`
	ast.SetCaseSensitive()
	parser := NewParser(query, "dc", "ddb")
	for priv, tbs := range parser.GetActionTables() {
		for tb, tbCount := range tbs {
			fmt.Printf("priv[%s] tb[%s] count[%d]\n", priv, tb.String(), tbCount)
		}
	}
}

func TestRecoverParser(t *testing.T) {

	query1 := `
RECOVER TABLE aaa.aaa;
`
	ast.SetCaseSensitive()
	parser1 := NewParser(query1, "dc", "ddb")
	for priv, tbs := range parser1.GetActionTables() {
		for tb, tbCount := range tbs {
			fmt.Printf("priv[%s] tb[%s] count[%d]\n", priv, tb.String(), tbCount)
		}
	}
	query2 := `
RECOVER Database eva;
`
	ast.SetCaseSensitive()
	parser2 := NewParser(query2, "dc", "ddb")
	for priv, tbs := range parser2.GetActionTables() {
		for tb, tbCount := range tbs {
			fmt.Printf("priv[%s] tb[%s] count[%d]\n", priv, tb.String(), tbCount)
		}
	}
}

func TestDescParser(t *testing.T) {
	query1 := `
desc violet.eva.test1
`
	ast.SetCaseSensitive()
	parser1 := NewParser(query1, "dc", "ddb")
	for priv, tbs := range parser1.GetActionTables() {
		for tb, tbCount := range tbs {
			fmt.Printf("priv[%s] tb[%s] count[%d]\n", priv, tb.String(), tbCount)
		}
	}
}

func TestShowColumn(t *testing.T) {
	query1 := `
show full columns from violet.eva.test1;
`
	ast.SetCaseSensitive()
	parser1 := NewParser(query1, "dc", "ddb")
	for priv, tbs := range parser1.GetActionTables() {
		for tb, tbCount := range tbs {
			fmt.Printf("priv[%s] tb[%s] count[%d]\n", priv, tb.String(), tbCount)
		}
	}
}

func TestRefresh(t *testing.T) {
	query1 := `
refresh external table violet.eva.test1;
`
	ast.SetCaseSensitive()
	parser1 := NewParser(query1, "dc", "ddb")
	for priv, tbs := range parser1.GetActionTables() {
		for tb, tbCount := range tbs {
			fmt.Printf("priv[%s] tb[%s] count[%d]\n", priv, tb.String(), tbCount)
		}
	}
	query2 := `
REFRESH MATERIALIZED VIEW eva.test2;
`
	ast.SetCaseSensitive()
	parser2 := NewParser(query2, "dc", "ddb")
	for priv, tbs := range parser2.GetActionTables() {
		for tb, tbCount := range tbs {
			fmt.Printf("priv[%s] tb[%s] count[%d]\n", priv, tb.String(), tbCount)
		}
	}
}

func TestCreateView(t *testing.T) {
	query1 := `
CREATE VIEW example_db.example_view (k1, k2, k3, v1)
AS
SELECT c1 as k1, k2, k3, SUM(v1) FROM example_table
WHERE k1 = 20160112 GROUP BY k1,k2,k3;
`
	ast.SetCaseSensitive()
	parser1 := NewParser(query1, "dc", "ddb")
	for priv, tbs := range parser1.GetActionTables() {
		for tb, tbCount := range tbs {
			fmt.Printf("priv[%s] tb[%s] count[%d]\n", priv, tb.String(), tbCount)
		}
	}
	query2 := `
CREATE VIEW example_db.example_view
(
    k1 COMMENT 'first key',
    k2 COMMENT 'second key',
    k3 COMMENT 'third key',
    v1 COMMENT 'first value'
)
COMMENT 'my first view'
AS
SELECT c1 as k1, k2, k3, SUM(v1) FROM example_table
WHERE k1 = 20160112 GROUP BY k1,k2,k3;
`
	ast.SetCaseSensitive()
	parser2 := NewParser(query2, "dc", "ddb")
	for priv, tbs := range parser2.GetActionTables() {
		for tb, tbCount := range tbs {
			fmt.Printf("priv[%s] tb[%s] count[%d]\n", priv, tb.String(), tbCount)
		}
	}
	query3 := `
CREATE VIEW example_db.example_view (k1, k2, k3, v1)
COMMENT 'my secure view'
SECURITY INVOKER
AS
SELECT c1 as k1, k2, k3, SUM(v1) FROM example_table
WHERE k1 = 20160112 GROUP BY k1,k2,k3;
`
	ast.SetCaseSensitive()
	parser3 := NewParser(query3, "dc", "ddb")
	for priv, tbs := range parser3.GetActionTables() {
		for tb, tbCount := range tbs {
			fmt.Printf("priv[%s] tb[%s] count[%d]\n", priv, tb.String(), tbCount)
		}
	}
}

func TestAlterView(t *testing.T) {
	query1 := `
ALTER VIEW example_db.example_view
(
    c1 COMMENT "column 1",
    c2 COMMENT "column 2",
    c3 COMMENT "column 3"
)
AS SELECT k1, k2, SUM(v1) 
FROM example_table
GROUP BY k1, k2;
`
	ast.SetCaseSensitive()
	parser1 := NewParser(query1, "dc", "ddb")
	for priv, tbs := range parser1.GetActionTables() {
		for tb, tbCount := range tbs {
			fmt.Printf("priv[%s] tb[%s] count[%d]\n", priv, tb.String(), tbCount)
		}
	}
}

func TestDropView(t *testing.T) {
	query1 := `
DROP VIEW IF EXISTS example_db.example_view;
`
	ast.SetCaseSensitive()
	parser1 := NewParser(query1, "dc", "ddb")
	for priv, tbs := range parser1.GetActionTables() {
		for tb, tbCount := range tbs {
			fmt.Printf("priv[%s] tb[%s] count[%d]\n", priv, tb.String(), tbCount)
		}
	}
}

func TestShowCreateView(t *testing.T) {
	query1 := `
SHOW CREATE VIEW example_view;
`
	ast.SetCaseSensitive()
	parser1 := NewParser(query1, "dc", "ddb")
	for priv, tbs := range parser1.GetActionTables() {
		for tb, tbCount := range tbs {
			fmt.Printf("priv[%s] tb[%s] count[%d]\n", priv, tb.String(), tbCount)
		}
	}
}

func TestShowPartitions(t *testing.T) {
	query1 := `
SHOW PARTITIONS FROM test.site_access WHERE PartitionName = "p1";
`
	ast.SetCaseSensitive()
	parser1 := NewParser(query1, "dc", "ddb")
	for priv, tbs := range parser1.GetActionTables() {
		for tb, tbCount := range tbs {
			fmt.Printf("priv[%s] tb[%s] count[%d]\n", priv, tb.String(), tbCount)
		}
	}
	query2 := `
SHOW TEMPORARY PARTITIONS FROM test.site_access WHERE PartitionName = "p1";
`
	ast.SetCaseSensitive()
	parser2 := NewParser(query2, "dc", "ddb")
	for priv, tbs := range parser2.GetActionTables() {
		for tb, tbCount := range tbs {
			fmt.Printf("priv[%s] tb[%s] count[%d]\n", priv, tb.String(), tbCount)
		}
	}
}

func TestRecoverPartition(t *testing.T) {
	query1 := `
RECOVER PARTITION p1 FROM example_tbl;
`
	ast.SetCaseSensitive()
	parser1 := NewParser(query1, "dc", "ddb")
	for priv, tbs := range parser1.GetActionTables() {
		for tb, tbCount := range tbs {
			fmt.Printf("priv[%s] tb[%s] count[%d]\n", priv, tb.String(), tbCount)
		}
	}
}

func TestCreateIndex(t *testing.T) {
	query1 := `
CREATE INDEX index3 ON sales_records (item_id) USING BITMAP COMMENT '';
`
	ast.SetCaseSensitive()
	parser1 := NewParser(query1, "dc", "ddb")
	for priv, tbs := range parser1.GetActionTables() {
		for tb, tbCount := range tbs {
			fmt.Printf("priv[%s] tb[%s] count[%d]\n", priv, tb.String(), tbCount)
		}
	}
}

func TestDropIndex(t *testing.T) {
	query1 := `
DROP INDEX index3 ON sales_records;
`
	ast.SetCaseSensitive()
	parser1 := NewParser(query1, "dc", "ddb")
	for priv, tbs := range parser1.GetActionTables() {
		for tb, tbCount := range tbs {
			fmt.Printf("priv[%s] tb[%s] count[%d]\n", priv, tb.String(), tbCount)
		}
	}
}

func TestShowIndex(t *testing.T) {
	query1 := `
show index from eva.test1;
`
	ast.SetCaseSensitive()
	parser1 := NewParser(query1, "dc", "ddb")
	for priv, tbs := range parser1.GetActionTables() {
		for tb, tbCount := range tbs {
			fmt.Printf("priv[%s] tb[%s] count[%d]\n", priv, tb.String(), tbCount)
		}
	}
	query2 := `
show key from test.test2 from eva;
`
	ast.SetCaseSensitive()
	parser2 := NewParser(query2, "dc", "ddb")
	for priv, tbs := range parser2.GetActionTables() {
		for tb, tbCount := range tbs {
			fmt.Printf("priv[%s] tb[%s] count[%d]\n", priv, tb.String(), tbCount)
		}
	}
}

func TestSubmitTask(t *testing.T) {
	query1 := `
SUBMIT TASK AS INSERT OVERWRITE tbl3 SELECT * FROM src_tbl;
`
	ast.SetCaseSensitive()
	parser1 := NewParser(query1, "dc", "ddb")
	for priv, tbs := range parser1.GetActionTables() {
		for tb, tbCount := range tbs {
			fmt.Printf("priv[%s] tb[%s] count[%d]\n", priv, tb.String(), tbCount)
		}
	}
}

func TestCreateMV(t *testing.T) {
	query1 := `
create materialized view k1_k2 as
select k1, k2 from duplicate_table;
`
	ast.SetCaseSensitive()
	parser1 := NewParser(query1, "dc", "ddb")
	for priv, tbs := range parser1.GetActionTables() {
		for tb, tbCount := range tbs {
			fmt.Printf("priv[%s] tb[%s] count[%d]\n", priv, tb.String(), tbCount)
		}
	}
}

func TestAlterMV(t *testing.T) {
	query1 := `
ALTER MATERIALIZED VIEW lo_mv1 RENAME lo_mv1_new_name;
`
	ast.SetCaseSensitive()
	parser1 := NewParser(query1, "dc", "ddb")
	for priv, tbs := range parser1.GetActionTables() {
		for tb, tbCount := range tbs {
			fmt.Printf("priv[%s] tb[%s] count[%d]\n", priv, tb.String(), tbCount)
		}
	}
}

func TestCancelRefreshMV(t *testing.T) {
	query1 := `
CANCEL REFRESH MATERIALIZED VIEW lo_mv1 FORCE;
`
	ast.SetCaseSensitive()
	parser1 := NewParser(query1, "dc", "ddb")
	for priv, tbs := range parser1.GetActionTables() {
		for tb, tbCount := range tbs {
			fmt.Printf("priv[%s] tb[%s] count[%d]\n", priv, tb.String(), tbCount)
		}
	}
}

func TestDropMV(t *testing.T) {
	query1 := `
DROP MATERIALIZED VIEW IF EXISTS k1_k2;
`
	ast.SetCaseSensitive()
	parser1 := NewParser(query1, "dc", "ddb")
	for priv, tbs := range parser1.GetActionTables() {
		for tb, tbCount := range tbs {
			fmt.Printf("priv[%s] tb[%s] count[%d]\n", priv, tb.String(), tbCount)
		}
	}
}

func TestRefreshMV(t *testing.T) {
	query1 := `
REFRESH MATERIALIZED VIEW lo_mv1 
PARTITION START ("2020-02-01") END ("2020-03-01") FORCE;
`
	ast.SetCaseSensitive()
	parser1 := NewParser(query1, "dc", "ddb")
	for priv, tbs := range parser1.GetActionTables() {
		for tb, tbCount := range tbs {
			fmt.Printf("priv[%s] tb[%s] count[%d]\n", priv, tb.String(), tbCount)
		}
	}
}

func TestShowCreateMV(t *testing.T) {
	query1 := `
SHOW CREATE MATERIALIZED VIEW lo_mv1
`
	ast.SetCaseSensitive()
	parser1 := NewParser(query1, "dc", "ddb")
	for priv, tbs := range parser1.GetActionTables() {
		for tb, tbCount := range tbs {
			fmt.Printf("priv[%s] tb[%s] count[%d]\n", priv, tb.String(), tbCount)
		}
	}
}

func TestCreateRoutineLoad(t *testing.T) {
	query1 := `
CREATE ROUTINE LOAD example_db.example_tbl1_ordertest1 ON example_tbl1
COLUMNS TERMINATED BY ",",
COLUMNS (order_id, pay_dt, customer_name, nationality, gender, price)
FROM KAFKA
(
    "kafka_broker_list" ="<kafka_broker1_ip>:<kafka_broker1_port>,<kafka_broker2_ip>:<kafka_broker2_port>",
    "kafka_topic" = "ordertest1",
    "kafka_partitions" ="0,1,2,3", -- 指定分区
    "kafka_offsets" = "1000, OFFSET_BEGINNING, OFFSET_END, 2000" -- 指定起始位点
);
`
	ast.SetCaseSensitive()
	parser1 := NewParser(query1, "dc", "ddb")
	for priv, tbs := range parser1.GetActionTables() {
		for tb, tbCount := range tbs {
			fmt.Printf("priv[%s] tb[%s] count[%d]\n", priv, tb.String(), tbCount)
		}
	}
}

func TestAnalyze(t *testing.T) {
	query1 := `
ANALYZE TABLE tbl_name(c1, c2, c3);
`
	ast.SetCaseSensitive()
	parser1 := NewParser(query1, "dc", "ddb")
	for priv, tbs := range parser1.GetActionTables() {
		for tb, tbCount := range tbs {
			fmt.Printf("priv[%s] tb[%s] count[%d]\n", priv, tb.String(), tbCount)
		}
	}
}

func TestCreateAnalyze(t *testing.T) {
	query1 := `
CREATE ANALYZE TABLE tbl_name(c1, c2, c3); 
`
	ast.SetCaseSensitive()
	parser1 := NewParser(query1, "dc", "ddb")
	for priv, tbs := range parser1.GetActionTables() {
		for tb, tbCount := range tbs {
			fmt.Printf("priv[%s] tb[%s] count[%d]\n", priv, tb.String(), tbCount)
		}
	}
	query2 := `
CREATE ANALYZE DATABASE db_name;
`
	ast.SetCaseSensitive()
	parser2 := NewParser(query2, "dc", "ddb")
	for priv, tbs := range parser2.GetActionTables() {
		for tb, tbCount := range tbs {
			fmt.Printf("priv[%s] tb[%s] count[%d]\n", priv, tb.String(), tbCount)
		}
	}
}

func TestDropStatus(t *testing.T) {
	query1 := `
DROP STATS b.a.test
`
	ast.SetCaseSensitive()
	parser1 := NewParser(query1, "dc", "ddb")
	for priv, tbs := range parser1.GetActionTables() {
		for tb, tbCount := range tbs {
			fmt.Printf("priv[%s] tb[%s] count[%d]\n", priv, tb.String(), tbCount)
		}
	}
	query2 := `
ANALYZE TABLE tbl_name DROP HISTOGRAM ON col_name1 , col_name2
`
	ast.SetCaseSensitive()
	parser2 := NewParser(query2, "dc", "ddb")
	for priv, tbs := range parser2.GetActionTables() {
		for tb, tbCount := range tbs {
			fmt.Printf("priv[%s] tb[%s] count[%d]\n", priv, tb.String(), tbCount)
		}
	}
}

func TestSetUserVar(t *testing.T) {
	query1 := `
set exec_mem_limit=1024*1024,exec_mem_limit=2048*10,@var1="test1";
`
	ast.SetCaseSensitive()
	parser1 := NewParser(query1, "dc", "ddb")
	for priv, tbs := range parser1.GetSystemVariables() {
		fmt.Println(priv, ":", tbs)
	}
	for priv, tbs := range parser1.GetUserProperties() {
		fmt.Println(priv, ":", tbs)
	}
	query2 := `
/* ApplicationName=DBeaver 22.0.1 - SQLEditor <Script-72.sql> */ SET SQL_SELECT_LIMIT=DEFAULT
`
	ast.SetCaseSensitive()
	parser2 := NewParser(query2, "dc", "ddb")
	for priv, tbs := range parser2.GetSystemVariables() {
		fmt.Println(priv, ":", tbs)
	}
	for priv, tbs := range parser2.GetUserProperties() {
		fmt.Println(priv, ":", tbs)
	}
	query3 := `
set @@sql_select_limit=DEFAULT;SET SQL_SELECT_LIMIT=DEFAULT;
`
	ast.SetCaseSensitive()
	ast.SetSysVarCaseSensitive()
	parser3 := NewParser(query3, "dc", "ddb")
	for priv, tbs := range parser3.GetSystemVariables() {
		fmt.Println(priv, ":", tbs)
	}
	for priv, tbs := range parser3.GetUserProperties() {
		fmt.Println(priv, ":", tbs)
	}
}

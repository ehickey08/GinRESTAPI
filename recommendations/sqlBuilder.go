package recommendations

import (
	"fmt"
	"strings"
)

var Sql = sqlBuilder(ColumnsWanted, TablesNeeded, fmt.Sprintf("dbo.HD2018"), " where main.sector < 3 and main.INSTCAT = 2 and main.STABBR = ?")

func buildColumns(columnsWanted []map[string]string) string {
	columnStrs := make([]string, len(columnsWanted))
	for i, col := range columnsWanted {
		for key, val := range col {
			columnStrs[i] = fmt.Sprintf("%s as %s", val, key)
		}
	}
	return strings.Join(columnStrs, ", ")
}

func buildJoins(tablesNeeded []map[string][]string) string {
	tableStrs := make([]string, len(tablesNeeded))
	for i, col := range tablesNeeded {
		for key, val := range col {
			dbName := val[0]
			matchOn := val[1]
			matchTo := val[2]
			tableStr := fmt.Sprintf("LEFT JOIN %s as %s ON %s = %s", dbName, key, matchOn, matchTo)
			if len(val) > 3 {
				tableStr += fmt.Sprintf(" %s", val[3])
			}
			tableStrs[i] = tableStr
		}
	}
	return strings.Join(tableStrs, " ")
}
func sqlBuilder(columns []map[string]string, tables []map[string][]string, mainTable string, addedFilter string) string {
	columnStrs := buildColumns(columns)
	tableStrs := buildJoins(tables)
	return fmt.Sprintf("SELECT distinct %s FROM %s as main %s %s", columnStrs, mainTable, tableStrs, addedFilter)
}

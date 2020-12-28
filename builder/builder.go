package builder

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"unicode"
)

type SQLBuilder struct {
	_select          string
	_insert          string
	_insertAll       string
	_update          string
	_delete          string
	_limit           string
	_orderBy         string
	_groupBy         string
	_table           string
	_join            string
	_where           string
	_having          string
	_insertParams    []interface{}
	_insertAllParams []interface{}
	_updateParams    []interface{}
	_whereParams     []interface{}
	_joinParams      []interface{}
	_havingParams    []interface{}
	_limitParams     []interface{}
}

var (
	ErrTableEmpty      = errors.New("table empty")
	ErrInsertStatement = errors.New("insert statement empty")
	ErrUpdateStatement = errors.New("update statement empty")
	ErrElementStatement = errors.New("element type error")
	m                  reflect.Type
	ignoreKey = "Id"
)

func NewSQLBuilder() *SQLBuilder {
	return &SQLBuilder{}
}

//SELECT `t1`.`name`,`t1`.`age`,`t2`.`teacher`,`t3`.`address` FROM `test` as t1 LEFT
//JOIN `test2` as `t2` ON `t1`.`class` = `t2`.`class` INNER JOIN `test3` as t3 ON
//`t1`.`school` = `t3`.`school` WHERE `t1`.`age` >= 20 GROUP BY `t1`.`age`
//HAVING COUNT(`t1`.`age`) > 2 ORDER BY `t1`.`age` DESC LIMIT ?,?

func (sb *SQLBuilder) GetQuerySql() (string, error) {
	if sb._table == "" {
		return "", ErrTableEmpty
	}

	var buf strings.Builder

	buf.WriteString("SELECT ")
	if sb._select != "" {
		buf.WriteString(sb._select)
	} else {
		buf.WriteString("*")
	}

	buf.WriteString(" FROM ")
	buf.WriteString(sb._table)

	if sb._join != "" {
		buf.WriteString(" ")
		buf.WriteString(sb._join)
	}

	if sb._where != "" {
		buf.WriteString(" ")
		buf.WriteString(sb._where)
	}

	if sb._groupBy != "" {
		buf.WriteString(" ")
		buf.WriteString(sb._groupBy)
	}

	if sb._having != "" {
		buf.WriteString(" ")
		buf.WriteString(sb._having)
	}

	if sb._orderBy != "" {
		buf.WriteString(" ")
		buf.WriteString(sb._orderBy)
	}

	if sb._limit != "" {
		buf.WriteString(" ")
		buf.WriteString(sb._limit)
	}

	return buf.String(), nil
}

func (sb *SQLBuilder) GetQueryParams() []interface{} {
	Params := []interface{}{}
	Params = append(Params, sb._joinParams...)
	Params = append(Params, sb._whereParams...)
	Params = append(Params, sb._havingParams...)
	Params = append(Params, sb._limitParams...)

	return Params
}

func (sb *SQLBuilder) Select(cols ...string) *SQLBuilder {
	var buf strings.Builder

	for key, col := range cols {
		buf.WriteString(col)
		if key != len(cols)-1 {
			buf.WriteString(",")
		}
	}

	sb._select = buf.String()

	return sb
}

func (sb *SQLBuilder) Table(table string) *SQLBuilder {
	sb._table = table
	return sb
}

func (sb *SQLBuilder) Insert(fields []string, values ...interface{}) *SQLBuilder {
	var buf strings.Builder

	buf.WriteString("(")

	for key, field := range fields {
		buf.WriteString(field)
		if key != len(fields)-1 {
			buf.WriteString(",")
		}
	}
	buf.WriteString(") VALUES (")

	for key := range fields {
		buf.WriteString("?")
		if key != len(fields)-1 {
			buf.WriteString(",")
		}
	}

	buf.WriteString(")")
	sb._insert = buf.String()

	for _, value := range values {
		sb._insertParams = append(sb._insertParams, value)
	}

	return sb
}

func (sb *SQLBuilder) GetInsertSql() (string, error) {
	if sb._table == "" {
		return "", ErrTableEmpty
	}

	if sb._insert == "" {
		return "", ErrInsertStatement
	}

	var buf strings.Builder
	buf.WriteString("INSERT INTO ")
	buf.WriteString(sb._table)
	buf.WriteString(" ")
	buf.WriteString(sb._insert)

	return buf.String(), nil
}

func (sb *SQLBuilder) GetInsertParams() []interface{} {
	return sb._insertParams
}

func (sb *SQLBuilder) Update(fields []string, values ...interface{}) *SQLBuilder {
	var buf strings.Builder

	for key, val := range fields {
		buf.WriteString(val)
		buf.WriteString(" = ")
		buf.WriteString("?")
		if key != len(fields)-1 {
			buf.WriteString(",")
		}
	}

	sb._update = buf.String()

	for _, val := range values {
		sb._updateParams = append(sb._updateParams, val)
	}

	return sb
}

func (sb *SQLBuilder) GetUpdateSql() (string, error) {
	if sb._table == "" {
		return "", ErrTableEmpty
	}

	if sb._update == "" {
		return "", ErrUpdateStatement
	}

	var buf strings.Builder
	buf.WriteString("UPDATE ")
	buf.WriteString(sb._table)
	buf.WriteString(" SET ")
	buf.WriteString(sb._update)
	buf.WriteString(" ")

	if sb._where != "" {
		buf.WriteString(sb._where)
	}

	return buf.String(), nil
}

func (sb *SQLBuilder) GetUpdateParams() []interface{} {
	Params := []interface{}{}
	Params = append(Params, sb._updateParams...)
	Params = append(Params, sb._whereParams...)
	return Params
}

func (sb *SQLBuilder) GetDeleteSql() (string, error) {
	if sb._table == "" {
		return "", ErrTableEmpty
	}

	var buf strings.Builder
	buf.WriteString("DELETE FROM ")
	buf.WriteString(sb._table)
	buf.WriteString(" ")

	if sb._where != "" {
		buf.WriteString(sb._where)
	}

	return buf.String(), nil
}

func (sb *SQLBuilder) GetDeleteParams() []interface{} {
	return sb._whereParams
}

func (sb *SQLBuilder) Where(field, condition string, value interface{}) *SQLBuilder {
	return sb.where("AND", field, condition, value)
}

func (sb *SQLBuilder) WhereOr(field string, condition string, value interface{}) *SQLBuilder {
	return sb.where("OR", field, condition, value)
}

func (sb *SQLBuilder) WhereRaw(s string, values []interface{}) *SQLBuilder {
	return sb.whereRaw("AND", s, values)
}

func (sb *SQLBuilder) WhereOrRaw(s string, values []interface{}) *SQLBuilder {
	return sb.whereRaw("OR", s, values)
}

func (sb *SQLBuilder) where(operator string, field string, condition string, value interface{}) *SQLBuilder {
	var buf strings.Builder

	buf.WriteString(sb._where)

	if buf.Len() == 0 {
		buf.WriteString("WHERE ")
	} else {
		buf.WriteString(" ")
		buf.WriteString(operator)
		buf.WriteString(" ")
	}

	buf.WriteString(field)
	buf.WriteString(" ")
	buf.WriteString(condition)
	buf.WriteString(" ")
	buf.WriteString("?")

	sb._where = buf.String()
	sb._whereParams = append(sb._whereParams, value)

	return sb
}

func (sb *SQLBuilder) whereRaw(operator string, s string, values []interface{}) *SQLBuilder {
	var buf strings.Builder
	buf.WriteString(sb._where)

	if buf.Len() == 0 {
		buf.WriteString("WHERE ")
	} else {
		buf.WriteString(" ")
		buf.WriteString(operator)
		buf.WriteString(" ")
	}

	buf.WriteString(s)
	sb._where = buf.String()

	for _, value := range values {
		sb._whereParams = append(sb._whereParams, value)
	}

	return sb
}

func (sb *SQLBuilder) WhereIn(field string, value []interface{}) *SQLBuilder {
	return sb.whereIn("AND", "IN", field, value)
}

func (sb *SQLBuilder) WhereNotIn(field string, value []interface{}) *SQLBuilder {
	return sb.whereIn("AND", "NOT IN", field, value)
}

func (sb *SQLBuilder) WhereOrIn(field string, value []interface{}) *SQLBuilder {
	return sb.whereIn("OR", "IN", field, value)
}

func (sb *SQLBuilder) WhereOrNotIn(field string, value []interface{}) *SQLBuilder {
	return sb.whereIn("OR", "NOT IN", field, value)
}

func (sb *SQLBuilder) whereIn(operator string, condition string, field string, values []interface{}) *SQLBuilder {
	var buf strings.Builder
	buf.WriteString(sb._where)

	if buf.Len() == 0 {
		buf.WriteString("WHERE ")
	} else {
		buf.WriteString(" ")
		buf.WriteString(operator)
		buf.WriteString(" ")
	}

	buf.WriteString(field)
	buf.WriteString(" ")
	buf.WriteString(condition)
	buf.WriteString(" ")
	buf.WriteString("(")
	for i := 0; i < len(values); i++ {
		buf.WriteString("?")
		if i != len(values)-1 {
			buf.WriteString(",")
		}
	}
	buf.WriteString(")")

	sb._where = buf.String()

	for _, val := range values {
		sb._whereParams = append(sb._whereParams, val)
	}

	return sb
}

func (sb *SQLBuilder) Limit(offset, num interface{}) *SQLBuilder {
	var buf strings.Builder
	buf.WriteString("LIMIT ?,?")
	sb._limit = buf.String()
	sb._limitParams = append(sb._limitParams, offset, num)
	return sb
}

func (sb *SQLBuilder) OrderBy(order string, fields ...string) *SQLBuilder {
	var buf strings.Builder
	buf.WriteString("ORDER BY ")
	for key, val := range fields {
		buf.WriteString(val)
		if key != len(fields)-1 {
			buf.WriteString(",")
		}
	}

	buf.WriteString(" ")
	buf.WriteString(order)

	sb._orderBy = buf.String()

	return sb
}

func (sb *SQLBuilder) GroupBy(field string) *SQLBuilder {
	var buf strings.Builder
	buf.WriteString("GROUP BY ")
	buf.WriteString(field)

	sb._groupBy = buf.String()
	return sb
}

func (sb *SQLBuilder) JoinRaw(s string, values ...interface{}) *SQLBuilder {
	var buf strings.Builder
	buf.WriteString(sb._join)

	if buf.Len() != 0 {
		buf.WriteString(" ")
	}
	buf.WriteString(s)
	sb._join = buf.String()

	for _, value := range values {
		sb._joinParams = append(sb._joinParams, value)
	}

	return sb
}

func (sb *SQLBuilder) HavingRaw(s string, values ...interface{}) *SQLBuilder {
	return sb.havingRaw("AND", s, values)
}

func (sb *SQLBuilder) HavingRawOr(s string, values ...interface{}) *SQLBuilder {
	return sb.havingRaw("OR", s, values)
}

func (sb *SQLBuilder) havingRaw(operator string, s string, values ...interface{}) *SQLBuilder {
	var buf strings.Builder
	buf.WriteString(sb._having)

	if buf.Len() == 0 {
		buf.WriteString("HAVING ")
	} else {
		buf.WriteString(" ")
		buf.WriteString(operator)
		buf.WriteString(" ")
	}
	buf.WriteString(s)
	sb._having = buf.String()

	for _, value := range values {
		sb._havingParams = append(sb._havingParams, value)
	}

	return sb
}

func (sb *SQLBuilder) Having(field string, condition string, value interface{}) *SQLBuilder {
	return sb.having("AND", field, condition, value)
}

func (sb *SQLBuilder) HavingOr(field string, condition string, value interface{}) *SQLBuilder {
	return sb.having("OR", field, condition, value)
}

func (sb *SQLBuilder) having(operator string, field string, condition string, value interface{}) *SQLBuilder {
	if sb._groupBy == "" {
		return sb
	}

	var buf strings.Builder
	buf.WriteString(sb._having)

	if buf.Len() == 0 {
		buf.WriteString("HAVING ")
	} else {
		buf.WriteString(" ")
		buf.WriteString(operator)
		buf.WriteString(" ")
	}

	buf.WriteString(field)
	buf.WriteString(" ")
	buf.WriteString(condition)
	buf.WriteString(" ")
	buf.WriteString("?")

	sb._having = buf.String()

	sb._havingParams = append(sb._havingParams, value)
	return sb
}

func (sb *SQLBuilder) InsertAll(fields []string, values ...[]interface{}) *SQLBuilder {
	var buf strings.Builder

	buf.WriteString("(")

	for key, field := range fields {
		buf.WriteString(field)
		if key != len(fields)-1 {
			buf.WriteString(",")
		}
	}
	buf.WriteString(") VALUES ")

	for key := range values {
		buf.WriteString("(")
		for key := range fields {
			buf.WriteString("?")
			if key != len(fields)-1 {
				buf.WriteString(",")
			}
		}
		buf.WriteString(")")
		if key != len(values)-1 {
			buf.WriteString(",")
		}
	}

	sb._insertAll = buf.String()

	for _, value := range values {
		sb._insertAllParams = append(sb._insertAllParams, value...)
	}

	return sb
}

func (sb *SQLBuilder) GetInsertAllSql() (string, error) {
	if sb._table == "" {
		return "", ErrTableEmpty
	}

	if sb._insertAll == "" {
		return "", ErrInsertStatement
	}

	var buf strings.Builder
	buf.WriteString("INSERT INTO ")
	buf.WriteString(sb._table)
	buf.WriteString(" ")
	buf.WriteString(sb._insertAll)

	return buf.String(), nil
}

func (sb *SQLBuilder) GetInsertAllParams() []interface{} {
	return sb._insertAllParams
}

//传数组或者切片,默认忽略Id字段，
//结构体字段添加标签`gdb:"ignore"`可以忽略该字段，
//添加`gdb:"require"`可以不忽略该字段
//字段示例：`gdb:"column:userAge;underline"`
func (sb *SQLBuilder) InsertAllModel(models interface{}) (string,error) {
	fields,values,err := sb.reflectElementInfo(models)
	if err != nil {
		return "",nil
	}

	return sb.InsertAll(fields,values...).GetInsertAllSql()
}

func (sb *SQLBuilder) reflectElementInfo(elem interface{}) ([]string,[][]interface{},error) {
	m = reflect.TypeOf(elem)
	values := make([]interface{}, 0)
	allValues := make([][]interface{}, 0)
	v := reflect.ValueOf(elem)
	switch m.Kind() {
	case reflect.Slice, reflect.Array:
		//fmt.Println(values)
		m = reflect.TypeOf(elem).Elem() //获取slice或array内的元素类型
		if m.Kind() == reflect.Ptr {
			//如果是指针，则指向其所指的元素
			m = m.Elem()
		}

		var e reflect.Value
		for i := 0; i < v.Len(); i++ {
			vals := []interface{}{}
			if v.Index(i).Kind() == reflect.Ptr {
				e = v.Index(i).Elem()
			}else{
				e = v.Index(i)
			}

			for i2 := 0; i2 < e.NumField(); i2++ {
				fn := m.Field(i2).Name
				tag := m.Field(i2).Tag.Get("gdb")
				if getTagIsContinue(fn, tag) {
					continue
				}
				vals = append(vals, e.FieldByName(fn).Interface())
			}
			allValues = append(allValues, vals)
		}
		break
	case reflect.Ptr:
		//如果是指针，则指向其所指的对象
		m = reflect.TypeOf(elem).Elem()
		values = getParams(v.Elem(), m)
		allValues = append(allValues, values)
		break
	case reflect.Struct:
		values = getParams(v, m)
		allValues = append(allValues, values)
		break
	default:
		return nil,nil,ErrElementStatement
	}

	fields := make([]string, 0)
	for i := 0; i < m.NumField(); i++ {
		fn := m.Field(i).Name
		tag := m.Field(i).Tag.Get("gdb")
		fn = getFnForTag(fn, tag)
		if getTagIsContinue(fn, tag) {
			continue
		}
		fields = append(fields, fmt.Sprintf("%s", fn))
	}

	return fields,allValues,nil
}

func getParams(v reflect.Value, m reflect.Type) []interface{} {
	values := make([]interface{}, 0)
	for i := 0; i < m.NumField(); i++ {
		fn := m.Field(i).Name
		if fn == "Id" {
			continue
		}
		values = append(values, v.FieldByName(fn).Interface())
	}

	return values
}

//根据标签处理字段名
func getFnForTag(fn, tag string) string {
	fn = getTagColumnName(fn, tag)
	fn = getTagUnderLine(fn, tag)

	return fn
}

//是否跳过
func getTagIsContinue(fn, tag string) bool {
	if strings.Contains(tag, "ignore") {
		return true
	}

	if strings.Contains(tag, "require") {
		return false
	}
	//fmt.Println("tag",tag,"unignore", strings.Contains(tag, "unignore"))
	if strings.Contains(ignoreKey, fn){
		return true
	}

	return false
}

//获取标签中的字段名
func getTagColumnName(fn, tag string) string {
	tags := strings.Split(tag, ";")
	for _, val := range tags {
		if strings.Contains(val, "column") {
			columnName := strings.Split(val, ":")[1]
			return columnName
		}
	}

	return fn
}

//字段名转下划线，在获取标签字段名之后处理
func getTagUnderLine(fn, tag string) string  {
	var buffer bytes.Buffer
	if strings.Contains(tag, "underline") {
		//Camel2Case
		for i, i2 := range fn {
			if unicode.IsUpper(i2){
				if i != 0 {
					buffer.WriteString("_")
				}

				buffer.WriteString(string(unicode.ToLower(i2)))
			}else{
				buffer.WriteString(string(i2))
			}
		}
		return buffer.String()
	}

	return fn
}
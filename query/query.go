package query

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

const (
	QuerySelect = 1
	QueryUpdate = 2
	QueryInsert = 3
	QueryDelete = 4
)

func ComposeQuery(queryType int, input interface{}, looksValue bool) (string, error) {

	switch queryType {
	case QuerySelect:
		return ComposeSelectQuery(input, looksValue)
	case QueryUpdate:
		return ComposeUpdateQuery(input, looksValue)
	case QueryInsert:
		return ComposeInsertQuery(input, looksValue)
	case QueryDelete:
		return ComposeDeleteQuery(input, looksValue)
	default:
		return "", fmt.Errorf("unsupported query type: %d", queryType)
	}

}

func ComposeSelectQuery(input interface{}, looksValue bool) (string, error) {
	ivalue := reflect.Indirect(reflect.ValueOf(input))
	itype := ivalue.Type()

	var strQuery *bytes.Buffer

	if ivalue.Kind() == reflect.Struct {
		queryType := NewQueryDataType()

		if err := queryType.Iterate(input, looksValue); err != nil {
			return "", err
		}

		// generate query
		strQuery = bytes.NewBufferString("SELECT ")
		if len(queryType.Cols) > 0 {
			strQuery.WriteString(strings.Join(queryType.Cols, ", "))
		} else {
			return "", errors.New("invalid struct format, please check your tag columns")
		}
		method := ivalue.MethodByName("Table")
		if method.IsValid() {
			vals := method.Call([]reflect.Value{})
			if len(vals) > 0 {
				nm := vals[0].String()
				strQuery.WriteString(" FROM " + strings.ToLower(nm))
			}
		} else {
			strQuery.WriteString(" FROM " + strings.ToLower(itype.Name()) + "s")
		}

		var globalIdx int

		if len(queryType.Foreigns) > 0 && len(queryType.Foreigns) > 0 && len(queryType.Primaries) == len(queryType.Primaries) {
			if !strings.Contains(strQuery.String(), " WHERE ") {
				strQuery.WriteString(" WHERE ")
			}

			for _, p := range queryType.Primaries {
				psplits := strings.Split(p, "#")
				for idx, f := range queryType.Foreigns {
					globalIdx += idx
					fsplits := strings.Split(f, "#")
					if psplits[1] == fsplits[1] {
						pair := fmt.Sprintf("%s = %s", psplits[0], fsplits[0])
						if globalIdx == 0 {
							strQuery.WriteString(pair)
							globalIdx += 1
						} else {
							strQuery.WriteString(" AND " + pair)
						}
					}
				}
			}
		}

		if len(queryType.Wheres) > 0 {
			if !strings.Contains(strQuery.String(), " WHERE ") {
				strQuery.WriteString(" WHERE ")
			}
			for idx, s := range queryType.Wheres {
				globalIdx += idx
				sv, found := queryType.WhereVal[s]
				var dt string
				switch sv.(type) {
				case string:
					dt = fmt.Sprintf("'%v'", sv)
				default:
					dt = fmt.Sprintf("%v", sv)
				}
				if globalIdx == 0 {
					if found {
						strQuery.WriteString(s + " = " + dt)
					} else {
						strQuery.WriteString(s + " = ?")
					}
				} else {
					if found {
						strQuery.WriteString(" AND " + s + " = " + dt)
					} else {
						strQuery.WriteString(" AND " + s + " = ?")
					}
				}
			}
		}

		if len(queryType.OrderBy) > 0 {
			strQuery.WriteString(" ORDER BY ")
			strQuery.WriteString(strings.Join(queryType.OrderBy, ", "))
		}

		return strQuery.String(), nil

	} else {
		return "", errors.New("only accepted struct type")
	}

}

func ComposeUpdateQuery(input interface{}, looksValue bool) (string, error) {
	ivalue := reflect.Indirect(reflect.ValueOf(input))
	itype := ivalue.Type()

	var strQuery *bytes.Buffer

	if ivalue.Kind() == reflect.Struct {
		queryType := NewQueryDataType()
		if err := queryType.Iterate(input, looksValue); err != nil {
			return "", err
		}

		// generate query
		strQuery = bytes.NewBufferString("UPDATE ")
		method := ivalue.MethodByName("Table")
		if method.IsValid() {
			vals := method.Call([]reflect.Value{})
			if len(vals) > 0 {
				nm := vals[0].String()
				strQuery.WriteString(strings.ToLower(nm) + " ")
			}
		} else {
			strQuery.WriteString(strings.ToLower(itype.Name()) + "s ")
		}
		if len(queryType.Cols) > 0 {
			strQuery.WriteString("SET ")
			for idx, s := range queryType.Cols {
				sv, found := queryType.ColVal[s]

				var dt string
				switch sv.(type) {
				case string:
					dt = fmt.Sprintf("'%v'", sv)
				default:
					dt = fmt.Sprintf("%v", sv)
				}

				if idx == 0 {
					if found {
						strQuery.WriteString(s + " = " + dt)
					} else {
						strQuery.WriteString(s + " = ?")
					}
				} else {
					if found {
						strQuery.WriteString(", " + s + " = " + dt)
					} else {
						strQuery.WriteString(", " + s + " = ?")
					}
				}

			}
		}
		if len(queryType.Wheres) > 0 {
			if !strings.Contains(strQuery.String(), " WHERE ") {
				strQuery.WriteString(" WHERE ")
			}
			for idx, s := range queryType.Wheres {
				sv, found := queryType.WhereVal[s]

				var dt string
				switch sv.(type) {
				case string:
					dt = fmt.Sprintf("'%v'", sv)
				default:
					dt = fmt.Sprintf("%v", sv)
				}
				if idx == 0 {
					if found {
						strQuery.WriteString(s + " = " + dt)
					} else {
						strQuery.WriteString(s + " = ?")
					}
				} else {
					if found {
						strQuery.WriteString(" AND " + s + " = " + dt)
					} else {
						strQuery.WriteString(" AND " + s + " = ?")
					}
				}
			}
		}

		return strQuery.String(), nil

	} else {
		return "", errors.New("only accepted struct type")
	}
}

func ComposeInsertQuery(input interface{}, looksValue bool) (string, error) {
	ivalue := reflect.Indirect(reflect.ValueOf(input))
	itype := ivalue.Type()

	var strQuery *bytes.Buffer

	if ivalue.Kind() == reflect.Struct {
		queryType := NewQueryDataType()
		if err := queryType.Iterate(input, looksValue); err != nil {
			return "", err
		}

		// generate query
		strQuery = bytes.NewBufferString("INSERT INTO ")
		method := ivalue.MethodByName("Table")
		if method.IsValid() {
			vals := method.Call([]reflect.Value{})
			if len(vals) > 0 {
				nm := vals[0].String()
				strQuery.WriteString(strings.ToLower(nm) + "s")
			}
		} else {
			strQuery.WriteString(strings.ToLower(itype.Name()) + "s")
		}

		if len(queryType.Cols) > 0 {
			strQuery.WriteString(" ( ")
			strQuery.WriteString(strings.Join(queryType.Cols, ", "))
			strQuery.WriteString(" ) ")
		}

		strQuery.WriteString("VALUES ( ")
		for idx, s := range queryType.Cols {
			sv, found := queryType.ColVal[s]

			var dt string
			switch sv.(type){
			case string:
				dt = fmt.Sprintf("'%v'", sv)
			default:
				dt = fmt.Sprintf("%v", sv)
			}

			if idx == 0 {
				if found {
					strQuery.WriteString(dt)
				} else {
					strQuery.WriteString("?")
				}
			} else {
				if found {
					strQuery.WriteString(", " + dt)
				} else {
					strQuery.WriteString(", ?")
				}
			}
		}
		strQuery.WriteString(" )")

		return strQuery.String(), nil
	} else {
		return "", errors.New("only accepted struct type")
	}
}

func ComposeDeleteQuery(input interface{}, looksValue bool) (string, error) {
	ivalue := reflect.Indirect(reflect.ValueOf(input))
	itype := ivalue.Type()

	var strQuery *bytes.Buffer

	if ivalue.Kind() == reflect.Struct {
		queryType := NewQueryDataType()
		if err := queryType.Iterate(input, looksValue); err != nil {
			return "", err
		}

		// generate query
		strQuery = bytes.NewBufferString("DELETE FROM ")
		method := ivalue.MethodByName("Table")
		if method.IsValid() {
			vals := method.Call([]reflect.Value{})
			if len(vals) > 0 {
				nm := vals[0].String()
				strQuery.WriteString(strings.ToLower(nm))
			}
		} else {
			strQuery.WriteString(strings.ToLower(itype.Name()) + "s")
		}
		if len(queryType.Wheres) > 0 {
			if !strings.Contains(strQuery.String(), " WHERE ") {
				strQuery.WriteString(" WHERE ")
			}
			for idx, s := range queryType.Wheres {
				sv, found := queryType.WhereVal[s]

				var dt string
				switch sv.(type) {
				case string:
					dt = fmt.Sprintf("'%v'", sv)
				default:
					dt = fmt.Sprintf("%v", sv)
				}
				if idx == 0 {
					if found {
						strQuery.WriteString(s + " = " + dt)
					} else {
						strQuery.WriteString(s + " = ?")
					}
				} else {
					if found {
						strQuery.WriteString(" AND " + s + " = " + dt)
					} else {
						strQuery.WriteString(" AND " + s + " = ?")
					}
				}
			}
		}

		return strQuery.String(), nil
	} else {
		return "", errors.New("only accepted struct type")
	}
}

func IsEmpty(v interface{}) bool {
	return reflect.DeepEqual(v, reflect.Zero(reflect.TypeOf(v)).Interface())
}

type QueryDataType struct {
	Cols      []string
	Wheres    []string
	OrderBy   []string
	Primaries []string
	Foreigns  []string
	ColVal    map[string]interface{}
	WhereVal  map[string]interface{}
}

func NewQueryDataType() *QueryDataType {
	q := new(QueryDataType)
	q.ColVal = make(map[string]interface{})

	return q
}

func (q *QueryDataType) Iterate(input interface{}, looksValue bool) error {
	ivalue := reflect.Indirect(reflect.ValueOf(input))
	itype := ivalue.Type()

	if ivalue.Kind() == reflect.Struct {
		var tmpCols string

		for i := 0; i < ivalue.NumField(); i++ {
			ftype := itype.Field(i)
			fval := ivalue.Field(i)

			tags := ftype.Tag.Get("query")

			if tags != "" {
				stags := strings.Split(tags, ",")
				for idx, s := range stags {
					if idx == 0 {
						tmpCols = s
					} else if s == "col" {
						q.Cols = append(q.Cols, tmpCols)
						if ftype.Type.Kind() == reflect.Struct || ftype.Type.Kind() == reflect.Map || ftype.Type.Kind() == reflect.Slice {
							return errors.New("unsupported struct with field of type collection")
						} else {
							if looksValue {
								if !IsEmpty(reflect.Indirect(fval).Interface()) {
									if q.ColVal == nil {
										q.ColVal = make(map[string]interface{})
									}
									ifield := reflect.Indirect(fval)
									q.ColVal[tmpCols] = ifield.Interface()
								}
							}
						}
					} else if s == "where" {
						q.Wheres = append(q.Wheres, tmpCols)
						if ftype.Type.Kind() == reflect.Struct || ftype.Type.Kind() == reflect.Map || ftype.Type.Kind() == reflect.Slice {
							return errors.New("unsupported struct with field of type collection")
						} else {
							if looksValue {
								if !IsEmpty(reflect.Indirect(fval).Interface()) {
									if q.WhereVal == nil {
										q.WhereVal = make(map[string]interface{})
									}
									ifield := reflect.Indirect(fval)
									q.WhereVal[tmpCols] = ifield.Interface()
								}
							}
						}
					} else if strings.Contains(s, "order") {
						ssplit := strings.Split(s, "#")
						if len(ssplit) == 2 {
							q.OrderBy = append(q.OrderBy, tmpCols+" "+ssplit[1])
						}
					} else if strings.Contains(s, "primary") {
						ssplit := strings.Split(s, "#")
						if len(ssplit) == 2 {
							q.Primaries = append(q.Primaries, tmpCols+"#"+ssplit[1])
						}
					} else if strings.Contains(s, "foreign") {
						ssplit := strings.Split(s, "#")
						if len(ssplit) == 2 {
							q.Foreigns = append(q.Foreigns, tmpCols+"#"+ssplit[1])
						}
					}
				}
			}
		}
		return nil
	} else {
		return errors.New("only accepted struct type")
	}
}

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
	QueryCreate = 3
	QueryDelete = 4
)

func ComposeQuery(queryType int, input interface{}) (string, error) {

	switch queryType {
	case QuerySelect:
		return ComposeSelectQuery(input)
	case QueryUpdate:
		return ComposeUpdateQuery(input)
	case QueryCreate:
		return ComposeCreateQuery(input)
	case QueryDelete:
		return ComposeDeleteQuery(input)
	default:
		return "", fmt.Errorf("unsupported query type: %d", queryType)
	}

}

func ComposeSelectQuery(input interface{}) (string, error) {

	defer func() {
		if perr := recover(); perr != nil {
			fmt.Println("Recover ", perr)
			fmt.Println("Please check your type")
		}
	}()

	ivalue := reflect.Indirect(reflect.ValueOf(input))
	itype := ivalue.Type()

	var (
		strQuery  *bytes.Buffer
		cols      []string
		wheres    []string
		orderBy   []string
		primaries []string
		foreigns  []string
		colVal    map[string]interface{}
	)

	if ivalue.Kind() == reflect.Struct {

		/*for i := 0; i < ivalue.NumField(); i++ {
			ftype := itype.Field(i)
			fval := ivalue.Field(i)

			tags := ftype.Tag.Get("query")

			if tags != "" {
				stags := strings.Split(tags, ",")
				for idx, s := range stags {
					if idx == 0 {
						tmpCols = s
					} else if s == "col" {
						cols = append(cols, tmpCols)
						if ftype.Type.Kind() == reflect.Struct || ftype.Type.Kind() == reflect.Map || ftype.Type.Kind() == reflect.Slice {
							return "", errors.New("unsupported struct with field of type collection")
						} else {
							if !IsEmpty(reflect.Indirect(fval).Interface()) {
								if colVal == nil {
									colVal = make(map[string]interface{})
								}
								ifield := reflect.Indirect(fval)
								colVal[tmpCols]=ifield.Interface()
							}
						}
					} else if s == "where" {
						wheres = append(wheres, tmpCols)
					} else if s == "order" {
						orderBy = append(orderBy, tmpCols)
					} else if strings.Contains(s, "primary") {
						ssplit := strings.Split(s, "#")
						primaries = append(primaries, tmpCols+"#"+ssplit[1])
					} else if strings.Contains(s, "foreign") {
						ssplit := strings.Split(s, "#")
						foreigns = append(foreigns, tmpCols+"#"+ssplit[1])
					}
				}
			}
		}*/

		if err := structIterate(input, cols, wheres, orderBy, primaries, foreigns, colVal); err != nil {
			return "", err
		}

		// generate query
		strQuery = bytes.NewBufferString("SELECT ")
		if len(cols) > 0 {
			strQuery.WriteString(strings.Join(cols, ", "))
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

		if len(primaries) > 0 && len(foreigns) > 0 && len(primaries) == len(primaries) {
			if !strings.Contains(strQuery.String(), " WHERE ") {
				strQuery.WriteString(" WHERE ")
			}

			for _, p := range primaries {
				psplits := strings.Split(p, "#")
				for idx, f := range foreigns {
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

		if len(wheres) > 0 {
			if !strings.Contains(strQuery.String(), " WHERE ") {
				strQuery.WriteString(" WHERE ")
			}
			for idx, s := range wheres {
				globalIdx += idx
				sv, found := colVal[s]
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

		if len(orderBy) > 0 {
			strQuery.WriteString(" ORDER BY ")
			strQuery.WriteString(strings.Join(orderBy, ", "))
		}

		return strQuery.String(), nil

	} else {
		return "", errors.New("only accepted struct type")
	}

}

func ComposeUpdateQuery(input interface{}) (string, error) {
	return "", errors.New("unimplemented")
}

func ComposeCreateQuery(input interface{}) (string, error) {
	return "", errors.New("unimplemented")
}

func ComposeDeleteQuery(input interface{}) (string, error) {
	return "", errors.New("unimplemented")
}

func IsEmpty(v interface{}) bool {
	return reflect.DeepEqual(v, reflect.Zero(reflect.TypeOf(v)).Interface())
}

func structIterate(input interface{}, cols []string, wheres []string, orderBy []string, primaries []string, foreigns []string, colVal map[string]interface{}) error {
	ivalue := reflect.Indirect(reflect.ValueOf(input))
	itype := ivalue.Type()

	var tmpCols string

	if ivalue.Kind() == reflect.Struct {
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
						cols = append(cols, tmpCols)
						if ftype.Type.Kind() == reflect.Struct || ftype.Type.Kind() == reflect.Map || ftype.Type.Kind() == reflect.Slice {
							return errors.New("unsupported struct with field of type collection")
						} else {
							if !IsEmpty(reflect.Indirect(fval).Interface()) {
								if colVal == nil {
									colVal = make(map[string]interface{})
								}
								ifield := reflect.Indirect(fval)
								colVal[tmpCols] = ifield.Interface()
							}
						}
					} else if s == "where" {
						wheres = append(wheres, tmpCols)
					} else if s == "order" {
						orderBy = append(orderBy, tmpCols)
					} else if strings.Contains(s, "primary") {
						ssplit := strings.Split(s, "#")
						primaries = append(primaries, tmpCols+"#"+ssplit[1])
					} else if strings.Contains(s, "foreign") {
						ssplit := strings.Split(s, "#")
						foreigns = append(foreigns, tmpCols+"#"+ssplit[1])
					}
				}
			}
		}
		return nil
	} else {
		return errors.New("only accepted struct type")
	}
}

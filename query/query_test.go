package query

import (
	"reflect"
	"testing"
)

const (
	success = "\u2713"
	failed  = "\u2717"
)

type User struct {
	Id      int64  `query:"id,col,primary#default,where"`
	Name    string `query:"name,col"`
	Address string `query:"address,col"`
	Email   string `query:"email,col"`
	Photo   string `query:"photo,col"`
}

type User2 struct {
	Id      int64  `query:"id,primary#default,where"`
	Name    string `query:"name,col"`
	Address string `query:"address,col"`
	Email   string `query:"email,col"`
	Photo   string `query:"photo,col"`
}

type UserSelect struct {
	Id      int64  `query:"id,col,where"`
	Name    string `query:"name,where,col"`
	Address string `query:"address,col"`
	Email   string `query:"email,col"`
	Photo   string `query:"photo,col"`
}

type UserSelectWithOrder struct {
	Id      int64  `query:"id,col,where"`
	Name    string `query:"name,where,col,order#desc"`
	Address string `query:"address,col"`
	Email   string `query:"email,col"`
	Photo   string `query:"photo,col"`
}

func (u UserSelect) Table() string {
	return "users"
}

type UserAddress struct {
	UserId     int64  `query:"users.id,col,primary#default"`
	Name       string `query:"users.name,col"`
	Email      string `query:"users.email,col"`
	Photo      string `query:"users.photo,col"`
	AddressId  int64  `query:"addresses.id,col,foreign#default"`
	StreetName string `query:"addresses.street_name,col"`
	PostCode   string `query:"addresses.post_code,col,where"`
}

func (u UserAddress) Table() string {
	return "users,addresses"
}

func TestQueryComposerSelect(t *testing.T) {
	t.Log("Testing select query from a struct")
	{
		user := User{}

		query, err := ComposeQuery(1, user, true)

		if err != nil {
			t.Errorf("%s expected error nil, got %s", failed, err.Error())
		} else {
			t.Logf("%s expected error nil", success)
		}

		expectedQuery := `SELECT id, name, address, email, photo FROM users WHERE id = ?`

		if query == expectedQuery {
			t.Logf("%s expected query %s", success, expectedQuery)
		} else {
			t.Errorf("%s expected query %s, got %s", failed, expectedQuery, query)
		}
	}

	t.Log("Testing select query from a struct, table name from Table method")
	{
		user := UserSelect{
			Id:   1234,
			Name: "Bilal Muhammad",
		}
		query, err := ComposeQuery(1, user, true)

		if err != nil {
			t.Errorf("%s expected error nil, got %s", failed, err.Error())
		} else {
			t.Logf("%s expected error nil", success)
		}

		expectedQuery := `SELECT id, name, address, email, photo FROM users WHERE id = 1234 AND name = 'Bilal Muhammad'`

		if query == expectedQuery {
			t.Logf("%s expected query %s", success, expectedQuery)
		} else {
			t.Errorf("%s expected query %s, got %s", failed, expectedQuery, query)
		}
	}

	t.Log("Testing select query with multiple tables")
	{
		userAddress := UserAddress{}
		query, err := ComposeQuery(1, userAddress, true)

		if err != nil {
			t.Errorf("%s expected error nil, got %s", failed, err.Error())
		} else {
			t.Logf("%s expected error nil", success)
		}

		expectedQuery := `SELECT users.id, users.name, users.email, users.photo, addresses.id, addresses.street_name, addresses.post_code FROM users,addresses WHERE users.id = addresses.id AND addresses.post_code = ?`

		if query == expectedQuery {
			t.Logf("%s expected query [%s]", success, expectedQuery)
		} else {
			t.Errorf("%s expected query [%s], got [%s]", failed, expectedQuery, query)
		}
	}

	t.Log("Testing select query with order")
	{
		user := UserSelectWithOrder{}
		query, err := ComposeQuery(1, user, true)

		if err != nil {
			t.Errorf("%s expected error nil, got %s", failed, err.Error())
		} else {
			t.Logf("%s expected error nil", success)
		}

		expectedQuery := `SELECT id, name, address, email, photo FROM userselectwithorders WHERE id = ? AND name = ? ORDER BY name desc`
		if query == expectedQuery {
			t.Logf("%s expected query %s", success, expectedQuery)
		} else {
			t.Errorf("%s expected query %s, got %s", failed, expectedQuery, query)
		}
	}
}

func TestComposeUpdateQuery(t *testing.T) {
	t.Log("Testing ComposeUpdateQuery, empty struct")
	{
		user := User2{}

		query, err := ComposeQuery(QueryUpdate, user, true)

		if err != nil {
			t.Errorf("%s expected error nil, got %s", failed, err.Error())
		} else {
			t.Logf("%s expected error nil", success)
		}

		expectedQuery := `UPDATE user2s SET name = ?, address = ?, email = ?, photo = ? WHERE id = ?`

		if query == expectedQuery {
			t.Logf("%s expecged query %s", success, expectedQuery)
		} else {
			t.Errorf("%s expected query [%s], got [%s]", failed, expectedQuery, query)
		}
	}

	t.Log("Testing ComposeUpdateQuery, filled struct")
	{
		user := User2{
			Id: 23,
			Name: "Bilal Muhammad",
			Address: "Example street No. 4",
			Email: "bilal.muhammad@example.com",
			Photo: "https://s3.amazonws.com/abucket/aphoto.png",
		}

		query, err := ComposeQuery(QueryUpdate, user, true)

		if err != nil {
			t.Errorf("%s expected error nil, got %s", failed, err.Error())
		} else {
			t.Logf("%s expected error nil", success)
		}

		expectedQuery := `UPDATE user2s SET name = 'Bilal Muhammad', address = 'Example street No. 4', email = 'bilal.muhammad@example.com', photo = 'https://s3.amazonws.com/abucket/aphoto.png' WHERE id = 23`

		if query == expectedQuery {
			t.Logf("%s expected query [%s]", success, expectedQuery)
		} else {
			t.Errorf("%s expected query [%s], got [%s]", failed, expectedQuery, query)
		}
	}
}

func TestComposeInsertQuery(t *testing.T) {
	t.Log("Test ComposeInsertQuery empty struct")
	{
		query, err := ComposeQuery(QueryInsert, User{}, true)
		if err != nil {
			t.Errorf("%s expected error nil, got %s", failed, err.Error())
		} else {
			t.Logf("%s expected error nil", success)
		}

		expectedQuery := `INSERT INTO users ( id, name, address, email, photo ) VALUES ( ?, ?, ?, ?, ? )`

		if query == expectedQuery {
			t.Logf("%s expecged query [%s]", success, expectedQuery)
		} else {
			t.Errorf("%s expected query [%s], got [%s]", failed, expectedQuery, query)
		}
	}

	t.Log("Test ComposeInsertQuery partially filled struct")
	{
		user := User{}
		user.Id = 12345
		user.Name = "Testing Admin"
		user.Photo = "https://ec2.aws.com/432445"
		query, err := ComposeQuery(QueryInsert, user, true)
		if err != nil {
			t.Errorf("%s expected error nil, got %s", failed, err.Error())
		} else {
			t.Logf("%s expected error nil", success)
		}

		expecteQuery := `INSERT INTO users ( id, name, address, email, photo ) VALUES ( 12345, 'Testing Admin', ?, ?, 'https://ec2.aws.com/432445' )`
		if query == expecteQuery {
			t.Logf("%s expected query [%s]", success, expecteQuery)
		} else {
			t.Errorf("%s expected query [%s], got [%s]", failed, expecteQuery, query)
		}
	}
}

type DeleteTypeTest struct {
	Id int `query:"id,where"`
}

func (d DeleteTypeTest) Table() string {
	return "types"
}

func TestComposeDeleteQuery(t *testing.T) {
	t.Log("Test ComposeDeleteQuery empty struct")
	{
		user := User{}
		query, err := ComposeQuery(QueryDelete, user, true)
		if err != nil {
			t.Errorf("%s expected error nil, got %s", failed, err.Error())
		} else {
			t.Logf("%s expected error nil", success)
		}

		expectedQuery := `DELETE FROM users WHERE id = ?`
		if query == expectedQuery {
			t.Logf("%s expected query %s", success, expectedQuery)
		} else {
			t.Errorf("%s expected query %s, got %s", failed, expectedQuery, query)
		}
	}

	t.Log("Test ComposeDeleteQuery filled struct and defined Table method")
	{
		types := DeleteTypeTest{}
		types.Id = 34512
		query, err := ComposeQuery(QueryDelete, types, true)
		if err != nil {
			t.Errorf("%s expected error nil, got %s", failed, err.Error())
		} else {
			t.Logf("%s expected error nil", success)
		}

		expectedQuery := `DELETE FROM types WHERE id = 34512`
		if query == expectedQuery {
			t.Logf("%s expected query %s", success, expectedQuery)
		} else {
			t.Errorf("%s expected query %s, got %s", failed, expectedQuery, query)
		}
	}
}

func TestIsEmpty(t *testing.T) {
	t.Log("Testing the IsEmpty function")
	{
		i1 := 0

		result := reflect.DeepEqual(i1, reflect.Zero(reflect.TypeOf(i1)).Interface())

		t.Log("Result", result)
		t.Log("zero value", reflect.Zero(reflect.TypeOf(i1)).Interface())

		i2 := 0

		var ii interface{}
		ii = i2

		result = reflect.DeepEqual(i1, reflect.Zero(reflect.TypeOf(ii)).Interface())

		t.Log("Result", result)
		t.Log("zero value", reflect.Zero(reflect.TypeOf(i1)).Interface())
	}
}

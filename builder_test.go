package gosql

import (
	"time"
	"testing"
	"strconv"
	"fmt"
	"encoding/json"
)

var (
	createSchema = `
CREATE TABLE users (
	id int(11) unsigned NOT NULL AUTO_INCREMENT,
	name  varchar(100) NOT NULL DEFAULT '',
	email  varchar(100) NOT NULL DEFAULT '',
	created_at datetime NOT NULL,
	updated_at datetime NOT NULL,
  	PRIMARY KEY (id)
)ENGINE=InnoDB CHARSET=utf8;
`

	dropSchema = `
	drop table users
`
)

type Users struct {
	Id        int       `db:"id"`
	Name      string    `db:"name"`
	Email     string    `db:"email"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (u *Users) DbName() string {
	return "default"
}

func (u *Users) TableName() string {
	return "users"
}

func (u *Users) PK() string {
	return "id"
}

func RunWithSchema(t *testing.T, test func(t *testing.T)) {
	db := DB()
	defer func() {
		db.Exec(dropSchema)
	}()

	_, err := db.Exec(createSchema)

	if err != nil {
		t.Fatalf("create schema error:%s", err)
	}

	test(t)
}

func insert(id int) {
	user := &Users{
		Id:    id,
		Name:  "test" + strconv.Itoa(id),
		Email: "test" + strconv.Itoa(id) + "@test.com",
	}
	Model(user).Create()
}

func TestBuilder_Get(t *testing.T) {
	RunWithSchema(t, func(t *testing.T) {
		insert(1)
		user := &Users{}
		err := Model(user).Where("id = ?", 1).Get()

		if err != nil {
			t.Error(err)
		}
		//fmt.Println(user)
	})
}

func json_encode(i interface{}) string {
	ret, _ := json.Marshal(i)
	return string(ret)
}

func TestBuilder_All(t *testing.T) {
	RunWithSchema(t, func(t *testing.T) {
		insert(1)
		insert(2)

		user := make([]*Users, 0)
		err := Model(&user).All()

		if err != nil {
			t.Error(err)
		}

		fmt.Println(json_encode(user))
	})
}

func TestBuilder_Update(t *testing.T) {
	RunWithSchema(t, func(t *testing.T) {
		insert(1)

		user := &Users{
			Name: "test2",
		}

		_, err := Model(user).Where("id=1").Update()

		if err != nil {
			t.Error("update user error", err)
		}
	})
}

func TestBuilder_Delete(t *testing.T) {
	RunWithSchema(t, func(t *testing.T) {
		insert(1)
		_, err := Model(&Users{}).Where("id=1").Delete()

		if err != nil {
			t.Error("delete user error", err)
		}
	})
}

func TestBuilder_Count(t *testing.T) {
	RunWithSchema(t, func(t *testing.T) {
		insert(1)

		num, err := Model(&Users{}).Count()

		if err != nil {
			t.Error(err)
		}

		if num != 1 {
			t.Error("count user error")
		}
	})
}

func TestBuilder_Create(t *testing.T) {
	RunWithSchema(t, func(t *testing.T) {
		user := &Users{
			Id:    1,
			Name:  "test",
			Email: "test@test.com",
		}
		id, err := Model(user).Create()

		if err != nil {
			t.Error(err)
		}

		if id != 1 {
			t.Error("lastInsertId error", id)
		}
	})
}

func TestBuilder_Limit(t *testing.T) {
	RunWithSchema(t, func(t *testing.T) {
		insert(1)
		insert(2)
		insert(3)
		user := &Users{}
		err := Model(user).Limit(1).Get()

		if err != nil {
			t.Error(err)
		}
	})
}

func TestBuilder_Offset(t *testing.T) {
	RunWithSchema(t, func(t *testing.T) {
		insert(1)
		insert(2)
		insert(3)
		user := &Users{}
		err := Model(user).Limit(1).Offset(1).Get()

		if err != nil {
			t.Error(err)
		}
	})
}

func TestBuilder_OrderBy(t *testing.T) {
	RunWithSchema(t, func(t *testing.T) {
		insert(1)
		insert(2)
		insert(3)
		user := &Users{}
		err := Model(user).OrderBy("id desc").Limit(1).Offset(1).Get()

		if err != nil {
			t.Error(err)
		}

		if user.Id != 2 {
			t.Error("order by error")
		}

		//fmt.Println(user)
	})
}

func TestBuilder_Where(t *testing.T) {
	RunWithSchema(t, func(t *testing.T) {
		insert(1)
		insert(2)
		insert(3)
		user := make([]*Users,0)
		err := Model(&user).Where("id in(?,?)", 2, 3).OrderBy("id desc").All()

		if err != nil {
			t.Error(err)
		}

		if len(user) != 2 {
			t.Error("where error")
		}

		//fmt.Println(user)
	})
}
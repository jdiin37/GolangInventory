package user

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
)

type Getter interface {
	Get(offset, limit int) ItemList
	QueryByEmail(email string) Item
	QueryById(id int) (Item, error)
}

type Adder interface {
	Add(item Item) error
}

type Updater interface {
	Update(item Item) error
	UpdateNoAvatar(Item Item) error
}

type User struct {
	DB *sql.DB
}

type Item struct {
	ID         int
	Email      string `form:"email"`
	Password   string `form:"password"`
	Name       string `form:"name"`
	Avatar     sql.NullString
	CreateTime string
	UpdateTime string
}

type ItemList struct {
	List  []Item
	Count int
}

func (user *User) Get(offset, limit int) ItemList {
	items := []Item{}

	resp := ItemList{}
	stmt, _ := user.DB.Prepare("SELECT * FROM user ORDER BY ID LIMIT ?,?")
	rows, _ := stmt.Query(offset, limit)

	var id int
	var email string
	var password string

	for rows.Next() {
		rows.Scan(&id, &email, &password)
		item := Item{
			ID:       id,
			Email:    email,
			Password: password,
		}
		items = append(items, item)
	}

	resp.List = items
	user.DB.QueryRow(`
		SELECT COUNT(*) FROM user
	 `).Scan(&resp.Count)

	return resp
}

func (user *User) Add(item Item) error {
	stmt, _ := user.DB.Prepare(`
		INSERT INTO user (Email,Password,CreateTime) values (?,?,datetime('now','localtime'))
	`)

	result, err := stmt.Exec(item.Email, item.Password)
	fmt.Println(result, err)
	return err
}

func (user *User) Update(item Item) error {
	var stmt, e = user.DB.Prepare("UPDATE user SET password=?,avatar=?,name=?,UpdateTime = datetime('now','localtime') WHERE id=? ")
	if e != nil {
		log.Panicln("發生了錯誤", e.Error())
	}
	_, e = stmt.Exec(item.Password, item.Avatar.String, item.Name, item.ID)
	if e != nil {
		log.Panicln("錯誤 e", e.Error())
	}
	return e
}
func (user *User) UpdateNoAvatar(item Item) error {
	var stmt, e = user.DB.Prepare("UPDATE user SET password=?,name=?, UpdateTime = datetime('now','localtime') WHERE id=? ")
	if e != nil {
		log.Panicln("發生了錯誤", e.Error())
	}
	_, e = stmt.Exec(item.Password, item.Name, item.ID)
	if e != nil {
		log.Panicln("錯誤 e", e.Error())
	}
	return e
}

func (user *User) QueryByEmail(email string) Item {
	u := Item{}
	row := user.DB.QueryRow("SELECT ID,Email,Password,Name FROM user WHERE email = ?;", email)
	e := row.Scan(&u.ID, &u.Email, &u.Password, &u.Name)
	if e != nil {
		//log.Panicln(e)
	}
	return u
}

func (user *User) QueryById(id int) (Item, error) {
	u := Item{}
	row := user.DB.QueryRow("SELECT ID,Email,Password,Avatar,Name  FROM user WHERE id = ?;", id)
	e := row.Scan(&u.ID, &u.Email, &u.Password, &u.Avatar, &u.Name)
	if e != nil {
		//log.Panicln(e)
		return u, errors.New(e.Error())
	}
	return u, nil
}

func NewUser(db *sql.DB) *User {
	stmt, _ := db.Prepare(`
	CREATE TABLE IF NOT EXISTS "user" (
		"ID" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		"Email" TEXT UNIQUE,
		"Name" TEXT,
		"Password" TEXT,
		"Avatar" TEXT,
		"CreateTime" TEXT NOT NULL,
		"UpdateTime" TEXT
	);
`)
	stmt.Exec()
	return &User{
		DB: db,
	}
}

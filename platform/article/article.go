package article

import (
	"database/sql"
	"log"
)

type Getter interface {
	GetAll(offset, limit int) ItemList
	FindById(id int) (Item, error)
}

type Adder interface {
	Add(item Item) int
}

type Updater interface {
	Update(item Item) error
}

type Deleter interface {
	DeleteOne(id int)
}

type Article struct {
	DB *sql.DB
}

type Item struct {
	ID         int
	Type       string `form:"type"`
	Content    string `form:"content"`
	CreateTime string
	UpdateTime string
}

type ItemList struct {
	List  []Item
	Count int
}

func NewArticle(db *sql.DB) *Article {
	stmt, _ := db.Prepare(`
		CREATE TABLE IF NOT EXISTS "article" (
			"ID" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			"Type" varchar(20) NULL,
			"Content" TEXT NOT NULL,
			"CreateTime" TEXT NOT NULL,
			"UpdateTime" TEXT
		);
	`)
	stmt.Exec()
	return &Article{
		DB: db,
	}
}

func (article *Article) Add(item Item) int {
	stmt, _ := article.DB.Prepare(`
		INSERT INTO article (type,content,CreateTime) values (?,?,datetime('now','localtime'))
	`)
	result, err := stmt.Exec(item.Type, item.Content)
	if err != nil {
		log.Panicln("文章添加失败", err.Error())
	}
	i, _ := result.LastInsertId()
	return int(i)
}

func (article *Article) GetAll(offset, limit int) ItemList {
	items := []Item{}

	resp := ItemList{}
	stmt, _ := article.DB.Prepare("SELECT ID,Type,Content FROM article ORDER BY ID LIMIT ?,?")
	rows, _ := stmt.Query(offset, limit)

	for rows.Next() {
		item := Item{}
		rows.Scan(&item.ID, &item.Type, &item.Content)
		items = append(items, item)
	}

	resp.List = items
	article.DB.QueryRow(`
		SELECT COUNT(*) FROM article
	 `).Scan(&resp.Count)

	return resp
}

func (article *Article) FindById(id int) (Item, error) {
	u := Item{}
	row := article.DB.QueryRow("SELECT ID,Type,Content FROM article WHERE id = ?;", id)
	e := row.Scan(&u.ID, &u.Type, &u.Content)
	if e != nil {
		log.Println(e.Error())
		return u, e
	}
	return u, nil
}

func (article *Article) DeleteOne(id int) {
	if _, e := article.DB.Exec("DELETE FROM article WHERE id = ?;", id); e != nil {
		log.Panicln("数据发生错误，无法删除")
	}
}

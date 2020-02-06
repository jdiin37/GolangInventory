package inventory

import "database/sql"

type Getter interface {
	Get(offset, limit int) ItemList
}

type Adder interface {
	Add(item Item)
}

type Inventory struct {
	DB *sql.DB
}

type Item struct {
	ID    int
	Title string `json:"title"`
	Post  string `json:"post"`
}

type ItemList struct {
	List  []Item
	Count int
}

func (inventory *Inventory) Get(offset, limit int) ItemList {
	items := []Item{}

	resp := ItemList{}
	// rows, _ := inventory.DB.Query(`
	// 	SELECT * FROM inventory LIMIT (?,?)
	// `, offset, limit)

	stmt, _ := inventory.DB.Prepare("SELECT * FROM inventory ORDER BY ID LIMIT ?,?")
	rows, _ := stmt.Query(offset, limit)

	var id int
	var title string
	var post string

	for rows.Next() {
		rows.Scan(&id, &title, &post)
		item := Item{
			ID:    id,
			Title: title,
			Post:  post,
		}
		items = append(items, item)
	}

	resp.List = items
	inventory.DB.QueryRow(`
		SELECT COUNT(*) FROM inventory
	 `).Scan(&resp.Count)

	return resp
}

func (inventory *Inventory) Add(item Item) {
	stmt, _ := inventory.DB.Prepare(`
		INSERT INTO inventory (title,post) values (?,?)
	`)

	stmt.Exec(item.Title, item.Post)
}

func NewInventory(db *sql.DB) *Inventory {
	stmt, _ := db.Prepare(`
		CREATE TABLE IF NOT EXISTS "inventory" (
			"ID" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			"title" TEXT,
			"post" TEXT
		);
	`)
	stmt.Exec()
	return &Inventory{
		DB: db,
	}
}

package activity

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

type Item struct {
	ID             int
	StartTime      string `form:"startTime"`
	EndTime        string `form:"endTime"`
	County         string `form:"county"`
	Location       string `form:"location"`
	CourtType      string `form:"courtType"`
	Memo           string `form:"memo"`
	Context        string `form:"context"`
	CoverImage     sql.NullString
	CreateTime     string
	UpdateTime     string
	DeleteTime     string
	CreateUserID   int
	DeleteUserID   int
	CreateUserName string
	ApplyCount     int
}

type ItemList struct {
	List  []Item
	Count int
}

type Activity struct {
	DB *sql.DB
}

func NewActivity(db *sql.DB) *Activity {
	stmt, _ := db.Prepare(`
		CREATE TABLE IF NOT EXISTS "activity" (
			"ID" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			"StartTime" TEXT NOT NULL,
			"EndTime" TEXT NOT NULL,
			"County" TEXT NOT NULL,
			"Location" TEXT NOT NULL,
			"CourtType" TEXT NOT NULL,
			"Memo" TEXT NOT NULL,
			"Context" TEXT NOT NULL,
			"CoverImage" TEXT,
			"CreateTime" TEXT NOT NULL,
			"UpdateTime" TEXT,
			"DeleteTime" TEXT,
			"CreateUserID" INTEGER,
			"DeleteUserID" INTEGER
		);
	`)
	stmt.Exec()
	return &Activity{
		DB: db,
	}
}

func (activity *Activity) Add(item Item) int {
	stmt, _ := activity.DB.Prepare(`
		INSERT INTO activity (startTime,endTime,county,location,courtType,memo,context,coverImage,createUserID,CreateTime) 
		values (?,?,?,?,?,?,?,?,?,datetime('now','localtime'))
	`)
	result, err := stmt.Exec(item.StartTime,
		item.EndTime,
		item.County,
		item.Location,
		item.CourtType,
		item.Memo,
		item.Context,
		item.CoverImage,
		item.CreateUserID)
	if err != nil {
		log.Panicln("Add fail", err.Error())
	}
	i, _ := result.LastInsertId()
	return int(i)
}

func (activity *Activity) Update(item Item) error {
	var stmt, e = activity.DB.Prepare("UPDATE activity SET startTime=?,endTime=?,county=?,location=?,courtType=?,memo=?,context=?,coverImage=?,UpdateTime = datetime('now','localtime') WHERE id=? ")
	if e != nil {
		log.Panicln("Update error ", e.Error())
	}
	_, e = stmt.Exec(item.StartTime, item.EndTime, item.County, item.Location, item.CourtType, item.Memo, item.Context, item.CoverImage.String, item.ID)
	if e != nil {
		log.Panicln("Update error", e.Error())
	}
	return e
}

func (activity *Activity) GetAll(offset, limit int) ItemList {
	items := []Item{}

	resp := ItemList{}
	stmt, _ := activity.DB.Prepare(`SELECT a.ID,a.startTime,a.endTime,a.county,a.location,a.courtType,a.memo,a.context,a.coverImage,a.createUserID,b.name as createUserName 
	 FROM activity a,user b
	 WHERE a.deleteTime IS NULL AND a.createUserID = b.ID
	 ORDER BY a.ID DESC LIMIT ?,?`)
	rows, err := stmt.Query(offset, limit)

	if err == nil {
		for rows.Next() {
			item := Item{}
			rows.Scan(&item.ID, &item.StartTime, &item.EndTime, &item.County, &item.Location, &item.CourtType, &item.Memo, &item.Context, &item.CoverImage, &item.CreateUserID, &item.CreateUserName)
			activity.DB.QueryRow(`
				SELECT COUNT(*) FROM activity_apply WHERE ActivityID = ? 
			`, item.ID).Scan(&item.ApplyCount)
			items = append(items, item)
		}
	}

	resp.List = items
	activity.DB.QueryRow(`
		SELECT COUNT(*) FROM activity WHERE deleteTime IS NOT NULL
	 `).Scan(&resp.Count)

	return resp
}

func (activity *Activity) FindById(id int) (Item, error) {
	item := Item{}
	row := activity.DB.QueryRow(`SELECT a.ID,a.startTime,a.endTime,a.county,a.location,a.courtType,a.memo,a.context,a.coverImage,a.createUserID,b.name as createUserName FROM activity a,user b
	WHERE a.deleteTime IS NULL AND a.createUserID = b.ID AND a.ID = ?;`, id)
	e := row.Scan(&item.ID, &item.StartTime, &item.EndTime, &item.County, &item.Location, &item.CourtType, &item.Memo, &item.Context, &item.CoverImage, &item.CreateUserID, &item.CreateUserName)
	if e != nil {
		log.Println(e.Error())
		return item, e
	}
	return item, nil
}

func (activity *Activity) DeleteOne(id int) {
	if _, e := activity.DB.Exec("UPDATE activity SET DeleteTime = datetime('now','localtime') WHERE id = ?;", id); e != nil {
		log.Panicln("deleteOne error")
	}
}

package activityApply

import (
	"database/sql"
	"log"
)

type Action interface {
	UserApply(item Item) int
	GetApplyList(activityID int) ItemList
	CancelApply(applyID int) error
	CheckCanApply(item Item) bool
}

type Item struct {
	ID            int
	ActivityID    int    `form:"activityID"`
	UserID        int    `form:"userID"`
	Remark        string `form:"remark"`
	ApplyStatus   int
	CreateTime    string
	UpdateTime    string
	CheckInTime   string
	ApplyUserName string
}

type ItemList struct {
	List  []Item
	Count int
}

type ActivityApply struct {
	DB *sql.DB
}

func NewActivityApply(db *sql.DB) *ActivityApply {
	stmt, _ := db.Prepare(`
		CREATE TABLE IF NOT EXISTS "activity_apply" (
			"ID" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			"ActivityID" INTEGER NOT NULL,
			"UserID" INTEGER NOT NULL,
			"Remark" TEXT,
			"ApplyStatus" INTEGER NOT NULL,
			"CreateTime" TEXT NOT NULL,
			"UpdateTime" TEXT,
			"CheckInTime" TEXT
		);
	`)
	stmt.Exec()
	return &ActivityApply{
		DB: db,
	}
}
func (activityApply *ActivityApply) CheckCanApply(item Item) bool {
	var count int
	activityApply.DB.QueryRow(`
	SELECT COUNT(*) FROM activity_apply WHERE ActivityID = ? AND UserID = ? 
 `, item.ActivityID, item.UserID).Scan(&count)
	if count > 0 {
		return false
	}
	return true
}

func (activityApply *ActivityApply) UserApply(item Item) int {

	stmt, _ := activityApply.DB.Prepare(`
		INSERT INTO activity_apply (ActivityID,UserID,Remark,ApplyStatus,CreateTime) 
		values (?,?,?,?,datetime('now','localtime'))
	`)
	result, err := stmt.Exec(item.ActivityID,
		item.UserID,
		item.Remark,
		0)
	if err != nil {
		log.Panicln("UserApply fail", err.Error())
	}
	i, _ := result.LastInsertId()
	return int(i)
}

func (activityApply *ActivityApply) GetApplyList(activityID int) ItemList {
	items := []Item{}
	resp := ItemList{}
	stmt, _ := activityApply.DB.Prepare(`SELECT a.ID,a.ActivityID,a.UserID,a.Remark,a.ApplyStatus,a.CreateTime,b.name as ApplyUserName 
	 FROM activity_apply a,user b 
	 WHERE a.ActivityID = ? AND a.UserID = b.ID AND a.ApplyStatus != 99
	 ORDER BY a.ID DESC`)
	rows, err := stmt.Query(activityID)

	if err == nil {
		for rows.Next() {
			item := Item{}
			rows.Scan(&item.ID, &item.ActivityID, &item.UserID, &item.Remark, &item.ApplyStatus, &item.CreateTime, &item.ApplyUserName)
			items = append(items, item)
		}
	}

	resp.List = items
	activityApply.DB.QueryRow(`
		SELECT COUNT(*) FROM activity_apply WHERE ActivityID = ?
	 `, activityID).Scan(&resp.Count)

	return resp
}

func (activityApply *ActivityApply) CancelApply(applyID int) error {
	var stmt, e = activityApply.DB.Prepare("UPDATE activity_apply SET applyStatus=99,UpdateTime = datetime('now','localtime') WHERE id=? ")
	if e != nil {
		log.Panicln("CancelApply error ", e.Error())
	}
	_, e = stmt.Exec(applyID)
	if e != nil {
		log.Panicln("CancelApply error", e.Error())
	}
	return e
}

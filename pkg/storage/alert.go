package storage

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type AlertTable struct {
	gorm.Model
	Appid       uuid.UUID
	Level       string
	Type        string
	Info        string `sql:"type:text"`
	PostContact string
}

func CreateAlert(alert AlertTable, db *gorm.DB) (AlertTable, error) {
	db.Create(alert)
	return alert, nil
}

func UpdateAlert(alert AlertTable, db *gorm.DB) (AlertTable, error) {
	exist := db.First(&alert, "id = ?", alert.ID).RecordNotFound()
	if !exist {
		return alert, fmt.Errorf("can not find id %s", alert.ID)
	}
	db.Model(&alert).Where("id = ?", alert.ID).Updates(
		map[string]interface{}{
			"appid":        alert.Appid,
			"level":        alert.Level,
			"type":         alert.Type,
			"info":         alert.Info,
			"post_contact": alert.PostContact,
		})

	return alert, nil
}

//func SearchAlert(id int, appid string, start_time string, end_time string, alertType string){}

func ListAlert(from, count int, db *gorm.DB) ([]AlertTable, error) {
	var alerts []AlertTable
	exist := db.Limit(count).Where("id", from).Find(&alerts).RecordNotFound()
	if !exist {
		return alerts, fmt.Errorf("list app failed")
	}
	return alerts, nil}

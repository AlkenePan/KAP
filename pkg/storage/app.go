package storage


import (
	"fmt"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"youzoo/why/pkg/app"
)

type AppTable struct {
	gorm.Model
	Appid uuid.UUID
}

type SourceTable struct {
	gorm.Model
	Appid uuid.UUID
	Language string
}

type ExecutableTable struct {
	gorm.Model
	Appid uuid.UUID
	AbsPath string
}

func CreateApp(app app.App, db *gorm.DB) (error) {
	db.Create(&AppTable{Appid:app.Appid})
	db.Create(&SourceTable{Appid:app.Appid, Language:app.SourceInfo.Language})
	db.Create(&ExecutableTable{Appid:app.Appid, AbsPath:app.ExecInfo.AbsPath})
	return nil
}

func UpdateApp(app app.App, db *gorm.DB) (error) {
	var appTable AppTable
	var sourceTable SourceTable
	var executableTable ExecutableTable
	exist := db.First(&appTable, "appid = ?", app.Appid).RecordNotFound()
	if !exist {
		return fmt.Errorf("can not find appid %s", app.Appid)
	}
	db.Model(&sourceTable).Where("appid = ?", app.Appid).Update("Language", app.SourceInfo.Language)
	db.Model(&executableTable).Where("appid = ?", app.Appid).Update("AbsPath", app.ExecInfo.AbsPath)
	return nil
}

func FindApp(appid string, db *gorm.DB) (app.App, error) {
	var app app.App
	exist := db.Where("appid = ?", appid).First(&app).RecordNotFound()
	if !exist {
		return app, fmt.Errorf("can not find appid %s", app.Appid)
	}
	return app, nil
}


func ListApp(from , count int, db *gorm.DB) ([]app.App, error) {
	var apps []app.App
	exist := db.Limit(count).Where("id", from).Find(&apps).RecordNotFound()
	if !exist {
		return apps, fmt.Errorf("list app failed")
	}
	return apps, nil
}
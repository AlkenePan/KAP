package storage

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"youzoo/why/pkg/app"
)

type CryptoTable struct {
	gorm.Model
	Appid uuid.UUID
	PubKey string `sql:"type:text"`
	PrivateKey string `sql:"type:text"`
}


func NewKeyPair(app app.App, db *gorm.DB) (error) {
	db.Create(&ExecutableTable{Appid:app.Appid, AbsPath:app.ExecInfo.AbsPath})
	return nil
}

func UpdateKeyPair(pubkey string, prikey string, db *gorm.DB) (error) {
	var cryptoTable CryptoTable
	exist := db.First(&appTable, "appid = ?", app.Appid).RecordNotFound()
	if !exist {
		return fmt.Errorf("can not find appid %s", app.Appid)
	}
	db.Model(&sourceTable).Update("Language", app.SourceInfo.Language)
	db.Model(&executableTable).Update("AbsPath", app.ExecInfo.AbsPath)
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
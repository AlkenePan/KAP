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
	DNS   string
}

type SourceTable struct {
	gorm.Model
	Appid    uuid.UUID
	Language string
}

type ExecutableTable struct {
	gorm.Model
	Appid    uuid.UUID
	AbsPath  string
	Argv     string
	Envv     string
	Ptrace   bool
	UserName string
}

func CreateApp(app app.App, db *gorm.DB) (error) {
	db.Create(&AppTable{Appid: app.Appid, DNS: app.DNS})
	db.Create(&SourceTable{Appid: app.Appid, Language: app.SourceInfo.Language})
	db.Create(&ExecutableTable{Appid: app.Appid, AbsPath: app.ExecInfo.AbsPath,
		Argv: app.ExecInfo.Argv, Envv: app.ExecInfo.Envv,
		Ptrace: app.ExecInfo.Ptrace, UserName: app.ExecInfo.UserName})
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
	var appTable AppTable
	var sourceTable SourceTable
	var executableTable ExecutableTable

	var appInfo app.App
	exist := db.Where("appid = ?", appid).First(&appTable).RecordNotFound()
	if exist {
		return appInfo, fmt.Errorf("can not find appid %s", appid)
	}
	db.Where("appid = ?", appid).First(&sourceTable).RecordNotFound()
	db.Where("appid = ?", appid).First(&executableTable).RecordNotFound()
	appInfo.Appid = appTable.Appid
	appInfo.DNS = appTable.DNS
	appInfo.ExecInfo.AbsPath = executableTable.AbsPath
	appInfo.ExecInfo.Argv = executableTable.Argv
	appInfo.ExecInfo.Envv = executableTable.Envv
	appInfo.ExecInfo.Ptrace = executableTable.Ptrace
	appInfo.ExecInfo.UserName = executableTable.UserName
	appInfo.SourceInfo.Language = sourceTable.Language
	return appInfo, nil
}

func ListApp(from, count int, db *gorm.DB) ([]app.App, error) {
	var apps []AppTable
	var appInfos []app.App

	exist := db.Limit(count).Where("id", from).Find(&apps).RecordNotFound()
	if exist {
		return appInfos, fmt.Errorf("list app failed")
	}
	for _, app_ := range apps {
		var sourceTable SourceTable
		var executableTable ExecutableTable

		var appInfo app.App
		db.Where("appid = ?", app_.Appid).First(&sourceTable).RecordNotFound()
		db.Where("appid = ?", app_.Appid).First(&executableTable).RecordNotFound()
		appInfo.Appid = app_.Appid
		appInfo.DNS = app_.DNS
		appInfo.ExecInfo.AbsPath = executableTable.AbsPath
		appInfo.ExecInfo.Argv = executableTable.Argv
		appInfo.ExecInfo.Envv = executableTable.Envv
		appInfo.ExecInfo.Ptrace = executableTable.Ptrace
		appInfo.ExecInfo.UserName = executableTable.UserName
		appInfo.SourceInfo.Language = sourceTable.Language
		appInfos = append(appInfos, appInfo)
	}
	return appInfos, nil
}

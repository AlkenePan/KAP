package storage

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

type BuildStatusJson struct {
	ID     int
	Status string
}

type BuildTable struct {
	gorm.Model
	Appid  uuid.UUID
	Status string
}

func CreateBuild(buildTable BuildTable, db *gorm.DB) (BuildTable, error) {
	db.Create(buildTable)
	return buildTable, nil
}

func SetBuildStatus(buildStatusJson BuildStatusJson, db *gorm.DB) (BuildTable, error) {
	var buildTable BuildTable
	exist := db.First(&buildTable, "id = ?", buildStatusJson.ID).RecordNotFound()
	if !exist {
		return buildTable, fmt.Errorf("can not find id %s", buildStatusJson.ID)
	}
	db.Model(&buildTable).Where("id = ?", buildStatusJson.ID).Update("Status", buildStatusJson.Status)
	return buildTable, nil
}

func GetBuildStatus(buildID int, db *gorm.DB) (BuildTable, error) {
	var buildTable BuildTable
	exist := db.First(&buildTable, "id = ?", buildID).RecordNotFound()
	if !exist {
		return buildTable, fmt.Errorf("can not find id %s", buildID)
	}
	return buildTable, nil

}

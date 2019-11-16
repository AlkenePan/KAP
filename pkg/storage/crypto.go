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
	PriKey string `sql:"type:text"`
}

//func save_key(appid) {
//
//}
//func find_key(appid ) {
//
//}

func NewKeyPair(app app.App, pub string, pri string, db *gorm.DB) (error) {
	db.Create(&CryptoTable{Appid: app.Appid, PubKey: pub, PriKey: pri})
	return nil
}

func UpdateKeyPair(appid string, pubkey string, prikey string, db *gorm.DB) (error) {
	var cryptoTable CryptoTable
	exist := db.First(&cryptoTable, "Appid = ?", appid).RecordNotFound()
	if !exist {
		return fmt.Errorf("can not find appid %s", appid)
	}
	db.Model(&cryptoTable).Where("appid = ?", appid).Update("PubKey", pubkey)
	db.Model(&cryptoTable).Where("appid = ?", appid).Update("PriKey", prikey)
	return nil
}

func FindKeyPair(appid string, db *gorm.DB) (CryptoTable, error) {
	var cryptoTable CryptoTable
	exist := db.Where("appid = ?", appid).First(&cryptoTable).RecordNotFound()
	if exist {
		return cryptoTable, fmt.Errorf("can not find appid %s", appid)
	}
	return cryptoTable, nil
}

func FindPubKey(appid string, db *gorm.DB) (CryptoTable, error) {
	cryptoTable, err := FindKeyPair(appid, db)
	if err != nil {
		return CryptoTable{}, fmt.Errorf("can not find appid %s", appid)
	}
	cryptoTable.PriKey = ""
	return cryptoTable, nil
}

func FindPriKey(appid string, db *gorm.DB) (CryptoTable, error) {
	cryptoTable, err := FindKeyPair(appid, db)
	if err != nil {
		return CryptoTable{}, fmt.Errorf("can not find appid %s", appid)
	}
	cryptoTable.PubKey = ""
	return cryptoTable, nil
}

package models

import (
	"github.com/CalvoM/account-management-ussd/db"
	log "github.com/sirupsen/logrus"
)

//Vault holds a single vault object
type Vault struct {
	ID      int64
	UserID  int64
	Content string
}

//AddToVault Add a vault object to DB
func (v *Vault) AddToVault() (id int64, err error) {
	stmt, err := db.DbClient.Prepare("insert into vault(user_id,content) values($1,$2) returning vault_id;")
	if err != nil {
		log.Error(err)
		return
	}
	defer stmt.Close()
	err = stmt.QueryRow(v.UserID, v.Content).Scan(&id)
	if err != nil {
		log.Error(err)
		return
	}
	v.ID = id
	return
}

func (v *Vault) GetVaultByUserID(id int64) (vaults []Vault, err error) {
	stmt, err := db.DbClient.Prepare("select * from vault where user_id= $1;")
	if err != nil {
		log.Error(err)
		return
	}
	defer stmt.Close()
	rows, err := stmt.Query(id)
	for rows.Next() {
		v := Vault{}
		if err = rows.Scan(&v.ID, &v.UserID, &v.Content); err != nil {
			log.Error(err)
			return
		}
		vaults = append(vaults, v)
	}
	return
}

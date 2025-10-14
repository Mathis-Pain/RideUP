package getdata

import (
	"RideUP/models"
	"database/sql"
)

func GetUserFromLogin(db *sql.DB, login string) (models.User, error) {
	// Préparation de la requête SQL : récupérer id, username et password
	sql := `SELECT id, username, password FROM users WHERE username =?`
	row := db.QueryRow(sql, login)
	var user models.User
	// Parcour la base de données en cherchant le username correspondant
	err := row.Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

package dbHelper

import (
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"jayant/database"
	"jayant/models"
	"jayant/utils"
)

func CreateUser(body models.Register, password string) error {
	SQL := `insert into 
					users(name,email,password)
			values 
				($1, TRIM(LOWER($2)), $3)`
	_, err := database.Todo.Exec(SQL, body.Name, body.Email, password)
	if err != nil {
		logrus.Errorf("CreateUser : failed to creating user: %v", err)
		return err
	}
	return nil
}

//IsEmailExits function returns bool
func IsEmailExits(email string) (bool, error) {
	SQL := `SELECT 
    			id 
			FROM 
			    users 
			WHERE 	
			    email = $1 
			  	AND archived_at IS NULL `

	var id int
	err := database.Todo.Get(&id, SQL, email)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		logrus.Errorf("EmailExits : failed to check email exits: %v", err)
		return false, err
	}
	return true, nil
}

func IsSameEmailUsedInOtherUser(email string, id int) (bool, error) {
	SQL := `SELECT 
    			id
			FROM 
			    users
			WHERE 
			    email = $1
			  	AND id != $2
			  	AND archived_at IS NULL `
	var uid int
	err := database.Todo.Get(&uid, SQL, email, id)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		logrus.Errorf("EmailExits : error while check email exits: %v", err)
		return false, err
	}
	return true, nil
}

// UserInfo function returns user information
func UserInfo(id int) (models.RetrieveUserInfo, error) {
	var body models.RetrieveUserInfo
	SQL := `SELECT id,
       				name,
       				email,
       				created_at,
       				updated_at 
			FROM 
			    users
			WHERE	
			    id = $1 
			  	AND archived_at IS NULL `
	err := database.Todo.Get(&body, SQL, id)
	if err != nil {
		logrus.Errorf("InfoUser : failed to retrieve user info: %v", err)
		return body, err
	}
	return body, nil
}

func DeleteUser(db sqlx.Tx, id int) error {
	SQL := `UPDATE 
    				users 
			SET 
			    	archived_at = now()
			WHERE   
			    id = $1 
			  	AND archived_at IS NULL `
	_, err := db.Exec(SQL, id)
	if err != nil {
		logrus.Errorf("DeleteUser : failed to deleting user: %v", err)
		return err
	}
	return nil
}

func DeleteSession(db sqlx.Tx, userid int) error {
	SQL := `UPDATE
					user_session
			SET
			    	archived_at = now()
			WHERE
			    	user_id = $1
			    	AND archived_at > now()`
	_, err := db.Exec(SQL, userid)
	if err != nil {
		logrus.Errorf("Session_User : failed to deleting session user: %v", err)
		return err
	}
	return nil

}

func Logout(token string) error {
	SQL := `UPDATE
					user_session
			SET
			    	archived_at = now()
			WHERE
			    	session_token =$1`
	_, err := database.Todo.Exec(SQL, token)
	if err != nil {
		logrus.Errorf("Session_User : failed to deleting session user: %v", err)
		return err
	}
	return nil
}

func UpdateUser(db sqlx.Tx, userid int, body models.Register) error {

	SQL := `UPDATE 
    			users
			SET 	
			    name = $2,
			    email = $3,
			    updated_at = now()
			WHERE 	
			    id = $1 
			  	AND archived_at IS NULL `
	_, err := db.Exec(SQL, userid, body.Name, body.Email)
	if err != nil {
		logrus.Errorf("UpdateUser : failed to updating user data: %v", err)
		return err
	}
	return nil
}

func CreateSession(id int, SessionToken string) error {
	SQL := `insert into 
    					user_session(user_Id,session_token)
			values 
			    	($1,$2)`
	_, err := database.Todo.Exec(SQL, id, SessionToken)
	if err != nil {
		logrus.Errorf("CreateUserSession : failed to Creating session: %v", err)
		return err
	}
	return nil
}

func RetrieveUserInfo(body models.Login) (int, error) {
	s := `SELECT 
    			id  as user_id,
				password
		  FROM 
				users
		  WHERE 
				archived_at IS NULL 
				AND email = TRIM(LOWER($1))`
	userDetails := models.UserPass{}
	err := database.Todo.Get(&userDetails, s, body.Email)
	if err == sql.ErrNoRows {
		err = errors.New("failed to login")
		return 0, err
	}
	passwordErr := utils.CheckPassword(body.Password, userDetails.Password)
	if passwordErr != nil {
		err = errors.New("failed to login")
		return 0, err
	}

	return userDetails.UserID, nil
}

func GetUserBySession(sessionToken string) (*models.User, error) {
	SQL := `SELECT 
       			u.id, 
       			u.name, 
       			u.email, 
       			u.created_at 
			FROM 
			    users u
			JOIN 
			    user_session us on u.id = us.user_id
			WHERE 	
			    u.archived_at IS NULL 
			  	AND us.session_token = $1
			  	AND us.archived_at > now()`
	var user models.User
	err := database.Todo.Get(&user, SQL, sessionToken)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

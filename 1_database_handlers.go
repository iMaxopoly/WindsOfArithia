package main

import (
	"fmt"

	//3rd Party
	"github.com/go-xorm/xorm"
	"strings"
)

/***
 *      _____        _        _                      _    _                 _ _
 *     |  __ \      | |      | |                    | |  | |               | | |
 *     | |  | | __ _| |_ __ _| |__   __ _ ___  ___  | |__| | __ _ _ __   __| | | ___ _ __ ___
 *     | |  | |/ _` | __/ _` | '_ \ / _` / __|/ _ \ |  __  |/ _` | '_ \ / _` | |/ _ \ '__/ __|
 *     | |__| | (_| | || (_| | |_) | (_| \__ \  __/ | |  | | (_| | | | | (_| | |  __/ |  \__ \
 *     |_____/ \__,_|\__\__,_|_.__/ \__,_|___/\___| |_|  |_|\__,_|_| |_|\__,_|_|\___|_|  |___/
 *
 *
 */

var database__db_guru *xorm.Engine

type database__user_model struct {
	Col_nUserNo       int     `xorm:"int pk autoincr 'nEMID'"`
	Col_sUserID       string  `xorm:"nvarchar(50) notnull unique 'sUsername'" json:"username"`
	Col_sUserPW       string  `xorm:"nvarchar(100) notnull 'sUserPass'" json:"password"`
	Col_sUserPassSalt string  `xorm:"nvarchar(50) notnull 'sUserPassSalt'" json:"-"`
	Col_sUserIP       string  `xorm:"nvarchar(30) 'sIP'" json:"-"`
	Col_sEmail        string  `xorm:"nvarchar(100) unique 'sEmail'" json:"email"`
	//Col_sGoogle           string  `gorm:"column:s_google" sql:"type:varchar(1024);unique" json:"googleid,omitempty"`
	//Col_sFacebook         string  `gorm:"column:s_facebook" sql:"type:varchar(1024);unique" json:"facebookid,omitempty"`
	//Col_sTwitter          string  `gorm:"column:s_twitter" sql:"type:varchar(1024);unique" json:"twitterid,omitempty"`
	//Col_sVerified     string  `gorm:"column:s_Verified" sql:"type:nvarchar(50);default:false" json:"verified,omitempty"`
	Col_sLastSeen     string  `xorm:"nvarchar(50) 's_LastSeen'" json:"lastseen,omitempty"`
	Col_sLastToken    string  `xorm:"nvarchar(1024) 's_LastToken'" json:"lastToken,omitempty"`
	Col_sSpecialCoins string  `xorm:"nvarchar(512) 's_SpecialCoins'" json:"specialcoins,omitempty"`
	RecaptchaResponse string  `xorm:"-" json:"recaptcharesponse"`
}

type database__token_model struct {
	Col_nUserNo int     `xorm:"int unique 'nEMID'"`
	Col_sToken  string  `xorm:"nvarchar(50) notnull unique 'sToken'" json:"-"`
}

func (c database__user_model) TableName() string {
	return "tAccounts"
}

func (c database__token_model) TableName() string {
	return "tTokens"
}

func database__update_verified(email string, username string, userno int) (error) {
	if _, err := database__db_guru.
	Exec("update tAccounts set sUsername = ? where sEmail = ? AND sUsername = ? AND nEMID = ?",
		strings.TrimPrefix(username, "unverified_"), email, username, userno); err != nil {
		return err
	}
	return nil
}

func database__update_daily(userNo int, email string, ip string, token string, time string) (error) {
	user := database__user_model{}
	user.Col_sUserIP = ip
	user.Col_sLastToken = token
	user.Col_sLastSeen = time
	if _, err := database__db_guru.Where("nEMID = ? AND sEmail = ?", userNo, email).Update(&user); err != nil {
		return err
	}
	return nil
}

func database__get_password_and_nEMID(username string) (string, int) {
	user := database__user_model{}
	_, err := database__db_guru.Where("sUsername = ?", username).Get(&user)
	if err != nil {
		return "", 0
	}
	return user.Col_sUserPW, user.Col_nUserNo
}

func database__save(username string, password string, email string, ip string, token string, time string) (error) {
	user := database__user_model{
		Col_sUserID : "unverified_" + username,
		Col_sUserPW : util__aes_encrypt([]byte(server_config.aes_key), password),
		Col_sUserPassSalt : "Arithia",
		Col_sEmail : email,
		Col_sUserIP : ip,
		Col_sLastToken : token,
		Col_sLastSeen : time,
	}
	if _, err := database__db_guru.Insert(user); err != nil {
		return err
	}
	return nil
}

//func database__save_oauth(email string, ip string, token string, time string) (err error) {
//	user := database__user_model{
//		Col_sEmail: email,
//		Col_sUserIP: ip,
//		Col_sLastToken: token,
//		Col_sLastSeen: time,
//		Col_sVerified: "true",
//	}
//	database__db_guru.NewRecord(user)
//	if err := database__db_guru.Create(&user); err != nil {
//		return err.Error
//	}
//	return nil
//}

func database__compare_login(email string) (database__user_model, bool, error) {
	user := database__user_model{}
	has, err := database__db_guru.Where("sEmail = ?", email).Get(&user)
	if err != nil {
		return database__user_model{}, false, err
	}
	return user, has, nil
}

func database__is_existing_userid(userid string) (bool, error) {
	user := database__user_model{}
	has, err := database__db_guru.Where("sUsername = ?", userid).Get(&user)
	if err != nil {
		return false, err
	}
	return has, nil
}

func database__is_currently_unverified(email string, userid string) (bool, error) {
	user := database__user_model{}
	fmt.Println("database__is_currently_unverified", email, userid)
	has, err := database__db_guru.Where("sEmail = ?", email).And("sUsername = ?", "unverified_" + userid).Get(&user)
	if err != nil {
		return false, err
	}
	return has, nil
}

func database__is_existing_email(email string) (bool, error) {
	user := database__user_model{}
	has, err := database__db_guru.Where("sEmail = ?", email).Get(&user)
	if err != nil {
		return false, err
	}
	return has, nil
}

func database__insert_game_token(nEMID int, gamekey string) {
	user := database__token_model{Col_nUserNo: nEMID, Col_sToken: gamekey}
	_, err := database__db_guru.Insert(&user)
	if err != nil {
		util__debug_info(fmt.Sprintf("[DB]database__insert_game_token| %s", err))
	}
}

func database__delete_game_token(nEMID int) {
	user := database__token_model{Col_nUserNo: nEMID}
	_, err := database__db_guru.Delete(&user)
	if err != nil {
		util__debug_info(fmt.Sprintf("[DB]database__delete_game_token| %s", err))
	}
}

func database__restservice_verify_token(token string) (int, string, bool) {
	tokenDb := database__token_model{}
	accountDb := database__user_model{}
	//[tAccounts].[nEMID] = [tTokens].[nEMID] and [tTokens].[sToken] = ?"
	database__db_guru.Where("sToken = ?", token).Get(&tokenDb)
	database__db_guru.Where("nEMID = ?", tokenDb.Col_nUserNo).Get(&accountDb)
	if accountDb.Col_sEmail == "" {
		return 0, "", true
	}
	return accountDb.Col_nUserNo, accountDb.Col_sUserID, false
}
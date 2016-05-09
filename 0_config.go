package main

import (
	"fmt"
	"os"
	"bufio"
	"io/ioutil"

	// 3rd Party
	"github.com/fogcreek/mini"
	_"github.com/denisenkom/go-mssqldb"
	"github.com/go-xorm/xorm"
)

/***
 *       _____             __ _
 *      / ____|           / _(_)
 *     | |     ___  _ __ | |_ _  __ _
 *     | |    / _ \| '_ \|  _| |/ _` |
 *     | |___| (_) | | | | | | | (_| |
 *      \_____\___/|_| |_|_| |_|\__, |
 *                               __/ |
 *                              |___/
 */

// WoA Specific Launcher related Constants
const conf__launcher_useragent string = "g4gs4yjhs4u54a#@^gdsh4G4"
const conf__launcher_err_loginfail string = "5wbr473543^#^4bv"
const conf__launcher_gamekey_salt string = "g53wh8ns"

var server_config struct {
	// Website configs
	port                     string `json:"port"`
	debug                    string `json:"debug"`
	domain                   string `json:"domain"`
	project_name             string `json:"project_name"`

	// Database
	db_instance_name         string `json:"db_instance_name"`
	db_username              string `json:"db_username"`
	db_password              string `json:"db_password"`
	db_port                  string `json:"db_port"`
	db_database              string `json:"db_database"`
	aes_key                  string `json:"aes_key"`

	// Sendgrid API key
	sendgrid_api_key         string `json:"sendgrid_api_key"`

	// XenAPI key
	xenapi_key               string `json:"xenapi_key"`

	// Facebook oauth
	facebook_oauth_client_id string `json:"facebook_oauth_client_id"`
	facebook_oauth_secret    string `json:"facebook_oauth_secret"`

	// Google oauth
	google_oauth_client_id   string `json:"google_oauth_client_id"`
	google_oauth_secret      string `json:"google_oauth_secret"`

	// Recaptcha API
	google_recaptcha_secret  string `json:"google_recaptcha_secret"`

	// JWT related
	jwt_privatekey_path      string `json:"jwt_privatekey_path"`
	jwt_privatekey           []byte
	jwt_publickey_path       string `json:"jwt_publickey_path"`
	jwt_publickey            []byte

	// Disposable block list
	disposable_emails_path   string `json:"disposable_emails_path"`
	disposable_emails        []string

	// File Security
	files_digest             string `json:"files_digest"`
}

func load__project_config() {
	var configFile, err = mini.LoadConfiguration("project_settings.ini")
	util__critical_error_check(err)

	// Website configs
	server_config.port =
	configFile.StringFromSection("Website", "port", "1234")
	server_config.debug =
	configFile.StringFromSection("Website", "debug", "false")
	server_config.domain =
	configFile.StringFromSection("Website", "domain", "windsofarithia")
	server_config.project_name =
	configFile.StringFromSection("Website", "project_name", "Winds of Arithia")

	// Database
	server_config.db_instance_name =
	configFile.StringFromSection("Database", "db_instance_name", "sqlserver")
	server_config.db_username =
	configFile.StringFromSection("Database", "db_username", "sa")
	server_config.db_password =
	configFile.StringFromSection("Database", "db_password", "12345")
	server_config.db_port =
	configFile.StringFromSection("Database", "db_port", "1337")
	server_config.db_database =
	configFile.StringFromSection("Database", "db_database", "Account")
	server_config.aes_key =
	configFile.StringFromSection("Database", "aes_key", "k19\"!3zeV75mvSEJS4M'8,?~v|xYpg|)")

	// Sendgrid API key
	server_config.sendgrid_api_key =
	configFile.StringFromSection("Mail", "sendgrid_api_key", "")

	// XenAPI key
	server_config.xenapi_key =
	configFile.StringFromSection("Xenforo", "xenapi_key", "")

	// Facebook oauth
	server_config.facebook_oauth_client_id =
	configFile.StringFromSection("Oauth", "facebook_oauth_client_id", "")
	server_config.facebook_oauth_secret =
	configFile.StringFromSection("Oauth", "facebook_oauth_secret", "")

	// Google oauth
	server_config.google_oauth_client_id =
	configFile.StringFromSection("Oauth", "google_oauth_client_id", "")
	server_config.google_oauth_secret =
	configFile.StringFromSection("Oauth", "google_oauth_secret", "")

	// Recaptcha API
	server_config.google_recaptcha_secret =
	configFile.StringFromSection("Captcha", "google_recaptcha_secret", "")

	// JWT related
	server_config.jwt_privatekey_path =
	configFile.StringFromSection("Json Web Tokens", "jwt_privatekey_path", "")
	server_config.jwt_publickey_path =
	configFile.StringFromSection("Json Web Tokens", "jwt_publickey_path", "")

	// Disposable block list
	server_config.disposable_emails_path =
	configFile.StringFromSection("Disposable Email", "disposable_emails_path", "")

	// File Security
	// Disposable block list
	server_config.files_digest =
	configFile.StringFromSection("Files Digest", "files_digest", "")
	util__debug_info("[CONFIG]load__project_config| Finished Loading Project Config")
}

func load__disposable_email_and_jwt_lists() {
	file, err := os.Open(server_config.disposable_emails_path)
	if err != nil {
		util__debug_info(fmt.Sprintf("[CONFIG]| %s", err))
		return
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		server_config.disposable_emails = append(server_config.disposable_emails, scanner.Text())
	}
	util__debug_info("[CONFIG]load__disposable_email_and_jwt_lists| Finished Disposable Email List")

	server_config.jwt_publickey, _ = ioutil.ReadFile(server_config.jwt_publickey_path)
	server_config.jwt_privatekey, _ = ioutil.ReadFile(server_config.jwt_privatekey_path)
	util__debug_info("[CONFIG]load__disposable_email_and_jwt_lists| Finished Loading JWT Files")
}

func load__database_init() {
	var err error
	database__db_guru, err = xorm.NewEngine("mssql",
		"Server=" + server_config.db_instance_name +
		";User Id=" + server_config.db_username +
		";Password=" + server_config.db_password +
		";Port=" + server_config.db_port +
		";Database=" + server_config.db_database)
	util__critical_error_check(err)

	if server_config.debug == "true" {
		database__db_guru.ShowSQL(true)
	}
	util__debug_info("[CONFIG]load__database_init| Finished Initializing Database Connection.")
}
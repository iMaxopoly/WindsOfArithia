package main

import (
	"fmt"
	"html/template"
	"bytes"
	"encoding/json"

	// 3rd Party
	"github.com/sendgrid/sendgrid-go"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"io/ioutil"
	"strings"
	"github.com/kataras/iris"
)

/***
 *      ______                 _ _   _    _                 _ _
 *     |  ____|               (_) | | |  | |               | | |
 *     | |__   _ __ ___   __ _ _| | | |__| | __ _ _ __   __| | | ___ _ __ ___
 *     |  __| | '_ ` _ \ / _` | | | |  __  |/ _` | '_ \ / _` | |/ _ \ '__/ __|
 *     | |____| | | | | | (_| | | | | |  | | (_| | | | | (_| | |  __/ |  \__ \
 *     |______|_| |_| |_|\__,_|_|_| |_|  |_|\__,_|_| |_|\__,_|_|\___|_|  |___/
 *
 *
 */

type EmailContent struct {
	Title     string
	SubTitle  string
	Body      string
	VerifyUrl string
}

func api_route__verify_email(c *iris.Context) {
	client_ip := string(c.Request.Header.Peek("X-Forwarded-For"))
	gtoken := c.URLParam("token")
	token, err := jwt.Parse(gtoken, func(token *jwt.Token) (interface{}, error) {
		return server_config.jwt_publickey, nil
	})

	if err != nil || !token.Valid {
		util__debug_info(fmt.Sprintf("[EMAIL]api_route__verify_email| %s - %s OR Token Invalid", client_ip, err))
		c.JSON(401, json_post{"message": "Authentication failed, unable to verify."})
		return
	}

	email := token.Claims["sub"].(string)

	s_user, s_exists, err := database__compare_login(email)
	if err != nil {
		util__debug_info(fmt.Sprintf("[EMAIL]api_route__verify_email| %s - %s", client_ip, err))
		c.JSON(401, json_post{"message": "We encountered an error."})
		return
	}
	if !s_exists {
		util__debug_info(fmt.Sprintf("[EMAIL]api_route__verify_email| %s - User not found", client_ip))
		c.JSON(401, json_post{"message": "User not found. Please register first."})
		return
	}

	util__debug_info(fmt.Sprintf("[EMAIL]api_route__verify_email| Creating forum account - %s", client_ip))

	response, err := http.Get(fmt.Sprintf(
		"https://windsofarithia.com/community/api.php?action=register&hash=%s&username=%s&password=%s&email=%s&user_state=%s",
		server_config.xenapi_key, strings.TrimPrefix(s_user.Col_sUserID, "unverified_"), s_user.Col_sUserPW,
		s_user.Col_sEmail, "valid"))
	if err != nil {
		util__debug_info(fmt.Sprintf("[EMAIL]api_route__verify_email| %s - %s", client_ip, err))
		c.JSON(401, json_post{"message": "We encountered an error."})
		return
	}
	defer response.Body.Close()

	// Read response
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		util__debug_info(fmt.Sprintf("[EMAIL]api_route__verify_email| %s - %s", client_ip, err))
		c.JSON(401, json_post{"message": "We encountered an error."})
		return
	}

	// Read if error
	var forum_json_response map[string]interface{}
	if err := json.Unmarshal(contents, &forum_json_response); err != nil {
		util__debug_info(fmt.Sprintf("[EMAIL]api_route__verify_email| %s - %s", client_ip, err))
		c.JSON(401, json_post{"message": "We encountered an error."})
		return
	}

	// Checking if there was an error creating forum account
	if _, ok := forum_json_response["error"]; ok {
		util__debug_info(fmt.Sprintf("[EMAIL]api_route__verify_email| %s - %s", client_ip, forum_json_response))
		c.JSON(401, json_post{"message": "We encountered an error."})
		return
	}

	// Now if there's no error, we want to verify once last time the json returned contains user id
	forum_user_state, ok := forum_json_response["user_state"]
	if !ok {
		util__debug_info(fmt.Sprintf("[EMAIL]api_route__verify_email| Failed to verify json contained valid user_state - %s", client_ip))
		c.JSON(401, json_post{"message": "We encountered an error."})
		return
	}
	forum_user_name, ok := forum_json_response["username"]
	if !ok {
		util__debug_info(fmt.Sprintf("[EMAIL]api_route__verify_email| Failed to verify json contained valid user_state - %s", client_ip))
		c.JSON(401, json_post{"message": "We encountered an error."})
		return
	}

	if forum_user_state.(string) != "valid" || forum_user_name.(string) != strings.TrimPrefix(s_user.Col_sUserID, "unverified_") {
		util__debug_info(fmt.Sprintf("[EMAIL]api_route__verify_email| Failed to verify json contained valid user_state - %s", client_ip))
		c.JSON(401, json_post{"message": "We encountered an error."})
		return
	}

	util__debug_info(fmt.Sprintf("[EMAIL]api_route__verify_email| Forum account created for %s - %s", client_ip, forum_user_state.(string)))

	// Registering user as verified
	if dberr := database__update_verified(s_user.Col_sEmail, s_user.Col_sUserID, s_user.Col_nUserNo); dberr != nil {
		util__debug_info(fmt.Sprintf("[EMAIL]api_route__verify_email| %s - %s", client_ip, dberr))
		c.JSON(401, json_post{"message": "We encountered an error."})
		return
	}

	c.Redirect("https://" + server_config.domain + ".com/verified-mail", 301)
}

func SendVerificationEmail(mail string) (bool) {

	tokenString := util__create_JWT_token(mail)

	params := &EmailContent{
		Title: server_config.project_name,
		SubTitle: "Thanks for signing up!",
		Body: "Please verify your email by clicking the button below.",
		VerifyUrl: "https://" + server_config.domain + ".com/api/auth/verifyemail?token=" + tokenString,
	}

	var doc bytes.Buffer
	t, _ := template.ParseFiles("web/views/mail_verification.html")
	t.Execute(&doc, params)

	sg := sendgrid.NewSendGridClientWithApiKey(server_config.sendgrid_api_key)
	message := sendgrid.NewMail()
	message.AddTo(mail)
	message.SetSubject("Winds Of Arithia Email Verification")
	message.SetHTML(doc.String())
	message.SetFrom("contact@kryptodev.com")
	message.SetFromName("Accounts - Winds Of Arithia")
	if r := sg.Send(message); r == nil {
		util__debug_info("Email sent!")
		return true
	} else {
		util__debug_info(r.Error())
		return false
	}
}

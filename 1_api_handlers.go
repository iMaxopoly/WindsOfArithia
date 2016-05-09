package main

import (
	"fmt"
	"encoding/json"
	"strings"
	"time"
	"net/http"
	"io/ioutil"

	// 3rd Party
	"github.com/kataras/iris"
	"github.com/valyala/fasthttp"
)

/***
 *               _____ _____   _    _                 _ _
 *         /\   |  __ \_   _| | |  | |               | | |
 *        /  \  | |__) || |   | |__| | __ _ _ __   __| | | ___ _ __ ___
 *       / /\ \ |  ___/ | |   |  __  |/ _` | '_ \ / _` | |/ _ \ '__/ __|
 *      / ____ \| |    _| |_  | |  | | (_| | | | | (_| | |  __/ |  \__ \
 *     /_/    \_\_|   |_____| |_|  |_|\__,_|_| |_|\__,_|_|\___|_|  |___/
 *
 *
 */

/***
 *      _____         __                            _
 *      \_   \___    / /  ___   __ _  __ _  ___  __| |
 *       / /\/ __|  / /  / _ \ / _` |/ _` |/ _ \/ _` |
 *    /\/ /_ \__ \ / /__| (_) | (_| | (_| |  __/ (_| |
 *    \____/ |___/ \____/\___/ \__, |\__, |\___|\__,_|
 *                             |___/ |___/
 */

func api_route__is_logged(c *iris.Context) {
	// Get supplied JWT Auth Token from user.
	authToken := string(c.Request.Header.Peek("Authorization"))

	// Verify JWT token.
	if _, JWTBool, JWTErr := util__verify_JWT_token(authToken); !JWTBool {
		if JWTErr != nil {
			c.JSON(401, json_post{"message": JWTErr})
			return
		}
		c.JSON(401, json_post{"message": "Unauthorized."})
		return
	}
	c.JSON(200, json_post{"user": map[string]string{"Status": "Authorized"}})
}

/***
 *       __             _
 *      / /  ___   __ _(_)_ __
 *     / /  / _ \ / _` | | '_ \
 *    / /__| (_) | (_| | | | | |
 *    \____/\___/ \__, |_|_| |_|
 *                |___/
 */

func api_route__login(c *iris.Context) {
	client_ip := string(c.Request.Header.Peek("X-Forwarded-For"))

	// Creating an instance of Account Struct.
	var tAccounts database__user_model

	// Unmarshalling Request Body into tAccounts.
	e := json.Unmarshal(c.Request.Body(), &tAccounts)
	if e != nil {
		util__debug_info(fmt.Sprintf("[LOGIN]api_route__login| %s - %s", client_ip, e))
		c.JSON(401, json_post{"message": "We encountered an error."})
		return
	}

	// Validating Recaptcha Input.
	if (!util__recaptcha_confirm(string(c.Request.Header.Peek("X-Forwarded-For")), tAccounts.RecaptchaResponse)) {
		util__debug_info(fmt.Sprintf("[LOGIN]api_route__login| %s - Invalid captcha", client_ip))
		c.JSON(401, json_post{"message": "Invalid captcha provided."})
		return
	}

	// Sanitizing received Email and Password.
	tAccounts.Col_sEmail = strings.ToLower(strings.Split(tAccounts.Col_sEmail, "@")[0]) + "@" + strings.ToLower(strings.Split(tAccounts.Col_sEmail, "@")[1])
	tAccounts.Col_sUserPW = strings.ToLower(tAccounts.Col_sUserPW)

	// Testing if Email is valid Email REGEX.
	if !util__validate_email_format(tAccounts.Col_sEmail) {
		util__debug_info(fmt.Sprintf("[LOGIN]api_route__login| %s - Invalid email: %s", client_ip, tAccounts.Col_sEmail))
		c.JSON(401, json_post{"message": "Please enter a valid email address."})
		return
	}

	// Looking up database to find the user based on received data.
	s_user, s_exists, err := database__compare_login(tAccounts.Col_sEmail)
	if err != nil {
		util__debug_info(fmt.Sprintf("[LOGIN]api_route__login| %s - %s", client_ip, err))
		c.JSON(401, json_post{"message": "We encountered an error."})
		return
	}

	// If user not found we tell them wrong email or password.
	if !s_exists {
		util__debug_info(fmt.Sprintf("[LOGIN]api_route__login| %s - Wrong Credentials[1]", client_ip))
		c.JSON(401, json_post{"message": "Wrong email/password!"})
		return
	}

	// Checking if unverified
	if strings.Contains(s_user.Col_sUserID, "unverified_") {
		util__debug_info(fmt.Sprintf("[LOGIN]api_route__login| %s - Unverified user", client_ip))
		if SendVerificationEmail(tAccounts.Col_sEmail) == false {
			util__debug_info(fmt.Sprintf("[LOGIN]api_route__login| %s - Email Send Fail", client_ip))
			c.JSON(401, json_post{"message": "We encountered an error."})
			return
		}
		c.JSON(401, json_post{"message": "We sent a new verification email, please verify it!"})
		return
	}

	// Authenticating.
	if tAccounts.Col_sEmail != s_user.Col_sEmail ||
		tAccounts.Col_sUserPW != util__aes_decrypt([]byte(server_config.aes_key), s_user.Col_sUserPW) {
		util__debug_info(fmt.Sprintf("[LOGIN]api_route__login| %s - Wrong Credentials[2]", client_ip))
		c.JSON(401, json_post{"message": "Wrong email/password."})
		return
	}

	// Creating JWT Token after Auth Success
	cToken := util__create_JWT_token(tAccounts.Col_sEmail)
	if cToken == "" {
		util__debug_info(fmt.Sprintf("[LOGIN]api_route__login| %s - Empty Token Generated", client_ip))
		c.JSON(401, json_post{"message": "We encountered an error."})
		return
	}

	// Update JWT on db for user.
	if dberr := database__update_daily(s_user.Col_nUserNo, s_user.Col_sEmail, client_ip, cToken, string(time.Now().UTC().Unix())); dberr != nil {
		util__debug_info(fmt.Sprintf("[LOGIN]api_route__login| %s - %s", client_ip, dberr.Error()))
		c.JSON(401, json_post{"message": "We encountered an error."})
		return
	}

	// Get cookies for forums
	response, err := http.Get(fmt.Sprintf("https://windsofarithia.com/community/api.php?action=login&username=%s&password=%s&ip_address=%s",
		s_user.Col_sUserID, s_user.Col_sUserPW, string(c.Request.Header.Peek("X-Forwarded-For"))))
	if err != nil {
		util__debug_info(fmt.Sprintf("[LOGIN]api_route__login| %s - %s", client_ip, err))
		c.JSON(401, json_post{"message": "We encountered an error."})
		return
	}
	defer response.Body.Close()

	// Read response
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		util__debug_info(fmt.Sprintf("[LOGIN]api_route__login| %s - %s", client_ip, err))
		c.JSON(401, json_post{"message": "We encountered an error."})
		return
	}

	// Read if error
	var forum_json_response map[string]interface{}
	if err := json.Unmarshal(contents, &forum_json_response); err != nil {
		util__debug_info(fmt.Sprintf("[LOGIN]api_route__login| %s - %s", client_ip, err))
		c.JSON(401, json_post{"message": "We encountered an error."})
		return
	}

	// Checking if there was an error using api to login
	if _, ok := forum_json_response["error"]; ok {
		util__debug_info(fmt.Sprintf("[LOGIN]api_route__login| %s - %s", client_ip, forum_json_response))
		c.JSON(401, json_post{"message": "We encountered an error."})
		return
	}

	// Now if there's no error, we want to verify once last time the json returned contains user id
	_, ok := forum_json_response["cookie_id"]
	if !ok {
		util__debug_info(fmt.Sprintf("[LOGIN]api_route__login| Failed to verify json contained valid user_state - %s", client_ip))
		c.JSON(401, json_post{"message": "We encountered an error."})
		return
	}

	if server_config.debug == "true" {
		fmt.Println(forum_json_response)
	}

	cookie1 := &fasthttp.Cookie{}
	cookie1.SetKey(forum_json_response["cookie_name"].(string))
	cookie1.SetValue(forum_json_response["cookie_id"].(string))
	cookie1.SetSecure(true)
	cookie1.SetPath("/")
	cookie1.SetHTTPOnly(true)

	c.Response.Header.SetCookie(cookie1)

	// All success, we send JSON response.
	util__debug_info(fmt.Sprintf("[LOGIN]api_route__login| %s - Login Granted", client_ip))
	c.JSON(200, json_post{
		"token": cToken,
		"user": map[string]string{
			"email": tAccounts.Col_sEmail,
			"verified": func() (string) {
				if strings.Contains(s_user.Col_sUserID, "Unverified_") {
					return "false"
				}
				return "true"
			}()}})
}


/***
 *       __            _     _
 *      /__\ ___  __ _(_)___| |_ ___ _ __
 *     / \/// _ \/ _` | / __| __/ _ \ '__|
 *    / _  \  __/ (_| | \__ \ ||  __/ |
 *    \/ \_/\___|\__, |_|___/\__\___|_|
 *               |___/
 */

func api_route__register(c *iris.Context) {
	client_ip := string(c.Request.Header.Peek("X-Forwarded-For"))

	// Creating an instance of Account Struct.
	var tAccounts database__user_model

	// Read JSON body

	err := json.Unmarshal(c.Request.Body(), &tAccounts)
	if err != nil {
		util__debug_info(fmt.Sprintf("[REGISTER]api_route__register| Failed to read body %s - %s", client_ip, err))
		c.JSON(401, json_post{"message": "We encountered an error."})
		return
	}
	fmt.Println(tAccounts)
	// Unmarshalling Request Body into tAccounts.
	tAccounts.Col_sUserIP = string(c.Request.Header.Peek("X-Forwarded-For"))
	if (!util__recaptcha_confirm(tAccounts.Col_sUserIP, tAccounts.RecaptchaResponse)) {
		util__debug_info(fmt.Sprintf("[REGISTER]api_route__register| %s - Invalid captcha", client_ip))
		c.JSON(401, json_post{"message": "Invalid captcha provided."})
		return
	}

	// Sanitizing user input
	sanitizedUserID, sanitized := util__sanitize_parameter(tAccounts.Col_sUserID, 10)
	if sanitized == false {
		util__debug_info(fmt.Sprintf("[REGISTER]api_route__register| %s - Garbage Input Username - %s", client_ip, tAccounts.Col_sUserID))
		c.JSON(401, json_post{"message": "Please enter userid between 4 and 10 chars(No Special Characters)."})
		return
	}

	sanitizedUserPW, sanitized := util__sanitize_parameter(tAccounts.Col_sUserPW, 20)
	if sanitized == false {
		util__debug_info(fmt.Sprintf("[REGISTER]api_route__register| %s - Garbage Input Password - %s", client_ip, tAccounts.Col_sUserPW))
		c.JSON(401, json_post{"message": "Please enter password between 4 and 20 chars(No Special Characters)."})
		return
	}

	// Validating email regex
	if !util__validate_email_format(tAccounts.Col_sEmail) {
		util__debug_info(fmt.Sprintf("[REGISTER]api_route__register| %s - Garbage Input Email - %s", client_ip, tAccounts.Col_sEmail))
		c.JSON(401, json_post{"message": "Please enter a valid email address."})
		return
	}

	// Checking if disposable email
	if util__test_disposable_email(string(strings.Split(tAccounts.Col_sEmail, "@")[1])) {
		util__debug_info(fmt.Sprintf("[REGISTER]api_route__register| %s - Disposable Input Email - %s", client_ip, tAccounts.Col_sEmail))
		c.JSON(401, json_post{"message": "Please enter a non-disposable email address."})
		return
	}

	// Checking if user is already verifying
	s_exists, err := database__is_currently_unverified(tAccounts.Col_sEmail, sanitizedUserID)
	if err != nil {
		util__debug_info(fmt.Sprintf("[REGISTER]api_route__register| Check already verified %s - %s", client_ip, err))
		c.JSON(401, json_post{"message": "We encountered an error."})
		return
	}
	if s_exists {
		util__debug_info(fmt.Sprintf("[REGISTER]api_route__register| %s - Unverified trying to register", client_ip))
		if SendVerificationEmail(tAccounts.Col_sEmail) == false {
			util__debug_info(fmt.Sprintf("[REGISTER]api_route__register| %s - Email Send Fail", client_ip))
			c.JSON(401, json_post{"message": "We encountered an error."})
			return
		}
		c.JSON(401, json_post{"message": "You are unverified, we have sent a new verification email, please respond."})
		return
	}

	// Verifying user doesn't exist in database by email
	s_exists, err = database__is_existing_email(tAccounts.Col_sEmail)
	if err != nil {
		util__debug_info(fmt.Sprintf("[REGISTER]api_route__register| %s - %s", client_ip, err))
		c.JSON(401, json_post{"message": "We encountered an error."})
		return
	}
	if s_exists {
		util__debug_info(fmt.Sprintf("[REGISTER]api_route__register| %s - Already Registered", client_ip))
		c.JSON(401, json_post{"message": "You are already registered."})
		return
	}

	// Verifying user doesn't exist in database by userid
	s_exists, err = database__is_existing_userid(sanitizedUserID)
	if err != nil {
		util__debug_info(fmt.Sprintf("[REGISTER]api_route__register| %s - %s", client_ip, err))
		c.JSON(401, json_post{"message": "We encountered an error."})
		return
	}
	if s_exists {
		util__debug_info(fmt.Sprintf("[REGISTER]api_route__register| %s - Already Registered", client_ip))
		c.JSON(401, json_post{"message": "You are already registered."})
		return
	}

	cToken := util__create_JWT_token(tAccounts.Col_sEmail)
	if cToken == "" {
		util__debug_info(fmt.Sprintf("[REGISTER]api_route__register| %s - Empty Token Generated", client_ip))
		c.JSON(401, json_post{"message": "We encountered an error."})
		return
	}

	if dberr := database__save(sanitizedUserID,
		sanitizedUserPW, tAccounts.Col_sEmail, tAccounts.Col_sUserIP,
		cToken, time.Now().Format(time.RFC3339)); dberr != nil {
		util__debug_info(fmt.Sprintf("[REGISTER]api_route__register| %s - %s", client_ip, dberr.Error()))
		c.JSON(401, json_post{"message": "We encountered an error."})
		return
	}

	if SendVerificationEmail(tAccounts.Col_sEmail) == false {
		util__debug_info(fmt.Sprintf("[REGISTER]api_route__register| %s - Email Send Fail", client_ip))
		c.JSON(401, json_post{"message": "We encountered an error."})
		return
	}

	util__debug_info(fmt.Sprintf("[REGISTER]api_route__register| %s - Registration Granted", client_ip))
	c.JSON(200, json_post{"token": cToken, "user": map[string]string{"email": tAccounts.Col_sEmail, "verified": "false"}})
}

/***
 *       ___                     _       ___                                    _
 *      / __\__  _ __ __ _  ___ | |_    / _ \__ _ ___ _____      _____  _ __ __| |
 *     / _\/ _ \| '__/ _` |/ _ \| __|  / /_)/ _` / __/ __\ \ /\ / / _ \| '__/ _` |
 *    / / | (_) | | | (_| | (_) | |_  / ___/ (_| \__ \__ \\ V  V / (_) | | | (_| |
 *    \/   \___/|_|  \__, |\___/ \__| \/    \__,_|___/___/ \_/\_/ \___/|_|  \__,_|
 *                   |___/
 */

func api_route__forgot_password(c *iris.Context) {
	util__debug_info(string(c.Request.Body()))
	c.JSON(401, json_post{"message": "We encountered an error."})
}

/***
 *       ___ _                                ___                                    _
 *      / __\ |__   __ _ _ __   __ _  ___    / _ \__ _ ___ _____      _____  _ __ __| |
 *     / /  | '_ \ / _` | '_ \ / _` |/ _ \  / /_)/ _` / __/ __\ \ /\ / / _ \| '__/ _` |
 *    / /___| | | | (_| | | | | (_| |  __/ / ___/ (_| \__ \__ \\ V  V / (_) | | | (_| |
 *    \____/|_| |_|\__,_|_| |_|\__, |\___| \/    \__,_|___/___/ \_/\_/ \___/|_|  \__,_|
 *                             |___/
 */

func api_route__change_password(c *iris.Context) {
	c.JSON(401, json_post{"message": "We encountered an error."})
}

/***
 *       __                         _
 *      / /  ___   __ _  ___  _   _| |_
 *     / /  / _ \ / _` |/ _ \| | | | __|
 *    / /__| (_) | (_| | (_) | |_| | |_
 *    \____/\___/ \__, |\___/ \__,_|\__|
 *                |___/
 */

func api_route__logout(c *iris.Context) {
	cookie1 := &fasthttp.Cookie{}
	cookie1.SetKey("xf_session")
	cookie1.SetSecure(true)
	cookie1.SetPath("/")
	cookie1.SetHTTPOnly(true)
	cookie1.SetExpire(time.Unix(0, 0))

	c.Response.Header.SetCookie(cookie1)

	c.JSON(200, json_post{"message": "Logged out"})
}
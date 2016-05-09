package main

import (
	"time"
	"fmt"
	"strings"
	"errors"
	"regexp"
	"net/http"
	"log"
	"encoding/json"
	"io/ioutil"
	"net/url"
	"encoding/base64"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
	"crypto/sha512"
	"encoding/hex"
	"os/exec"

	//3rd Party
	"github.com/dgrijalva/jwt-go"
)

/***
 *      _    _ _   _ _
 *     | |  | | | (_) |
 *     | |  | | |_ _| |___
 *     | |  | | __| | / __|
 *     | |__| | |_| | \__ \
 *      \____/ \__|_|_|___/
 *
 *
 */

// Json post syntactic sugar
type json_post map[string]interface{}

func util__critical_error_check(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func util__debug_info(msg string) {
	if server_config.debug == "true" {
		fmt.Println(msg)
	}
}

func util__sanitize_parameter(parameter string, maxlen int) (string, bool) {
	if r, err := regexp.MatchString("[a-zA-Z0-9]", parameter); r != true || err != nil {
		return "", false
	}
	if len([]rune(parameter)) > maxlen || len([]rune(parameter)) < 4 {
		return "", false
	}

	return parameter, true
}

func util__validate_email_format(email string) bool {
	if matched, _ := regexp.MatchString("^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\\.[a-zA-Z0-9-.]+$", email); matched {
		return true;
	}
	return false;
}

func util__game_firewall_ip(job <-chan struct_sentry) {
	for {
		select {
		case work_order := <-job:
		// Debug Output
			util__debug_info(
				fmt.Sprintf("[Util]util__game_firewall_ip| Function called with:  %s %s %s params",
					work_order.cmd, work_order.ip, work_order.rulename))

			cmd, err := exec.Command("util/sentry.exe", work_order.cmd, work_order.ip, work_order.rulename).Output()
			if err != nil {
				util__debug_info(fmt.Sprintf("[Util]util__game_firewall_ip| %s", err.Error()))
			}

			output := strings.TrimSpace(string(cmd))
			switch output {
			case "":
				util__debug_info(fmt.Sprintf("[Util]util__game_firewall_ip| Firewall %s %s failed, Unknown error", work_order.cmd, work_order.ip))
				break
			case "BadIP":
				util__debug_info(fmt.Sprintf("[Util]util__game_firewall_ip| Firewall %s %s failed, Bad IP", work_order.cmd, work_order.ip))
				break
			case "True":
				util__debug_info(fmt.Sprintf("[Util]util__game_firewall_ip| Firewall %s %s succeeded", work_order.cmd, work_order.ip))
				break
			case "False":
				util__debug_info(fmt.Sprintf("[Util]util__game_firewall_ip| Firewall %s %s failed, bad rule maybe?", work_order.cmd, work_order.ip))
				break
			default:
				util__debug_info(fmt.Sprintf("[Util]util__game_firewall_ip| Don't Know:  %s %s failed, bad rule maybe? ~ %s", work_order.cmd,
					work_order.ip, output))
				break;
			}
		}
	}
}

/*
util__recaptcha_check uses the client ip address and the client's response input to determine whether or not
the client answered the reCaptcha input question correctly.
It returns a boolean value indicating whether or not the client answered correctly.
 */
func util__recaptcha_check(remoteip, response string) (body []byte) {
	vals := url.Values{"secret": {server_config.google_recaptcha_secret}, "remoteip": {remoteip}, "response": {response}}
	resp, err := http.Get("https://www.google.com/recaptcha/api/siteverify?" + vals.Encode())
	if err != nil {
		util__debug_info(fmt.Sprintf("[Util]util__recaptcha_check| %s", err.Error()))
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		util__debug_info(fmt.Sprintf("[Util]util__recaptcha_check| Read error: could not read body: %s", err.Error()))
	}
	return
}

/*
Confirm is the public interface function.
It calls check, which the client ip address, the challenge code from the reCaptcha form,
and the client's response input to that challenge to determine whether or not
the client answered the reCaptcha input question correctly.
It returns a boolean value indicating whether or not the client answered correctly.
*/
func util__recaptcha_confirm(remoteip, response string) bool {
	var v map[string]interface{}
	if err := json.Unmarshal(util__recaptcha_check(remoteip, response), &v); err != nil {
		util__debug_info(fmt.Sprintf("[Util]util__recaptcha_confirm| JSON error: %s", err.Error()))
	}
	if errors, found := v["error-codes"]; found {
		util__debug_info(fmt.Sprintf("[Util]util__recaptcha_confirm| Recaptcha errors: %s", errors))
	}
	success, ok := v["success"].(bool)
	return ok && success
}

//Email Check
func util__test_disposable_email(x string) (bool) {
	for _, b := range server_config.disposable_emails {
		if b == x {
			return true
		}
	}
	return false
}

func util__create_JWT_token(email string) string {
	token := jwt.New(jwt.GetSigningMethod("RS256")) // Create a Token that will be signed with RSA 256.
	token.Claims["sub"] = email
	token.Claims["exp"] = time.Now().UTC().Add(time.Hour * 12).Unix()
	tokenString, err := token.SignedString(server_config.jwt_privatekey)
	if err != nil {
		fmt.Println(err.Error())
	}
	return tokenString
}

func util__verify_JWT_token(input string) (int, bool, error) {
	if input == "" {
		return 0, false, errors.New("We encountered an error.")
	}
	// Should be a bearer token
	if len(input) < 6 || strings.ToUpper(input[0:6]) != "BEARER" {
		return 0, false, errors.New("Invalid Token.")
	}
	token, err := jwt.Parse(string(input[7:]), func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return false, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return server_config.jwt_publickey, nil
	})

	if err == nil && token.Valid && int(token.Claims["exp"].(float64)) > int(time.Now().UTC().Unix()) {
		s_user, s_exists, err := database__compare_login(token.Claims["sub"].(string))
		if err != nil {
			return 0, false, errors.New("We encountered an error.")
		}

		if !s_exists {
			return 0, false, errors.New("Wrong email/password!")
		}
		return s_user.Col_nUserNo, true, nil
	}

	return 0, false, errors.New("Invalid Token.")
}

//func util__oauth_factory(body []byte, provider string, ip string) (interface{}, bool) {
//	var x map[string]interface{}
//	var y string
//
//	e := json.Unmarshal(body, &x)
//	if e != nil {
//		return map[string]interface{}{}, true
//	}
//
//	if x["code"].(string) == "" {
//		return map[string]interface{}{}, true
//	}
//
//	var conf = &oauth2.Config{
//		ClientID:     "",
//		ClientSecret: "",
//		RedirectURL: "",
//		Scopes:       []string{""},
//		Endpoint: oauth2.Endpoint{},
//	}
//
//	if provider == "google" {
//		conf.Endpoint = google.Endpoint
//		conf.ClientID = server_config.google_oauth_client_id
//		conf.ClientSecret = server_config.google_oauth_secret
//		conf.Scopes = []string{"openid", "profile"}
//	} else if provider == "facebook" {
//		conf.Endpoint = facebook.Endpoint
//		conf.ClientID = server_config.facebook_oauth_client_id
//		conf.ClientSecret = server_config.facebook_oauth_secret
//		conf.Scopes = []string{"public_profile"}
//	}
//	conf.RedirectURL = x["redirectUri"].(string)
//
//	tok, err := conf.Exchange(oauth2.NoContext, x["code"].(string))
//	if err != nil {
//		fmt.Println(err)
//		return map[string]interface{}{}, true
//	}
//
//	client := conf.Client(oauth2.NoContext, tok)
//
//	if provider == "google" {
//		y = "https://www.googleapis.com/plus/v1/people/me/openIdConnect"
//	} else if provider == "facebook" {
//		y = "https://graph.facebook.com/me?access_token=" + tok.AccessToken
//	}
//
//	resp, err := client.Get(y)
//	if err != nil {
//		fmt.Println(err)
//		return map[string]interface{}{}, true
//	}
//
//	raw, err := ioutil.ReadAll(resp.Body)
//	defer resp.Body.Close()
//	if err != nil {
//		return map[string]interface{}{}, true
//	}
//
//	var profile map[string]interface{}
//	if err := json.Unmarshal(raw, &profile); err != nil {
//		return map[string]interface{}{}, true
//	}
//
//	fmt.Println(profile)
//
//	var LoginResponse struct {
//		Token string `json:"token"`
//		User  struct {
//			      Email       string `json:"email"`
//			      GoogleId    string `json:"googleId,omitempty"`
//			      FacebookId  string `json:"facebookId,omitempty"`
//			      DisplayName string `json:"displayName"`
//			      Gender      string `json:"gender,omitempty"`
//		      } `json:"user"`
//	}
//
//	rBool, rUser, rUserNo, rErr := database__exists(profile["email"].(string))
//	if rErr != nil {
//		return map[string]interface{}{}, true
//	}
//
//	if rBool {
//		ew := json.Unmarshal([]byte(rUser), &LoginResponse.User)
//		if ew != nil {
//			fmt.Println(ew)
//			return map[string]interface{}{}, true
//		}
//
//		LoginResponse.Token = util__create_JWT_token(profile["email"].(string))
//		if dberr := database__update_daily(rUserNo, ip, LoginResponse.Token, string(time.Now().UTC().Unix())); dberr != nil {
//			fmt.Println(dberr.Error())
//			return map[string]interface{}{}, true
//		}
//		return LoginResponse, false
//	}
//
//	LoginResponse.User.DisplayName = profile["name"].(string)
//	if profile["gender"].(string) != "" {
//		LoginResponse.User.Gender = profile["gender"].(string)
//	}
//
//	LoginResponse.User.Email = profile["email"].(string)
//	if provider == "google" {
//		LoginResponse.User.GoogleId = profile["sub"].(string)
//	} else if provider == "facebook" {
//		LoginResponse.User.FacebookId = profile["id"].(string)
//	}
//
//	LoginResponse.Token = util__create_JWT_token(LoginResponse.User.Email)
//
//	if LoginResponse.Token == "" {
//		return map[string]interface{}{}, true
//	}
//
//	displayname := strings.Split(LoginResponse.User.DisplayName, " ")
//	fname := displayname[0]
//	lname := displayname[1]
//
//	if strings.ToLower(LoginResponse.User.Gender) == "male" {
//		LoginResponse.User.Gender = string("0")
//	} else if strings.ToLower(LoginResponse.User.Gender) == "female" {
//		LoginResponse.User.Gender = string("1")
//	} else {
//		LoginResponse.User.Gender = string("2")
//	}
//
//	gender, err := strconv.Atoi(LoginResponse.User.Gender)
//	if err != nil {
//		return map[string]interface{}{}, true
//	}
//
//	database__save_oauth(LoginResponse.User.Email, fname, lname, gender,
//		LoginResponse.User.FacebookId, LoginResponse.User.GoogleId, ip, LoginResponse.Token, string(time.Now().UTC().Unix()))
//
//	return map[string]string{"token":LoginResponse.Token, "email": LoginResponse.User.Email}, false
//}

func util__aes_encrypt(key []byte, text string) string {
	plaintext := []byte(text)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	ciphertext := make([]byte, aes.BlockSize + len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	// convert to base64
	return base64.URLEncoding.EncodeToString(ciphertext)
}

func util__aes_decrypt(key []byte, cryptoText string) string {
	ciphertext, _ := base64.URLEncoding.DecodeString(cryptoText)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	if len(ciphertext) < aes.BlockSize {
		panic("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)

	// XORKeyStream can work in-place if the two arguments are the same.
	stream.XORKeyStream(ciphertext, ciphertext)

	return fmt.Sprintf("%s", ciphertext)
}
// END AES encryption and decryption

// Simple Xor encrypt-decrypt
func util__xor_encryptdecrypt(toEncrypt []byte) ([]byte) {
	var key byte = 'W'; //Any char will work
	output := toEncrypt;
	for i := 0; i < len(toEncrypt); i++ {
		output[i] = toEncrypt[i] ^ key;
	}
	return output;
}

// SHA-512 utility
func util__sha512_tohex(text string) string {
	hasher := sha512.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}
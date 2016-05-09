package main

import (
	"fmt"
	"os"
	"strconv"
	"io/ioutil"

	// 3rd party
	"github.com/kataras/iris"
)

/***
 *     __          ________ ____   _____  ____   _____ _  ________ _______
 *     \ \        / /  ____|  _ \ / ____|/ __ \ / ____| |/ /  ____|__   __|
 *      \ \  /\  / /| |__  | |_) | (___ | |  | | |    | ' /| |__     | |
 *       \ \/  \/ / |  __| |  _ < \___ \| |  | | |    |  < |  __|    | |
 *        \  /\  /  | |____| |_) |____) | |__| | |____| . \| |____   | |
 *         \/  \/   |______|____/|_____/ \____/ \_____|_|\_\______|  |_|
 *
 *
 */

// [Localhost Only]
func api__game_token_confirm(c *iris.Context) {
	clientIP := string(c.Request.Header.Peek("X-Forwarded-For"))

	if clientIP != "127.0.0.1" {
		util__debug_info(fmt.Sprintf("[GAME]api__game_token_confirm| %s - Illegal Connection[1]", clientIP))
		c.Redirect("/", 302)
		return
	}

	//clientIP, _, _ := net.SplitHostPort(c.RemoteIP().String())
	//if clientIP != "127.0.0.1" || clientIP != "::1" {
	//	util__debug_info(fmt.Sprintf("[GAME]api__game_token_confirm| %s - Illegal Connection[2]", clientIP))
	//	c.Redirect("/", 302)
	//	return
	//}

	//sanitizer
	tokenfromlauncher, xbool := util__sanitize_parameter(c.URLParam("token"), 50)
	if tokenfromlauncher == "" || xbool == false {
		util__debug_info(fmt.Sprintf("[GAME]api__game_token_confirm| %s - Empty token confirmation call", clientIP))
		c.Text(401, "")
		return
	}

	//evaluating user exists
	fmt.Println("Talking to the game evaluating")
	nEMID, sUsername, isError := database__restservice_verify_token(tokenfromlauncher)
	if isError {
		util__debug_info(fmt.Sprintf("[GAME]api__game_token_confirm| %s - No such player", clientIP))
		c.Text(401, "NoToken")
		return
	}
	RestServiceJson := "{\"token_age\":0,\"user_id\":" + strconv.Itoa(nEMID) + ",\"login\":\"" + sUsername + "\",\"user_role\":\"user\",\"blocked\":false}"
	c.Text(200, RestServiceJson)
}

func api__game_launcher_version(c *iris.Context) {
	client_ip := string(c.Request.Header.Peek("X-Forwarded-For"))
	if string(c.UserAgent()) != conf__launcher_useragent {
		util__debug_info(fmt.Sprintf("[GAME]api__game_launcher_version| %s - Invalid U-A", client_ip))
		c.Redirect("/", 302)
		return
	}
	dat, err := ioutil.ReadFile("patches/launcher.version")
	if err != nil {
		util__debug_info(fmt.Sprintf("[GAME]api__game_launcher_version| %s - File doesn't exist - %s", client_ip, err))
	}
	c.Text(200, string(dat))
}

func api__game_patch_version(c *iris.Context) {
	client_ip := string(c.Request.Header.Peek("X-Forwarded-For"))
	if string(c.UserAgent()) != conf__launcher_useragent {
		util__debug_info(fmt.Sprintf("[GAME]api__game_patch_version| %s - Invalid U-A", client_ip))
		c.Redirect("/", 302)
		return
	}
	dat, err := ioutil.ReadFile("patches/game.version")
	if err != nil {
		util__debug_info(fmt.Sprintf("[GAME]api__game_patch_version| %s - File doesn't exist - %s", client_ip, err))
	}
	c.Text(200, string(dat))
}

func api__game_patch_download(c *iris.Context) {
	client_ip := string(c.Request.Header.Peek("X-Forwarded-For"))

	version := c.Param("id")
	fileRequested := "patches/game" + version + ".zip"
	if _, err := os.Stat(fileRequested); os.IsNotExist(err) {
		fmt.Println(client_ip, " - (StartGamePatchDownload) - no such file:", fileRequested)
		c.Redirect("/", 302)
		return
	}
	c.SendFile(fileRequested, "game"+version+".zip")
}

func api__game_launcher_download(c *iris.Context) {
	client_ip := string(c.Request.Header.Peek("X-Forwarded-For"))
	if string(c.UserAgent()) != conf__launcher_useragent {
		util__debug_info(fmt.Sprintf("[GAME]api__game_launcher_download| %s - Invalid U-A", client_ip))
		c.Redirect("/", 302)
		return
	}
	version := c.Param("id")
	fileRequested := "patches/launcher" + version + ".zip"
	if _, err := os.Stat(fileRequested); os.IsNotExist(err) {
		util__debug_info(fmt.Sprintf("[GAME]api__game_launcher_download| %s - File doesn't exist - %s", client_ip, fileRequested))
		c.Redirect("/", 302)
		return
	}
	c.SendFile(fileRequested, "launcher" + version + ".zip")
}

func api__game_launcher_news(c *iris.Context) {
	client_ip := string(c.Request.Header.Peek("X-Forwarded-For"))
	if string(c.UserAgent()) != conf__launcher_useragent {
		util__debug_info(fmt.Sprintf("[GAME]api__game_launcher_news| %s - Invalid U-A", client_ip))
		c.Redirect("/", 302)
		return
	}
	dat, err := ioutil.ReadFile("patch_news.view")
	if err != nil {
		util__debug_info(fmt.Sprintf("[GAME]api__game_launcher_news| %s - File doesn't exist - %s", client_ip, err))
		c.Text(200, "Failed to retrieve news.")
	}
	c.Text(200, string(dat))
}
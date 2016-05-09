package main

import (
	// 3rd Party
	"github.com/kataras/iris"
)

/***
 *      _____   ____  _    _ _______ ______
 *     |  __ \ / __ \| |  | |__   __|  ____|
 *     | |__) | |  | | |  | |  | |  | |__
 *     |  _  /| |  | | |  | |  | |  |  __|
 *     | | \ \| |__| | |__| |  | |  | |____
 *     |_|  \_\\____/ \____/_ _|_|_ |______|______ _____   _____
 *     | |  | |   /\   | \ | |  __ \| |    |  ____|  __ \ / ____|
 *     | |__| |  /  \  |  \| | |  | | |    | |__  | |__) | (___
 *     |  __  | / /\ \ | . ` | |  | | |    |  __| |  _  / \___ \
 *     | |  | |/ ____ \| |\  | |__| | |____| |____| | \ \ ____) |
 *     |_|  |_/_/    \_\_| \_|_____/|______|______|_|  \_\_____/
 *
 *
 */

func route__SEO_Sitemap(c *iris.Context) {
	c.XML(200, []byte(
	//`<?xml version="1.0" encoding="UTF-8"?>` +
		`<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">` +
			`<url><loc>https://windsofarithia.com/home</loc></url>` +
			`<url><loc>https://windsofarithia.com/news</loc></url>` +
			`<url><loc>https://windsofarithia.com/community/</loc></url>` +
			`<url><loc>https://windsofarithia.com/login</loc></url>` +
			`<url><loc>https://windsofarithia.com/register</loc></url>` +
			`<url><loc>https://windsofarithia.com/store</loc></url>` +
			`<url><loc>https://windsofarithia.com/chat</loc></url>` +
			`<url><loc>https://windsofarithia.com</loc></url>` +
			`</urlset>`))
}

func route__SEO_Robotstxt(c *iris.Context) {
	c.Text(200, `User-Agent: *
Disallow:
Allow: /
Sitemap: https://windsofarithia.com/sitemap.xml`)
}

func route__base(c *iris.Context) {
	c.Render("base.html", nil)
}

func route__footer(c *iris.Context) {
	c.Render("footer.html", nil)
}

func route__index(c *iris.Context) {
	c.Render("index.html", nil)
}

func route__header(c *iris.Context) {
	c.Render("header.html", nil)
}

func route__home(c *iris.Context) {
	c.Render("home.html", nil)
}

func route__store(c *iris.Context) {
	authToken := string(c.Request.Header.Peek("Authorization"))
	if _, JWTBool, JWTErr := util__verify_JWT_token(authToken); !JWTBool {
		if JWTErr != nil {
			c.JSON(401, json_post{"message": JWTErr})
			return
		}
		c.JSON(401, json_post{"message": "Unauthorized."})
		return
	}
	c.Render("store.html", nil)
}

func route__donate(c *iris.Context) {
	authToken := string(c.Request.Header.Peek("Authorization"))
	if _, JWTBool, JWTErr := util__verify_JWT_token(authToken); !JWTBool {
		if JWTErr != nil {
			c.JSON(401, json_post{"message": JWTErr})
			return
		}
		c.JSON(401, json_post{"message": "Unauthorized."})
		return
	}
	c.Render("donate.html", nil)
}

func route__register(c *iris.Context) {
	authToken := string(c.Request.Header.Peek("Authorization"))
	if _, JWTBool, _ := util__verify_JWT_token(authToken); !JWTBool {
		c.Render("register.html", nil)
		return
	}
	c.JSON(401, json_post{"message": "You are already logged in."})
}

func route__login(c *iris.Context) {
	authToken := string(c.Request.Header.Peek("Authorization"))
	if _, JWTBool, _ := util__verify_JWT_token(authToken); !JWTBool {
		c.Render("login.html", nil)
		return
	}
	c.JSON(401, json_post{"message": "You are already logged in."})
}

func route__forgot_password(c *iris.Context) {
	authToken := string(c.Request.Header.Peek("Authorization"))
	if _, JWTBool, _ := util__verify_JWT_token(authToken); !JWTBool {
		c.Render("forgot_password.html", nil)
		return
	}
	c.JSON(401, json_post{"message": "You are already logged in."})
}

func route__change_password(c *iris.Context) {
	authToken := string(c.Request.Header.Peek("Authorization"))
	if _, JWTBool, _ := util__verify_JWT_token(authToken); !JWTBool {
		c.Render("change_password.html", nil)
		return
	}
	c.JSON(401, json_post{"message": "You are already logged in."})
}

func route__verified_mail(c *iris.Context) {
	authToken := string(c.Request.Header.Peek("Authorization"))
	if _, JWTBool, _ := util__verify_JWT_token(authToken); !JWTBool {
		c.Render("verified_email.html", nil)
		return
	}
	c.JSON(401, json_post{"message": "You are already logged in."})
}

func route__unverified_mail(c *iris.Context) {
	authToken := string(c.Request.Header.Peek("Authorization"))
	if _, JWTBool, _ := util__verify_JWT_token(authToken); !JWTBool {
		c.Render("unverified_email.html", nil)
		return
	}
	c.JSON(401, json_post{"message": "You are already logged in."})
}

func route__disclaimer(c *iris.Context) {
	c.Render("disclaimer.html", nil)
}

func route__privacy_policy(c *iris.Context) {
	c.Render("privacy_policy.html", nil)
}

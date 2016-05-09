package main

/***
 *                                                               .-'''-.
 *                                 _______                      '   _    \
 *                  .--.   _..._   \  ___ `'.                 /   /` '.   \
 *           _     _|__| .'     '.  ' |--.\  \               .   |     \  '   _.._
 *     /\    \\   //.--..   .-.   . | |    \  '              |   '      |  '.' .._|
 *     `\\  //\\ // |  ||  '   '  | | |     |  '             \    \     / / | '
 *       \`//  \'/  |  ||  |   |  | | |     |  |    _         `.   ` ..' /__| |__
 *        \|   |/   |  ||  |   |  | | |     ' .'  .' |           '-...-'`|__   __|
 *         '        |  ||  |   |  | | |___.' /'  .   | /                    | |
 *                  |__||  |   |  |/_______.'/ .'.'| |//                    | |
 *                      |  |   |  |\_______|/.'.'.-'  /                     | |
 *                      |  |   |  |          .'   \_.'                      | |
 *                      '--'   '--'                                         |_|
 *                                .--.           .        .--.
 *                                |__|         .'|        |__|
 *                        .-,.--. .--.     .| <  |        .--.
 *                  __    |  .-. ||  |   .' |_ | |        |  |    __
 *               .:--.'.  | |  | ||  | .'     || | .'''-. |  | .:--.'.
 *              / |   \ | | |  | ||  |'--.  .-'| |/.'''. \|  |/ |   \ |
 *              `" __ | | | |  '- |  |   |  |  |  /    | ||  |`" __ | |
 *               .'.''| | | |     |__|   |  |  | |     | ||__| .'.''| |
 *              / /   | |_| |            |  '.'| |     | |    / /   | |_
 *              \ \._,\ '/|_|            |   / | '.    | '.   \ \._,\ '/
 *               `--'  `"                `'-'  '---'   '---'   `--'  `"
 */

import (
	"fmt"

	// 3rd Party
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/logger"
	"github.com/kataras/iris/middleware/recovery"
)


// Setup Sentry
type struct_sentry struct {
	cmd string
	ip string
	rulename string
}
var chan_sentry chan struct_sentry

/***
 *      __  __       _
 *     |  \/  |     (_)
 *     | \  / | __ _ _ _ __
 *     | |\/| |/ _` | | '_ \
 *     | |  | | (_| | | | | |
 *     |_|  |_|\__,_|_|_| |_|
 *
 *
 */
func main() {
	load__project_config()
	load__disposable_email_and_jwt_lists()
	load__database_init()
	util__debug_info("[MAIN]| Project Modules Loaded")

	// Loading views
	iris.Config().Render.Directory = "./web/views/*.html"
	util__debug_info("[MAIN]| Views Loaded")

	// Middleware
	if server_config.debug == "true" {
		iris.UseFunc(logger.Default(logger.Options{IP: true}))
	}
	errorLogger := logger.Default(logger.Options{Latency: false, IP: true})
	iris.Use(recovery.Recovery())
	util__debug_info("[MAIN]| Middlewares Loaded")

	// Define custom HTTP handler
	iris.OnError(404,func (c *iris.Context){
	    errorLogger.Serve(c)
	    c.Render("base.html", nil)
        c.SetStatusCode(301)
	})

	/***
	 *      _____   ____  _    _ _______ _____ _   _  _____
	 *     |  __ \ / __ \| |  | |__   __|_   _| \ | |/ ____|
	 *     | |__) | |  | | |  | |  | |    | | |  \| | |  __
	 *     |  _  /| |  | | |  | |  | |    | | | . ` | | |_ |
	 *     | | \ \| |__| | |__| |  | |   _| |_| |\  | |__| |
	 *     |_|  \_\\____/ \____/   |_|  |_____|_| \_|\_____|
	 *
	 *
	 */

	// Static Files
	iris.Static("/static", "./web/static/", 1)

	// SEO for the website
	iris.Get("/sitemap.xml", route__SEO_Sitemap)
	iris.Get("/robots.txt", route__SEO_Robotstxt)

	// Website-wide Basic API for core functions here.
	iris.Post("/api/islogged", api_route__is_logged)
	iris.Get("/api/auth/verifyemail", api_route__verify_email)
	iris.Post("/api/auth/register", api_route__register)
	iris.Post("/api/auth/login", api_route__login)
	iris.Post("/api/auth/forgotpassword", api_route__forgot_password)
	iris.Post("/api/auth/changepassword", api_route__change_password)
	iris.Post("/api/auth/logout", api_route__logout)
	//webserver.Post("/api/auth/google", api_route_oauth__google)
	//webserver.Post("/api/auth/facebook", api_route_oauth__facebook)

	// Website-wide Routing configuration here.
	iris.Get("/", route__base)
	iris.Get("/init-index", route__index)
	iris.Get("/init-footer", route__footer)
	iris.Get("/init-header", route__header)
	iris.Get("/woa-home", route__home)
	iris.Get("/woa-store", route__store)
	iris.Get("/woa-donate", route__donate)
	iris.Get("/woa-register", route__register)
	iris.Get("/woa-login", route__login)
	iris.Get("/woa-forgot-password", route__forgot_password)
	iris.Get("/woa-change-password", route__change_password)
	iris.Get("/woa-verified-mail", route__verified_mail)
	iris.Get("/woa-unverified-mail", route__unverified_mail)
	iris.Get("/woa-disclaimer", route__disclaimer)
	iris.Get("/woa-privacy-policy", route__privacy_policy)

	// Game-related API Routing configuration here.
	iris.Get("/api/game/launchernews", api__game_launcher_news)
	iris.Get("/api/game/launcherversion", api__game_launcher_version)
	iris.Get("/api/game/patchversion", api__game_patch_version)
	iris.Get("/api/game/launcherdownload/:id", api__game_launcher_download)
	iris.Get("/api/game/patchdownload/:id", api__game_patch_download)
	iris.Get("/user/v1/getInfo", api__game_token_confirm) //Locahost only
	util__debug_info("[MAIN]| Routes Defined")


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

	iris.Get("/ws", websocket__route)
	util__debug_info(fmt.Sprintf("[MAIN]| Started WS server on port %s @ /ws", server_config.port))

	/***
	 *      _              _    _ _   _  _____ _    _
	 *     | |        /\  | |  | | \ | |/ ____| |  | |
	 *     | |       /  \ | |  | |  \| | |    | |__| |
	 *     | |      / /\ \| |  | | . ` | |    |  __  |
	 *     | |____ / ____ \ |__| | |\  | |____| |  | |
	 *     |______/_/    \_\____/|_| \_|\_____|_|  |_|
	 *
	 *
	 */

	// Channel definitions
	chan_sentry = make(chan struct_sentry, 1)

	fmt.Println(util__aes_decrypt([]byte(server_config.aes_key), "0j-MHhUcKzshF5icYUrKUtH-QASBN02L_zkRfOneXlOc"))

	// Ready our goroutines
	go util__game_firewall_ip(chan_sentry)

	util__debug_info(fmt.Sprintf("[MAIN]| Started http server on port %s", server_config.port))
	iris.Listen(":" + server_config.port)
}
package main

/***
 *       ____         _    _ _______ _    _
 *      / __ \   /\  | |  | |__   __| |  | |
 *     | |  | | /  \ | |  | |  | |  | |__| |
 *     | |  | |/ /\ \| |  | |  | |  |  __  |
 *     | |__| / ____ \ |__| |  | |  | |  | |
 *      \____/_/    \_\____/   |_|  |_|  |_|
 *
 *
 */

//func api_route_oauth__google(c *echo.Context)(error) {
//	body, err := ioutil.ReadAll(c.Request().Body)
//	if err != nil {
//		fmt.Println(err)
//		c.JSON(401, json_post{"message": "We encountered an error."})
//		return nil
//	}
//
//	something, isError := util__oauth_factory(body, "google", c.Request().Header.Get("X-Forwarded-For"))
//	if isError {
//		c.JSON(401, json_post{"message": "We lol encountered an error."})
//		return nil
//	}
//	md, _ := something.(map[string]string)
//	c.JSON(200, json_post{"token": string(md["token"]), "user": map[string]string{"email": string(md["email"])}})
//	return nil
//}
//
//func api_route_oauth__facebook(c *echo.Context)(error) {
//	body, err := ioutil.ReadAll(c.Request().Body)
//	if err != nil {
//		fmt.Println(err)
//		c.JSON(401, json_post{"message": "We encountered an error."})
//		return nil
//	}
//
//	something, isError := util__oauth_factory(body, "facebook", c.Request().Header.Get("X-Forwarded-For"))
//	if isError {
//		c.JSON(401, json_post{"message": "We encountered an error."})
//		return nil
//	}
//
//	md, _ := something.(map[string]string)
//	c.JSON(200, json_post{"token": string(md["token"]), "user": map[string]string{"email": string(md["email"])}})
//	return nil
//}
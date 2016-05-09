define(function (require) {
	'use strict';
	var angular = require('angular');
	return angular.module('app.controllers.ChatCtrl', [])
		.controller('ChatCtrl', ["$window", function ($window) {
			$window.location.href = '/community/chat/';
		}]);
});
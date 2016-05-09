define(function (require) {
	'use strict';
	var angular = require('angular');
	return angular.module('app.controllers.NewsCtrl', [])
		.controller('NewsCtrl', ["$window", function ($window) {
			$window.location.href = '/community/forums/news-and-announcements.2/';
		}]);
});
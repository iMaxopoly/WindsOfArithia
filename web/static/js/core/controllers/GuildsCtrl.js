define(function (require) {
	'use strict';
	var angular = require('angular');
	return angular.module('app.controllers.GuildsCtrl', [])
		.controller('GuildsCtrl', ["$window", function ($window) {
			$window.location.href = '/community/forums/guild-discussion.17/';
		}]);
});
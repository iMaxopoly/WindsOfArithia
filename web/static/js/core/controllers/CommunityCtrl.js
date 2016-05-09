define(function (require) {
	'use strict';
	var angular = require('angular');
	return angular.module('app.controllers.CommunityCtrl', [])
		.controller('CommunityCtrl', ["$window", function ($window) {
			$window.location.href = '/community/';
		}]);
});
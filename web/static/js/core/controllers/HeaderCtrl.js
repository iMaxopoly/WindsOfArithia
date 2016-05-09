define(function (require) {
	'use strict';
	var angular = require('angular');
	return angular.module('app.controllers.HeaderCtrl', []).controller('HeaderCtrl',
		['$scope', '$auth', '$state',
			function ($scope, $auth, $state) {
				$scope.isAuthenticated = $auth.isAuthenticated;
			}]);
});

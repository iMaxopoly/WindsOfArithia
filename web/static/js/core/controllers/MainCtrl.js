define(function (require) {
	'use strict';
	var angular = require('angular');
	return angular.module('app.controllers.MainCtrl', [])
		.controller('MainCtrl', ["$scope", "$uibModal", function ($scope, $uibModal) {
			$scope.disclaimerModel = function (size) {
				$uibModal.open({
					templateUrl: 'woa-disclaimer',
					size: size
				});
			};
			$scope.privacyPolicyModel = function (size) {
				$uibModal.open({
					templateUrl: 'woa-privacy-policy',
					size: size
				});
			};
		}]);
});
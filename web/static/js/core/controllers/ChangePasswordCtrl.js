define(function (require) {
	'use strict';
	var angular = require('angular');
	return angular.module('app.controllers.ChangePasswordCtrl', [])
		.controller('ChangePasswordCtrl', ["$scope", "alert", "$http", "vcRecaptchaService",
			function ($scope, alert, $http, vcRecaptchaService) {
			$scope.gRecaptchResponse = null;
			$scope.gRecaptchWidgetId = null;
			$scope.gRecaptchModel = {
				key: '6LcQ4xYTAAAAAFonPEBORwlB7r3m2J1ctVcLF5cu'
			};
			$scope.setgRecaptchResponse = function (response) {
				$scope.gRecaptchResponse = response;
			};
			$scope.setgRecaptchWidgetId = function (widgetId) {
				$scope.gRecaptchWidgetId = widgetId;
			};
			$scope.user = {};
			user.gRecaptchResponse = $scope.gRecaptchResponse;
			$scope.submit = function () {
				$http({
					method: 'POST',
					url: '//windsofarithia.com/api/auth/woa-change-password',
					data: $scope.user,
					headers: {'Content-Type': 'application/x-www-form-urlencoded'}
				})
					.success(function (data) {
						console.log(data);
					})
					.catch(handleError);
			};

			$scope.authenticate = function (provider) {
				$auth.authenticate(provider).then(function (res) {
					alert('success', 'Welcome!', 'Thanks for coming back!');
				}, handleError);
			};

			function handleError(err) {
				vcRecaptchaService.reload($scope.gRecaptchWidgetId);
				alert('warning', 'Oops! ', err.data.message);
			}
		}]);
});
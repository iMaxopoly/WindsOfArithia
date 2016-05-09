define(function (require) {
	'use strict';
	var angular = require('angular');
	return angular.module('app.controllers.HomeCtrl', [])
		.controller('HomeCtrl', ["$scope", "$uibModal", function ($scope, $uibModal) {
			$scope.DownloadModel = function (size) {
				$uibModal.open({
					template: 'This feature is not available yet. Please stay tuned with our social-networks for information related.',
					size: size
				});
			};
			$scope.slideInterval = 4000;
			$scope.slides = [
				{
					image: '/static/img/slide1.png'
				},
				{
					image: '/static/img/slide2.png'
				},
				{
					image: '/static/img/slide3.png'
				},
				{
					image: '/static/img/slide4.png'
				}
			];
		}]);
});
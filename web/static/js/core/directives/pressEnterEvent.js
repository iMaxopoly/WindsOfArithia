define(function (require) {
	'use strict';
	var angular = require('angular');
	return angular.module('app.directives.pressEnterEvent', []).directive('ngEnter', function () {
		return function (scope, element, attrs) {
			element.bind("keydown keypress", function (event) {
				if (event.which === 13) {
					scope.$apply(function () {
						scope.$eval(attrs.ngEnter, {'event': event});
					});

					event.preventDefault();
				}
			});
		};
	});
});
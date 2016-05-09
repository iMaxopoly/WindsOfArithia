define(function (require) {
	'use strict';
	var angular = require('angular');
	return angular.module('app.directives.toLowerCase', []).directive('toLowerCase', function ($parse) {
		return {
			require: 'ngModel',
			link: function (scope, element, attrs, modelCtrl) {

				modelCtrl.$parsers.push(function (inputValue) {

					var transformedInput = inputValue.toLowerCase().replace(/ /g, '');

					if (transformedInput != inputValue) {
						modelCtrl.$setViewValue(transformedInput);
						modelCtrl.$render();
					}

					return transformedInput;
				});
			}
		}
	});
});
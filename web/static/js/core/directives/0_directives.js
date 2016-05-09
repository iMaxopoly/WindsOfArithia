define(function (require) {
	'use strict';

	var angular = require('angular'),
		validateEquals = require('./validateEquals'),
		pressEnterEvent = require('./pressEnterEvent'),
		toLowerCase = require('./toLowerCase');

	return angular.module('app.directives', [validateEquals.name, pressEnterEvent.name, toLowerCase.name]);
});
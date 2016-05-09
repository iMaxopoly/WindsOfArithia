define(function (require) {
	'use strict';

	var angular = require('angular'),
		alert = require('./alert');

	return angular.module('app.services', [
		alert.name]);
});
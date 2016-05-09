define(function () {
	"use strict";

	var moment = require('moment');
	return $(function () {
		setInterval(function () {
			$('#divUTC').text(moment.utc().format("Do MMM, YYYY, h:mm a"));
		}, 1000);
	});
});
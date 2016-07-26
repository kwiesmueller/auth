'use strict';

angular.module('authFilters', []);

angular.module('authFilters').filter('length', function () {
	return function (text) {
		return ('' + (text || '')).length;
	}
});

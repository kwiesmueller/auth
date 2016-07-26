'use strict';

angular.module('authServices', ['ngResource', 'ngCookies']);

angular.module('authServices').factory('UserService', ['$resource', '$http', 'UserDao', function ($resource, $http, UserDao) {
	var service = {};
	service.list = function () {
		return UserDao.query().$promise;
	};
	service.create = function (userName) {
		return $http({
			method: 'POST',
			url: 'api/user',
			data: {
				Name: userName
			}
		});
	};
	return service;
}]);

angular.module('authServices').factory('UserDao', ['$resource', function ($resource) {
	return $resource('api/user/:Id', {}, {
		query: {
			method: 'GET', params: {}, isArray: true
		},
	});
}]);

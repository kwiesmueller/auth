'use strict';

angular.module('authControllers', []);

angular.module('authControllers').controller('UserListCtrl', ['$scope', '$log', 'UserService', function ($scope, $log, UserService) {
	$log.debug('init UserListCtrl');

	$scope.reset = function () {
		$log.debug('reset');
		$scope.messages = [];
		$scope.users = [];
	};

	$scope.loadUsers = function () {
		$log.debug('loadUsers');
		UserService.list().then(function (result) {
			$log.debug('list users success');
			$scope.users = result;
		}, function (error) {
			$scope.users = [];
			$log.debug('list users failed');
			$scope.messages.push('list users failed: ' + error);
		});
	};

	$scope.reset();
	$scope.loadUsers();
}]);

angular.module('authControllers').controller('UserCreateCtrl', ['$scope', '$log', '$location', 'UserService', function ($scope, $log, $location, UserService) {
	$log.debug('init UserCreateCtrl');

	$scope.reset = function () {
		$log.debug('reset');
		$scope.messages = [];
		$scope.user = '';
	};

	$scope.submit = function () {
		$log.debug('submit');
		UserService.create($scope.user).then(function (response) {
			$log.debug('create user successful');
			$location.path('/packages/' + $scope.user);
		}, function (error) {
			$log.warn('create user failed: ' + error);
			$scope.messages.push('create user failed: ' + error);
		});
	};

	$scope.reset();
}]);

angular.module('authControllers').controller('NaviTopCtrl', ['$scope', '$log', function ($scope, $log) {
	$log.debug('init NaviTopCtrl');
}]);

angular.module('authControllers').controller('NaviBottomCtrl', ['$scope', '$log', function ($scope, $log) {
	$log.debug('init NaviBottomCtrl');
}]);

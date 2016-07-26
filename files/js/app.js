'use strict';

angular.module('authApp', [
  'ngRoute',
  'authControllers',
  'authDirectives',
  'authFilters',
  'authServices'
]);

angular.module('authApp').config(['$routeProvider', function ($routeProvider) {
  $routeProvider.
    when('/users', {
      templateUrl: 'partials/user/list.html',
      controller: 'UserListCtrl'
    }).
    when('/user/create', {
      templateUrl: 'partials/user/create.html',
      controller: 'UserCreateCtrl'
    }).
    otherwise({
      redirectTo: '/users'
    });
}]);

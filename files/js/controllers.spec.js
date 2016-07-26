'use strict';

describe('UserListCtrl TestSuite', function () {
  beforeEach(module('authApp'));

  beforeEach(inject(function (_$controller_) {
    $controller = _$controller_;
  }));

  it('controller not null', function () {
    var $scope = {};
    var controller = $controller('UserListCtrl', {$scope: $scope});
    expect(controller).not.toBe(null);
  });
});

describe('UserCreateCtrl TestSuite', function () {
  beforeEach(module('authApp'));

  beforeEach(inject(function (_$controller_) {
    $controller = _$controller_;
  }));

  it('controller not null', function () {
    var $scope = {};
    var controller = $controller('UserCreateCtrl', {$scope: $scope});
    expect(controller).not.toBe(null);
  });
});

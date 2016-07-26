'use strict';

describe('UserService TestSuite', function () {
	beforeEach(module('authApp'));

	it('convert return not null', inject(['UserService', function (UserService) {
		var result = UserService.get(1);
		expect(result).not.toBe(null);
	}]));
});

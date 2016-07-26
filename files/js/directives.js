'use strict';

angular.module('authDirectives', []);

angular.module('authDirectives').directive('confirmationNeeded', ['$log', function ($log) {
  return {
    link: function (scope, element, attr) {
      var msg = attr.confirmationNeeded || "sure?";
      var clickAction = attr.ngClick;
      element.off('click');
      element.bind('click', function () {
        if (window.confirm(msg)) {
          $log.debug('confirmed');
          scope.$eval(clickAction);
        }
      });
    }
  };
}]);

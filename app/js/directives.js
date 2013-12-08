'use strict';

/* Directives */

var angular = angular || {}; // To shut JSHint
angular.module('mustacheApp.directives', []).
    directive('appVersion', ['version', function (version) {
        return function (scope, elm, attrs) {
            elm.text(version);
        };
    }]);

var directives = angular.module('mustacheApp.directives');


directives.directive('accessLevel', ['$rootScope', 'Auth', function ($rootScope, Auth) {

    return {
        restrict: 'A',
        link: function (scope, element, attrs) {
            var prevDisp = element.css('display');
            $rootScope.$watch('User.role', function (role) {
                if (!Auth.authorize(attrs.accessLevel)) {
                    element.css('display', 'none');
                } else {
                    element.css('display', prevDisp);
                }
            });
        }
    };
}]);


directives.directive('integer', function () {
    return {
        require: 'ngModel',
        link: function (scope, ele, attr, ctrl) {
            ctrl.$parsers.unshift(function (viewValue) {
                return parseInt(viewValue);
            });
        }
    };
});

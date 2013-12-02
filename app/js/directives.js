'use strict';

/* Directives */

var angular = angular || {}; // To shut JSHint
angular.module('mustacheApp.directives', []).
  directive('appVersion', ['version', function(version) {
    return function(scope, elm, attrs) {
      elm.text(version);
    };
  }]);

angular.module('mustacheApp.directives')
	.directive('accessLevel', ['$rootScope', 'Auth', function($rootScope, Auth) {
	
	return { 
		restrict: 'A',
		link: function(scope, element, attrs) {
			var prevDisp = element.css('display');
			$rootScope.$watch('user.role', function(role) {
				if(!Auth.authorize(attrs.accessLevel))
					element.css('display', 'none');
				else
					element.css('display', prevDisp);
			});
		}
	};
}]);

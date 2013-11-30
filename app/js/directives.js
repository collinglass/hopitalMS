'use strict';

/* Directives */

var angular = angular || {}; // To shut JSHint
angular.module('mustacheApp.directives', []).
  directive('appVersion', ['version', function(version) {
    return function(scope, elm, attrs) {
      elm.text(version);
    };
  }]);

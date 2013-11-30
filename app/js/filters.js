'use strict';

/* Filters */
var angular = angular || {}; // To shut JSHint
var mustacheFilters = angular.module('mustacheApp.filters', []);

mustacheFilters.
  filter('interpolate', ['version', function(version) {
    return function(text) {
      return String(text).replace(/\%VERSION\%/mg, version);
    };
  }]);


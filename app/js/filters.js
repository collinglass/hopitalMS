'use strict';

/* Filters */

var mustacheFilters = angular.module('mustacheApp.filters', []);

mustacheFilters.
  filter('interpolate', ['version', function(version) {
    return function(text) {
      return String(text).replace(/\%VERSION\%/mg, version);
    };
  }]);


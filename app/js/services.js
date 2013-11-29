'use strict';

/* Services */


var mustacheServices = angular.module('mustacheApp.services', ['ngResource']);

mustacheServices.factory("Ward", ["$resource", function($resource){
    return $resource('/api/v0.1/wards/:wardId', {wardId:'@id'}, {
        query: {method:'GET', params:{wardId:'@id'}, isArray:true}
    });
}]);
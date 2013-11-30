'use strict';

/* Services */


var mustacheServices = angular.module('mustacheApp.services', ['ngResource']);

mustacheServices.factory("Ward", ["$resource", function($resource){
    return $resource('/api/v0.1/wards/:wardId', {wardId:'@id'}, {
        query: {method:'GET', params:{wardId:'@id'}, isArray:true}
    });
}]);

mustacheServices.factory("Patient", ["$resource", function($resource){
    return $resource('/api/v0.1/patients/:patientId', {patientId:'@id'}, {
        query: {method:'GET', params:{patientId:'@id'}, isArray:true}
    });
}]);

mustacheServices.factory("NextOfKin", ["$resource", function($resource){
    return $resource('/api/v0.1/nextOfKins/:nokId', {nokId:'@id'}, {
        query: {method:'GET', params:{nokId:'@id'}, isArray:true}
    });
}]);

mustacheServices.factory("User", ["$resource", function($resource){
    return $resource('/api/v0.1/users/:userId', {userId:'@id'}, {
        query: {method:'GET', params:{userId:'@id'}, isArray:true}
    });
}]);
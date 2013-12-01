'use strict';

/* Services */

var angular = angular || {}; // To shut JSHint
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

mustacheServices.factory("Role", ["$resource", function($resource){
    return $resource('/api/v0.1/roles/:roleId', {roleId:'@id'}, {
        query: {method:'GET', params:{roleId:'@id'}, isArray:true}
    });
}]);


mustacheServices.factory("Employee", ["$resource", function($resource){
    return $resource('/api/v0.1/employees/:employeeId', {employeeId:'@id'}, {
        query: {method:'GET', params:{employeeId:'@id'}, isArray:true}
    });
}]);
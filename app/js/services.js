'use strict';

/* Services */

var angular = angular || {}; // To shut JSHint
var mustacheServices = angular.module('mustacheApp.services', ['ngResource']);

mustacheServices.factory("Ward", ["$resource", function ($resource) {
    return $resource('/api/v0.1/wards/:wardId', {wardId: '@id'}, {
        query: {method: 'GET', params: {wardId: '@id'}, isArray: true},
        save: {method: 'POST', params: {wardId: '@id'}}
    });
}]);


mustacheServices.factory("Patient", ["$resource", function ($resource) {
    return $resource('/api/v0.1/patients/:patientId', {patientId: '@id'}, {
        query: {method: 'GET', params: {patientId: '@id'}, isArray: true},
        save: {method: 'POST', params: {patientId: '@id'}},
        delete: {method: 'DELETE', params: {patientId: '@id'}}
    });
}]);

mustacheServices.factory("AdmissionRequest", ["$resource", function ($resource) {
    return $resource('/api/v0.1/admissionRequests/:admRequestId', {admRequestId: '@id'}, {
        query: {method: 'GET', params: {admRequestId: '@id'}, isArray: true}
    });
}]);

mustacheServices.factory("Employee", ["$resource", function ($resource) {
    return $resource('/api/v0.1/employees/:employeeId', {employeeId: '@id'}, {
        query: {method: 'GET', params: {employeeId: '@id'}, isArray: true}
    });
}]);

mustacheServices.factory('Auth', ["$http", "$rootScope", "Employee", function ($http, $rootScope, Employee) {

    var isLogged = false;

    return {
        logIn: function (employeeId, password, success, error) {
            $http.post('/api/v0.1/sessions', {employeeId: employeeId, password: password})
                .success(function (data, status, headers, config) {
                    Employee.get({employeeId: employeeId}, function (empl) {
                        $rootScope.User = empl;
                        success(data, status, headers, config);
                        isLogged = true;
                    });
                })
                .error(function (data, status, headers, config) {
                    error(data, status, headers, config);
                });
        },
        logOut: function () {
            var promise = $http.delete('/api/v0.1/sessions');
            promise.success(function () {
                $rootScope.User = undefined;
                isLogged = false;
            });

            promise.error(function (data, status) {
                window.console.log("Status:" + status + JSON.stringify(data));
            });
        },
        isLogged: function () {
            return isLogged;
        },
        getUser: function () {
            return $rootScope.User;
        },
        authorize: function (requiredRoles) {
            var user = $rootScope.User || {};
            user.roles = user.roles || {};

            // For each required role
            for (var role in requiredRoles) {
                // If the required role is set to false,
                if (!requiredRoles[role]) {
                    // skip to next role
                    continue;
                }
                // this role is mandatory!
                if (!user.roles[role]) {
                    // But this user doesn't have it
                    return false;
                }
                // User has this role
            }
            return true;
        }
    };
}]);


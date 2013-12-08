'use strict';

/* Services */

var angular = angular || {}; // To shut JSHint
var mustacheServices = angular.module('mustacheApp.services', ['ngResource']);

mustacheServices.factory("Ward", ["$resource", function ($resource) {
    return $resource('/api/v0.1/wards/:wardId', {wardId: '@id'}, {
        query: {method: 'GET', params: {wardId: '@id'}, isArray: true}
    });
}]);

mustacheServices.factory("Patient", ["$resource", function ($resource) {
    return $resource('/api/v0.1/patients/:patientId', {patientId: '@id'}, {
        query: {method: 'GET', params: {patientId: '@id'}, isArray: true}
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

    return {
        logIn: function (employeeId, password) {
            return $http.post('/api/v0.1/sessions', {employeeId: employeeId, password: password});
        },
        logOut: function () {
            var promise = $http.delete('/api/v0.1/sessions');
            promise.success(function () {
                $rootScope.User = undefined;
            });

            promise.error(function (data, status) {
                window.console.log("Status:" + status + JSON.stringify(data));
            });
        },
        isLogged: function () {
            return $rootScope.User !== undefined;
        },
        getUser: function () {
            return $rootScope.User;
        },
        authorize: function (requiredRoles) {
            if (!this.isLogged()) {
                return false;
            }
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


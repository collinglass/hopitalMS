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

mustacheServices.factory("Role", ["$resource", function ($resource) {
    return $resource('/api/v0.1/roles/:roleId', {roleId: '@id'}, {
        query: {method: 'GET', params: {roleId: '@id'}, isArray: true}
    });
}]);


mustacheServices.factory("Employee", ["$resource", function ($resource) {
    return $resource('/api/v0.1/employees/:employeeId', {employeeId: '@id'}, {
        query: {method: 'GET', params: {employeeId: '@id'}, isArray: true}
    });
}]);

mustacheServices.factory('Auth', function ($http, $rootScope, $cookieStore) {

    var currentUser = $cookieStore.get('User') || { username: "", roles: ["public"] };
    
    // *** Start of Dummy $rootScope data to make app work without backend

    $rootScope.User = { username: "", roles: {
            "public":true,
            "chargeNurse": false,
            "doctor": false,
            "medicalStaff": true
        } 
    };
    //    console.log($rootScope.User);

    // *** End

    return {
        authorize: function (accessRoles, roles) {
            if (roles === undefined) {
                roles = $rootScope.User.roles;
            }
            for(var requiredRole in accessRoles) {
                if (roles[requiredRole]) {
                    return true;
                }
            }
            return false;
        },

        isLoggedIn: function (User) {
            if (User === undefined) {
                User = $rootScope.User;
            }
            for( var i = 0; i < 3; i++ ) {
                if ( User.roles["medicalStaff"] || User.roles["doctor"] || User.roles["chargeNurse"] ) { // TODO This is not correct
                    return true;
                }
            }

            return false;
        },

        register: function (User, success, error) {
            $http.post('/register', User).success(success).error(error);
        },

        login: function (User, success, error) {
            $http.post('/login', User).success(function (User) {
                $rootScope.User = User;
                success(User);
            }).error(error);
        },

        logout: function (success, error) {
            $http.post('/logout').success(function () {
                $rootScope.User = {
                    username: "",
                    roles: ["public"]
                };
                success();
            }).error(error);
        },

        User: currentUser
    };
});


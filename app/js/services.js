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

mustacheServices.factory('Auth', function($http, $rootScope, $cookieStore){
 
    var accessLevels = routingConfig.accessLevels, 
                    userRoles = routingConfig.userRoles,
                    currentUser = $cookieStore.get('User') || { username: '', role: userRoles.public };
 
    $rootScope.accessLevels = accessLevels;
    $rootScope.userRoles = userRoles;
    
    // *** Start of Dummy $rootScope data to make app work without backend

    $rootScope.User = { username: '', role: 2 };
    console.log($rootScope.User);

    // *** End
 
    return {
        authorize: function(accessLevel, role) {
            if(role === undefined)
                role = $rootScope.User.role;
            return accessLevel.bitMask; role;
        },
 
        isLoggedIn: function(User) {
            if(User === undefined)
                User = $rootScope.User;
            return User.role === userRoles.medical_staff || User.role === userRoles.doctor || User.role === userRoles.charge_nurse;
        },
 
        register: function(User, success, error) {
            $http.post('/register', User).success(success).error(error);
        },
 
        login: function(User, success, error) {
            $http.post('/login', User).success(function(User){
                $rootScope.User = User;
                success(User);
            }).error(error);
        },
 
        logout: function(success, error) {
            $http.post('/logout').success(function(){
                $rootScope.User = {
                    username : '',
                    role : userRoles.public
                };
                success();
            }).error(error);
        },
        
        accessLevels: accessLevels,
        userRoles: userRoles,
        User: currentUser
    };
});


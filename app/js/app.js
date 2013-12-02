'use strict';


// Declare app level module which depends on filters, and services
angular.module('mustacheApp', [
        'ngRoute',
        'ngCookies',
        'mustacheApp.filters',
        'mustacheApp.services',
        'mustacheApp.directives',
        'mustacheApp.controllers'
    ]).
    config(['$routeProvider', '$locationProvider', '$httpProvider', 
        function ($routeProvider, $locationProvider, $httpProvider) {
        var access = routingConfig.accessLevels;

        $routeProvider.when('/login', {
            templateUrl: 'partials/login.html',
            controller: 'LoginCtrl',
            access: access.anon
        });
        $routeProvider.when('/register', {
            templateUrl: 'partials/register.html',
            controller: 'RegisterCtrl',
            access: access.anon
        });
        $routeProvider.when('/ward', {
            templateUrl: 'partials/ward_list.html',
            controller: 'WardListCtrl',
            access: access.medicalStaff
        });
        $routeProvider.when('/ward/:wardId', {
            templateUrl: 'partials/ward_detail.html',
            controller: 'WardDetailCtrl',
            access: access.medicalStaff
        });
        $routeProvider.when('/patients/new', {
            templateUrl: 'partials/patient.html',
            controller: 'PatientCtrl',
            access: access.medicalStaff
        });
        $routeProvider.when('/patients/:patientId', {
            templateUrl: 'partials/patient.html',
            controller: 'PatientCtrl',
            access: access.medicalStaff
        });
        $routeProvider.otherwise({
            redirectTo: '/login',
            access: access.anon
        });


        var interceptor = ['$location', '$q', function($location, $q) {
        function success(response) {
            return response;
        }
 
        function error(response) {
 
            if(response.status === 401) {
                $location.path('/login');
                return $q.reject(response);
            } else {
            return $q.reject(response);
            }
        }
 
        return function(promise) {
            return promise.then(success, error);
        }
    }];

    $httpProvider.responseInterceptors.push(interceptor);

}]).run(['$rootScope', '$location', 'Auth', function ($rootScope, $location, Auth) {
 
    $rootScope.$on("$routeChangeStart", function (event, next, current) {
        if (!Auth.authorize(next.access)) {
            if(Auth.isLoggedIn()) $location.path(next);                // TODO Works but brings up error
            else $location.path('/login');
        } // TODO else redirect if authorize failed
    });
 
}]);

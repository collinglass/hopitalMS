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

        $routeProvider.when('/login', {
            templateUrl: 'partials/login.html',
            controller: 'LoginCtrl',
            access: {
                    "public":true,
                    "chargeNurse": false,
                    "doctor": false,
                    "medicalStaff": false,
                }
        });
        $routeProvider.when('/register', {
            templateUrl: 'partials/register.html',
            controller: 'RegisterCtrl',
            access: {
                    "public":true,
                    "chargeNurse": false,
                    "doctor": false,
                    "medicalStaff": false,
                }
        });
        $routeProvider.when('/ward', {
            templateUrl: 'partials/ward_list.html',
            controller: 'WardListCtrl',
            access: {
                    "medicalStaff": true,
                }
        });
        $routeProvider.when('/ward/:wardId', {
            templateUrl: 'partials/ward_detail.html',
            controller: 'WardDetailCtrl',
            access: {
                    "medicalStaff": true,
                }
        });
        $routeProvider.when('/patients/new', {
            templateUrl: 'partials/patient.html',
            controller: 'PatientCtrl',
            access: {
                    "medicalStaff": true,
                }
        });
        $routeProvider.when('/patients/:patientId', {
            templateUrl: 'partials/patient.html',
            controller: 'PatientCtrl',
            access: {
                    "medicalStaff": true,
                }
        });
        $routeProvider.when('/admissions/:admRequestId', {
            templateUrl: 'partials/patient.html',
            controller: 'PatientCtrl',
            access: {
                    "medicalStaff": true,
                    "chargeNurse": true,
                }
        });
        $routeProvider.when('/rationale/:admRequestId', {
            templateUrl: 'partials/rationale.html',
            controller: 'RationaleCtrl',
            access: {
                    "medicalStaff": true,
                }
        });
        $routeProvider.otherwise({
            redirectTo: '/login',
            access: {
                    "public":true,
                    "chargeNurse": false,
                    "doctor": false,
                    "medicalStaff": false,
                }
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
        if (Auth.authorize(next.access)) {
            $location.path(next.$route);
        } else if (!Auth.isLoggedIn()) {
            $location.path('/login');
        //} else if ( next.$route === '/login') {       // TODO if loggedin and trying to access login page -> redirect
        //    $location.path('/ward');
        } else {
            // Do nothing, user doesn't have access to this location
            console.log("Tried to access " + next + " from " + current + " but unauthorized.")
        }
    });
 
}]);

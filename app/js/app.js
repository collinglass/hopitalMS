'use strict';


// Declare app level module which depends on filters, and services
angular.module('mustacheApp', [
        'ngRoute',
        'mustacheApp.filters',
        'mustacheApp.services',
        'mustacheApp.directives',
        'mustacheApp.controllers'
    ]).
    config(['$routeProvider', function ($routeProvider) {
        $routeProvider.when('/login', {
            templateUrl: 'partials/login.html',
            controller: 'LoginCtrl'
        });
        $routeProvider.when('/register', {
            templateUrl: 'partials/register.html',
            controller: 'RegisterCtrl'
        });
        $routeProvider.when('/ward', {
            templateUrl: 'partials/ward_list.html',
            controller: 'WardListCtrl'
        });
        $routeProvider.when('/ward/:wardId', {
            templateUrl: 'partials/ward_detail.html',
            controller: 'WardDetailCtrl'
        });
        $routeProvider.otherwise({redirectTo: '/login'});
    }]);

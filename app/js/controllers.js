'use strict';

/* Controllers */

var controllers = angular.module('mustacheApp.controllers', []);

controllers.controller('LoginCtrl', [function () {

}]);

controllers.controller('RegisterCtrl', [function () {

}]);

controllers.controller('WardListCtrl', ["$scope", "Ward", function ($scope, Ward) {
        $scope.wards = Ward.query();

    }]).controller('WardDetailCtrl', ["$scope", "$routeParams", "Ward", function ($scope, $routeParams, Ward) {
        $scope.ward = Ward.get({wardId: $routeParams.wardId});
    }]);

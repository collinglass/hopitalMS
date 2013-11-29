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
        Ward.get({wardId: $routeParams.wardId}, function(ward){
            $scope.ward = ward;
            $scope.patients = ward.patients;
            $scope.admissions = ward.admissions;
        });

        $scope.admitPatient = function() {

            angular.forEach($scope.admissions, function(obj) {
                if ( obj.selected === true ) {
                    $scope.patients.push({lastName: obj.lastName, firstName: obj.firstName, healthInsNum: obj.healthInsNum, roomNum: "00", bedNum: "00", status: "nominal",});
                    $scope.admissions.splice(obj, 1);
                }
            });
        };
    }]);



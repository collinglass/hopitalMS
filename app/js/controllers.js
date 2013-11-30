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
        Ward.get({wardId: $routeParams.wardId}, function (ward) {
            $scope.ward = ward;
            $scope.patients = ward.patients;
            $scope.admissions = ward.admissions;

            $scope.admisions.admit = function () {

                angular.forEach($scope.admissions, function (patient) {
                    if (!patient.selected) {
                        // ignore those that aren't selected
                        return;
                    }
                    $scope.patients.push({
                        lastName: patient.lastName,
                        firstName: patient.firstName,
                        healthInsNum: patient.healthInsNum,
                        roomNum: "00",
                        bedNum: "00",
                        status: "nominal"
                    });

                    $scope.admissions.splice(patient, 1);
                });
            };

            $scope.patients.discharge = function () {
                angular.forEach($scope.patients, function (patient) {
                    if (!patient.selected) {
                        // ignore those that aren't selected
                        return;
                    }
                    var index = $scope.patients.indexOf(patient);
                    $scope.patients.splice(index, 1);
                });
            };

        });


    }]);



'use strict';

/* Controllers */

var angular = angular || {}; // To shut JSHint
var controllers = angular.module('mustacheApp.controllers', []);

controllers.controller('LoginCtrl', [function () {

}]);

controllers.controller('RegisterCtrl', [function () {

}]);

controllers.controller('WardListCtrl', ["$scope", "Ward", function ($scope, Ward) {

        Ward.query(function(wards) {
            $scope.wards = wards;
        });


    }]).controller('WardDetailCtrl', ["$scope", "$routeParams", "Ward", "Patient", function ($scope, $routeParams, Ward, Patient) {

        Ward.get({wardId: $routeParams.wardId}, function (ward) {
            $scope.ward = ward;
            $scope.patients = ward.patients;
            $scope.admissionsRequest = ward.admissionsRequest;
            $scope.admissionsResponse = ward.admissionsResponse;



            $scope.patients.forEach(function(patient){
                Patient.get({patientId:patient.patientId}, function(patientDetails){
                    patient.details = patientDetails;

                });
            });

            $scope.admissionsRequest.forEach(function(request){
                Patient.get({patientId:request.patientId}, function(patientDetails){
                    request.patientDetails = patientDetails;

                });
            });

            $scope.admissionsResponse.forEach(function(response){
                Patient.get({patientId:response.patientId}, function(patientDetails){
                    response.patientDetails = patientDetails;

                });
            });


            $scope.admissionsRequest.admit = function () {

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



'use strict';

/* Controllers */

var angular = angular || {}; // To shut JSHint
var controllers = angular.module('mustacheApp.controllers', []);

controllers.controller('LoginCtrl', [function () {

}]);

controllers.controller('RegisterCtrl', [function () {

}]);

controllers.controller('WardListCtrl', ["$scope", "$location", "Ward", "Employee",
    function ($scope, $location, Ward, Employee) {

        $scope.wards = [];
        Ward.query(function (wardIds) {
            wardIds.forEach(function (wardId) {
                Ward.get({wardId: wardId.id}, function (ward) {

                    Employee.get({employeeId: ward.doctorId}, function(doctor) {
                        ward.doctor = doctor;
                    });

                    Employee.get({employeeId: ward.chargeNurseId}, function(nurse) {
                        ward.chargeNurse = nurse;
                    });

                    $scope.wards.push(ward);
                });
            });
        });
    }]);

controllers.controller('PatientCtrl', ["$scope", "$location", "$routeParams", "Ward", "Patient", "Employee",
    function ($scope, $location, $routeParams, Ward, Patient, Employee) {
    $scope.hello = "Hello from PatientCtrl";

    $scope.go = function (path) {
            $location.path(path);
    };

    Patient.get({patientId: $routeParams.patientId}, function (patient) {
        $scope.patient = patient;
        $scope.patientId = patient.patientId;
        $scope.lastName = patient.lastName;
        $scope.firstName = patient.firstName;
        $scope.healthInsNum = patient.healthInsNum;
        $scope.address = patient.address;
        $scope.phoneNum = patient.phoneNum;
        $scope.dateOfBirth = patient.dateOfBirth;
        $scope.gender = patient.gender;
        $scope.maritalStatus = patient.maritalStatus;

        $scope.nextOfKin.name = patient.nextOfKin.name;
        $scope.nextOfKin.relationship = patient.nextOfKin.relationship;
        $scope.nextOfKin.address = patient.nextOfKin.address;
        $scope.nextOfKin.phoneNum = patient.nextOfKin.phoneNum;


    });


}]);

controllers.controller('WardDetailCtrl', ["$scope", "$location", "$routeParams", "Ward", "Patient", "Employee",
    function ($scope, $location, $routeParams, Ward, Patient, Employee) {

        $scope.go = function (path) {
            $location.path(path);
        };

        Ward.get({wardId: $routeParams.wardId}, function (ward) {
            $scope.ward = ward;
            $scope.patients = ward.patients;
            $scope.admissionsRequest = ward.admissionsRequest;
            $scope.admissionsResponse = ward.admissionsResponse;

            $scope.patients.forEach(function (patient) {
                Patient.get({patientId: patient.patientId}, function (patientDetails) {
                    patient.details = patientDetails;

                });
                ward.beds.forEach(function (bed) {
                    if (bed.bedId === patient.bedId) {
                        patient.roomId = bed.roomId;
                    }
                });
            });

            $scope.admissionsRequest.forEach(function (request) {
                Patient.get({patientId: request.patientId}, function (patientDetails) {
                    request.patientDetails = patientDetails;

                });

                Ward.get({wardId: request.fromWardId}, function (fromWard) {
                    request.fromWard = fromWard;
                    Employee.get({employeeId: fromWard.chargeNurseId}, function (nurse) {
                        request.chargeNurseName = nurse.firstName + " " + nurse.lastName;
                    });
                });
            });

            $scope.admissionsResponse.forEach(function (response) {
                Patient.get({patientId: response.patientId}, function (patientDetails) {
                    response.patientDetails = patientDetails;
                });
            });

            $scope.admissionsRequest.admit = function () {
                window.console.log("I AINT ADMITIN NOTHIN");
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



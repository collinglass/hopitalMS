'use strict';

/* Controllers */

var angular = angular || {}; // To shut JSHint
var controllers = angular.module('mustacheApp.controllers', []);

controllers.controller('LoginCtrl', ["$scope", /*"$rootScope",*/ "$http", "$location", function ($scope, $rootScope, $http, $location) {
    $scope.onNewAccount = function () {
        window.console.log("User is: ", JSON.stringify($scope.User));
        $rootScope.User = $scope.User;
        $location.path("/register/");
    };

    $scope.onLogin = function () {
        var promise = $http.post("/api/v0.1/sessions/", {
            employeeId: $rootScope.User.employeeId,
            password: $rootScope.User.password
        });

        promise.success(function (data, status) {
            window.console.log("Login success:", status);
            $rootScope.User.tokenName = data.tokenName;
            $rootScope.User.token = data.token;
            $location.path("/wards/");
        });

        promise.error(function (data, status) {
            window.console.log("Login error:", status);
        });
    };
}]);

controllers.controller('NavCtlr', ["$scope", "$rootScope", function ($scope, $rootScope) {

    $scope.updateRole = function () {
        var role = $rootScope.User.role;

        if ($scope.publicRole) {
            role |= 1;
        } else {
            role &= (~1);
        }

        if ($scope.medStaffRole) {
            role |= 2;
        } else {
            role &= (~2);
        }

        if ($scope.doctorRole) {
            role |= 4;
        } else {
            role &= (~4);
        }

        if ($scope.nurseRole) {
            role |= 8;
        } else {
            role &= (~8);
        }

        $rootScope.User.role = role;
        console.log($rootScope.User);
    };
}]);

controllers.controller('RegisterCtrl', ["$scope", "$rootScope", "$http", "$location", "Employee", function ($scope, $rootScope, $http, $location, Employee) {

    $scope.User = $rootScope.User;
    window.console.log("RegisterCtrl, User is: ", JSON.stringify($scope.User));

    $scope.onRegister = function () {
        window.console.log("User is: ", JSON.stringify($scope.User));
        $rootScope.User = $scope.User;
        Employee.save($scope.User, function () {
            var promise = $http.post("/api/v0.1/sessions/", {
                employeeId: $rootScope.User.employeeId,
                password: $rootScope.User.password
            });

            promise.success(function (data, status) {
                window.console.log("Login success:", status);
                $rootScope.User.tokenName = data.tokenName;
                $rootScope.User.token = data.token;
                $location.path("/wards/");
            });

            promise.error(function (data, status) {
                window.console.log("Login error:", status);
            });
        });
    };
}]);

controllers.controller('WardListCtrl', ["$scope", "$location", "Ward", "Employee",
    function ($scope, $location, Ward, Employee) {

        $scope.wards = [];
        Ward.query(function (wardIds) {
            wardIds.forEach(function (wardId) {
                Ward.get({wardId: wardId.id}, function (ward) {

                    Employee.get({employeeId: ward.doctorId}, function (doctor) {
                        ward.doctor = doctor;
                    });

                    Employee.get({employeeId: ward.chargeNurseId}, function (nurse) {
                        ward.chargeNurse = nurse;
                    });

                    $scope.wards.push(ward);
                });
            });
        });
    }]);


controllers.controller('WardDetailCtrl', ["$scope", "$location", "$routeParams", "Ward", "Patient", "Employee", "Auth",
    function ($scope, $location, $routeParams, Ward, Patient, Employee, Auth) {

        $scope.authorize = function (accessRoles) {
            return Auth.authorize(accessRoles);
        }

        $scope.go = function (path) {
            $location.path(path);
        };

        Ward.get({wardId: $routeParams.wardId}, function (ward) {
            $scope.ward = ward;
            $scope.patients = ward.patients;
            $scope.admissionRequests = ward.admissionRequests;
            $scope.admissionResponses = ward.admissionResponses;

            console.log($scope.admissionRequests);

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

            $scope.admissionRequests.forEach(function (request) {
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

            $scope.admissionResponses.forEach(function (response) {
                Patient.get({patientId: response.patientId}, function (patientDetails) {
                    response.patientDetails = patientDetails;
                });
            });

            $scope.patients.view = function () {
                angular.forEach($scope.patients, function(obj) {
                    if ( obj.selected == true ) {
                        $scope.go("/patients/" + obj.patientId);
                    }
                });
            };

            $scope.admissionRequests.view = function () {
                angular.forEach($scope.admissionRequests, function(obj) {
                    if ( obj.selected == true ) {
                        console.log(obj); // TODO remove
                        $scope.go("/rationale/" + obj.admRequestId);
                    }
                });
            };

            $scope.admissionRequests.admit = function () {
                angular.forEach($scope.admissionRequests, function(obj) {
                    if ( obj.selected == true ) {
                        var patientPush = {
                            details: obj.patientDetails, roomId: "00", 
                            bedId: "00", status: "nominal" };
                            $scope.patients.push(patientPush);

                        $scope.admissionRequests.splice(obj, 1);
                        $scope.ward.save();                         // TODO test if save function works

                        obj.fromWard.patients.splice(obj, 1);
                        obj.fromWard.save();                        // TODO test if save function works
                    }
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

controllers.controller('PatientCtrl', ["$scope", "$location", "$routeParams", "Ward", "Patient",
    function ($scope, $location, $routeParams, Ward, Patient) {

        $scope.go = function (path) {
            $location.path(path);
        };

        $scope.save = function (patient) {
            Patient.save(patient);
        };

        $scope.newPath = function () {
            if ($location.path() === '/patients/new') {
                return true;
            }
            return false;
        };

        if ($location.path() === '/patients/new') {

            // Test with Ward 1
            Ward.get({wardId: 1}, function (ward) {  // TODO dynamic wardId
                $scope.ward = ward;
                $scope.patients = ward.patients;
            });

            window.console.log($scope.patients);

            $scope.admit = function () {
                window.console.log($scope.patients);
                var patientPush = { patientId: 205, lastName: $scope.lastName, firstName: $scope.firstName,
                    healthInsNum: $scope.healthInsNum, address: $scope.address, phoneNum: $scope.phoneNum,
                    dateOfBirth: $scope.dateOfBirth, gender: $scope.gender, maritalStatus: $scope.maritalStatus,
                    nextOfKin: { name: $scope.nextOfKin.name, relationship: $scope.nextOfKin.relationship,
                        address: $scope.nextOfKin.address, phoneNum: $scope.nextOfKin.phoneNum }
                    };

                    $scope.patients.push(patientPush);

                // Check out the ward after
                $scope.go('/ward/1');               // TODO dynamic wardId

            };

        } else {
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

                $scope.nextOfKin = patient.nextOfKin;

                $scope.update = function () {
                    window.console.log('update');
                    /*
                     $scope.patients.patient({ patientId: $scope.patientId, lastName: $scope.lastName, firstName: $scope.firstName,
                     healthInsNum: $scope.healthInsNum, address: $scope.address, phoneNum: $scope.phoneNum,
                     dateOfBirth: $scope.dateOfBirth, gender: $scope.gender, maritalStatus: $scope.maritalStatus,
                     nextOfKin: { name: $scope.nextOfKin.name, relationship: $scope.nextOfKin.relationship,
                     address: $scope.nextOfKin.address, phoneNum: $scope.nextOfKin.phoneNum }});
*/
};

});
}
}]);

controllers.controller('AdmissionCtrl', ["$scope", "$location", "$routeParams", "Ward", "Patient", "AdmissionRequest",
    function ($scope, $location, $routeParams, Ward, Patient, AdmissionRequest) {

        $scope.go = function (path) {
            $location.path(path);
        };

        $scope.save = function (patient) {
            Patient.save(patient);
        };



        AdmissionRequest.get({admRequestId: $routeParams.admRequestId}, function (admissionRequest) {

            console.log(admissionRequest);
            $scope.admissionRequest = admissionRequest;
            $scope.admRequestId = admissionRequest.admRequestId;
            $scope.patientId = admissionRequest.patientId;
            $scope.rationale = admissionRequest.rationale;
            $scope.priority = admissionRequest.priority;
            //$scope.lastName = admissionRequest.lastName;              // TODO using patientId get patient info
            //$scope.firstName = admissionRequest.firstName;
            //$scope.healthInsNum = admissionRequest.healthInsNum;      // TODO using fromWardId get ward outbound doctor, etc.
            //$scope.address = patient.address;
            //$scope.phoneNum = patient.phoneNum;
            //$scope.dateOfBirth = patient.dateOfBirth;
            //$scope.gender = patient.gender;
            //$scope.maritalStatus = patient.maritalStatus;


            $scope.update = function () {
                window.console.log('update');
                    /*
                     $scope.patients.patient({ patientId: $scope.patientId, lastName: $scope.lastName, firstName: $scope.firstName,
                     healthInsNum: $scope.healthInsNum, address: $scope.address, phoneNum: $scope.phoneNum,
                     dateOfBirth: $scope.dateOfBirth, gender: $scope.gender, maritalStatus: $scope.maritalStatus,
                     nextOfKin: { name: $scope.nextOfKin.name, relationship: $scope.nextOfKin.relationship,
                     address: $scope.nextOfKin.address, phoneNum: $scope.nextOfKin.phoneNum }});
                    */
            };

        });
}]);
controllers.controller('RationaleCtrl', ["$scope", "$location", "$routeParams", "Ward", "Patient", "AdmissionRequest",
    function ($scope, $location, $routeParams, Ward, Patient, AdmissionRequest) {

        $scope.go = function (path) {
            $location.path(path);
        };

        $scope.save = function (patient) {
            Patient.save(patient);
        };



        AdmissionRequest.get({admRequestId: $routeParams.admRequestId}, function (admissionRequest) {

            console.log(admissionRequest);
            $scope.admissionRequest = admissionRequest;
            $scope.admRequestId = admissionRequest.admRequestId;
            $scope.patientId = admissionRequest.patientId;
            $scope.rationale = admissionRequest.rationale;
            $scope.priority = admissionRequest.priority;
            //$scope.lastName = admissionRequest.lastName;              // TODO using patientId get patient info
            //$scope.firstName = admissionRequest.firstName;
            //$scope.healthInsNum = admissionRequest.healthInsNum;      // TODO using fromWardId get ward outbound doctor, etc.
            //$scope.address = patient.address;
            //$scope.phoneNum = patient.phoneNum;
            //$scope.dateOfBirth = patient.dateOfBirth;
            //$scope.gender = patient.gender;
            //$scope.maritalStatus = patient.maritalStatus;


            $scope.update = function () {
                window.console.log('update');
                    /*
                     $scope.patients.patient({ patientId: $scope.patientId, lastName: $scope.lastName, firstName: $scope.firstName,
                     healthInsNum: $scope.healthInsNum, address: $scope.address, phoneNum: $scope.phoneNum,
                     dateOfBirth: $scope.dateOfBirth, gender: $scope.gender, maritalStatus: $scope.maritalStatus,
                     nextOfKin: { name: $scope.nextOfKin.name, relationship: $scope.nextOfKin.relationship,
                     address: $scope.nextOfKin.address, phoneNum: $scope.nextOfKin.phoneNum }});
                    */
            };

        });
}]);
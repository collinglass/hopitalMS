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

controllers.controller('NavCtlr', ["$scope", "$rootScope", "Auth", function ($scope, $rootScope, Auth) {

    $scope.isLoggedIn = function () {
        return Auth.isLoggedIn();
    }

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


controllers.controller('WardDetailCtrl', ["$scope", "$location", "$routeParams", "Ward", "Patient", "Employee", "Auth", "$rootScope",
    function ($scope, $location, $routeParams, Ward, Patient, Employee, Auth, $rootScope) {

        Ward.get({wardId: $routeParams.wardId}, function (ward) {
            $scope.ward = ward;
            $scope.patients = ward.patients;
            $scope.admissionRequests = ward.admissionRequests;
            $scope.admissionResponses = ward.admissionResponses;

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
            

            if ( $scope.ward.wardId === $rootScope.User.wardId ) {

                $scope.authorize = function (accessRoles) {
                    return Auth.authorize(accessRoles);
                };

                $scope.go = function (path) {
                    $location.path(path);
                };

                $scope.patients.view = function () {
                    angular.forEach($scope.patients, function(obj) {
                        if ( obj.selected == true ) {
                            $scope.go("/patients/" + obj.patientId);
                        }
                    });
                };

                $scope.patients.transfer = function () {
                    angular.forEach($scope.patients, function(obj) {
                        if ( obj.selected == true ) {
                            $scope.go("/transfer/" + obj.patientId);
                        }
                    });
                };

                $scope.admissionRequests.refuse = function () {
                    angular.forEach($scope.admissionRequests, function(obj) {
                        if ( obj.selected == true ) {
                            $scope.go("/refusal/" + obj.admRequestId);
                        }
                    });
                };

                $scope.admissionRequests.view = function () {
                    angular.forEach($scope.admissionRequests, function(obj) {
                        if ( obj.selected == true ) {
                            $scope.go("/rationale/" + obj.admRequestId);
                        }
                    });
                };

                $scope.admissionRequests.admit = function () {
                    angular.forEach($scope.admissionRequests, function(obj) {
                        if ( obj.selected ) {
                            console.log(obj);
                            $scope.go("/admissions/" + obj.admRequestId);
                        }
                    });
                };

                $scope.patients.discharge = function () {
                    angular.forEach($scope.patients, function (patient) {
                        if (patient.selected) {
                            var index = $scope.patients.indexOf(patient);
                            $scope.patients.splice(index, 1);                   // TODO Free up bed
                        }
                    });
                };
            };
        });
}]);

controllers.controller('PatientCtrl', ["$scope", "$location", "$routeParams", "Ward", "Patient", "AdmissionRequest",
    function ($scope, $location, $routeParams, Ward, Patient, AdmissionRequest) {

        $scope.go = function (path) {
            $location.path(path);
        };

        $scope.save = function () {
            Patient.save($scope.patient);
        };

        $scope.newPath = function () {
            if ($location.path() === '/patients/new') {
                return true;
            }
            return false;
        };
                                                                    // TODO dynamically show available rooms and beds
                                                                    if ($location.path() === '/patients/new') {
            // Test with Ward 1
            Ward.get({wardId: 1}, function (ward) {  // TODO dynamic wardId
                $scope.ward = ward;
                $scope.patients = ward.patients;
                $scope.beds = ward.beds;
            });
            Patient.query(function (patientIds) {
                $scope.newPatientId = patientIds.length;
            });

            //window.console.log($scope.patients);

            $scope.admit = function () {
                var newPatient = new Patient({ 
                    patientId: $scope.newPatientId, 
                    lastName: $scope.lastName, 
                    firstName: $scope.firstName,
                    healthInsNum: $scope.healthInsNum, 
                    address: $scope.address, 
                    phoneNum: $scope.phoneNum,
                    dateOfBirth: $scope.dateOfBirth, 
                    gender: $scope.gender, 
                    maritalStatus: $scope.maritalStatus,
                    nextOfKin: { 
                        name: $scope.nextOfKin.name, 
                        relationship: $scope.nextOfKin.relationship,
                        address: $scope.nextOfKin.address, 
                        phoneNum: $scope.nextOfKin.phoneNum }
                    });
                $scope.patient = newPatient;
                    $scope.save();                      // TODO get working with database
                    var wardPush = {
                        patientId: $scope.newPatientId,
                        bedId: $scope.bedId,
                        status: "nominal"
                    };
                    console.log(wardPush);

                    $scope.patients.push(wardPush);

                // Check out the ward after
                $scope.go('/ward/' + $scope.ward.wardId);

            };

        } else if ( $location.path() === ('/admissions/' + $routeParams.admRequestId) ) {
            console.log($routeParams);

            AdmissionRequest.get({admRequestId: $routeParams.admRequestId}, function (admissionRequest) {

                $scope.admissionRequest = admissionRequest;
                $scope.admRequestId = admissionRequest.admRequestId;
                $scope.patientId = admissionRequest.patientId;
                $scope.rationale = admissionRequest.rationale;
                $scope.priority = admissionRequest.priority;

                console.log(admissionRequest);


                Patient.get({patientId: admissionRequest.patientId}, function (patient) {  // TODO dynamic wardId
                    $scope.patient = patient;
                });
                Ward.get({wardId: 1}, function (ward) {  // TODO dynamic wardId
                    $scope.ward = ward;
                    $scope.patients = ward.patients;
                    $scope.beds = ward.beds;
                    $scope.admissionRequests = ward.admissionRequests;
                });
            });

            var wardPush = {
                patientId: $scope.patient.patientId,
                bedId: $scope.bedId,
                status: "nominal"
            };
            console.log(wardPush);
            $scope.patients.push(wardPush);
            var index = $scope.admissionRequests.indexOf(admissionRequest);
            $scope.admissionRequests.splice(index, 1);
            
        } else {
            Patient.get({patientId: $routeParams.patientId}, function (patient) {
                $scope.patient = patient;

                Ward.get({wardId: 1}, function (ward) {  // TODO dynamic wardId
                    $scope.ward = ward;
                    $scope.patients = ward.patients;
                    $scope.beds = ward.beds;
                });
            });
        };
    }]);

controllers.controller('RationaleCtrl', ["$scope", "$location", "$routeParams", "Ward", "Patient", "Employee", "AdmissionRequest",
    function ($scope, $location, $routeParams, Ward, Patient, Employee, AdmissionRequest) {
        $scope.back = function() {
            history.go(-1);
        };


        AdmissionRequest.get({admRequestId: $routeParams.admRequestId}, function (admissionRequest) {

            $scope.admissionRequest = admissionRequest;
            $scope.admRequestId = admissionRequest.admRequestId;
            $scope.patientId = admissionRequest.patientId;
            $scope.rationale = admissionRequest.rationale;
            $scope.priority = admissionRequest.priority;

            Patient.get({patientId: admissionRequest.patientId}, function (patient) {
                $scope.firstName = patient.firstName;
                $scope.lastName = patient.lastName;
            });
            

            Ward.get({wardId: admissionRequest.fromWardId}, function (fromWard) {
                $scope.wardName = fromWard.name;
                Employee.get({employeeId: fromWard.chargeNurseId}, function (nurse) {
                    $scope.chargeNurseName = nurse.firstName + " " + nurse.lastName;
                });
                Employee.get({employeeId: fromWard.doctorId}, function (doctor) {
                    $scope.doctorName = doctor.firstName + " " + doctor.lastName;
                });
            });

        });
}]);

controllers.controller('RefusalCtrl', ["$scope", "$location", "$routeParams", "Ward", "Patient", "Employee", "AdmissionRequest",
    function ($scope, $location, $routeParams, Ward, Patient, Employee, AdmissionRequest) {
        $scope.back = function() {
            history.go(-1);
        };

        $scope.refuse = function() {
            var response = {
                patientId: $scope.patientId,
                toWardId: $scope.fromWard.wardId,
                inProgress: false,
                refusal: $scope.refusal
            };
            console.log(response);
            $scope.fromWard.admissionResponses.push(response);
            var index = $scope.toWard.admissionRequests.indexOf($scope.admissionRequest);
            $scope.toWard.admissionRequests.splice(index, 1);
        }

        AdmissionRequest.get({admRequestId: $routeParams.admRequestId}, function (admissionRequest) {

            $scope.admissionRequest = admissionRequest;
            $scope.admRequestId = admissionRequest.admRequestId;
            $scope.patientId = admissionRequest.patientId;
            $scope.rationale = admissionRequest.rationale;
            $scope.priority = admissionRequest.priority;

            Patient.get({patientId: admissionRequest.patientId}, function (patient) {
                $scope.firstName = patient.firstName;
                $scope.lastName = patient.lastName;
            });
            
            Ward.get({wardId: admissionRequest.toWardId}, function (toWard) {           
                $scope.toWard = toWard;
            });

            Ward.get({wardId: admissionRequest.fromWardId}, function (fromWard) {
                $scope.wardName = fromWard.name;
                $scope.fromWard = fromWard;
                Employee.get({employeeId: fromWard.chargeNurseId}, function (nurse) {
                    $scope.chargeNurseName = nurse.firstName + " " + nurse.lastName;
                });
                Employee.get({employeeId: fromWard.doctorId}, function (doctor) {
                    $scope.doctorName = doctor.firstName + " " + doctor.lastName;
                });

            });

        });

}]);

controllers.controller('TransferCtrl', ["$scope", "$location", "$routeParams", "Ward", "Patient", "Employee", "AdmissionRequest",
    function ($scope, $location, $routeParams, Ward, Patient, Employee, AdmissionRequest) {
        $scope.back = function() {
            history.go(-1);
        };

        $scope.transfer = function() {
            var request = {
                admRequestId: $scope.newAdmissionId,
                patientId: $scope.patient.patientId,
                toWardId: $scope.toWardId,
                fromWardId: $scope.fromWardId,
                priority: $scope.priority,
                rationale: $scope.rationale
            };
            console.log(request);
            $scope.admissionRequests.push(request);
            // TODO save
        }

        Patient.get({patientId: $routeParams.patientId}, function (patient) {
            $scope.patient = patient;
            var toWards = [];
            Ward.query(function (wardIds) {
                $scope.toWards = wardIds;
                wardIds.forEach(function (wardId) {
                    Ward.get({wardId: wardId.id}, function (ward) {
                        ward.patients.forEach(function (wardPatient) {
                            if ( patient.patientId === wardPatient.patientId ) {
                                $scope.admissionRequests = ward.admissionRequests;
                                $scope.fromWardId = ward.wardId;
                            }
                        });
                    });
                });
            });
        });
    }]);


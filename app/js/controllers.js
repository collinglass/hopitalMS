'use strict';

/* Controllers */

var angular = angular || {}; // To shut JSHint
var controllers = angular.module('mustacheApp.controllers', []);

controllers.controller('LoginCtrl', ['$scope', '$rootScope', '$location', 'Auth', function ($scope, $rootScope, $location, Auth) {
    $scope.User = {};
    $scope.User.employeeId = '';
    $scope.User.password = '';

    $scope.isLogged = function () {
        return Auth.isLogged();
    };

    $scope.onNewAccount = function () {
        var employeeId = $scope.User.employeeId || '';
        var password = $scope.User.password || '';

        if (employeeId && password) {
            $rootScope.User = $scope.User;
            $location.path('/register');
        }
    };

    $scope.onLogin = function () {
        $scope.User = $scope.User || {};
        var employeeId = $scope.User.employeeId;
        var password = $scope.User.password;

        var success = function () {
            $location.path('/ward/' + Auth.getUser().wardId);
        };

        var error = function (data) {
            if (data.code !== 400) {
                $scope.errorMsg = data.message;
            }
            window.console.log('Status: ' + status + ', message: ' + angular.toJson(data));
        };

        Auth.logIn(employeeId, password, success, error);
    };
}]);

controllers.controller('RegisterCtrl', ['$scope', '$rootScope', '$location', 'Employee', 'Auth', function ($scope, $rootScope, $location, Employee, Auth) {

    var getRoleObject = function () {
        switch ($scope.RoleChoice) {
            case 'Medical Staff':
                return {medicalStaff: true};
            case 'Charge Nurse':
                return {medicalStaff: true, chargeNurse: true};
            case 'Doctor':
                return {medicalStaff: true, doctor: true};
        }
    };

    $scope.onRegister = function () {
        $scope.User.roles = getRoleObject();
        Employee.save($scope.User, function () {
            $scope.User = $scope.User || {};
            var employeeId = $scope.User.employeeId;
            var password = $scope.User.password;

            var success = function () {
                $location.path('/ward/' + Auth.getUser().wardId);
            };

            var error = function (data) {
                window.console.log('Status: ' + status + ', message: ' + angular.toJson(data));
            };
            Auth.logIn(employeeId, password, success, error);
        });
    };
}]);

controllers.controller('NavCtlr', ['$scope', '$location', 'Auth', function ($scope, $location, Auth) {

    $scope.isLogged = function () {
        return Auth.isLogged();
    };

    $scope.logOut = function () {
        Auth.logOut();
        $location.path('/login');
    };

    $scope.getUser = function () {
        return Auth.getUser();
    };
}]);

controllers.controller('WardListCtrl', ['$scope', '$location', 'Ward', 'Employee',
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


controllers.controller('WardDetailCtrl', ['$scope', '$location', '$routeParams', 'Ward', 'Patient', 'Employee', 'Auth', '$rootScope',
    function ($scope, $location, $routeParams, Ward, Patient, Employee, Auth, $rootScope) {

        Ward.get({wardId: $routeParams.wardId}, function (ward) {
            $scope.ward = ward || {};
            $scope.patients = ward.patients || [];
            $scope.admissionRequests = ward.admissionRequests || [];
            $scope.admissionResponses = ward.admissionResponses || [];

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
                        request.chargeNurseName = nurse.firstName + ' ' + nurse.lastName;
                    });
                });
            });

            $scope.admissionResponses.forEach(function (response) {
                Patient.get({patientId: response.patientId}, function (patientDetails) {
                    response.patientDetails = patientDetails;
                });
            });


            if ($scope.ward.wardId !== $rootScope.User.wardId) {
                // Stop here
                return;
            }

            $scope.authorize = function (accessRoles) {
                return Auth.authorize(accessRoles);
            };

            $scope.go = function (path) {
                $location.path(path);
            };

            $scope.patients.view = function () {
                angular.forEach($scope.patients, function (obj) {
                    if (obj.selected) {
                        $scope.go('/patients/' + obj.patientId);
                    }
                });
            };

            $scope.patients.transfer = function () {
                angular.forEach($scope.patients, function (obj) {
                    if (obj.selected) {
                        $scope.go('/transfer/' + obj.patientId);
                    }
                });
            };

            $scope.admissionRequests.refuse = function () {
                angular.forEach($scope.admissionRequests, function (obj) {
                    if (obj.selected) {
                        $scope.go('/refusal/' + obj.admRequestId);
                    }
                });
            };

            $scope.admissionRequests.view = function () {
                angular.forEach($scope.admissionRequests, function (obj) {
                    if (obj.selected) {
                        $scope.go('/rationale/' + obj.admRequestId);
                    }
                });
            };

            $scope.admissionRequests.admit = function () {
                angular.forEach($scope.admissionRequests, function (obj) {
                    if (obj.selected) {
                        //window.console.log(obj);
                        $scope.go('/admissions/' + obj.admRequestId);
                    }
                });
            };

            $scope.patients.discharge = function () {
                angular.forEach($scope.patients, function (patient) {
                    if (patient.selected) {
                        patient.details.$delete({patientId: patient.patientId});
                        var ward = $scope.ward;
                        var index = ward.patients.indexOf(patient);
                        ward.patients.splice(index, 1);
                        ward.$save({wardId: ward.wardId});
                    }
                });
            };
        });
    }]);

controllers.controller('PatientCtrl', ['$scope', '$location', '$routeParams', '$rootScope', 'Ward', 'Patient', 'Auth',
    function ($scope, $location, $routeParams, $rootScope, Ward, Patient, Auth) {

        var filterFreeBeds = function(ward) {
            var busyBeds = {};
            angular.forEach(ward.patients, function(inPatient) {
               busyBeds[inPatient.bedId] = true;
            });

            var freeBeds = [];
            angular.forEach(ward.beds, function(bed) {
                if (busyBeds[bed.bedId]) {
                    // skip this one
                    return;
                }
                freeBeds.push(bed);
            });
            return freeBeds;
        };

        $scope.go = function (path) {
            $location.path(path);
        };

        $scope.save = function () {
            var patient = $scope.patient;
            window.console.log('Saving patient: ' + angular.toJson(patient));
            patient.$save({patientId: patient.patientId});

            $scope.go('/ward/' + $rootScope.User.wardId);
        };

        $scope.newPath = function () {
            return $location.path() === '/patients/new';
        };
        $scope.admissionPath = function () {
            return $location.path() === '/admissions/' + $routeParams.admRequestId;

        };
        $scope.updatePath = function () {
            return $location.path() === '/patients/' + $routeParams.patientId;
        };

        // TODO dynamically show available rooms and beds
        if ($location.path() === '/patients/new') {
            var wardId = Auth.getUser().wardId;
            Ward.get({wardId: wardId}, function (ward) {
                $scope.ward = ward;
                $scope.patients = ward.patients;
                $scope.freeBeds = filterFreeBeds(ward);
            });
            Patient.query(function (patientIds) {
                $scope.newPatientId = patientIds.length;
            });


            $scope.admit = function () {

                var newPatient = $scope.patient;
                var patient = new Patient();
                patient.lastName = newPatient.lastName;
                patient.firstName = newPatient.firstName;
                patient.healthInsNum = newPatient.healthInsNum;
                patient.address = newPatient.address;
                patient.phoneNum = newPatient.phoneNum;
                patient.dateOfBirth = newPatient.dateOfBirth;
                patient.gender = newPatient.gender;
                patient.maritalStatus = newPatient.maritalStatus;
                patient.nextOfKin = newPatient.nextOfKin;

                patient.$save(function (savedPt) {
                    var bedId = $scope.bedId;
                    var ward = $scope.ward;
                    ward.patients.push({
                        patientId: savedPt.patientId,
                        bedId: bedId,
                        status: 'nominal'
                    });
                    ward.$save({wardId: ward.wardId});

                    // Check out the ward after
                    $scope.go('/ward/' + $scope.ward.wardId);
                });

            };

        } else if ($location.path() === ('/admissions/' + $routeParams.admRequestId)) {
            Ward.get({wardId: $rootScope.User.wardId}, function (ward) {
                window.console.log(angular.toJson(ward));
                $scope.ward = ward;
                ward.admissionRequests.forEach( function (admissionRequest) {
                    if( $routeParams.admRequestId == admissionRequest.admRequestId ) {
                        $scope.admissionRequest = admissionRequest;
                        Patient.get({patientId: $scope.admissionRequest.patientId}, function (patient) {
                            $scope.patient = patient;
                            $scope.freeBeds = filterFreeBeds(ward);
                        });
                        return;
                    }
                });
            });

            $scope.admit = function () {
                var wardPush = {
                    patientId: $scope.patient.patientId,
                    bedId: $scope.bedId,
                    status: 'nominal'
                };
                window.console.log(wardPush);
                $scope.ward.patients.push(wardPush);
                var index = $scope.ward.admissionRequests.indexOf($scope.admissionRequest);
                window.console.log('Splice at ' + index);
                $scope.ward.admissionRequests.splice(index, 1);
                $scope.ward.$save({wardId: $scope.ward.wardId});

                $scope.go('/ward/' + $scope.ward.wardId);
            };
        } else {
            Patient.get({patientId: $routeParams.patientId}, function (patient) {
                $scope.patient = patient;

                Ward.get({wardId: 1}, function (ward) {  // TODO dynamic wardId
                    $scope.ward = ward;
                    $scope.patients = ward.patients;
                    $scope.beds = ward.beds;
                });
            });
        }
    }]);

controllers.controller('RationaleCtrl', ['$scope', '$location', '$routeParams', '$rootScope', 'Ward', 'Patient', 'Employee',
    function ($scope, $location, $routeParams, $rootScope, Ward, Patient, Employee) {
        $scope.back = function () {
            history.go(-1);
        };

        Ward.get({wardId: $rootScope.User.wardId}, function (ward) {
            window.console.log(ward);
            $scope.ward = ward;
            Employee.get({employeeId: ward.chargeNurseId}, function (employee) {
                $scope.chargeNurseName = employee.firstName + ' ' + employee.lastName;
            });
            Employee.get({employeeId: ward.doctorId}, function (employee) {
                $scope.doctorName = employee.firstName + ' ' + employee.lastName;
            });
            ward.admissionRequests.forEach( function (admissionRequest) {
                if( $routeParams.admRequestId == admissionRequest.admRequestId ) {
                    $scope.admissionRequest = admissionRequest;
                    window.console.log(admissionRequest);
                    Ward.get({wardId: admissionRequest.fromWardId}, function (ward) {
                        $scope.fromWard = ward;
                    });
                    Patient.get({patientId: admissionRequest.patientId}, function (patient) {
                        $scope.patient = patient;
                    });
                    return;
                }
            });
        });
    }]);

controllers.controller('RefusalCtrl', ['$scope', '$location', '$routeParams', '$rootScope', 'Ward', 'Patient', 'Employee',
    function ($scope, $location, $routeParams, $rootScope, Ward, Patient, Employee) {
        $scope.go = function (path) {
            $location.path(path);
        };

        $scope.back = function () {
            history.go(-1);
        };

        Ward.get({wardId: $rootScope.User.wardId}, function (ward) {
            window.console.log(ward);
            $scope.ward = ward;
            Employee.get({employeeId: ward.chargeNurseId}, function (employee) {
                $scope.chargeNurseName = employee.firstName + ' ' + employee.lastName;
            });
            Employee.get({employeeId: ward.doctorId}, function (employee) {
                $scope.doctorName = employee.firstName + ' ' + employee.lastName;
            });
            ward.admissionRequests.forEach( function (admissionRequest) {
                if( $routeParams.admRequestId == admissionRequest.admRequestId ) {
                    $scope.admissionRequest = admissionRequest;
                    window.console.log(admissionRequest);
                    Ward.get({wardId: admissionRequest.fromWardId}, function (ward) {
                        $scope.fromWard = ward;
                    });
                    Patient.get({patientId: admissionRequest.patientId}, function (patient) {
                        $scope.patient = patient;
                    });
                }
            });
        });

        $scope.refuse = function () {
            console.log("Ready!");
            var response = {
                patientId: $scope.patient.patientId,
                toWardId: $scope.fromWard.wardId,
                inProgress: false,
                refusal: $scope.refusal
            };
            $scope.fromWard.admissionResponses.push(response);
            var index = $scope.ward.admissionRequests.indexOf($scope.admissionRequest);
            $scope.ward.admissionRequests.splice(index, 1);
            $scope.ward.$save({wardId: $scope.ward.wardId});

            $scope.go('/ward/' + $scope.ward.wardId);
        };

    }]);

controllers.controller('TransferCtrl', ['$scope', '$location', '$routeParams', '$rootScope', 'Ward', 'Patient',
    function ($scope, $location, $routeParams, $rootScope, Ward, Patient) {

        $scope.go = function (path) {
            $location.path(path);
        };

        $scope.back = function () {
            history.go(-1);
        };

        String.prototype.hashCode = function () {
            var hash = 0;
            if (this.length === 0) return hash;
            for (var i = 0; i < this.length; i++) {
                var char = this.charCodeAt(i);
                hash = ((hash<<5)-hash)+char;
                hash = hash & hash; // Convert to 32bit integer
            }
            return hash;
        };

        Patient.get({patientId: $routeParams.patientId}, function (patient) {
            var wardList = [];
            $scope.patient = patient;
            var newAdmissionId = $scope.patient.patientId + ' ' + new Date().getTime();

            Ward.query(function (wards) {
                window.console.log(wards);
                $scope.wards = wards;
                
                wards.forEach(function (ward) {
                    if ( $rootScope.User.wardId == ward.wardId ) {
                        $scope.ward = ward;
                    } else {
                        $scope.toWard = ward;
                        wardList.push(ward);
                    }
                    $scope.ward.patients.forEach( function (patient) {
                        if ( patient.patientId == $scope.patient.patientId ) {
                            patient.status = 'transfer';
                        }
                    });
                });

                $scope.wardList = wardList;

                $scope.transfer = function () {
                    var request = {
                        admRequestId: newAdmissionId.hashCode(),
                        patientId: $scope.patient.patientId,
                        toWardId: $scope.wardId,
                        fromWardId: $scope.User.wardId,
                        priority: $scope.priority,
                        rationale: $scope.rationale
                    };
                    
                    $scope.toWard.admissionRequests.push(request);
                    
                    $scope.toWard.$save({wardId: $scope.toWard.wardId});
                    $scope.ward.$save({wardId: $scope.ward.wardId});

                    $scope.go('/ward/' + $scope.ward.wardId);
                };
            });
        });
    }]);


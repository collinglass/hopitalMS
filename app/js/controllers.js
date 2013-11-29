'use strict';

/* Controllers */

angular.module('mustacheApp.controllers', []).
    controller('LoginCtrl', [function () {

    }])
    .controller('RegisterCtrl', [function () {

    }]).controller('WardListCtrl', [function () {
        // Doing nothing

    }]).controller('WardDetailCtrl', ["$scope", function ($scope) {

        $scope.admitPatient = function() {
            angular.forEach($scope.admissions, function(obj) {
                if ( obj.selected == true ) {
                    $scope.patients.push({lastName: obj.lastName, firstName: obj.firstName, healthInsNum: obj.healthInsNum, roomNum: "00", bedNum: "00", status: "nominal",});
                    $scope.admissions.splice(obj, 1);
                }
            });
        };

        $scope.patients = [
            {
                lastName: "Smith",
                firstName: "John",
                healthInsNum: "0123 4567 89AB CDEF",
                roomNum: 9000,
                bedNum: 42,
                status: "nominal"
            }
        ];
        $scope.admissions = [
            {
                lastName: "Balls",
                firstName: "Harry",
                healthInsNum: "9876 5432 10FE DCBA",
                outboundWard: 9001,
                outboundNurse: 41,
                priority: "nominal",
                rationale: "The patient..."
            }
        ];

    }]);

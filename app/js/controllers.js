'use strict';

/* Controllers */

angular.module('mustacheApp.controllers', []).
    controller('LoginCtrl', [function () {

    }])
    .controller('RegisterCtrl', [function () {

    }]).controller('WardListCtrl', [function () {
        // Doing nothing

    }]).controller('WardDetailCtrl', ["$scope", function ($scope) {
        $scope.patients = [
            {
                lastName: "Smith",
                firstName: "John",
                heathInsNum: "0123 4567 89AB CDEF",
                roomNum: 9000,
                bedNum: 42,
                status: "nominal"
            }
        ];

    }]);

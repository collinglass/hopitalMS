(function (routingConfig) {
    'use strict';

    var userRoles = {
        public: 1, // 0001
        medicalStaff: 2, // 0010
        doctor: 4, // 0100
        chargeNurse: 8 // 1000
    };

    routingConfig.userRoles = userRoles;
    routingConfig.accessLevels = {
        public: userRoles.public | // 1111
            userRoles.medicalStaff |
            userRoles.doctor |
            userRoles.chargeNurse,
        anon: userRoles.public, // 0001
        medicalStaff: userRoles.medicalStaff | // 1110
            userRoles.doctor |
            userRoles.chargeNurse,
        doctor: userRoles.doctor, // 0100
        chargeNurse: userRoles.chargeNurse // 1000
    };
})(this.routingConfig);
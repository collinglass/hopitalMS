(function(exports){

	var userRoles = {
		public: 1, // 0001
		medical_staff: 2, // 0010
		doctor: 4, // 0100
		charge_nurse: 8 // 1000
	};

	exports.userRoles = userRoles;
	exports.accessLevels = {
		public: userRoles.public | // 1111
				userRoles.medical_staff |
				userRoles.doctor |
				userRoles.charge_nurse,
		anon: userRoles.public, // 0001
		medical_staff: userRoles.medical_staff | // 1110
						userRoles.doctor |
						userRoles.charge_nurse,
		doctor: userRoles.doctor, // 0100
		charge_nurse: userRoles.charge_nurse // 1000
	};
})(typeof exports === 'undefined'? this['routingConfig']={}: exports);
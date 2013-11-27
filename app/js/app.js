function LoginCtrl($scope) {

  var User = function(login) {
    var wasNotNull = login || false;
    login = login ||{};
    return {
      name: login.employeeId || "",
      employeeId: login.employeeId || "",
      loggedIn: wasNotNull && login.employeeId === login.password
    };
  };

  $scope.User = new User();

  $scope.Login = {
    employeeId: undefined,
    password: undefined,
    employeeIdPlaceholder: "Employee ID",
    passwordPlaceholder: "Password"
  };


  $scope.onNewAccount = function() {
    var login = $scope.Login;
    console.log("New account for :" + JSON.stringify(login));
  };

  $scope.onLogin = function() {
    var login = $scope.Login;
    console.log("Logging in with :" + JSON.stringify(login));

    if (login.employeeId === login.password) {
      $scope.User = new User(login);
    } else {
      // show "invalid username or password"
      console.log("Invalid username or password :" + JSON.stringify(login));
    }
  };

  $scope.onLogOut = function() {
    console.log("Logging out from :" + $scope.User.name);
    $scope.User = new User({});
  };
}
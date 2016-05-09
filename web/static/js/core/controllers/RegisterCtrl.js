define(function (require) {
    'use strict';
    var angular = require('angular');
    return angular.module('app.controllers.RegisterCtrl', [])
        .controller('RegisterCtrl', ["$scope", "alert", "$auth", "vcRecaptchaService", "$state",
            function ($scope, alert, $auth, vcRecaptchaService, $state) {
                $scope.gRecaptchResponse = null;
                $scope.gRecaptchWidgetId = null;
                $scope.gRecaptchModel = {
                    key: '6LcQ4xYTAAAAAFonPEBORwlB7r3m2J1ctVcLF5cu'
                };
                $scope.setgRecaptchResponse = function (response) {
                    $scope.gRecaptchResponse = response;
                };
                $scope.setgRecaptchWidgetId = function (widgetId) {
                    $scope.gRecaptchWidgetId = widgetId;
                };
                $scope.submit = function () {
                    $auth.signup({
                        email: $scope.email,
                        password: $scope.password,
                        username: $scope.username,
                        recaptcharesponse: $scope.gRecaptchResponse
                    }).then(function (res) {
                        alert('success', 'Account Created!', 'Welcome, ' + res.data.user.email + '! Please verify your email to avoid account termination. :)');
                        $state.go('home');
                    }).catch(function (err) {
                        vcRecaptchaService.reload($scope.gRecaptchWidgetId);
                        alert('warning', 'Something went wrong :( ', err.data.message);
                    });
                };
            }]);
});
define(function (require) {
    'use strict';
    var angular = require('angular');
    return angular.module('app.controllers.LoginCtrl', [])
        .controller('LoginCtrl', ["$scope", "alert", "$auth", "vcRecaptchaService", "$state",
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
                    console.log(widgetId);
                };
                $scope.submit = function () {
                    $auth.login({
                        email: $scope.email,
                        password: $scope.password,
                        recaptcharesponse: $scope.gRecaptchResponse
                    }).then(function (res) {
                        var message = 'Thanks for coming back, ' + res.data.user.email + '!';

                        if (res.data.user.verified == "false") {
                            message = 'Just a reminder, please verify your email to avoid account termination. :)';
                        }

                        alert('success', 'Welcome!', message, 6000);
                        $state.go('home');
                    }).catch(handleError);
                };

                $scope.authenticate = function (provider) {
                    $auth.authenticate(provider).then(function (res) {
                        alert('success', 'Welcome!', 'Thanks for coming back!');
                    }, handleError);
                };

                function handleError(err) {
                    vcRecaptchaService.reload($scope.gRecaptchWidgetId);
                    alert('warning', 'Oops! ', err.data.message);
                }
            }]);
});
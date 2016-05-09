define(function (require) {
    'use strict';
    var angular = require('angular');
    return angular.module('app.controllers.ForgotPasswordCtrl', [])
        .controller('ForgotPasswordCtrl', ["$scope", "alert", "$http", "vcRecaptchaService",
            function ($scope, alert, $http, vcRecaptchaService) {
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
                var user_post = {
                    gRecaptchResponse: $scope.gRecaptchResponse,
                    email: $scope.email,
                    username: $scope.username
                };

                $scope.submit = function () {
                    $http({
                        method: 'POST',
                        url: '//windsofarithia.com/api/auth/forgotpassword',
                        data: user_post,
                        headers: {'Content-Type': 'application/json'}
                    })
                        .success(function (data) {
                            console.log(data);
                        })
                        .catch(handleError);
                };

                function handleError(err) {
                    vcRecaptchaService.reload($scope.gRecaptchWidgetId);
                    alert('warning', 'Oops! ', err.data.message);
                }
            }]);
});
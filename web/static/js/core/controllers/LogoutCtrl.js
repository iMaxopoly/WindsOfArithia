define(function (require) {
    'use strict';
    var angular = require('angular');
    return angular.module('app.controllers.LogoutCtrl', [])
        .controller('LogoutCtrl', ["$auth", "$state", 'alert', "$http", function ($auth, $state, alert, $http) {
            $http({
                method: 'POST',
                url: '//windsofarithia.com/api/auth/logout'
            })
                .success(function (data) {
                    console.log(data);
                })
                .catch(handleError);

            function handleError(err) {
                alert('warning', 'Oops! ', err.data.message);
            }

            $auth.logout();
            $state.go('home');
        }]);
});
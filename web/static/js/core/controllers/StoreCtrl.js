define(function (require) {
    'use strict';
    var angular = require('angular');
    return angular.module('app.controllers.StoreCtrl', [])
        .controller('StoreCtrl', ["$scope", "$window", function ($scope, $window) {
            $scope.tabs = [
                { title:'Dynamic Title 1', content:'Dynamic content 1' },
                { title:'Dynamic Title 2', content:'Dynamic content 2', disabled: true }
            ];

            $scope.alertMe = function() {
                setTimeout(function() {
                    $window.alert('You\'ve selected the alert tab!');
                });
            };

            $scope.model = {
                name: 'Tabs'
            };
        }]);
});
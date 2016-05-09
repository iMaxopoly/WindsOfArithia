define(['angular', 'uirouter', 'uibootstrap', 'satellizer', 'angularAnimate', 'angularCookies',
    'angularSanitizer', 'angularMoment', 'preloadImage', 'loadingbar',
    'angularRecaptcha', './core/services/0_services', './core/controllers/0_controllers', './core/directives/0_directives',
    'util', 'time', 'apiGoogleRecaptcha'], function (angular) {
    'use strict';

    var app = angular.module('app', ['ui.router', 'ngAnimate', 'ngSanitize', 'ngCookies','ui.bootstrap', 'angular-preload-image',
        'satellizer', 'angularMoment', 'angular-loading-bar', 'vcRecaptcha', 'app.services', 'app.directives', 'app.controllers']);

    app.init = function () {
        angular.bootstrap(document, ['app']);
    };

    app.config(['$stateProvider', '$urlRouterProvider', '$interpolateProvider', '$authProvider', '$locationProvider',
        '$httpProvider', 'cfpLoadingBarProvider', 'API_URL',
        function ($stateProvider, $urlRouterProvider, $interpolateProvider, $authProvider, $locationProvider,
                  $httpProvider, cfpLoadingBarProvider, API_URL) {

            cfpLoadingBarProvider.includeSpinner = false;
            //AnalyticsProvider.setAccount('UA-34308301-5');
            //AnalyticsProvider.trackPages(true);
            //AnalyticsProvider.setPageEvent('$stateChangeSuccess');

            $interpolateProvider.startSymbol('{[').endSymbol(']}');
            $urlRouterProvider.otherwise('/');
            $locationProvider.html5Mode(true);
            $stateProvider
                .state('home', {url: '/', templateUrl: 'woa-home', controller: 'HomeCtrl'})
                .state('news', {url: '/news', controller: 'NewsCtrl'})
                .state('community', {url: '/community', controller: 'CommunityCtrl'})
                .state('login', {url: '/login', templateUrl: 'woa-login', controller: 'LoginCtrl'})
                .state('verifiedmail', {
                    url: '/verified-mail',
                    templateUrl: 'woa-verified-mail'
                })
                .state('unverifiedmail', {
                    url: '/unverified-mail',
                    templateUrl: 'woa-unverified-mail'
                })
                .state('forgot_password', {
                    url: '/forgot-password',
                    templateUrl: 'woa-forgot-password',
                    data: {requireLogin: false},
                    controller: 'ForgotPasswordCtrl'
                })
                .state('change_password', {
                    url: '/change-password',
                    templateUrl: 'woa-change-password',
                    controller: 'ChangePasswordCtrl'
                })
                .state('register', {url: '/register', templateUrl: 'woa-register', controller: 'RegisterCtrl'})
                .state('logout', {url: '/logout', data: {requireLogin: true}, controller: 'LogoutCtrl'})
                .state('store', {url: '/store', templateUrl: 'woa-store', data: {requireLogin: true}, controller: 'StoreCtrl'})
                .state('donate', {url: '/donate', templateUrl: 'woa-donate', data: {requireLogin: true}})
                .state('support', {url: '/support', templateUrl: 'woa-support', data: {requireLogin: true}})
                .state('chat', {
                    url: '/chat',
                    data: {requireLogin: true},
                    controller: 'ChatCtrl'
                });

            $authProvider.baseUrl = '/api/';
            $authProvider.facebook({
                clientId: '1616506938606346',
                url: 'auth/facebook'
            });
            $authProvider.google({
                clientId: '836322842190-7j8j00vi4tuhinc7mqngerbl8ev15as9.apps.googleusercontent.com',
                url: 'auth/google'
            });
            $authProvider.loginUrl = 'auth/login';
            $authProvider.signupUrl = 'auth/register';
            $authProvider.tokenPrefix = 'woa';
        }
    ]);

    app.constant('API_URL', '//windsofarithia.com/api/');
    //app.constant('WSOCKET_URL', 'wss://windsofarithia.com/ws');

    app.run(function ($state, $window) {
        var params = $window.location.search.substring(1);

        if (params && $window.opener && $window.opener.location.origin === $window.location.origin) {
            var pair = params.split('=');
            var code = decodeURIComponent(pair[1]);

            $window.opener.postMessage(code, $window.location.origin);
        }
    });

    app.run(function ($rootScope, $location, $state, $auth, alert) {
        $rootScope.$on('$stateChangeStart', function (e, toState, toParams, fromState, fromParams) {
            if (!(toState.name === "home")) {
                $rootScope.ariGirl = false;
            } else {
                $rootScope.ariGirl = true;
            }

            var isLogged = $auth.isAuthenticated();
            var isLogin = toState.name === "login";
            if (isLogin && isLogged === false) {
                return; // no need to redirect
            }

            // now, redirect only not authenticated

            if (isLogged === false && toState.data !== undefined && toState.data.requireLogin === true) {
                e.preventDefault(); // stop current execution
                $state.go('login'); // go to login
                return;
            }

            if (isLogged === true && (toState.name === "register" || toState.name === "login")) {
                e.preventDefault(); // stop current execution
                $state.go('home'); // go to login
            }
        });
        $rootScope.$on('$stateChangeError', function (event) {
            event.preventDefault();
            alert('warning', 'Unable to resolve.');
            $state.go('home');
        });
    });
    app.directive('imgPreload', ['$rootScope', function ($rootScope) {
        return {
            restrict: 'A',
            scope: {
                ngSrc: '@'
            },
            link: function (scope, element, attrs) {
                element.on('load', function () {
                    element.addClass('in');
                }).on('error', function () {
                    //
                });

                scope.$watch('ngSrc', function (newVal) {
                    element.removeClass('in');
                });
            }
        };
    }]);
    return app;

});
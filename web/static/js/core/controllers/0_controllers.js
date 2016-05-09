define(function (require) {
    'use strict';

    var angular = require('angular'),
        MainCtrl = require('./MainCtrl'),
        ChatCtrl = require('./ChatCtrl'),
        NewsCtrl = require('./NewsCtrl'),
        GuildsCtrl = require('./GuildsCtrl'),
        HeaderCtrl = require('./HeaderCtrl'),
        LoginCtrl = require('./LoginCtrl'),
        ForgotPasswordCtrl = require('./ForgotPasswordCtrl'),
        ChangePasswordCtrl = require('./ChangePasswordCtrl'),
        LogoutCtrl = require('./LogoutCtrl'),
        HomeCtrl = require('./HomeCtrl'),
        CommunityCtrl = require('./CommunityCtrl'),
        RegisterCtrl = require('./RegisterCtrl'),
        StoreCtrl = require('./StoreCtrl');

    return angular.module('app.controllers', [MainCtrl.name, ChatCtrl.name, StoreCtrl.name, NewsCtrl.name, GuildsCtrl.name, HeaderCtrl.name, ForgotPasswordCtrl.name,
        ChangePasswordCtrl.name, LoginCtrl.name, LogoutCtrl.name, HomeCtrl.name, CommunityCtrl.name, RegisterCtrl.name]);
});
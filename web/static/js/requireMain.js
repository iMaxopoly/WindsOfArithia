require.config({
    paths: {
        'apiGoogleRecaptcha': 'https://www.google.com/recaptcha/api.js?onload=vcRecaptchaApiLoaded&render=explicit',
        //vendor mods

        //bower components
        'angular': 'bower_components/angular/angular',
        'jquery': 'bower_components/jquery/dist/jquery',
        'uirouter': 'bower_components/angular-ui-router/release/angular-ui-router',
        'satellizer': 'bower_components/magnet-satellizer/satellizer',
        'angularMoment': 'bower_components/angular-moment/angular-moment',
        'angularSanitizer': 'bower_components/angular-sanitize/angular-sanitize',
        'angularAnimate': 'bower_components/angular-animate/angular-animate',
        'angularCookies': 'bower_components/angular-cookies/angular-cookies',
        'uibootstrap': 'bower_components/angular-bootstrap/ui-bootstrap-tpls',
        'angularRecaptcha': 'bower_components/angular-recaptcha/release/angular-recaptcha',
        //'ganalytics': 'bower_components/angular-google-analytics/dist/angular-google-analytics',
        'loadingbar': 'bower_components/angular-loading-bar/build/loading-bar',
        'preloadImage': 'bower_components/angular-preload-image/angular-preload-image',
        'moment': 'bower_components/moment/moment',
        'panorama': 'core/panorama',
        'util': 'core/util',
        'time': 'core/time'
    },
    shim: {
        'jquery': {exports: '$'},
        'angular': {deps: ['jquery'], exports: 'angular'},
        'moment': {deps: ['jquery'], exports: 'moment'},
        'uirouter': ['angular'],
        'uibootstrap': ['angular'],
        'satellizer': ['angular'],
        'loadingbar': ['angular'],
        'angularRecaptcha': ['angular'],
        'preloadImage': ['angular'],
        'apiGoogleRecaptcha': ['angularRecaptcha'],
        'angularSanitizer': {exports: 'ngSanitize', deps: ['angular']},
        'angularAnimate': {exports: 'ngAnimate', deps: ['angular']},
        'angularCookies': {exports: 'ngCookies', deps: ['angular']},
        'angularMoment': ['angular', 'moment'],
        //'ganalytics': ['angular'],
        'panorama': ['angular'],
        'util': ['panorama'],
        'time': ['angular', 'moment']
    },
    priority: [
        "jquery",
        "angular"
    ]
});


require(['app'], function (app) {
    app.init();
});

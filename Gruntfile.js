module.exports = function (grunt) {
    "use strict";
    grunt.initConfig({
        cssmin: {
            options: {
                shorthandCompacting: false,
                roundingPrecision: -1,
                keepSpecialComments: 0

            },
            target: {
                files: {
                    'web/dist/css/lunastyle.css': [
                        "web/static/js/bower_components/bootstrap/dist/css/bootstrap.css",
                        'web/static/js/bower_components/font-awesome/css/font-awesome.css',
                        'web/static/js/bower_components/animate.css/animate.css',
                        'web/static/css/loading-bar.css',
                        'web/static/css/sidetabs.css',
                        'web/static/css/aristyle.css'
                    ]
                }
            }
        },
        htmlmin: {                                     // Task
            dist: {                                      // Target
                options: {                                 // Target options
                    removeComments: true,
                    collapseWhitespace: true
                },
                files: {                                   // Dictionary of files
                    'web/dist/views/base.html': 'web/views/base.html',
                    'web/dist/views/change_password.html': 'web/views/change_password.html',
                    'web/dist/views/disclaimer.html': 'web/views/disclaimer.html',
                    'web/dist/views/donate.html': 'web/views/donate.html',
                    'web/dist/views/store.html': 'web/views/store.html',
                    'web/dist/views/footer.html': 'web/views/footer.html',
                    'web/dist/views/forgot_password.html': 'web/views/forgot_password.html',
                    'web/dist/views/header.html': 'web/views/header.html',
                    'web/dist/views/home.html': 'web/views/home.html',
                    'web/dist/views/index.html': 'web/views/index.html',
                    'web/dist/views/login.html': 'web/views/login.html',
                    'web/dist/views/mail_verification.html': 'web/views/mail_verification.html',
                    'web/dist/views/news.html': 'web/views/news.html',
                    'web/dist/views/privacy_policy.html': 'web/views/privacy_policy.html',
                    'web/dist/views/register.html': 'web/views/register.html',
                    'web/dist/views/support.html': 'web/views/support.html',
                    'web/dist/views/unverified_email.html': 'web/views/unverified_email.html',
                    'web/dist/views/verified_email.html': 'web/views/verified_email.html',
                }
            }
        },
        requirejs: {
            compile: {
                options: {
                    uglify2: {
                        mangle: false
                    },
                    baseUrl: "./web/static/js/",
                    mainConfigFile: "web/static/js/requireMain.js",
                    out: "web/dist/js/requireMain.js",
                    optimize: "uglify2",
                    name: 'requireMain',
                    removeCombined: true,
                    findNestedDependencies: true,
                    preserveLicenseComments: false
                }
            }
        }
    });
    grunt.loadNpmTasks('grunt-contrib-htmlmin');
    grunt.loadNpmTasks('grunt-contrib-requirejs');
    grunt.loadNpmTasks('grunt-contrib-cssmin');
    grunt.registerTask('default', ['htmlmin', 'cssmin', 'requirejs']);
};
'use strict';

angular.module('lytup', [
  'ngRoute',
  'ngSocket',
  'restangular',
  'angularFileUpload',
  'ui.bootstrap',
  'angular-loading-bar',
  'ngAnimate',
  'lytup.filters',
  'lytup.services',
  'lytup.directives',
  'lytup.controllers'
])
  .config(['$locationProvider',
    '$httpProvider',
    '$routeProvider',
    'RestangularProvider',
    function($locationProvider, $httpProvider, $routeProvider, RestangularProvider) {
      $locationProvider.html5Mode(true);

      $httpProvider.interceptors.push('AuthInterceptor');

      // $routeProvider.when('/', {
      //   controller: 'LandingCtrl',
      //   templateUrl: '/tpl/landing.html'
      // });

      $routeProvider.when('/', {
        templateUrl: '/tpl/home.html'
      });

      $routeProvider.when('/join', {
        controller: 'SignupCtrl',
        templateUrl: '/tpl/home.html'
      });

      $routeProvider.when('/login', {
        controller: 'SigninCtrl',
        templateUrl: '/tpl/home.html'
      });

      $routeProvider.when('/confirm/:key', {
        controller: 'ConfirmCtrl',
        template: ''
      });

      $routeProvider.when('/forgot', {
        controller: 'ForgotPwdCtrl',
        templateUrl: '/tpl/home.html'
      });

      $routeProvider.when('/reset/:key', {
        controller: 'ResetPwdCtrl',
        templateUrl: '/tpl/home.html'
      });

      $routeProvider.when('/:id', {
        controller: 'FolderCtrl',
        templateUrl: '/tpl/folder.html'
      });

      $routeProvider.when('/i/:id', {
        controller: 'FileCtrl',
        templateUrl: '/tpl/file.html'
      });

      $routeProvider.otherwise({
        redirectTo: '/'
      });

      RestangularProvider.setBaseUrl('/api');
    }
  ]).run(['$rootScope',
    function($rootScope) {
      $rootScope.BASE_URI = location.protocol + '//' + location.hostname + (location.port && ':' + location.port);

      $rootScope.MESSAGE = function(code) {
        return {
          blankFirstName: 'First name is required',
          blankLastName: 'Last name is required',
          blankEmail: 'Email is required',
          invalidEmail: 'Email is invalid',
          registeredEmail: 'This email is already registered',
          blankPassword: 'Password cannot be blank',
          invalidPassword: 'Password must be at least 6 characters',
          mismatchPasswords: 'Passwords don\'t match',
          blankExpiry: 'Expiry is required',
          loginFailed: 'Login faild, invalid email or password'
        }[code];
      };

      toastr.options = {
        'positionClass': 'toast-bottom-right'
      };
    }
  ]);

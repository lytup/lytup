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

      $routeProvider.when('/verify/:code', {
        controller: 'VerifyEmailCtrl',
        template: ''
      });

      $routeProvider.when('/forgot', {
        controller: 'ForgotPwdCtrl',
        templateUrl: '/tpl/home.html'
      });

      $routeProvider.when('/reset/:code', {
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

      $rootScope.MESSAGE = function(key) {
        // TODO: Get it from server
        return {
          blankFirstName: 'First name is required',
          blankLastName: 'Last name is required',
          blankEmail: 'Email is required',
          invalidEmail: 'Email is invalid',
          blankPassword: 'Password cannot be blank',
          invalidPassword: 'Password must be at least 6 characters',
          passwordMistach: 'Passwords do not match',
          blankExpiry: 'Expiry is required'
        }[key];
      };

      toastr.options = {
        'positionClass': 'toast-bottom-right'
      };
    }
  ]);

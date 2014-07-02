'use strict';

angular.module('lytup', [
  'ngRoute',
  'ngSocket',
  'restangular',
  'angularFileUpload',
  'ui.bootstrap',
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
        controller: 'LandingCtrl',
        templateUrl: '/tpl/landing.html'
      });

      $routeProvider.when('/home', {
        controller: 'HomeCtrl',
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
    }
  ]);

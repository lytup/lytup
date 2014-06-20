'use strict';

angular.module('lytup', [
  'ngRoute',
  'ngSocket',
  'restangular',
  'angularFileUpload',
  'lytup.filters',
  'lytup.services',
  'lytup.directives',
  'lytup.controllers'
])
  .config(['$locationProvider', '$routeProvider', 'RestangularProvider',
    function($locationProvider, $routeProvider, RestangularProvider) {
      $locationProvider.html5Mode(true);
      RestangularProvider.setBaseUrl('/api');

      $routeProvider.when('/', {
        templateUrl: 'partials/landing.html'
      });

      $routeProvider.when('/home', {
        controller: 'HomeCtrl',
        templateUrl: 'partials/home.html'
      });

      $routeProvider.when('/:id', {
        controller: 'FolderCtrl',
        templateUrl: 'partials/folder.html'
      });

      $routeProvider.otherwise({
        redirectTo: '/'
      });
    }
  ]).run(['$rootScope',
    function($rootScope) {
      $rootScope.BASE_URI = location.protocol + '//' + location.hostname
          + (location.port && ':' + location.port);
    }
  ]);

'use strict';

angular.module('lytup.services', [])
  .factory('AuthInterceptor', function($rootScope, $window, $q) {
    return {
      request: function(cfg) {
        if ($window.sessionStorage.token) {
          cfg.headers.Authorization = 'Bearer ' + $window.sessionStorage.token;
        }
        return cfg;
      },
      response: function(res) {
        if (res.status === 401) {
          // Handle the case where the user is not authenticated
        }
        return res || $q.when(res);
      }
    };
  })
  .value('version', '0.0.1');

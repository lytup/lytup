'use strict';

angular.module('lytup.services', [])
  .factory('AuthInterceptor', [
    '$window',
    '$q',
    function($window, $q) {
      return {
        request: function(cfg) {
          var token = $window.localStorage.getItem('token');
          if (token) {
            cfg.headers.Authorization = 'Bearer ' + token;
          }
          return cfg;
        },
        responseError: function(res) {
          if (res.status === 401) {
            $window.localStorage.removeItem('token');
          }
          return $q.reject(res);
        }
      };
    }
  ])
  .value('version', '0.0.1');

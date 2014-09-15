'use strict';

angular.module('lytup.services', [])
  .factory('AuthInterceptor', [
    '$window',
    '$q',
    'Notification',
    function($window, $q, Notification) {
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
          Notification.error(res.data.error);
          return $q.reject(res);
        }
      };
    }
  ])
  .factory('Notification', [
    function() {
      return {
        info: function(msg) {
          toastr.info(msg);
        },
        success: function(msg) {
          toastr.success(msg);
        },
        warning: function(msg) {
          toastr.warning(msg);
        },
        error: function(msg) {
          toastr.error(msg);
        }
      }
    }
  ])
  .value('version', '0.0.1');

'use strict';

angular.module('lytup.filters', [])
  .filter('interpolate', ['version',
    function(version) {
      return function(text) {
        return String(text).replace(/\%VERSION\%/mg, version);
      };
    }
  ])
  .filter('bytes', [
    function() {
      return function(bytes) {
        var unit = 1024;
        if (!bytes || bytes === -1) return '--';
        if (bytes < unit) return bytes + ' B';
        var exp = ~~ (Math.log(bytes) / Math.log(unit));
        var pre = 'KMGTPE'.charAt(exp - 1);
        return (bytes / Math.pow(unit, exp)).toFixed(1) + ' ' + pre + 'B';
      };
    }
  ]);

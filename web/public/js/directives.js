'use strict';

angular.module('lytup.directives', [])
  .directive('lyCopy', [

    function() {
      return {
        link: function(scope, elm) {
          var clip = new ZeroClipboard(elm);

          elm.tooltip();

          elm.on('mouseover', function() {
            elm.attr('data-original-title', 'Copy Link')
              .tooltip('show');
          });

          clip.on('aftercopy', function() {
            elm.attr('data-original-title', 'Copied!')
              .tooltip('show');
          });

          clip.on('error', function(event) {
            ZeroClipboard.destroy();
          });
        }
      };
    }
  ])
  .directive('lyDownload', [
    function() {
      return {
        link: function(scope, elm, attrs) {
          var iframe = $('<iframe>').hide().appendTo(elm);
          elm.on('click', function() {
            iframe.attr('src', attrs.uri);
          });
        }
      };
    }
  ])
  .directive('lyKnob', [
    function() {
      return {
        link: function(scope, elm, attrs) {
          elm.knob({
            readOnly: true
          });

          attrs.$observe('value', function(val) {
            elm.val(val).trigger('change');
          });
        }
      }
    }
  ])
  .directive('bsTooltip', [function() {
      return {
        link: function(scope, elm) {
          elm.tooltip();
        }
      };
    }
  ])
  .directive('appVersion', ['version',
    function(version) {
      return function(scope, elm, attrs) {
        elm.text(version);
      };
    }
  ]);

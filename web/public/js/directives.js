'use strict';

angular.module('lytup.directives', [])
  .directive('lyCopy', [

    function() {
      return {
        link: function(scope, elm, attrs) {
          var clt = new ZeroClipboard(elm);

          elm.on('mouseover', function() {
            scope.text = attrs.copyText;
            scope.$apply();
          });

          clt.on('aftercopy', function() {
            scope.text = attrs.copiedText;
            scope.$apply();
          });

          clt.on('error', function() {
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
          elm.on('click', function(evt) {
            iframe.attr('src', attrs.href);
            event.preventDefault();
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
  .directive('lyFromNow', [

    function() {
      return {
        link: function(scope, elm, attrs) {
          attrs.$observe('date', function(val) {
            if (val) {
              elm.text(moment(attrs.date).fromNow());
            }
          });
        }
      }
    }
  ])
  .directive('appVersion', ['version',
    function(version) {
      return function(scope, elm, attrs) {
        elm.text(version);
      };
    }
  ]);

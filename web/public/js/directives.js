'use strict';

angular.module('lytup.directives', [])
// HTML5 autofocus works only once during the page load,
// so re-opening a modal doesn't loose the focus
// https://github.com/angular-ui/bootstrap/issues/1696
.directive('lyFocus', ['$timeout',
  function($timeout) {
    return {
      link: function(scope, elm) {
        $timeout(function() {
          elm.focus();
        }, 500);
      }
    };
  }
])
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
  .directive('lyDeadCenter', [

    function() {
      return {
        link: function(scope, elm) {
          elm.parent().css('position', 'relative');
          elm.css('position', 'absolute');
          elm.css('top', '50%');
          elm.css('left', '50%');
          elm.css('margin-top', '-' + elm.height() / 2 + 'px');
          elm.css('margin-left', '-' + elm.width() / 2 + 'px');
        }
      }
    }
  ])
  .directive('appVersion', ['version',
    function(version) {
      return function(scope, elm) {
        elm.text(version);
      };
    }
  ]);

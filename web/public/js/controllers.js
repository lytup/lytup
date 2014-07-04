'use strict';

angular.module('lytup.controllers', [])
  .controller('HomeCtrl', [
    '$scope',
    '$window',
    '$location',
    '$log',
    'ngSocket',
    '$modal',
    'Restangular',
    function($scope, $window, $location, $log, ngSocket, $modal, Restangular) {
      $log.info('Home controller');
      var token = $window.localStorage.getItem('token');
      var ws = ngSocket('ws://localhost:1431/ws');

      ws.onMessage(function(msg) {
        $log.info('message received', msg.data);
      });

      if (token) {
        Restangular.one('users').get().then(function(usr) {
          $scope.user = usr;
          $scope.user.token = token;
          $scope.folders = Restangular.all('folders').getList().$object;
        });
      }

      $scope.openSignup = function() {
        $modal.open({
          scope: $scope,
          controller: 'SignupCtrl',
          templateUrl: '/tpl/signup.html',
          size: 'sm',
          resolve: {
            'user': function() {
              return $scope.user;
            }
          }
        });
      };

      $scope.openSignin = function() {
        $modal.open({
          scope: $scope,
          controller: 'SigninCtrl',
          templateUrl: '/tpl/signin.html',
          size: 'sm',
          resolve: {
            'user': function() {
              return $scope.user;
            },
            'folders': function() {
              return $scope.folders;
            }
          }
        });
      };

      $scope.signout = function() {
        $window.localStorage.removeItem('token');
        $scope.user = {};
        $location.path('/');
      };

      $scope.createFolder = function() {
        Restangular.all('folders').post({}).then(function(fol) {
          $scope.folder = fol;
          $scope.folders.push(fol);
          $location.path('/' + fol.id);
        });
      };

      $scope.deleteFolder = function(fol) {
        fol.remove().then(function() {
          _.remove($scope.folders, {
            'id': fol.id
          });
        });
      };
    }
  ])
  .controller('FolderCtrl', [
    '$scope',
    '$routeParams',
    '$log',
    'Restangular',
    '$upload',
    function($scope, $routeParams, $log, Restangular, $upload) {
      $log.info('Folder controller');
      $scope.folder = Restangular.one('folders', $routeParams.id).get().$object;

      $scope.addFiles = function(files) {
        var fol = $scope.folder;

        _.forEach(files, function(file) {
          var f = _.pick(file, 'name', 'size', 'type');
          var i = fol.files.push(f) - 1;

          file.i = i; // Store the index for later user

          // Create file
          fol.files[i] = fol.post('files', f).$object;
        });

        uploadFiles(files);
      };

      $scope.fileIconClass = function(type) {
        return /image/.test(type) ? 'fa-file-image-o' :
          /audio/.test(type) ? 'fa-file-audio-o' :
          /video/.test(type) ? 'fa-file-video-o' :
          /document/.test(type) ? 'fa-file-word-o' :
          /sheet/.test(type) ? 'fa-file-excel-o' :
          /presentation/.test(type) ? 'fa-file-powerpoint-o' :
          /pdf/.test(type) ? 'fa-file-pdf-o' :
          /zip/.test(type) ? 'fa-file-archive-o' :
          'fa-file-o';
      };

      $scope.deleteFile = function(id) {
        $scope.folder.one('files', id).remove().then(function() {
          _.remove($scope.folder.files, {
            'id': id
          });
        });
      };

      function uploadFiles(files) {
        var fol = $scope.folder;

        _.forEach(files, function(file, i) {
          var f = fol.files[file.i];

          $upload.upload({
            url: '/u/' + fol.id,
            file: file
          }).progress(function(evt) {
            f.loaded = Math.round(evt.loaded / evt.total * 100);
          }).success(function() {
            // Update file
            f.patch({
              loaded: 100
            });
          });
        });
      }
    }
  ])
  .controller('FileCtrl', [
    '$scope',
    '$routeParams',
    '$log',
    'Restangular',
    function($scope, $routeParams, $log, Restangular) {
      $log.info('File controller');
      $scope.file = Restangular.one('files', $routeParams.id).get().$object;
    }
  ])
  .controller('SignupCtrl', [
    '$scope',
    '$window',
    '$log',
    '$modalInstance',
    'Restangular',
    function($scope, $window, $log, $modalInstance, Restangular) {
      $log.info('Signup controller');

      $scope.signup = function(user) {
        Restangular.all('users').post(user).then(function(usr) {
          $window.localStorage.setItem('token', usr.token);
          $scope.$parent.user = usr;
          $modalInstance.close();
        });
      };

      $scope.cancel = function() {
        $modalInstance.dismiss('cancel');
      };
    }
  ])
  .controller('SigninCtrl', [
    '$scope',
    '$window',
    '$log',
    '$modalInstance',
    'Restangular',
    function($scope, $window, $log, $modalInstance, Restangular) {
      $log.info('Signin controller');

      $scope.signin = function(user) {
        Restangular.all('users').all('login').post(user)
          .then(function(usr) {
            $window.localStorage.setItem('token', usr.token);
            $scope.$parent.user = usr;
            $scope.$parent.folders = Restangular.all('folders').getList().$object;
            $modalInstance.close();
          });
      };

      $scope.cancel = function() {
        $modalInstance.dismiss('cancel');
      };
    }
  ]);

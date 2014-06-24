'use strict';

angular.module('lytup.controllers', [])
  .controller('Controller', [
    '$scope',
    'ngSocket',
    function($scope, ngSocket) {
      console.log("Lytup controller")
      var ws = ngSocket('ws://localhost:1431/ws');

      ws.onMessage(function(msg) {
        console.log('message received', msg.data);
      });

      $scope.expiresIn = function(date) {
        return moment(date).fromNow();
      };
    }
  ])
  .controller('HomeCtrl', [
    '$scope',
    '$location',
    'Restangular',
    function($scope, $location, Restangular) {
      console.log("Home controller");
      $scope.folders = Restangular.all('folders').getList().$object;
      $scope.folder = {};

      $scope.createFolder = function() {
        $scope.folders.post({}).then(function(fol) {
          $scope.folder = fol;
          $scope.folders.push(fol);
          $location.path(fol.id);
        });
      };

      $scope.deleteFolder = function(fol, i) {
        fol.remove().then(function() {
          $scope.folders.splice(i, 1);
        });
      };
    }
  ])
  .controller('FolderCtrl', [
    '$scope',
    '$routeParams',
    'Restangular',
    '$upload',
    function($scope, $routeParams, Restangular, $upload) {
      console.log("Folder controller");
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

      $scope.deleteFile = function(file, i) {
        $scope.folder.one('files', file.id).remove().then(function() {
          $scope.folder.files.splice(i, 1);
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
  'Restangular',
    function($scope, $routeParams, Restangular) {
      console.log('File controller');
      $scope.file = Restangular.one('files', $routeParams.id).get().$object;
    }
  ]);

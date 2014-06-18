'use strict';

angular.module('lytup.controllers', [])
  .controller('Controller', ['$scope', '$location', '$routeParams', 'ngSocket', 'Restangular', '$upload',
    function($scope, $location, $routeParams, ngSocket, Restangular, $upload) {
      console.log("Lytup controller")
      // var ws = ngSocket('ws://localhost:1431/ws');
      //
      // ws.onMessage(function(msg) {
      //   console.log('message received', msg.data);
      // });
      //
      // ws.send({
      //   foo: 'bar'
      // });

      // Get all folders
      $scope.folders = Restangular.all('folders').getList().$object;
      $scope.folder = {};

      $scope.createFolder = function() {
        Restangular.all('folders').post({}).then(function(fol) {
          $scope.folder = fol;
          $scope.folders.push(fol);
          $location.path(fol.id);
        });
      };

      $scope.addFiles = function(files) {
        var fol = $scope.folder;

        _.forEach(files, function(file) {
          var f = _.pick(file, 'name', 'size', 'type');
          file.i = fol.files.push(f) - 1; // Store index for mapping later
        });

        // Update folder
        fol.patch({
          'files': fol.files.slice(-files.length)
        }).then(function() {
          uploadFiles(files);
        });
      };

      function uploadFiles(files, mFiles) {
        var fol = $scope.folder;
        var n = 0;

        _.forEach(files, function(file, i) {
          $upload.upload({
            url: '/u/' + fol.id,
            file: file
          }).progress(function(evt) {
            // var f = fol.files[l + i];
            var f = fol.files[file.i];
            f.loaded = Math.round(evt.loaded / evt.total * 100);
          }).success(function() {
            if (fol.files.length === ++n) {
              console.log('Files uploaded');
            }
          });
        });
      }
    }
  ]).controller('FolderCtrl', ['$scope', '$routeParams', 'Restangular',
    function($scope, $routeParams, Restangular) {
      console.log("Folder controller");

      Restangular.one('folders', $routeParams.id).get().then(function(fol) {
        _.assign($scope.folder, fol);
      });

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
    }
  ]);

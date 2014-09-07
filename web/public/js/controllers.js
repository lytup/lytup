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
      $scope.user = {};
      $scope.folders = [];

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

      $scope.signin = function(user) {
        return Restangular.all('users').all('login').post(user)
          .then(function(usr) {
            $window.localStorage.setItem('token', usr.token);
            $scope.user = usr;
            $scope.folders = Restangular.all('folders').getList().$object;
          });
      };

      $scope.signout = function() {
        $window.localStorage.removeItem('token');
        $scope.user = {};
        $location.path('/');
      };

      $scope.folderModal = function() {
        $modal.open({
          scope: $scope,
          controller: 'FolderModalCtrl',
          templateUrl: '/tpl/modals/folder.html',
          size: 'sm'
        });
      };

      $scope.deleteFolder = function(fol) {
        fol.remove().then(function() {
          _.remove($scope.folders, {
            'id': fol.id
          });
        });
      };

      /************
      / Validations
      /************/
      function validateFirstName(form, errors) {
        errors.blankFirstName = !form.firstName.$viewValue;
      };

      function validateLastName(form, errors) {
        errors.blankLastName = !form.lastName.$viewValue;
      };

      function validateEmail(form, errors) {
        errors.blankEmail = !form.email.$viewValue;
        errors.invalidEmail = form.email.$viewValue && form.email.$invalid;
      };

      function validatePassword(form, errors) {
        errors.blankPassword = !form.password.$viewValue;
        errors.invalidPassword = form.password.$viewValue && form.password.$invalid;
      };

      function validateExpiry(form, errors) {
        errors.blankExpiry = !form.expiry.$viewValue;
      };

      $scope.validate = function(form, errors) {
        if (form.firstName) {
          validateFirstName(form, errors);
        }
        if (form.lastName) {
          validateLastName(form, errors);
        }
        if (form.email) {
          validateEmail(form, errors);
        }
        if (form.password) {
          validatePassword(form, errors);
        }
        if (form.expiry) {
          validateExpiry(form, errors);
        }
        // Continue if form is pristine with expiry field
        if (form.$pristine && !form.expiry || form.$invalid) {
          return false;
        }
        return true;
      };
    }
  ])
  .controller('FolderCtrl', [
    '$scope',
    '$routeParams',
    '$log',
    '$modal',
    'Restangular',
    '$upload',
    function($scope, $routeParams, $log, $modal, Restangular, $upload) {
      $log.info('Folder controller');
      // Look into folders or get it from the server
      $scope.folder = _.find($scope.folders, {
        id: $routeParams.id
      }) || Restangular.one('folders', $routeParams.id).get().$object;

      $scope.folderModal = function() {
        $modal.open({
          scope: $scope,
          controller: 'FolderModalCtrl',
          templateUrl: '/tpl/modals/folder.html',
          size: 'sm'
        });
      };

      $scope.addFiles = function(files) {
        var fol = $scope.folder;

        // Concatenate in the same order to display on top
        fol.files = files.concat(fol.files)

        _.forEach(files, function(f, i) {
          var file = _.pick(f, 'name', 'size', 'type');

          // Create file
          fol.post('files', file).then(function(file) {
            // Replace file from the server
            fol.files[i] = file;
            upload(file, f);
          });
        });
      };

      $scope.fileIcon = function(typ) {
        return /image/.test(typ) ? 'fa-file-image-o' :
          /audio/.test(typ) ? 'fa-file-audio-o' :
          /video/.test(typ) ? 'fa-file-video-o' :
          /wordprocessingml/.test(typ) ? 'fa-file-word-o' :
          /spreadsheetml/.test(typ) ? 'fa-file-excel-o' :
          /presentationml/.test(typ) ? 'fa-file-powerpoint-o' :
          /pdf/.test(typ) ? 'fa-file-pdf-o' :
          /text/.test(typ) ? 'fa-file-text-o' :
          /zip/.test(typ) ? 'fa-file-archive-o' :
          '';
      };

      $scope.deleteFile = function(id) {
        $scope.folder.one('files', id).remove().then(function() {
          _.remove($scope.folder.files, {
            'id': id
          });
        });
      };

      function upload(file, f) {
        var fol = $scope.folder;

        $upload.upload({
          url: '/u/',
          file: f,
          data: {
            folId: fol.id,
            fileId: file.id
          }
        }).progress(function(evt) {
          file.loaded = Math.round(evt.loaded / evt.total * 100);
        }).success(function(f) {
          _.assign(file, _.omit(f, 'createdAt')); // https://code.google.com/p/go/issues/detail?id=5218
          // Update file
          file.patch(_.pick(file, 'loaded'));
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
    '$location',
    '$log',
    '$modal',
    'Restangular',
    function($scope, $window, $location, $log, $modal, Restangular) {
      $log.info('Signup controller');
      $scope.user = {};
      $scope.errors = {};

      var modal = $modal.open({
        scope: $scope,
        templateUrl: '/tpl/modals/signup.html',
        size: 'sm'
      });

      modal.result.finally(function() {
        $location.path('/');
      });

      $scope.submit = function() {
        Restangular.all('users').post($scope.user).then(function(usr) {
          $scope.signin($scope.user);
          modal.close();
        }, function(res) {
          if (res.data.error === 'duplicate') {
            $scope.errors.registeredEmail = true;
          }
        });
      };
    }
  ])
  .controller('SigninCtrl', [
    '$scope',
    '$window',
    '$location',
    '$log',
    '$modal',
    'Restangular',
    function($scope, $window, $location, $log, $modal, Restangular) {
      $log.info('Signin controller');
      $scope.user = {};
      $scope.errors = {};

      var modal = $modal.open({
        scope: $scope,
        templateUrl: '/tpl/modals/signin.html',
        size: 'sm'
      });

      modal.result.finally(function() {
        if ($location.path() !== '/forgot') {
          $location.path('/');
        }
      });

      $scope.forgotPassword = function() {
        modal.close();
        $location.path('/forgot');
      }

      $scope.submit = function() {
        $scope.signin($scope.user).then(function() {
          modal.close();
        }, function(res) {
          if (res.status === 404) {
            $scope.errors.loginFailed = true;
          }
        })
      }
    }
  ]).controller('FolderModalCtrl', [
    '$scope',
    '$location',
    '$log',
    '$modalInstance',
    'Restangular',
    function($scope, $location, $log, $modalInstance, Restangular) {
      $log.info('Folder modal controller');
      $scope.expiries = [{
        val: 1,
        lbl: '1 hour'
      }, {
        val: 4,
        lbl: '4 hours'
      }, {
        val: 24,
        lbl: '1 day'
      }, {
        val: 72,
        lbl: '3 days'
      }, {
        val: 168,
        lbl: '1 week'
      }];
      $scope.errors = {};

      $scope.save = function(folder) {
        if (!folder.id) {
          // Create
          Restangular.all('folders').post(folder).then(function(fol) {
            $scope.folders.push(fol);
            $modalInstance.close();
            $location.path('/' + fol.id);
          });
        } else {
          // Update
          folder.patch(_.pick(folder, 'name', 'expiry'))
            .then(function(fol) {
              _.assign(folder, _.pick(fol, 'name', 'expiry', 'expiresAt'))
              $modalInstance.close();
              $location.path('/' + folder.id);
            });
        }
      };
    }
  ]).controller('ConfirmCtrl', [
    '$scope',
    '$location',
    '$routeParams',
    '$log',
    'Restangular',
    function($scope, $location, $routeParams, $log, Restangular) {
      $log.info('Confirm controller');
      // Send confirmation request
      Restangular.all('users').one('confirm', $routeParams.key).get().then(function() {
        // $window.localStorage.setItem('token', usr.token);
        // $scope.user = usr;
        // $location.path('/');
      });
    }
  ]).controller('ForgotPwdCtrl', [
    '$scope',
    '$location',
    '$log',
    '$modal',
    'Restangular',
    function($scope, $location, $log, $modal, Restangular) {
      $log.info('Forgot password controller');
      $scope.user = {};
      $scope.errors = {};

      var modal = $modal.open({
        scope: $scope,
        templateUrl: '/tpl/modals/forgotpwd.html',
        size: 'sm'
      });

      modal.result.finally(function() {
        $location.path('/');
      });

      $scope.submit = function() {
        Restangular.all('users').all('forgot').post($scope.user).then(function() {
          // $window.localStorage.setItem('token', usr.token);
          // $scope.user = usr;
          // $location.path('/');
        });
      }
    }
  ]).controller('ResetPwdCtrl', [
    '$scope',
    '$window',
    '$location',
    '$routeParams',
    '$log',
    'Restangular',
    function($scope, $window, $location, $routeParams, $log, Restangular) {
      $log.info('Reset password controller');
      Restangular.all('users').one('reset', $routeParams.key).get().then(function(usr) {
        $window.localStorage.setItem('token', usr.token);
        _.assign($scope.user, usr);
        // $location.path('/');
      });
    }
  ])

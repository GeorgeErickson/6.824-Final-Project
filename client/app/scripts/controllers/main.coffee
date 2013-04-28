'use strict'

angular.module('clientApp')
  .controller 'MainCtrl', ($scope) ->
    $scope.rooms = [
      'HTML5 Boilerplate',
      'AngularJS',
      'Karma'
    ]

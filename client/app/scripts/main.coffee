#global require
"use strict"
require.config
  shim:
    underscore:
      exports: "_"

    backbone:
      deps: ["underscore", "jquery"]
      exports: "Backbone"

    bootstrap:
      deps: ["jquery"]
      exports: "jquery"

  paths:
    jquery: "../components/jquery/jquery"
    backbone: "../components/backbone-amd/backbone"
    underscore: "../components/underscore-amd/underscore"
    bootstrap: "vendor/bootstrap"

require ['router'], (router) ->
  window.app = new Thorax.LayoutView()
  app.appendTo '#app-main'
  app.router = router

  Backbone.history.start()

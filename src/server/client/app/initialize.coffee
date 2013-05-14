# Load App Helpers
require 'lib/helpers'
app = require 'app'

$ ->
  $.fn.editable.defaults.mode = 'inline'
  app.start()
  

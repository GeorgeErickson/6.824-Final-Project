Backbone.Marionette.Renderer.render = (template, data) ->
  template(data)

# Tell Swag where to look for partials
Swag.Config.partialsPath = '../views/templates/'

# Put your handlebars.js helpers here.

class DefaultNavView extends Backbone.Marionette.ItemView
  className: 'DefaultNav'

  template: require 'views/templates/DefaultNav'

module.exports = new DefaultNavView()
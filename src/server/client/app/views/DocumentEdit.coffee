module.exports = class DocumentEditView extends Backbone.Marionette.ItemView
  className: 'DocumentEdit'
  template: require 'views/templates/DocumentEdit'

  onRender: ->
    setTimeout =>
      @editor = ace.edit("editor")
    , 5

  
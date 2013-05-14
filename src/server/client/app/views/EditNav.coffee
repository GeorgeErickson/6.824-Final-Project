module.exports = class EditNavView extends Backbone.Marionette.ItemView
  className: 'EditNav'
  ui:
    title: "#title"


  initialize: ->
    @model.on 'change:Title', =>
      @render()

  onRender: =>
    @ui.title.editable
      emptytext: 'Untitled document'

    @ui.title.on 'save', (e, params) =>
      @model.save
        title: params.newValue

  template: require 'views/templates/EditNav'

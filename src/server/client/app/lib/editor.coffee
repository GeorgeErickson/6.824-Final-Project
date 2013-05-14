class Editor extends Backbone.View
  className: 'hide'
  render: ->
    @$el.append '<div id="editor"></div>'
    @editor = ace.edit("editor")
    @session = @editor.getSession()


  show: (model) ->
    snapshot = model.get 'Snapshot'
    @editor.setValue snapshot
    @$el.removeClass 'hide'

  hide: ->
    @$el.addClass 'hide'

editor = new Editor()
editor.$el.appendTo 'body'
editor.render()
module.exports = editor


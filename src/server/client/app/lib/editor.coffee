class Editor extends Backbone.View
  className: 'hide'
  render: ->
    @$el.append '<div id="editor"></div>'
    @editor = ace.edit("editor")

  show: ->
    @$el.removeClass 'hide'

  hide: ->
    @$el.addClass 'hide'

editor = new Editor()
editor.$el.appendTo 'body'
editor.render()
module.exports = editor


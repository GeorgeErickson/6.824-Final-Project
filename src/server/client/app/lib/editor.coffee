class Editor extends Backbone.View
  className: 'hide'
  render: ->
    @$el.append '<div id="editor"></div>'
    @editor = ace.edit("editor")

  show: (model) ->
    @editor.removeAllListeners 'change'
    snapshot = model.get 'Snapshot'
    @editor.setValue snapshot
    @editor.on 'change', @onchange
    @$el.removeClass 'hide'

  onchange: (e, editor) ->
    data = e.data
    console.log data

  hide: ->
    @$el.addClass 'hide'

editor = new Editor()
editor.$el.appendTo 'body'
editor.render()
module.exports = editor


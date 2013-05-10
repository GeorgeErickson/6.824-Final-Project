#Document
define [], () ->
  class Document extends Thorax.View 
    template: 'document'
    name: 'document'


    set_editor: (editor) ->
      @editor = editor
      @session = @editor.getSession()

  new Document()
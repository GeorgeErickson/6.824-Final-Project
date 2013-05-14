DocumentsListView = require 'views/DocumentsList'

describe 'DocumentsListView', ->
    beforeEach ->
        @view = new DocumentsListView()

    it 'should exist', ->
        expect(@view).to.be.ok

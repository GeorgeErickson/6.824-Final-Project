DocumentView = require 'views/Document'

describe 'DocumentView', ->
    beforeEach ->
        @view = new DocumentView()

    it 'should exist', ->
        expect(@view).to.be.ok

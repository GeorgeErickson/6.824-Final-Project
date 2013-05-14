DocumentsCollection = require 'collections/Documents'

describe 'DocumentsCollection', ->
    beforeEach ->
        @collection = new DocumentsCollection()

    it 'should exist', ->
        expect(@collection).to.be.ok

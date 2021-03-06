package store

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	collectionName = "collection"
	whoID          = 1
	whomeID        = 2
	anotherWhome   = 3

	simpleTest       = "simple test success"
	simpleDeleteTest = "simple delete test"
)

func Test_registry_AddToList(t *testing.T) {
	collection, _ := prepareCollection(collectionName)
	st := NewRegistry(collection)
	t.Run(simpleTest, func(t *testing.T) {
		err := st.AddEvent(Entry{Who: whoID, Whome: whomeID, Event: EventLike})
		assert.NoError(t, err)
		likes, err := st.GetEvents(Options{bson.E{"who", whoID}})
		assert.NoError(t, err)
		assert.Equal(t, int64(whoID), likes[0].Who)
		filter := bson.D{
			{
				"who", whoID,
			},
		}
		deleteRes, err := st.(*registry).collection.DeleteOne(context.TODO(), filter)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), deleteRes.DeletedCount)
	})

	t.Run(simpleTest, func(t *testing.T) {
		err := st.AddEvent(Entry{Who: whoID, Whome: whomeID, Event: EventLike})
		assert.NoError(t, err)
		err = st.AddEvent(Entry{Who: whoID, Whome: anotherWhome, Event: EventLike})
		assert.NoError(t, err)
		likes, err := st.GetEvents(Options{bson.E{"who", whoID}})
		assert.NoError(t, err)
		assert.Equal(t, int64(whoID), likes[0].Who)
		assert.Equal(t, int64(whomeID), likes[0].Whome)
		assert.Equal(t, int64(anotherWhome), likes[1].Whome)
		filter := bson.D{
			{
				"who", whoID,
			},
		}
		deleteRes, err := st.(*registry).collection.DeleteOne(context.TODO(), filter)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), deleteRes.DeletedCount)
	})
}

func TestRegistry_DeleteItem(t *testing.T) {
	collection, _ := prepareCollection(collectionName)
	reg := NewRegistry(collection)
	t.Run(simpleDeleteTest, func(t *testing.T) {
		err := reg.AddEvent(Entry{Who: 1, Whome: 2, Event: EventLike})
		assert.NoError(t, err)
		entry, err := reg.GetEvents(Options{bson.E{"who", 1}})
		assert.NoError(t, err)
		assert.NotZero(t, len(entry))
		//assert.Equal(t, 1, int(entry[0].Who))
		//assert.Equal(t, 2, int(entry[0].Whome))
		err = reg.DeleteEvents(Options{bson.E{"who", 1}})
		assert.NoError(t, err)
		entry, err = reg.GetEvents(Options{bson.E{"who", 1}})
		assert.Zero(t, len(entry))
		assert.Nil(t, entry)
	})
}

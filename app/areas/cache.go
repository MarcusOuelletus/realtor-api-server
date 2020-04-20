package areas

/*

BuildCache - queries all Residential Property listings, grabs their (locations) community, municipality and area
columns and puts all unique values in a slice called Slice.

OneLetterMap - every unique letter (made lower case) from all locations are added to this map as keys,
the value of each key is a slice filled with pointers to all values from Slice that contain the unique
letter in question.

i.e. index "a" contains a slice with oshawa, ancaster, bowmanville and every other location containing an "a"

TwoLetterMap - every unique pair of letters (made lower case) within a location string is add to this map
as a key. Each key maps to a slice of pointers to locations with Slice that container the letter pairing.

i.e. index "ca" contains a slice with newcastle, ancaster, etc.

-----------------------------------

There are approximately 1400 unique locations in the database, when a user searches a specific string,
if the string is one letter, it searches the OneLetterMap, this reduces the amount of comparisons by a
decent amount if the letter is less popular, like g or p.

If two letters or more are input then the system uses the TwoLetterMap for lookup, this significantly reduces
the amount of comparisons, in cases such as aj for ajax, comparisons go from 1400 to < 10.

*/

import (
	"context"
	"strings"

	"github.com/MarcusOuelletus/rets-server/database"
	"github.com/MarcusOuelletus/rets-server/templates"
	"github.com/golang/glog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type cacheData struct {
	Slice        []string
	OneLetterMap map[string][]*string
	TwoLetterMap map[string][]*string
}

var cacheObject = &cacheData{
	Slice:        nil,
	OneLetterMap: make(map[string][]*string),
	TwoLetterMap: make(map[string][]*string),
}

type cache struct {
	db *mongo.Database
}

type row struct {
	Area         string `bson:"Area"`
	Municipality string `bson:"Municipality"`
	Community    string `bson:"Community"`
}

// BuildCache - sets Slice of all unique locations as well as the one and two letter maps.
// this allows areas.QueryCache to perform much more efficiently than querying the database
// on every key stroke with the locations search box.
func BuildCache() error {
	var err error

	var c = new(cache)

	c.db, err = database.OpenDatabase()

	if err != nil {
		return err
	}

	defer c.db.Client().Disconnect(context.TODO())

	locationMap, err := c.createLocationMap()

	if err != nil {
		return err
	}

	c.populateCacheSlice(locationMap)

	return nil
}

func (c *cache) populateCacheSlice(m map[string]bool) {
	if cacheObject.Slice == nil {
		cacheObject.Slice = make([]string, len(m))
	}

	var index = -1

	for location := range m {
		index++

		// Add to cache slice
		cacheObject.Slice[index] = location
		delete(m, location)

		c.addLettersToCacheSingle(index)
		c.addLettersToCacheDouble(index)
	}
}

func (c *cache) addLettersToCacheSingle(indexOfString int) {
	lettersMap := make(map[rune]bool)

	for _, theRune := range cacheObject.Slice[indexOfString] {
		if _, ok := lettersMap[theRune]; ok {
			continue
		}

		lettersMap[theRune] = true

		letter := string(theRune)
		letter = strings.ToLower(letter)

		cacheObject.OneLetterMap[letter] = append(cacheObject.OneLetterMap[letter], &cacheObject.Slice[indexOfString])
	}
}

func (c *cache) addLettersToCacheDouble(indexOfString int) {
	pairMap := make(map[string]bool)

	var pair string

	location := cacheObject.Slice[indexOfString]

	for i := 0; i < len(location)-1; i++ {
		pair = string(location[i]) + string(location[i+1])
		pair = strings.ToLower(pair)

		if _, ok := pairMap[pair]; !ok {
			pairMap[pair] = true

			cacheObject.TwoLetterMap[pair] = append(cacheObject.TwoLetterMap[pair], &cacheObject.Slice[indexOfString])
		}
	}
}

func (c *cache) createLocationMap() (map[string]bool, error) {
	var m = make(map[string]bool)

	for _, collection := range [3]string{"res", "con", "com"} {
		connection := c.db.Collection(collection)

		cursor, err := c.getQueryCursor(connection)

		if err != nil {
			return nil, err
		}

		if err := c.addRowsToMap(cursor, m); err != nil {
			return nil, err
		}
	}

	return m, nil
}

func (c *cache) getQueryCursor(connection *mongo.Collection) (*mongo.Cursor, error) {
	var queryOptions = &options.FindOptions{
		Projection: bson.M{
			templates.Fields.Area:         1,
			templates.Fields.Community:    1,
			templates.Fields.Municipality: 1,
		},
	}

	cursor, err := connection.Find(context.Background(), primitive.M{}, queryOptions)

	if err != nil {
		glog.Errorln("error performing Select, Find() failed")
		return nil, err
	}

	return cursor, nil
}

func (c *cache) addRowsToMap(cursor *mongo.Cursor, m map[string]bool) error {
	for cursor.Next(context.TODO()) {
		var r = new(row)

		if err := cursor.Decode(r); err != nil {
			return err
		}

		for _, location := range [3]string{r.Area, r.Community, r.Municipality} {
			if location != "" {
				m[location] = true
			}
		}
	}

	return nil
}

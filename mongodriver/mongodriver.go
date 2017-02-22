package mongodriver

import (
	"fmt"
	"sync"

	"../logger"

	mgo "gopkg.in/mgo.v2-unstable"
	"gopkg.in/mgo.v2-unstable/bson"
)

type server struct {
	dbName         string
	collectionName string
	session        *mgo.Session
	m              *sync.Mutex
}

func init() {
	logger.LogStorageRegister("mongo", func(stHost string, stPort int, stName, appName string) (s logger.LogStorage, err error) {
		session, sessionErr := mgo.Dial(fmt.Sprintf("mongodb://%s:%d", stHost, stPort))
		if sessionErr != nil {
			err = sessionErr
			return
		}
		s = &server{
			dbName:         stName,
			collectionName: appName,
			session:        session,
			m:              &sync.Mutex{}}
		return
	})
}

func (s *server) Write(content *logger.LogContent) (err error) {
	s.m.Lock()
	defer s.m.Unlock()
	s.session.SetMode(mgo.Monotonic, true)
	c := s.session.DB(s.dbName).C(s.collectionName)
	err = c.Insert(content)
	if err != nil {
		fmt.Println(err)
	}
	return
}

//ReadLogALL ...
func (s *server) Read(conditions bson.M) ([]logger.LogContent, error) {
	s.session.SetMode(mgo.Monotonic, true)
	c := s.session.DB(s.dbName).C(s.collectionName)
	result := []logger.LogContent{}
	err := c.Find(conditions).All(&result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

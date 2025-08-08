package hue

const EVENT_BUFFER_SIZE = 25

// A Database stores Hue resource (e.g., light) data locally.
// A Database is intended to be synchronized with a Hue bridge.
type Database struct {
	conn      *BridgeConnection
	resources []Resource
	eventChan chan *eventContainer
}

func NewDatabase(conn *BridgeConnection) *Database {
	return &Database{
		conn:      conn,
		resources: make([]Resource, 0),
		eventChan: make(chan *eventContainer, EVENT_BUFFER_SIZE),
	}
}

// Initializes the Database with the latest data from the Hue bridge for all
// supported devices and sets up an event listener for detecting remote changes.
func (d *Database) Initialize() error {
	var err error

	err = d.fetchResources()
	if err != nil {
		return err
	}

	err = d.startEventListener()
	if err != nil {
		return err
	}

	return nil
}

// Pushes any changes made to the local database to the Hue bridge.
func (d *Database) PushChanges() error {
	for _, res := range d.resources {
		if err := res.SubmitChanges(d.conn); err != nil {
			return err
		}
	}

	return nil
}

// Pulls any events/changes from the Hue bridge and updates the local database.
func (d *Database) PullChanges() error {
	d.processEvents()
	return nil
}

// Retrieves Resources matching a given type.
func (d Database) GetResourcesByType(resType string) []Resource {
	matchingResources := make([]Resource, 0)

	for _, res := range d.resources {
		if res.Type() == resType {
			matchingResources = append(matchingResources, res)
		}
	}

	return matchingResources
}

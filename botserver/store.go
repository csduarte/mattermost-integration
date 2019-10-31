package botserver

import "github.com/csduarte/mattermost-integration/platform"

// const (
// 	//StorageLocation file to store integration data
// 	StorageLocation = "integrations.json"
// )
//
type store struct {
	Triggers []*Trigger
	clients  map[clientKey]*platform.Client
}

type clientKey struct {
	server   string
	user     string
	password string
}

//
// func restoreIntegrationStore(path string) (*integrationStore, error) {
// 	store := integrationStore{}
// 	store.clients = make(map[clientKey]*platform.Client)
// 	err := store.restore()
// 	if err != nil {
// 		switch {
// 		case os.IsNotExist(err):
// 			err = store.create()
// 		case os.IsPermission(err):
// 			return nil, fmt.Errorf("Cannot restore Integrations file - %v", err)
// 		default:
// 			return nil, fmt.Errorf("Unknown error - %v", err)
// 		}
// 	}
// 	return &store, err
// }
//
// func (is *integrationStore) restore() error {
// 	file, err := os.Open(StorageLocation)
// 	defer file.Close()
// 	if err != nil {
// 		return err
// 	}
// 	decoder := json.NewDecoder(file)
// 	err = decoder.Decode(&is)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }
//
// func (is *integrationStore) create() error {
// 	file, err := os.Create(StorageLocation)
// 	defer file.Close()
// 	if err != nil {
// 		return fmt.Errorf("Could not create store - %v", err)
// 	}
// 	if is.Integrations == nil {
// 		is.Integrations = []*Integration{}
// 	}
// 	encoder := json.NewEncoder(file)
// 	err = encoder.Encode(is)
// 	if err != nil {
// 		return fmt.Errorf("Count not encode store for save - %v", err)
// 	}
// 	return nil
// }
//
// func (is *integrationStore) save() error {
// 	file, err := os.OpenFile(StorageLocation, os.O_WRONLY, 0666)
// 	defer file.Close()
// 	if err != nil {
// 		return fmt.Errorf("Could not open store for save - %v", err)
// 	}
// 	if is.Integrations == nil {
// 		is.Integrations = []*Integration{}
// 	}
// 	encoder := json.NewEncoder(file)
// 	err = encoder.Encode(is)
// 	if err != nil {
// 		return fmt.Errorf("Count not encode store for save - %v", err)
// 	}
// 	return nil
// }
//
// func (is *integrationStore) matchIntegrations(isc Config) error {
// 	for _, ic := range isc.Integrations {
// 		if found := is.findByConfig(ic); found == nil {
// 			newInt, err := is.addIntegration(isc, *ic)
// 			if err != nil {
// 				return err
// 			}
// 			is.Integrations = append(is.Integrations, newInt)
// 		} else {
// 			c, err := is.clientForConfig(*ic)
// 			if err != nil {
// 				return err
// 			}
// 			found.initialize(c)
// 		}
// 	}
// 	is.save()
// 	return nil
// }
//
// func (is *integrationStore) addIntegration(isc Config, inc integrationConfig) (*Integration, error) {
// 	i := NewIntegrationFromConfig(isc, inc)
// 	c, err := is.clientForConfig(inc)
// 	if err != nil {
// 		return nil, fmt.Errorf("Failed to add webhooks for %q - %v", inc.Name, err)
// 	}
// 	fmt.Println("Initializing")
// 	err = i.initialize(c)
// 	if err != nil {
// 		return nil, fmt.Errorf("Failed to add webhooks for %q - %v", inc.Name, err)
// 	}
// 	return &i, nil
// }
//
// func (is *integrationStore) clientForConfig(inc integrationConfig) (*platform.Client, error) {
// 	key := clientKey{inc.Server, inc.Username, inc.Password}
// 	if c, ok := is.clients[key]; ok {
// 		return c, nil
// 	}
// 	c := platform.NewClient(inc.Server)
// 	_, err := c.Login(inc.Username, inc.Password)
// 	if err != nil {
// 		return nil, fmt.Errorf("Login error for %q - %v", inc.Name, err)
// 	}
// 	_, err = c.GetAllTeams()
// 	if err != nil {
// 		return nil, fmt.Errorf("Get Team Fails for %q - %v", inc.Name, err)
// 	}
// 	is.clients[key] = c
// 	return c, nil
// }
//
// func (is *integrationStore) findByConfig(c *integrationConfig) *Integration {
// 	for _, i := range is.Integrations {
// 		if i.Name == c.Name && i.ChatServer == c.Server {
// 			return i
// 		}
// 	}
// 	return nil
// }
//
// func (is *integrationStore) findByName(name string) *Integration {
// 	for _, i := range is.Integrations {
// 		if i.Name == name {
// 			return i
// 		}
// 	}
// 	return nil
// }
//
// func (is *integrationStore) findByToken(token string) *Integration {
// 	for _, i := range is.Integrations {
// 		if i.FromMM.Token == token {
// 			return i
// 		}
// 	}
// 	return nil
// }

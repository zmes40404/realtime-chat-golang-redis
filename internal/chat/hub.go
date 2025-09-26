 package chat

 import(
	"sync"
 )
 // Manage all rooms and client connections (central controller) -> Multi-person room support architecture


// Hub manage all the chat rooms
 type Hub struct {
	Rooms map[string]*Room 	// key of each room is room name and value is the memory address of Room
	mu sync.RWMutex		// Protect Rooms map
 }

 // Room is a struct of chat room(simple Ver)
 type Room struct {
	Name string
	Clients map[*Client] bool	// Check there any clients in the Room
	mu sync.RWMutex
 }

 func NewHub() *Hub {	// There will only be one Hub instance in the entire system (singleton pattern)
	return &Hub{
		Rooms: make(map[string]*Room),
	}
 }

 // GetOrCreateRoom: get or create a room
 func (h *Hub) GetOrCreateRoom(name string) *Room {
	h.mu.Lock()
	defer h.mu.Unlock()		// unlock before the func return the paramaters

	room, exists := h.Rooms[name]	// return mamory adress of room and bool value
	if !exists {	// if room of specific name isn't exist, create the Room with specific name 
		room = &Room{
			Name: name,
			Clients: make(map[*Client] bool),
		}
		h.Rooms[name] = room	// assign the memomry addr(value of map) to the map which key is equal to name
	}

	return room
 }
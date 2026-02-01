package main

import (
	"mymodule/struktury_i_funkcje"
	"sync"
	"time"
	"image/color"
	"math/rand"
	"math"

	_ "image/png"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
	"github.com/faiface/pixel/imdraw"
)

type EnemyType int

const (
	Fly EnemyType = iota
	Slime
)

type Enemy struct {
	Position pixel.Vec
	Type     EnemyType
	Health   int
	Active   bool
	LastShot time.Time
}

var sciany pixel.Rect
var itemPos pixel.Vec


func DrawMiniMap(win *pixelgl.Window) {
	im := imdraw.New(nil)

	roomSize := 18.0
	spacing := 26.0

	// pozycja minimapy
	origin := pixel.V(70, 950)

	// przesunięcie względem aktualnego pokoju
	offsetX := float64(currentRoom.X) * spacing
	offsetY := float64(currentRoom.Y) * spacing


	for _, room := range gameMap {
		x := origin.X + float64(room.X)*spacing - offsetX
		y := origin.Y + float64(room.Y)*spacing - offsetY

		// kolor pokoju
		switch room.Type {
		case StartRoom:
			im.Color = colornames.Dodgerblue
		case NormalRoom:
			im.Color = colornames.White
		case BossRoom:
			im.Color = colornames.Red
		case ItemRoom:
			im.Color = colornames.Gold
		}

		im.Push(
			pixel.V(x-roomSize/2, y-roomSize/2),
			pixel.V(x+roomSize/2, y+roomSize/2),
		)
		im.Rectangle(0)

		/* Pozycje drzwi – jak w Binding of Isaac
		doorPositions := map[string]pixel.Vec{
			"up":    pixel.V(960, sciany.Max.Y),
			"down":  pixel.V(960, sciany.Min.Y),
			"left":  pixel.V(sciany.Min.X, 560),
			"right": pixel.V(sciany.Max.X, 560),
		}*/

		// połączenia (drzwi)
		im.Color = colornames.Gray
		for dir := range room.Doors {
			switch dir {
			case "up":
				im.Push(pixel.V(x, y+roomSize/2), pixel.V(x, y+spacing-roomSize/2))
			case "down":
				im.Push(pixel.V(x, y-roomSize/2), pixel.V(x, y-spacing+roomSize/2))
			case "left":
				im.Push(pixel.V(x-roomSize/2, y), pixel.V(x-spacing+roomSize/2, y))
			case "right":
				im.Push(pixel.V(x+roomSize/2, y), pixel.V(x+spacing-roomSize/2, y))
			}
			im.Line(2)
		}

		// aktualny pokój – zielona ramka
		if room == currentRoom {
			im.Color = colornames.Limegreen
			im.Push(
				pixel.V(x-roomSize/2-2, y-roomSize/2-2),
				pixel.V(x+roomSize/2+2, y+roomSize/2+2),
			)
			im.Rectangle(2)
		}
	}

	im.Draw(win)
}


// Deklaracja lez bossa
var (
	bossTears []*struktury_i_funkcje.BossTear
	bossTearsMux sync.Mutex
)

// Deklaracja Husha
var (
	hushes []*struktury_i_funkcje.Hush
	hushMux sync.Mutex
	hushDrawn = false
	hushDead = false
	bossHush *struktury_i_funkcje.Hush
)

// Deklaracja Gurdiego
var (
	gurdys []*struktury_i_funkcje.Gurdy
	gurdyMux sync.Mutex
	gurdyRespawn = time.Now()
	gurdyDrawn = false
	gurdyDead = true
	bossGurdy *struktury_i_funkcje.Gurdy
)

var (
	enemies    []*Enemy
	enemiesMux sync.Mutex
)


var itemTaken = false
var itemCanBeTaken = true
var itemGeneratedThisFloor = false
var itemCollected = false
var floorItem = make(map[int]*pixel.Sprite)
var availableItems = []*pixel.Sprite{}
var roomItemTaken = make(map[*Room]bool)
var itemEffectApplied = make(map[*pixel.Sprite]bool)
var losowyItem *pixel.Sprite
var hasPiercing = false
var bossItemSpawned = false
var bossItemTaken = false

var roomJustEntered = false
var roomCleared = make(map[*Room]bool)
var roomVisited = make(map[*Room]bool)

var currentFloor = 1
var exitActive = false
var exitPos = pixel.V(960, 540)



func OnEnterRoom() {

	// bossowe pociski
	bossTears = nil
	bossGurdy = nil

	enemiesMux.Lock()
	enemies = nil
	enemiesMux.Unlock()

	// ---------- TREASURE ROOM ----------
	if currentRoom.Type == ItemRoom {
		itemTaken = roomItemTaken[currentRoom]
	
		itemPos = pixel.V(
			(sciany.Min.X+sciany.Max.X)/2,
			((sciany.Min.Y+sciany.Max.Y)/2) + 40,
		)
	}
	


	if currentRoom.Type == NormalRoom || currentRoom.Type == BossRoom {
		doorsLocked = !roomCleared[currentRoom]
	} else {
		doorsLocked = false
	}	

	// ===== BOSS ROOM =====
	if currentRoom.Type == BossRoom && !roomVisited[currentRoom] {

		exitActive = false
		bossTears = nil
	
		if currentFloor == 1 {
			gurdyMux.Lock()
			bossGurdy = &struktury_i_funkcje.Gurdy{
				Position: pixel.V(920, 520),
				Active:   true,
				Health:   300,
			}
			gurdys = []*struktury_i_funkcje.Gurdy{bossGurdy}
			gurdyDead = false
			gurdyDrawn = true
			gurdyMux.Unlock()
		}
	
		if currentFloor == 2 {
			hushMux.Lock()
			bossHush = &struktury_i_funkcje.Hush{
				Position: pixel.V(920, 520),
				Active:   true,
				Health:   400,
			}
			hushes = []*struktury_i_funkcje.Hush{bossHush}
			hushDead = false
			hushDrawn = true
			hushMux.Unlock()
		}
	}	

}



// ================= MAPA / POKOJE =================

type RoomType int

const (
	StartRoom RoomType = iota
	NormalRoom
	BossRoom
	ItemRoom
)

type Room struct {
	X, Y  int
	Type  RoomType
	Doors map[string]bool // "up","down","left","right"
}

var directions = map[string][2]int{
	"up":    {0, 1},
	"down":  {0, -1},
	"left":  {-1, 0},
	"right": {1, 0},
}

func opposite(dir string) string {
	switch dir {
	case "up":
		return "down"
	case "down":
		return "up"
	case "left":
		return "right"
	case "right":
		return "left"
	}
	return ""
}


func GenerateMap() map[[2]int]*Room {
	rand.Seed(time.Now().UnixNano())

	rooms := make(map[[2]int]*Room)

	start := &Room{
		X: 0, Y: 0,
		Type:  StartRoom,
		Doors: make(map[string]bool),
	}
	rooms[[2]int{0, 0}] = start

	frontier := []*Room{start}

	for len(rooms) < 8 {
		base := frontier[rand.Intn(len(frontier))]

		var possible []string
		for dir := range directions {
			dx, dy := directions[dir][0], directions[dir][1]
			pos := [2]int{base.X + dx, base.Y + dy}
			if _, exists := rooms[pos]; !exists {
				possible = append(possible, dir)
			}
		}

		if len(possible) == 0 {
			continue
		}

		dir := possible[rand.Intn(len(possible))]
		dx, dy := directions[dir][0], directions[dir][1]

		newRoom := &Room{
			X: base.X + dx,
			Y: base.Y + dy,
			Type:  NormalRoom,
			Doors: make(map[string]bool),
		}

		base.Doors[dir] = true
		newRoom.Doors[opposite(dir)] = true

		rooms[[2]int{newRoom.X, newRoom.Y}] = newRoom
		frontier = append(frontier, newRoom)
	}

	assignSpecialRooms(rooms)
	return rooms
}

func assignSpecialRooms(rooms map[[2]int]*Room) {
	var start *Room
	for _, r := range rooms {
		if r.Type == StartRoom {
			start = r
			break
		}
	}

	dist := make(map[*Room]int)
	queue := []*Room{start}
	dist[start] = 0

	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]

		for dir := range cur.Doors {
			dx, dy := directions[dir][0], directions[dir][1]
			n := rooms[[2]int{cur.X + dx, cur.Y + dy}]
			if _, ok := dist[n]; !ok {
				dist[n] = dist[cur] + 1
				queue = append(queue, n)
			}
		}
	}

	var boss *Room
	maxDist := -1
	for r, d := range dist {
		if d > maxDist && r.Type == NormalRoom {
			maxDist = d
			boss = r
		}
	}
	boss.Type = BossRoom

	for _, r := range rooms {
		if r.Type == NormalRoom && len(r.Doors) == 1 {
			r.Type = ItemRoom
			break
		}
	}
}

func ChangeRoom(dir string) {
	if !currentRoom.Doors[dir] {
		return
	}

	dx, dy := directions[dir][0], directions[dir][1]
	currentRoom = gameMap[[2]int{currentRoom.X + dx, currentRoom.Y + dy}]

	OnEnterRoom()
	roomJustEntered = true

	if !roomVisited[currentRoom] {
		SpawnEnemies()
		roomVisited[currentRoom] = true
	}
}

func NextFloor() {

	roomItemTaken = make(map[*Room]bool)
	currentFloor++

	if currentFloor > 2 {
		return // koniec gry
	}

	// reset mapy
	gameMap = GenerateMap()
	currentRoom = gameMap[[2]int{0, 0}]

	// reset stanów
	roomVisited = make(map[*Room]bool)
	roomCleared = make(map[*Room]bool)
	enemies = nil
	bossTears = nil

	// reset itemów
	itemTaken = false
	itemCollected = false

	// reset bossów
	gurdyDead = true
	hushDead = true
	gurdyDrawn = false
	hushDrawn = false

	exitActive = false

	AssignFloorItem()
	losowyItem = floorItem[currentFloor]
	OnEnterRoom()
}


func DrawDoors(win *pixelgl.Window, room *Room,
	drzwiGora, drzwiDol, drzwiLewo, drzwiPrawo *pixel.Sprite) {

	for dir, exists := range room.Doors {
		if !exists {
			continue
		}
		
		switch dir {
		case "up":
			drzwiGora.Draw(win, pixel.IM.Moved(
				pixel.V(940, sciany.Max.Y+18),
			))
		case "down":
			drzwiDol.Draw(win, pixel.IM.Moved(
				pixel.V(940, sciany.Min.Y-105),
			))
		case "left":
			drzwiLewo.Draw(win, pixel.IM.Moved(
				pixel.V(sciany.Min.X-80, 540),
			))
		case "right":
			drzwiPrawo.Draw(win, pixel.IM.Moved(
				pixel.V(sciany.Max.X+80, 540),
			))
		}
	}
}


func SpawnEnemies() {
	if currentRoom.Type != NormalRoom {
		return
	}

	count := rand.Intn(4) + 3 // 3–6

	for i := 0; i < count; i++ {
		t := Fly
		if rand.Intn(2) == 0 {
			t = Slime
		}

		x := rand.Float64()*(sciany.Max.X-sciany.Min.X) + sciany.Min.X
		y := rand.Float64()*(sciany.Max.Y-sciany.Min.Y) + sciany.Min.Y

		e := &Enemy{
			Position: pixel.V(x, y),
			Type:     t,
			Health:   20,
			Active:   true,
			LastShot: time.Now(),
		}

		enemies = append(enemies, e)
	}
}

var doorsLocked bool

func IsRoomCleared() bool {
	if currentRoom.Type == BossRoom {
		return gurdyDead && hushDead
	}

	for _, e := range enemies {
		if e.Active {
			return false
		}
	}
	return true
}



// ===== MAPA =====
var (
	gameMap     map[[2]int]*Room
	currentRoom *Room
)

func AssignFloorItem() {
	if floorItem[currentFloor] != nil {
		return
	}

	if len(availableItems) == 0 {
		return
	}

	index := rand.Intn(len(availableItems))
	floorItem[currentFloor] = availableItems[index]

	// usuń z puli globalnej
	availableItems = append(
		availableItems[:index],
		availableItems[index+1:]...,
	)
}

func SpriteHitbox(pos pixel.Vec, sprite *pixel.Sprite, scale float64) pixel.Rect {
	frame := sprite.Frame()

	halfW := frame.W() * scale / 2
	halfH := frame.H() * scale / 2

	return pixel.R(
		pos.X-halfW, pos.Y-halfH,
		pos.X+halfW, pos.Y+halfH,
	)
}




func run() {	

	gameMap = GenerateMap()
	currentRoom = gameMap[[2]int{0, 0}]
	itemGeneratedThisFloor = false

	// Wczytywanie okna
	cfg := pixelgl.WindowConfig{
		Title:     "Isaac",
		Bounds:    pixel.R(0, 0, 1920, 1080),
		VSync:     true,
		Monitor:   nil,
		Resizable: true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	// Wczytywanie tła
	pic1 := struktury_i_funkcje.LoadPicture("tlo.png")
	tlo := pixel.NewSprite(pic1, pic1.Bounds())

	// Wczytywanie postaci w stanie podstawowym/idacej w dol
	pic2 := struktury_i_funkcje.LoadPicture("w_dol.png")
	w_dol := pixel.NewSprite(pic2, pic2.Bounds())

	// Wczytywanie postaci idacej do przodu
	pic3 := struktury_i_funkcje.LoadPicture("do_przodu.png")
	do_przodu := pixel.NewSprite(pic3, pic3.Bounds())

	// Wczytywanie postaci idacej w prawo
	pic4 := struktury_i_funkcje.LoadPicture("w_prawo.png")
	w_prawo := pixel.NewSprite(pic4, pic4.Bounds())

	// Wczytywanie postaci idacej w lewo
	pic5 := struktury_i_funkcje.LoadPicture("w_lewo.png")
	w_lewo := pixel.NewSprite(pic5, pic5.Bounds())

	// Wczytywanie lzy
	pic6 := struktury_i_funkcje.LoadPicture("lza.png")
	lza := pixel.NewSprite(pic6, pic6.Bounds())

	// Wczytywanie gurdiego
	pic7 := struktury_i_funkcje.LoadPicture("gurdy.png")
	gurdy := pixel.NewSprite(pic7, pic7.Bounds())

	// Wczytywanie serca
	pic8 := struktury_i_funkcje.LoadPicture("serce.png")
	serce := pixel.NewSprite(pic8, pic8.Bounds())

	// Wczytywanie paska zycia bossa
	pic9 := struktury_i_funkcje.LoadPicture("pasek.png")
	pasek := pixel.NewSprite(pic9, pic9.Bounds())

	// Wczytywanie drzwi
	pic10 := struktury_i_funkcje.LoadPicture("drzwiGora.png")
	drzwiGora := pixel.NewSprite(pic10, pic10.Bounds())
	pic17 := struktury_i_funkcje.LoadPicture("drzwiDol.png")
	drzwiDol := pixel.NewSprite(pic17, pic17.Bounds())
	pic18 := struktury_i_funkcje.LoadPicture("drzwiLewo.png")
	drzwiLewo := pixel.NewSprite(pic18, pic18.Bounds())
	pic19 := struktury_i_funkcje.LoadPicture("drzwiPrawo.png")
	drzwiPrawo := pixel.NewSprite(pic19, pic19.Bounds())

	// Wczytywanie drugiego bossa
	pic11 := struktury_i_funkcje.LoadPicture("hush.png")
	hush := pixel.NewSprite(pic11, pic11.Bounds())

	// Wczytywanie piedestalu
	pic12 := struktury_i_funkcje.LoadPicture("piedestal.png")
	piedestal := pixel.NewSprite(pic12, pic12.Bounds())

	// Wczytywanie soku
	pic13 := struktury_i_funkcje.LoadPicture("tearsUp.png")
	sok := pixel.NewSprite(pic13, pic13.Bounds())

	// Wczytywanie strzykawki
	pic14 := struktury_i_funkcje.LoadPicture("dmgUp.png")
	strzykawka := pixel.NewSprite(pic14, pic14.Bounds())

	// Wczytywanie buta
	pic15 := struktury_i_funkcje.LoadPicture("rangeUp.png")
	but := pixel.NewSprite(pic15, pic15.Bounds())

	// Wczytywanie miesa
	pic16 := struktury_i_funkcje.LoadPicture("healthUp.png")
	mieso := pixel.NewSprite(pic16, pic16.Bounds())

	// Wczytywanie przejscia
	pic20 := struktury_i_funkcje.LoadPicture("drop.png")
	drop := pixel.NewSprite(pic20, pic20.Bounds())

	// Wczytywanie potworow
	pic21 := struktury_i_funkcje.LoadPicture("mucha.png")
	mucha := pixel.NewSprite(pic21, pic21.Bounds())
	pic22 := struktury_i_funkcje.LoadPicture("follower.png")
	follower := pixel.NewSprite(pic22, pic22.Bounds())
	muchaScale := 0.42
	followerScale := 1.0

	// Wczytywanie pierca
	pic23 := struktury_i_funkcje.LoadPicture("pierceUp.png")
	pierceUp := pixel.NewSprite(pic23, pic23.Bounds())
	pic25 := struktury_i_funkcje.LoadPicture("pierceDown.png")
	pierceDown := pixel.NewSprite(pic25, pic25.Bounds())
	pic26 := struktury_i_funkcje.LoadPicture("pierceRight.png")
	pierceRight := pixel.NewSprite(pic26, pic26.Bounds())
	pic27 := struktury_i_funkcje.LoadPicture("pierceLeft.png")
	pierceLeft := pixel.NewSprite(pic27, pic27.Bounds())
	pic24 := struktury_i_funkcje.LoadPicture("arrow.png")
	arrow := pixel.NewSprite(pic24, pic24.Bounds())

	pos := pixel.V(960, 340)
	poslza := pos
	speed := 7.5
	MoveUp := false
	MoveDown := false
	MoveRight := false
	MoveLeft := false


	// Ograniczenie scian dla postaci
	sciany = pixel.R(275, 210, 1585, 935)
	// Ograniczenie scian dla lez
	scianyLzy := pixel.R(275, 210, 1585, 935)
	// Ograniczenie scian dla znikania lez
	scianyLzyOff := pixel.R(240, 153, 1620, 935)
	// Ograniczenie ruchu postaci dla bossa
	bossGurdySize := pixel.R(636, 557, 1245, 735)
	// Ograniczenie ruchu lez dla bossa
	bossGurdySizeLzy := pixel.R(636, 557, 1245, 935)
	// Kierunek ataku gurdiego
	//attackLeft := pixel.R(275, 210, 636, 935)
	//attackMiddle := pixel.R(637, 210, 1244, 557)
	//attackRight := pixel.R(1245, 210, 1585, 935)
	// Czarny porstokat na pasek zycia bossa
	czarneZycie := pixel.R(705, 87, 1180, 127)
	// Czerwony porstokat na pasek zycia bossa
	czerwoneZycie := pixel.R(705, 87, 1180, 127)
	// Czarny porstokat na pasek zycia drugiego bossa
	czarneZycie2 := pixel.R(705, 87, 1180, 127)
	// Czerwony porstokat na pasek zycia drugiego bossa
	czerwoneZycie2 := pixel.R(705, 87, 1180, 127)
	// Kwadrat na drzwi
	doorUp := pixel.R(820, sciany.Max.Y-5, 1100, sciany.Max.Y+20)
	doorDown := pixel.R(820, sciany.Min.Y-20, 1100, sciany.Min.Y+5)
	doorLeft := pixel.R(sciany.Min.X-20, 480, sciany.Min.X+5, 650)
	doorRight := pixel.R(sciany.Max.X-5, 480, sciany.Max.X+20, 650)

	// Ograniczenie dla drugiego bossa
	bossHushSize := pixel.R(716, 407, 1136, 707)
	// Piedestal
	piedestalHitbox := pixel.R(880, 540, 980, 640)
	wziacItem := pixel.R(870, 530, 990, 650)

	availableItems = []*pixel.Sprite{
		sok,
		strzykawka,
		but,
		mieso,
	}
	AssignFloorItem()
	losowyItem = floorItem[currentFloor]


	gameMap = GenerateMap()
	currentRoom = gameMap[[2]int{0, 0}]

	if doorsLocked && IsRoomCleared() {
		doorsLocked = false
	}	

	// Kolekcja lez
	var (
		tears []*struktury_i_funkcje.Tear
		tearsMux sync.Mutex
		tear *struktury_i_funkcje.Tear
		sleep time.Duration
		damage int
		lifetime time.Duration
	)

	// Losowanie itemu po bossie
	var (
		itemMux sync.Mutex
		plansza1 = true
		isaacHp = 3
	)

	/*items := []*pixel.Sprite{
		sok,
		strzykawka,
		but,
		mieso,
	}

	if floorItem[currentFloor] == nil {
		index := rand.Intn(len(items))
		floorItem[currentFloor] = items[index]
	}
		
	losowyItem = floorItem[currentFloor]*/

	
	// Funkcja uzywajaca itemow
	
	go func() {
		for {

			itemMux.Lock()
			if losowyItem == mieso && itemTaken && !itemEffectApplied[losowyItem] {
				isaacHp++
				itemEffectApplied[losowyItem] = true
			}
			if losowyItem == sok && itemTaken {
				sleep = 250*time.Millisecond
			} else {
				sleep = 333*time.Millisecond
			}
			if losowyItem == but && itemTaken {
				lifetime = 2000*time.Millisecond
			} else {
				lifetime = 1300*time.Millisecond
			}
			if losowyItem == strzykawka && itemTaken {
				damage = 15.0
			} else {
				damage = 10.0
			}
			itemMux.Unlock()
		}
	}()

	// Funkcja generujaca lzy
	go func() {
		for {

			if roomJustEntered {
				roomJustEntered = false
				time.Sleep(50 * time.Millisecond)
				continue
			}

			if win.Pressed(pixelgl.KeyUp) {
				tearsMux.Lock()

				tear = &struktury_i_funkcje.Tear{
					Position: pos,
					Velocity: pixel.V(0, 500),
					StartTime: time.Now(),
					Lifetime: lifetime,
					Active: true,
					Damage: damage,
					Direction: "up",
				}				
				tears = append(tears, tear)
				tearsMux.Unlock()
			} else if win.Pressed(pixelgl.KeyDown) {
				tearsMux.Lock()

				tear = &struktury_i_funkcje.Tear {
					Position: pos,
					Velocity: pixel.V(0, -500),
					StartTime: time.Now(),
					Lifetime: lifetime,
					Active: true,
					Damage: damage,
					Direction: "down",
				}
				tears = append(tears, tear)
				tearsMux.Unlock()
			} else if win.Pressed(pixelgl.KeyRight) {
				tearsMux.Lock()

				tear = &struktury_i_funkcje.Tear {
					Position: pos,
					Velocity: pixel.V(500, 0),
					StartTime: time.Now(),
					Lifetime: lifetime,
					Active: true,
					Damage: damage,
					Direction: "right",
				}
				tears = append(tears, tear)
				tearsMux.Unlock()
			} else if win.Pressed(pixelgl.KeyLeft) {
				tearsMux.Lock()

				tear = &struktury_i_funkcje.Tear {
					Position: pos,
					Velocity: pixel.V(-500, 0),
					StartTime: time.Now(),
					Lifetime: lifetime,
					Active: true,
					Damage: damage,
					Direction: "left",
				}
				tears = append(tears, tear)
				tearsMux.Unlock()
			}
			time.Sleep(sleep)
		}
	}()


	// Funkcja odejmujaca hp gurdiego
	var ( 
		gurdyHpMux sync.Mutex
		gurdyGotHit = false
	)

	go func() {
		for {
			if gurdyGotHit {
				gurdyHpMux.Lock()
				bossGurdy.Health = bossGurdy.Health - tear.Damage
				czerwoneZycie.Max.X = czerwoneZycie.Max.X - 15
				gurdyGotHit = false
				gurdyHpMux.Unlock()
			}
			time.Sleep(300*time.Millisecond)
		}
	}()

	// Zmiana planszy
	var (
		zmiana = true
		zmianaMux sync.Mutex
	)

	go func() {
		for {
			zmianaMux.Lock()
			zmiana = currentRoom.Type != BossRoom || gurdyDead
			

			if zmiana && plansza1 {
				if doorUp.Contains(pos) && currentRoom.Doors["up"] && !doorsLocked {
					ChangeRoom("up")
					pos = pixel.V(960, sciany.Min.Y+40)
				}
				if doorDown.Contains(pos) && currentRoom.Doors["down"] && !doorsLocked {
					ChangeRoom("down")
					pos = pixel.V(960, sciany.Max.Y-40)
				}
				if doorLeft.Contains(pos) && currentRoom.Doors["left"] && !doorsLocked {
					ChangeRoom("left")
					pos = pixel.V(sciany.Max.X-40, 560)
				}
				if doorRight.Contains(pos) && currentRoom.Doors["right"] && !doorsLocked {
					ChangeRoom("right")
					pos = pixel.V(sciany.Min.X+40, 560)
				}
			}
			
			zmianaMux.Unlock()
		}
	}()
	
	// Funkcja strzelajaca lzami bossa
	go func() {
		for {
			if gurdyDead || !gurdyDrawn || currentRoom.Type != BossRoom {
				time.Sleep(100 * time.Millisecond)
				continue
			}
	
			bossTearsMux.Lock()
	
			// kierunek do gracza
			dir := pos.Sub(bossGurdy.Position)
			if dir.Len() == 0 {
				dir = pixel.V(0, -1)
			}
			dir = dir.Unit()
	
			// parametry fali
			bulletSpeed := 350.0
			bulletCount := 7
			spread := 0.6 // im większe, tym szerszy wachlarz
	
			for i := 0; i < bulletCount; i++ {
				offset := float64(i-(bulletCount/2)) * spread
				angle := math.Atan2(dir.Y, dir.X) + offset
	
				vel := pixel.V(
					math.Cos(angle),
					math.Sin(angle),
				).Scaled(bulletSpeed)
	
				bossTears = append(bossTears, &struktury_i_funkcje.BossTear{
					Position: bossGurdy.Position,
					Velocity: vel,
					Active:   true,
				})
			}
	
			bossTearsMux.Unlock()
			time.Sleep(900 * time.Millisecond)
		}
	}()
	

	// Funkcja odejmujaca hp postaci
	var (
		isaacHpMux sync.Mutex
		isaacGotHit = false
		isaacDead = false
	)

	go func() {
		for {
			if isaacGotHit {
				isaacHpMux.Lock()
				isaacHp = isaacHp - 1
				isaacGotHit = false
				if isaacHp == 0 {
					isaacDead = true
				}
				isaacHpMux.Unlock()
			}
			time.Sleep(500*time.Millisecond)
		}
	}()


	// Funkcja odejmujaca hp Husha
	var ( 
		hushHpMux sync.Mutex
		hushGotHit = false
	)

	go func() {
		for {
			if hushGotHit {
				hushHpMux.Lock()
				bossHush.Health = bossHush.Health - tear.Damage
				if losowyItem == strzykawka && itemTaken {
					czerwoneZycie2.Max.X = czerwoneZycie2.Max.X - 17
				} else {
					czerwoneZycie2.Max.X = czerwoneZycie2.Max.X - 11
				}
				hushGotHit = false
				hushHpMux.Unlock()
			}
			time.Sleep(300*time.Millisecond)
		}
	}()

	// Funkcja losujaca punkt w obrebie planszy
	rand.Seed(time.Now().UnixNano())

	go func() {
		for {
			if !hushDead && hushDrawn {
				bossTearsMux.Lock()
				x := rand.Float64()*(sciany.Max.X-sciany.Min.X) + sciany.Min.X
				y := rand.Float64()*(sciany.Max.Y-sciany.Min.Y) + sciany.Min.Y
				bossTear := &struktury_i_funkcje.BossTear {
					Position: pixel.V(x, y),
					Velocity: pixel.V(-100, 400),
					Active: true,
				}
				bossTears = append(bossTears, bossTear)
				bossTear = &struktury_i_funkcje.BossTear {
					Position: pixel.V(x, y),
					Velocity: pixel.V(100, 400),
					Active: true,
				}
				bossTears = append(bossTears, bossTear)
				bossTear = &struktury_i_funkcje.BossTear {
					Position: pixel.V(x, y),
					Velocity: pixel.V(200, 300),
					Active: true,
				}
				bossTears = append(bossTears, bossTear)
				bossTear = &struktury_i_funkcje.BossTear {
					Position: pixel.V(x, y),
					Velocity: pixel.V(400, 100),
					Active: true,
				}
				bossTears = append(bossTears, bossTear)
				bossTear = &struktury_i_funkcje.BossTear {
					Position: pixel.V(x, y),
					Velocity: pixel.V(400, -100),
					Active: true,
				}
				bossTears = append(bossTears, bossTear)
				bossTear = &struktury_i_funkcje.BossTear {
					Position: pixel.V(x, y),
					Velocity: pixel.V(200, -300),
					Active: true,
				}
				bossTears = append(bossTears, bossTear)
				bossTear = &struktury_i_funkcje.BossTear {
					Position: pixel.V(x, y),
					Velocity: pixel.V(-100, -400),
					Active: true,
				}
				bossTears = append(bossTears, bossTear)
				bossTear = &struktury_i_funkcje.BossTear {
					Position: pixel.V(x, y),
					Velocity: pixel.V(100, -400),
					Active: true,
				}
				bossTears = append(bossTears, bossTear)
				bossTear = &struktury_i_funkcje.BossTear {
					Position: pixel.V(x, y),
					Velocity: pixel.V(-200, -300),
					Active: true,
				}
				bossTears = append(bossTears, bossTear)
				bossTear = &struktury_i_funkcje.BossTear {
					Position: pixel.V(x, y),
					Velocity: pixel.V(-400, -100),
					Active: true,
				}
				bossTears = append(bossTears, bossTear)
				bossTear = &struktury_i_funkcje.BossTear {
					Position: pixel.V(x, y),
					Velocity: pixel.V(-400, 100),
					Active: true,
				}
				bossTears = append(bossTears, bossTear)
				bossTear = &struktury_i_funkcje.BossTear {
					Position: pixel.V(x, y),
					Velocity: pixel.V(-200, 300),
					Active: true,
				}
				bossTears = append(bossTears, bossTear)
				bossTearsMux.Unlock()
			}
			time.Sleep(1500*time.Millisecond)
		}
	}()



	last := time.Now()

	for !win.Closed() {
		dt := time.Since(last).Seconds()
		last = time.Now()

		// Obsługa poruszania się postaci
		anim := w_dol

		if win.Pressed(pixelgl.KeyW) {
			MoveUp = true
			MoveDown = false
		} else if win.Pressed(pixelgl.KeyS) {
			MoveUp = false
			MoveDown = true
		} else {
			MoveUp = false
			MoveDown = false
		}

		if win.Pressed(pixelgl.KeyD) {
			MoveRight = true
			MoveLeft = false
		} else if win.Pressed(pixelgl.KeyA) {
			MoveRight = false
			MoveLeft = true
		} else {
			MoveRight = false
			MoveLeft = false
		}

		enemiesMux.Lock()

		for _, e := range enemies {
			if !e.Active {
				continue
			}

			// --- HITBOXY ---
			var hitbox pixel.Rect

			switch e.Type {
			case Fly:
				// 🟣 mucha – hitbox = sprite mucha
				hitbox = SpriteHitbox(e.Position, mucha, muchaScale)

			case Slime:
				// 🟢 follower – hitbox = sprite follower
				hitbox = SpriteHitbox(e.Position, follower, followerScale)

				// kolizja gracza z followerem
				if hitbox.Contains(pos) {
					isaacGotHit = true
				}
			}


			// (opcjonalnie) kolizja łez z enemy
			for _, tear := range tears {
				if tear.Active && hitbox.Contains(tear.Position) {
					e.Health -= tear.Damage
					if e.Health <= 0 {
						e.Active = false
					}
			
					// TYLKO jeśli NIE ma piercingu – łza znika
					if !hasPiercing {
						tear.Active = false
					}
				}
			}			
		}
		enemiesMux.Unlock()

		

		// Obsluga animacji przy chodzeniu
		if MoveLeft {
			anim = w_lewo
		} else if MoveRight {
			anim = w_prawo
		} else if MoveUp {
			anim = do_przodu
		} else if MoveDown {
			anim = w_dol
		}

		// Obsluga animacji przy strzelaniu
		if win.Pressed(pixelgl.KeyUp) {
			anim = do_przodu
		} else if win.Pressed(pixelgl.KeyDown) {
			anim = w_dol
		} else if win.Pressed(pixelgl.KeyRight) {
			anim = w_prawo
		} else if win.Pressed(pixelgl.KeyLeft) {
			anim = w_lewo
		}

		// Aktualizacja pozycji postaci
		movement := pixel.ZV

		if MoveUp && !MoveDown {
			movement.Y += 1
		} else if MoveDown && !MoveUp {
			movement.Y -= 1
		}

		if MoveRight && !MoveLeft {
			movement.X += 1
		} else if MoveLeft && !MoveRight {
			movement.X -= 1
		}

		// Normalizacja ruchu diagonalnego
		if movement.Len() > 0 {
			movement = movement.Unit()
		}

		// Aktualizacja pozycji postaci przy zderzeniu ze sciana
		oldPos := pos
		NewPos := pos.Add(movement.Scaled(speed))

		if sciany.Contains(NewPos) {
			pos = NewPos
			poslza = pos
		} else {
			if NewPos.X < sciany.Min.X {
         	   NewPos.X = sciany.Min.X
        	} else if NewPos.X > sciany.Max.X {
        		NewPos.X = sciany.Max.X
        	}
        	if NewPos.Y < sciany.Min.Y {
            	NewPos.Y = sciany.Min.Y
        	} else if NewPos.Y > sciany.Max.Y {
            	NewPos.Y = sciany.Max.Y
        	}
			
			pos = NewPos
		}

		// ===== WEJŚCIE DO WYJŚCIA (NASTĘPNE PIĘTRO) =====
		exitHitbox := pixel.R(
			exitPos.X-40, exitPos.Y-40,
			exitPos.X+40, exitPos.Y+40,
		)

		if exitActive && exitHitbox.Contains(pos) {
			NextFloor()
			roomJustEntered = true
			pos = pixel.V(960, 340)
		}


		// Aktualizacja pozycji postaci przy zderzeniu z bossem
		NewPos = pos
		
		if !gurdyDead {
			if bossGurdySize.Contains(NewPos) {
				pos = oldPos
				isaacHpMux.Lock()
				isaacGotHit = true
				isaacHpMux.Unlock()
			}
		}

		if gurdyDead && plansza1 && piedestalHitbox.Contains(NewPos) && currentRoom.Type == ItemRoom {
			pos = oldPos
		}
		if currentRoom.Type == ItemRoom && !itemTaken {
			if wziacItem.Contains(pos) && win.JustPressed(pixelgl.KeyE) {
				itemTaken = true
				itemCollected = true
				roomItemTaken[currentRoom] = true

			}
		}

		// ===== PODNOSZENIE ITEMU PO BOSSIE =====
		if bossItemSpawned && !bossItemTaken {
			bossItemHitbox := SpriteHitbox(itemPos, arrow, 1.0)

			if bossItemHitbox.Contains(pos) && win.JustPressed(pixelgl.KeyE) {
				bossItemTaken = true
				hasPiercing = true
			}
		}		
		
		if !hushDead && hushDrawn {
			if bossHushSize.Contains(NewPos) {
				pos = oldPos
				isaacHpMux.Lock()
				isaacGotHit = true
				isaacHpMux.Unlock()
			}
		}
		// Aktualizacja pozycji lez przy zderzeniu ze sciana
		oldPosLzy := poslza
		NewPosLza := poslza.Add(movement.Scaled(speed))

		if scianyLzy.Contains(NewPosLza) {
			poslza = NewPosLza
		} else {
			if NewPosLza.X < scianyLzy.Min.X {
           		NewPosLza.X = scianyLzy.Min.X
        	} else if NewPosLza.X > scianyLzy.Max.X {
            	NewPosLza.X = scianyLzy.Max.X
        	}
        	if NewPosLza.Y < scianyLzy.Min.Y {
            NewPosLza.Y = scianyLzy.Min.Y
        	} else if NewPosLza.Y > scianyLzy.Max.Y {
            	NewPosLza.Y = scianyLzy.Max.Y
       		}
			poslza = NewPosLza
		}
		
		// Aktualizacja pozycji lez przy zderzeniu z bossem
		NewPosLza = poslza

		if currentRoom.Type == BossRoom && !gurdyDead {
			if bossGurdySize.Contains(NewPos) {
				pos = oldPos
				isaacGotHit = true
			}
		}
		

		if gurdyDead && plansza1 && piedestalHitbox.Contains(NewPosLza) {
			poslza = oldPosLzy
		}
		
		if !hushDead && hushDrawn {
			if bossHushSize.Contains(NewPosLza) {
				poslza = oldPosLzy
			}
		}

		win.Clear(color.Black)

		// Rysowanie tła
		zmianaMux.Lock()
		if plansza1 {
			tlo.Draw(win, pixel.IM.Moved(pixel.V(930, 525)))
		} else {
			return	//tlo2.Draw(win, pixel.IM.Moved(pixel.V(930, 525)))
		}
		zmianaMux.Unlock()

		DrawDoors(win, currentRoom,
			drzwiGora, drzwiDol, drzwiLewo, drzwiPrawo)
		

		// Rysowanie Gurdiego
		gurdyMux.Lock()
		
		for _, bossGurdy := range gurdys {
    		if bossGurdy.Active && gurdyDrawn && !gurdyDead {
        		if bossGurdy.Health <= 0 && !gurdyDead {
        		    gurdyDead = true
        		}
        		gurdy.Draw(win, pixel.IM.Moved(bossGurdy.Position))
   			}
		}
	
		gurdyMux.Unlock()

		// ===== EXIT + ITEM PO BOSSIE =====
		if currentRoom.Type == BossRoom &&
		currentFloor == 1 &&
		gurdyDead {

			exitActive = true
			drop.Draw(win, pixel.IM.Moved(exitPos))

			// spawn itemu bossa (raz)
			if !bossItemSpawned {
				bossItemSpawned = true
				bossItemTaken = false
				itemPos = exitPos.Add(pixel.V(0, -80))
			}

			// rysowanie piedestału + arrow
			if !bossItemTaken {
				piedestal.Draw(win, pixel.IM.Moved(pixel.V(960, 340)))
				arrow.Draw(win, pixel.IM.Moved(itemPos))
			}
		}



		// Rysowanie Husha
		hushMux.Lock()
		for _, bossHush := range hushes {
			if bossHush.Active && hushDrawn && !hushDead {
				if bossHush.Health <= 0 && !hushDead {
					hushDead = true
				}
				hush.Draw(win, pixel.IM.Moved(bossHush.Position))
			}
		}
		hushMux.Unlock()

	enemiesMux.Lock()
	for _, e := range enemies {
		if !e.Active {
			continue
		}

		if e.Type == Fly && time.Since(e.LastShot) > 1*time.Second {
			bossTearsMux.Lock()
			bossTears = append(bossTears, &struktury_i_funkcje.BossTear{
				Position: e.Position,
				Velocity: pos.Sub(e.Position).Unit().Scaled(300),
				Active:   true,
			})
			bossTearsMux.Unlock()
			e.LastShot = time.Now()
		}
		

		dir := pos.Sub(e.Position)
		if dir.Len() > 0 {
			dir = dir.Unit()
		}

		switch e.Type {
		case Fly:
			e.Position = e.Position.Add(dir.Scaled(1.5))
		case Slime:
			e.Position = e.Position.Add(dir.Scaled(1.0))
		}
	}
	enemiesMux.Unlock()


		if IsRoomCleared() {
			doorsLocked = false
			roomCleared[currentRoom] = true
		}
	

		// Rysowanie piedestalu i itemu
		if currentRoom.Type == ItemRoom && !itemTaken {
			piedestal.Draw(win, pixel.IM.Moved(pixel.V(930, 540)))
		}
		
		itemMux.Lock()
		if currentRoom.Type == ItemRoom && !itemTaken && losowyItem != nil {
			losowyItem.Draw(win, pixel.IM.Moved(itemPos))
		}				
		itemMux.Unlock()
		
		// Rysowanie postaci na tle
		isaacHpMux.Lock()
		if !isaacDead {
			anim.Draw(win, pixel.IM.Moved(pos))
		}
		isaacHpMux.Unlock()

		// Rysowanie lez
		isaacHpMux.Lock()
		tearsMux.Lock()

		if !isaacDead {
			for _, tear := range tears {
				if tear.Active {
					if time.Since(tear.StartTime) <= tear.Lifetime {
						tear.Position = tear.Position.Add(tear.Velocity.Scaled(dt))
						if hasPiercing {
							switch tear.Direction {
							case "up":
								pierceUp.Draw(win, pixel.IM.Moved(tear.Position))
							case "down":
								pierceDown.Draw(win, pixel.IM.Moved(tear.Position))
							case "left":
								pierceLeft.Draw(win, pixel.IM.Moved(tear.Position))
							case "right":
								pierceRight.Draw(win, pixel.IM.Moved(tear.Position))
							}
						} else {
							lza.Draw(win, pixel.IM.Moved(tear.Position))
						}						

						// Sprawdzanie kolizji ze sciana lub bossem
						if !scianyLzyOff.Contains(tear.Position){
							tear.Active = false
						}
						if !gurdyDead && gurdyDrawn && currentRoom.Type == BossRoom {
							if bossGurdySizeLzy.Contains(tear.Position) || bossGurdySizeLzy.Contains(tear.Position) {
								tear.Active = false
								gurdyHpMux.Lock()
								gurdyGotHit = true
								gurdyHpMux.Unlock()
							}
						}
						if !hushDead && hushDrawn {
							if bossHushSize.Contains(tear.Position) || bossHushSize.Contains(tear.Position) {
								tear.Active = false
								hushHpMux.Lock()
								hushGotHit = true
								hushHpMux.Unlock()

							}
						}

					} else {
						tear.Active = false
					}
				}
			}
		}
		
		tearsMux.Unlock()
		isaacHpMux.Unlock()

		// Rysowanie lez bossa
		isaacHpMux.Lock()
		bossTearsMux.Lock()

		if !isaacDead {
			for _, bossTear := range bossTears {
				if bossTear.Active {
					bossTear.Position = bossTear.Position.Add(bossTear.Velocity.Scaled(dt))
					lza.Draw(win, pixel.IM.Moved(bossTear.Position))

					// Sprawdzanie kolizji ze sciana lub postacia
					if !scianyLzyOff.Contains(bossTear.Position) {
						bossTear.Active = false
					}
					posHitbox := pos
					posHitbox.Y -= 30
					roznicaY := posHitbox.Y - bossTear.Position.Y
					roznicaX := posHitbox.X - bossTear.Position.X
				
					if struktury_i_funkcje.Abs(roznicaY) < 30 && struktury_i_funkcje.Abs(roznicaX) < 30 {
						bossTear.Active = false
						isaacGotHit = true
					}
				}
			}
		}
		bossTearsMux.Unlock()
		isaacHpMux.Unlock()

		// Rysowanie serc
		isaacHpMux.Lock()
		if isaacHp >= 1 {
			serce.Draw(win, pixel.IM.Moved(pixel.V(175, 935)))
		}
		if isaacHp >= 2 {
			serce.Draw(win, pixel.IM.Moved(pixel.V(225, 935)))
		}
		if isaacHp >= 3 {
			serce.Draw(win, pixel.IM.Moved(pixel.V(275, 935)))
		}
		if isaacHp >= 4 {
			serce.Draw(win, pixel.IM.Moved(pixel.V(325, 935)))
		}
		isaacHpMux.Unlock()

		// Rysowanie paska zycia bossa
		if !gurdyDead && gurdyDrawn {
			czarny := imdraw.New(nil)
			czarny.Color = colornames.Black
			czarny.Push(czarneZycie.Min, czarneZycie.Max)
			czarny.Rectangle(0)
			czarny.Draw(win)
			czerwony := imdraw.New(nil)
			czerwony.Color = colornames.Red
			czerwony.Push(czerwoneZycie.Min, czerwoneZycie.Max)
			czerwony.Rectangle(0)
			czerwony.Draw(win)
			pasek.Draw(win, pixel.IM.Moved(pixel.V(930, 110)))
		}
		if !hushDead && hushDrawn {
			czarny2 := imdraw.New(nil)
			czarny2.Color = colornames.Black
			czarny2.Push(czarneZycie2.Min, czarneZycie2.Max)
			czarny2.Rectangle(0)
			czarny2.Draw(win)
			czerwony2 := imdraw.New(nil)
			czerwony2.Color = colornames.Red
			czerwony2.Push(czerwoneZycie2.Min, czerwoneZycie2.Max)
			czerwony2.Rectangle(0)
			czerwony2.Draw(win)
			pasek.Draw(win, pixel.IM.Moved(pixel.V(930, 110)))
		}

		DrawMiniMap(win)

		enemiesMux.Lock()
		for _, e := range enemies {
			if !e.Active {
				continue
			}

			switch e.Type {
			case Fly:
				// 🟣 mucha
				mucha.Draw(win, pixel.IM.Scaled(pixel.ZV, 0.42).Moved(e.Position))

			case Slime:
				// 🟢 follower
				follower.Draw(win, pixel.IM.Moved(e.Position))
			}
		}
enemiesMux.Unlock()



		win.Update()
	}
}

func main() {
	pixelgl.Run(run)
}
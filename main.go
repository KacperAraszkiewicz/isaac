package main

import (
	"mymodule/struktury_i_funkcje"
	"sync"
	"time"
	"image/color"
	"math/rand"

	_ "image/png"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
	"github.com/faiface/pixel/imdraw"
)

func run() {
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

	// Wczytywanie drugiego tla
	pic10 := struktury_i_funkcje.LoadPicture("tlo2.png")
	tlo2 := pixel.NewSprite(pic10, pic10.Bounds())

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

	pos := pixel.V(960, 340)
	poslza := pos
	speed := 7.5
	MoveUp := false
	MoveDown := false
	MoveRight := false
	MoveLeft := false


	// Ograniczenie scian dla postaci
	sciany := pixel.R(275, 210, 1585, 935)
	// Ograniczenie scian dla lez
	scianyLzy := pixel.R(275, 210, 1585, 935)
	// Ograniczenie scian dla znikania lez
	scianyLzyOff := pixel.R(240, 153, 1620, 935)
	// Ograniczenie ruchu postaci dla bossa
	bossGurdySize := pixel.R(636, 557, 1245, 935)
	// Ograniczenie ruchu lez dla bossa
	bossGurdySizeLzy := pixel.R(636, 557, 1245, 935)
	// Kierunek ataku gurdiego
	attackLeft := pixel.R(275, 210, 636, 935)
	attackMiddle := pixel.R(637, 210, 1244, 557)
	attackRight := pixel.R(1245, 210, 1585, 935)
	// Czarny porstokat na pasek zycia bossa
	czarneZycie := pixel.R(705, 87, 1180, 127)
	// Czerwony porstokat na pasek zycia bossa
	czerwoneZycie := pixel.R(705, 87, 1180, 127)
	// Czarny porstokat na pasek zycia drugiego bossa
	czarneZycie2 := pixel.R(705, 87, 1180, 127)
	// Czerwony porstokat na pasek zycia drugiego bossa
	czerwoneZycie2 := pixel.R(705, 87, 1180, 127)
	// Kwadrat na drzwi
	drzwi := pixel.R(911, 200, 951, 220)
	// Ograniczenie dla drugiego bossa
	bossHushSize := pixel.R(716, 407, 1136, 707)
	// Piedestal
	piedestalHitbox := pixel.R(880, 540, 980, 640)
	wziacItem := pixel.R(881, 541, 981, 641)

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
			losuj = true
			itemTaken = false
			losowyItem *pixel.Sprite
			gurdyDead = false
			plansza1 = true
			isaacHp = 3
			hpUp = false
		)

		items := []*pixel.Sprite{
			sok,
			strzykawka,
			but,
			mieso,
		}

		if losuj {
		index := rand.Intn(len(items))
		losowyItem = items[index]
		losuj = false
		}
	
		// Funkcja uzywajaca itemow
	
		go func() {
			for {
				itemMux.Lock()
				if losowyItem == mieso && itemTaken {
					if !hpUp {
						isaacHp = isaacHp + 1
						hpUp = true
					}
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
			if win.Pressed(pixelgl.KeyUp) {
				tearsMux.Lock()

				tear = &struktury_i_funkcje.Tear {
					Position: poslza,
					Velocity: pixel.V(0, 500),
					StartTime: time.Now(),
					Lifetime: lifetime,
					Active: true,
					Damage: damage,
				}
				tears = append(tears, tear)
				tearsMux.Unlock()
			} else if win.Pressed(pixelgl.KeyDown) {
				tearsMux.Lock()

				tear = &struktury_i_funkcje.Tear {
					Position: poslza,
					Velocity: pixel.V(0, -500),
					StartTime: time.Now(),
					Lifetime: lifetime,
					Active: true,
					Damage: damage,
				}
				tears = append(tears, tear)
				tearsMux.Unlock()
			} else if win.Pressed(pixelgl.KeyRight) {
				tearsMux.Lock()

				tear = &struktury_i_funkcje.Tear {
					Position: poslza,
					Velocity: pixel.V(500, 0),
					StartTime: time.Now(),
					Lifetime: lifetime,
					Active: true,
					Damage: damage,
				}
				tears = append(tears, tear)
				tearsMux.Unlock()
			} else if win.Pressed(pixelgl.KeyLeft) {
				tearsMux.Lock()

				tear = &struktury_i_funkcje.Tear {
					Position: poslza,
					Velocity: pixel.V(-500, 0),
					StartTime: time.Now(),
					Lifetime: lifetime,
					Active: true,
					Damage: damage,
				}
				tears = append(tears, tear)
				tearsMux.Unlock()
			}
			time.Sleep(sleep)
		}
	}()

	// Deklaracja Gurdiego
	var (
		gurdys []*struktury_i_funkcje.Gurdy
		gurdyMux sync.Mutex
		gurdyRespawn = time.Now()
		gurdyDrawn = false
		bossGurdy *struktury_i_funkcje.Gurdy
	)

	// Funkcja generujaca Gurdiego
	go func() {
		for {
			if time.Since(gurdyRespawn) >= 1*time.Millisecond && !gurdyDrawn && plansza1 {
				gurdyMux.Lock()
				bossGurdy = &struktury_i_funkcje.Gurdy {
					Position: pixel.V(920, 520),
					Active: true,
					Health: 300,
				}
				gurdys = append(gurdys, bossGurdy)
				gurdyDrawn = true
				gurdyMux.Unlock()
			}
			time.Sleep(1*time.Millisecond)
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
		zmiana = false
		zmianaMux sync.Mutex
	)

	go func() {
		for {
			zmianaMux.Lock()
			if gurdyDead {
				zmiana = true
			}
			if zmiana && plansza1 {
				if drzwi.Contains(pos) {
					zmiana = false
					plansza1 = false
					pos = pixel.V(931, 935)
					poslza = pixel.V(931, 935)
				}
			}
			zmianaMux.Unlock()
		}
	}()

	// Deklaracja lez bossa
	var (
		bossTears []*struktury_i_funkcje.BossTear
		bossTearsMux sync.Mutex
	)

	// Funkcja strzelajaca lzami bossa
	go func() {
		for {
			if !gurdyDead && gurdyDrawn && plansza1{
				if attackLeft.Contains(pos) {
					bossTearsMux.Lock()
					bossTear := &struktury_i_funkcje.BossTear {
						Position: pixel.V(686, 780),
						Velocity: pixel.V(-425, 75),
						Active: true,
					}
					bossTears = append(bossTears, bossTear)
					bossTear = &struktury_i_funkcje.BossTear {
						Position: pixel.V(686, 740),
						Velocity: pixel.V(-425, -75),
						Active: true,
					}
					bossTears = append(bossTears, bossTear)
					bossTear = &struktury_i_funkcje.BossTear {
						Position: pixel.V(686, 820),
						Velocity: pixel.V(-350, 150),
						Active: true,
					}
					bossTears = append(bossTears, bossTear)
					bossTear = &struktury_i_funkcje.BossTear {
						Position: pixel.V(686, 700),
						Velocity: pixel.V(-350, -150),
						Active: true,
					}
					bossTears = append(bossTears, bossTear)
					bossTear = &struktury_i_funkcje.BossTear {
						Position: pixel.V(686, 860),
						Velocity: pixel.V(-275, 225),
						Active: true,
					}
					bossTears = append(bossTears, bossTear)
					bossTear = &struktury_i_funkcje.BossTear {
						Position: pixel.V(686, 660),
						Velocity: pixel.V(-275, -225),
						Active: true,
					}
					bossTears = append(bossTears, bossTear)
					bossTearsMux.Unlock()
				}
				time.Sleep(700*time.Millisecond)
				if attackMiddle.Contains(pos) {
					bossTearsMux.Lock()
					bossTear := &struktury_i_funkcje.BossTear {
						Position: pixel.V(931, 607),
						Velocity: pixel.V(0, -500),
						Active: true,
					}
					bossTears = append(bossTears, bossTear)
					bossTear = &struktury_i_funkcje.BossTear {
						Position: pixel.V(866, 607),
						Velocity: pixel.V(-75, -425),
						Active: true,
					}
					bossTears = append(bossTears, bossTear)
					bossTear = &struktury_i_funkcje.BossTear {
						Position: pixel.V(996, 607),
						Velocity: pixel.V(75, -425),
						Active: true,
					}
					bossTears = append(bossTears, bossTear)
					bossTear = &struktury_i_funkcje.BossTear {
						Position: pixel.V(801, 607),
						Velocity: pixel.V(-150, -350),
						Active: true,
					}
					bossTears = append(bossTears, bossTear)
					bossTear = &struktury_i_funkcje.BossTear {
						Position: pixel.V(1061, 607),
						Velocity: pixel.V(150, -350),
						Active: true,
					}
					bossTears = append(bossTears, bossTear)
					bossTear = &struktury_i_funkcje.BossTear {
						Position: pixel.V(736, 607),
						Velocity: pixel.V(-225, -275),
						Active: true,
					}
					bossTears = append(bossTears, bossTear)
					bossTear = &struktury_i_funkcje.BossTear {
						Position: pixel.V(1126, 607),
						Velocity: pixel.V(225, -275),
						Active: true,
					}
					bossTears = append(bossTears, bossTear)
					bossTearsMux.Unlock()
				}
				time.Sleep(700*time.Millisecond)
				if attackRight.Contains(pos) {
					bossTearsMux.Lock()
					bossTear := &struktury_i_funkcje.BossTear {
						Position: pixel.V(1195, 780),
						Velocity: pixel.V(425, 75),
						Active: true,
					}
					bossTears = append(bossTears, bossTear)
					bossTear = &struktury_i_funkcje.BossTear {
						Position: pixel.V(1195, 740),
						Velocity: pixel.V(425, -75),
						Active: true,
					}
					bossTears = append(bossTears, bossTear)
					bossTear = &struktury_i_funkcje.BossTear {
						Position: pixel.V(1195, 820),
						Velocity: pixel.V(350, 150),
						Active: true,
					}
					bossTears = append(bossTears, bossTear)
					bossTear = &struktury_i_funkcje.BossTear {
						Position: pixel.V(1195, 700),
						Velocity: pixel.V(350, -150),
						Active: true,
					}
					bossTears = append(bossTears, bossTear)
					bossTear = &struktury_i_funkcje.BossTear {
						Position: pixel.V(1195, 860),
						Velocity: pixel.V(275, 225),
						Active: true,
					}
					bossTears = append(bossTears, bossTear)
					bossTear = &struktury_i_funkcje.BossTear {
						Position: pixel.V(1195, 660),
						Velocity: pixel.V(275, -225),
						Active: true,
					}
					bossTears = append(bossTears, bossTear)
					bossTearsMux.Unlock()
				}
				time.Sleep(700*time.Millisecond)
			}
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

	// funkcja generujaca drugiego bossa
	var (
		hushes []*struktury_i_funkcje.Hush
		hushMux sync.Mutex
		hushDrawn = false
		hushDead = false
		bossHush *struktury_i_funkcje.Hush
	)

	go func() {
		for {
			if !plansza1 && !hushDrawn {
				hushMux.Lock()
				bossHush = &struktury_i_funkcje.Hush {
					Position: pixel.V(920, 520),
					Active: true,
					Health: 400,
				}
				hushes = append(hushes, bossHush)
				hushDrawn = true
				hushMux.Unlock()
			}
			time.Sleep(1*time.Millisecond)
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
			time.Sleep(3*time.Second)
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

		if gurdyDead && plansza1 && piedestalHitbox.Contains(NewPos) {
			pos = oldPos
		}
		if gurdyDead && plansza1 && wziacItem.Contains(NewPos) {
			if !itemTaken {
				itemTaken = true
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

		if !gurdyDead {
			if bossGurdySizeLzy.Contains(NewPosLza) {
				poslza = oldPosLzy
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
			tlo2.Draw(win, pixel.IM.Moved(pixel.V(930, 525)))
		}
		zmianaMux.Unlock()

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

		// Rysowanie piedestalu i itemu
		if gurdyDead && plansza1 {
			piedestal.Draw(win, pixel.IM.Moved(pixel.V(930, 540)))
		}
		
		itemMux.Lock()
		if gurdyDead && plansza1 && !itemTaken {
			losowyItem.Draw(win, pixel.IM.Moved(pixel.V(935, 620)))
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
						lza.Draw(win, pixel.IM.Moved(tear.Position))

						// Sprawdzanie kolizji ze sciana lub bossem
						if !scianyLzyOff.Contains(tear.Position){
							tear.Active = false
						}
						if !gurdyDead && gurdyDrawn {
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


		win.Update()
	}
}

func main() {
	pixelgl.Run(run)
}
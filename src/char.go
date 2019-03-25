package main

import (
	"fmt"
	"io"
	"math"
	"os"
	"strings"
)

const MaxPalNo = 12
const MaxQuotes = 100

type SystemCharFlag uint32

const (
	SCF_ko SystemCharFlag = 1 << iota
	SCF_ctrl
	SCF_standby
	SCF_guard
	SCF_airjump
	SCF_over
	SCF_ko_round_middle
)

type CharSpecialFlag uint32

const (
	CSF_nostandguard CharSpecialFlag = 1 << iota
	CSF_nocrouchguard
	CSF_noairguard
	CSF_noshadow
	CSF_invisible
	CSF_unguardable
	CSF_nojugglecheck
	CSF_noautoturn
	CSF_nowalk
	CSF_screenbound
	CSF_movecamera_x
	CSF_movecamera_y
	CSF_posfreeze
	CSF_playerpush
	CSF_angledraw
	CSF_destroy
	CSF_frontedge
	CSF_backedge
	CSF_frontwidth
	CSF_backwidth
	CSF_trans
	CSF_gethit
	CSF_assertspecial CharSpecialFlag = CSF_nostandguard | CSF_nocrouchguard |
		CSF_noairguard | CSF_noshadow | CSF_invisible | CSF_unguardable |
		CSF_nojugglecheck | CSF_noautoturn | CSF_nowalk
)

type GlobalSpecialFlag uint32

const (
	GSF_intro GlobalSpecialFlag = 1 << iota
	GSF_roundnotover
	GSF_nomusic
	GSF_nobardisplay
	GSF_nobg
	GSF_nofg
	GSF_globalnoshadow
	GSF_timerfreeze
	GSF_nokosnd
	GSF_nokoslow
	GSF_noko
)

type PosType int32

const (
	PT_P1 PosType = iota
	PT_P2
	PT_F
	PT_B
	PT_L
	PT_R
	PT_N
)

type Space int32

const (
	Space_none Space = iota
	Space_stage
	Space_screen
)

type ClsnRect [][4]float32

func (cr *ClsnRect) Add(clsn []float32, x, y, xs, ys float32) {
	x = (x - sys.cam.Pos[0]) * sys.cam.Scale
	y = (y-sys.cam.Pos[1])*sys.cam.Scale + sys.cam.GroundLevel()
	xs *= sys.cam.Scale
	ys *= sys.cam.Scale
	for i := 0; i+3 < len(clsn); i += 4 {
		rect := [...]float32{x + xs*clsn[i] + float32(sys.gameWidth)/2,
			y + ys*clsn[i+1] + float32(sys.gameHeight-240),
			xs * (clsn[i+2] - clsn[i]), ys * (clsn[i+3] - clsn[i+1])}
		if xs < 0 {
			rect[0] *= -1
		}
		if ys < 0 {
			rect[1] *= -1
		}
		*cr = append(*cr, rect)
	}
}
func (cr ClsnRect) draw(trans int32) {
	for _, c := range cr {
		RenderMugen(*sys.clsnSpr.Tex, sys.clsnSpr.Pal, -1, sys.clsnSpr.Size,
			-c[0]*sys.widthScale, -c[1]*sys.heightScale, &notiling,
			c[2]*sys.widthScale, c[2]*sys.widthScale, c[3]*sys.heightScale, 1, 0, 0,
			trans, &sys.scrrect, 0, 0)
	}
}

type CharData struct {
	life    int32
	power   int32
	attack  int32
	defence int32
	fall    struct {
		defence_mul float32
	}
	liedown struct {
		time int32
	}
	airjuggle int32
	sparkno   int32
	guard     struct {
		sparkno int32
	}
	ko struct {
		echo int32
	}
	volume            int32
	intpersistindex   int32
	floatpersistindex int32
}

func (cd *CharData) init() {
	*cd = CharData{}
	cd.life = 1000
	cd.power = 3000
	cd.attack = 100
	cd.defence = 100
	cd.fall.defence_mul = 1.5
	cd.liedown.time = 60
	cd.airjuggle = 15
	cd.sparkno = 2
	cd.guard.sparkno = 40
	cd.ko.echo = 0
	cd.volume = 256
	cd.intpersistindex = 0
	cd.floatpersistindex = 0
}

type CharSize struct {
	xscale float32
	yscale float32
	ground struct {
		back  float32
		front float32
	}
	air struct {
		back  float32
		front float32
	}
	height float32
	attack struct {
		dist float32
		z    struct {
			width [2]float32
		}
	}
	proj struct {
		attack struct {
			dist float32
		}
		doscale int32
	}
	head struct {
		pos [2]float32
	}
	mid struct {
		pos [2]float32
	}
	shadowoffset float32
	draw         struct {
		offset [2]float32
	}
	z struct {
		width float32
	}
}

func (cs *CharSize) init() {
	*cs = CharSize{}
	cs.xscale = 1
	cs.yscale = 1
	cs.ground.back = 15
	cs.ground.front = 16
	cs.air.back = 12
	cs.air.front = 12
	cs.height = 60
	cs.attack.dist = 160
	cs.proj.attack.dist = 90
	cs.proj.doscale = 0
	cs.head.pos = [...]float32{-5, -90}
	cs.mid.pos = [...]float32{-5, -60}
	cs.shadowoffset = 0
	cs.draw.offset = [...]float32{0, 0}
	cs.z.width = 3
	cs.attack.z.width = [...]float32{4, 4}
}

type CharVelocity struct {
	walk struct {
		fwd  float32
		back float32
		up   struct {
			x float32
		}
		down struct {
			x float32
		}
	}
	run struct {
		fwd  [2]float32
		back [2]float32
		up   struct {
			x float32
			y float32
		}
		down struct {
			x float32
			y float32
		}
	}
	jump struct {
		neu  [2]float32
		back float32
		fwd  float32
		up   struct {
			x float32
		}
		down struct {
			x float32
		}
	}
	runjump struct {
		back [2]float32
		fwd  [2]float32
		up   struct {
			x float32
		}
		down struct {
			x float32
		}
	}
	airjump struct {
		neu  [2]float32
		back float32
		fwd  float32
		up   struct {
			x float32
		}
		down struct {
			x float32
		}
	}
	air struct {
		gethit struct {
			groundrecover [2]float32
			airrecover    struct {
				mul  [2]float32
				add  [2]float32
				back float32
				fwd  float32
				up   float32
				down float32
			}
		}
	}
}

func (cv *CharVelocity) init() {
	*cv = CharVelocity{}
	cv.air.gethit.groundrecover = [...]float32{-0.15, -3.5}
	cv.air.gethit.airrecover.mul = [...]float32{0.5, 0.2}
	cv.air.gethit.airrecover.add = [...]float32{0.0, -4.5}
	cv.air.gethit.airrecover.back = -1.0
	cv.air.gethit.airrecover.fwd = 0.0
	cv.air.gethit.airrecover.up = -2.0
	cv.air.gethit.airrecover.down = 1.5
}

type CharMovement struct {
	airjump struct {
		num    int32
		height int32
	}
	yaccel float32
	stand  struct {
		friction           float32
		friction_threshold float32
	}
	crouch struct {
		friction           float32
		friction_threshold float32
	}
	air struct {
		gethit struct {
			groundlevel   float32
			groundrecover struct {
				ground struct {
					threshold float32
				}
				groundlevel float32
			}
			airrecover struct {
				threshold float32
				yaccel    float32
			}
			trip struct {
				groundlevel float32
			}
		}
	}
	down struct {
		bounce struct {
			offset      [2]float32
			yaccel      float32
			groundlevel float32
		}
		friction_threshold float32
	}
}

func (cm *CharMovement) init() {
	*cm = CharMovement{}
	cm.airjump.num = 0
	cm.airjump.height = 35
	cm.yaccel = 0.44
	cm.stand.friction = 0.85
	cm.stand.friction_threshold = 2.0
	cm.crouch.friction = 0.82
	cm.crouch.friction_threshold = 0.0
	cm.air.gethit.groundlevel = 10.0
	cm.air.gethit.groundrecover.ground.threshold = -20.0
	cm.air.gethit.groundrecover.groundlevel = 10.0
	cm.air.gethit.airrecover.threshold = -1.0
	cm.air.gethit.airrecover.yaccel = 0.35
	cm.air.gethit.trip.groundlevel = 15.0
	cm.down.bounce.offset = [...]float32{0.0, 20.0}
	cm.down.bounce.yaccel = 0.4
	cm.down.bounce.groundlevel = 12.0
	cm.down.friction_threshold = 0.05
}

type Reaction int32

const (
	RA_Light   Reaction = 0
	RA_Medium  Reaction = 1
	RA_Hard    Reaction = 2
	RA_Back    Reaction = 3
	RA_Up      Reaction = 4
	RA_Diagup  Reaction = 5
	RA_Unknown Reaction = -1
)

type HitType int32

const (
	HT_None    HitType = 0
	HT_High    HitType = 1
	HT_Low     HitType = 2
	HT_Trip    HitType = 3
	HT_Unknown HitType = -1
)

type AiuchiType int32

const (
	AT_Hit AiuchiType = iota
	AT_Miss
	AT_Dodge
)

type Fall struct {
	animtype       Reaction
	xvelocity      float32
	yvelocity      float32
	recover        bool
	recovertime    int32
	damage         int32
	kill           bool
	envshake_time  int32
	envshake_freq  float32
	envshake_ampl  int32
	envshake_phase float32
}

func (f *Fall) clear() {
	*f = Fall{animtype: RA_Unknown, xvelocity: float32(math.NaN()),
		yvelocity: -4.5}
}
func (f *Fall) setDefault() {
	*f = Fall{animtype: RA_Unknown, xvelocity: float32(math.NaN()),
		yvelocity: -4.5, recover: true, recovertime: 4, kill: true,
		envshake_freq: 60, envshake_ampl: -4, envshake_phase: float32(math.NaN())}
}
func (f *Fall) xvel() float32 {
	if math.IsNaN(float64(f.xvelocity)) {
		return -32760
	}
	return f.xvelocity
}

type HitDef struct {
	attr                       int32
	reversal_attr              int32
	hitflag                    int32
	guardflag                  int32
	affectteam                 int32
	animtype                   Reaction
	air_animtype               Reaction
	priority                   int32
	bothhittype                AiuchiType
	hitdamage                  int32
	guarddamage                int32
	pausetime                  int32
	shaketime                  int32
	guard_pausetime            int32
	guard_shaketime            int32
	sparkno                    int32
	guard_sparkno              int32
	sparkxy                    [2]float32
	hitsound                   [2]int32
	guardsound                 [2]int32
	ground_type                HitType
	air_type                   HitType
	ground_slidetime           int32
	guard_slidetime            int32
	ground_hittime             int32
	guard_hittime              int32
	air_hittime                int32
	guard_ctrltime             int32
	airguard_ctrltime          int32
	guard_dist                 int32
	yaccel                     float32
	ground_velocity            [2]float32
	guard_velocity             float32
	air_velocity               [2]float32
	airguard_velocity          [2]float32
	ground_cornerpush_veloff   float32
	air_cornerpush_veloff      float32
	down_cornerpush_veloff     float32
	guard_cornerpush_veloff    float32
	airguard_cornerpush_veloff float32
	air_juggle                 int32
	p1sprpriority              int32
	p2sprpriority              int32
	p1getp2facing              int32
	p1facing                   int32
	p2facing                   int32
	p1stateno                  int32
	p2stateno                  int32
	p2getp1state               bool
	forcestand                 int32
	ground_fall                bool
	air_fall                   bool
	down_velocity              [2]float32
	down_hittime               int32
	down_bounce                bool
	id                         int32
	chainid                    int32
	nochainid                  [2]int32
	hitonce                    int32
	numhits                    int32
	hitgetpower                int32
	guardgetpower              int32
	hitgivepower               int32
	guardgivepower             int32
	palfx                      PalFXDef
	envshake_time              int32
	envshake_freq              float32
	envshake_ampl              int32
	envshake_phase             float32
	mindist                    [2]float32
	maxdist                    [2]float32
	snap                       [2]float32
	snapt                      int32
	fall                       Fall
	playerNo                   int
	kill                       bool
	guard_kill                 bool
	forcenofall                bool
	lhit                       bool
}

func (hd *HitDef) clear() {
	*hd = HitDef{hitflag: int32(ST_S | ST_C | ST_A | ST_F), affectteam: 1,
		animtype: RA_Light, air_animtype: RA_Unknown, priority: 4,
		bothhittype: AT_Hit, sparkno: IErr, guard_sparkno: IErr,
		hitsound: [...]int32{IErr, -1}, guardsound: [...]int32{IErr, -1},
		ground_type: HT_High, air_type: HT_Unknown, air_hittime: 20,
		yaccel: float32(math.NaN()), guard_velocity: float32(math.NaN()),
		airguard_velocity: [...]float32{float32(math.NaN()),
			float32(math.NaN())},
		ground_cornerpush_veloff:   float32(math.NaN()),
		air_cornerpush_veloff:      float32(math.NaN()),
		down_cornerpush_veloff:     float32(math.NaN()),
		guard_cornerpush_veloff:    float32(math.NaN()),
		airguard_cornerpush_veloff: float32(math.NaN()), p1sprpriority: 1,
		p1stateno: -1, p2stateno: -1, forcestand: IErr,
		down_velocity: [...]float32{float32(math.NaN()), float32(math.NaN())},
		chainid:       -1, nochainid: [...]int32{-1, -1}, numhits: 1,
		hitgetpower: IErr, guardgetpower: IErr, hitgivepower: IErr,
		guardgivepower: IErr, envshake_freq: 60, envshake_ampl: -4,
		envshake_phase: float32(math.NaN()),
		mindist:        [...]float32{float32(math.NaN()), float32(math.NaN())},
		maxdist:        [...]float32{float32(math.NaN()), float32(math.NaN())},
		snap:           [...]float32{float32(math.NaN()), float32(math.NaN())},
		kill:           true, guard_kill: true, playerNo: -1}
	hd.palfx.mul, hd.palfx.color = [...]int32{255, 255, 255}, 1
	hd.fall.setDefault()
}
func (hd *HitDef) invalidate(stateType StateType) {
	hd.attr = hd.attr&^int32(ST_MASK) | int32(stateType) | -1<<31
	hd.reversal_attr |= -1 << 31
	hd.lhit = false
}
func (hd *HitDef) testAttr(attr int32) bool {
	attr &= hd.attr
	return attr&int32(ST_MASK) != 0 && attr&^int32(ST_MASK)&^(-1<<31) != 0
}

type GetHitVar struct {
	hitBy          [][2]int32
	hit1           [2]int32
	hit2           [2]int32
	attr           int32
	_type          HitType
	airanimtype    Reaction
	groundanimtype Reaction
	airtype        HitType
	groundtype     HitType
	damage         int32
	hitcount       int32
	fallcount      int32
	hitshaketime   int32
	hittime        int32
	slidetime      int32
	ctrltime       int32
	xvel           float32
	yvel           float32
	yaccel         float32
	hitid          int32
	xoff           float32
	yoff           float32
	fall           Fall
	playerNo       int
	fallf          bool
	guarded        bool
	p2getp1state   bool
	forcestand     bool
}

func (ghv *GetHitVar) clear() {
	*ghv = GetHitVar{_type: -1, hittime: -1, yaccel: float32(math.NaN()),
		xoff: ghv.xoff, yoff: ghv.yoff, hitid: -1, playerNo: -1}
	ghv.fall.clear()
}
func (ghv *GetHitVar) clearOff() {
	ghv.xoff, ghv.yoff = 0, 0
}
func (ghv GetHitVar) getYaccel(c *Char) float32 {
	if math.IsNaN(float64(ghv.yaccel)) {
		return 0.35 / c.localscl
	}
	return ghv.yaccel
}
func (ghv GetHitVar) chainId() int32 {
	if ghv.hitid > 0 {
		return ghv.hitid
	}
	return 0
}
func (ghv GetHitVar) idMatch(id int32) bool {
	for _, v := range ghv.hitBy {
		if v[0] == id {
			return true
		}
	}
	return false
}
func (ghv GetHitVar) getJuggle(id, defaultJuggle int32) int32 {
	for _, v := range ghv.hitBy {
		if v[0] == id {
			return v[1]
		}
	}
	return defaultJuggle
}
func (ghv *GetHitVar) dropId(id int32) {
	for i, v := range ghv.hitBy {
		if v[0] == id {
			ghv.hitBy = append(ghv.hitBy[:i], ghv.hitBy[i+1:]...)
			break
		}
	}
}
func (ghv *GetHitVar) addId(id, juggle int32) {
	juggle = ghv.getJuggle(id, juggle)
	ghv.dropId(id)
	ghv.hitBy = append(ghv.hitBy, [...]int32{id, juggle})
}

type HitBy struct {
	flag, time int32
}
type HitOverride struct {
	attr     int32
	stateno  int32
	time     int32
	forceair bool
	playerNo int
}

func (ho *HitOverride) clear() {
	*ho = HitOverride{stateno: -1, playerNo: -1}
}

type aimgImage struct {
	anim           Animation
	pos, scl, ascl [2]float32
	angle          float32
	oldVer         bool
}

type AfterImage struct {
	time       int32
	length     int32
	postbright [3]int32
	add        [3]int32
	mul        [3]float32
	timegap    int32
	framegap   int32
	alpha      [2]int32
	palfx      []PalFX
	imgs       [64]aimgImage
	imgidx     int32
	restgap    int32
	reccount   int32
}

func newAfterImage() *AfterImage {
	ai := &AfterImage{palfx: make([]PalFX, sys.afterImageMax)}
	for i := range ai.palfx {
		ai.palfx[i].enable, ai.palfx[i].negType = true, true
	}
	ai.clear()
	ai.timegap = 0
	return ai
}
func (ai *AfterImage) clear() {
	ai.time = 0
	ai.length = 20
	if len(ai.palfx) > 0 {
		ai.palfx[0].eColor = 1
		ai.palfx[0].eInvertall = false
		ai.palfx[0].eAdd = [...]int32{30, 30, 30}
		ai.palfx[0].eMul = [...]int32{120, 120, 220}
	}
	ai.postbright = [3]int32{}
	ai.add = [...]int32{10, 10, 25}
	ai.mul = [...]float32{0.65, 0.65, 0.75}
	ai.timegap = 1
	ai.framegap = 6
	ai.alpha = [...]int32{-1, 0}
	ai.imgidx = 0
	ai.restgap = 0
	ai.reccount = 0
}
func (ai *AfterImage) setPalColor(color int32) {
	if len(ai.palfx) > 0 {
		ai.palfx[0].eColor = float32(Max(0, Min(256, color))) / 256
	}
}
func (ai *AfterImage) setPalInvertall(invertall bool) {
	if len(ai.palfx) > 0 {
		ai.palfx[0].eInvertall = invertall
	}
}
func (ai *AfterImage) setPalBrightR(addr int32) {
	if len(ai.palfx) > 0 {
		ai.palfx[0].eAdd[0] = addr
	}
}
func (ai *AfterImage) setPalBrightG(addg int32) {
	if len(ai.palfx) > 0 {
		ai.palfx[0].eAdd[1] = addg
	}
}
func (ai *AfterImage) setPalBrightB(addb int32) {
	if len(ai.palfx) > 0 {
		ai.palfx[0].eAdd[2] = addb
	}
}
func (ai *AfterImage) setPalContrastR(mulr int32) {
	if len(ai.palfx) > 0 {
		ai.palfx[0].eMul[0] = mulr
	}
}
func (ai *AfterImage) setPalContrastG(mulg int32) {
	if len(ai.palfx) > 0 {
		ai.palfx[0].eMul[1] = mulg
	}
}
func (ai *AfterImage) setPalContrastB(mulb int32) {
	if len(ai.palfx) > 0 {
		ai.palfx[0].eMul[2] = mulb
	}
}
func (ai *AfterImage) setupPalFX() {
	pb := ai.postbright
	for i := 1; i < len(ai.palfx); i++ {
		ai.palfx[i].eColor = ai.palfx[i-1].eColor
		ai.palfx[i].eInvertall = ai.palfx[i-1].eInvertall
		for j := range pb {
			ai.palfx[i].eAdd[j] = ai.palfx[i-1].eAdd[j] + ai.add[j] + pb[j]
			ai.palfx[i].eMul[j] = int32(float32(ai.palfx[i-1].eMul[j]) * ai.mul[j])
		}
		pb = [3]int32{}
	}
}
func (ai *AfterImage) recAfterImg(sd *SprData) {
	if ai.time == 0 {
		ai.reccount, ai.timegap = 0, 0
		return
	}
	if ai.time > 0 {
		ai.time--
	}
	if ai.restgap <= 0 {
		img := &ai.imgs[ai.imgidx]
		img.anim = *sd.anim
		if sd.anim.spr != nil {
			img.anim.spr = newSprite()
			*img.anim.spr = *sd.anim.spr
			sd.anim.sff.palList.SwapPalMap(&sd.fx.remap)
			img.anim.spr.Pal = sd.anim.spr.GetPal(&sd.anim.sff.palList)
			sd.anim.sff.palList.SwapPalMap(&sd.fx.remap)
		}
		img.pos = sd.pos
		img.scl = sd.scl
		img.angle = sd.angle
		img.ascl = sd.ascl
		img.oldVer = sd.oldVer
		ai.imgidx = (ai.imgidx + 1) & 63
		if int(ai.reccount) < len(ai.imgs) {
			ai.reccount++
		}
		ai.restgap = ai.timegap
	}
	ai.restgap--
}
func (ai *AfterImage) recAndCue(sd *SprData, rec bool) {
	if ai.time == 0 || ai.timegap < 1 || ai.timegap > 32767 ||
		ai.framegap < 1 || ai.framegap > 32767 {
		ai.time = 0
		ai.reccount, ai.timegap = 0, 0
		return
	}
	end := Min(sys.afterImageMax,
		(Min(ai.reccount, ai.length)/ai.framegap)*ai.framegap)
	for i := ai.framegap; i <= end; i += ai.framegap {
		img := &ai.imgs[(ai.imgidx-i)&63]
		sys.sprites.add(&SprData{&img.anim, &ai.palfx[i/ai.framegap-1], img.pos,
			img.scl, ai.alpha, sd.priority - 2, img.angle, img.ascl,
			false, sd.bright, sd.oldVer, sd.facing}, 0, 0, 0, 0)
	}
	if rec {
		ai.recAfterImg(sd)
	}
}

type Explod struct {
	id             int32
	bindtime       int32
	scale          [2]float32
	time           int32
	removeongethit bool
	removetime     int32
	velocity       [2]float32
	accel          [2]float32
	sprpriority    int32
	postype        PosType
	space          Space
	offset         [2]float32
	relativef      int32
	pos            [2]float32
	facing         float32
	vfacing        float32
	shadow         [3]int32
	supermovetime  int32
	pausemovetime  int32
	anim           *Animation
	ontop          bool
	alpha          [2]int32
	ownpal         bool
	playerId       int32
	bindId         int32
	ignorehitpause bool
	angle          float32
	oldPos         [2]float32
	newPos         [2]float32
	palfx          *PalFX
	localscl       float32
}

func (e *Explod) clear() {
	*e = Explod{id: IErr, scale: [...]float32{1, 1}, removetime: -2,
		postype: PT_P1, relativef: 1, facing: 1, vfacing: 1, localscl: 1, space: Space_none,
		alpha: [...]int32{-1, 0}, playerId: -1, bindId: -1, ignorehitpause: true}
}
func (e *Explod) setX(x float32) {
	e.pos[0], e.oldPos[0], e.newPos[0] = x, x, x
}
func (e *Explod) setY(y float32) {
	e.pos[1], e.oldPos[1], e.newPos[1] = y, y, y
}
func (e *Explod) setPos(c *Char) {
	pPos := func(c *Char) {
		e.bindId, e.facing = c.id, c.facing*float32(e.relativef)
		e.offset[0] *= c.facing
		e.setX(c.pos[0]*c.localscl/e.localscl + c.offsetX()*c.localscl/e.localscl + e.offset[0])
		e.setY(c.pos[1]*c.localscl/e.localscl + c.offsetY()*c.localscl/e.localscl + e.offset[1])
		if e.bindtime == 0 {
			e.bindtime = 1
		}
	}
	lPos := func() {
		e.setX(sys.cam.ScreenPos[0]/e.localscl + e.offset[0]/sys.cam.Scale)
		e.setY(sys.cam.ScreenPos[1]/e.localscl + e.offset[1]/sys.cam.Scale)
		if e.bindtime == 0 {
			e.bindtime = 1
		}
	}
	rPos := func() {
		e.setX(sys.cam.ScreenPos[0]/e.localscl +
			(float32(sys.gameWidth)/e.localscl + e.offset[0]/sys.cam.Scale))
		e.setY(sys.cam.ScreenPos[1]/e.localscl + e.offset[1]/sys.cam.Scale)
		if e.bindtime == 0 {
			e.bindtime = 1
		}
	}
	if e.space >= Space_stage {
		e.postype = PT_N
	}
	if e.space <= Space_none {
		switch e.postype {
		case PT_P1:
			pPos(c)
		case PT_P2:
			if p2 := sys.charList.enemyNear(c, 0, true); p2 != nil {
				pPos(p2)
			}
		case PT_F, PT_B:
			e.facing = c.facing * float32(e.relativef)
			// front と back はバインドの都合で left か right になおす
			if c.facing > 0 && e.postype == PT_F || c.facing < 0 && e.postype == PT_B {
				if e.postype == PT_B {
					e.offset[0] *= -1
				}
				e.postype = PT_R
				rPos()
			} else {
				// explod の postype = front はキャラの向きで pos が反転しない
				//if e.postype == PT_F && c.gi().ver[0] != 1 {
				// 旧バージョンだと front は キャラの向きが facing に反映されない
				// 1.1でも反映されてない模様
				e.facing = float32(e.relativef)
				//}
				e.postype = PT_L
				lPos()
			}
		case PT_L:
			e.facing = float32(e.relativef)
			lPos()
		case PT_R:
			e.facing = float32(e.relativef)
			rPos()
		case PT_N:
			e.facing = float32(e.relativef)
			e.setX(e.offset[0])
			e.setY(e.offset[1])
		}
	} else {
		switch e.space {
		case Space_screen:
			e.facing = float32(e.relativef)
			lPos()
		case Space_stage:
			e.facing = float32(e.relativef)
			e.setX(e.offset[0])
			e.setY(e.offset[1])
		}
	}
}
func (e *Explod) matchId(eid, pid int32) bool {
	return e.id >= 0 && e.playerId == pid && (eid < 0 || e.id == eid)
}
func (e *Explod) update(oldVer bool, playerNo int) {
	if e.anim == nil {
		e.id = IErr
	}
	if e.id == IErr {
		e.anim = nil
		return
	}
	var c *Char
	if !e.ignorehitpause || e.removeongethit {
		c = sys.playerID(e.playerId)
	}
	p := false
	if sys.super > 0 {
		p = e.supermovetime >= 0 && e.time >= e.supermovetime
	} else if sys.pause > 0 {
		p = e.pausemovetime >= 0 && e.time >= e.pausemovetime
	}
	act := !p
	if act && !e.ignorehitpause {
		act = c == nil || c.acttmp%2 >= 0
	}
	if sys.tickFrame() {
		if c != nil && e.removeongethit && c.sf(CSF_gethit) ||
			e.removetime >= 0 && e.time >= e.removetime ||
			act && e.removetime < -1 && e.anim.loopend {
			e.id, e.anim = IErr, nil
			return
		}
	}
	screen := false
	if e.space == Space_screen || e.postype >= PT_L && e.postype != PT_N {
		screen = true
	}
	if e.bindtime != 0 {
		if e.space == Space_screen {
			e.pos[0] = e.offset[0]
			e.pos[1] = e.offset[1]
			e.pos[0] -= float32(sys.gameWidth) / e.localscl / 2
		} else if e.postype == PT_N && e.bindId < 0 {
			e.pos[0] = e.offset[0]
			e.pos[1] = e.offset[1]
			e.bindtime = 0
		} else if e.postype >= PT_L && e.postype != PT_N {
			e.pos[0] = e.offset[0]
			e.pos[1] = e.offset[1]
			if e.postype == PT_L {
				e.pos[0] -= float32(sys.gameWidth) / e.localscl / 2
			} else {
				e.pos[0] += float32(sys.gameWidth) / e.localscl / 2
			}
		} else {
			if c := sys.playerID(e.bindId); c != nil {
				e.pos[0] = c.drawPos[0]*c.localscl/e.localscl + c.offsetX()*c.localscl/e.localscl + e.offset[0]
				e.pos[1] = c.drawPos[1]*c.localscl/e.localscl + c.offsetY()*c.localscl/e.localscl + e.offset[1]
			} else {
				e.bindtime = 0
			}
		}
	} else {
		for i := range e.pos {
			e.pos[i] = e.newPos[i] -
				(e.newPos[i]-e.oldPos[i])*(1-sys.tickInterpola())
		}
	}
	if sys.tickFrame() && act {
		e.anim.UpdateSprite()
	}
	sprs := &sys.sprites
	if e.ontop {
		sprs = &sys.topSprites
	}
	var pfx *PalFX
	if e.anim.sff != sys.lifebar.fsff {
		pfx = e.palfx
	} else if !e.ownpal {
		pfx = &PalFX{}
		*pfx = *e.palfx
		pfx.remap = nil
	}
	alp := e.alpha
	if alp[0] < 0 {
		alp[0] = -1
	}
	agl := e.angle
	if (e.facing < 0) != (e.vfacing < 0) {
		agl *= -1
	}
	sdwalp := 255 - alp[1]
	if sdwalp < 0 {
		sdwalp = 256
	}
	var epos = [2]float32{e.pos[0] * e.localscl, e.pos[1] * e.localscl}
	sprs.add(&SprData{e.anim, pfx, epos, [...]float32{e.facing * e.scale[0] * e.localscl,
		e.vfacing * e.scale[1] * e.localscl}, alp, e.sprpriority, agl, [...]float32{1, 1},
		screen, playerNo == sys.superplayer, oldVer, e.facing},
		e.shadow[0]<<16|e.shadow[1]&0xff<<8|e.shadow[0]&0xff, sdwalp, 0, 0)
	if sys.tickNextFrame() {
		if e.bindtime > 0 {
			e.bindtime--
		}
		//if screen && e.bindtime == 0 {
		//	if e.space <= Space_none {
		//		switch e.postype {
		//		case PT_L:
		//			for i := range e.pos {
		//				e.pos[i] = sys.cam.ScreenPos[i] + e.offset[i]/sys.cam.Scale
		//			}
		//		case PT_R:
		//			e.pos[0] = sys.cam.ScreenPos[0] +
		//				(float32(sys.gameWidth)+e.offset[0])/sys.cam.Scale
		//			e.pos[1] = sys.cam.ScreenPos[1] + e.offset[1]/sys.cam.Scale
		//		}
		//	} else if e.space == Space_screen {
		//		for i := range e.pos {
		//			e.pos[i] = sys.cam.ScreenPos[i] + e.offset[i]/sys.cam.Scale
		//		}
		//	}
		//}
		if act {
			if e.bindtime == 0 {
				e.oldPos = e.pos
				e.newPos[0] = e.pos[0] + e.velocity[0]*e.facing*float32(e.relativef)
				e.newPos[1] = e.pos[1] + e.velocity[1]
				for i := range e.velocity {
					e.velocity[i] += e.accel[i]
				}
			}
			e.anim.Action()
			e.time++
		} else {
			e.setX(e.pos[0])
			e.setY(e.pos[1])
		}
	}
}

type Projectile struct {
	hitdef        HitDef
	id            int32
	anim          int32
	hitanim       int32
	remanim       int32
	cancelanim    int32
	scale         [2]float32
	clsnScale     [2]float32
	remove        bool
	removetime    int32
	velocity      [2]float32
	remvelocity   [2]float32
	accel         [2]float32
	velmul        [2]float32
	hits          int32
	misstime      int32
	priority      int32
	prioritypoint int32
	sprpriority   int32
	edgebound     int32
	stagebound    int32
	heightbound   [2]int32
	pos           [2]float32
	facing        float32
	shadow        [3]int32
	supermovetime int32
	pausemovetime int32
	ani           *Animation
	timemiss      int32
	hitpause      int32
	oldPos        [2]float32
	newPos        [2]float32
	aimg          AfterImage
	palfx         *PalFX
	localscl      float32
}

func newProjectile() *Projectile {
	p := &Projectile{}
	p.clear()
	return p
}
func (p *Projectile) clear() {
	*p = Projectile{id: IErr, hitanim: -1, remanim: IErr, cancelanim: IErr,
		scale: [...]float32{1, 1}, clsnScale: [...]float32{1, 1}, remove: true, localscl: 1,
		removetime: -1, velmul: [...]float32{1, 1}, hits: 1, priority: 1,
		prioritypoint: 1, sprpriority: 3, edgebound: 40, stagebound: 40,
		heightbound: [...]int32{-240, 1}, facing: 1, aimg: *newAfterImage()}
	p.hitdef.clear()
}
func (p *Projectile) setPos(pos [2]float32) {
	p.pos, p.oldPos, p.newPos = pos, pos, pos
}
func (p *Projectile) paused(playerNo int) bool {
	if !sys.chars[playerNo][0].pause() {
		if sys.super > 0 {
			if p.supermovetime == 0 {
				return true
			}
		} else if sys.pause > 0 {
			if p.pausemovetime == 0 {
				return true
			}
		}
	}
	return false
}
func (p *Projectile) update(playerNo int) {
	if sys.tickFrame() {
		rem := true
		if p.anim >= 0 {
			if p.hits < 0 && p.remove {
				if p.hits == -1 {
					if p.hitanim != p.anim {
						p.ani = sys.chars[playerNo][0].getAnim(p.hitanim, false)
					}
				} else if p.cancelanim != p.anim {
					p.ani = sys.chars[playerNo][0].getAnim(p.cancelanim, false)
				}
			} else if p.pos[0] < sys.xmin/p.localscl-float32(p.edgebound) ||
				p.pos[0] > sys.xmax/p.localscl+float32(p.edgebound) ||
				p.velocity[0]*p.facing < 0 &&
					p.pos[0] < sys.cam.XMin/p.localscl-float32(p.stagebound) ||
				p.velocity[0]*p.facing > 0 &&
					p.pos[0] > sys.cam.XMax/p.localscl+float32(p.stagebound) ||
				p.velocity[1] > 0 && p.pos[1] > float32(p.heightbound[1]) ||
				p.velocity[1] < 0 && p.pos[1] < float32(p.heightbound[0]) ||
				p.removetime == 0 ||
				p.removetime <= -2 && (p.ani == nil || p.ani.loopend) {
				if p.remanim != p.anim {
					p.ani = sys.chars[playerNo][0].getAnim(p.remanim, false)
				}
			} else {
				rem = false
			}
			if rem {
				if p.ani != nil {
					p.ani.UpdateSprite()
				}
				p.velocity = p.remvelocity
				p.accel, p.velmul, p.anim = [2]float32{}, [...]float32{1, 1}, -1
				if p.hits >= 0 {
					p.hits = -1
				}
			}
		}
		if rem {
			if p.ani != nil && (p.ani.totaltime <= 0 || p.ani.AnimTime() == 0) {
				p.ani = nil
			}
			if p.ani == nil && p.id >= 0 {
				p.id = ^p.id
			}
		}
	}
	if p.paused(playerNo) || p.hitpause != 0 {
		return
	}
	if sys.tickFrame() {
		p.oldPos = p.pos
		p.newPos = [...]float32{p.pos[0] + p.velocity[0]*p.facing,
			p.pos[1] + p.velocity[1]}
	}
	ti := sys.tickInterpola()
	for i, np := range p.newPos {
		p.pos[i] = np - (np-p.oldPos[i])*(1-ti)
	}
	if sys.tickNextFrame() {
		for i := range p.velocity {
			p.velocity[i] += p.accel[i]
			p.velocity[i] *= p.velmul[i]
		}
		if p.velocity[0] < 0 {
			p.facing *= -1
			p.velocity[0] *= -1
			p.accel[0] *= -1
		}
	}
}
func (p *Projectile) clsn(playerNo int) {
	if p.ani == nil || len(p.ani.frames) == 0 {
		return
	}
	cancel := func(prioritypoint *int32, priority int32, hits *int32,
		oppprioritypoint int32) {
		if oppprioritypoint > *prioritypoint || *hits <= 0 {
			*hits = -2
			return
		}
		if oppprioritypoint == *prioritypoint {
			if *hits <= 1 {
				*hits = -2
			} else {
				*hits = int32(int64(*hits)*3/4 - 1)
				if *hits <= 0 {
					(*prioritypoint)--
				} else {
					*prioritypoint = priority
				}
			}
		} else {
			(*prioritypoint)--
		}
	}
	for i := 0; i < playerNo && p.hits >= 0; i++ {
		for j, pr := range sys.projs[i] {
			if pr.hits < 0 || pr.id < 0 || pr.hitdef.affectteam != 0 &&
				(playerNo&1 != i&1) != (pr.hitdef.affectteam > 0) ||
				pr.ani == nil || len(pr.ani.frames) == 0 {
				continue
			}
			clsn1 := pr.ani.CurrentFrame().Clsn2()
			clsn2 := p.ani.CurrentFrame().Clsn2()
			if sys.clsnHantei(clsn1, [...]float32{pr.clsnScale[0] * pr.localscl, pr.clsnScale[1] * pr.localscl},
				[...]float32{pr.pos[0] * pr.localscl, pr.pos[1] * pr.localscl}, pr.facing,
				clsn2, [...]float32{p.clsnScale[0] * p.localscl, p.clsnScale[1] * p.localscl},
				[...]float32{p.pos[0] * p.localscl, p.pos[1] * p.localscl}, p.facing) {
				opp, pp := &sys.projs[i][j], p.prioritypoint
				cancel(&p.prioritypoint, p.priority, &p.hits, opp.prioritypoint)
				cancel(&opp.prioritypoint, opp.priority, &opp.hits, pp)
				if p.hits < 0 {
					break
				}
			}
		}
	}
}
func (p *Projectile) tick(playerNo int) {
	if p.timemiss < 0 {
		p.timemiss = ^p.timemiss
		if p.hits >= 0 {
			if p.timemiss <= 0 && p.hitpause == 0 {
				p.hits = -1
			} else {
				p.hits--
				if p.hits <= 0 {
					p.hits = -1
				}
			}
		}
		p.hitdef.air_juggle = 0
	}
	if p.hits < 0 {
		p.hitpause = 0
	}
	if !p.paused(playerNo) {
		if p.hitpause <= 0 {
			if p.removetime > 0 {
				p.removetime--
			}
			if p.timemiss > 0 {
				p.timemiss--
			}
		}
	}
}
func (p *Projectile) cueDraw(oldVer bool, playerNo int) {
	notpause := p.hitpause <= 0 && !p.paused(playerNo)
	if sys.tickFrame() && p.ani != nil && notpause {
		p.ani.UpdateSprite()
	}
	if sys.clsnDraw && p.ani != nil {
		if frm := p.ani.drawFrame(); frm != nil {
			xs := p.facing * p.clsnScale[0] * p.localscl
			if clsn := frm.Clsn1(); len(clsn) > 0 {
				sys.drawc1.Add(clsn, p.pos[0]*p.localscl, p.pos[1]*p.localscl, xs, p.clsnScale[1]*p.localscl)
			}
			if clsn := frm.Clsn2(); len(clsn) > 0 {
				sys.drawc2.Add(clsn, p.pos[0]*p.localscl, p.pos[1]*p.localscl, xs, p.clsnScale[1]*p.localscl)
			}
		}
	}
	if sys.tickNextFrame() && (notpause || !p.paused(playerNo)) {
		if p.ani != nil && notpause {
			p.ani.Action()
		}
		if p.hitpause > 0 {
			p.hitpause--
		} else {
			if p.supermovetime > 0 {
				p.supermovetime--
			}
			if p.pausemovetime > 0 {
				p.pausemovetime--
			}
		}
	}
	if p.ani != nil {
		sd := &SprData{p.ani, p.palfx, [...]float32{p.pos[0] * p.localscl, p.pos[1] * p.localscl},
			[...]float32{p.facing * p.scale[0] * p.localscl, p.scale[1] * p.localscl}, [2]int32{-1},
			p.sprpriority, 0, [...]float32{1, 1}, false, playerNo == sys.superplayer,
			sys.cgi[playerNo].ver[0] != 1, p.facing}
		p.aimg.recAndCue(sd, sys.tickNextFrame() && notpause)
		sys.sprites.add(sd,
			p.shadow[0]<<16|p.shadow[1]&255<<8|p.shadow[2]&255, 256, 0, 0)
	}
}

type MoveContact int32

const (
	MC_Hit MoveContact = iota
	MC_Guarded
	MC_Reversed
)

type ProjContact int32

const (
	PC_Hit ProjContact = iota
	PC_Guarded
	PC_Cancel
)

type CharGlobalInfo struct {
	def              string
	displayname      string
	author           string
	nameLow          string
	authorLow        string
	palkeymap        [MaxPalNo]int32
	sff              *Sff
	snd              *Snd
	anim             AnimationTable
	palno, drawpalno int32
	pal              [MaxPalNo]string
	palExist         [MaxPalNo]bool
	palSelectable    [MaxPalNo]bool
	ver              [2]uint16
	data             CharData
	velocity         CharVelocity
	movement         CharMovement
	states           map[int32]StateBytecode
	wakewakaLength   int32
	pctype           ProjContact
	pctime, pcid     int32
	unhittable       int32
	quotes           [MaxQuotes]string
}

func (cgi *CharGlobalInfo) clearPCTime() {
	cgi.pctype = PC_Hit
	cgi.pctime = -1
	cgi.pcid = 0
}

type StateState struct {
	stateType       StateType
	moveType        MoveType
	physics         StateType
	ps              []int32
	wakegawakaranai [len(sys.cgi)][]bool
	no, prevno      int32
	time            int32
	sb              StateBytecode
}

func (ss *StateState) clear() {
	ss.stateType, ss.moveType, ss.physics = ST_S, MT_I, ST_N
	ss.ps = nil
	for i, v := range ss.wakegawakaranai {
		if len(v) < int(sys.cgi[i].wakewakaLength) {
			ss.wakegawakaranai[i] = make([]bool, sys.cgi[i].wakewakaLength)
		} else {
			for i := range v {
				v[i] = false
			}
		}
	}
	ss.clearWw()
	ss.no, ss.prevno = 0, 0
	ss.time = 0
	ss.sb = StateBytecode{}
}
func (ss *StateState) clearWw() {
	for _, v := range ss.wakegawakaranai {
		for i := range v {
			v[i] = false
		}
	}
}

type HMF int32

const (
	HMF_H HMF = iota
	HMF_M
	HMF_F
)

type CharSystemVar struct {
	airJumpCount  int32
	hitCount      int32
	uniqHitCount  int32
	pauseMovetime int32
	superMovetime int32
	bindTime      int32
	bindToId      int32
	bindPos       [2]float32
	bindFacing    float32
	hitPauseTime  int32
	angle         float32
	angleScalse   [2]float32
	alpha         [2]int32
	recoverTime   int32
	systemFlag    SystemCharFlag
	specialFlag   CharSpecialFlag
	sprPriority   int32
	getcombo      int32
	veloff        float32
	width, edge   [2]float32
	attackMul     float32
	defenceMul    float32
}
type Char struct {
	name            string
	palfx           *PalFX
	anim            *Animation
	curFrame        *AnimFrame
	cmd             []CommandList
	ss              StateState
	key             int
	id              int32
	helperId        int32
	helperIndex     int32
	parentIndex     int32
	playerNo        int
	keyctrl         bool
	player          bool
	animPN          int
	animNo          int32
	life            int32
	lifeMax         int32
	power           int32
	powerMax        int32
	juggle          int32
	fallTime        int32
	localcoord      int32
	localscl        float32
	oldLocalscl     float32
	size            CharSize
	clsnScale       [2]float32
	hitdef          HitDef
	ghv             GetHitVar
	hitby           [2]HitBy
	ho              [8]HitOverride
	hoIdx           int
	mctype          MoveContact
	mctime          int32
	children        []*Char
	targets         []int32
	targetsOfHitdef []int32
	enemynear       [2][]*Char
	pos             [2]float32
	drawPos         [2]float32
	oldPos          [2]float32
	vel             [2]float32
	facing          float32
	ivar            [NumVar + NumSysVar]int32
	fvar            [NumFvar + NumSysFvar]float32
	CharSystemVar
	aimg          AfterImage
	sounds        Sounds
	p1facing      float32
	cpucmd        int32
	attackDist    float32
	offset        [2]float32
	stchtmp       bool
	inguarddist   bool
	pushed        bool
	hitdefContact bool
	atktmp        int8
	hittmp        int8
	acttmp        int8
	minus         int8
	winquote      int32
	memberNo      int
	selectNo      int
}

func newChar(n int, idx int32) (c *Char) {
	c = &Char{aimg: *newAfterImage()}
	c.init(n, idx)
	return c
}
func (c *Char) panic() {
	if sys.workingState != &c.ss.sb {
		sys.errLog.Panicf("%v\n%v\n%v\n%+v\n", c.gi().def, c.name,
			sys.cgi[sys.workingState.playerNo].def, sys.workingState)
	}
	sys.errLog.Panicf("%v\n%v\n%v\n%+v\n", c.gi().def, c.name,
		sys.cgi[c.ss.sb.playerNo].def, c.ss)
}
func (c *Char) init(n int, idx int32) {
	c.clear1()
	c.playerNo, c.helperIndex = n, idx
	c.animPN = c.playerNo
	if c.helperIndex == 0 {
		c.keyctrl, c.player = true, true
	}
	c.key = n
	if n >= 0 && n < len(sys.com) && sys.com[n] != 0 {
		c.key ^= -1
	}
}
func (c *Char) clearState() {
	c.ss.clear()
	c.hitdef.clear()
	c.ghv.clear()
	c.ghv.clearOff()
	c.hitby = [2]HitBy{}
	for i := range c.ho {
		c.ho[i].clear()
	}
	c.mctype = MC_Hit
	c.mctime = 0
	c.fallTime = 0
	c.hitdefContact = false
}
func (c *Char) clear1() {
	c.anim = nil
	c.cmd = nil
	c.curFrame = nil
	c.clearState()
	c.hoIdx = -1
	c.mctype, c.mctime = MC_Hit, 0
	c.fallTime = 0
	c.varRangeSet(0, int32(NumVar)-1, 0)
	c.fvarRangeSet(0, int32(NumFvar)-1, 0)
	c.key = -1
	c.id = -1
	c.helperId = 0
	c.helperIndex = -1
	c.parentIndex = IErr
	c.playerNo = -1
	c.facing = 1
	c.keyctrl = false
	c.player = false
	c.animPN = -1
	c.animNo = 0
	c.stchtmp = false
	c.inguarddist = false
	c.p1facing = 0
	c.pushed = false
	c.atktmp, c.hittmp, c.acttmp, c.minus = 0, 0, 0, 2
	c.winquote = -1
}
func (c *Char) copyParent(p *Char) {
	c.parentIndex = p.helperIndex
	c.name, c.key, c.size = p.name+"'s helper", p.key, p.size
	c.life, c.lifeMax, c.power, c.powerMax = p.lifeMax, p.lifeMax, 0, p.powerMax
	c.clear2()
}
func (c *Char) addChild(ch *Char) {
	for i, chi := range c.children {
		if chi == nil {
			c.children[i] = ch
			return
		}
	}
	c.children = append(c.children, ch)
}
func (c *Char) enemyNearClear() {
	c.enemynear[0] = c.enemynear[0][:0]
	c.enemynear[1] = c.enemynear[1][:0]
}
func (c *Char) clear2() {
	c.sysVarRangeSet(0, int32(NumSysVar)-1, 0)
	c.sysFvarRangeSet(0, int32(NumSysFvar)-1, 0)
	c.CharSystemVar = CharSystemVar{bindToId: -1,
		angleScalse: [...]float32{1, 1}, alpha: [...]int32{255, 0},
		width:      [...]float32{c.defFW(), c.defBW()},
		attackMul:  float32(c.gi().data.attack) / 100,
		defenceMul: float32(c.gi().data.defence) / 100}
	c.oldPos, c.drawPos = c.pos, c.pos
	if c.helperIndex == 0 {
		if sys.roundsExisted[c.playerNo&1] > 0 {
			c.palfx.clear()
		} else {
			c.palfx = newPalFX()
		}
	} else {
		c.palfx = nil
	}
	c.aimg.timegap = -1
	c.enemyNearClear()
	c.targets = c.targets[:0]
	c.cpucmd = -1
}
func (c *Char) gi() *CharGlobalInfo {
	return &sys.cgi[c.playerNo]
}
func (c *Char) stCgi() *CharGlobalInfo {
	return &sys.cgi[c.ss.sb.playerNo]
}
func (c *Char) load(def string) error {
	gi := &sys.cgi[c.playerNo]
	gi.def, gi.displayname, gi.author, gi.sff, gi.snd, gi.quotes = def, "", "", nil, nil, [MaxQuotes]string{}
	gi.anim = NewAnimationTable()
	for i := range gi.palkeymap {
		gi.palkeymap[i] = int32(i)
	}
	str, err := LoadText(def)
	if err != nil {
		return err
	}
	lines, i := SplitAndTrim(str, "\n"), 0
	cns, sprite, anim, sound := "", "", "", ""
	info, files, keymap := true, true, true
	c.localcoord = 320
	c.localscl = 1
	for i < len(lines) {
		is, name, subname := ReadIniSection(lines, &i)
		switch name {
		case "info":
			if info {
				info = false
				c.name, _, _ = is.getText("name")
				var ok bool
				gi.displayname, ok, _ = is.getText("displayname")
				if !ok {
					gi.displayname = c.name
				}
				gi.author, _, _ = is.getText("author")
				gi.authorLow = strings.ToLower(gi.author)
				gi.nameLow = strings.ToLower(c.name)
				is.ReadI32("localcoord", &c.localcoord)
				c.localscl = 320 / float32(c.localcoord)
			}
		case "files":
			if files {
				files = false
				cns, sprite = is["cns"], is["sprite"]
				anim, sound = is["anim"], is["sound"]
				for i := range gi.pal {
					gi.pal[i] = is[fmt.Sprintf("pal%v", i+1)]
				}
			}
		case "palette ":
			if keymap &&
				len(subname) >= 6 && strings.ToLower(subname[:6]) == "keymap" {
				keymap = false
				for i, v := range [12]string{"a", "b", "c", "x", "y", "z",
					"a2", "b2", "c2", "x2", "y2", "z2"} {
					var i32 int32
					if is.ReadI32(v, &i32) {
						if i32 < 1 || int(i32) > len(gi.palkeymap) {
							i32 = 1
						}
						gi.palkeymap[i] = i32 - 1
					}
				}
			}
		}
	}
	if err := LoadFile(&cns, def, func(filename string) error {
		str, err := LoadText(filename)
		if err != nil {
			return err
		}
		lines, i = SplitAndTrim(str, "\n"), 0
		return nil
	}); err != nil {
		return err
	}
	gi.data.init()
	c.size.init()
	c.size.ground.back = c.size.ground.back / c.localscl
	c.size.ground.front = c.size.ground.front / c.localscl
	c.size.air.back = c.size.air.back / c.localscl
	c.size.air.front = c.size.air.front / c.localscl
	c.size.height = c.size.height / c.localscl
	c.size.attack.dist = c.size.attack.dist / c.localscl
	c.size.proj.attack.dist = c.size.proj.attack.dist / c.localscl
	c.size.head.pos[0] = c.size.head.pos[0] / c.localscl
	c.size.head.pos[1] = c.size.head.pos[1] / c.localscl
	c.size.mid.pos[0] = c.size.mid.pos[0] / c.localscl
	c.size.mid.pos[1] = c.size.mid.pos[1] / c.localscl
	c.size.shadowoffset = c.size.shadowoffset / c.localscl
	c.size.draw.offset[0] = c.size.draw.offset[0] / c.localscl
	c.size.draw.offset[1] = c.size.draw.offset[1] / c.localscl
	c.size.z.width = c.size.z.width / c.localscl
	c.size.attack.z.width[0] = c.size.attack.z.width[0] / c.localscl
	c.size.attack.z.width[1] = c.size.attack.z.width[1] / c.localscl

	gi.velocity.init()

	gi.velocity.air.gethit.groundrecover[0] /= c.localscl
	gi.velocity.air.gethit.groundrecover[1] /= c.localscl
	gi.velocity.air.gethit.airrecover.add[0] /= c.localscl
	gi.velocity.air.gethit.airrecover.add[1] /= c.localscl
	gi.velocity.air.gethit.airrecover.back /= c.localscl
	gi.velocity.air.gethit.airrecover.fwd /= c.localscl
	gi.velocity.air.gethit.airrecover.up /= c.localscl
	gi.velocity.air.gethit.airrecover.down /= c.localscl

	gi.movement.init()

	gi.movement.airjump.height = int32(float32(gi.movement.airjump.height) / c.localscl)
	gi.movement.yaccel /= c.localscl
	gi.movement.stand.friction_threshold /= c.localscl
	gi.movement.crouch.friction_threshold /= c.localscl
	gi.movement.air.gethit.groundlevel /= c.localscl
	gi.movement.air.gethit.groundrecover.ground.threshold /= c.localscl
	gi.movement.air.gethit.groundrecover.groundlevel /= c.localscl
	gi.movement.air.gethit.airrecover.threshold /= c.localscl
	gi.movement.air.gethit.airrecover.yaccel /= c.localscl
	gi.movement.air.gethit.trip.groundlevel /= c.localscl
	gi.movement.down.bounce.offset[0] /= c.localscl
	gi.movement.down.bounce.offset[1] /= c.localscl
	gi.movement.down.bounce.yaccel /= c.localscl
	gi.movement.down.bounce.groundlevel /= c.localscl
	gi.movement.down.friction_threshold /= c.localscl

	data, size, velocity, movement, quotes := true, true, true, true, true
	for i < len(lines) {
		is, name, _ := ReadIniSection(lines, &i)
		switch name {
		case "data":
			if data {
				data = false
				var tmp int32
				tmp = Atoi(sys.cmdFlags[fmt.Sprintf("-p%v.life", c.playerNo+1)])
				if tmp != 0 {
					gi.data.life = tmp
				} else {
					is.ReadI32("life", &gi.data.life)
				}
				c.lifeMax = gi.data.life
				tmp = Atoi(sys.cmdFlags[fmt.Sprintf("-p%v.power", c.playerNo+1)])
				if tmp != 0 {
					gi.data.power = tmp
				} else {
					is.ReadI32("power", &gi.data.power)
				}
				c.powerMax = gi.data.power
				is.ReadI32("attack", &gi.data.attack)
				is.ReadI32("defence", &gi.data.defence)
				var i32 int32
				if is.ReadI32("fall.defence_up", &i32) {
					gi.data.fall.defence_mul = (float32(i32) + 100) / 100
				}
				if is.ReadI32("liedown.time", &i32) {
					gi.data.liedown.time = Max(1, i32)
				}
				is.ReadI32("airjuggle", &gi.data.airjuggle)
				is.ReadI32("sparkno", &gi.data.sparkno)
				if gi.data.sparkno < 0 {
					gi.data.sparkno = ^IErr
				}
				is.ReadI32("guard.sparkno", &gi.data.guard.sparkno)
				if gi.data.guard.sparkno < 0 {
					gi.data.guard.sparkno = ^IErr
				}
				is.ReadI32("ko.echo", &gi.data.ko.echo)
				if gi.ver[0] == 1 {
					if is.ReadI32("volumescale", &i32) {
						gi.data.volume = i32 * 64 / 25
					}
				} else if is.ReadI32("volume", &i32) {
					gi.data.volume = i32 + 256
				}
				is.ReadI32("intpersistindex", &gi.data.intpersistindex)
				is.ReadI32("floatpersistindex", &gi.data.floatpersistindex)
			}
		case "size":
			if size {
				size = false

				is.ReadF32("xscale", &c.size.xscale)
				is.ReadF32("yscale", &c.size.yscale)
				is.ReadF32("ground.back", &c.size.ground.back)
				is.ReadF32("ground.front", &c.size.ground.front)
				is.ReadF32("air.back", &c.size.air.back)
				is.ReadF32("air.front", &c.size.air.front)
				is.ReadF32("height", &c.size.height)
				is.ReadF32("attack.dist", &c.size.attack.dist)
				is.ReadF32("proj.attack.dist", &c.size.proj.attack.dist)
				is.ReadI32("proj.doscale", &c.size.proj.doscale)
				is.ReadF32("head.pos", &c.size.head.pos[0], &c.size.head.pos[1])
				is.ReadF32("mid.pos", &c.size.mid.pos[0], &c.size.mid.pos[1])
				is.ReadF32("shadowoffset", &c.size.shadowoffset)
				is.ReadF32("draw.offset",
					&c.size.draw.offset[0], &c.size.draw.offset[1])
				is.ReadF32("z.width", &c.size.z.width)
				is.ReadF32("attack.z.width",
					&c.size.attack.z.width[0], &c.size.attack.z.width[1])
			}
		case "velocity":
			if velocity {
				velocity = false
				is.ReadF32("walk.fwd", &gi.velocity.walk.fwd)
				is.ReadF32("walk.back", &gi.velocity.walk.back)
				is.ReadF32("walk.up.x", &gi.velocity.walk.up.x)
				is.ReadF32("walk.down.x", &gi.velocity.walk.down.x)
				is.ReadF32("run.fwd", &gi.velocity.run.fwd[0], &gi.velocity.run.fwd[1])
				is.ReadF32("run.back",
					&gi.velocity.run.back[0], &gi.velocity.run.back[1])
				is.ReadF32("run.up.x", &gi.velocity.run.up.x)
				is.ReadF32("run.up.y", &gi.velocity.run.up.y)
				is.ReadF32("run.down.x", &gi.velocity.run.down.x)
				is.ReadF32("run.down.y", &gi.velocity.run.down.y)
				is.ReadF32("jump.neu",
					&gi.velocity.jump.neu[0], &gi.velocity.jump.neu[1])
				is.ReadF32("jump.back", &gi.velocity.jump.back)
				is.ReadF32("jump.fwd", &gi.velocity.jump.fwd)
				is.ReadF32("jump.up.x", &gi.velocity.jump.up.x)
				is.ReadF32("jump.down.x", &gi.velocity.jump.down.x)
				is.ReadF32("runjump.back",
					&gi.velocity.runjump.back[0], &gi.velocity.runjump.back[1])
				is.ReadF32("runjump.fwd",
					&gi.velocity.runjump.fwd[0], &gi.velocity.runjump.fwd[1])
				is.ReadF32("runjump.up.x", &gi.velocity.runjump.up.x)
				is.ReadF32("runjump.down.x", &gi.velocity.runjump.down.x)
				is.ReadF32("airjump.neu",
					&gi.velocity.airjump.neu[0], &gi.velocity.airjump.neu[1])
				is.ReadF32("airjump.back", &gi.velocity.airjump.back)
				is.ReadF32("airjump.fwd", &gi.velocity.airjump.fwd)
				is.ReadF32("airjump.up.x", &gi.velocity.airjump.up.x)
				is.ReadF32("airjump.down.x", &gi.velocity.airjump.down.x)
				is.ReadF32("air.gethit.groundrecover",
					&gi.velocity.air.gethit.groundrecover[0],
					&gi.velocity.air.gethit.groundrecover[1])
				is.ReadF32("air.gethit.airrecover.mul",
					&gi.velocity.air.gethit.airrecover.mul[0],
					&gi.velocity.air.gethit.airrecover.mul[1])
				is.ReadF32("air.gethit.airrecover.add",
					&gi.velocity.air.gethit.airrecover.add[0],
					&gi.velocity.air.gethit.airrecover.add[1])
				is.ReadF32("air.gethit.airrecover.back",
					&gi.velocity.air.gethit.airrecover.back)
				is.ReadF32("air.gethit.airrecover.fwd",
					&gi.velocity.air.gethit.airrecover.fwd)
				is.ReadF32("air.gethit.airrecover.up",
					&gi.velocity.air.gethit.airrecover.up)
				is.ReadF32("air.gethit.airrecover.down",
					&gi.velocity.air.gethit.airrecover.down)
			}
		case "movement":
			if movement {
				movement = false
				is.ReadI32("airjump.num", &gi.movement.airjump.num)
				is.ReadI32("airjump.height", &gi.movement.airjump.height)
				is.ReadF32("yaccel", &gi.movement.yaccel)
				is.ReadF32("stand.friction", &gi.movement.stand.friction)
				is.ReadF32("stand.friction.threshold",
					&gi.movement.stand.friction_threshold)
				is.ReadF32("crouch.friction", &gi.movement.crouch.friction)
				is.ReadF32("crouch.friction.threshold",
					&gi.movement.crouch.friction_threshold)
				is.ReadF32("air.gethit.groundlevel",
					&gi.movement.air.gethit.groundlevel)
				is.ReadF32("air.gethit.groundrecover.ground.threshold",
					&gi.movement.air.gethit.groundrecover.ground.threshold)
				is.ReadF32("air.gethit.groundrecover.groundlevel",
					&gi.movement.air.gethit.groundrecover.groundlevel)
				is.ReadF32("air.gethit.airrecover.threshold",
					&gi.movement.air.gethit.airrecover.threshold)
				is.ReadF32("air.gethit.airrecover.yaccel",
					&gi.movement.air.gethit.airrecover.yaccel)
				is.ReadF32("air.gethit.trip.groundlevel",
					&gi.movement.air.gethit.trip.groundlevel)
				is.ReadF32("down.bounce.offset",
					&gi.movement.down.bounce.offset[0],
					&gi.movement.down.bounce.offset[1])
				is.ReadF32("down.bounce.yaccel", &gi.movement.down.bounce.yaccel)
				is.ReadF32("down.bounce.groundlevel",
					&gi.movement.down.bounce.groundlevel)
				is.ReadF32("down.friction.threshold",
					&gi.movement.down.friction_threshold)
			}
		case "quotes":
			if quotes {
				quotes = false
				for i := 0; i < MaxQuotes; i++ {
					if is[fmt.Sprintf("victory%v", i)] != "" {
						gi.quotes[i], _, _ = is.getText(fmt.Sprintf("victory%v", i))
					}
				}
			}
		}
	}
	if LoadFile(&sprite, def, func(filename string) error {
		var err error
		gi.sff, err = loadSff(filename, true)
		return err
	}); err != nil {
		return err
	}
	if LoadFile(&anim, def, func(filename string) error {
		str, err := LoadText(filename)
		if err != nil {
			return err
		}
		str = str + sys.commonAir
		lines, i := SplitAndTrim(str, "\n"), 0
		gi.anim = ReadAnimationTable(gi.sff, lines, &i)
		return nil
	}); err != nil {
		return err
	}
	if len(sound) > 0 {
		if LoadFile(&sound, def, func(filename string) error {
			var err error
			gi.snd, err = LoadSnd(filename)
			return err
		}); err != nil {
			return err
		}
	} else {
		gi.snd = newSnd()
	}
	return nil
}
func (c *Char) loadPallet() {
	if c.gi().sff.header.Ver0 == 1 {
		c.gi().sff.palList.ResetRemap()
		tmp := 0
		for i := 0; i < MaxPalNo; i++ {
			pl := c.gi().sff.palList.Get(i)
			var f *os.File
			var err error
			if LoadFile(&c.gi().pal[i], c.gi().def, func(file string) error {
				f, err = os.Open(file)
				return err
			}) == nil {
				for i := 255; i >= 0; i-- {
					var rgb [3]byte
					if _, err = io.ReadFull(f, rgb[:]); err != nil {
						break
					}
					pl[i] = uint32(rgb[2])<<16 | uint32(rgb[1])<<8 | uint32(rgb[0])
				}
				chk(f.Close())
				if err == nil {
					if tmp == 0 && i > 0 {
						copy(c.gi().sff.palList.Get(0), pl)
					}
					tmp = i + 1
					c.gi().palExist[i] = true
				}
			}
			if err != nil {
				c.gi().palExist[i] = false
				if i > 0 {
					delete(c.gi().sff.palList.PalTable, [...]int16{1, int16(i + 1)})
				}
			}
		}
		if tmp == 0 {
			delete(c.gi().sff.palList.PalTable, [...]int16{1, 1})
		}
	} else {
		for i := 0; i < MaxPalNo; i++ {
			_, c.gi().palExist[i] =
				c.gi().sff.palList.PalTable[[...]int16{1, int16(i + 1)}]
		}
	}
	for i := range c.gi().palSelectable {
		c.gi().palSelectable[i] = false
	}
	for i := 0; i < MaxPalNo; i++ {
		startj := c.gi().palkeymap[i]
		if !c.gi().palExist[startj] {
			startj %= 6
		}
		j := startj
		for {
			if c.gi().palExist[j] {
				c.gi().palSelectable[j] = true
				break
			}
			j++
			if j >= MaxPalNo {
				j = 0
			}
			if j == startj {
				break
			}
		}
	}
	c.gi().drawpalno = c.gi().palno
	starti := c.gi().palno - 1
	if !c.gi().palExist[starti] {
		starti %= 6
	}
	i := starti
	for {
		if c.gi().palExist[i] {
			j := 0
			for ; j < len(sys.chars); j++ {
				if j != c.playerNo && len(sys.chars[j]) > 0 &&
					sys.cgi[j].def == c.gi().def && sys.cgi[j].drawpalno == i+1 {
					break
				}
			}
			if j >= len(sys.chars) {
				c.gi().drawpalno = i + 1
				if !c.gi().palExist[c.gi().palno-1] {
					c.gi().palno = c.gi().drawpalno
				}
				break
			}
		}
		i++
		if i >= MaxPalNo {
			i = 0
		}
		if i == starti {
			if !c.gi().palExist[c.gi().palno-1] {
				i := 0
				for ; i < len(c.gi().palExist); i++ {
					if c.gi().palExist[i] {
						c.gi().palno, c.gi().drawpalno = int32(i+1), int32(i+1)
						break
					}
				}
				if i >= len(c.gi().palExist) {
					c.gi().palno, c.gi().palExist[0] = 1, true
					c.gi().palSelectable[0] = true
				}
			}
			break
		}
	}
}
func (c *Char) clearHitCount() {
	c.hitCount, c.uniqHitCount = 0, 0
}
func (c *Char) clearMoveHit() {
	c.mctime = 0
}
func (c *Char) clearHitDef() {
	c.hitdef.clear()
}
func (c *Char) setSprPriority(sprpriority int32) {
	c.sprPriority = sprpriority
}
func (c *Char) setJuggle(juggle int32) {
	c.juggle = juggle
}
func (c *Char) setXV(xv float32) {
	c.vel[0] = xv
}
func (c *Char) setYV(yv float32) {
	c.vel[1] = yv
}
func (c *Char) changeAnim(animNo int32) {
	if a := c.getAnim(animNo, false); a != nil {
		c.anim = a
		c.animPN = c.playerNo
		c.animNo = animNo
		c.clsnScale = [...]float32{sys.chars[c.animPN][0].size.xscale,
			sys.chars[c.animPN][0].size.yscale}
		if c.hitPause() {
			c.curFrame = a.CurrentFrame()
		}
	}
}
func (c *Char) changeAnim2(animNo int32) {
	if a := sys.chars[c.ss.sb.playerNo][0].getAnim(animNo, false); a != nil {
		c.anim = a
		c.animPN = c.ss.sb.playerNo
		c.animNo = animNo
		c.clsnScale = [...]float32{sys.chars[c.animPN][0].size.xscale,
			sys.chars[c.animPN][0].size.yscale}
		a.sff = sys.cgi[c.playerNo].sff
		if c.hitPause() {
			c.curFrame = a.CurrentFrame()
		}
	}
}
func (c *Char) setAnimElem(e int32) {
	if c.anim != nil {
		c.anim.SetAnimElem(e)
		c.curFrame = c.anim.CurrentFrame()
	}
}
func (c *Char) setCtrl(ctrl bool) {
	if ctrl {
		c.setSCF(SCF_ctrl)
	} else {
		c.unsetSCF(SCF_ctrl)
	}
}
func (c *Char) scf(scf SystemCharFlag) bool {
	return c.systemFlag&scf != 0
}
func (c *Char) setSCF(scf SystemCharFlag) {
	c.systemFlag |= scf
}
func (c *Char) unsetSCF(scf SystemCharFlag) {
	c.systemFlag &^= scf
}
func (c *Char) sf(csf CharSpecialFlag) bool {
	return c.specialFlag&csf != 0
}
func (c *Char) setSF(csf CharSpecialFlag) {
	c.specialFlag |= csf
}
func (c *Char) unsetSF(csf CharSpecialFlag) {
	c.specialFlag &^= csf
}
func (c *Char) parent() *Char {
	if c.parentIndex == IErr {
		return nil
	}
	if c.parentIndex < 0 && !sys.ignoreMostErrors {
		sys.errLog.Println(c.name + " によるすでに削除された親ヘルパーへのリダイレクト")
	}
	return sys.chars[c.playerNo][Abs(c.parentIndex)]
}
func (c *Char) root() *Char {
	if c.helperIndex == 0 {
		return nil
	}
	return sys.chars[c.playerNo][0]
}
func (c *Char) helper(id int32) *Char {
	for _, h := range sys.chars[c.playerNo][1:] {
		if !h.sf(CSF_destroy) && (id <= 0 || id == h.helperId) {
			return h
		}
	}
	return nil
}
func (c *Char) target(id int32) *Char {
	for _, tid := range c.targets {
		if t := sys.playerID(tid); t != nil && (id < 0 || id == t.ghv.hitid) {
			return t
		}
	}
	return nil
}
func (c *Char) partner(n int32) *Char {
	n = Max(0, n)
	if int(n) > len(sys.chars)/2-2 {
		return nil
	}
	var p int
	if int(n) == c.playerNo>>1 {
		p = c.playerNo + 2
	} else {
		p = c.playerNo&1 + int(n)<<1
		if int(n) > c.playerNo>>1 {
			p += 2
		}
	}
	if len(sys.chars[p]) > 0 {
		return sys.chars[p][0]
	}
	return nil
}
func (c *Char) enemy(n int32) *Char {
	if n < 0 || n >= c.numEnemy() {
		return nil
	}
	return sys.chars[n*2+int32(^c.playerNo&1)][0]
}
func (c *Char) enemyNear(n int32) *Char {
	return sys.charList.enemyNear(c, n, false)
}
func (c *Char) p2() *Char {
	p2 := sys.charList.enemyNear(c, 0, true)
	if p2 != nil && p2.scf(SCF_ko) && p2.scf(SCF_over) {
		return nil
	}
	return p2
}
func (c *Char) aiLevel() int32 {
	if c.helperIndex != 0 && c.gi().ver[0] == 1 {
		return 0
	}
	return sys.com[c.playerNo]
}
func (c *Char) alive() bool {
	return !c.scf(SCF_ko)
}
func (c *Char) animElemNo(time int32) BytecodeValue {
	if c.anim != nil && time >= -c.anim.sumtime {
		return BytecodeInt(c.anim.AnimElemNo(time))
	}
	return BytecodeSF()
}
func (c *Char) animElemTime(e int32) BytecodeValue {
	if e >= 1 && c.anim != nil && int(e) <= len(c.anim.frames) {
		return BytecodeInt(c.anim.AnimElemTime(e))
	}
	return BytecodeSF()
}
func (c *Char) animExist(wc *Char, anim BytecodeValue) BytecodeValue {
	if anim.IsSF() {
		return BytecodeSF()
	}
	if c != wc {
		return c.selfAnimExist(anim)
	}
	return sys.chars[c.ss.sb.playerNo][0].selfAnimExist(anim)
}
func (c *Char) animTime() int32 {
	if c.anim != nil {
		return c.anim.AnimTime()
	}
	return 0
}
func (c *Char) backEdge() float32 {
	if c.facing < 0 {
		return c.rightEdge()
	}
	return c.leftEdge()
}
func (c *Char) backEdgeBodyDist() float32 {
	return c.backEdgeDist() - c.getEdge(c.edge[1], false)
}
func (c *Char) backEdgeDist() float32 {
	if c.facing < 0 {
		return sys.xmax/c.localscl - c.pos[0]
	}
	return c.pos[0] - sys.xmin/c.localscl
}
func (c *Char) bottomEdge() float32 {
	return sys.cam.ScreenPos[1]/c.localscl + c.gameHeight()
}
func (c *Char) canRecover() bool {
	return c.ghv.fall.recover && c.fallTime >= c.ghv.fall.recovertime
}
func (c *Char) command(pn, i int) bool {
	if !c.keyctrl || c.cmd == nil {
		return false
	}
	cl := c.cmd[pn].At(i)
	if len(cl) > 0 && c.key < 0 {
		if c.gi().ver[0] == 1 && c.helperIndex != 0 {
			return false
		}
		if c.helperIndex != 0 || len(cl[0].cmd) != 1 || len(cl[0].cmd[0].key) !=
			1 || int(Btoi(cl[0].cmd[0].slash)) != len(cl[0].hold) {
			return i == int(c.cpucmd)
		}
		if c.helperIndex != 0 {
			return false
		}
	}
	for _, c := range cl {
		if c.curbuftime > 0 {
			return true
		}
	}
	return false
}
func (c *Char) commandByName(name string) bool {
	if c.cmd == nil {
		return false
	}
	i, ok := c.cmd[c.playerNo].Names[name]
	return ok && c.command(c.playerNo, i)
}
func (c *Char) ctrl() bool {
	return c.scf(SCF_ctrl) && !c.ctrlOver()
}
func (c *Char) drawgame() bool {
	return c.roundState() >= 3 && sys.winTeam < 0
}
func (c *Char) frontEdge() float32 {
	if c.facing > 0 {
		return c.rightEdge()
	}
	return c.leftEdge()
}
func (c *Char) frontEdgeBodyDist() float32 {
	return c.frontEdgeDist() - c.getEdge(c.edge[0], false)
}
func (c *Char) frontEdgeDist() float32 {
	if c.facing > 0 {
		return sys.xmax/c.localscl - c.pos[0]
	}
	return c.pos[0] - sys.xmin/c.localscl
}
func (c *Char) gameHeight() float32 {
	return 240 / c.localscl / sys.cam.Scale
}
func (c *Char) gameWidth() float32 {
	return float32(sys.gameWidth) / c.localscl / sys.cam.Scale
}
func (c *Char) hitDefAttr(attr int32) bool {
	return c.ss.moveType == MT_A && c.hitdef.testAttr(attr)
}
func (c *Char) hitOver() bool {
	return c.ghv.hittime < 0
}
func (c *Char) hitShakeOver() bool {
	return c.ghv.hitshaketime <= 0
}
func (c *Char) hitVelX() float32 {
	if c.ss.moveType != MT_H {
		return 0
	}
	return -c.ghv.xvel
}
func (c *Char) hitVelY() float32 {
	if c.ss.moveType != MT_H {
		return 0
	}
	return -c.ghv.yvel
}
func (c *Char) isHelper(hid BytecodeValue) BytecodeValue {
	if hid.IsSF() {
		return BytecodeSF()
	}
	id := hid.ToI()
	return BytecodeBool(c.helperIndex != 0 && (id <= 0 || c.helperId == id))
}
func (c *Char) leftEdge() float32 {
	return sys.cam.ScreenPos[0] / c.localscl
}
func (c *Char) lose() bool {
	return sys.winTeam == ^c.playerNo&1
}
func (c *Char) loseKO() bool {
	return c.lose() && sys.finish == FT_KO
}
func (c *Char) loseTime() bool {
	return c.lose() && sys.finish == FT_TO
}
func (c *Char) moveContact() int32 {
	if c.mctype != MC_Reversed {
		return Abs(c.mctime)
	}
	return 0
}
func (c *Char) moveGuarded() int32 {
	if c.mctype == MC_Guarded {
		return Abs(c.mctime)
	}
	return 0
}
func (c *Char) moveHit() int32 {
	if c.mctype == MC_Hit {
		return Abs(c.mctime)
	}
	return 0
}
func (c *Char) moveReversed() int32 {
	if c.mctype == MC_Reversed {
		return Abs(c.mctime)
	}
	return 0
}
func (c *Char) numEnemy() int32 {
	if sys.tmode[^c.playerNo&1] != TM_Simul {
		return 1
	}
	return sys.numSimul[^c.playerNo&1]
}
func (c *Char) numExplod(eid BytecodeValue) BytecodeValue {
	if eid.IsSF() {
		return BytecodeSF()
	}
	var id, n int32 = eid.ToI(), 0
	for _, e := range sys.explods[c.playerNo] {
		if e.matchId(id, c.id) {
			n++
		}
	}
	return BytecodeInt(n)
}
func (c *Char) numHelper(hid BytecodeValue) BytecodeValue {
	if hid.IsSF() {
		return BytecodeSF()
	}
	var id, n int32 = hid.ToI(), 0
	for _, h := range sys.chars[c.playerNo][1:] {
		if !h.sf(CSF_destroy) && (id <= 0 || h.helperId == id) {
			n++
		}
	}
	return BytecodeInt(n)
}
func (c *Char) numPartner() int32 {
	if sys.tmode[c.playerNo&1] != TM_Simul {
		return 0
	}
	return sys.numSimul[c.playerNo&1] - 1
}
func (c *Char) numProj() int32 {
	n := int32(0)
	for _, p := range sys.projs[c.playerNo] {
		if p.id >= 0 && p.hits >= 0 {
			n++
		}
	}
	return n
}
func (c *Char) numProjID(pid BytecodeValue) BytecodeValue {
	if pid.IsSF() {
		return BytecodeSF()
	}
	if c.helperIndex != 0 {
		return BytecodeInt(0)
	}
	var id, n int32 = Max(0, pid.ToI()), 0
	for _, p := range sys.projs[c.playerNo] {
		if p.id == id && p.hits >= 0 {
			n++
		}
	}
	return BytecodeInt(n)
}
func (c *Char) numTarget(hid BytecodeValue) BytecodeValue {
	if hid.IsSF() {
		return BytecodeSF()
	}
	var id, n int32 = hid.ToI(), 0
	for _, tid := range c.targets {
		if tid >= 0 {
			if id < 0 {
				n++
			} else if t := sys.playerID(tid); t != nil && t.ghv.hitid == id {
				n++
			}
		}
	}
	return BytecodeInt(n)
}
func (c *Char) palno() int32 {
	if c.helperIndex != 0 && c.gi().ver[0] != 1 {
		return 1
	}
	return c.gi().palno
}
func (c *Char) getPower() int32 {
	if sys.powerShare[c.playerNo&1] {
		return sys.chars[c.playerNo&1][0].power
	}
	return sys.chars[c.playerNo][0].power
}
func (c *Char) projCancelTime(pid BytecodeValue) BytecodeValue {
	if pid.IsSF() {
		return BytecodeSF()
	}
	id := pid.ToI()
	if id > 0 && id != c.gi().pcid || c.gi().pctype != PC_Cancel {
		return BytecodeInt(-1)
	}
	return BytecodeInt(c.gi().pctime)
}
func (c *Char) projContactTime(pid BytecodeValue) BytecodeValue {
	if pid.IsSF() {
		return BytecodeSF()
	}
	id := pid.ToI()
	if id > 0 && id != c.gi().pcid {
		return BytecodeInt(-1)
	}
	return BytecodeInt(c.gi().pctime)
}
func (c *Char) projGuardedTime(pid BytecodeValue) BytecodeValue {
	if pid.IsSF() {
		return BytecodeSF()
	}
	id := pid.ToI()
	if id > 0 && id != c.gi().pcid || c.gi().pctype != PC_Guarded {
		return BytecodeInt(-1)
	}
	return BytecodeInt(c.gi().pctime)
}
func (c *Char) projHitTime(pid BytecodeValue) BytecodeValue {
	if pid.IsSF() {
		return BytecodeSF()
	}
	id := pid.ToI()
	if id > 0 && id != c.gi().pcid || c.gi().pctype != PC_Hit {
		return BytecodeInt(-1)
	}
	return BytecodeInt(c.gi().pctime)
}
func (c *Char) rightEdge() float32 {
	return sys.cam.ScreenPos[0]/c.localscl + c.gameWidth()
}
func (c *Char) roundsExisted() int32 {
	return sys.roundsExisted[c.playerNo&1]
}
func (c *Char) roundState() int32 {
	switch {
	case sys.intro > sys.lifebar.ro.ctrl_time+1:
		return 0
	case sys.lifebar.ro.cur == 0:
		return 1
	case !sys.roundEnd():
		return 2
	case sys.intro < -(sys.lifebar.ro.over_hittime+
		sys.lifebar.ro.over_waittime) && (sys.chars[c.playerNo][0].scf(SCF_over) ||
		sys.chars[c.playerNo][0].scf(SCF_ko)):
		return 4
	default:
		return 3
	}
}
func (c *Char) screenPosX() float32 {
	return (c.pos[0]*c.localscl - sys.cam.ScreenPos[0]) // * sys.cam.Scale
}
func (c *Char) screenPosY() float32 {
	return (c.pos[1]*c.localscl - sys.cam.ScreenPos[1]) // * sys.cam.Scale
}
func (c *Char) selfAnimExist(anim BytecodeValue) BytecodeValue {
	if anim.IsSF() {
		return BytecodeSF()
	}
	return BytecodeBool(c.gi().anim.get(anim.ToI()) != nil)
}
func (c *Char) time() int32 {
	return c.ss.time
}
func (c *Char) topEdge() float32 {
	return sys.cam.ScreenPos[1] / c.localscl
}
func (c *Char) win() bool {
	return sys.winTeam == c.playerNo&1
}
func (c *Char) winKO() bool {
	return c.win() && sys.finish == FT_KO
}
func (c *Char) winTime() bool {
	return c.win() && sys.finish == FT_TO
}
func (c *Char) winPerfect() bool {
	return c.win() && sys.winType[c.playerNo&1] >= WT_PN
}
func (c *Char) newChannel(ch int32, lowpriority bool) *Sound {
	ch = Min(255, ch)
	if ch >= 0 {
		if lowpriority {
			if len(c.sounds) > int(ch) && c.sounds[ch].sound != nil {
				return nil
			}
		}
		if len(c.sounds) < int(ch+1) {
			c.sounds = append(c.sounds, newSounds(int(ch+1)-len(c.sounds))...)
		}
		return &c.sounds[ch]
	}
	if len(c.sounds) < 256 {
		c.sounds = append(c.sounds, newSounds(256-len(c.sounds))...)
	}
	for i := 255; i >= 0; i-- {
		if c.sounds[i].sound == nil {
			return &c.sounds[i]
		}
	}
	return nil
}
func (c *Char) playSound(f, lowpriority, loop bool, g, n, chNo, vol int32,
	_, freqmul float32, _ *float32) {
	if g < 0 {
		return
	}
	var w *Wave
	if f {
		if sys.lifebar.fsnd != nil {
			w = sys.lifebar.fsnd.Get([...]int32{g, n})
		}
	} else {
		if c.gi().snd != nil {
			w = c.gi().snd.Get([...]int32{g, n})
		}
	}
	if w == nil && !sys.ignoreMostErrors {
		str := "存在しないサウンド: "
		if f {
			str += "F:"
		} else {
			str += fmt.Sprintf("P%v:", c.playerNo+1)
		}
		sys.errLog.Printf("%v%v,%v\n", str, g, n)
		return
	}
	if ch := c.newChannel(chNo, lowpriority); ch != nil {
		ch.sound, ch.loop, ch.freqmul = w, loop, freqmul
		vol = Max(-25600, Min(25600, vol))
		if c.gi().ver[0] == 1 {
			if f {
				ch.SetVolume(256)
			} else {
				ch.SetVolume(c.gi().data.volume * vol / 100)
			}
		} else {
			if f {
				ch.SetVolume(vol + 256)
			} else {
				ch.SetVolume(c.gi().data.volume + vol)
			}
		}
	}
}
func (c *Char) furimuki() {
	if c.scf(SCF_ctrl) && c.helperIndex == 0 {
		if c.rdDistX(sys.charList.enemyNear(c, 0, true), c).ToF() < 0 {
			switch c.ss.stateType {
			case ST_S:
				c.changeAnim(5)
			case ST_C:
				c.changeAnim(6)
			}
			c.setFacing(-c.facing)
		}
	}
}
func (c *Char) stateChange1(no int32, pn int) bool {
	if sys.changeStateNest >= 2500 {
		sys.errLog.Printf("2500 loops: %v, %v -> %v -> %v\n",
			c.name, c.ss.prevno, c.ss.no, no)
		return false
	}
	c.ss.no, c.ss.prevno, c.ss.time = Max(0, no), c.ss.no, 0
	if c.ss.sb.playerNo != c.playerNo && pn != c.ss.sb.playerNo {
		c.enemyExplodsRemove(c.ss.sb.playerNo)
	}
	if c.localscl != 320/float32(sys.chars[pn][0].localcoord) {
		c.pos[0] *= c.localscl / (320 / float32(sys.chars[pn][0].localcoord))
		c.pos[1] *= c.localscl / (320 / float32(sys.chars[pn][0].localcoord))
		c.oldPos = c.pos

		c.ghv.xvel *= c.localscl / (320 / float32(sys.chars[pn][0].localcoord))
		c.ghv.yvel *= c.localscl / (320 / float32(sys.chars[pn][0].localcoord))
		c.ghv.fall.xvelocity *= c.localscl / (320 / float32(sys.chars[pn][0].localcoord))
		c.ghv.fall.yvelocity *= c.localscl / (320 / float32(sys.chars[pn][0].localcoord))
		c.ghv.yaccel *= c.localscl / (320 / float32(sys.chars[pn][0].localcoord))

		c.localscl = (320 / float32(sys.chars[pn][0].localcoord))
	}
	var ok bool
	if c.ss.sb, ok = sys.cgi[pn].states[no]; !ok {
		sys.errLog.Printf("存在しないステート: P%v:%v\n", pn+1, no)
		c.ss.sb = *newStateBytecode(pn)
		c.ss.sb.stateType, c.ss.sb.moveType, c.ss.sb.physics = ST_U, MT_U, ST_U
	}
	c.stchtmp = true
	return true
}
func (c *Char) stateChange2() bool {
	if c.stchtmp && !c.hitPause() {
		c.ss.sb.init(c)
		c.stchtmp = false
		return true
	}
	return false
}
func (c *Char) changeStateEx(no int32, pn int, anim, ctrl int32) {
	if c.minus <= 0 && (c.ss.stateType == ST_S || c.ss.stateType == ST_C) {
		c.furimuki()
	}
	if anim >= 0 {
		c.changeAnim(anim)
	}
	if ctrl >= 0 {
		c.setCtrl(ctrl != 0)
	}
	if c.stateChange1(no, pn) && sys.changeStateNest == 0 && c.minus == 0 {
		for c.stchtmp && sys.changeStateNest < 2500 {
			c.stateChange2()
			sys.changeStateNest++
			if !c.ss.sb.run(c) {
				break
			}
		}
		sys.changeStateNest = 0
	}
}
func (c *Char) changeState(no, anim, ctrl int32) {
	c.changeStateEx(no, c.ss.sb.playerNo, anim, ctrl)
}
func (c *Char) selfState(no, anim, ctrl int32) {
	c.changeStateEx(no, c.playerNo, anim, ctrl)
}
func (c *Char) destroy() {
	if c.helperIndex > 0 {
		c.exitTarget(true)
		c.getcombo = 0
		for _, tid := range c.targets {
			if t := sys.playerID(tid); t != nil {
				t.gethitBindClear()
				t.ghv.dropId(c.id)
			}
		}
		if c.parentIndex >= 0 {
			if p := c.parent(); p != nil {
				for i, ch := range p.children {
					if ch == c {
						p.children[i] = nil
					}
				}
			}
		}
		for _, ch := range c.children {
			if ch != nil {
				ch.parentIndex *= -1
			}
		}
		c.children = c.children[:0]
		sys.charList.delete(c)
		c.helperIndex = -1
		c.setSF(CSF_destroy)
	}
}
func (c *Char) destroySelf(recursive, removeexplods bool) bool {
	if c.helperIndex <= 0 {
		return false
	}
	c.setSF(CSF_destroy)
	if removeexplods {
		c.removeExplod(-1)
	}
	if recursive {
		for _, ch := range c.children {
			if ch != nil {
				ch.destroySelf(recursive, removeexplods)
			}
		}
	}
	return true
}
func (c *Char) newHelper() (h *Char) {
	i := int32(0)
	for ; int(i) < len(sys.chars[c.playerNo]); i++ {
		if sys.chars[c.playerNo][i].helperIndex < 0 {
			h = sys.chars[c.playerNo][i]
			h.init(c.playerNo, i)
			break
		}
	}
	if int(i) >= len(sys.chars[c.playerNo]) {
		if i >= sys.helperMax {
			return
		}
		h = newChar(c.playerNo, i)
		sys.chars[c.playerNo] = append(sys.chars[c.playerNo], h)
	}
	h.id, h.helperId = sys.newCharId(), 0
	h.copyParent(c)
	c.addChild(h)
	sys.charList.add(h)
	return
}
func (c *Char) helperPos(pt PosType, pos [2]float32, facing int32,
	dstFacing *float32, localscl float32) (p [2]float32) {
	if facing < 0 {
		*dstFacing *= -1
	}
	switch pt {
	case PT_P1:
		p[0] = c.pos[0]*c.localscl/localscl + pos[0]*c.facing
		p[1] = c.pos[1]*c.localscl/localscl + pos[1]
		*dstFacing *= c.facing
	case PT_P2:
		if p2 := sys.charList.enemyNear(c, 0, true); p2 != nil {
			p[0] = p2.pos[0]*p2.localscl/localscl + pos[0]*p2.facing
			p[1] = p2.pos[1]*p2.localscl/localscl + pos[1]
			*dstFacing *= p2.facing
		}
	case PT_F, PT_B:
		if c.facing > 0 && pt == PT_F || c.facing < 0 && pt == PT_B {
			p[0] = c.rightEdge() * c.localscl / localscl
		} else {
			p[0] = c.leftEdge() * c.localscl / localscl
		}
		if c.facing > 0 {
			p[0] += pos[0]
		} else {
			p[0] -= pos[0]
		}
		p[1] = pos[1]
		*dstFacing *= c.facing
	case PT_L:
		p[0] = c.leftEdge()*c.localscl/localscl + pos[0]
		p[1] = pos[1]
	case PT_R:
		p[0] = c.rightEdge()*c.localscl/localscl + pos[0]
		p[1] = pos[1]
	case PT_N:
		p = pos
	}
	return
}
func (c *Char) helperInit(h *Char, st int32, pt PosType, x, y float32,
	facing int32, ownpal bool) {
	p := c.helperPos(pt, [...]float32{x, y}, facing, &h.facing, h.localscl)
	h.setX(p[0])
	h.setY(p[1])
	h.vel = [2]float32{}
	if ownpal {
		h.palfx = newPalFX()
		tmp := c.getPalfx().remap
		h.palfx.remap = make([]int, len(tmp))
		copy(h.palfx.remap, tmp)
	}
	h.changeStateEx(st, c.playerNo, 0, 1)
}
func (c *Char) newExplod() (*Explod, int) {
	explinit := func(expl *Explod) *Explod {
		expl.clear()
		expl.id, expl.playerId, expl.palfx = -1, c.id, c.getPalfx()
		return expl
	}
	for i := range sys.explods[c.playerNo] {
		if sys.explods[c.playerNo][i].id == IErr {
			return explinit(&sys.explods[c.playerNo][i]), i
		}
	}
	i := len(sys.explods[c.playerNo])
	if i < sys.explodMax {
		sys.explods[c.playerNo] = append(sys.explods[c.playerNo], Explod{})
		return explinit(&sys.explods[c.playerNo][i]), i
	}
	return nil, -1
}
func (c *Char) getExplods(id int32) (expls []*Explod) {
	for i, e := range sys.explods[c.playerNo] {
		if e.matchId(id, c.id) {
			expls = append(expls, &sys.explods[c.playerNo][i])
		}
	}
	return
}
func (c *Char) insertExplodEx(i int, rp [2]int32) {
	e := &sys.explods[c.playerNo][i]
	if e.anim == nil {
		e.id = IErr
		return
	}
	e.anim.UpdateSprite()
	if e.ownpal && e.anim.sff != sys.lifebar.fsff {
		remap := make([]int, len(e.palfx.remap))
		copy(remap, e.palfx.remap)
		e.palfx = newPalFX()
		e.palfx.remap = remap
		c.remapPal(e.palfx, [...]int32{1, 1}, rp)
	}
	if e.ontop {
		td := &sys.topexplDrawlist[c.playerNo]
		for ii, te := range *td {
			if te < 0 {
				(*td)[ii] = i
				return
			}
		}
		*td = append(*td, i)
	} else {
		ed := &sys.explDrawlist[c.playerNo]
		for ii, ex := range *ed {
			pid := sys.explods[c.playerNo][ex].playerId
			if pid >= c.id && (pid > c.id || ex < i) {
				*ed = append(*ed, 0)
				copy((*ed)[ii+1:], (*ed)[ii:])
				(*ed)[ii] = i
				return
			}
		}
		*ed = append(*ed, i)
	}
}
func (c *Char) insertExplod(i int) {
	c.insertExplodEx(i, [...]int32{-1, 0})
}
func (c *Char) explodBindTime(id, time int32) {
	for i, e := range sys.explods[c.playerNo] {
		if e.matchId(id, c.id) {
			sys.explods[c.playerNo][i].bindtime = time
		}
	}
}
func (c *Char) removeExplod(id int32) {
	remove := func(drawlist *[]int, drop bool) {
		for i := len(*drawlist) - 1; i >= 0; i-- {
			ei := (*drawlist)[i]
			if ei >= 0 && sys.explods[c.playerNo][ei].matchId(id, c.id) {
				sys.explods[c.playerNo][ei].id = IErr
				if drop {
					*drawlist = append((*drawlist)[:i], (*drawlist)[i+1:]...)
				} else {
					(*drawlist)[i] = -1
				}
			}
		}
	}
	remove(&sys.explDrawlist[c.playerNo], true)
	remove(&sys.topexplDrawlist[c.playerNo], false)
}
func (c *Char) enemyExplodsRemove(en int) {
	remove := func(drawlist *[]int, drop bool) {
		for i := len(*drawlist) - 1; i >= 0; i-- {
			ei := (*drawlist)[i]
			if ei >= 0 && sys.explods[en][ei].bindtime != 0 &&
				sys.explods[en][ei].bindId == c.id {
				sys.explods[en][ei].id = IErr
				if drop {
					*drawlist = append((*drawlist)[:i], (*drawlist)[i+1:]...)
				} else {
					(*drawlist)[i] = -1
				}
			}
		}
	}
	remove(&sys.explDrawlist[en], true)
	remove(&sys.topexplDrawlist[en], false)
}
func (c *Char) getAnim(n int32, ffx bool) (a *Animation) {
	if n < 0 {
		return nil
	}
	if ffx {
		a = sys.lifebar.fat.get(n)
	} else {
		a = c.gi().anim.get(n)
	}
	if a == nil && !sys.ignoreMostErrors {
		str := "存在しないアニメ: "
		if ffx {
			str += "F:"
		} else {
			str += fmt.Sprintf("P%v:", c.playerNo+1)
		}
		sys.errLog.Printf("%v%v\n", str, n)
	}
	return
}
func (c *Char) setPosX(x float32) {
	if c.pos[0] != x {
		c.pos[0] = x
		c.enemyNearClear()
		if c.player {
			for i := ^c.playerNo & 1; i < len(sys.chars); i += 2 {
				for j := range sys.chars[i] {
					sys.chars[i][j].enemyNearClear()
				}
			}
		}
	}
}
func (c *Char) setPosY(y float32) {
	c.pos[1] = y
}
func (c *Char) posReset() {
	c.facing = 1 - 2*float32(c.playerNo&1)
	c.setX((float32(sys.stage.p[c.playerNo&1].startx-sys.cam.startx)*
		sys.stage.localscl - c.facing*float32(c.playerNo>>1)*P1P3Dist) / c.localscl)
	c.setY(0)
	c.setXV(0)
	c.setYV(0)
}
func (c *Char) setX(x float32) {
	c.oldPos[0], c.drawPos[0] = x, x
	c.setPosX(x)
}
func (c *Char) setY(y float32) {
	c.oldPos[1], c.drawPos[1] = y, y
	c.setPosY(y)
}
func (c *Char) addX(x float32) {
	c.setX(c.pos[0] + c.facing*x)
}
func (c *Char) addY(y float32) {
	c.setY(c.pos[1] + y)
}
func (c *Char) addXV(xv float32) {
	c.vel[0] += xv
}
func (c *Char) addYV(yv float32) {
	c.vel[1] += yv
}
func (c *Char) mulXV(xv float32) {
	c.vel[0] *= xv
}
func (c *Char) mulYV(yv float32) {
	c.vel[1] *= yv
}
func (c *Char) hitAdd(h int32) {
	c.hitCount += h
	for _, tid := range c.targets {
		if t := sys.playerID(tid); t != nil {
			t.getcombo += h
		}
	}
}
func (c *Char) newProj() *Projectile {
	for i, p := range sys.projs[c.playerNo] {
		if p.id < 0 {
			sys.projs[c.playerNo][i].clear()
			sys.projs[c.playerNo][i].id = 0
			sys.projs[c.playerNo][i].palfx = c.getPalfx()
			return &sys.projs[c.playerNo][i]
		}
	}
	if i := len(sys.projs[c.playerNo]); i < sys.playerProjectileMax {
		sys.projs[c.playerNo] = append(sys.projs[c.playerNo], *newProjectile())
		p := &sys.projs[c.playerNo][i]
		p.id, p.palfx = 0, c.getPalfx()
		return p
	}
	return nil
}
func (c *Char) projInit(p *Projectile, pt PosType, x, y float32,
	op bool, rpg, rpn int32) {
	p.setPos(c.helperPos(pt, [...]float32{x, y}, 1, &p.facing, p.localscl))
	if p.anim < -1 {
		p.anim = 0
	}
	p.ani = c.getAnim(p.anim, false)
	if p.ani == nil && c.anim != nil {
		p.ani = &Animation{}
		*p.ani = *c.anim
		p.ani.SetAnimElem(1)
		p.anim = c.animNo
	}
	if p.ani != nil {
		p.ani.UpdateSprite()
	}
	if c.size.proj.doscale != 0 {
		p.scale[0] *= c.size.xscale
		p.scale[1] *= c.size.yscale
	}
	p.clsnScale = c.clsnScale
	if p.velocity[0] < 0 {
		p.facing *= -1
		p.velocity[0] *= -1
		p.accel[0] *= -1
	}
	if op {
		remap := make([]int, len(p.palfx.remap))
		copy(remap, p.palfx.remap)
		p.palfx = newPalFX()
		p.palfx.remap = remap
		c.remapPal(p.palfx, [...]int32{1, 1}, [...]int32{rpg, rpn})
	}
}
func (c *Char) setHitdefDefault(hd *HitDef, proj bool) {
	if !proj {
		c.targetsOfHitdef = c.targetsOfHitdef[:0]
	}
	if hd.attr&^int32(ST_MASK) == 0 {
		hd.attr = 0
	}
	if hd.hitonce < 0 || hd.attr&int32(AT_AT) != 0 {
		hd.hitonce = 1
	}
	ifnanset := func(dst *float32, src float32) {
		if math.IsNaN(float64(*dst)) {
			*dst = src
		}
	}
	ifierrset := func(dst *int32, src int32) {
		if *dst == IErr {
			*dst = src
		}
	}
	ifnanset(&hd.ground_velocity[0], 0)
	ifnanset(&hd.ground_velocity[1], 0)
	ifnanset(&hd.air_velocity[0], 0)
	ifnanset(&hd.air_velocity[1], 0)
	ifnanset(&hd.guard_velocity, hd.ground_velocity[0])
	ifnanset(&hd.airguard_velocity[0], hd.air_velocity[0]*1.5)
	ifnanset(&hd.airguard_velocity[1], hd.air_velocity[1]*0.5)
	ifnanset(&hd.down_velocity[0], hd.air_velocity[0])
	ifnanset(&hd.down_velocity[1], hd.air_velocity[1])
	if hd.fall.animtype == RA_Unknown {
		if hd.air_animtype != RA_Unknown {
			hd.fall.animtype = hd.air_animtype
		} else if hd.animtype < RA_Back {
			hd.fall.animtype = RA_Back
		} else {
			hd.fall.animtype = hd.animtype
		}
	}
	if hd.air_animtype == RA_Unknown {
		hd.air_animtype = hd.animtype
	}
	if hd.air_type == HT_Unknown {
		if hd.ground_type == HT_Trip {
			hd.air_type = HT_High
		} else {
			hd.air_type = hd.ground_type
		}
	}
	ifierrset(&hd.forcestand, Btoi(hd.ground_velocity[1] != 0))
	if hd.attr&int32(ST_A) != 0 {
		ifnanset(&hd.ground_cornerpush_veloff, 0)
	} else {
		ifnanset(&hd.ground_cornerpush_veloff, hd.guard_velocity*1.3)
	}
	ifnanset(&hd.air_cornerpush_veloff, hd.ground_cornerpush_veloff)
	ifnanset(&hd.down_cornerpush_veloff, hd.ground_cornerpush_veloff)
	ifnanset(&hd.guard_cornerpush_veloff, hd.ground_cornerpush_veloff)
	ifnanset(&hd.airguard_cornerpush_veloff, hd.ground_cornerpush_veloff)
	ifierrset(&hd.hitgetpower,
		int32(sys.attack_LifeToPowerMul*float32(hd.hitdamage)))
	ifierrset(&hd.guardgetpower,
		int32(sys.attack_LifeToPowerMul*float32(hd.hitdamage)*0.5))
	ifierrset(&hd.hitgivepower,
		int32(sys.getHit_LifeToPowerMul*float32(hd.hitdamage)))
	ifierrset(&hd.guardgivepower,
		int32(sys.getHit_LifeToPowerMul*float32(hd.hitdamage)*0.5))
	if !math.IsNaN(float64(hd.snap[0])) {
		hd.maxdist[0], hd.mindist[0] = hd.snap[0], hd.snap[0]
	}
	if !math.IsNaN(float64(hd.snap[1])) {
		hd.maxdist[1], hd.mindist[1] = hd.snap[1], hd.snap[1]
	}
	hd.playerNo = c.ss.sb.playerNo
}
func (c *Char) setFEdge(fe float32) {
	c.edge[0] = fe
	c.setSF(CSF_frontedge)
}
func (c *Char) setBEdge(be float32) {
	c.edge[1] = be
	c.setSF(CSF_backedge)
}
func (c *Char) setFWidth(fw float32) {
	c.width[0] = c.defFW()*(320/float32(c.localcoord))/c.localscl + fw
	c.setSF(CSF_frontwidth)
}
func (c *Char) setBWidth(bw float32) {
	c.width[1] = c.defBW()*(320/float32(c.localcoord))/c.localscl + bw
	c.setSF(CSF_backwidth)
}
func (c *Char) gethitAnimtype() Reaction {
	if c.ghv.fallf {
		return c.ghv.fall.animtype
	} else if c.ss.stateType == ST_A {
		return c.ghv.airanimtype
	}
	return c.ghv.groundanimtype
}
func (c *Char) isBound() bool {
	return c.ghv.idMatch(c.bindToId)
}
func (c *Char) varGet(i int32) BytecodeValue {
	if i >= 0 && i < int32(NumVar) {
		return BytecodeInt(c.ivar[i])
	}
	return BytecodeSF()
}
func (c *Char) fvarGet(i int32) BytecodeValue {
	if i >= 0 && i < int32(NumFvar) {
		return BytecodeFloat(c.fvar[i])
	}
	return BytecodeSF()
}
func (c *Char) sysVarGet(i int32) BytecodeValue {
	if i >= 0 && i < int32(NumSysVar) {
		return BytecodeInt(c.ivar[i+int32(NumVar)])
	}
	return BytecodeSF()
}
func (c *Char) sysFvarGet(i int32) BytecodeValue {
	if i >= 0 && i < int32(NumSysFvar) {
		return BytecodeFloat(c.fvar[i+int32(NumFvar)])
	}
	return BytecodeSF()
}
func (c *Char) varSet(i, v int32) BytecodeValue {
	if i >= 0 && i < int32(NumVar) {
		c.ivar[i] = v
		return BytecodeInt(v)
	}
	return BytecodeSF()
}
func (c *Char) fvarSet(i int32, v float32) BytecodeValue {
	if i >= 0 && i < int32(NumFvar) {
		c.fvar[i] = v
		return BytecodeFloat(v)
	}
	return BytecodeSF()
}
func (c *Char) sysVarSet(i, v int32) BytecodeValue {
	if i >= 0 && i < int32(NumSysVar) {
		c.ivar[i+int32(NumVar)] = v
		return BytecodeInt(v)
	}
	return BytecodeSF()
}
func (c *Char) sysFvarSet(i int32, v float32) BytecodeValue {
	if i >= 0 && i < int32(NumSysFvar) {
		c.fvar[i+int32(NumFvar)] = v
		return BytecodeFloat(v)
	}
	return BytecodeSF()
}
func (c *Char) varAdd(i, v int32) BytecodeValue {
	if i >= 0 && i < int32(NumVar) {
		c.ivar[i] += v
		return BytecodeInt(c.ivar[i])
	}
	return BytecodeSF()
}
func (c *Char) fvarAdd(i int32, v float32) BytecodeValue {
	if i >= 0 && i < int32(NumFvar) {
		c.fvar[i] += v
		return BytecodeFloat(c.fvar[i])
	}
	return BytecodeSF()
}
func (c *Char) sysVarAdd(i, v int32) BytecodeValue {
	if i >= 0 && i < int32(NumSysVar) {
		c.ivar[i+int32(NumVar)] += v
		return BytecodeInt(c.ivar[i+int32(NumVar)])
	}
	return BytecodeSF()
}
func (c *Char) sysFvarAdd(i int32, v float32) BytecodeValue {
	if i >= 0 && i < int32(NumSysFvar) {
		c.fvar[i+int32(NumFvar)] += v
		return BytecodeFloat(c.fvar[i+int32(NumFvar)])
	}
	return BytecodeSF()
}
func (c *Char) varRangeSet(s, e, v int32) {
	if s >= 0 {
		for i := s; i <= e && i < int32(NumVar); i++ {
			c.ivar[i] = v
		}
	}
}
func (c *Char) fvarRangeSet(s, e int32, v float32) {
	if s >= 0 {
		for i := s; i <= e && i < int32(NumFvar); i++ {
			c.fvar[i] = v
		}
	}
}
func (c *Char) sysVarRangeSet(s, e, v int32) {
	if s >= 0 {
		for i := s; i <= e && i < int32(NumSysVar); i++ {
			c.ivar[i+int32(NumVar)] = v
		}
	}
}
func (c *Char) sysFvarRangeSet(s, e int32, v float32) {
	if s >= 0 {
		for i := s; i <= e && i < int32(NumSysFvar); i++ {
			c.fvar[i+int32(NumFvar)] = v
		}
	}
}
func (c *Char) setFacing(f float32) {
	if f != 0 {
		if (c.facing < 0) != (f < 0) {
			c.facing *= -1
			c.vel[0] *= -1
			c.ghv.xvel *= -1
		}
	}
}
func (c *Char) getTarget(id int32) []int32 {
	if id < 0 {
		return c.targets
	}
	var tg []int32
	for _, tid := range c.targets {
		if t := sys.playerID(tid); t != nil {
			if t.ghv.hitid == id {
				tg = append(tg, tid)
			}
		}
	}
	return tg
}
func (c *Char) targetFacing(tar []int32, f int32) {
	tf := c.facing
	if f < 0 {
		tf *= -1
	}
	for _, tid := range tar {
		if t := sys.playerID(tid); t != nil {
			t.setFacing(tf)
		}
	}
}
func (c *Char) targetBind(tar []int32, time int32, x, y float32) {
	for _, tid := range tar {
		if t := sys.playerID(tid); t != nil {
			t.setBindToId(c)
			t.setBindTime(time)
			t.bindFacing = 0
			x *= c.localscl / t.localscl
			y *= c.localscl / t.localscl
			t.bindPos = [...]float32{x, y}
		}
	}
}
func (c *Char) bindToTarget(tar []int32, time int32, x, y float32, hmf HMF) {
	if len(tar) > 0 {
		if t := sys.playerID(tar[0]); t != nil {
			switch hmf {
			case HMF_M:
				x += float32(t.size.mid.pos[0])
				y += float32(t.size.mid.pos[1])
			case HMF_H:
				x += float32(t.size.head.pos[0])
				y += float32(t.size.head.pos[1])
			}
			if !math.IsNaN(float64(x)) {
				c.setX(t.pos[0] + t.facing*x)
			}
			if !math.IsNaN(float64(y)) {
				c.setY(t.pos[1] + y)
			}
			c.targetBind(tar[:1], time, c.facing*c.distX(t, c), (t.pos[1]*t.localscl/c.localscl)-(c.pos[1]*c.localscl/t.localscl))
		}
	}
}
func (c *Char) targetLifeAdd(tar []int32, add int32, kill, absolute bool) {
	for _, tid := range tar {
		if t := sys.playerID(tid); t != nil {
			t.lifeAdd(-float64(t.computeDamage(-float64(add), kill, absolute, 1)),
				true, true)
		}
	}
}
func (c *Char) targetState(tar []int32, state int32) {
	if state >= 0 {
		pn := c.ss.sb.playerNo
		if c.minus == -2 || c.minus == -20 {
			pn = c.playerNo
		}
		for _, tid := range tar {
			if t := sys.playerID(tid); t != nil {
				t.setCtrl(false)
				t.stateChange1(state, pn)
			}
		}
	}
}
func (c *Char) targetVelSetX(tar []int32, x float32) {
	for _, tid := range tar {
		if t := sys.playerID(tid); t != nil {
			x *= c.localscl / t.localscl
			t.setXV(x)
		}
	}
}
func (c *Char) targetVelSetY(tar []int32, y float32) {
	for _, tid := range tar {
		if t := sys.playerID(tid); t != nil {
			y *= c.localscl / t.localscl
			t.setYV(y)
		}
	}
}
func (c *Char) targetVelAddX(tar []int32, x float32) {
	for _, tid := range tar {
		if t := sys.playerID(tid); t != nil {
			x *= c.localscl / t.localscl
			t.vel[0] += x
		}
	}
}
func (c *Char) targetVelAddY(tar []int32, y float32) {
	for _, tid := range tar {
		if t := sys.playerID(tid); t != nil {
			y *= c.localscl / t.localscl
			t.vel[1] += y
		}
	}
}
func (c *Char) targetPowerAdd(tar []int32, power int32) {
	for _, tid := range tar {
		if t := sys.playerID(tid); t != nil && t.player {
			t.powerAdd(power)
		}
	}
}
func (c *Char) targetDrop(excludeid int32, keepone bool) {
	var tg []int32
	if excludeid < 0 {
		tg = c.targets
	} else {
		for _, tid := range c.targets {
			if t := sys.playerID(tid); t != nil {
				if t.ghv.hitid == excludeid {
					tg = append(tg, tid)
				} else {
					t.gethitBindClear()
					t.ghv.dropId(c.id)
				}
			}
		}
	}
	if (keepone || excludeid < 0) && len(tg) > 0 {
		c.targets = nil
		r := -1
		if keepone && excludeid >= 0 {
			r = int(Rand(0, int32(len(tg))-1))
		}
		for i, tid := range tg {
			if i == r {
				c.targets = append(c.targets, tid)
			} else if t := sys.playerID(tid); t != nil {
				t.gethitBindClear()
				t.ghv.dropId(c.id)
			}
		}
	} else {
		c.targets = tg
	}
}
func (c *Char) computeDamage(damage float64, kill, absolute bool,
	atkmul float32) int32 {
	if damage == 0 || !absolute && atkmul == 0 {
		return 0
	}
	if !absolute {
		damage *= float64(atkmul / c.defenceMul)
	}
	damage = math.Ceil(damage)
	min, max := float64(c.life-c.lifeMax), float64(Max(0, c.life-Btoi(!kill)))
	if damage < min {
		damage = min
	}
	if damage > max {
		damage = max
	}
	return int32(damage)
}
func (c *Char) lifeAdd(add float64, kill, absolute bool) {
	if add != 0 && c.roundState() != 3 {
		if !absolute {
			add /= float64(c.defenceMul)
		}
		add = math.Floor(add)
		max := float64(c.lifeMax - c.life)
		if add > max {
			add = max
		}
		min := float64(Btoi(!kill && c.life > 0) - c.life)
		if add < min {
			add = min
		}
		c.lifeSet(c.life + int32(add))
	}
}
func (c *Char) lifeSet(life int32) {
	if c.life = Max(0, Min(c.lifeMax, life)); c.life == 0 {
		if c.player {
			if c.alive() && c.helperIndex == 0 {
				if c.ss.moveType != MT_H {
					if c.playerNo == c.ss.sb.playerNo {
						sys.winType[^c.playerNo&1] = WT_Suicide
					} else if c.playerNo&1 == c.ss.sb.playerNo&1 {
						sys.winType[^c.playerNo&1] = WT_Teammate
					}
				} else if c.playerNo == c.ghv.playerNo {
					sys.winType[^c.playerNo&1] = WT_Suicide
				} else if c.ghv.playerNo >= 0 && c.playerNo&1 == c.ghv.playerNo&1 {
					sys.winType[^c.playerNo&1] = WT_Teammate
				} else if c.ghv.guarded {
					sys.winType[^c.playerNo&1] = WT_C
				} else if c.ghv.attr&int32(AT_AH) != 0 {
					sys.winType[^c.playerNo&1] = WT_H
				} else if c.ghv.attr&int32(AT_AS) != 0 {
					sys.winType[^c.playerNo&1] = WT_S
				} else if c.ghv.attr&int32(AT_AT) != 0 {
					sys.winType[^c.playerNo&1] = WT_Throw
				} else {
					sys.winType[^c.playerNo&1] = WT_N
				}
			}
		} else {
			c.life = 1
		}
	}
}
func (c *Char) setPower(pow int32) {
	if !sys.roundEnd() {
		c.power = Max(0, Min(c.powerMax, pow))
	}
}
func (c *Char) powerAdd(add int32) {
	if sys.powerShare[c.playerNo&1] {
		sys.chars[c.playerNo&1][0].setPower(c.getPower() + add)
	} else {
		sys.chars[c.playerNo][0].setPower(c.getPower() + add)
	}
}
func (c *Char) powerSet(pow int32) {
	if sys.powerShare[c.playerNo&1] {
		sys.chars[c.playerNo&1][0].setPower(pow)
	} else {
		sys.chars[c.playerNo][0].setPower(pow)
	}
}
func (c *Char) distX(opp *Char, oc *Char) float32 {
	return (opp.pos[0]*opp.localscl - c.pos[0]*c.localscl) / oc.localscl
}
func (c *Char) bodyDistX(opp *Char, oc *Char) float32 {
	dist := c.distX(opp, oc)
	var oppw float32
	if dist == 0 || (dist < 0) != (opp.facing < 0) {
		oppw = opp.facing * opp.width[0] * (320 / float32(opp.localcoord)) / oc.localscl
	} else {
		oppw = -opp.facing * opp.width[1] * (320 / float32(opp.localcoord)) / oc.localscl
	}
	return dist + oppw - c.facing*c.width[0]
}
func (c *Char) rdDistX(rd *Char, oc *Char) BytecodeValue {
	if rd == nil {
		return BytecodeSF()
	}
	dist := c.facing * c.distX(rd, oc)
	if c.stCgi().ver[0] != 1 {
		dist = float32(int32(dist)) //旧バージョンでは小数点切り捨て
	}
	return BytecodeFloat(dist)
}
func (c *Char) rdDistY(rd *Char, oc *Char) BytecodeValue {
	if rd == nil {
		return BytecodeSF()
	}
	dist := (rd.pos[1]*rd.localscl - c.pos[1]*c.localscl) / oc.localscl
	return BytecodeFloat(dist)
}
func (c *Char) p2BodyDistX(oc *Char) BytecodeValue {
	if p2 := c.p2(); p2 == nil {
		return BytecodeSF()
	} else {
		dist := c.facing * c.bodyDistX(p2, oc)
		if c.stCgi().ver[0] != 1 {
			dist = float32(int32(dist)) //旧バージョンでは小数点切り捨て
		}
		return BytecodeFloat(dist)
	}
}
func (c *Char) hitVelSetX() {
	if c.ss.moveType == MT_H {
		c.setXV(c.ghv.xvel)
	}
}
func (c *Char) hitVelSetY() {
	if c.ss.moveType == MT_H {
		c.setYV(c.ghv.yvel)
	}
}
func (c *Char) getEdge(base float32, actually bool) float32 {
	if !actually || c.stCgi().ver[0] != 1 {
		switch c.ss.stateType {
		case ST_A:
			return base + 1
		case ST_L:
			return base + 2
		}
	}
	return base
}
func (c *Char) defFW() float32 {
	if c.ss.stateType == ST_A {
		return float32(c.size.air.front)
	}
	return float32(c.size.ground.front)
}
func (c *Char) defBW() float32 {
	if c.ss.stateType == ST_A {
		return float32(c.size.air.back)
	}
	return float32(c.size.ground.back)
}
func (c *Char) height() float32 {
	return float32(c.size.height)
}
func (c *Char) setPauseTime(pausetime, movetime int32) {
	if ^pausetime < sys.pausetime || c.playerNo != c.ss.sb.playerNo ||
		sys.pauseplayer == c.playerNo {
		sys.pausetime = ^pausetime
		sys.pauseplayer = c.playerNo
		if sys.pauseendcmdbuftime < 0 || sys.pauseendcmdbuftime > pausetime {
			sys.pauseendcmdbuftime = 0
		}
	}
	c.pauseMovetime = Max(0, movetime)
	if c.pauseMovetime > pausetime {
		c.pauseMovetime = 0
	} else if sys.pause > 0 && c.pauseMovetime > 0 {
		c.pauseMovetime--
	}
}
func (c *Char) setSuperPauseTime(pausetime, movetime int32, unhittable bool) {
	if ^pausetime < sys.supertime || c.playerNo != c.ss.sb.playerNo ||
		sys.superplayer == c.playerNo {
		sys.supertime = ^pausetime
		sys.superplayer = c.playerNo
		if sys.superendcmdbuftime < 0 || sys.superendcmdbuftime > pausetime {
			sys.superendcmdbuftime = 0
		}
	}
	c.superMovetime = Max(0, movetime)
	if c.superMovetime > pausetime {
		c.superMovetime = 0
	} else if sys.super > 0 && c.superMovetime > 0 {
		c.superMovetime--
	}
	if unhittable {
		c.gi().unhittable = pausetime + Btoi(pausetime > 0)
	}
}
func (c *Char) getPalfx() *PalFX {
	if c.palfx != nil {
		return c.palfx
	}
	if c.parentIndex >= 0 {
		if p := c.parent(); p != nil {
			return p.getPalfx()
		}
	}
	c.palfx = newPalFX()
	return c.palfx
}
func (c *Char) getPalMap() []int {
	return c.getPalfx().remap
}
func (c *Char) pause() bool {
	return c.acttmp <= -2
}
func (c *Char) hitPause() bool {
	return c.hitPauseTime > 0
}
func (c *Char) angleSet(a float32) {
	c.angle = a
}
func (c *Char) ctrlOver() bool {
	return sys.time == 0 ||
		sys.intro < -(sys.lifebar.ro.over_hittime+sys.lifebar.ro.over_waittime)
}
func (c *Char) over() bool {
	return c.scf(SCF_over) || (c.ctrlOver() && c.scf(SCF_ctrl) &&
		c.ss.stateType != ST_A && c.ss.physics != ST_A)
}
func (c *Char) makeDust(x, y float32) {
	if e, i := c.newExplod(); e != nil {
		e.anim = c.getAnim(120, true)
		e.sprpriority = math.MaxInt32
		e.ownpal = true
		e.offset = [...]float32{x, y}
		e.setPos(c)
		c.insertExplod(i)
	}
}
func (c *Char) hitFallDamage() {
	if c.ss.moveType == MT_H {
		c.lifeAdd(-float64(c.ghv.fall.damage), c.ghv.fall.kill, false)
	}
}
func (c *Char) hitFallVel() {
	if c.ss.moveType == MT_H {
		if !math.IsNaN(float64(c.ghv.fall.xvelocity)) {
			c.setXV(c.ghv.fall.xvelocity)
		}
		c.setYV(c.ghv.fall.yvelocity)
	}
}
func (c *Char) hitFallSet(f int32, xv, yv float32) {
	if c.ss.moveType == MT_H {
		if f >= 0 {
			c.ghv.fallf = f != 0
		}
		if !math.IsNaN(float64(xv)) {
			c.ghv.fall.xvelocity = xv
		}
		if !math.IsNaN(float64(yv)) {
			c.ghv.fall.yvelocity = yv
		}
	}
}
func (c *Char) remapPal(pfx *PalFX, src [2]int32, dst [2]int32) {
	if src[0] < 0 || src[1] < 0 || dst[0] < 0 || dst[1] < 0 {
		return
	}
	si, ok := c.gi().sff.palList.PalTable[[...]int16{int16(src[0]),
		int16(src[1])}]
	if !ok {
		return
	}
	var di int
	di, ok = c.gi().sff.palList.PalTable[[...]int16{int16(dst[0]),
		int16(dst[1])}]
	if !ok {
		di = si
	}
	if pfx.remap == nil {
		pfx.remap = c.gi().sff.palList.GetPalMap()
	}
	if c.gi().sff.palList.SwapPalMap(&pfx.remap) {
		c.gi().sff.palList.Remap(si, di)
		if src[0] == 1 && src[1] == 1 && c.gi().sff.header.Ver0 == 1 {
			spr := c.gi().sff.GetSprite(0, 0)
			if spr != nil {
				c.gi().sff.palList.Remap(spr.palidx, di)
			}
			spr = c.gi().sff.GetSprite(9000, 0)
			if spr != nil {
				c.gi().sff.palList.Remap(spr.palidx, di)
			}
		}
		c.gi().sff.palList.SwapPalMap(&pfx.remap)
	}
}
func (c *Char) inGuardState() bool {
	return c.ss.no == 120 || (c.ss.no >= 130 && c.ss.no <= 132) ||
		c.ss.no == 140 || (c.ss.no >= 150 && c.ss.no <= 155)
}
func (c *Char) gravity() {
	c.vel[1] += c.gi().movement.yaccel
}
func (c *Char) posUpdate() {
	nobind := [...]bool{c.bindTime == 0 || math.IsNaN(float64(c.bindPos[0])),
		c.bindTime == 0 || math.IsNaN(float64(c.bindPos[1]))}
	for i := range nobind {
		if nobind[i] {
			c.oldPos[i], c.drawPos[i] = c.pos[i], c.pos[i]
		}
	}
	if c.sf(CSF_posfreeze) {
		if nobind[0] {
			c.setPosX(c.oldPos[0] + c.veloff)
		}
	} else {
		if nobind[0] {
			c.setPosX(c.oldPos[0] + c.vel[0]*c.facing + c.veloff)
		}
		if nobind[1] {
			c.setPosY(c.oldPos[1] + c.vel[1])
		}
		switch c.ss.physics {
		case ST_S:
			c.vel[0] *= c.gi().movement.stand.friction
			if AbsF(c.vel[0]) < 1 {
				c.vel[0] = 0
			}
		case ST_C:
			c.vel[0] *= c.gi().movement.crouch.friction
		case ST_A:
			c.gravity()
		}
	}
	c.veloff *= 0.7
	if AbsF(c.veloff) < 1 {
		c.veloff = 0
	}
	c.oldLocalscl = c.localscl
}
func (c *Char) addTarget(id int32) {
	if !c.hasTarget(id) {
		c.targets = append(c.targets, id)
	}
}
func (c *Char) hasTarget(id int32) bool {
	for _, tid := range c.targets {
		if tid == id {
			return true
		}
	}
	return false
}
func (c *Char) hasTargetOfHitdef(id int32) bool {
	for _, tid := range c.targetsOfHitdef {
		if tid == id {
			return true
		}
	}
	return false
}
func (c *Char) setBindTime(time int32) {
	c.bindTime = time
	if time == 0 {
		c.bindToId = -1
	}
}
func (c *Char) setBindToId(to *Char) {
	c.bindToId = to.id
	if to.bindToId == c.id {
		to.setBindTime(0)
	}
	c.bindFacing = to.facing * 2
}
func (c *Char) bind() {
	if c.bindTime == 0 {
		return
	}
	if bt := sys.playerID(c.bindToId); bt != nil {
		if bt.hasTarget(c.id) {
			if bt.sf(CSF_destroy) {
				c.selfState(5050, -1, -1)
				c.setBindTime(0)
				return
			}
			if !math.IsNaN(float64(c.bindPos[0])) {
				c.setXV(c.facing * bt.facing * bt.vel[0])
			}
			if !math.IsNaN(float64(c.bindPos[1])) {
				c.setYV(bt.vel[1])
			}
		}
		if !math.IsNaN(float64(c.bindPos[0])) {
			f := bt.facing
			if AbsF(c.bindFacing) == 2 {
				f = c.bindFacing / 2
			}
			c.setX(bt.pos[0] + f*c.bindPos[0])
			c.drawPos[0] += bt.drawPos[0] - bt.pos[0]
			c.oldPos[0] += bt.oldPos[0] - bt.pos[0]
			c.pushed = c.pushed || bt.pushed
			c.ghv.xoff = 0
		}
		if !math.IsNaN(float64(c.bindPos[1])) {
			c.setY(bt.pos[1] + c.bindPos[1])
			c.drawPos[1] += bt.drawPos[1] - bt.pos[1]
			c.oldPos[1] += bt.oldPos[1] - bt.pos[1]
			c.ghv.yoff = 0
		}
		if AbsF(c.bindFacing) == 1 {
			if c.bindFacing > 0 {
				c.setFacing(bt.facing)
			} else {
				c.setFacing(-bt.facing)
			}
		}
	} else {
		c.setBindTime(0)
		return
	}
}
func (c *Char) xScreenBound() {
	x := c.pos[0]
	if c.sf(CSF_screenbound) {
		min, max := c.getEdge(c.edge[0], true), -c.getEdge(c.edge[1], true)
		if c.facing > 0 {
			min, max = -max, -min
		}
		x = MaxF(min+sys.xmin/c.localscl, MinF(max+sys.xmax/c.localscl, x))
	}
	x = MaxF(sys.stage.leftbound/c.localscl, MinF(sys.stage.rightbound/c.localscl, x))
	c.setPosX(x)
}
func (c *Char) gethitBindClear() {
	if c.isBound() {
		c.setBindTime(0)
	}
}
func (c *Char) removeTarget(pid int32) {
	for i, t := range c.targets {
		if t == pid {
			c.targets = append(c.targets[:i], c.targets[i+1:]...)
			break
		}
	}
}
func (c *Char) exitTarget(explremove bool) {
	if c.hittmp >= 0 {
		for _, hb := range c.ghv.hitBy {
			if e := sys.playerID(hb[0]); e != nil {
				e.removeTarget(c.id)
				if explremove {
					c.enemyExplodsRemove(e.playerNo)
				}
			}
		}
		c.gethitBindClear()
	}
	c.ghv.hitBy = c.ghv.hitBy[:0]
}
func (c *Char) offsetX() float32 {
	return float32(c.size.draw.offset[0])*c.facing + c.offset[0]
}
func (c *Char) offsetY() float32 {
	return float32(c.size.draw.offset[1]) + c.offset[1]
}
func (c *Char) projClsnCheck(p *Projectile, gethit bool) bool {
	if p.ani == nil || c.curFrame == nil {
		return false
	}
	frm := p.ani.CurrentFrame()
	if frm == nil {
		return false
	}
	var clsn1, clsn2 []float32
	if gethit {
		clsn1, clsn2 = frm.Clsn1(), c.curFrame.Clsn2()
	} else {
		clsn1, clsn2 = frm.Clsn2(), c.curFrame.Clsn1()
	}
	return sys.clsnHantei(clsn1, [...]float32{p.clsnScale[0] * p.localscl, p.clsnScale[1] * p.localscl},
		[...]float32{p.pos[0] * p.localscl, p.pos[1] * p.localscl}, p.facing,
		clsn2, [...]float32{c.clsnScale[0] * (320 / float32(sys.chars[c.animPN][0].localcoord)), c.clsnScale[1] * (320 / float32(sys.chars[c.animPN][0].localcoord))},
		[...]float32{c.pos[0]*c.localscl + c.offsetX()*c.localscl,
			c.pos[1]*c.localscl + c.offsetY()*c.localscl}, c.facing)
}
func (c *Char) clsnCheck(atk *Char, c1atk, c1slf bool) bool {
	if atk.curFrame == nil || c.curFrame == nil {
		return false
	}
	var clsn1, clsn2 []float32
	if c1atk {
		clsn1 = atk.curFrame.Clsn1()
	} else {
		clsn1 = atk.curFrame.Clsn2()
	}
	if c1slf {
		clsn2 = c.curFrame.Clsn1()
	} else {
		clsn2 = c.curFrame.Clsn2()
	}
	return sys.clsnHantei(clsn1, [...]float32{sys.chars[atk.animPN][0].clsnScale[0] * (320 / float32(sys.chars[atk.animPN][0].localcoord)), sys.chars[atk.animPN][0].clsnScale[1] * (320 / float32(sys.chars[atk.animPN][0].localcoord))},
		[...]float32{atk.pos[0]*atk.localscl + atk.offsetX()*atk.localscl,
			atk.pos[1]*atk.localscl + atk.offsetY()*atk.localscl},
		atk.facing, clsn2, [...]float32{sys.chars[c.animPN][0].clsnScale[0] * (320 / float32(sys.chars[c.animPN][0].localcoord)), sys.chars[c.animPN][0].clsnScale[1] * (320 / float32(sys.chars[c.animPN][0].localcoord))},
		[...]float32{c.pos[0]*c.localscl + c.offsetX()*c.localscl,
			c.pos[1]*c.localscl + c.offsetY()*c.localscl}, c.facing)
}
func (c *Char) hitCheck(e *Char) bool {
	return c.clsnCheck(e, true, e.hitdef.reversal_attr > 0)
}
func (c *Char) attrCheck(h *HitDef, pid int32, st StateType) bool {
	if c.gi().unhittable > 0 || h.chainid >= 0 && c.ghv.hitid != h.chainid {
		return false
	}
	if len(c.ghv.hitBy) > 0 && c.ghv.hitBy[len(c.ghv.hitBy)-1][0] == pid {
		for _, nci := range h.nochainid {
			if nci >= 0 && c.ghv.hitid == nci {
				return false
			}
		}
	}
	if h.reversal_attr > 0 {
		return c.atktmp != 0 && c.hitdef.attr > 0 &&
			(c.hitdef.attr&h.reversal_attr&int32(ST_MASK)) != 0 &&
			(c.hitdef.attr&h.reversal_attr&^int32(ST_MASK)) != 0
	}
	if h.attr <= 0 || h.hitflag&int32(c.ss.stateType) == 0 ||
		h.hitflag&int32(ST_F) == 0 && c.hittmp >= 2 ||
		h.hitflag&int32(MT_MNS) != 0 && c.hittmp > 0 ||
		h.hitflag&int32(MT_PLS) != 0 && c.hittmp <= 0 {
		return false
	}
	if h.chainid < 0 {
		var styp int32
		if st == ST_N {
			styp = h.attr & int32(ST_MASK)
		} else {
			styp = int32(st)
		}
		for _, hb := range c.hitby {
			if hb.time != 0 &&
				(hb.flag&styp == 0 || hb.flag&h.attr&^int32(ST_MASK) == 0) {
				return false
			}
		}
	}
	return true
}
func (c *Char) hittable(h *HitDef, e *Char, st StateType,
	countercheck func(*HitDef) bool) bool {
	if !c.attrCheck(h, e.id, st) {
		return false
	}
	if c.atktmp != 0 && (c.hitdef.attr > 0 && c.ss.stateType != ST_L ||
		c.hitdef.reversal_attr > 0) {
		switch {
		case c.hitdef.reversal_attr > 0:
			if h.reversal_attr > 0 {
				if countercheck(&c.hitdef) {
					c.atktmp = -1
					return e.atktmp < 0
				}
				return true
			}
		case h.reversal_attr > 0:
			return true
		case h.priority < c.hitdef.priority:
		case h.priority == c.hitdef.priority:
			switch {
			case c.hitdef.bothhittype == AT_Dodge:
			case h.bothhittype != AT_Hit:
			case c.hitdef.bothhittype == AT_Hit:
				if (c.hitdef.p1stateno >= 0 || c.hitdef.attr&int32(AT_AT) != 0 &&
					h.hitonce != 0) && countercheck(&c.hitdef) {
					c.atktmp = -1
					return e.atktmp < 0 || Rand(0, 1) == 1
				}
				return true
			default:
				return true
			}
		default:
			return true
		}
		return !countercheck(&c.hitdef)
	}
	return true
}
func (c *Char) action() {
	if c.minus != 2 || c.sf(CSF_destroy) {
		return
	}
	p := false
	if c.cmd != nil {
		if sys.super > 0 {
			p = c.superMovetime == 0
		} else if sys.pause > 0 && c.pauseMovetime == 0 {
			p = true
		}
	}
	c.acttmp = -int8(Btoi(p)) * 2
	c.unsetSCF(SCF_guard)
	if !(c.scf(SCF_ko) || c.ctrlOver()) && (c.scf(SCF_ctrl) || c.ss.no == 52) &&
		c.ss.moveType == MT_I && c.cmd != nil &&
		(sys.autoguard[c.playerNo] || c.cmd[0].Buffer.B > 0) &&
		(c.ss.stateType == ST_S && !c.sf(CSF_nostandguard) ||
			c.ss.stateType == ST_C && !c.sf(CSF_nocrouchguard) ||
			c.ss.stateType == ST_A && !c.sf(CSF_noairguard)) {
		c.setSCF(SCF_guard)
	}
	if !p {
		if c.palfx != nil {
			c.palfx.step()
		}
		if c.keyctrl && c.cmd != nil {
			if c.ss.stateType == ST_A {
				if c.cmd[0].Buffer.U < 0 {
					c.setSCF(SCF_airjump)
				}
			} else {
				c.airJumpCount = 0
				c.unsetSCF(SCF_airjump)
			}
			if c.ctrl() && (c.key >= 0 || c.helperIndex == 0) {
				if !sys.roundEnd() && c.ss.stateType == ST_S && c.cmd[0].Buffer.U > 0 {
					if c.ss.no != 40 {
						c.changeState(40, -1, -1)
					}
				} else if c.ss.stateType == ST_A && c.scf(SCF_airjump) &&
					c.pos[1] <= float32(c.gi().movement.airjump.height) &&
					c.airJumpCount < c.gi().movement.airjump.num &&
					c.cmd[0].Buffer.U > 0 {
					if c.ss.no != 45 {
						c.airJumpCount++
						c.unsetSCF(SCF_airjump)
						c.changeState(45, -1, -1)
					}
				} else {
					if c.ss.stateType == ST_S && c.cmd[0].Buffer.D > 0 {
						if c.ss.no != 10 {
							c.changeState(10, -1, -1)
						}
					} else if c.ss.stateType == ST_C && c.cmd[0].Buffer.D < 0 {
						if c.ss.no != 12 {
							c.changeState(12, -1, -1)
						}
					} else if !c.sf(CSF_nowalk) && c.ss.stateType == ST_S &&
						(c.cmd[0].Buffer.F > 0 || !(c.inguarddist && c.scf(SCF_guard)) &&
							c.cmd[0].Buffer.B > 0) {
						if c.ss.no != 20 {
							c.changeState(20, -1, -1)
						}
					} else if c.ss.no == 20 &&
						c.cmd[0].Buffer.B < 0 && c.cmd[0].Buffer.F < 0 {
						c.changeState(0, -1, -1)
					}
					if c.inguarddist && c.scf(SCF_guard) && c.cmd[0].Buffer.B > 0 &&
						!c.inGuardState() {
						c.changeState(120, -1, -1)
					}
				}
			} else if c.scf(SCF_ctrl) {
				switch c.ss.no {
				case 11:
					c.changeState(12, -1, -1)
				case 20:
					c.changeState(0, -1, -1)
				}
			}
		}
		if !c.hitPause() {
			if !c.sf(CSF_noautoturn) && c.ss.no == 52 {
				c.furimuki()
			}
			if !sys.roundEnd() {
				if c.alive() && c.life > 0 {
					c.unsetSCF(SCF_over | SCF_ko_round_middle)
				}
				if c.ss.no == 5150 || c.scf(SCF_over) {
					c.setSCF(SCF_ko_round_middle)
				}
			}
			if c.ss.no == 5150 {
				c.setSCF(SCF_over)
			}
			c.specialFlag = 0
			if c.player {
				if c.alive() || !c.scf(SCF_over) || !c.scf(SCF_ko_round_middle) {
					c.setSF(CSF_screenbound | CSF_movecamera_x | CSF_movecamera_y)
					if (c.alive() || !c.scf(SCF_over)) && c.roundState() > 0 {
						c.setSF(CSF_playerpush)
					}
				}
			}
			c.angleScalse = [...]float32{1, 1}
			c.attackDist = float32(c.size.attack.dist)
			c.offset = [2]float32{}
			for i, hb := range c.hitby {
				if hb.time > 0 {
					c.hitby[i].time--
				}
			}
			for i, ho := range c.ho {
				if ho.time > 0 {
					c.ho[i].time--
				}
			}
			if sys.super > 0 {
				if c.superMovetime > 0 {
					c.superMovetime--
				}
			} else if sys.pause > 0 && c.pauseMovetime > 0 {
				c.pauseMovetime--
			}
		}
		c.unsetSF(CSF_noautoturn)
		if c.gi().ver[0] == 1 {
			c.unsetSF(CSF_assertspecial | CSF_angledraw)
			c.angleScalse = [...]float32{1, 1}
			c.offset = [2]float32{}
		}
		c.minus = -30
		if c.ss.sb.playerNo == c.playerNo && c.player {
			if sb, ok := c.gi().states[-30]; ok {
				sb.run(c)
			}
		}
		c.minus = -20
		if c.player {
			if sb, ok := c.gi().states[-20]; ok {
				sb.run(c)
			}
		}
		c.minus = -10
		if c.keyctrl && c.ss.sb.playerNo == c.playerNo {
			if sb, ok := c.gi().states[-10]; ok {
				sb.run(c)
			}
		}
		c.minus = -3
		if c.ss.sb.playerNo == c.playerNo && c.player {
			if sb, ok := c.gi().states[-3]; ok {
				sb.run(c)
			}
		}
		c.minus = -2
		if c.player {
			if sb, ok := c.gi().states[-2]; ok {
				sb.run(c)
			}
		}
		c.minus = -1
		if c.keyctrl && c.ss.sb.playerNo == c.playerNo {
			if sb, ok := c.gi().states[-1]; ok {
				sb.run(c)
			}
		}
		c.stateChange2()
		c.minus = 0
		c.ss.sb.run(c)
		if !c.hitPause() {
			if c.ss.no == 5110 && c.recoverTime <= 0 && c.alive() {
				c.changeState(5120, -1, -1)
			}
			for c.ss.no == 140 && (c.anim == nil || len(c.anim.frames) == 0 ||
				c.ss.time >= c.anim.totaltime) {
				c.changeState(Btoi(c.ss.stateType == ST_C)*11+
					Btoi(c.ss.stateType == ST_A)*51, -1, -1)
			}
			for {
				c.posUpdate()
				if c.ss.physics != ST_A || c.vel[1] <= 0 || c.pos[1] < 0 ||
					c.ss.no == 105 {
					break
				}
				c.changeState(52, -1, -1)
			}
			c.ss.time++
			if c.mctime > 0 {
				c.mctime++
			}
			c.setFacing(c.p1facing)
			c.p1facing = 0
			if c.anim != nil {
				c.curFrame = c.anim.CurrentFrame()
			} else {
				c.curFrame = nil
			}
		}
		if c.ghv.damage != 0 {
			if c.ss.moveType == MT_H {
				c.lifeAdd(-float64(c.ghv.damage), true, true)
			}
			c.ghv.damage = 0
		}
		if c.helperIndex == 0 && c.gi().pctime >= 0 {
			c.gi().pctime++
		}
	}
	c.xScreenBound()
	if !p {
		for _, tid := range c.targets {
			if t := sys.playerID(tid); t != nil && t.bindToId == c.id {
				t.bind()
			}
		}
	}
	c.minus = 1
	c.acttmp += int8(Btoi(!c.pause() && !c.hitPause())) -
		int8(Btoi(c.hitPause()))
	if !c.hitPause() {
		if !c.sf(CSF_frontwidth) {
			c.width[0] = c.defFW() * (320 / float32(c.localcoord)) / c.localscl
		}
		if !c.sf(CSF_backwidth) {
			c.width[1] = c.defBW() * (320 / float32(c.localcoord)) / c.localscl
		}
		if !c.sf(CSF_frontedge) {
			c.edge[0] = 0
		}
		if !c.sf(CSF_backedge) {
			c.edge[1] = 0
		}
	}
}
func (c *Char) update(cvmin, cvmax,
	highest, lowest, leftest, rightest *float32) {
	if sys.tickFrame() {
		if c.sf(CSF_destroy) {
			c.destroy()
			return
		}
		if !c.pause() && !c.isBound() {
			c.bind()
		}
		if c.acttmp > 0 {
			if c.anim != nil {
				c.anim.UpdateSprite()
			}
			if c.ss.moveType == MT_H {
				if c.ghv.xoff != 0 {
					c.setPosX(c.pos[0] + c.ghv.xoff)
					c.ghv.xoff = 0
				}
				if c.ghv.yoff != 0 {
					c.setPosY(c.pos[1] + c.ghv.yoff)
					c.ghv.yoff = 0
				}
			}
		}
		if c.ss.moveType == MT_H {
			if sys.super <= 0 && sys.pause <= 0 {
				c.superMovetime, c.pauseMovetime = 0, 0
			}
			c.hittmp = int8(Btoi(c.ghv.fallf)) + 1
			if c.acttmp > 0 && (c.ss.no == 5100 || c.ss.no == 5070) &&
				c.ss.time == 1 {
				c.defenceMul *= c.gi().data.fall.defence_mul
				c.ghv.fallcount++
			}
		}
		if c.acttmp > 0 && c.ss.moveType != MT_H || c.roundState() == 2 &&
			c.scf(SCF_ko) && c.scf(SCF_over) {
			c.exitTarget(true)
		}
		c.atktmp = int8(Btoi((c.ss.moveType != MT_I ||
			c.hitdef.reversal_attr > 0) && !c.hitPause()))
		c.hoIdx = -1
		if c.acttmp > 0 {
			if c.inGuardState() {
				c.setSCF(SCF_guard)
			}
			if c.ss.moveType == MT_H {
				if c.ghv.guarded {
					c.getcombo = 0
				}
				if c.ghv.hitshaketime > 0 {
					c.ghv.hitshaketime--
				}
				if c.ghv.hitshaketime <= 0 && c.ghv.hittime >= 0 {
					c.ghv.hittime--
				}
				if c.ghv.fallf {
					c.fallTime++
				}
			} else {
				if c.hittmp > 0 {
					c.hittmp = 0
				}
				c.defenceMul = float32(c.gi().data.defence) / 100
				c.ghv.hittime = -1
				c.ghv.hitshaketime = 0
				c.ghv.fallf = false
				c.ghv.fallcount = 0
				c.ghv.hitid = -1
				c.getcombo = 0
			}
			if (c.ss.moveType == MT_H || c.ss.no == 52) && c.pos[1] == 0 &&
				AbsF(c.pos[0]-c.oldPos[0]) >= 1 && c.ss.time%3 == 0 {
				c.makeDust(0, 0)
			}
		}
	}
	if sys.tickNextFrame() {
		c.pushed = false
	}
	if c.acttmp > 0 {
		spd := sys.tickInterpola()
		if c.pushed {
			spd = 0
		}
		if !c.sf(CSF_posfreeze) {
			for i := 0; i < 2; i++ {
				c.drawPos[i] = c.pos[i] - (c.pos[i]-c.oldPos[i])*(1-spd)
			}
		}
	}
	min, max := c.getEdge(c.edge[0], true), -c.getEdge(c.edge[1], true)
	if c.facing > 0 {
		min, max = -max, -min
	}
	if c.sf(CSF_screenbound) {
		c.drawPos[0] = MaxF(min+sys.xmin/c.localscl, MinF(max+sys.xmax/c.localscl, c.drawPos[0]))
	}
	if c.sf(CSF_movecamera_x) {
		*leftest = MaxF(sys.xmin, MinF(c.drawPos[0]*c.localscl-min*c.localscl, *leftest))
		*rightest = MinF(sys.xmax, MaxF(c.drawPos[0]*c.localscl-max*c.localscl, *rightest))
		if c.acttmp > 0 && !c.sf(CSF_posfreeze) &&
			(c.bindTime == 0 || math.IsNaN(float64(c.bindPos[0]))) {
			*cvmin = MinF(*cvmin, c.vel[0]*c.localscl*c.facing)
			*cvmax = MaxF(*cvmax, c.vel[0]*c.localscl*c.facing)
		}
	}
	if c.sf(CSF_movecamera_y) {
		*highest = MinF(c.drawPos[1]*c.localscl, *highest)
		*lowest = MinF(0, MaxF(c.drawPos[1]*c.localscl, *lowest))
	}
}
func (c *Char) tick() {
	if c.acttmp > 0 && c.anim != nil {
		c.anim.Action()
	}
	if c.bindTime > 0 {
		if c.isBound() {
			if bt := sys.playerID(c.bindToId); bt != nil && !bt.pause() {
				c.setBindTime(c.bindTime - 1)
			}
		} else {
			if !c.pause() {
				c.setBindTime(c.bindTime - 1)
			}
		}
	}
	if c.cmd == nil {
		if c.keyctrl {
			c.cmd = make([]CommandList, len(sys.chars))
			c.cmd[0].Buffer = NewCommandBuffer()
			for i := range c.cmd {
				c.cmd[i].Buffer = c.cmd[0].Buffer
				c.cmd[i].CopyList(sys.chars[c.playerNo][0].cmd[i])
				c.cmd[i].BufReset()
			}
		} else {
			c.cmd = sys.chars[c.playerNo][0].cmd
		}
	}
	if c.hitdefContact {
		if c.hitdef.hitonce != 0 || c.moveReversed() != 0 {
			c.hitdef.invalidate(c.ss.stateType)
		}
		c.hitdefContact = false
	} else if c.hitdef.lhit {
		c.hitdef.attr = c.hitdef.attr&^int32(ST_MASK) | int32(c.ss.stateType)
		c.hitdef.lhit = false
	}
	if c.mctime < 0 {
		c.mctime = 1
		if c.mctype == MC_Hit {
			c.juggle = 0
			c.hitCount += c.hitdef.numhits
		}
	}
	if c.sf(CSF_gethit) {
		c.ss.moveType = MT_H
		if c.hitPauseTime > 0 {
			c.ss.clearWw()
		}
		c.hitPauseTime = 0
		if c.hoIdx >= 0 && c.ho[c.hoIdx].forceair {
			c.ss.stateType = ST_A
		}
		pn := c.playerNo
		if c.ghv.p2getp1state && !c.ghv.guarded {
			pn = c.ghv.playerNo
		}
		if c.stchtmp {
			c.ss.prevno = 0
		} else if c.ss.stateType == ST_L {
			c.changeStateEx(5080, pn, -1, 0)
		} else if c.ghv.guarded && (c.ghv.damage < c.life || sys.sf(GSF_noko)) {
			switch c.ss.stateType {
			case ST_S:
				c.selfState(150, -1, 0)
			case ST_C:
				c.selfState(152, -1, 0)
			case ST_A:
				c.selfState(154, -1, 0)
			}
		} else if c.ghv._type == HT_Trip {
			c.changeStateEx(5070, pn, -1, 0)
		} else {
			if c.ghv.forcestand && c.ss.stateType == ST_C {
				c.ss.stateType = ST_S
			}
			switch c.ss.stateType {
			case ST_S:
				c.changeStateEx(5000, pn, -1, 0)
			case ST_C:
				c.changeStateEx(5010, pn, -1, 0)
			case ST_A:
				c.changeStateEx(5020, pn, -1, 0)
			}
		}
		if c.hoIdx >= 0 {
			c.stateChange1(c.ho[c.hoIdx].stateno, c.ho[c.hoIdx].playerNo)
		}
	}
	if !c.pause() {
		if c.hitPauseTime > 0 {
			c.hitPauseTime--
			if c.hitPauseTime == 0 {
				c.ss.clearWw()
			}
		}
		if c.hitPauseTime <= 0 && c.ss.stateType == ST_L && c.recoverTime > 0 &&
			c.ss.sb.playerNo == c.playerNo && (c.cmd[0].Buffer.Bb == 1 ||
			c.cmd[0].Buffer.Db == 1 || c.cmd[0].Buffer.Fb == 1 ||
			c.cmd[0].Buffer.Ub == 1 || c.cmd[0].Buffer.ab == 1 ||
			c.cmd[0].Buffer.bb == 1 || c.cmd[0].Buffer.cb == 1 ||
			c.cmd[0].Buffer.xb == 1 || c.cmd[0].Buffer.yb == 1 ||
			c.cmd[0].Buffer.zb == 1 || c.cmd[0].Buffer.sb == 1) {
			c.recoverTime -= RandI(1, (c.recoverTime+1)/2)
		}
		if !c.stchtmp {
			if c.helperIndex == 0 && (c.alive() || c.ss.no == 0) && c.life <= 0 &&
				c.ss.moveType != MT_H && !sys.sf(GSF_noko) {
				c.ghv.fallf = true
				c.selfState(5030, -1, -1)
				c.ss.time = 1
			} else if c.ss.no == 5150 && c.ss.time >= 90 && c.alive() {
				c.selfState(5120, -1, -1)
			}
		}
	}
	if !c.hitPause() {
		if c.life <= 0 && !sys.sf(GSF_noko) {
			if !sys.sf(GSF_nokosnd) && c.alive() {
				vo := int32(0)
				if c.gi().ver[0] == 1 {
					vo = 100
				}
				c.playSound(false, false, false, 11, 0, -1, vo, 0, 1, &c.pos[0])
			}
			c.setSCF(SCF_ko)
		}
		if c.ss.moveType != MT_H {
			c.recoverTime = c.gi().data.liedown.time
		}
		if c.ss.no == 5110 && c.recoverTime > 0 && !c.pause() {
			c.recoverTime--
		}
	}
}
func (c *Char) cueDraw() {
	if c.helperIndex < 0 {
		return
	}
	if sys.clsnDraw && c.curFrame != nil {
		x, y := c.pos[0]*c.oldLocalscl+c.offsetX()*c.oldLocalscl, c.pos[1]*c.oldLocalscl+c.offsetY()*c.oldLocalscl
		xs, ys := c.facing*c.clsnScale[0]*(320/float32(sys.chars[c.animPN][0].localcoord)), c.clsnScale[1]*(320/float32(sys.chars[c.animPN][0].localcoord))
		if clsn := c.curFrame.Clsn1(); len(clsn) > 0 && c.atktmp != 0 {
			sys.drawc1.Add(clsn, x, y, xs, ys)
		}
		if clsn := c.curFrame.Clsn2(); len(clsn) > 0 {
			hb, mtk := false, false
			for _, h := range c.hitby {
				if h.time != 0 {
					hb = true
					mtk = mtk || h.flag&int32(ST_SCA) == 0 || h.flag&int32(AT_ALL) == 0
				}
			}
			if mtk {
				sys.drawc2mtk.Add(clsn, x, y, xs, ys)
			} else if hb {
				sys.drawc2sp.Add(clsn, x, y, xs, ys)
			} else {
				sys.drawc2.Add(clsn, x, y, xs, ys)
			}
		}
		if c.sf(CSF_playerpush) {
			sys.drawwh.Add([]float32{-c.width[1] * c.oldLocalscl, -c.height() * (320 / float32(c.localcoord)), c.width[0] * c.oldLocalscl, 0},
				c.pos[0]*c.oldLocalscl, c.pos[1]*c.oldLocalscl, c.facing, 1)
		}
	}
	if c.anim != nil {
		pos := [...]float32{c.drawPos[0]*c.oldLocalscl + c.offsetX()*c.oldLocalscl, c.drawPos[1]*c.oldLocalscl + c.offsetY()*c.oldLocalscl}
		scl := [...]float32{c.facing * c.size.xscale * (320 / float32(c.localcoord)), c.size.yscale * (320 / float32(c.localcoord))}
		agl := float32(0)
		if c.sf(CSF_angledraw) {
			agl = c.angle
			if agl == 0 {
				agl = 360
			} else if c.facing < 0 {
				agl *= -1
			}
		}
		rec := sys.tickNextFrame() && c.acttmp > 0
		sdf := func() *SprData {
			sd := &SprData{c.anim, c.getPalfx(), pos,
				scl, c.alpha, c.sprPriority, agl, c.angleScalse, false,
				c.playerNo == sys.superplayer, c.gi().ver[0] != 1, c.facing}
			if !c.sf(CSF_trans) {
				sd.alpha[0] = -1
			}
			return sd
		}
		if c.sf(CSF_invisible) {
			if rec {
				c.aimg.recAfterImg(sdf())
			}
		} else {
			//if c.gi().ver[0] != 1 && c.sf(CSF_angledraw) && !c.sf(CSF_trans) {
			//	c.setSF(CSF_trans)
			//	c.alpha = [...]int32{255, 0}
			//}
			sd := sdf()
			c.aimg.recAndCue(sd, rec)
			if c.ghv.hitshaketime > 0 && c.ss.time&1 != 0 {
				sd.pos[0] -= c.facing
			}
			var sc, sa int32 = -1, 255
			if c.sf(CSF_noshadow) {
				sc = 0
			}
			if c.sf(CSF_trans) {
				sa = 255 - c.alpha[1]
			}
			sys.sprites.add(sd, sc, sa, float32(c.size.shadowoffset), c.offsetY())
		}
	}
	if sys.tickNextFrame() {
		if c.roundState() == 4 {
			c.exitTarget(false)
		}
		if sys.supertime < 0 && c.playerNo&1 != sys.superplayer&1 {
			c.defenceMul *= sys.superp2defmul
		}
		c.minus = 2
		c.oldPos = c.pos
	}
}
func (c *Char) victoryQuote(v int32) {
	c.winquote = v
}

type CharList struct {
	runOrder, drawOrder []*Char
	idMap               map[int32]*Char
}

func (cl *CharList) clear() {
	*cl = CharList{idMap: make(map[int32]*Char)}
	sys.nextCharId = sys.helperMax
}
func (cl *CharList) add(c *Char) {
	cl.runOrder = append(cl.runOrder, c)
	i := 0
	for ; i < len(cl.drawOrder); i++ {
		if cl.drawOrder[i] == nil {
			cl.drawOrder[i] = c
			break
		}
	}
	if i >= len(cl.drawOrder) {
		cl.drawOrder = append(cl.drawOrder, c)
	}
	cl.idMap[c.id] = c
}
func (cl *CharList) delete(dc *Char) {
	for i, c := range cl.runOrder {
		if c == dc {
			delete(cl.idMap, c.id)
			cl.runOrder = append(cl.runOrder[:i], cl.runOrder[i+1:]...)
			break
		}
	}
	for i, c := range cl.drawOrder {
		if c == dc {
			cl.drawOrder[i] = nil
			break
		}
	}
}
func (cl *CharList) action(x float32, cvmin, cvmax,
	highest, lowest, leftest, rightest *float32) {
	sys.commandUpdate()
	for i := 0; i < len(cl.runOrder); i++ {
		if cl.runOrder[i].ss.moveType == MT_A {
			cl.runOrder[i].action()
		}
	}
	for i := 0; i < len(cl.runOrder); i++ {
		cl.runOrder[i].action()
	}
	sys.charUpdate(cvmin, cvmax, highest, lowest, leftest, rightest)
}
func (cl *CharList) update(cvmin, cvmax,
	highest, lowest, leftest, rightest *float32) {
	ro := make([]*Char, len(cl.runOrder))
	copy(ro, cl.runOrder)
	for _, c := range ro {
		c.update(cvmin, cvmax, highest, lowest, leftest, rightest)
	}
}
func (cl *CharList) clsn(getter *Char, proj bool) {
	var gxmin, gxmax float32
	hit := func(c *Char, hd *HitDef, pos [2]float32,
		projf float32, hits int32) (hitType int32) {
		if !proj && c.ss.stateType == ST_L && hd.reversal_attr <= 0 {
			c.hitdef.lhit = true
			return 0
		}
		if getter.stchtmp && getter.ss.sb.playerNo != hd.playerNo && func() bool {
			if getter.sf(CSF_gethit) {
				return hd.p2stateno >= 0
			}
			return getter.acttmp > 0
		}() || hd.p1stateno >= 0 && (c.sf(CSF_gethit) || c.stchtmp &&
			c.ss.sb.playerNo != hd.playerNo) {
			return 0
		}
		guard := (proj || !c.sf(CSF_unguardable)) && getter.scf(SCF_guard) &&
			(!getter.sf(CSF_gethit) || getter.ghv.guarded)
		if guard && sys.autoguard[getter.playerNo] &&
			getter.acttmp > 0 && !getter.sf(CSF_gethit) &&
			(getter.ss.stateType == ST_S || getter.ss.stateType == ST_C) &&
			int32(getter.ss.stateType)&hd.guardflag == 0 {
			if int32(ST_S)&hd.guardflag != 0 && !getter.sf(CSF_nostandguard) {
				getter.ss.stateType = ST_S
			} else if int32(ST_C)&hd.guardflag != 0 &&
				!getter.sf(CSF_nocrouchguard) {
				getter.ss.stateType = ST_C
			}
		}
		hitType = 1
		if guard && int32(getter.ss.stateType)&hd.guardflag != 0 {
			hitType = 2
		}
		if hd.reversal_attr > 0 {
			hitType *= -1
		} else if getter.ss.stateType == ST_A {
			if hd.air_type == HT_None {
				hitType *= -1
			}
		} else if hd.ground_type == HT_None {
			hitType *= -1
		}
		p2s := false
		if !getter.stchtmp || !getter.sf(CSF_gethit) {
			_break := false
			for i, ho := range getter.ho {
				if ho.time == 0 || ho.attr&hd.attr&^int32(ST_MASK) == 0 {
					continue
				}
				if proj {
					if ho.attr&hd.attr&int32(ST_MASK) == 0 {
						continue
					}
				} else {
					if ho.attr&int32(c.ss.stateType) == 0 {
						continue
					}
				}
				if !proj && Abs(hitType) == 1 &&
					(hd.p2stateno >= 0 || hd.p1stateno >= 0) {
					return 0
				}
				if ho.stateno >= 0 {
					getter.hoIdx = i
					_break = true
					break
				}
			}
			if !_break {
				if Abs(hitType) == 1 && hd.p2stateno >= 0 {
					pn := getter.playerNo
					if hd.p2getp1state {
						pn = hd.playerNo
					}
					if getter.stateChange1(hd.p2stateno, pn) {
						getter.setCtrl(false)
						p2s = true
						getter.hoIdx = -1
					}
				}
			}
		}
		c.targetsOfHitdef = append(c.targetsOfHitdef, getter.id)
		ghvset := !getter.stchtmp || p2s || !getter.sf(CSF_gethit)
		if ghvset {
			if !proj {
				c.sprPriority = hd.p1sprpriority
			}
			getter.sprPriority = hd.p2sprpriority
			getter.ghv.hitid = hd.id
			getter.ghv.playerNo = hd.playerNo
		}
		if Abs(hitType) == 1 {
			if hd.pausetime > 0 {
				hits = 1
			}
		} else if hd.guard_pausetime > 0 {
			hits = 1
		}
		byf := c.facing
		if proj {
			byf = projf
		}
		if !proj && hitType == 1 {
			if hd.p1getp2facing != 0 {
				byf = getter.facing
				if hd.p1getp2facing < 0 {
					byf *= -1
				}
			} else if hd.p1facing < 0 {
				byf *= -1
			}
		}
		if hitType > 0 {
			if hitType == 1 && len(getter.sounds) > 0 {
				getter.sounds[0].sound = nil
			}
			if getter.bindToId == c.id {
				getter.setBindTime(0)
			}
			absdamage := int32(0)
			if ghvset {
				ghv := &getter.ghv
				cmb := (getter.ss.moveType == MT_H || getter.sf(CSF_gethit)) &&
					!ghv.guarded
				fall, hc, fc, by := ghv.fallf, ghv.hitcount, ghv.fallcount, ghv.hitBy
				ghv.clear()
				ghv.hitBy = by
				ghv.attr = hd.attr
				ghv.hitid = hd.id
				ghv.playerNo = hd.playerNo
				ghv.p2getp1state = hd.p2getp1state
				ghv.forcestand = hd.forcestand != 0
				ghv.fall = hd.fall
				getter.fallTime = 0
				ghv.fall.xvelocity = hd.fall.xvelocity * c.localscl / getter.localscl
				ghv.fall.yvelocity = hd.fall.yvelocity * c.localscl / getter.localscl
				ghv.yaccel = hd.yaccel * c.localscl / getter.localscl
				if hd.forcenofall {
					fall = false
				}
				ghv.groundtype = hd.ground_type
				ghv.airtype = hd.air_type
				if getter.ss.stateType == ST_A {
					ghv._type = ghv.airtype
				} else {
					ghv._type = ghv.groundtype
				}
				ghv.airanimtype = hd.air_animtype
				ghv.groundanimtype = hd.animtype
				if guard && int32(getter.ss.stateType)&hd.guardflag != 0 {
					ghv.hitshaketime = Max(0, hd.guard_shaketime)
					ghv.hittime = Max(0, hd.guard_hittime)
					ghv.slidetime = hd.guard_slidetime
					ghv.guarded = true
					if getter.ss.stateType == ST_A {
						ghv.ctrltime = hd.airguard_ctrltime
						ghv.xvel = hd.airguard_velocity[0] * c.localscl / getter.localscl
						ghv.yvel = hd.airguard_velocity[1] * c.localscl / getter.localscl
					} else {
						ghv.ctrltime = hd.guard_ctrltime
						ghv.xvel = hd.guard_velocity * c.localscl / getter.localscl
						ghv.yvel = hd.ground_velocity[1] * c.localscl / getter.localscl
					}
					absdamage = hd.guarddamage
					ghv.hitcount = hc
				} else {
					ghv.hitshaketime = Max(0, hd.shaketime)
					ghv.slidetime = hd.ground_slidetime
					if getter.ss.stateType == ST_A {
						ghv.hittime = hd.air_hittime
						ghv.ctrltime = hd.air_hittime
						ghv.xvel = hd.air_velocity[0] * c.localscl / getter.localscl
						ghv.yvel = hd.air_velocity[1] * c.localscl / getter.localscl
						ghv.fallf = hd.air_fall
					} else if getter.ss.stateType == ST_L {
						ghv.hittime = hd.down_hittime
						ghv.ctrltime = hd.down_hittime
						ghv.xvel = hd.down_velocity[0] * c.localscl / getter.localscl
						ghv.yvel = hd.down_velocity[1] * c.localscl / getter.localscl
						if !hd.down_bounce {
							ghv.fall.xvelocity = float32(math.NaN())
							ghv.fall.yvelocity = 0
						}
					} else {
						ghv.hittime = hd.ground_hittime
						ghv.ctrltime = hd.ground_hittime
						ghv.xvel = hd.ground_velocity[0] * c.localscl / getter.localscl
						ghv.yvel = hd.ground_velocity[1] * c.localscl / getter.localscl
						ghv.fallf = hd.ground_fall
						if ghv.fallf && ghv.yvel == 0 {
							ghv.yvel = -0.001 * c.localscl / getter.localscl //新MUGENだとウィンドウサイズを大きくするとここに入る数値が小さくなるが、再現しないほうがよいと思う。
						}
					}
					if ghv.hittime < 0 {
						ghv.hittime = 0
					}
					absdamage = hd.hitdamage
					if cmb {
						ghv.hitcount = hc + 1
					} else {
						ghv.hitcount = 1
					}
					ghv.fallcount = fc
					ghv.fallf = ghv.fallf || fall
				}
				byPos := c.pos
				if proj {
					for i, p := range pos {
						byPos[i] += p
					}
				}
				snap := [...]float32{float32(math.NaN()), float32(math.NaN())}
				if !math.IsNaN(float64(hd.mindist[0])) {
					if byf < 0 {
						if getter.pos[0] > byPos[0]-hd.mindist[0] {
							snap[0] = byPos[0] - hd.mindist[0]
						}
					} else {
						if getter.pos[0] < byPos[0]+hd.mindist[0] {
							snap[0] = byPos[0] + hd.mindist[0]
						}
					}
				}
				if !math.IsNaN(float64(hd.maxdist[0])) {
					if byf < 0 {
						if getter.pos[0]*getter.localscl/c.localscl < byPos[0]-hd.maxdist[0] {
							snap[0] = byPos[0] - hd.maxdist[0]
						}
					} else {
						if getter.pos[0]*getter.localscl/c.localscl > byPos[0]+hd.maxdist[0] {
							snap[0] = byPos[0] + hd.maxdist[0]
						}
					}
				}
				if hitType == 1 || getter.ss.stateType == ST_A {
					if !math.IsNaN(float64(hd.mindist[1])) {
						if getter.pos[1]*getter.localscl/c.localscl < byPos[1]+hd.mindist[1] {
							snap[1] = byPos[1] + hd.mindist[1]
						}
					}
					if !math.IsNaN(float64(hd.maxdist[1])) {
						if getter.pos[1]*getter.localscl/c.localscl > byPos[1]+hd.maxdist[1] {
							snap[1] = byPos[1] + hd.maxdist[1]
						}
					}
				}
				if !math.IsNaN(float64(snap[0])) {
					ghv.xoff = snap[0]*c.localscl/getter.localscl - getter.pos[0]
				}
				if !math.IsNaN(float64(snap[1])) {
					ghv.yoff = snap[1]*c.localscl/getter.localscl - getter.pos[1]
				}
				if hd.snapt != 0 && getter.hoIdx < 0 {
					getter.setBindToId(c)
					getter.setBindTime(hd.snapt + Btoi(hd.snapt > 0 && !c.pause()))
					getter.bindFacing = 0
					if !math.IsNaN(float64(snap[0])) {
						getter.bindPos[0] = hd.mindist[0] * c.localscl / getter.localscl
					} else {
						getter.bindPos[0] = float32(math.NaN())
					}
					if !math.IsNaN(float64(snap[1])) &&
						(hitType == 1 || getter.ss.stateType == ST_A) {
						getter.bindPos[1] = hd.mindist[1] * c.localscl / getter.localscl
					} else {
						getter.bindPos[1] = float32(math.NaN())
					}
				}
			} else if hitType == 1 {
				absdamage = hd.hitdamage
			} else {
				absdamage = hd.guarddamage
			}
			if sys.super > 0 {
				getter.superMovetime =
					Max(getter.superMovetime, getter.ghv.hitshaketime)
			} else if sys.pause > 0 {
				getter.pauseMovetime =
					Max(getter.pauseMovetime, getter.ghv.hitshaketime)
			}
			if !p2s && !getter.sf(CSF_gethit) {
				getter.stchtmp = false
			}
			getter.setSF(CSF_gethit)
			live, kill := getter.life > 0, hd.kill
			if hitType == 2 {
				kill = hd.guard_kill
			}
			getter.ghv.damage += getter.computeDamage(
				float64(absdamage)*float64(hits), kill, false, c.attackMul)
			if ghvset && getter.ghv.damage >= getter.life {
				if kill || !live {
					getter.ghv.fallf = true
					if getter.ghv.fall.animtype < RA_Back {
						getter.ghv.fall.animtype = RA_Back
					}
					if getter.ss.stateType == ST_A {
						if getter.ghv.xvel < 0 {
							getter.ghv.xvel -= 2 / c.localscl * getter.localscl
						}
						if getter.ghv.yvel <= 0 {
							getter.ghv.yvel -= 2 / c.localscl * getter.localscl
							if getter.ghv.yvel > -3 {
								getter.ghv.yvel = -3 / c.localscl * getter.localscl
							}
						}
					} else {
						if getter.ghv.yvel == 0 {
							getter.ghv.xvel *= 0.66
						}
						if getter.ghv.xvel < 0 {
							getter.ghv.xvel -= 2.5 / c.localscl * getter.localscl
						}
						if getter.ghv.yvel <= 0 {
							getter.ghv.yvel -= 2 / c.localscl * getter.localscl
							if getter.ghv.yvel > -6 {
								getter.ghv.yvel = -6 / c.localscl * getter.localscl
							}
						}
					}
					getter.ghv.damage = getter.life
				} else {
					getter.ghv.damage = getter.life - 1
				}
			}
		}
		hitspark := func(p1, p2 *Char, animNo int32) {
			ffx := animNo < 0
			if ffx {
				animNo ^= -1
			}
			off := pos
			if !proj {
				off[0] = p2.pos[0]*p2.localscl - p1.pos[0]*p1.localscl
				if (p1.facing < 0) != (p2.facing < 0) {
					off[0] += p2.facing * p2.width[0] * p2.localscl
				} else {
					off[0] -= p2.facing * p2.width[1] * p2.localscl
				}
			}
			off[0] *= p1.facing
			if proj {
				off[0] *= c.localscl
				off[1] *= c.localscl
				off[0] += hd.sparkxy[0] * projf * p1.facing * c.localscl
			} else {
				off[0] -= hd.sparkxy[0] * c.localscl
			}
			off[1] += hd.sparkxy[1] * c.localscl
			if c.id != p1.id {
				off[1] += p1.hitdef.sparkxy[1] * c.localscl
			}
			if e, i := c.newExplod(); e != nil {
				e.anim = c.getAnim(animNo, ffx)
				e.ontop = true
				e.sprpriority = math.MinInt32
				e.ownpal = true
				e.offset = off
				e.supermovetime = -1
				e.pausemovetime = -1
				e.localscl = 1
				if !ffx {
					e.scale = [...]float32{c.localscl, c.localscl}
				}
				e.setPos(p1)
				c.insertExplod(i)
			}
		}
		if Abs(hitType) == 1 {
			if hd.sparkno != IErr {
				if hd.reversal_attr > 0 {
					hitspark(getter, c, hd.sparkno)
				} else {
					hitspark(c, getter, hd.sparkno)
				}
			}
			if hd.hitsound[0] != IErr {
				sg := hd.hitsound[0]
				f := sg < 0
				if f {
					sg ^= -1
				}
				vo := int32(0)
				if c.gi().ver[0] == 1 {
					vo = 100
				}
				c.playSound(f, false, false, sg, hd.hitsound[1],
					-1, vo, 0, 1, &getter.pos[0])
			}
			if hitType > 0 {
				c.powerAdd(hd.hitgetpower)
				if getter.player {
					getter.powerAdd(hd.hitgivepower)
				}
			}
		} else {
			if hd.guard_sparkno != IErr {
				if hd.reversal_attr > 0 {
					hitspark(getter, c, hd.guard_sparkno)
				} else {
					hitspark(c, getter, hd.guard_sparkno)
				}
			}
			if hd.guardsound[0] != IErr {
				sg := hd.guardsound[0]
				f := sg < 0
				if f {
					sg ^= -1
				}
				vo := int32(0)
				if c.gi().ver[0] == 1 {
					vo = 100
				}
				c.playSound(f, false, false, sg, hd.guardsound[1],
					-1, vo, 0, 1, &getter.pos[0])
			}
			if hitType > 0 {
				c.powerAdd(hd.guardgetpower)
				if getter.player {
					getter.powerAdd(hd.guardgivepower)
				}
			}
		}
		if !ghvset {
			return
		}
		getter.p1facing = 0
		if getter.hoIdx >= 0 {
			return
		}
		if !proj && hd.hitonce > 0 {
			c.targetDrop(-1, false)
		}
		c.addTarget(getter.id)
		getter.ghv.addId(c.id, c.gi().data.airjuggle)
		xmi, xma := gxmin+2, gxmax-2
		if c.stCgi().ver[0] != 1 {
			xmi += 2
			xma -= 2
		}
		if Abs(hitType) == 1 {
			if !proj && (hd.p1getp2facing != 0 || hd.p1facing < 0) &&
				c.facing != byf {
				c.p1facing = byf
			}
			if hd.p2facing < 0 {
				getter.p1facing = byf
			} else if hd.p2facing > 0 {
				getter.p1facing = -byf
			}
			if getter.p1facing == getter.facing {
				getter.p1facing = 0
			}
			if hd.p1stateno >= 0 && c.stateChange1(hd.p1stateno, hd.playerNo) {
				c.setCtrl(false)
			}
			if getter.ghv.fallf && !c.sf(CSF_nojugglecheck) {
				jug := &getter.ghv.hitBy[len(getter.ghv.hitBy)-1][1]
				if proj {
					*jug -= hd.air_juggle
				} else {
					*jug -= c.juggle
				}
			}
			if hd.palfx.time > 0 && getter.palfx != nil {
				getter.palfx.clear2(true)
				getter.palfx.PalFXDef = hd.palfx
			}
			if hd.envshake_time > 0 {
				sys.envShake.time = hd.envshake_time
				sys.envShake.freq = hd.envshake_freq * float32(math.Pi) / 180
				sys.envShake.ampl = int32(float32(hd.envshake_ampl) * c.localscl)
				sys.envShake.phase = hd.envshake_phase * c.localscl
				sys.envShake.setDefPhase()
			}
			getter.getcombo += hd.numhits * hits
			if hitType > 0 && !proj && getter.sf(CSF_screenbound) &&
				(c.facing < 0 && getter.pos[0] <= xmi ||
					c.facing > 0 && getter.pos[0] >= xma) {
				switch getter.ss.stateType {
				case ST_S, ST_C:
					c.veloff = hd.ground_cornerpush_veloff * c.facing
				case ST_A:
					c.veloff = hd.air_cornerpush_veloff * c.facing
				case ST_L:
					c.veloff = hd.down_cornerpush_veloff * c.facing
				}
			}
		} else {
			if hitType > 0 && !proj && getter.sf(CSF_screenbound) &&
				(c.facing < 0 && getter.pos[0] <= xmi ||
					c.facing > 0 && getter.pos[0] >= xma) {
				switch getter.ss.stateType {
				case ST_S, ST_C:
					c.veloff = hd.guard_cornerpush_veloff * c.facing
				case ST_A:
					c.veloff = hd.airguard_cornerpush_veloff * c.facing
				}
			}
		}
		if !proj {
			if c.p1facing != 0 {
				byf = c.p1facing
			} else {
				byf = c.facing
			}
		}
		if (getter.facing < 0) == (byf < 0) {
			getter.ghv.xvel *= -1
			if getter.ghv.groundtype == 1 || getter.ghv.groundtype == 2 {
				getter.ghv.groundtype += 3 - getter.ghv.groundtype*2
			}
			if getter.ghv.airtype == 1 || getter.ghv.airtype == 2 {
				getter.ghv.airtype += 3 - getter.ghv.airtype*2
			}
		}
		return
	}
	if proj {
		for i, pr := range sys.projs {
			if i == getter.playerNo || len(sys.projs[i]) == 0 {
				continue
			}
			c := sys.chars[i][0]
			orgatktmp := c.atktmp
			c.atktmp = -1
			for j := range pr {
				p := &pr[j]
				if p.id < 0 || p.hits < 0 || p.hitdef.affectteam != 0 &&
					(getter.playerNo&1 != i&1) != (p.hitdef.affectteam > 0) {
					continue
				}
				if dist := (getter.pos[0]*getter.localscl - p.pos[0]*p.localscl) * p.facing; dist >= 0 &&
					dist <= float32(c.size.proj.attack.dist)*c.localscl {
					getter.inguarddist = true
				}
				if p.hits == 0 {
					continue
				}
				if getter.atktmp != 0 && (getter.hitdef.affectteam == 0 ||
					(i&1 != getter.playerNo&1) == (getter.hitdef.affectteam > 0)) &&
					getter.hitdef.hitflag&int32(ST_P) != 0 &&
					getter.projClsnCheck(p, false) {
					p.hits = -2
					sys.cgi[i].pctype = PC_Cancel
					sys.cgi[i].pctime = 0
					sys.cgi[i].pcid = p.id
					getter.hitdefContact = true
					continue
				}
				if !(getter.stchtmp && (getter.sf(CSF_gethit) || getter.acttmp > 0)) &&
					(c.sf(CSF_nojugglecheck) || getter.ghv.getJuggle(c.id,
						c.gi().data.airjuggle) >= p.hitdef.air_juggle) &&
					p.timemiss <= 0 && p.hitpause <= 0 && getter.hittable(&p.hitdef,
					c, ST_N, func(h *HitDef) bool { return false }) {
					orghittmp := getter.hittmp
					if getter.sf(CSF_gethit) {
						getter.hittmp = int8(Btoi(getter.ghv.fallf)) + 1
					}
					if dist := -getter.distX(c, getter) * c.facing; dist >= 0 &&
						dist <= float32(p.hitdef.guard_dist) {
						getter.inguarddist = true
					}
					if getter.projClsnCheck(p, true) {
						hits := p.hits
						if p.misstime > 0 {
							hits = 1
						}
						if ht := hit(c, &p.hitdef, [...]float32{p.pos[0] - c.pos[0]*c.localscl/p.localscl,
							p.pos[1] - c.pos[1]*c.localscl/p.localscl}, p.facing, hits); ht != 0 {
							p.timemiss = ^Max(0, p.misstime)
							if Abs(ht) == 1 {
								sys.cgi[i].pctype = PC_Hit
								sys.cgi[i].pctime = 0
								sys.cgi[i].pcid = p.id
								p.hitpause = Max(0, p.hitdef.pausetime)
							} else {
								sys.cgi[i].pctype = PC_Guarded
								sys.cgi[i].pctime = 0
								sys.cgi[i].pcid = p.id
								p.hitpause = Max(0, p.hitdef.guard_pausetime)
							}
						}
					}
					getter.hittmp = orghittmp
				}
			}
			c.atktmp = orgatktmp
		}
	} else {
		gxmin = getter.getEdge(getter.edge[0], true)
		gxmax = -getter.getEdge(getter.edge[1], true)
		if getter.facing > 0 {
			gxmin, gxmax = -gxmax, -gxmin
		}
		gxmin += sys.xmin / getter.localscl
		gxmax += sys.xmax / getter.localscl
		getter.inguarddist = false
		getter.unsetSF(CSF_gethit)
		gl, gr := -getter.width[0]*getter.localscl, getter.width[1]*getter.localscl
		if getter.facing > 0 {
			gl, gr = -gr, -gl
		}
		gl += getter.pos[0] * getter.localscl
		gr += getter.pos[0] * getter.localscl
		getter.enemyNearClear()
		for _, c := range cl.runOrder {
			contact := 0
			if c.atktmp != 0 && c.id != getter.id && (c.hitdef.affectteam == 0 ||
				(getter.playerNo&1 != c.playerNo&1) == (c.hitdef.affectteam > 0)) {
				dist := -getter.distX(c, getter) * c.facing
				if c.ss.moveType == MT_A && dist >= 0 && dist <= c.attackDist {
					getter.inguarddist = true
				}
				if c.hitdef.hitonce >= 0 && !c.hasTargetOfHitdef(getter.id) &&
					(c.hitdef.reversal_attr <= 0 || !getter.hasTargetOfHitdef(c.id)) &&
					(getter.hittmp < 2 || c.sf(CSF_nojugglecheck) ||
						getter.ghv.getJuggle(c.id, c.gi().data.airjuggle) >= c.juggle) &&
					getter.hittable(&c.hitdef, c, c.ss.stateType, func(h *HitDef) bool {
						return (c.atktmp >= 0 || !getter.hasTarget(c.id)) &&
							c.attrCheck(h, getter.id, getter.ss.stateType) &&
							c.hitCheck(getter)
					}) {
					if c.ss.moveType == MT_A && dist >= 0 &&
						dist <= float32(c.hitdef.guard_dist) {
						getter.inguarddist = true
					}
					if c.hitdef.reversal_attr <= 0 {
						contact = -1
					}
					if getter.hitCheck(c) {
						if contact < 0 {
							contact = 1
						}
						if ht := hit(c, &c.hitdef, [2]float32{}, 0, 1); ht != 0 {
							mvh := ht > 0 || c.hitdef.reversal_attr > 0
							if Abs(ht) == 1 {
								if mvh {
									c.mctype = MC_Hit
								}
								if c.hitdef.reversal_attr > 0 {
									getter.mctype = MC_Reversed
									getter.mctime = -1
									getter.hitdefContact = true
									getter.targetsOfHitdef = append(getter.targetsOfHitdef, c.id)
									if getter.hittmp == 0 {
										getter.hittmp = -1
									}
									if !getter.sf(CSF_gethit) {
										getter.hitPauseTime = Max(1, c.hitdef.shaketime+
											Btoi(c.gi().ver[0] == 1))
									}
								}
								if !c.sf(CSF_gethit) {
									c.hitPauseTime = Max(1, c.hitdef.pausetime+
										Btoi(c.gi().ver[0] == 1))
								}
								c.uniqHitCount++
							} else {
								if mvh {
									c.mctype = MC_Guarded
								}
								if !c.sf(CSF_gethit) {
									c.hitPauseTime = Max(1, c.hitdef.guard_pausetime+
										Btoi(c.gi().ver[0] == 1))
								}
							}
							if c.hitdef.hitonce > 0 {
								c.hitdef.hitonce = -1
							}
							if mvh {
								c.mctime = -1
							}
							c.hitdefContact = true
						}
					}
				}
			}
			if getter.playerNo&1 != c.playerNo&1 && getter.sf(CSF_playerpush) &&
				c.sf(CSF_playerpush) && (getter.ss.stateType == ST_A ||
				getter.pos[1]*getter.localscl-c.pos[1]*c.localscl < getter.height()*c.localscl) &&
				(c.ss.stateType == ST_A || c.pos[1]*c.localscl-getter.pos[1]*getter.localscl < c.height()*(320/float32(c.localcoord))) {
				cl, cr := -c.width[0]*c.localscl, c.width[1]*c.localscl
				if c.facing > 0 {
					cl, cr = -cr, -cl
				}
				cl += c.pos[0] * c.localscl
				cr += c.pos[0] * c.localscl
				if gl < cr && cl < gr && (contact > 0 ||
					getter.clsnCheck(c, false, false)) {
					getter.pushed, c.pushed = true, true
					tmp := getter.distX(c, getter)
					if tmp == 0 {
						if getter.pos[1]*getter.localscl > c.pos[1]*c.localscl {
							tmp = getter.facing
						} else {
							tmp = -c.facing
						}
					}
					if tmp > 0 {
						getter.pos[0] -= ((gr - cl) * 0.5) / getter.localscl
						c.pos[0] += ((gr - cl) * 0.5) / c.localscl
					} else {
						getter.pos[0] += ((cr - gl) * 0.5) / getter.localscl
						c.pos[0] -= ((cr - gl) * 0.5) / c.localscl
					}
					if getter.sf(CSF_screenbound) {
						getter.pos[0] = MaxF(gxmin, MinF(gxmax, getter.pos[0]))
					}
					if c.sf(CSF_screenbound) {
						l, r := c.getEdge(c.edge[0], true), -c.getEdge(c.edge[1], true)
						if c.facing > 0 {
							l, r = -r, -l
						}
						c.pos[0] = MaxF(l+sys.xmin/c.localscl, MinF(r+sys.xmax/c.localscl, c.pos[0]))
					}
					getter.pos[0] = MaxF(sys.stage.leftbound/getter.localscl, MinF(sys.stage.rightbound/getter.localscl,
						getter.pos[0]))
					c.pos[0] = MaxF(sys.stage.leftbound/c.localscl, MinF(sys.stage.rightbound/c.localscl,
						c.pos[0]))
					getter.drawPos[0], c.drawPos[0] = getter.pos[0], c.pos[0]
				}
			}
		}
	}
}
func (cl *CharList) getHit() {
	for _, c := range cl.runOrder {
		cl.clsn(c, false)
	}
	for _, c := range cl.runOrder {
		cl.clsn(c, true)
	}
}
func (cl *CharList) tick() {
	sys.gameTime++
	for i := range sys.cgi {
		if sys.cgi[i].unhittable > 0 {
			sys.cgi[i].unhittable--
		}
	}
	for _, c := range cl.runOrder {
		c.tick()
	}
}
func (cl *CharList) cueDraw() {
	for _, c := range cl.drawOrder {
		if c != nil {
			c.cueDraw()
		}
	}
}
func (cl *CharList) get(id int32) *Char {
	if id < 0 {
		return nil
	}
	return cl.idMap[id]
}
func (cl *CharList) enemyNear(c *Char, n int32, p2 bool) *Char {
	if n < 0 {
		return nil
	}
	cache := &c.enemynear[Btoi(p2)]
	if int(n) < len(*cache) {
		return (*cache)[n]
	}
	*cache = (*cache)[:0]
	var add func(*Char, int)
	add = func(e *Char, idx int) {
		for i := idx; i <= int(n); i++ {
			if i >= len(*cache) {
				*cache = append(*cache, e)
				return
			}
			if p2 && !e.scf(SCF_ko_round_middle) &&
				(*cache)[i].scf(SCF_ko_round_middle) || (!p2 ||
				e.scf(SCF_ko_round_middle) == (*cache)[i].scf(SCF_ko_round_middle)) &&
				AbsF(c.distX(e, c)) < AbsF(c.distX((*cache)[i], c)) {
				add((*cache)[i], i+1)
				(*cache)[i] = e
			}
		}
	}
	for _, e := range cl.runOrder {
		if e.player && e.playerNo&1 != c.playerNo&1 {
			add(e, 0)
		}
	}
	if int(n) >= len(*cache) {
		return nil
	}
	return (*cache)[n]
}

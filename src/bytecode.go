package main

import (
	"math"
	"unsafe"
)

type StateType int32

const (
	ST_S StateType = 1 << iota
	ST_C
	ST_A
	ST_L
	ST_N
	ST_U
	ST_MASK = 1<<iota - 1
	ST_D    = ST_L
	ST_F    = ST_N
	ST_P    = ST_U
	ST_SCA  = ST_S | ST_C | ST_A
)

type AttackType int32

const (
	AT_NA AttackType = 1 << (iota + 6)
	AT_NT
	AT_NP
	AT_SA
	AT_ST
	AT_SP
	AT_HA
	AT_HT
	AT_HP
	AT_AA  = AT_NA | AT_SA | AT_HA
	AT_AT  = AT_NT | AT_ST | AT_HT
	AT_AP  = AT_NP | AT_SP | AT_HP
	AT_ALL = AT_AA | AT_AT | AT_AP
	AT_AN  = AT_NA | AT_NT | AT_NP
	AT_AS  = AT_SA | AT_ST | AT_SP
	AT_AH  = AT_HA | AT_HT | AT_HP
)

type MoveType int32

const (
	MT_I MoveType = 1 << (iota + 15)
	MT_H
	MT_A
	MT_U
	MT_MNS = MT_I
	MT_PLS = MT_H
)

type ValueType int

const (
	VT_None ValueType = iota
	VT_Float
	VT_Int
	VT_Bool
	VT_SFalse
)

type OpCode byte

const (
	OC_var OpCode = iota + 110
	OC_sysvar
	OC_fvar
	OC_sysfvar
	OC_localvar
	OC_int8
	OC_int
	OC_float
	OC_pop
	OC_dup
	OC_swap
	OC_run
	OC_nordrun
	OC_jsf8
	OC_jmp8
	OC_jz8
	OC_jnz8
	OC_jmp
	OC_jz
	OC_jnz
	OC_eq
	OC_ne
	OC_gt
	OC_le
	OC_lt
	OC_ge
	OC_neg
	OC_blnot
	OC_bland
	OC_blxor
	OC_blor
	OC_not
	OC_and
	OC_xor
	OC_or
	OC_add
	OC_sub
	OC_mul
	OC_div
	OC_mod
	OC_pow
	OC_abs
	OC_exp
	OC_ln
	OC_log
	OC_cos
	OC_sin
	OC_tan
	OC_acos
	OC_asin
	OC_atan
	OC_floor
	OC_ceil
	OC_ifelse
	OC_time
	OC_animtime
	OC_animelemtime
	OC_animelemno
	OC_statetype
	OC_movetype
	OC_ctrl
	OC_command
	OC_random
	OC_pos_x
	OC_pos_y
	OC_vel_x
	OC_vel_y
	OC_screenpos_x
	OC_screenpos_y
	OC_facing
	OC_anim
	OC_animexist
	OC_selfanimexist
	OC_alive
	OC_life
	OC_lifemax
	OC_power
	OC_powermax
	OC_canrecover
	OC_roundstate
	OC_ishelper
	OC_numhelper
	OC_numexplod
	OC_numprojid
	OC_numproj
	OC_teammode
	OC_teamside
	OC_hitdefattr
	OC_inguarddist
	OC_movecontact
	OC_movehit
	OC_moveguarded
	OC_movereversed
	OC_projcontacttime
	OC_projhittime
	OC_projguardedtime
	OC_projcanceltime
	OC_backedge
	OC_backedgedist
	OC_backedgebodydist
	OC_frontedge
	OC_frontedgedist
	OC_frontedgebodydist
	OC_leftedge
	OC_rightedge
	OC_topedge
	OC_bottomedge
	OC_camerapos_x
	OC_camerapos_y
	OC_camerazoom
	OC_gamewidth
	OC_gameheight
	OC_screenwidth
	OC_screenheight
	OC_stateno
	OC_prevstateno
	OC_id
	OC_playeridexist
	OC_gametime
	OC_numtarget
	OC_numenemy
	OC_numpartner
	OC_ailevel
	OC_palno
	OC_hitcount
	OC_uniqhitcount
	OC_hitpausetime
	OC_hitover
	OC_hitshakeover
	OC_hitfall
	OC_hitvel_x
	OC_hitvel_y
	OC_roundsexisted
	OC_parent
	OC_root
	OC_helper
	OC_target
	OC_partner
	OC_enemy
	OC_enemynear
	OC_playerid
	OC_p2
	OC_rdreset
	OC_const_
	OC_st_
	OC_ex_
	OC_var0     = 0
	OC_sysvar0  = 60
	OC_fvar0    = 65
	OC_sysfvar0 = 105
)
const (
	OC_const_data_life OpCode = iota
	OC_const_data_power
	OC_const_data_attack
	OC_const_data_defence
	OC_const_data_fall_defence_mul
	OC_const_data_liedown_time
	OC_const_data_airjuggle
	OC_const_data_sparkno
	OC_const_data_guard_sparkno
	OC_const_data_ko_echo
	OC_const_data_intpersistindex
	OC_const_data_floatpersistindex
	OC_const_size_xscale
	OC_const_size_yscale
	OC_const_size_ground_back
	OC_const_size_ground_front
	OC_const_size_air_back
	OC_const_size_air_front
	OC_const_size_z_width
	OC_const_size_height
	OC_const_size_attack_dist
	OC_const_size_attack_z_width_back
	OC_const_size_attack_z_width_front
	OC_const_size_proj_attack_dist
	OC_const_size_proj_doscale
	OC_const_size_head_pos_x
	OC_const_size_head_pos_y
	OC_const_size_mid_pos_x
	OC_const_size_mid_pos_y
	OC_const_size_shadowoffset
	OC_const_size_draw_offset_x
	OC_const_size_draw_offset_y
	OC_const_velocity_walk_fwd_x
	OC_const_velocity_walk_back_x
	OC_const_velocity_walk_up_x
	OC_const_velocity_walk_down_x
	OC_const_velocity_run_fwd_x
	OC_const_velocity_run_fwd_y
	OC_const_velocity_run_back_x
	OC_const_velocity_run_back_y
	OC_const_velocity_run_up_x
	OC_const_velocity_run_up_y
	OC_const_velocity_run_down_x
	OC_const_velocity_run_down_y
	OC_const_velocity_jump_y
	OC_const_velocity_jump_neu_x
	OC_const_velocity_jump_back_x
	OC_const_velocity_jump_fwd_x
	OC_const_velocity_jump_up_x
	OC_const_velocity_jump_down_x
	OC_const_velocity_runjump_back_x
	OC_const_velocity_runjump_back_y
	OC_const_velocity_runjump_y
	OC_const_velocity_runjump_fwd_x
	OC_const_velocity_runjump_up_x
	OC_const_velocity_runjump_down_x
	OC_const_velocity_airjump_y
	OC_const_velocity_airjump_neu_x
	OC_const_velocity_airjump_back_x
	OC_const_velocity_airjump_fwd_x
	OC_const_velocity_airjump_up_x
	OC_const_velocity_airjump_down_x
	OC_const_velocity_air_gethit_groundrecover_x
	OC_const_velocity_air_gethit_groundrecover_y
	OC_const_velocity_air_gethit_airrecover_mul_x
	OC_const_velocity_air_gethit_airrecover_mul_y
	OC_const_velocity_air_gethit_airrecover_add_x
	OC_const_velocity_air_gethit_airrecover_add_y
	OC_const_velocity_air_gethit_airrecover_back
	OC_const_velocity_air_gethit_airrecover_fwd
	OC_const_velocity_air_gethit_airrecover_up
	OC_const_velocity_air_gethit_airrecover_down
	OC_const_movement_airjump_num
	OC_const_movement_airjump_height
	OC_const_movement_yaccel
	OC_const_movement_stand_friction
	OC_const_movement_crouch_friction
	OC_const_movement_stand_friction_threshold
	OC_const_movement_crouch_friction_threshold
	OC_const_movement_air_gethit_groundlevel
	OC_const_movement_air_gethit_groundrecover_ground_threshold
	OC_const_movement_air_gethit_groundrecover_groundlevel
	OC_const_movement_air_gethit_airrecover_threshold
	OC_const_movement_air_gethit_airrecover_yaccel
	OC_const_movement_air_gethit_trip_groundlevel
	OC_const_movement_down_bounce_offset_x
	OC_const_movement_down_bounce_offset_y
	OC_const_movement_down_bounce_yaccel
	OC_const_movement_down_bounce_groundlevel
	OC_const_movement_down_friction_threshold
	OC_const_name
	OC_const_p2name
	OC_const_p3name
	OC_const_p4name
	OC_const_authorname
	OC_const_stagevar_info_author
	OC_const_stagevar_info_displayname
	OC_const_stagevar_info_name
)
const (
	OC_st_var OpCode = iota + OC_var*2
	OC_st_sysvar
	OC_st_fvar
	OC_st_sysfvar
	OC_st_varadd
	OC_st_sysvaradd
	OC_st_fvaradd
	OC_st_sysfvaradd
	OC_st_var0        = OC_var0
	OC_st_sysvar0     = OC_sysvar0
	OC_st_fvar0       = OC_fvar0
	OC_st_sysfvar0    = OC_sysfvar0
	OC_st_var0add     = OC_var + OC_var0
	OC_st_sysvar0add  = OC_var + OC_sysvar0
	OC_st_fvar0add    = OC_var + OC_fvar0
	OC_st_sysfvar0add = OC_var + OC_sysfvar0
)
const (
	OC_ex_p2dist_x OpCode = iota
	OC_ex_p2dist_y
	OC_ex_p2bodydist_x
	OC_ex_parentdist_x
	OC_ex_parentdist_y
	OC_ex_rootdist_x
	OC_ex_rootdist_y
	OC_ex_win
	OC_ex_winko
	OC_ex_wintime
	OC_ex_winperfect
	OC_ex_lose
	OC_ex_loseko
	OC_ex_losetime
	OC_ex_drawgame
	OC_ex_matchover
	OC_ex_matchno
	OC_ex_roundno
	OC_ex_ishometeam
	OC_ex_tickspersecond
	OC_ex_timemod
	OC_ex_majorversion
	OC_ex_drawpalno
	OC_ex_gethitvar_animtype
	OC_ex_gethitvar_airtype
	OC_ex_gethitvar_groundtype
	OC_ex_gethitvar_damage
	OC_ex_gethitvar_hitcount
	OC_ex_gethitvar_fallcount
	OC_ex_gethitvar_hitshaketime
	OC_ex_gethitvar_hittime
	OC_ex_gethitvar_slidetime
	OC_ex_gethitvar_ctrltime
	OC_ex_gethitvar_recovertime
	OC_ex_gethitvar_xoff
	OC_ex_gethitvar_yoff
	OC_ex_gethitvar_xvel
	OC_ex_gethitvar_yvel
	OC_ex_gethitvar_yaccel
	OC_ex_gethitvar_chainid
	OC_ex_gethitvar_guarded
	OC_ex_gethitvar_isbound
	OC_ex_gethitvar_fall
	OC_ex_gethitvar_fall_damage
	OC_ex_gethitvar_fall_xvel
	OC_ex_gethitvar_fall_yvel
	OC_ex_gethitvar_fall_recover
	OC_ex_gethitvar_fall_time
	OC_ex_gethitvar_fall_recovertime
	OC_ex_gethitvar_fall_kill
	OC_ex_gethitvar_fall_envshake_time
	OC_ex_gethitvar_fall_envshake_freq
	OC_ex_gethitvar_fall_envshake_ampl
	OC_ex_gethitvar_fall_envshake_phase
)
const (
	NumVar     = OC_sysvar0 - OC_var0
	NumSysVar  = OC_fvar0 - OC_sysvar0
	NumFvar    = OC_sysfvar0 - OC_fvar0
	NumSysFvar = OC_var - OC_sysfvar0
)

type StringPool struct {
	List []string
	Map  map[string]int
}

func NewStringPool() *StringPool {
	return &StringPool{Map: make(map[string]int)}
}
func (sp *StringPool) Clear() {
	sp.List, sp.Map = nil, make(map[string]int)
}
func (sp *StringPool) Add(s string) int {
	i, ok := sp.Map[s]
	if !ok {
		i = len(sp.List)
		sp.List = append(sp.List, s)
		sp.Map[s] = i
	}
	return i
}

type BytecodeValue struct {
	t ValueType
	v float64
}

func (bv BytecodeValue) IsNone() bool { return bv.t == VT_None }
func (bv BytecodeValue) IsSF() bool   { return bv.t == VT_SFalse }
func (bv BytecodeValue) ToF() float32 {
	if bv.IsSF() {
		return 0
	}
	return float32(bv.v)
}
func (bv BytecodeValue) ToI() int32 {
	if bv.IsSF() {
		return 0
	}
	return int32(bv.v)
}
func (bv BytecodeValue) ToB() bool {
	if bv.IsSF() || bv.v == 0 {
		return false
	}
	return true
}
func (bv *BytecodeValue) SetF(f float32) {
	if math.IsNaN(float64(f)) {
		*bv = BytecodeSF()
	} else {
		*bv = BytecodeValue{VT_Float, float64(f)}
	}
}
func (bv *BytecodeValue) SetI(i int32) {
	*bv = BytecodeValue{VT_Int, float64(i)}
}
func (bv *BytecodeValue) SetB(b bool) {
	bv.t = VT_Bool
	if b {
		bv.v = 1
	} else {
		bv.v = 0
	}
}

func bvNone() BytecodeValue {
	return BytecodeValue{VT_None, 0}
}
func BytecodeSF() BytecodeValue {
	return BytecodeValue{VT_SFalse, math.NaN()}
}
func BytecodeFloat(f float32) BytecodeValue {
	return BytecodeValue{VT_Float, float64(f)}
}
func BytecodeInt(i int32) BytecodeValue {
	return BytecodeValue{VT_Int, float64(i)}
}
func BytecodeBool(b bool) BytecodeValue {
	return BytecodeValue{VT_Bool, float64(Btoi(b))}
}

type BytecodeStack []BytecodeValue

func (bs *BytecodeStack) Clear()                { *bs = (*bs)[:0] }
func (bs *BytecodeStack) Push(bv BytecodeValue) { *bs = append(*bs, bv) }
func (bs *BytecodeStack) PushI(i int32)         { bs.Push(BytecodeInt(i)) }
func (bs *BytecodeStack) PushF(f float32)       { bs.Push(BytecodeFloat(f)) }
func (bs *BytecodeStack) PushB(b bool)          { bs.Push(BytecodeBool(b)) }
func (bs BytecodeStack) Top() *BytecodeValue {
	return &bs[len(bs)-1]
}
func (bs *BytecodeStack) Pop() (bv BytecodeValue) {
	bv, *bs = *bs.Top(), (*bs)[:len(*bs)-1]
	return
}
func (bs *BytecodeStack) Dup() {
	bs.Push(*bs.Top())
}
func (bs *BytecodeStack) Swap() {
	*bs.Top(), (*bs)[len(*bs)-2] = (*bs)[len(*bs)-2], *bs.Top()
}
func (bs *BytecodeStack) Alloc(size int) []BytecodeValue {
	if len(*bs)+size > cap(*bs) {
		tmp := *bs
		*bs = make(BytecodeStack, len(*bs)+size)
		copy(*bs, tmp)
	} else {
		*bs = (*bs)[:len(*bs)+size]
		for i := len(*bs) - size; i < len(*bs); i++ {
			(*bs)[i] = bvNone()
		}
	}
	return (*bs)[len(*bs)-size:]
}

type BytecodeExp []OpCode

func (be *BytecodeExp) append(op ...OpCode) {
	*be = append(*be, op...)
}
func (be *BytecodeExp) appendValue(bv BytecodeValue) (ok bool) {
	switch bv.t {
	case VT_Float:
		be.append(OC_float)
		f := float32(bv.v)
		be.append((*(*[4]OpCode)(unsafe.Pointer(&f)))[:]...)
	case VT_Int:
		if bv.v >= -128 && bv.v <= 127 {
			be.append(OC_int8, OpCode(bv.v))
		} else {
			be.append(OC_int)
			i := int32(bv.v)
			be.append((*(*[4]OpCode)(unsafe.Pointer(&i)))[:]...)
		}
	case VT_Bool:
		if bv.v != 0 {
			be.append(OC_int8, 1)
		} else {
			be.append(OC_int8, 0)
		}
	case VT_SFalse:
		be.append(OC_int8, 0)
	default:
		return false
	}
	return true
}
func (be *BytecodeExp) appendI32Op(op OpCode, addr int32) {
	be.append(op)
	be.append((*(*[4]OpCode)(unsafe.Pointer(&addr)))[:]...)
}
func (_ BytecodeExp) neg(v *BytecodeValue) {
	if v.t == VT_Float {
		v.v *= -1
	} else {
		v.SetI(-v.ToI())
	}
}
func (_ BytecodeExp) not(v *BytecodeValue) {
	v.SetI(^v.ToI())
}
func (_ BytecodeExp) blnot(v *BytecodeValue) {
	v.SetB(!v.ToB())
}
func (_ BytecodeExp) pow(v1 *BytecodeValue, v2 BytecodeValue, pn int) {
	if ValueType(Min(int32(v1.t), int32(v2.t))) == VT_Float {
		v1.SetF(Pow(v1.ToF(), v2.ToF()))
	} else if v2.ToF() < 0 {
		v1.SetF(Pow(v1.ToF(), v2.ToF()))
	} else {
		i1, i2, hb := v1.ToI(), v2.ToI(), int32(-1)
		for uint32(i2)>>uint(hb+1) != 0 {
			hb++
		}
		var i, bit, tmp int32 = 1, 0, i1
		for ; bit <= hb; bit++ {
			var shift uint
			if bit == hb || sys.cgi[pn].ver[0] == 1 {
				shift = uint(bit)
			} else {
				shift = uint((hb - 1) - bit)
			}
			if i2&(1<<shift) != 0 {
				i *= tmp
			}
			tmp *= tmp
		}
		v1.SetI(i)
	}
}
func (_ BytecodeExp) mul(v1 *BytecodeValue, v2 BytecodeValue) {
	if ValueType(Min(int32(v1.t), int32(v2.t))) == VT_Float {
		v1.SetF(v1.ToF() * v2.ToF())
	} else {
		v1.SetI(v1.ToI() * v2.ToI())
	}
}
func (_ BytecodeExp) div(v1 *BytecodeValue, v2 BytecodeValue) {
	if ValueType(Min(int32(v1.t), int32(v2.t))) == VT_Float {
		v1.SetF(v1.ToF() / v2.ToF())
	} else if v2.ToI() == 0 {
		*v1 = BytecodeSF()
	} else {
		v1.SetI(v1.ToI() / v2.ToI())
	}
}
func (_ BytecodeExp) mod(v1 *BytecodeValue, v2 BytecodeValue) {
	if v2.ToI() == 0 {
		*v1 = BytecodeSF()
	} else {
		v1.SetI(v1.ToI() % v2.ToI())
	}
}
func (_ BytecodeExp) add(v1 *BytecodeValue, v2 BytecodeValue) {
	if ValueType(Min(int32(v1.t), int32(v2.t))) == VT_Float {
		v1.SetF(v1.ToF() + v2.ToF())
	} else {
		v1.SetI(v1.ToI() + v2.ToI())
	}
}
func (_ BytecodeExp) sub(v1 *BytecodeValue, v2 BytecodeValue) {
	if ValueType(Min(int32(v1.t), int32(v2.t))) == VT_Float {
		v1.SetF(v1.ToF() - v2.ToF())
	} else {
		v1.SetI(v1.ToI() - v2.ToI())
	}
}
func (_ BytecodeExp) gt(v1 *BytecodeValue, v2 BytecodeValue) {
	if ValueType(Min(int32(v1.t), int32(v2.t))) == VT_Float {
		v1.SetB(v1.ToF() > v2.ToF())
	} else {
		v1.SetB(v1.ToI() > v2.ToI())
	}
}
func (_ BytecodeExp) ge(v1 *BytecodeValue, v2 BytecodeValue) {
	if ValueType(Min(int32(v1.t), int32(v2.t))) == VT_Float {
		v1.SetB(v1.ToF() >= v2.ToF())
	} else {
		v1.SetB(v1.ToI() >= v2.ToI())
	}
}
func (_ BytecodeExp) lt(v1 *BytecodeValue, v2 BytecodeValue) {
	if ValueType(Min(int32(v1.t), int32(v2.t))) == VT_Float {
		v1.SetB(v1.ToF() < v2.ToF())
	} else {
		v1.SetB(v1.ToI() < v2.ToI())
	}
}
func (_ BytecodeExp) le(v1 *BytecodeValue, v2 BytecodeValue) {
	if ValueType(Min(int32(v1.t), int32(v2.t))) == VT_Float {
		v1.SetB(v1.ToF() <= v2.ToF())
	} else {
		v1.SetB(v1.ToI() <= v2.ToI())
	}
}
func (_ BytecodeExp) eq(v1 *BytecodeValue, v2 BytecodeValue) {
	if ValueType(Min(int32(v1.t), int32(v2.t))) == VT_Float {
		v1.SetB(v1.ToF() == v2.ToF())
	} else {
		v1.SetB(v1.ToI() == v2.ToI())
	}
}
func (_ BytecodeExp) ne(v1 *BytecodeValue, v2 BytecodeValue) {
	if ValueType(Min(int32(v1.t), int32(v2.t))) == VT_Float {
		v1.SetB(v1.ToF() != v2.ToF())
	} else {
		v1.SetB(v1.ToI() != v2.ToI())
	}
}
func (_ BytecodeExp) and(v1 *BytecodeValue, v2 BytecodeValue) {
	v1.SetI(v1.ToI() & v2.ToI())
}
func (_ BytecodeExp) xor(v1 *BytecodeValue, v2 BytecodeValue) {
	v1.SetI(v1.ToI() ^ v2.ToI())
}
func (_ BytecodeExp) or(v1 *BytecodeValue, v2 BytecodeValue) {
	v1.SetI(v1.ToI() | v2.ToI())
}
func (_ BytecodeExp) bland(v1 *BytecodeValue, v2 BytecodeValue) {
	v1.SetB(v1.ToB() && v2.ToB())
}
func (_ BytecodeExp) blxor(v1 *BytecodeValue, v2 BytecodeValue) {
	v1.SetB(v1.ToB() != v2.ToB())
}
func (_ BytecodeExp) blor(v1 *BytecodeValue, v2 BytecodeValue) {
	v1.SetB(v1.ToB() || v2.ToB())
}
func (_ BytecodeExp) abs(v1 *BytecodeValue) {
	if v1.t == VT_Float {
		v1.v = math.Abs(v1.v)
	} else {
		v1.SetI(Abs(v1.ToI()))
	}
}
func (_ BytecodeExp) exp(v1 *BytecodeValue) {
	v1.SetF(float32(math.Exp(v1.v)))
}
func (_ BytecodeExp) ln(v1 *BytecodeValue) {
	if v1.v <= 0 {
		*v1 = BytecodeSF()
	} else {
		v1.SetF(float32(math.Log(v1.v)))
	}
}
func (_ BytecodeExp) log(v1 *BytecodeValue, v2 BytecodeValue) {
	if v1.v <= 0 || v2.v <= 0 {
		*v1 = BytecodeSF()
	} else {
		v1.SetF(float32(math.Log(v1.v) / math.Log(v2.v)))
	}
}
func (_ BytecodeExp) cos(v1 *BytecodeValue) {
	v1.SetF(float32(math.Cos(v1.v)))
}
func (_ BytecodeExp) sin(v1 *BytecodeValue) {
	v1.SetF(float32(math.Sin(v1.v)))
}
func (_ BytecodeExp) tan(v1 *BytecodeValue) {
	v1.SetF(float32(math.Tan(v1.v)))
}
func (_ BytecodeExp) acos(v1 *BytecodeValue) {
	v1.SetF(float32(math.Acos(v1.v)))
}
func (_ BytecodeExp) asin(v1 *BytecodeValue) {
	v1.SetF(float32(math.Asin(v1.v)))
}
func (_ BytecodeExp) atan(v1 *BytecodeValue) {
	v1.SetF(float32(math.Atan(v1.v)))
}
func (_ BytecodeExp) floor(v1 *BytecodeValue) {
	if v1.t == VT_Float {
		f := math.Floor(v1.v)
		if math.IsNaN(f) {
			*v1 = BytecodeSF()
		} else {
			v1.SetI(int32(f))
		}
	}
}
func (_ BytecodeExp) ceil(v1 *BytecodeValue) {
	if v1.t == VT_Float {
		f := math.Ceil(v1.v)
		if math.IsNaN(f) {
			*v1 = BytecodeSF()
		} else {
			v1.SetI(int32(f))
		}
	}
}
func (be BytecodeExp) run(c *Char) BytecodeValue {
	oc := c
	for i := 1; i <= len(be); i++ {
		switch be[i-1] {
		case OC_jsf8:
			if sys.bcStack.Top().IsSF() {
				if be[i] == 0 {
					i = len(be)
				} else {
					i += int(uint8(be[i])) + 1
				}
			} else {
				i++
			}
		case OC_jz8, OC_jnz8:
			if sys.bcStack.Top().ToB() == (be[i-1] == OC_jz8) {
				i++
				break
			}
			fallthrough
		case OC_jmp8:
			if be[i] == 0 {
				i = len(be)
			} else {
				i += int(uint8(be[i])) + 1
			}
		case OC_jz, OC_jnz:
			if sys.bcStack.Top().ToB() == (be[i-1] == OC_jz) {
				i += 4
				break
			}
			fallthrough
		case OC_jmp:
			i += int(*(*int32)(unsafe.Pointer(&be[i]))) + 4
		case OC_parent:
			if c = c.parent(); c != nil {
				i += 4
				continue
			}
			sys.bcStack.Push(BytecodeSF())
			i += int(*(*int32)(unsafe.Pointer(&be[i]))) + 4
		case OC_root:
			if c = c.root(); c != nil {
				i += 4
				continue
			}
			sys.bcStack.Push(BytecodeSF())
			i += int(*(*int32)(unsafe.Pointer(&be[i]))) + 4
		case OC_helper:
			if c = c.helper(sys.bcStack.Pop().ToI()); c != nil {
				i += 4
				continue
			}
			sys.bcStack.Push(BytecodeSF())
			i += int(*(*int32)(unsafe.Pointer(&be[i]))) + 4
		case OC_target:
			if c = c.target(sys.bcStack.Pop().ToI()); c != nil {
				i += 4
				continue
			}
			sys.bcStack.Push(BytecodeSF())
			i += int(*(*int32)(unsafe.Pointer(&be[i]))) + 4
		case OC_partner:
			if c = c.partner(sys.bcStack.Pop().ToI()); c != nil {
				i += 4
				continue
			}
			sys.bcStack.Push(BytecodeSF())
			i += int(*(*int32)(unsafe.Pointer(&be[i]))) + 4
		case OC_enemy:
			if c = c.enemy(sys.bcStack.Pop().ToI()); c != nil {
				i += 4
				continue
			}
			sys.bcStack.Push(BytecodeSF())
			i += int(*(*int32)(unsafe.Pointer(&be[i]))) + 4
		case OC_enemynear:
			if c = c.enemyNear(sys.bcStack.Pop().ToI()); c != nil {
				i += 4
				continue
			}
			sys.bcStack.Push(BytecodeSF())
			i += int(*(*int32)(unsafe.Pointer(&be[i]))) + 4
		case OC_playerid:
			if c = sys.playerID(sys.bcStack.Pop().ToI()); c != nil {
				i += 4
				continue
			}
			sys.bcStack.Push(BytecodeSF())
			i += int(*(*int32)(unsafe.Pointer(&be[i]))) + 4
		case OC_p2:
			if c = c.p2(); c != nil {
				i += 4
				continue
			}
			sys.bcStack.Push(BytecodeSF())
			i += int(*(*int32)(unsafe.Pointer(&be[i]))) + 4
		case OC_rdreset:
			// NOP
		case OC_run:
			l := int(*(*int32)(unsafe.Pointer(&be[i])))
			sys.bcStack.Push(be[i+4 : i+4+l].run(c))
			i += 4 + l
		case OC_nordrun:
			l := int(*(*int32)(unsafe.Pointer(&be[i])))
			sys.bcStack.Push(be[i+4 : i+4+l].run(oc))
			i += 4 + l
			continue
		case OC_int8:
			sys.bcStack.PushI(int32(int8(be[i])))
			i++
		case OC_int:
			sys.bcStack.PushI(*(*int32)(unsafe.Pointer(&be[i])))
			i += 4
		case OC_float:
			sys.bcStack.PushF(*(*float32)(unsafe.Pointer(&be[i])))
			i += 4
		case OC_neg:
			be.neg(sys.bcStack.Top())
		case OC_not:
			be.not(sys.bcStack.Top())
		case OC_blnot:
			be.blnot(sys.bcStack.Top())
		case OC_pow:
			v2 := sys.bcStack.Pop()
			be.pow(sys.bcStack.Top(), v2, sys.workingChar.ss.sb.playerNo)
		case OC_mul:
			v2 := sys.bcStack.Pop()
			be.mul(sys.bcStack.Top(), v2)
		case OC_div:
			v2 := sys.bcStack.Pop()
			be.div(sys.bcStack.Top(), v2)
		case OC_mod:
			v2 := sys.bcStack.Pop()
			be.mod(sys.bcStack.Top(), v2)
		case OC_add:
			v2 := sys.bcStack.Pop()
			be.add(sys.bcStack.Top(), v2)
		case OC_sub:
			v2 := sys.bcStack.Pop()
			be.sub(sys.bcStack.Top(), v2)
		case OC_gt:
			v2 := sys.bcStack.Pop()
			be.gt(sys.bcStack.Top(), v2)
		case OC_ge:
			v2 := sys.bcStack.Pop()
			be.ge(sys.bcStack.Top(), v2)
		case OC_lt:
			v2 := sys.bcStack.Pop()
			be.lt(sys.bcStack.Top(), v2)
		case OC_le:
			v2 := sys.bcStack.Pop()
			be.le(sys.bcStack.Top(), v2)
		case OC_eq:
			v2 := sys.bcStack.Pop()
			be.eq(sys.bcStack.Top(), v2)
		case OC_ne:
			v2 := sys.bcStack.Pop()
			be.ne(sys.bcStack.Top(), v2)
		case OC_and:
			v2 := sys.bcStack.Pop()
			be.and(sys.bcStack.Top(), v2)
		case OC_xor:
			v2 := sys.bcStack.Pop()
			be.xor(sys.bcStack.Top(), v2)
		case OC_or:
			v2 := sys.bcStack.Pop()
			be.or(sys.bcStack.Top(), v2)
		case OC_bland:
			v2 := sys.bcStack.Pop()
			be.bland(sys.bcStack.Top(), v2)
		case OC_blxor:
			v2 := sys.bcStack.Pop()
			be.blxor(sys.bcStack.Top(), v2)
		case OC_blor:
			v2 := sys.bcStack.Pop()
			be.blor(sys.bcStack.Top(), v2)
		case OC_abs:
			be.abs(sys.bcStack.Top())
		case OC_exp:
			be.exp(sys.bcStack.Top())
		case OC_ln:
			be.ln(sys.bcStack.Top())
		case OC_log:
			v2 := sys.bcStack.Pop()
			be.log(sys.bcStack.Top(), v2)
		case OC_cos:
			be.cos(sys.bcStack.Top())
		case OC_sin:
			be.sin(sys.bcStack.Top())
		case OC_tan:
			be.tan(sys.bcStack.Top())
		case OC_acos:
			be.acos(sys.bcStack.Top())
		case OC_asin:
			be.asin(sys.bcStack.Top())
		case OC_atan:
			be.atan(sys.bcStack.Top())
		case OC_floor:
			be.floor(sys.bcStack.Top())
		case OC_ceil:
			be.ceil(sys.bcStack.Top())
		case OC_ifelse:
			v3 := sys.bcStack.Pop()
			v2 := sys.bcStack.Pop()
			if sys.bcStack.Top().ToB() {
				*sys.bcStack.Top() = v2
			} else {
				*sys.bcStack.Top() = v3
			}
		case OC_pop:
			sys.bcStack.Pop()
		case OC_dup:
			sys.bcStack.Dup()
		case OC_swap:
			sys.bcStack.Swap()
		case OC_ailevel:
			sys.bcStack.PushI(c.aiLevel())
		case OC_alive:
			sys.bcStack.PushB(c.alive())
		case OC_anim:
			sys.bcStack.PushI(c.animNo)
		case OC_animelemno:
			*sys.bcStack.Top() = c.animElemNo(sys.bcStack.Top().ToI())
		case OC_animelemtime:
			*sys.bcStack.Top() = c.animElemTime(sys.bcStack.Top().ToI())
		case OC_animexist:
			*sys.bcStack.Top() = c.animExist(sys.workingChar, *sys.bcStack.Top())
		case OC_animtime:
			sys.bcStack.PushI(c.animTime())
		case OC_backedge:
			sys.bcStack.PushF(c.backEdge())
		case OC_backedgebodydist:
			sys.bcStack.PushI(int32(c.backEdgeBodyDist()))
		case OC_backedgedist:
			sys.bcStack.PushI(int32(c.backEdgeDist()))
		case OC_bottomedge:
			sys.bcStack.PushF(c.bottomEdge())
		case OC_camerapos_x:
			sys.bcStack.PushF(sys.cam.Pos[0] / oc.localscl)
		case OC_camerapos_y:
			sys.bcStack.PushF(sys.cam.Pos[1] / oc.localscl)
		case OC_camerazoom:
			sys.bcStack.PushF(sys.cam.Scale)
		case OC_canrecover:
			sys.bcStack.PushB(c.canRecover())
		case OC_command:
			sys.bcStack.PushB(c.command(sys.workingState.playerNo,
				int(*(*int32)(unsafe.Pointer(&be[i])))))
			i += 4
		case OC_ctrl:
			sys.bcStack.PushB(c.ctrl())
		case OC_facing:
			sys.bcStack.PushI(int32(c.facing))
		case OC_frontedge:
			sys.bcStack.PushF(c.frontEdge())
		case OC_frontedgebodydist:
			sys.bcStack.PushI(int32(c.frontEdgeBodyDist()))
		case OC_frontedgedist:
			sys.bcStack.PushI(int32(c.frontEdgeDist()))
		case OC_gameheight:
			sys.bcStack.PushF(c.gameHeight())
		case OC_gametime:
			sys.bcStack.PushI(sys.gameTime)
		case OC_gamewidth:
			sys.bcStack.PushF(c.gameWidth())
		case OC_hitcount:
			sys.bcStack.PushI(c.hitCount)
		case OC_hitdefattr:
			sys.bcStack.PushB(c.hitDefAttr(*(*int32)(unsafe.Pointer(&be[i]))))
			i += 4
		case OC_hitfall:
			sys.bcStack.PushB(c.ghv.fallf)
		case OC_hitover:
			sys.bcStack.PushB(c.hitOver())
		case OC_hitpausetime:
			sys.bcStack.PushI(c.hitPauseTime)
		case OC_hitshakeover:
			sys.bcStack.PushB(c.hitShakeOver())
		case OC_hitvel_x:
			sys.bcStack.PushF(c.hitVelX() * c.localscl / oc.localscl)
		case OC_hitvel_y:
			sys.bcStack.PushF(c.hitVelY() * c.localscl / oc.localscl)
		case OC_id:
			sys.bcStack.PushI(c.id)
		case OC_inguarddist:
			sys.bcStack.PushB(c.inguarddist)
		case OC_ishelper:
			*sys.bcStack.Top() = c.isHelper(*sys.bcStack.Top())
		case OC_leftedge:
			sys.bcStack.PushF(c.leftEdge())
		case OC_life:
			sys.bcStack.PushI(c.life)
		case OC_lifemax:
			sys.bcStack.PushI(c.lifeMax)
		case OC_movecontact:
			sys.bcStack.PushI(c.moveContact())
		case OC_moveguarded:
			sys.bcStack.PushI(c.moveGuarded())
		case OC_movehit:
			sys.bcStack.PushI(c.moveHit())
		case OC_movereversed:
			sys.bcStack.PushI(c.moveReversed())
		case OC_movetype:
			sys.bcStack.PushB(c.ss.moveType == MoveType(be[i])<<15)
			i++
		case OC_numenemy:
			sys.bcStack.PushI(c.numEnemy())
		case OC_numexplod:
			*sys.bcStack.Top() = c.numExplod(*sys.bcStack.Top())
		case OC_numhelper:
			*sys.bcStack.Top() = c.numHelper(*sys.bcStack.Top())
		case OC_numpartner:
			sys.bcStack.PushI(c.numPartner())
		case OC_numproj:
			sys.bcStack.PushI(c.numProj())
		case OC_numprojid:
			*sys.bcStack.Top() = c.numProjID(*sys.bcStack.Top())
		case OC_numtarget:
			*sys.bcStack.Top() = c.numTarget(*sys.bcStack.Top())
		case OC_palno:
			sys.bcStack.PushI(c.palno())
		case OC_pos_x:
			sys.bcStack.PushF((c.pos[0]*c.localscl/oc.localscl - sys.cam.Pos[0]/oc.localscl))
		case OC_pos_y:
			sys.bcStack.PushF(c.pos[1] * c.localscl / oc.localscl)
		case OC_power:
			sys.bcStack.PushI(c.getPower())
		case OC_powermax:
			sys.bcStack.PushI(c.powerMax)
		case OC_playeridexist:
			*sys.bcStack.Top() = sys.playerIDExist(*sys.bcStack.Top())
		case OC_prevstateno:
			sys.bcStack.PushI(c.ss.prevno)
		case OC_projcanceltime:
			*sys.bcStack.Top() = c.projCancelTime(*sys.bcStack.Top())
		case OC_projcontacttime:
			*sys.bcStack.Top() = c.projContactTime(*sys.bcStack.Top())
		case OC_projguardedtime:
			*sys.bcStack.Top() = c.projGuardedTime(*sys.bcStack.Top())
		case OC_projhittime:
			*sys.bcStack.Top() = c.projHitTime(*sys.bcStack.Top())
		case OC_random:
			sys.bcStack.PushI(Rand(0, 999))
		case OC_rightedge:
			sys.bcStack.PushF(c.rightEdge())
		case OC_roundsexisted:
			sys.bcStack.PushI(c.roundsExisted())
		case OC_roundstate:
			sys.bcStack.PushI(c.roundState())
		case OC_screenheight:
			sys.bcStack.PushF(sys.screenHeight() / oc.localscl)
		case OC_screenpos_x:
			sys.bcStack.PushF((c.screenPosX()) / oc.localscl)
		case OC_screenpos_y:
			sys.bcStack.PushF((c.screenPosY()) / oc.localscl)
		case OC_screenwidth:
			sys.bcStack.PushF(sys.screenWidth() / oc.localscl)
		case OC_selfanimexist:
			*sys.bcStack.Top() = c.selfAnimExist(*sys.bcStack.Top())
		case OC_stateno:
			sys.bcStack.PushI(c.ss.no)
		case OC_statetype:
			sys.bcStack.PushB(c.ss.stateType == StateType(be[i]))
			i++
		case OC_teammode:
			sys.bcStack.PushB(sys.tmode[c.playerNo&1] == TeamMode(be[i]))
			i++
		case OC_teamside:
			sys.bcStack.PushI(int32(c.playerNo)&1 + 1)
		case OC_time:
			sys.bcStack.PushI(c.time())
		case OC_topedge:
			sys.bcStack.PushF(c.topEdge())
		case OC_uniqhitcount:
			sys.bcStack.PushI(c.uniqHitCount)
		case OC_vel_x:
			sys.bcStack.PushF(c.vel[0] * c.localscl / oc.localscl)
		case OC_vel_y:
			sys.bcStack.PushF(c.vel[1] * c.localscl / oc.localscl)
		case OC_st_:
			be.run_st(c, &i)
		case OC_const_:
			be.run_const(c, &i, oc)
		case OC_ex_:
			be.run_ex(c, &i, oc)
		case OC_var:
			*sys.bcStack.Top() = c.varGet(sys.bcStack.Top().ToI())
		case OC_sysvar:
			*sys.bcStack.Top() = c.sysVarGet(sys.bcStack.Top().ToI())
		case OC_fvar:
			*sys.bcStack.Top() = c.fvarGet(sys.bcStack.Top().ToI())
		case OC_sysfvar:
			*sys.bcStack.Top() = c.sysFvarGet(sys.bcStack.Top().ToI())
		case OC_localvar:
			sys.bcStack.Push(sys.bcVar[uint8(be[i])])
			i++
		default:
			vi := be[i-1]
			if vi < OC_sysvar0+NumSysVar {
				sys.bcStack.PushI(c.ivar[vi-OC_var0])
			} else {
				sys.bcStack.PushF(c.fvar[vi-OC_fvar0])
			}
		}
		c = oc
	}
	return sys.bcStack.Pop()
}
func (be BytecodeExp) run_st(c *Char, i *int) {
	(*i)++
	switch be[*i-1] {
	case OC_st_var:
		v := sys.bcStack.Pop().ToI()
		*sys.bcStack.Top() = c.varSet(sys.bcStack.Top().ToI(), v)
	case OC_st_sysvar:
		v := sys.bcStack.Pop().ToI()
		*sys.bcStack.Top() = c.sysVarSet(sys.bcStack.Top().ToI(), v)
	case OC_st_fvar:
		v := sys.bcStack.Pop().ToF()
		*sys.bcStack.Top() = c.fvarSet(sys.bcStack.Top().ToI(), v)
	case OC_st_sysfvar:
		v := sys.bcStack.Pop().ToF()
		*sys.bcStack.Top() = c.sysFvarSet(sys.bcStack.Top().ToI(), v)
	case OC_st_varadd:
		v := sys.bcStack.Pop().ToI()
		*sys.bcStack.Top() = c.varAdd(sys.bcStack.Top().ToI(), v)
	case OC_st_sysvaradd:
		v := sys.bcStack.Pop().ToI()
		*sys.bcStack.Top() = c.sysVarAdd(sys.bcStack.Top().ToI(), v)
	case OC_st_fvaradd:
		v := sys.bcStack.Pop().ToF()
		*sys.bcStack.Top() = c.fvarAdd(sys.bcStack.Top().ToI(), v)
	case OC_st_sysfvaradd:
		v := sys.bcStack.Pop().ToF()
		*sys.bcStack.Top() = c.sysFvarAdd(sys.bcStack.Top().ToI(), v)
	default:
		vi := be[*i-1]
		if vi < OC_st_sysvar0+NumSysVar {
			c.ivar[vi-OC_st_var0] = sys.bcStack.Top().ToI()
			sys.bcStack.Top().SetI(c.ivar[vi-OC_st_var0])
		} else if vi < OC_st_sysfvar0+NumSysFvar {
			c.fvar[vi-OC_st_fvar0] = sys.bcStack.Top().ToF()
			sys.bcStack.Top().SetF(c.fvar[vi-OC_st_fvar0])
		} else if vi < OC_st_sysvar0add+NumSysVar {
			c.ivar[vi-OC_st_var0add] += sys.bcStack.Top().ToI()
			sys.bcStack.Top().SetI(c.ivar[vi-OC_st_var0add])
		} else if vi < OC_st_sysfvar0add+NumSysFvar {
			c.fvar[vi-OC_st_fvar0add] += sys.bcStack.Top().ToF()
			sys.bcStack.Top().SetF(c.fvar[vi-OC_st_fvar0add])
		} else {
			sys.errLog.Printf("%v\n", be[*i-1])
			c.panic()
		}
	}
}
func (be BytecodeExp) run_const(c *Char, i *int, oc *Char) {
	(*i)++
	switch be[*i-1] {
	case OC_const_data_life:
		sys.bcStack.PushI(c.gi().data.life)
	case OC_const_data_power:
		sys.bcStack.PushI(c.gi().data.power)
	case OC_const_data_attack:
		sys.bcStack.PushI(c.gi().data.attack)
	case OC_const_data_defence:
		sys.bcStack.PushI(c.gi().data.defence)
	case OC_const_data_fall_defence_mul:
		sys.bcStack.PushF(c.gi().data.fall.defence_mul)
	case OC_const_data_liedown_time:
		sys.bcStack.PushI(c.gi().data.liedown.time)
	case OC_const_data_airjuggle:
		sys.bcStack.PushI(c.gi().data.airjuggle)
	case OC_const_data_sparkno:
		sys.bcStack.PushI(c.gi().data.sparkno)
	case OC_const_data_guard_sparkno:
		sys.bcStack.PushI(c.gi().data.guard.sparkno)
	case OC_const_data_ko_echo:
		sys.bcStack.PushI(c.gi().data.ko.echo)
	case OC_const_data_intpersistindex:
		sys.bcStack.PushI(c.gi().data.intpersistindex)
	case OC_const_data_floatpersistindex:
		sys.bcStack.PushI(c.gi().data.floatpersistindex)
	case OC_const_size_xscale:
		sys.bcStack.PushF(c.size.xscale)
	case OC_const_size_yscale:
		sys.bcStack.PushF(c.size.yscale)
	case OC_const_size_ground_back:
		sys.bcStack.PushF(c.size.ground.back * (320 / float32(c.localcoord)) / oc.localscl)
	case OC_const_size_ground_front:
		sys.bcStack.PushF(c.size.ground.front * (320 / float32(c.localcoord)) / oc.localscl)
	case OC_const_size_air_back:
		sys.bcStack.PushF(c.size.air.back * (320 / float32(c.localcoord)) / oc.localscl)
	case OC_const_size_air_front:
		sys.bcStack.PushF(c.size.air.front * (320 / float32(c.localcoord)) / oc.localscl)
	case OC_const_size_z_width:
		sys.bcStack.PushF(c.size.z.width * (320 / float32(c.localcoord)) / oc.localscl)
	case OC_const_size_height:
		sys.bcStack.PushF(c.size.height * (320 / float32(c.localcoord)) / oc.localscl)
	case OC_const_size_attack_dist:
		sys.bcStack.PushF(c.size.attack.dist * (320 / float32(c.localcoord)) / oc.localscl)
	case OC_const_size_attack_z_width_back:
		sys.bcStack.PushF(c.size.attack.z.width[1] * (320 / float32(c.localcoord)) / oc.localscl)
	case OC_const_size_attack_z_width_front:
		sys.bcStack.PushF(c.size.attack.z.width[0] * (320 / float32(c.localcoord)) / oc.localscl)
	case OC_const_size_proj_attack_dist:
		sys.bcStack.PushF(c.size.proj.attack.dist * (320 / float32(c.localcoord)) / oc.localscl)
	case OC_const_size_proj_doscale:
		sys.bcStack.PushI(c.size.proj.doscale)
	case OC_const_size_head_pos_x:
		sys.bcStack.PushF(c.size.head.pos[0] * (320 / float32(c.localcoord)) / oc.localscl)
	case OC_const_size_head_pos_y:
		sys.bcStack.PushF(c.size.head.pos[1] * (320 / float32(c.localcoord)) / oc.localscl)
	case OC_const_size_mid_pos_x:
		sys.bcStack.PushF(c.size.mid.pos[0] * (320 / float32(c.localcoord)) / oc.localscl)
	case OC_const_size_mid_pos_y:
		sys.bcStack.PushF(c.size.mid.pos[1] * (320 / float32(c.localcoord)) / oc.localscl)
	case OC_const_size_shadowoffset:
		sys.bcStack.PushF(c.size.shadowoffset * (320 / float32(c.localcoord)) / oc.localscl)
	case OC_const_size_draw_offset_x:
		sys.bcStack.PushF(c.size.draw.offset[0] * (320 / float32(c.localcoord)) / oc.localscl)
	case OC_const_size_draw_offset_y:
		sys.bcStack.PushF(c.size.draw.offset[1] * (320 / float32(c.localcoord)) / oc.localscl)
	case OC_const_velocity_walk_fwd_x:
		sys.bcStack.PushF(c.gi().velocity.walk.fwd * (320 / float32(c.localcoord)) / oc.localscl)
	case OC_const_velocity_walk_back_x:
		sys.bcStack.PushF(c.gi().velocity.walk.back * (320 / float32(c.localcoord)) / oc.localscl)
	case OC_const_velocity_walk_up_x:
		sys.bcStack.PushF(c.gi().velocity.walk.up.x * (320 / float32(c.localcoord)) / oc.localscl)
	case OC_const_velocity_walk_down_x:
		sys.bcStack.PushF(c.gi().velocity.walk.down.x * (320 / float32(c.localcoord)) / oc.localscl)
	case OC_const_velocity_run_fwd_x:
		sys.bcStack.PushF(c.gi().velocity.run.fwd[0] * (320 / float32(c.localcoord)) / oc.localscl)
	case OC_const_velocity_run_fwd_y:
		sys.bcStack.PushF(c.gi().velocity.run.fwd[1] * (320 / float32(c.localcoord)) / oc.localscl)
	case OC_const_velocity_run_back_x:
		sys.bcStack.PushF(c.gi().velocity.run.back[0] * (320 / float32(c.localcoord)) / oc.localscl)
	case OC_const_velocity_run_back_y:
		sys.bcStack.PushF(c.gi().velocity.run.back[1] * (320 / float32(c.localcoord)) / oc.localscl)
	case OC_const_velocity_run_up_x:
		sys.bcStack.PushF(c.gi().velocity.run.up.x * (320 / float32(c.localcoord)) / oc.localscl)
	case OC_const_velocity_run_up_y:
		sys.bcStack.PushF(c.gi().velocity.run.up.y * (320 / float32(c.localcoord)) / oc.localscl)
	case OC_const_velocity_run_down_x:
		sys.bcStack.PushF(c.gi().velocity.run.down.x * (320 / float32(c.localcoord)) / oc.localscl)
	case OC_const_velocity_run_down_y:
		sys.bcStack.PushF(c.gi().velocity.run.down.y * (320 / float32(c.localcoord)) / oc.localscl)
	case OC_const_velocity_jump_y:
		sys.bcStack.PushF(c.gi().velocity.jump.neu[1] * (320 / float32(c.localcoord)) / oc.localscl)
	case OC_const_velocity_jump_neu_x:
		sys.bcStack.PushF(c.gi().velocity.jump.neu[0] * (320 / float32(c.localcoord)) / oc.localscl)
	case OC_const_velocity_jump_back_x:
		sys.bcStack.PushF(c.gi().velocity.jump.back * (320 / float32(c.localcoord)) / oc.localscl)
	case OC_const_velocity_jump_fwd_x:
		sys.bcStack.PushF(c.gi().velocity.jump.fwd * (320 / float32(c.localcoord)) / oc.localscl)
	case OC_const_velocity_jump_up_x:
		sys.bcStack.PushF(c.gi().velocity.jump.up.x * (320 / float32(c.localcoord)) / oc.localscl)
	case OC_const_velocity_jump_down_x:
		sys.bcStack.PushF(c.gi().velocity.jump.down.x * (320 / float32(c.localcoord)) / oc.localscl)
	case OC_const_velocity_runjump_back_x:
		sys.bcStack.PushF(c.gi().velocity.runjump.back[0] * (320 / float32(c.localcoord)) / oc.localscl)
	case OC_const_velocity_runjump_back_y:
		sys.bcStack.PushF(c.gi().velocity.runjump.back[1] * (320 / float32(c.localcoord)) / oc.localscl)
	case OC_const_velocity_runjump_y:
		sys.bcStack.PushF(c.gi().velocity.runjump.fwd[1] * (320 / float32(c.localcoord)) / oc.localscl)
	case OC_const_velocity_runjump_fwd_x:
		sys.bcStack.PushF(c.gi().velocity.runjump.fwd[0] * (320 / float32(c.localcoord)) / oc.localscl)
	case OC_const_velocity_runjump_up_x:
		sys.bcStack.PushF(c.gi().velocity.runjump.up.x * (320 / float32(c.localcoord)) / oc.localscl)
	case OC_const_velocity_runjump_down_x:
		sys.bcStack.PushF(c.gi().velocity.runjump.down.x * (320 / float32(c.localcoord)) / oc.localscl)
	case OC_const_velocity_airjump_y:
		sys.bcStack.PushF(c.gi().velocity.airjump.neu[1] * (320 / float32(c.localcoord)) / oc.localscl)
	case OC_const_velocity_airjump_neu_x:
		sys.bcStack.PushF(c.gi().velocity.airjump.neu[0] * (320 / float32(c.localcoord)) / oc.localscl)
	case OC_const_velocity_airjump_back_x:
		sys.bcStack.PushF(c.gi().velocity.airjump.back * (320 / float32(c.localcoord)) / oc.localscl)
	case OC_const_velocity_airjump_fwd_x:
		sys.bcStack.PushF(c.gi().velocity.airjump.fwd * (320 / float32(c.localcoord)) / oc.localscl)
	case OC_const_velocity_airjump_up_x:
		sys.bcStack.PushF(c.gi().velocity.airjump.up.x * (320 / float32(c.localcoord)) / oc.localscl)
	case OC_const_velocity_airjump_down_x:
		sys.bcStack.PushF(c.gi().velocity.airjump.down.x * (320 / float32(c.localcoord)) / oc.localscl)
	case OC_const_velocity_air_gethit_groundrecover_x:
		sys.bcStack.PushF(c.gi().velocity.air.gethit.groundrecover[0] * (320 / float32(c.localcoord)) / oc.localscl)
	case OC_const_velocity_air_gethit_groundrecover_y:
		sys.bcStack.PushF(c.gi().velocity.air.gethit.groundrecover[1] * (320 / float32(c.localcoord)) / oc.localscl)
	case OC_const_velocity_air_gethit_airrecover_mul_x:
		sys.bcStack.PushF(c.gi().velocity.air.gethit.airrecover.mul[0])
	case OC_const_velocity_air_gethit_airrecover_mul_y:
		sys.bcStack.PushF(c.gi().velocity.air.gethit.airrecover.mul[1])
	case OC_const_velocity_air_gethit_airrecover_add_x:
		sys.bcStack.PushF(c.gi().velocity.air.gethit.airrecover.add[0] * (320 / float32(c.localcoord)) / oc.localscl)
	case OC_const_velocity_air_gethit_airrecover_add_y:
		sys.bcStack.PushF(c.gi().velocity.air.gethit.airrecover.add[1] * (320 / float32(c.localcoord)) / oc.localscl)
	case OC_const_velocity_air_gethit_airrecover_back:
		sys.bcStack.PushF(c.gi().velocity.air.gethit.airrecover.back * (320 / float32(c.localcoord)) / oc.localscl)
	case OC_const_velocity_air_gethit_airrecover_fwd:
		sys.bcStack.PushF(c.gi().velocity.air.gethit.airrecover.fwd * (320 / float32(c.localcoord)) / oc.localscl)
	case OC_const_velocity_air_gethit_airrecover_up:
		sys.bcStack.PushF(c.gi().velocity.air.gethit.airrecover.up * (320 / float32(c.localcoord)) / oc.localscl)
	case OC_const_velocity_air_gethit_airrecover_down:
		sys.bcStack.PushF(c.gi().velocity.air.gethit.airrecover.down * (320 / float32(c.localcoord)) / oc.localscl)
	case OC_const_movement_airjump_num:
		sys.bcStack.PushI(c.gi().movement.airjump.num)
	case OC_const_movement_airjump_height:
		sys.bcStack.PushI(int32(float32(c.gi().movement.airjump.height) * (320 / float32(c.localcoord)) / oc.localscl))
	case OC_const_movement_yaccel:
		sys.bcStack.PushF(c.gi().movement.yaccel * (320 / float32(c.localcoord)) / oc.localscl)
	case OC_const_movement_stand_friction:
		sys.bcStack.PushF(c.gi().movement.stand.friction)
	case OC_const_movement_crouch_friction:
		sys.bcStack.PushF(c.gi().movement.crouch.friction)
	case OC_const_movement_stand_friction_threshold:
		sys.bcStack.PushF(c.gi().movement.stand.friction_threshold * (320 / float32(c.localcoord)) / oc.localscl)
	case OC_const_movement_crouch_friction_threshold:
		sys.bcStack.PushF(c.gi().movement.crouch.friction_threshold * (320 / float32(c.localcoord)) / oc.localscl)
	case OC_const_movement_air_gethit_groundlevel:
		sys.bcStack.PushF(c.gi().movement.air.gethit.groundlevel * (320 / float32(c.localcoord)) / oc.localscl)
	case OC_const_movement_air_gethit_groundrecover_ground_threshold:
		sys.bcStack.PushF(
			c.gi().movement.air.gethit.groundrecover.ground.threshold * (320 / float32(c.localcoord)) / oc.localscl)
	case OC_const_movement_air_gethit_groundrecover_groundlevel:
		sys.bcStack.PushF(c.gi().movement.air.gethit.groundrecover.groundlevel * (320 / float32(c.localcoord)) / oc.localscl)
	case OC_const_movement_air_gethit_airrecover_threshold:
		sys.bcStack.PushF(c.gi().movement.air.gethit.airrecover.threshold * (320 / float32(c.localcoord)) / oc.localscl)
	case OC_const_movement_air_gethit_airrecover_yaccel:
		sys.bcStack.PushF(c.gi().movement.air.gethit.airrecover.yaccel * (320 / float32(c.localcoord)) / oc.localscl)
	case OC_const_movement_air_gethit_trip_groundlevel:
		sys.bcStack.PushF(c.gi().movement.air.gethit.trip.groundlevel * (320 / float32(c.localcoord)) / oc.localscl)
	case OC_const_movement_down_bounce_offset_x:
		sys.bcStack.PushF(c.gi().movement.down.bounce.offset[0] * (320 / float32(c.localcoord)) / oc.localscl)
	case OC_const_movement_down_bounce_offset_y:
		sys.bcStack.PushF(c.gi().movement.down.bounce.offset[1] * (320 / float32(c.localcoord)) / oc.localscl)
	case OC_const_movement_down_bounce_yaccel:
		sys.bcStack.PushF(c.gi().movement.down.bounce.yaccel * (320 / float32(c.localcoord)) / oc.localscl)
	case OC_const_movement_down_bounce_groundlevel:
		sys.bcStack.PushF(c.gi().movement.down.bounce.groundlevel * (320 / float32(c.localcoord)) / oc.localscl)
	case OC_const_movement_down_friction_threshold:
		sys.bcStack.PushF(c.gi().movement.down.friction_threshold * (320 / float32(c.localcoord)) / oc.localscl)
	case OC_const_authorname:
		sys.bcStack.PushB(c.gi().authorLow ==
			sys.stringPool[sys.workingState.playerNo].List[*(*int32)(
				unsafe.Pointer(&be[*i]))])
		*i += 4
	case OC_const_name:
		sys.bcStack.PushB(c.gi().nameLow ==
			sys.stringPool[sys.workingState.playerNo].List[*(*int32)(
				unsafe.Pointer(&be[*i]))])
		*i += 4
	case OC_const_p2name:
		p2 := c.p2()
		sys.bcStack.PushB(p2 != nil && p2.gi().nameLow ==
			sys.stringPool[sys.workingState.playerNo].List[*(*int32)(
				unsafe.Pointer(&be[*i]))])
		*i += 4
	case OC_const_p3name:
		p3 := c.partner(0)
		sys.bcStack.PushB(p3 != nil && p3.gi().nameLow ==
			sys.stringPool[sys.workingState.playerNo].List[*(*int32)(
				unsafe.Pointer(&be[*i]))])
		*i += 4
	case OC_const_p4name:
		p4 := sys.charList.enemyNear(c, 1, true)
		sys.bcStack.PushB(p4 != nil && !(p4.scf(SCF_ko) && p4.scf(SCF_over)) &&
			p4.gi().nameLow ==
				sys.stringPool[sys.workingState.playerNo].List[*(*int32)(
					unsafe.Pointer(&be[*i]))])
		*i += 4
	case OC_const_stagevar_info_name:
		sys.bcStack.PushB(sys.stage.nameLow ==
			sys.stringPool[sys.workingState.playerNo].List[*(*int32)(
				unsafe.Pointer(&be[*i]))])
		*i += 4
	case OC_const_stagevar_info_displayname:
		sys.bcStack.PushB(sys.stage.displaynameLow ==
			sys.stringPool[sys.workingState.playerNo].List[*(*int32)(
				unsafe.Pointer(&be[*i]))])
		*i += 4
	case OC_const_stagevar_info_author:
		sys.bcStack.PushB(sys.stage.authorLow ==
			sys.stringPool[sys.workingState.playerNo].List[*(*int32)(
				unsafe.Pointer(&be[*i]))])
		*i += 4
	default:
		sys.errLog.Printf("%v\n", be[*i-1])
		c.panic()
	}
}
func (be BytecodeExp) run_ex(c *Char, i *int, oc *Char) {
	(*i)++
	switch be[*i-1] {
	case OC_ex_drawgame:
		sys.bcStack.PushB(c.drawgame())
	case OC_ex_ishometeam:
		sys.bcStack.PushB(c.playerNo&1 == sys.home)
	case OC_ex_lose:
		sys.bcStack.PushB(c.lose())
	case OC_ex_loseko:
		sys.bcStack.PushB(c.loseKO())
	case OC_ex_losetime:
		sys.bcStack.PushB(c.loseTime())
	case OC_ex_matchno:
		sys.bcStack.PushI(sys.match)
	case OC_ex_matchover:
		sys.bcStack.PushB(sys.matchOver())
	case OC_ex_roundno:
		sys.bcStack.PushI(sys.round)
	case OC_ex_tickspersecond:
		sys.bcStack.PushI(FPS)
	case OC_ex_win:
		sys.bcStack.PushB(c.win())
	case OC_ex_winko:
		sys.bcStack.PushB(c.winKO())
	case OC_ex_wintime:
		sys.bcStack.PushB(c.winTime())
	case OC_ex_winperfect:
		sys.bcStack.PushB(c.winPerfect())
	case OC_ex_p2dist_x:
		sys.bcStack.Push(c.rdDistX(c.p2(), oc))
	case OC_ex_p2dist_y:
		sys.bcStack.Push(c.rdDistY(c.p2(), oc))
	case OC_ex_p2bodydist_x:
		sys.bcStack.Push(c.p2BodyDistX(oc))
	case OC_ex_rootdist_x:
		sys.bcStack.Push(c.rdDistX(c.root(), oc))
	case OC_ex_rootdist_y:
		sys.bcStack.Push(c.rdDistY(c.root(), oc))
	case OC_ex_parentdist_x:
		sys.bcStack.Push(c.rdDistX(c.parent(), oc))
	case OC_ex_parentdist_y:
		sys.bcStack.Push(c.rdDistY(c.parent(), oc))
	case OC_ex_gethitvar_animtype:
		sys.bcStack.PushI(int32(c.gethitAnimtype()))
	case OC_ex_gethitvar_airtype:
		sys.bcStack.PushI(int32(c.ghv.airtype))
	case OC_ex_gethitvar_groundtype:
		sys.bcStack.PushI(int32(c.ghv.groundtype))
	case OC_ex_gethitvar_damage:
		sys.bcStack.PushI(c.ghv.damage)
	case OC_ex_gethitvar_hitcount:
		sys.bcStack.PushI(c.ghv.hitcount)
	case OC_ex_gethitvar_fallcount:
		sys.bcStack.PushI(c.ghv.fallcount)
	case OC_ex_gethitvar_hitshaketime:
		sys.bcStack.PushI(c.ghv.hitshaketime)
	case OC_ex_gethitvar_hittime:
		sys.bcStack.PushI(c.ghv.hittime)
	case OC_ex_gethitvar_slidetime:
		sys.bcStack.PushI(c.ghv.slidetime)
	case OC_ex_gethitvar_ctrltime:
		sys.bcStack.PushI(c.ghv.ctrltime)
	case OC_ex_gethitvar_recovertime:
		sys.bcStack.PushI(c.recoverTime)
	case OC_ex_gethitvar_xoff:
		sys.bcStack.PushF(c.ghv.xoff * c.localscl / oc.localscl)
	case OC_ex_gethitvar_yoff:
		sys.bcStack.PushF(c.ghv.yoff * c.localscl / oc.localscl)
	case OC_ex_gethitvar_xvel:
		sys.bcStack.PushF(c.ghv.xvel * c.facing * c.localscl / oc.localscl)
	case OC_ex_gethitvar_yvel:
		sys.bcStack.PushF(c.ghv.yvel * c.localscl / oc.localscl)
	case OC_ex_gethitvar_yaccel:
		sys.bcStack.PushF(c.ghv.getYaccel(oc) * c.localscl / oc.localscl)
	case OC_ex_gethitvar_chainid:
		sys.bcStack.PushI(c.ghv.chainId())
	case OC_ex_gethitvar_guarded:
		sys.bcStack.PushB(c.ghv.guarded)
	case OC_ex_gethitvar_isbound:
		sys.bcStack.PushB(c.isBound())
	case OC_ex_gethitvar_fall:
		sys.bcStack.PushB(c.ghv.fallf)
	case OC_ex_gethitvar_fall_damage:
		sys.bcStack.PushI(c.ghv.fall.damage)
	case OC_ex_gethitvar_fall_xvel:
		sys.bcStack.PushF(c.ghv.fall.xvel() * c.localscl / oc.localscl)
	case OC_ex_gethitvar_fall_yvel:
		sys.bcStack.PushF(c.ghv.fall.yvelocity * c.localscl / oc.localscl)
	case OC_ex_gethitvar_fall_recover:
		sys.bcStack.PushB(c.ghv.fall.recover)
	case OC_ex_gethitvar_fall_time:
		sys.bcStack.PushI(c.fallTime)
	case OC_ex_gethitvar_fall_recovertime:
		sys.bcStack.PushI(c.ghv.fall.recovertime)
	case OC_ex_gethitvar_fall_kill:
		sys.bcStack.PushB(c.ghv.fall.kill)
	case OC_ex_gethitvar_fall_envshake_time:
		sys.bcStack.PushI(c.ghv.fall.envshake_time)
	case OC_ex_gethitvar_fall_envshake_freq:
		sys.bcStack.PushF(c.ghv.fall.envshake_freq)
	case OC_ex_gethitvar_fall_envshake_ampl:
		sys.bcStack.PushI(int32(float32(c.ghv.fall.envshake_ampl) * c.localscl / oc.localscl))
	case OC_ex_gethitvar_fall_envshake_phase:
		sys.bcStack.PushF(c.ghv.fall.envshake_phase * c.localscl / oc.localscl)
	case OC_ex_majorversion:
		sys.bcStack.PushI(int32(c.gi().ver[0]))
	case OC_ex_drawpalno:
		sys.bcStack.PushI(c.gi().drawpalno)
	default:
		sys.errLog.Printf("%v\n", be[*i-1])
		c.panic()
	}
}
func (be BytecodeExp) evalF(c *Char) float32 {
	return be.run(c).ToF()
}
func (be BytecodeExp) evalI(c *Char) int32 {
	return be.run(c).ToI()
}
func (be BytecodeExp) evalB(c *Char) bool {
	return be.run(c).ToB()
}

type StateController interface {
	Run(c *Char, ps []int32) (changeState bool)
}
type NullStateController struct{}

func (_ NullStateController) Run(_ *Char, _ []int32) bool { return false }

var nullStateController NullStateController

type bytecodeFunction struct {
	numVars int32
	numRets int32
	numArgs int32
	ctrls   []StateController
}

func (bf bytecodeFunction) run(c *Char, ret []uint8) (changeState bool) {
	oldv, oldvslen := sys.bcVar, len(sys.bcVarStack)
	sys.bcVar = sys.bcVarStack.Alloc(int(bf.numVars))
	if len(sys.bcStack) != int(bf.numArgs) {
		c.panic()
	}
	copy(sys.bcVar, sys.bcStack)
	sys.bcStack.Clear()
	for _, sc := range bf.ctrls {
		switch sc.(type) {
		case StateBlock:
		default:
			if c.hitPause() {
				continue
			}
		}
		if sc.Run(c, nil) {
			changeState = true
			break
		}
	}
	if !changeState {
		if len(ret) > 0 {
			if len(ret) != int(bf.numRets) {
				c.panic()
			}
			for i, r := range ret {
				oldv[r] = sys.bcVar[int(bf.numArgs)+i]
			}
		}
	}
	sys.bcVar, sys.bcVarStack = oldv, sys.bcVarStack[:oldvslen]
	return
}

type callFunction struct {
	bytecodeFunction
	arg BytecodeExp
	ret []uint8
}

func (cf callFunction) Run(c *Char, _ []int32) (changeState bool) {
	if len(cf.arg) > 0 {
		sys.bcStack.Push(cf.arg.run(c))
	}
	return cf.run(c, cf.ret)
}

type StateBlock struct {
	persistent          int32
	persistentIndex     int32
	ignorehitpause      int32
	ctrlsIgnorehitpause bool
	trigger             BytecodeExp
	elseBlock           *StateBlock
	ctrls               []StateController
}

func newStateBlock() *StateBlock {
	return &StateBlock{persistent: 1, persistentIndex: -1, ignorehitpause: -2}
}
func (b StateBlock) Run(c *Char, ps []int32) (changeState bool) {
	if c.hitPause() {
		if b.ignorehitpause < -1 {
			return false
		}
		if b.ignorehitpause >= 0 {
			ww := &c.ss.wakegawakaranai[sys.workingState.playerNo][b.ignorehitpause]
			*ww = !*ww
			if !*ww {
				return false
			}
		}
	}
	if b.persistentIndex >= 0 {
		ps[b.persistentIndex]--
		if ps[b.persistentIndex] > 0 {
			return false
		}
	}
	sys.workingChar = c
	if len(b.trigger) > 0 && !b.trigger.evalB(c) {
		if b.elseBlock != nil {
			return b.elseBlock.Run(c, ps)
		}
		return false
	}
	for _, sc := range b.ctrls {
		switch sc.(type) {
		case StateBlock:
		default:
			if !b.ctrlsIgnorehitpause && c.hitPause() {
				continue
			}
		}
		if sc.Run(c, ps) {
			return true
		}
	}
	if b.persistentIndex >= 0 {
		ps[b.persistentIndex] = b.persistent
	}
	return false
}

type StateExpr BytecodeExp

func (se StateExpr) Run(c *Char, _ []int32) (changeState bool) {
	BytecodeExp(se).run(c)
	return false
}

type varAssign struct {
	vari uint8
	be   BytecodeExp
}

func (va varAssign) Run(c *Char, _ []int32) (changeState bool) {
	sys.bcVar[va.vari] = va.be.run(c)
	return false
}

type StateControllerBase []byte

func newStateControllerBase() *StateControllerBase {
	return (*StateControllerBase)(&[]byte{})
}
func (_ StateControllerBase) beToExp(be ...BytecodeExp) []BytecodeExp {
	return be
}
func (_ StateControllerBase) fToExp(f ...float32) (exp []BytecodeExp) {
	for _, v := range f {
		var be BytecodeExp
		be.appendValue(BytecodeFloat(v))
		exp = append(exp, be)
	}
	return
}
func (_ StateControllerBase) iToExp(i ...int32) (exp []BytecodeExp) {
	for _, v := range i {
		var be BytecodeExp
		be.appendValue(BytecodeInt(v))
		exp = append(exp, be)
	}
	return
}
func (scb *StateControllerBase) add(id byte, exp []BytecodeExp) {
	*scb = append(*scb, id, byte(len(exp)))
	for _, e := range exp {
		l := int32(len(e))
		*scb = append(*scb, (*(*[4]byte)(unsafe.Pointer(&l)))[:]...)
		*scb = append(*scb, *(*[]byte)(unsafe.Pointer(&e))...)
	}
}
func (scb StateControllerBase) run(c *Char,
	f func(byte, []BytecodeExp) bool) {
	for i := 0; i < len(scb); {
		id := scb[i]
		i++
		n := scb[i]
		i++
		if cap(sys.workBe) < int(n) {
			sys.workBe = make([]BytecodeExp, n)
		} else {
			sys.workBe = sys.workBe[:n]
		}
		for m := 0; m < int(n); m++ {
			l := *(*int32)(unsafe.Pointer(&scb[i]))
			i += 4
			sys.workBe[m] = (*(*BytecodeExp)(unsafe.Pointer(&scb)))[i : i+int(l)]
			i += int(l)
		}
		if !f(id, sys.workBe) {
			break
		}
	}
}

type stateDef StateControllerBase

const (
	stateDef_hitcountpersist byte = iota
	stateDef_movehitpersist
	stateDef_hitdefpersist
	stateDef_sprpriority
	stateDef_facep2
	stateDef_juggle
	stateDef_velset
	stateDef_anim
	stateDef_ctrl
	stateDef_poweradd
)

func (sc stateDef) Run(c *Char) {
	StateControllerBase(sc).run(c, func(id byte, exp []BytecodeExp) bool {
		switch id {
		case stateDef_hitcountpersist:
			if !exp[0].evalB(c) {
				c.clearHitCount()
			}
		case stateDef_movehitpersist:
			if !exp[0].evalB(c) {
				c.clearMoveHit()
			}
		case stateDef_hitdefpersist:
			if !exp[0].evalB(c) {
				c.clearHitDef()
			}
		case stateDef_sprpriority:
			c.setSprPriority(exp[0].evalI(c))
		case stateDef_facep2:
			if exp[0].evalB(c) && c.rdDistX(c.p2(), c).ToF() < 0 {
				c.setFacing(-c.facing)
			}
		case stateDef_juggle:
			c.setJuggle(exp[0].evalI(c))
		case stateDef_velset:
			c.setXV(exp[0].evalF(c))
			if len(exp) > 1 {
				c.setYV(exp[1].evalF(c))
				if len(exp) > 2 {
					exp[2].run(c)
				}
			}
		case stateDef_anim:
			c.changeAnim(exp[0].evalI(c))
		case stateDef_ctrl:
			c.setCtrl(exp[0].evalB(c))
		case stateDef_poweradd:
			c.powerAdd(exp[0].evalI(c))
		}
		return true
	})
}

type hitBy StateControllerBase

const (
	hitBy_value byte = iota
	hitBy_value2
	hitBy_time
)

func (sc hitBy) Run(c *Char, _ []int32) bool {
	time := int32(1)
	StateControllerBase(sc).run(c, func(id byte, exp []BytecodeExp) bool {
		switch id {
		case hitBy_time:
			time = exp[0].evalI(c)
		case hitBy_value:
			c.hitby[0].time = time
			c.hitby[0].flag = exp[0].evalI(c)
		case hitBy_value2:
			c.hitby[1].time = time
			c.hitby[1].flag = exp[0].evalI(c)
		}
		return true
	})
	return false
}

type notHitBy hitBy

func (sc notHitBy) Run(c *Char, _ []int32) bool {
	time := int32(1)
	StateControllerBase(sc).run(c, func(id byte, exp []BytecodeExp) bool {
		switch id {
		case hitBy_time:
			time = exp[0].evalI(c)
		case hitBy_value:
			c.hitby[0].time = time
			c.hitby[0].flag = ^exp[0].evalI(c)
		case hitBy_value2:
			c.hitby[1].time = time
			c.hitby[1].flag = ^exp[0].evalI(c)
		}
		return true
	})
	return false
}

type assertSpecial StateControllerBase

const (
	assertSpecial_flag byte = iota
	assertSpecial_flag_g
)

func (sc assertSpecial) Run(c *Char, _ []int32) bool {
	StateControllerBase(sc).run(c, func(id byte, exp []BytecodeExp) bool {
		switch id {
		case assertSpecial_flag:
			c.setSF(CharSpecialFlag(exp[0].evalI(c)))
		case assertSpecial_flag_g:
			sys.setSF(GlobalSpecialFlag(exp[0].evalI(c)))
		}
		return true
	})
	return false
}

type playSnd StateControllerBase

const (
	playSnd_value = iota
	playSnd_channel
	playSnd_lowpriority
	playSnd_pan
	playSnd_abspan
	playSnd_volume
	playSnd_freqmul
	playSnd_loop
)

func (sc playSnd) Run(c *Char, _ []int32) bool {
	f, lw, lp := false, false, false
	var g, n, ch, vo int32 = -1, 0, -1, 0
	if c.gi().ver[0] == 1 {
		vo = 100
	}
	var p, fr float32 = 0, 1
	x := &c.pos[0]
	StateControllerBase(sc).run(c, func(id byte, exp []BytecodeExp) bool {
		switch id {
		case playSnd_value:
			f = exp[0].evalB(c)
			g = exp[1].evalI(c)
			if len(exp) > 2 {
				n = exp[2].evalI(c)
			}
		case playSnd_channel:
			ch = exp[0].evalI(c)
		case playSnd_lowpriority:
			lw = exp[0].evalB(c)
		case playSnd_pan:
			p = exp[0].evalF(c)
		case playSnd_abspan:
			x = nil
			p = exp[0].evalF(c)
		case playSnd_volume:
			vo = exp[0].evalI(c)
		case playSnd_freqmul:
			fr = exp[0].evalF(c)
		case playSnd_loop:
			lp = exp[0].evalB(c)
		}
		return true
	})
	c.playSound(f, lw, lp, g, n, ch, vo, p, fr, x)
	return false
}

type changeState StateControllerBase

const (
	changeState_value byte = iota
	changeState_ctrl
	changeState_anim
)

func (sc changeState) Run(c *Char, _ []int32) bool {
	var v, a, ctrl int32 = -1, -1, -1
	StateControllerBase(sc).run(c, func(id byte, exp []BytecodeExp) bool {
		switch id {
		case changeState_value:
			v = exp[0].evalI(c)
		case changeState_ctrl:
			ctrl = exp[0].evalI(c)
		case changeState_anim:
			a = exp[0].evalI(c)
		}
		return true
	})
	c.changeState(v, a, ctrl)
	return true
}

type selfState changeState

func (sc selfState) Run(c *Char, _ []int32) bool {
	var v, a, ctrl int32 = -1, -1, -1
	StateControllerBase(sc).run(c, func(id byte, exp []BytecodeExp) bool {
		switch id {
		case changeState_value:
			v = exp[0].evalI(c)
		case changeState_ctrl:
			ctrl = exp[0].evalI(c)
		case changeState_anim:
			a = exp[0].evalI(c)
		}
		return true
	})
	c.selfState(v, a, ctrl)
	return true
}

type tagIn StateControllerBase

const (
	tagIn_stateno = iota
	tagIn_partnerstateno
)

func (sc tagIn) Run(c *Char, _ []int32) bool {
	p := c.partner(0)
	if p == nil {
		return false
	}
	sn := int32(-1)
	ret := false
	StateControllerBase(sc).run(c, func(id byte, exp []BytecodeExp) bool {
		switch id {
		case tagIn_stateno:
			sn = exp[0].evalI(c)
		case tagIn_partnerstateno:
			if psn := exp[0].evalI(c); psn >= 0 {
				if sn >= 0 {
					c.changeState(sn, -1, -1)
				}
				p.unsetSCF(SCF_standby)
				p.changeState(psn, -1, -1)
				ret = true
			} else {
				return false
			}
		}
		return true
	})
	return ret
}

type tagOut StateControllerBase

const (
	tagOut_ = iota
)

func (sc tagOut) Run(c *Char, _ []int32) bool {
	StateControllerBase(sc).run(c, func(id byte, exp []BytecodeExp) bool {
		switch id {
		case tagOut_:
			c.setSCF(SCF_standby)
		}
		return true
	})
	return true
}

type destroySelf StateControllerBase

const (
	destroySelf_recursive = iota
	destroySelf_removeexplods
)

func (sc destroySelf) Run(c *Char, _ []int32) bool {
	rec, rem := false, false
	StateControllerBase(sc).run(c, func(id byte, exp []BytecodeExp) bool {
		switch id {
		case destroySelf_recursive:
			rec = exp[0].evalB(c)
		case destroySelf_removeexplods:
			rem = exp[0].evalB(c)
		}
		return true
	})
	return c.destroySelf(rec, rem)
}

type changeAnim StateControllerBase

const (
	changeAnim_elem byte = iota
	changeAnim_value
)

func (sc changeAnim) Run(c *Char, _ []int32) bool {
	var elem int32
	setelem := false
	StateControllerBase(sc).run(c, func(id byte, exp []BytecodeExp) bool {
		switch id {
		case changeAnim_elem:
			elem = exp[0].evalI(c)
			setelem = true
		case changeAnim_value:
			c.changeAnim(exp[0].evalI(c))
			if setelem {
				c.setAnimElem(elem)
			}
		}
		return true
	})
	return false
}

type changeAnim2 changeAnim

func (sc changeAnim2) Run(c *Char, _ []int32) bool {
	var elem int32
	setelem := false
	StateControllerBase(sc).run(c, func(id byte, exp []BytecodeExp) bool {
		switch id {
		case changeAnim_elem:
			elem = exp[0].evalI(c)
			setelem = true
		case changeAnim_value:
			c.changeAnim2(exp[0].evalI(c))
			if setelem {
				c.setAnimElem(elem)
			}
		}
		return true
	})
	return false
}

type helper StateControllerBase

const (
	helper_helpertype byte = iota
	helper_name
	helper_postype
	helper_ownpal
	helper_size_xscale
	helper_size_yscale
	helper_size_ground_back
	helper_size_ground_front
	helper_size_air_back
	helper_size_air_front
	helper_size_height
	helper_size_proj_doscale
	helper_size_head_pos
	helper_size_mid_pos
	helper_size_shadowoffset
	helper_stateno
	helper_keyctrl
	helper_id
	helper_pos
	helper_facing
	helper_pausemovetime
	helper_supermovetime
)

func (sc helper) Run(c *Char, _ []int32) bool {
	h := c.newHelper()
	if h == nil {
		return false
	}
	pt := PT_P1
	var f, st int32 = 1, 0
	op := false
	var x, y float32 = 0, 0
	StateControllerBase(sc).run(c, func(id byte, exp []BytecodeExp) bool {
		switch id {
		case helper_helpertype:
			h.player = exp[0].evalB(c)
		case helper_name:
			h.name = string(*(*[]byte)(unsafe.Pointer(&exp[0])))
		case helper_postype:
			pt = PosType(exp[0].evalI(c))
		case helper_ownpal:
			op = exp[0].evalB(c)
		case helper_size_xscale:
			h.size.xscale = exp[0].evalF(c)
		case helper_size_yscale:
			h.size.yscale = exp[0].evalF(c)
		case helper_size_ground_back:
			h.size.ground.back = exp[0].evalF(c)
		case helper_size_ground_front:
			h.size.ground.front = exp[0].evalF(c)
		case helper_size_air_back:
			h.size.air.back = exp[0].evalF(c)
		case helper_size_air_front:
			h.size.air.front = exp[0].evalF(c)
		case helper_size_height:
			h.size.height = exp[0].evalF(c)
		case helper_size_proj_doscale:
			h.size.proj.doscale = exp[0].evalI(c)
		case helper_size_head_pos:
			h.size.head.pos[0] = exp[0].evalF(c)
			if len(exp) > 1 {
				h.size.head.pos[1] = exp[1].evalF(c)
			}
		case helper_size_mid_pos:
			h.size.mid.pos[0] = exp[0].evalF(c)
			if len(exp) > 1 {
				h.size.mid.pos[1] = exp[1].evalF(c)
			}
		case helper_size_shadowoffset:
			h.size.shadowoffset = exp[0].evalF(c)
		case helper_stateno:
			st = exp[0].evalI(c)
		case helper_keyctrl:
			h.keyctrl = exp[0].evalB(c)
		case helper_id:
			h.helperId = exp[0].evalI(c)
		case helper_pos:
			x = exp[0].evalF(c)
			if len(exp) > 1 {
				y = exp[1].evalF(c)
			}
		case helper_facing:
			f = exp[0].evalI(c)
		case helper_pausemovetime:
			h.pauseMovetime = exp[0].evalI(c)
		case helper_supermovetime:
			h.superMovetime = exp[0].evalI(c)
		}
		h.localscl = c.localscl
		h.localcoord = c.localcoord
		return true
	})
	c.helperInit(h, st, pt, x, y, f, op)
	return false
}

type ctrlSet StateControllerBase

const (
	ctrlSet_value byte = iota
)

func (sc ctrlSet) Run(c *Char, _ []int32) bool {
	StateControllerBase(sc).run(c, func(id byte, exp []BytecodeExp) bool {
		switch id {
		case ctrlSet_value:
			c.setCtrl(exp[0].evalB(c))
		}
		return true
	})
	return false
}

type explod StateControllerBase

const (
	explod_ownpal byte = iota
	explod_remappal
	explod_id
	explod_facing
	explod_vfacing
	explod_pos
	explod_random
	explod_postype
	explod_velocity
	explod_accel
	explod_scale
	explod_bindtime
	explod_removetime
	explod_supermove
	explod_supermovetime
	explod_pausemovetime
	explod_sprpriority
	explod_ontop
	explod_strictontop
	explod_shadow
	explod_removeongethit
	explod_trans
	explod_anim
	explod_angle
	explod_yangle
	explod_xangle
	explod_ignorehitpause
	explod_bindid
	explod_space
)

func (sc explod) Run(c *Char, _ []int32) bool {
	e, i := c.newExplod()
	if e == nil {
		return false
	}
	e.id = 0
	rp := [...]int32{-1, 0}
	if c.stCgi().ver[1] == 1 && c.stCgi().ver[1] == 1 {
		e.postype = PT_N
	}
	StateControllerBase(sc).run(c, func(id byte, exp []BytecodeExp) bool {
		switch id {
		case explod_ownpal:
			e.ownpal = exp[0].evalB(c)
		case explod_remappal:
			rp[0] = exp[0].evalI(c)
			if len(exp) > 1 {
				rp[1] = exp[1].evalI(c)
			}
		case explod_id:
			e.id = Max(0, exp[0].evalI(c))
		case explod_facing:
			if exp[0].evalI(c) < 0 {
				e.relativef = -1
			} else {
				e.relativef = 1
			}
		case explod_vfacing:
			if exp[0].evalI(c) < 0 {
				e.vfacing = -1
			} else {
				e.vfacing = 1
			}
		case explod_pos:
			e.offset[0] = exp[0].evalF(c)
			if len(exp) > 1 {
				e.offset[1] = exp[1].evalF(c)
			}
		case explod_random:
			rndx := exp[0].evalF(c)
			e.offset[0] += RandF(-rndx, rndx)
			if len(exp) > 1 {
				rndy := exp[1].evalF(c)
				e.offset[1] += RandF(-rndy, rndy)
			}
		case explod_postype:
			e.postype = PosType(exp[0].evalI(c))
		case explod_space:
			e.space = Space(exp[0].evalI(c))
		case explod_velocity:
			e.velocity[0] = exp[0].evalF(c)
			if len(exp) > 1 {
				e.velocity[1] = exp[1].evalF(c)
			}
		case explod_accel:
			e.accel[0] = exp[0].evalF(c)
			if len(exp) > 1 {
				e.accel[1] = exp[1].evalF(c)
			}
		case explod_scale:
			e.scale[0] = exp[0].evalF(c)
			if len(exp) > 1 {
				e.scale[1] = exp[1].evalF(c)
			}
		case explod_bindtime:
			e.bindtime = exp[0].evalI(c)
		case explod_removetime:
			e.removetime = exp[0].evalI(c)
		case explod_supermove:
			if exp[0].evalB(c) {
				e.supermovetime = -1
			} else {
				e.supermovetime = 0
			}
		case explod_supermovetime:
			e.supermovetime = exp[0].evalI(c)
		case explod_pausemovetime:
			e.pausemovetime = exp[0].evalI(c)
		case explod_sprpriority:
			e.sprpriority = exp[0].evalI(c)
		case explod_ontop:
			e.ontop = exp[0].evalB(c)
		case explod_strictontop:
			if e.ontop {
				e.sprpriority = 0
			}
		case explod_shadow:
			e.shadow[0] = exp[0].evalI(c)
			if len(exp) > 1 {
				e.shadow[1] = exp[1].evalI(c)
				if len(exp) > 2 {
					e.shadow[2] = exp[2].evalI(c)
				}
			}
		case explod_removeongethit:
			e.removeongethit = exp[0].evalB(c)
		case explod_trans:
			e.alpha[0] = exp[0].evalI(c)
			e.alpha[1] = exp[1].evalI(c)
			if len(exp) >= 3 {
				e.alpha[0] = Max(0, Min(255, e.alpha[0]))
				e.alpha[1] = Max(0, Min(255, e.alpha[1]))
				if len(exp) >= 4 {
					e.alpha[1] = ^e.alpha[1]
				} else if e.alpha[0] == 1 && e.alpha[1] == 255 {
					e.alpha[0] = 0
				}
			}
		case explod_anim:
			e.anim = c.getAnim(exp[1].evalI(c), exp[0].evalB(c))
		case explod_angle:
			e.angle = exp[0].evalF(c)
		case explod_yangle:
			exp[0].run(c)
		case explod_xangle:
			exp[0].run(c)
		case explod_ignorehitpause:
			e.ignorehitpause = exp[0].evalB(c)
		case explod_bindid:
			e.bindId = exp[0].evalI(c)
		}
		return true
	})
	e.localscl = c.localscl
	e.setPos(c)
	c.insertExplodEx(i, rp)
	return false
}

type modifyExplod explod

func (sc modifyExplod) Run(c *Char, _ []int32) bool {
	eid := int32(-1)
	var expls []*Explod
	rp := [...]int32{-1, 0}
	eachExpl := func(f func(e *Explod)) {
		for _, e := range expls {
			f(e)
		}
	}
	StateControllerBase(sc).run(c, func(id byte, exp []BytecodeExp) bool {
		switch id {
		case explod_remappal:
			rp[0] = exp[0].evalI(c)
			if len(exp) > 1 {
				rp[1] = exp[1].evalI(c)
			}
		case explod_id:
			eid = exp[0].evalI(c)
		default:
			if len(expls) == 0 {
				expls = c.getExplods(eid)
				if len(expls) == 0 {
					return false
				}
				eachExpl(func(e *Explod) {
					if e.ownpal {
						c.remapPal(e.palfx, [...]int32{1, 1}, rp)
					}
				})
			}
			switch id {
			case explod_facing:
				if exp[0].evalI(c) < 0 {
					eachExpl(func(e *Explod) { e.relativef = -1 })
				} else {
					eachExpl(func(e *Explod) { e.relativef = 1 })
				}
			case explod_vfacing:
				if exp[0].evalI(c) < 0 {
					eachExpl(func(e *Explod) { e.vfacing = -1 })
				} else {
					eachExpl(func(e *Explod) { e.vfacing = 1 })
				}
			case explod_pos:
				x := exp[0].evalF(c)
				eachExpl(func(e *Explod) { e.offset[0] = x })
				if len(exp) > 1 {
					y := exp[1].evalF(c)
					eachExpl(func(e *Explod) { e.offset[1] = y })
				}
			case explod_random:
				rndx := exp[0].evalF(c)
				rndx = RandF(-rndx, rndx)
				eachExpl(func(e *Explod) { e.offset[0] += rndx })
				if len(exp) > 1 {
					rndy := exp[1].evalF(c)
					rndy = RandF(-rndy, rndy)
					eachExpl(func(e *Explod) { e.offset[1] += rndy })
				}
			case explod_postype:
				pt := PosType(exp[0].evalI(c))
				eachExpl(func(e *Explod) {
					e.postype = pt
					e.setPos(c)
				})
			case explod_velocity:
				x := exp[0].evalF(c)
				eachExpl(func(e *Explod) { e.velocity[0] = x })
				if len(exp) > 1 {
					y := exp[1].evalF(c)
					eachExpl(func(e *Explod) { e.velocity[1] = y })
				}
			case explod_accel:
				x := exp[0].evalF(c)
				eachExpl(func(e *Explod) { e.accel[0] = x })
				if len(exp) > 1 {
					y := exp[1].evalF(c)
					eachExpl(func(e *Explod) { e.accel[1] = y })
				}
			case explod_scale:
				x := exp[0].evalF(c)
				eachExpl(func(e *Explod) { e.scale[0] = x })
				if len(exp) > 1 {
					y := exp[1].evalF(c)
					eachExpl(func(e *Explod) { e.scale[1] = y })
				}
			case explod_bindtime:
				t := exp[0].evalI(c)
				eachExpl(func(e *Explod) { e.bindtime = t })
			case explod_removetime:
				t := exp[0].evalI(c)
				eachExpl(func(e *Explod) { e.removetime = t })
			case explod_supermove:
				if exp[0].evalB(c) {
					eachExpl(func(e *Explod) { e.supermovetime = -1 })
				} else {
					eachExpl(func(e *Explod) { e.supermovetime = 0 })
				}
			case explod_supermovetime:
				t := exp[0].evalI(c)
				eachExpl(func(e *Explod) { e.supermovetime = t })
			case explod_pausemovetime:
				t := exp[0].evalI(c)
				eachExpl(func(e *Explod) { e.pausemovetime = t })
			case explod_sprpriority:
				t := exp[0].evalI(c)
				eachExpl(func(e *Explod) { e.sprpriority = t })
			case explod_ontop:
				t := exp[0].evalB(c)
				eachExpl(func(e *Explod) {
					e.ontop = t
				})
			case explod_strictontop:
				eachExpl(func(e *Explod) {
					if e.ontop {
						e.sprpriority = 0
					}
				})
			case explod_shadow:
				r := exp[0].evalI(c)
				eachExpl(func(e *Explod) { e.shadow[0] = r })
				if len(exp) > 1 {
					g := exp[1].evalI(c)
					eachExpl(func(e *Explod) { e.shadow[1] = g })
					if len(exp) > 2 {
						b := exp[2].evalI(c)
						eachExpl(func(e *Explod) { e.shadow[2] = b })
					}
				}
			case explod_removeongethit:
				t := exp[0].evalB(c)
				eachExpl(func(e *Explod) { e.removeongethit = t })
			case explod_trans:
				s, d := exp[0].evalI(c), exp[1].evalI(c)
				if len(exp) >= 3 {
					s, d = Max(0, Min(255, s)), Max(0, Min(255, d))
					if len(exp) >= 4 {
						d = ^d
					} else if s == 1 && d == 255 {
						s = 0
					}
				}
				eachExpl(func(e *Explod) { e.alpha = [...]int32{s, d} })
			case explod_angle:
				a := exp[0].evalF(c)
				eachExpl(func(e *Explod) { e.angle = a })
			case explod_yangle:
				exp[0].run(c)
			case explod_xangle:
				exp[0].run(c)
			case explod_bindid:
				exp[0].evalI(c)
			}
		}
		return true
	})
	return false
}

type gameMakeAnim StateControllerBase

const (
	gameMakeAnim_pos byte = iota
	gameMakeAnim_random
	gameMakeAnim_under
	gameMakeAnim_anim
)

func (sc gameMakeAnim) Run(c *Char, _ []int32) bool {
	e, i := c.newExplod()
	if e == nil {
		return false
	}
	e.ontop, e.sprpriority, e.ownpal = true, math.MinInt32, true
	StateControllerBase(sc).run(c, func(id byte, exp []BytecodeExp) bool {
		switch id {
		case gameMakeAnim_pos:
			e.offset[0] = exp[0].evalF(c)
			if len(exp) > 1 {
				e.offset[1] = exp[1].evalF(c)
			}
		case gameMakeAnim_random:
			rndx := exp[0].evalF(c)
			e.offset[0] += RandF(-rndx, rndx)
			if len(exp) > 1 {
				rndy := exp[1].evalF(c)
				e.offset[1] += RandF(-rndy, rndy)
			}
		case gameMakeAnim_under:
			e.ontop = !exp[0].evalB(c)
		case gameMakeAnim_anim:
			e.anim = c.getAnim(exp[1].evalI(c), exp[0].evalB(c))
		}
		return true
	})
	e.offset[0] -= float32(c.size.draw.offset[0])
	e.offset[1] -= float32(c.size.draw.offset[1])
	//e.localscl = c.localscl
	e.setPos(c)
	c.insertExplod(i)
	return false
}

type posSet StateControllerBase

const (
	posSet_x byte = iota
	posSet_y
	posSet_z
)

func (sc posSet) Run(c *Char, _ []int32) bool {
	StateControllerBase(sc).run(c, func(id byte, exp []BytecodeExp) bool {
		switch id {
		case posSet_x:
			c.setX(sys.cam.Pos[0]/c.localscl + exp[0].evalF(c))
		case posSet_y:
			c.setY(exp[0].evalF(c))
		case posSet_z:
			exp[0].run(c)
		}
		return true
	})
	return false
}

type posAdd posSet

func (sc posAdd) Run(c *Char, _ []int32) bool {
	StateControllerBase(sc).run(c, func(id byte, exp []BytecodeExp) bool {
		switch id {
		case posSet_x:
			c.addX(exp[0].evalF(c))
		case posSet_y:
			c.addY(exp[0].evalF(c))
		case posSet_z:
			exp[0].run(c)
		}
		return true
	})
	return false
}

type velSet posSet

func (sc velSet) Run(c *Char, _ []int32) bool {
	StateControllerBase(sc).run(c, func(id byte, exp []BytecodeExp) bool {
		switch id {
		case posSet_x:
			c.setXV(exp[0].evalF(c))
		case posSet_y:
			c.setYV(exp[0].evalF(c))
		case posSet_z:
			exp[0].run(c)
		}
		return true
	})
	return false
}

type velAdd posSet

func (sc velAdd) Run(c *Char, _ []int32) bool {
	StateControllerBase(sc).run(c, func(id byte, exp []BytecodeExp) bool {
		switch id {
		case posSet_x:
			c.addXV(exp[0].evalF(c))
		case posSet_y:
			c.addYV(exp[0].evalF(c))
		case posSet_z:
			exp[0].run(c)
		}
		return true
	})
	return false
}

type velMul posSet

func (sc velMul) Run(c *Char, _ []int32) bool {
	StateControllerBase(sc).run(c, func(id byte, exp []BytecodeExp) bool {
		switch id {
		case posSet_x:
			c.mulXV(exp[0].evalF(c))
		case posSet_y:
			c.mulYV(exp[0].evalF(c))
		case posSet_z:
			exp[0].run(c)
		}
		return true
	})
	return false
}

type palFX StateControllerBase

const (
	palFX_time byte = iota
	palFX_color
	palFX_add
	palFX_mul
	palFX_sinadd
	palFX_invertall
	palFX_last = iota - 1
)

func (sc palFX) runSub(c *Char, pfd *PalFXDef,
	id byte, exp []BytecodeExp) bool {
	switch id {
	case palFX_time:
		pfd.time = exp[0].evalI(c)
	case palFX_color:
		pfd.color = MaxF(0, MinF(1, exp[0].evalF(c)/256))
	case palFX_add:
		pfd.add[0] = exp[0].evalI(c)
		pfd.add[1] = exp[1].evalI(c)
		pfd.add[2] = exp[2].evalI(c)
	case palFX_mul:
		pfd.mul[0] = exp[0].evalI(c)
		pfd.mul[1] = exp[1].evalI(c)
		pfd.mul[2] = exp[2].evalI(c)
	case palFX_sinadd:
		pfd.sinadd[0] = exp[0].evalI(c)
		pfd.sinadd[1] = exp[1].evalI(c)
		pfd.sinadd[2] = exp[2].evalI(c)
		if len(exp) > 3 {
			pfd.cycletime = exp[3].evalI(c)
		}
	case palFX_invertall:
		pfd.invertall = exp[0].evalB(c)
	default:
		return false
	}
	return true
}
func (sc palFX) Run(c *Char, _ []int32) bool {
	pf := c.palfx
	if pf == nil {
		pf = newPalFX()
	}
	pf.clear2(true)
	StateControllerBase(sc).run(c, func(id byte, exp []BytecodeExp) bool {
		sc.runSub(c, &pf.PalFXDef, id, exp)
		return true
	})
	return false
}

type allPalFX palFX

func (sc allPalFX) Run(c *Char, _ []int32) bool {
	sys.allPalFX.clear()
	StateControllerBase(sc).run(c, func(id byte, exp []BytecodeExp) bool {
		palFX(sc).runSub(c, &sys.allPalFX.PalFXDef, id, exp)
		return true
	})
	return false
}

type bgPalFX palFX

func (sc bgPalFX) Run(c *Char, _ []int32) bool {
	sys.bgPalFX.clear()
	StateControllerBase(sc).run(c, func(id byte, exp []BytecodeExp) bool {
		palFX(sc).runSub(c, &sys.bgPalFX.PalFXDef, id, exp)
		return true
	})
	return false
}

type afterImage palFX

const (
	afterImage_trans = iota + palFX_last + 1
	afterImage_time
	afterImage_length
	afterImage_timegap
	afterImage_framegap
	afterImage_palcolor
	afterImage_palinvertall
	afterImage_palbright
	afterImage_palcontrast
	afterImage_palpostbright
	afterImage_paladd
	afterImage_palmul
	afterImage_last = iota + palFX_last + 1 - 1
)

func (sc afterImage) runSub(c *Char, ai *AfterImage,
	id byte, exp []BytecodeExp) {
	switch id {
	case afterImage_trans:
		ai.alpha = [...]int32{exp[0].evalI(c), exp[1].evalI(c)}
	case afterImage_time:
		ai.time = exp[0].evalI(c)
	case afterImage_length:
		ai.length = exp[0].evalI(c)
	case afterImage_timegap:
		ai.timegap = Max(1, exp[0].evalI(c))
	case afterImage_framegap:
		ai.framegap = exp[0].evalI(c)
	case afterImage_palcolor:
		ai.setPalColor(exp[0].evalI(c))
	case afterImage_palinvertall:
		ai.setPalInvertall(exp[0].evalB(c))
	case afterImage_palbright:
		ai.setPalBrightR(exp[0].evalI(c))
		if len(exp) > 1 {
			ai.setPalBrightG(exp[1].evalI(c))
			if len(exp) > 2 {
				ai.setPalBrightB(exp[2].evalI(c))
			}
		}
	case afterImage_palcontrast:
		ai.setPalContrastR(exp[0].evalI(c))
		if len(exp) > 1 {
			ai.setPalContrastG(exp[1].evalI(c))
			if len(exp) > 2 {
				ai.setPalContrastB(exp[2].evalI(c))
			}
		}
	case afterImage_palpostbright:
		ai.postbright[0] = exp[0].evalI(c)
		if len(exp) > 1 {
			ai.postbright[1] = exp[1].evalI(c)
			if len(exp) > 2 {
				ai.postbright[2] = exp[2].evalI(c)
			}
		}
	case afterImage_paladd:
		ai.add[0] = exp[0].evalI(c)
		if len(exp) > 1 {
			ai.add[1] = exp[1].evalI(c)
			if len(exp) > 2 {
				ai.add[2] = exp[2].evalI(c)
			}
		}
	case afterImage_palmul:
		ai.mul[0] = exp[0].evalF(c)
		if len(exp) > 1 {
			ai.mul[1] = exp[1].evalF(c)
			if len(exp) > 2 {
				ai.mul[2] = exp[2].evalF(c)
			}
		}
	}
}
func (sc afterImage) Run(c *Char, _ []int32) bool {
	c.aimg.clear()
	c.aimg.time = 1
	StateControllerBase(sc).run(c, func(id byte, exp []BytecodeExp) bool {
		sc.runSub(c, &c.aimg, id, exp)
		return true
	})
	c.aimg.setupPalFX()
	return false
}

type afterImageTime StateControllerBase

const (
	afterImageTime_time byte = iota
)

func (sc afterImageTime) Run(c *Char, _ []int32) bool {
	StateControllerBase(sc).run(c, func(id byte, exp []BytecodeExp) bool {
		if c.aimg.timegap <= 0 {
			return false
		}
		switch id {
		case afterImageTime_time:
			c.aimg.time = exp[0].evalI(c)
		}
		return true
	})
	return false
}

type hitDef afterImage

const (
	hitDef_attr = iota + afterImage_last + 1
	hitDef_guardflag
	hitDef_hitflag
	hitDef_ground_type
	hitDef_air_type
	hitDef_animtype
	hitDef_air_animtype
	hitDef_fall_animtype
	hitDef_affectteam
	hitDef_id
	hitDef_chainid
	hitDef_nochainid
	hitDef_kill
	hitDef_guard_kill
	hitDef_fall_kill
	hitDef_hitonce
	hitDef_air_juggle
	hitDef_getpower
	hitDef_damage
	hitDef_givepower
	hitDef_numhits
	hitDef_hitsound
	hitDef_guardsound
	hitDef_priority
	hitDef_p1stateno
	hitDef_p2stateno
	hitDef_p2getp1state
	hitDef_p1sprpriority
	hitDef_p2sprpriority
	hitDef_forcestand
	hitDef_forcenofall
	hitDef_fall_damage
	hitDef_fall_xvelocity
	hitDef_fall_yvelocity
	hitDef_fall_recover
	hitDef_fall_recovertime
	hitDef_sparkno
	hitDef_guard_sparkno
	hitDef_sparkxy
	hitDef_down_hittime
	hitDef_p1facing
	hitDef_p1getp2facing
	hitDef_mindist
	hitDef_maxdist
	hitDef_snap
	hitDef_p2facing
	hitDef_air_hittime
	hitDef_fall
	hitDef_air_fall
	hitDef_air_cornerpush_veloff
	hitDef_down_bounce
	hitDef_down_velocity
	hitDef_down_cornerpush_veloff
	hitDef_ground_hittime
	hitDef_guard_hittime
	hitDef_guard_dist
	hitDef_pausetime
	hitDef_guard_pausetime
	hitDef_air_velocity
	hitDef_airguard_velocity
	hitDef_ground_slidetime
	hitDef_guard_slidetime
	hitDef_guard_ctrltime
	hitDef_airguard_ctrltime
	hitDef_ground_velocity_x
	hitDef_ground_velocity_y
	hitDef_ground_velocity
	hitDef_guard_velocity
	hitDef_ground_cornerpush_veloff
	hitDef_guard_cornerpush_veloff
	hitDef_airguard_cornerpush_veloff
	hitDef_yaccel
	hitDef_envshake_time
	hitDef_envshake_ampl
	hitDef_envshake_phase
	hitDef_envshake_freq
	hitDef_fall_envshake_time
	hitDef_fall_envshake_ampl
	hitDef_fall_envshake_phase
	hitDef_fall_envshake_freq
	hitDef_last = iota + afterImage_last + 1 - 1
)

func (sc hitDef) runSub(c *Char, hd *HitDef, id byte, exp []BytecodeExp) bool {
	switch id {
	case hitDef_attr:
		hd.attr = exp[0].evalI(c)
	case hitDef_guardflag:
		hd.guardflag = exp[0].evalI(c)
	case hitDef_hitflag:
		hd.hitflag = exp[0].evalI(c)
	case hitDef_ground_type:
		hd.ground_type = HitType(exp[0].evalI(c))
	case hitDef_air_type:
		hd.air_type = HitType(exp[0].evalI(c))
	case hitDef_animtype:
		hd.animtype = Reaction(exp[0].evalI(c))
	case hitDef_air_animtype:
		hd.air_animtype = Reaction(exp[0].evalI(c))
	case hitDef_fall_animtype:
		hd.fall.animtype = Reaction(exp[0].evalI(c))
	case hitDef_affectteam:
		hd.affectteam = exp[0].evalI(c)
	case hitDef_id:
		hd.id = Max(0, exp[0].evalI(c))
	case hitDef_chainid:
		hd.chainid = exp[0].evalI(c)
	case hitDef_nochainid:
		hd.nochainid[0] = exp[0].evalI(c)
		if len(exp) > 1 {
			hd.nochainid[1] = exp[1].evalI(c)
		}
	case hitDef_kill:
		hd.kill = exp[0].evalB(c)
	case hitDef_guard_kill:
		hd.guard_kill = exp[0].evalB(c)
	case hitDef_fall_kill:
		hd.fall.kill = exp[0].evalB(c)
	case hitDef_hitonce:
		hd.hitonce = Btoi(exp[0].evalB(c))
	case hitDef_air_juggle:
		hd.air_juggle = exp[0].evalI(c)
	case hitDef_getpower:
		hd.hitgetpower = Max(IErr+1, exp[0].evalI(c))
		if len(exp) > 1 {
			hd.guardgetpower = Max(IErr+1, exp[1].evalI(c))
		}
	case hitDef_damage:
		hd.hitdamage = exp[0].evalI(c)
		if len(exp) > 1 {
			hd.guarddamage = exp[1].evalI(c)
		}
	case hitDef_givepower:
		hd.hitgivepower = Max(IErr+1, exp[0].evalI(c))
		if len(exp) > 1 {
			hd.guardgivepower = Max(IErr+1, exp[1].evalI(c))
		}
	case hitDef_numhits:
		hd.numhits = exp[0].evalI(c)
	case hitDef_hitsound:
		n := exp[1].evalI(c)
		if n < 0 {
			hd.hitsound[0] = IErr
		} else if exp[0].evalB(c) {
			hd.hitsound[0] = ^n
		} else {
			hd.hitsound[0] = n
		}
		if len(exp) > 2 {
			hd.hitsound[1] = exp[2].evalI(c)
		}
	case hitDef_guardsound:
		n := exp[1].evalI(c)
		if n < 0 {
			hd.guardsound[0] = IErr
		} else if exp[0].evalB(c) {
			hd.guardsound[0] = ^n
		} else {
			hd.guardsound[0] = n
		}
		if len(exp) > 2 {
			hd.guardsound[1] = exp[2].evalI(c)
		}
	case hitDef_priority:
		hd.priority = exp[0].evalI(c)
		hd.bothhittype = AiuchiType(exp[1].evalI(c))
	case hitDef_p1stateno:
		hd.p1stateno = exp[0].evalI(c)
	case hitDef_p2stateno:
		hd.p2stateno = exp[0].evalI(c)
		hd.p2getp1state = true
	case hitDef_p2getp1state:
		hd.p2getp1state = exp[0].evalB(c)
	case hitDef_p1sprpriority:
		hd.p1sprpriority = exp[0].evalI(c)
	case hitDef_p2sprpriority:
		hd.p2sprpriority = exp[0].evalI(c)
	case hitDef_forcestand:
		hd.forcestand = Btoi(exp[0].evalB(c))
	case hitDef_forcenofall:
		hd.forcenofall = exp[0].evalB(c)
	case hitDef_fall_damage:
		hd.fall.damage = exp[0].evalI(c)
	case hitDef_fall_xvelocity:
		hd.fall.xvelocity = exp[0].evalF(c)
	case hitDef_fall_yvelocity:
		hd.fall.yvelocity = exp[0].evalF(c)
	case hitDef_fall_recover:
		hd.fall.recover = exp[0].evalB(c)
	case hitDef_fall_recovertime:
		hd.fall.recovertime = exp[0].evalI(c)
	case hitDef_sparkno:
		n := exp[1].evalI(c)
		if n < 0 {
			hd.sparkno = IErr
		} else if exp[0].evalB(c) {
			hd.sparkno = ^n
		} else {
			hd.sparkno = n
		}
	case hitDef_guard_sparkno:
		n := exp[1].evalI(c)
		if n < 0 {
			hd.guard_sparkno = IErr
		} else if exp[0].evalB(c) {
			hd.guard_sparkno = ^n
		} else {
			hd.guard_sparkno = n
		}
	case hitDef_sparkxy:
		hd.sparkxy[0] = exp[0].evalF(c)
		if len(exp) > 1 {
			hd.sparkxy[1] = exp[1].evalF(c)
		}
	case hitDef_down_hittime:
		hd.down_hittime = exp[0].evalI(c)
	case hitDef_p1facing:
		hd.p1facing = exp[0].evalI(c)
	case hitDef_p1getp2facing:
		hd.p1getp2facing = exp[0].evalI(c)
	case hitDef_mindist:
		hd.mindist[0] = exp[0].evalF(c)
		if len(exp) > 1 {
			hd.mindist[1] = exp[1].evalF(c)
			if len(exp) > 2 {
				exp[2].run(c)
			}
		}
	case hitDef_maxdist:
		hd.maxdist[0] = exp[0].evalF(c)
		if len(exp) > 1 {
			hd.maxdist[1] = exp[1].evalF(c)
			if len(exp) > 2 {
				exp[2].run(c)
			}
		}
	case hitDef_snap:
		hd.snap[0] = exp[0].evalF(c)
		if len(exp) > 1 {
			hd.snap[1] = exp[1].evalF(c)
			if len(exp) > 2 {
				exp[2].run(c)
				if len(exp) > 3 {
					hd.snapt = exp[3].evalI(c)
				}
			}
		}
	case hitDef_p2facing:
		hd.p2facing = exp[0].evalI(c)
	case hitDef_air_hittime:
		hd.air_hittime = exp[0].evalI(c)
	case hitDef_fall:
		hd.ground_fall = exp[0].evalB(c)
		hd.air_fall = hd.ground_fall
	case hitDef_air_fall:
		hd.air_fall = exp[0].evalB(c)
	case hitDef_air_cornerpush_veloff:
		hd.air_cornerpush_veloff = exp[0].evalF(c)
	case hitDef_down_bounce:
		hd.down_bounce = exp[0].evalB(c)
	case hitDef_down_velocity:
		hd.down_velocity[0] = exp[0].evalF(c)
		if len(exp) > 1 {
			hd.down_velocity[1] = exp[1].evalF(c)
		}
	case hitDef_down_cornerpush_veloff:
		hd.down_cornerpush_veloff = exp[0].evalF(c)
	case hitDef_ground_hittime:
		hd.ground_hittime = exp[0].evalI(c)
		hd.guard_hittime = hd.ground_hittime
	case hitDef_guard_hittime:
		hd.guard_hittime = exp[0].evalI(c)
	case hitDef_guard_dist:
		hd.guard_dist = exp[0].evalI(c)
	case hitDef_pausetime:
		hd.pausetime = exp[0].evalI(c)
		hd.guard_pausetime = hd.pausetime
		if len(exp) > 1 {
			hd.shaketime = exp[1].evalI(c)
			hd.guard_shaketime = hd.shaketime
		}
	case hitDef_guard_pausetime:
		hd.guard_pausetime = exp[0].evalI(c)
		if len(exp) > 1 {
			hd.guard_shaketime = exp[1].evalI(c)
		}
	case hitDef_air_velocity:
		hd.air_velocity[0] = exp[0].evalF(c)
		if len(exp) > 1 {
			hd.air_velocity[1] = exp[1].evalF(c)
		}
	case hitDef_airguard_velocity:
		hd.airguard_velocity[0] = exp[0].evalF(c)
		if len(exp) > 1 {
			hd.airguard_velocity[1] = exp[1].evalF(c)
		}
	case hitDef_ground_slidetime:
		hd.ground_slidetime = exp[0].evalI(c)
		hd.guard_slidetime = hd.ground_slidetime
		hd.guard_ctrltime = hd.ground_slidetime
		hd.airguard_ctrltime = hd.ground_slidetime
	case hitDef_guard_slidetime:
		hd.guard_slidetime = exp[0].evalI(c)
		hd.guard_ctrltime = hd.guard_slidetime
		hd.airguard_ctrltime = hd.guard_slidetime
	case hitDef_guard_ctrltime:
		hd.guard_ctrltime = exp[0].evalI(c)
		hd.airguard_ctrltime = hd.guard_ctrltime
	case hitDef_airguard_ctrltime:
		hd.airguard_ctrltime = exp[0].evalI(c)
	case hitDef_ground_velocity_x:
		hd.ground_velocity[0] = exp[0].evalF(c)
	case hitDef_ground_velocity_y:
		hd.ground_velocity[1] = exp[0].evalF(c)
	case hitDef_guard_velocity:
		hd.guard_velocity = exp[0].evalF(c)
	case hitDef_ground_cornerpush_veloff:
		hd.ground_cornerpush_veloff = exp[0].evalF(c)
	case hitDef_guard_cornerpush_veloff:
		hd.guard_cornerpush_veloff = exp[0].evalF(c)
	case hitDef_airguard_cornerpush_veloff:
		hd.airguard_cornerpush_veloff = exp[0].evalF(c)
	case hitDef_yaccel:
		hd.yaccel = exp[0].evalF(c)
	case hitDef_envshake_time:
		hd.envshake_time = exp[0].evalI(c)
	case hitDef_envshake_ampl:
		hd.envshake_ampl = exp[0].evalI(c)
	case hitDef_envshake_phase:
		hd.envshake_phase = exp[0].evalF(c)
	case hitDef_envshake_freq:
		hd.envshake_freq = MaxF(0, exp[0].evalF(c))
	case hitDef_fall_envshake_time:
		hd.fall.envshake_time = exp[0].evalI(c)
	case hitDef_fall_envshake_ampl:
		hd.fall.envshake_ampl = exp[0].evalI(c)
	case hitDef_fall_envshake_phase:
		hd.fall.envshake_phase = exp[0].evalF(c)
	case hitDef_fall_envshake_freq:
		hd.fall.envshake_freq = MaxF(0, exp[0].evalF(c))
	default:
		if !palFX(sc).runSub(c, &hd.palfx, id, exp) {
			return false
		}
	}
	return true
}
func (sc hitDef) Run(c *Char, _ []int32) bool {
	c.hitdef.clear()
	c.hitdef.playerNo = sys.workingState.playerNo
	c.hitdef.sparkno = ^c.gi().data.sparkno
	c.hitdef.guard_sparkno = ^c.gi().data.guard.sparkno
	StateControllerBase(sc).run(c, func(id byte, exp []BytecodeExp) bool {
		sc.runSub(c, &c.hitdef, id, exp)
		return true
	})
	c.setHitdefDefault(&c.hitdef, false)
	return false
}

type reversalDef hitDef

const (
	reversalDef_reversal_attr = iota + hitDef_last + 1
)

func (sc reversalDef) Run(c *Char, _ []int32) bool {
	c.hitdef.clear()
	c.hitdef.playerNo = sys.workingState.playerNo
	StateControllerBase(sc).run(c, func(id byte, exp []BytecodeExp) bool {
		switch id {
		case reversalDef_reversal_attr:
			c.hitdef.reversal_attr = exp[0].evalI(c)
		default:
			hitDef(sc).runSub(c, &c.hitdef, id, exp)
		}
		return true
	})
	c.setHitdefDefault(&c.hitdef, false)
	return false
}

type projectile hitDef

const (
	projectile_postype = iota + hitDef_last + 1
	projectile_projid
	projectile_projremove
	projectile_projremovetime
	projectile_projshadow
	projectile_projmisstime
	projectile_projhits
	projectile_projpriority
	projectile_projhitanim
	projectile_projremanim
	projectile_projcancelanim
	projectile_velocity
	projectile_velmul
	projectile_remvelocity
	projectile_accel
	projectile_projscale
	projectile_offset
	projectile_projsprpriority
	projectile_projstagebound
	projectile_projedgebound
	projectile_projheightbound
	projectile_projanim
	projectile_supermovetime
	projectile_pausemovetime
	projectile_ownpal
	projectile_remappal
)

func (sc projectile) Run(c *Char, _ []int32) bool {
	p := c.newProj()
	if p == nil {
		return false
	}
	p.hitdef.playerNo = sys.workingState.playerNo
	pt := PT_P1
	var x, y float32 = 0, 0
	op := false
	rp := [...]int32{-1, 0}
	StateControllerBase(sc).run(c, func(id byte, exp []BytecodeExp) bool {
		switch id {
		case projectile_postype:
			pt = PosType(exp[0].evalI(c))
		case projectile_projid:
			p.id = exp[0].evalI(c)
		case projectile_projremove:
			p.remove = exp[0].evalB(c)
		case projectile_projremovetime:
			p.removetime = exp[0].evalI(c)
		case projectile_projshadow:
			p.shadow[0] = exp[0].evalI(c)
			if len(exp) > 1 {
				p.shadow[1] = exp[1].evalI(c)
				if len(exp) > 2 {
					p.shadow[2] = exp[2].evalI(c)
				}
			}
		case projectile_projmisstime:
			p.misstime = exp[0].evalI(c)
		case projectile_projhits:
			p.hits = exp[0].evalI(c)
		case projectile_projpriority:
			p.priority = exp[0].evalI(c)
		case projectile_projhitanim:
			p.hitanim = exp[0].evalI(c)
		case projectile_projremanim:
			p.remanim = Max(-1, exp[0].evalI(c))
		case projectile_projcancelanim:
			p.cancelanim = Max(-1, exp[0].evalI(c))
		case projectile_velocity:
			p.velocity[0] = exp[0].evalF(c)
			if len(exp) > 1 {
				p.velocity[1] = exp[1].evalF(c)
			}
		case projectile_velmul:
			p.velmul[0] = exp[0].evalF(c)
			if len(exp) > 1 {
				p.velmul[1] = exp[1].evalF(c)
			}
		case projectile_remvelocity:
			p.remvelocity[0] = exp[0].evalF(c)
			if len(exp) > 1 {
				p.remvelocity[1] = exp[1].evalF(c)
			}
		case projectile_accel:
			p.accel[0] = exp[0].evalF(c)
			if len(exp) > 1 {
				p.accel[1] = exp[1].evalF(c)
			}
		case projectile_projscale:
			p.scale[0] = exp[0].evalF(c)
			if len(exp) > 1 {
				p.scale[1] = exp[1].evalF(c)
			}
		case projectile_offset:
			x = exp[0].evalF(c)
			if len(exp) > 1 {
				y = exp[1].evalF(c)
			}
		case projectile_projsprpriority:
			p.sprpriority = exp[0].evalI(c)
		case projectile_projstagebound:
			p.stagebound = exp[0].evalI(c)
		case projectile_projedgebound:
			p.edgebound = exp[0].evalI(c)
		case projectile_projheightbound:
			p.heightbound[0] = exp[0].evalI(c)
			if len(exp) > 1 {
				p.heightbound[1] = exp[1].evalI(c)
			}
		case projectile_projanim:
			p.anim = exp[0].evalI(c)
		case projectile_supermovetime:
			p.supermovetime = exp[0].evalI(c)
		case projectile_pausemovetime:
			p.pausemovetime = exp[0].evalI(c)
		case projectile_ownpal:
			op = exp[0].evalB(c)
		case projectile_remappal:
			rp[0] = exp[0].evalI(c)
			if len(exp) > 1 {
				rp[1] = exp[1].evalI(c)
			}
		default:
			if !hitDef(sc).runSub(c, &p.hitdef, id, exp) {
				afterImage(sc).runSub(c, &p.aimg, id, exp)
			}
		}
		return true
	})
	c.setHitdefDefault(&p.hitdef, true)
	if p.remanim == IErr {
		p.remanim = p.hitanim
	}
	if p.cancelanim == IErr {
		p.cancelanim = p.remanim
	}
	if p.aimg.time != 0 {
		p.aimg.setupPalFX()
	}
	p.localscl = c.localscl
	c.projInit(p, pt, x, y, op, rp[0], rp[1])
	return false
}

type width StateControllerBase

const (
	width_edge byte = iota
	width_player
	width_value
)

func (sc width) Run(c *Char, _ []int32) bool {
	StateControllerBase(sc).run(c, func(id byte, exp []BytecodeExp) bool {
		switch id {
		case width_edge:
			c.setFEdge(exp[0].evalF(c))
			if len(exp) > 1 {
				c.setBEdge(exp[1].evalF(c))
			}
		case width_player:
			c.setFWidth(exp[0].evalF(c))
			if len(exp) > 1 {
				c.setBWidth(exp[1].evalF(c))
			}
		case width_value:
			v1 := exp[0].evalF(c)
			c.setFEdge(v1)
			c.setFWidth(v1)
			if len(exp) > 1 {
				v2 := exp[1].evalF(c)
				c.setBEdge(v2)
				c.setBWidth(v2)
			}
		}
		return true
	})
	return false
}

type sprPriority StateControllerBase

const (
	sprPriority_value byte = iota
)

func (sc sprPriority) Run(c *Char, _ []int32) bool {
	StateControllerBase(sc).run(c, func(id byte, exp []BytecodeExp) bool {
		switch id {
		case sprPriority_value:
			c.setSprPriority(exp[0].evalI(c))
		}
		return true
	})
	return false
}

type varSet StateControllerBase

const (
	varSet_ byte = iota
)

func (sc varSet) Run(c *Char, _ []int32) bool {
	StateControllerBase(sc).run(c, func(id byte, exp []BytecodeExp) bool {
		switch id {
		case varSet_:
			exp[0].run(c)
		}
		return true
	})
	return false
}

type turn StateControllerBase

const (
	turn_ byte = iota
)

func (sc turn) Run(c *Char, _ []int32) bool {
	StateControllerBase(sc).run(c, func(id byte, exp []BytecodeExp) bool {
		switch id {
		case turn_:
			c.setFacing(-c.facing)
		}
		return true
	})
	return false
}

type targetFacing StateControllerBase

const (
	targetFacing_id byte = iota
	targetFacing_value
)

func (sc targetFacing) Run(c *Char, _ []int32) bool {
	tar := c.getTarget(-1)
	if len(tar) == 0 {
		return false
	}
	StateControllerBase(sc).run(c, func(id byte, exp []BytecodeExp) bool {
		switch id {
		case targetFacing_id:
			tar = c.getTarget(exp[0].evalI(c))
			if len(tar) == 0 {
				return false
			}
		case targetFacing_value:
			c.targetFacing(tar, exp[0].evalI(c))
		}
		return true
	})
	return false
}

type targetBind StateControllerBase

const (
	targetBind_id byte = iota
	targetBind_time
	targetBind_pos
)

func (sc targetBind) Run(c *Char, _ []int32) bool {
	tar := c.getTarget(-1)
	if len(tar) == 0 {
		return false
	}
	t := int32(1)
	var x, y float32 = 0, 0
	StateControllerBase(sc).run(c, func(id byte, exp []BytecodeExp) bool {
		switch id {
		case targetBind_id:
			tar = c.getTarget(exp[0].evalI(c))
			if len(tar) == 0 {
				return false
			}
		case targetBind_time:
			t = exp[0].evalI(c)
		case targetBind_pos:
			x = exp[0].evalF(c)
			if len(exp) > 1 {
				y = exp[1].evalF(c)
			}
		}
		return true
	})
	c.targetBind(tar, t, x, y)
	return false
}

type bindToTarget StateControllerBase

const (
	bindToTarget_id byte = iota
	bindToTarget_time
	bindToTarget_pos
)

func (sc bindToTarget) Run(c *Char, _ []int32) bool {
	tar := c.getTarget(-1)
	if len(tar) == 0 {
		return false
	}
	t, x, y, hmf := int32(1), float32(0), float32(math.NaN()), HMF_F
	StateControllerBase(sc).run(c, func(id byte, exp []BytecodeExp) bool {
		switch id {
		case bindToTarget_id:
			tar = c.getTarget(exp[0].evalI(c))
			if len(tar) == 0 {
				return false
			}
		case bindToTarget_time:
			t = exp[0].evalI(c)
		case bindToTarget_pos:
			x = exp[0].evalF(c)
			if len(exp) > 1 {
				y = exp[1].evalF(c)
				if len(exp) > 2 {
					hmf = HMF(exp[2].evalI(c))
				}
			}
		}
		return true
	})
	c.bindToTarget(tar, t, x, y, hmf)
	return false
}

type targetLifeAdd StateControllerBase

const (
	targetLifeAdd_id byte = iota
	targetLifeAdd_absolute
	targetLifeAdd_kill
	targetLifeAdd_value
)

func (sc targetLifeAdd) Run(c *Char, _ []int32) bool {
	tar, a, k := c.getTarget(-1), false, true
	if len(tar) == 0 {
		return false
	}
	StateControllerBase(sc).run(c, func(id byte, exp []BytecodeExp) bool {
		switch id {
		case targetLifeAdd_id:
			tar = c.getTarget(exp[0].evalI(c))
			if len(tar) == 0 {
				return false
			}
		case targetLifeAdd_absolute:
			a = exp[0].evalB(c)
		case targetLifeAdd_kill:
			k = exp[0].evalB(c)
		case targetLifeAdd_value:
			c.targetLifeAdd(tar, exp[0].evalI(c), k, a)
		}
		return true
	})
	return false
}

type targetState StateControllerBase

const (
	targetState_id byte = iota
	targetState_value
)

func (sc targetState) Run(c *Char, _ []int32) bool {
	tar := c.getTarget(-1)
	if len(tar) == 0 {
		return false
	}
	StateControllerBase(sc).run(c, func(id byte, exp []BytecodeExp) bool {
		switch id {
		case targetState_id:
			tar = c.getTarget(exp[0].evalI(c))
			if len(tar) == 0 {
				return false
			}
		case targetState_value:
			c.targetState(tar, exp[0].evalI(c))
		}
		return true
	})
	return false
}

type targetVelSet StateControllerBase

const (
	targetVelSet_id byte = iota
	targetVelSet_x
	targetVelSet_y
)

func (sc targetVelSet) Run(c *Char, _ []int32) bool {
	tar := c.getTarget(-1)
	if len(tar) == 0 {
		return false
	}
	StateControllerBase(sc).run(c, func(id byte, exp []BytecodeExp) bool {
		switch id {
		case targetVelSet_id:
			tar = c.getTarget(exp[0].evalI(c))
			if len(tar) == 0 {
				return false
			}
		case targetVelSet_x:
			c.targetVelSetX(tar, exp[0].evalF(c))
		case targetVelSet_y:
			c.targetVelSetY(tar, exp[0].evalF(c))
		}
		return true
	})
	return false
}

type targetVelAdd StateControllerBase

const (
	targetVelAdd_id byte = iota
	targetVelAdd_x
	targetVelAdd_y
)

func (sc targetVelAdd) Run(c *Char, _ []int32) bool {
	tar := c.getTarget(-1)
	if len(tar) == 0 {
		return false
	}
	StateControllerBase(sc).run(c, func(id byte, exp []BytecodeExp) bool {
		switch id {
		case targetVelAdd_id:
			tar = c.getTarget(exp[0].evalI(c))
			if len(tar) == 0 {
				return false
			}
		case targetVelAdd_x:
			c.targetVelAddX(tar, exp[0].evalF(c))
		case targetVelAdd_y:
			c.targetVelAddY(tar, exp[0].evalF(c))
		}
		return true
	})
	return false
}

type targetPowerAdd StateControllerBase

const (
	targetPowerAdd_id byte = iota
	targetPowerAdd_value
)

func (sc targetPowerAdd) Run(c *Char, _ []int32) bool {
	tar := c.getTarget(-1)
	if len(tar) == 0 {
		return false
	}
	StateControllerBase(sc).run(c, func(id byte, exp []BytecodeExp) bool {
		switch id {
		case targetPowerAdd_id:
			tar = c.getTarget(exp[0].evalI(c))
			if len(tar) == 0 {
				return false
			}
		case targetPowerAdd_value:
			c.targetPowerAdd(tar, exp[0].evalI(c))
		}
		return true
	})
	return false
}

type targetDrop StateControllerBase

const (
	targetDrop_excludeid byte = iota
	targetDrop_keepone
)

func (sc targetDrop) Run(c *Char, _ []int32) bool {
	tar, eid, ko := c.getTarget(-1), int32(-1), true
	if len(tar) == 0 {
		return false
	}
	StateControllerBase(sc).run(c, func(id byte, exp []BytecodeExp) bool {
		switch id {
		case targetDrop_excludeid:
			eid = exp[0].evalI(c)
		case targetDrop_keepone:
			ko = exp[0].evalB(c)
		}
		return true
	})
	c.targetDrop(eid, ko)
	return false
}

type lifeAdd StateControllerBase

const (
	lifeAdd_absolute byte = iota
	lifeAdd_kill
	lifeAdd_value
)

func (sc lifeAdd) Run(c *Char, _ []int32) bool {
	a, k := false, true
	StateControllerBase(sc).run(c, func(id byte, exp []BytecodeExp) bool {
		switch id {
		case lifeAdd_absolute:
			a = exp[0].evalB(c)
		case lifeAdd_kill:
			k = exp[0].evalB(c)
		case lifeAdd_value:
			c.lifeAdd(float64(exp[0].evalI(c)), k, a)
		}
		return true
	})
	return false
}

type lifeSet StateControllerBase

const (
	lifeSet_value byte = iota
)

func (sc lifeSet) Run(c *Char, _ []int32) bool {
	StateControllerBase(sc).run(c, func(id byte, exp []BytecodeExp) bool {
		switch id {
		case lifeSet_value:
			c.lifeSet(exp[0].evalI(c))
		}
		return true
	})
	return false
}

type powerAdd StateControllerBase

const (
	powerAdd_value byte = iota
)

func (sc powerAdd) Run(c *Char, _ []int32) bool {
	StateControllerBase(sc).run(c, func(id byte, exp []BytecodeExp) bool {
		switch id {
		case powerAdd_value:
			c.powerAdd(exp[0].evalI(c))
		}
		return true
	})
	return false
}

type powerSet StateControllerBase

const (
	powerSet_value byte = iota
)

func (sc powerSet) Run(c *Char, _ []int32) bool {
	StateControllerBase(sc).run(c, func(id byte, exp []BytecodeExp) bool {
		switch id {
		case powerSet_value:
			c.powerSet(exp[0].evalI(c))
		}
		return true
	})
	return false
}

type hitVelSet StateControllerBase

const (
	hitVelSet_x byte = iota
	hitVelSet_y
)

func (sc hitVelSet) Run(c *Char, _ []int32) bool {
	StateControllerBase(sc).run(c, func(id byte, exp []BytecodeExp) bool {
		switch id {
		case hitVelSet_x:
			if exp[0].evalB(c) {
				c.hitVelSetX()
			}
		case hitVelSet_y:
			if exp[0].evalB(c) {
				c.hitVelSetY()
			}
		}
		return true
	})
	return false
}

type screenBound StateControllerBase

const (
	screenBound_value byte = iota
	screenBound_movecamera
)

func (sc screenBound) Run(c *Char, _ []int32) bool {
	StateControllerBase(sc).run(c, func(id byte, exp []BytecodeExp) bool {
		switch id {
		case screenBound_value:
			if exp[0].evalB(c) {
				c.setSF(CSF_screenbound)
			} else {
				c.unsetSF(CSF_screenbound)
			}
		case screenBound_movecamera:
			if exp[0].evalB(c) {
				c.setSF(CSF_movecamera_x)
			} else {
				c.unsetSF(CSF_movecamera_x)
			}
			if len(exp) > 1 {
				if exp[1].evalB(c) {
					c.setSF(CSF_movecamera_y)
				} else {
					c.unsetSF(CSF_movecamera_y)
				}
			}
		}
		return true
	})
	return false
}

type posFreeze StateControllerBase

const (
	posFreeze_value byte = iota
)

func (sc posFreeze) Run(c *Char, _ []int32) bool {
	StateControllerBase(sc).run(c, func(id byte, exp []BytecodeExp) bool {
		switch id {
		case posFreeze_value:
			if exp[0].evalB(c) {
				c.setSF(CSF_posfreeze)
			}
		}
		return true
	})
	return false
}

type envShake StateControllerBase

const (
	envShake_time byte = iota
	envShake_ampl
	envShake_phase
	envShake_freq
)

func (sc envShake) Run(c *Char, _ []int32) bool {
	sys.envShake.clear()
	StateControllerBase(sc).run(c, func(id byte, exp []BytecodeExp) bool {
		switch id {
		case envShake_time:
			sys.envShake.time = exp[0].evalI(c)
		case envShake_ampl:
			sys.envShake.ampl = int32(float32(exp[0].evalI(c)) * c.localscl)
		case envShake_phase:
			sys.envShake.phase = MaxF(0, exp[0].evalF(c)*float32(math.Pi)/180) * c.localscl
		case envShake_freq:
			sys.envShake.freq = exp[0].evalF(c)
		}
		return true
	})
	sys.envShake.setDefPhase()
	return false
}

type hitOverride StateControllerBase

const (
	hitOverride_attr byte = iota
	hitOverride_slot
	hitOverride_stateno
	hitOverride_time
	hitOverride_forceair
)

func (sc hitOverride) Run(c *Char, _ []int32) bool {
	var a, s, st, t int32 = 0, 0, -1, 1
	f := false
	StateControllerBase(sc).run(c, func(id byte, exp []BytecodeExp) bool {
		switch id {
		case hitOverride_attr:
			a = exp[0].evalI(c)
		case hitOverride_slot:
			s = Max(0, exp[0].evalI(c))
			if s > 7 {
				s = 0
			}
		case hitOverride_stateno:
			st = exp[0].evalI(c)
		case hitOverride_time:
			t = exp[0].evalI(c)
			if t < -1 || t == 0 {
				t = 1
			}
		case hitOverride_forceair:
			f = exp[0].evalB(c)
		}
		return true
	})
	if st < 0 {
		t = 0
	}
	c.ho[s] = HitOverride{attr: a, stateno: st, time: t, forceair: f,
		playerNo: sys.workingState.playerNo}
	return false
}

type pause StateControllerBase

const (
	pause_time byte = iota
	pause_movetime
	pause_pausebg
	pause_endcmdbuftime
)

func (sc pause) Run(c *Char, _ []int32) bool {
	var t, mt int32 = 0, 0
	sys.pausebg, sys.pauseendcmdbuftime = true, 0
	StateControllerBase(sc).run(c, func(id byte, exp []BytecodeExp) bool {
		switch id {
		case pause_time:
			t = exp[0].evalI(c)
		case pause_movetime:
			mt = exp[0].evalI(c)
		case pause_pausebg:
			sys.pausebg = exp[0].evalB(c)
		case pause_endcmdbuftime:
			sys.pauseendcmdbuftime = exp[0].evalI(c)
		}
		return true
	})
	c.setPauseTime(t, mt)
	return false
}

type superPause StateControllerBase

const (
	superPause_time byte = iota
	superPause_movetime
	superPause_pausebg
	superPause_endcmdbuftime
	superPause_darken
	superPause_anim
	superPause_pos
	superPause_p2defmul
	superPause_poweradd
	superPause_unhittable
	superPause_sound
)

func (sc superPause) Run(c *Char, _ []int32) bool {
	var t, mt int32 = 30, 0
	uh := true
	sys.superanim, sys.superpmap.remap = c.getAnim(30, true), nil
	sys.superpos, sys.superfacing = [...]float32{c.pos[0] * c.localscl, c.pos[1] * c.localscl}, c.facing
	sys.superpausebg, sys.superendcmdbuftime, sys.superdarken = true, 0, true
	sys.superp2defmul = sys.super_TargetDefenceMul
	StateControllerBase(sc).run(c, func(id byte, exp []BytecodeExp) bool {
		switch id {
		case superPause_time:
			t = exp[0].evalI(c)
		case superPause_movetime:
			mt = exp[0].evalI(c)
		case superPause_pausebg:
			sys.superpausebg = exp[0].evalB(c)
		case superPause_endcmdbuftime:
			sys.superendcmdbuftime = exp[0].evalI(c)
		case superPause_darken:
			sys.superdarken = exp[0].evalB(c)
		case superPause_anim:
			f := exp[0].evalB(c)
			if sys.superanim = c.getAnim(exp[1].evalI(c), f); sys.superanim != nil {
				if f {
					sys.superpmap.remap = nil
				} else {
					sys.superpmap.remap = c.getPalMap()
				}
			}
		case superPause_pos:
			sys.superpos[0] += c.facing * exp[0].evalF(c) * c.localscl
			if len(exp) > 1 {
				sys.superpos[1] += exp[1].evalF(c) * c.localscl
			}
		case superPause_p2defmul:
			sys.superp2defmul = exp[0].evalF(c)
			if sys.superp2defmul == 0 {
				sys.superp2defmul = sys.super_TargetDefenceMul
			}
		case superPause_poweradd:
			c.powerAdd(exp[0].evalI(c))
		case superPause_unhittable:
			uh = exp[0].evalB(c)
		case superPause_sound:
			n := int32(0)
			if len(exp) > 2 {
				n = exp[2].evalI(c)
			}
			vo := int32(0)
			if c.gi().ver[0] == 1 {
				vo = 100
			}
			c.playSound(exp[0].evalB(c), false, false, exp[1].evalI(c), n, -1,
				vo, 0, 1, &c.pos[0])
		}
		return true
	})
	c.setSuperPauseTime(t, mt, uh)
	return false
}

type trans StateControllerBase

const (
	trans_trans byte = iota
)

func (sc trans) Run(c *Char, _ []int32) bool {
	c.alpha[1] = 255
	StateControllerBase(sc).run(c, func(id byte, exp []BytecodeExp) bool {
		switch id {
		case trans_trans:
			c.alpha[0] = exp[0].evalI(c)
			c.alpha[1] = exp[1].evalI(c)
			if len(exp) >= 3 {
				c.alpha[0] = Max(0, Min(255, c.alpha[0]))
				c.alpha[1] = Max(0, Min(255, c.alpha[1]))
				if len(exp) >= 4 {
					c.alpha[1] = ^c.alpha[1]
				} else if c.alpha[0] == 1 && c.alpha[1] == 255 {
					c.alpha[0] = 0
				}
			}
		}
		return true
	})
	c.setSF(CSF_trans)
	return false
}

type playerPush StateControllerBase

const (
	playerPush_value byte = iota
)

func (sc playerPush) Run(c *Char, _ []int32) bool {
	StateControllerBase(sc).run(c, func(id byte, exp []BytecodeExp) bool {
		switch id {
		case playerPush_value:
			if exp[0].evalB(c) {
				c.setSF(CSF_playerpush)
			} else {
				c.unsetSF(CSF_playerpush)
			}
		}
		return true
	})
	return false
}

type stateTypeSet StateControllerBase

const (
	stateTypeSet_statetype byte = iota
	stateTypeSet_movetype
	stateTypeSet_physics
)

func (sc stateTypeSet) Run(c *Char, _ []int32) bool {
	StateControllerBase(sc).run(c, func(id byte, exp []BytecodeExp) bool {
		switch id {
		case stateTypeSet_statetype:
			c.ss.stateType = StateType(exp[0].evalI(c))
		case stateTypeSet_movetype:
			c.ss.moveType = MoveType(exp[0].evalI(c))
		case stateTypeSet_physics:
			c.ss.physics = StateType(exp[0].evalI(c))
		}
		return true
	})
	return false
}

type angleDraw StateControllerBase

const (
	angleDraw_value byte = iota
	angleDraw_scale
)

func (sc angleDraw) Run(c *Char, _ []int32) bool {
	c.setSF(CSF_angledraw)
	StateControllerBase(sc).run(c, func(id byte, exp []BytecodeExp) bool {
		switch id {
		case angleDraw_value:
			c.angleSet(exp[0].evalF(c))
		case angleDraw_scale:
			c.angleScalse[0] *= exp[0].evalF(c)
			if len(exp) > 1 {
				c.angleScalse[1] *= exp[1].evalF(c)
			}
		}
		return true
	})
	return false
}

type angleSet StateControllerBase

const (
	angleSet_value byte = iota
)

func (sc angleSet) Run(c *Char, _ []int32) bool {
	StateControllerBase(sc).run(c, func(id byte, exp []BytecodeExp) bool {
		switch id {
		case angleSet_value:
			c.angleSet(exp[0].evalF(c))
		}
		return true
	})
	return false
}

type angleAdd StateControllerBase

const (
	angleAdd_value byte = iota
)

func (sc angleAdd) Run(c *Char, _ []int32) bool {
	StateControllerBase(sc).run(c, func(id byte, exp []BytecodeExp) bool {
		switch id {
		case angleAdd_value:
			c.angleSet(c.angle + exp[0].evalF(c))
		}
		return true
	})
	return false
}

type angleMul StateControllerBase

const (
	angleMul_value byte = iota
)

func (sc angleMul) Run(c *Char, _ []int32) bool {
	StateControllerBase(sc).run(c, func(id byte, exp []BytecodeExp) bool {
		switch id {
		case angleMul_value:
			c.angleSet(c.angle * exp[0].evalF(c))
		}
		return true
	})
	return false
}

type envColor StateControllerBase

const (
	envColor_value byte = iota
	envColor_time
	envColor_under
)

func (sc envColor) Run(c *Char, _ []int32) bool {
	sys.envcol = [...]int32{255, 255, 255}
	sys.envcol_time = 1
	sys.envcol_under = false
	StateControllerBase(sc).run(c, func(id byte, exp []BytecodeExp) bool {
		switch id {
		case envColor_value:
			sys.envcol[0] = exp[0].evalI(c)
			sys.envcol[1] = exp[1].evalI(c)
			sys.envcol[2] = exp[2].evalI(c)
		case envColor_time:
			sys.envcol_time = exp[0].evalI(c)
		case envColor_under:
			sys.envcol_under = exp[0].evalB(c)
		}
		return true
	})
	return false
}

type displayToClipboard StateControllerBase

const (
	displayToClipboard_params byte = iota
	displayToClipboard_text
)

func (sc displayToClipboard) Run(c *Char, _ []int32) bool {
	params := []interface{}{}
	StateControllerBase(sc).run(c, func(id byte, exp []BytecodeExp) bool {
		switch id {
		case displayToClipboard_params:
			for _, e := range exp {
				if bv := e.run(c); bv.t == VT_Float {
					params = append(params, bv.ToF())
				} else {
					params = append(params, bv.ToI())
				}
			}
		case displayToClipboard_text:
			sys.clipboardText[sys.workingState.playerNo] = nil
			sys.appendToClipboard(sys.workingState.playerNo,
				int(exp[0].evalI(c)), params...)
		}
		return true
	})
	return false
}

type appendToClipboard displayToClipboard

func (sc appendToClipboard) Run(c *Char, _ []int32) bool {
	params := []interface{}{}
	StateControllerBase(sc).run(c, func(id byte, exp []BytecodeExp) bool {
		switch id {
		case displayToClipboard_params:
			for _, e := range exp {
				if bv := e.run(c); bv.t == VT_Float {
					params = append(params, bv.ToF())
				} else {
					params = append(params, bv.ToI())
				}
			}
		case displayToClipboard_text:
			sys.appendToClipboard(sys.workingState.playerNo,
				int(exp[0].evalI(c)), params...)
		}
		return true
	})
	return false
}

type clearClipboard StateControllerBase

const (
	clearClipboard_ byte = iota
)

func (sc clearClipboard) Run(c *Char, _ []int32) bool {
	StateControllerBase(sc).run(c, func(id byte, exp []BytecodeExp) bool {
		switch id {
		case clearClipboard_:
			sys.clipboardText[sys.workingState.playerNo] = nil
		}
		return true
	})
	return false
}

type makeDust StateControllerBase

const (
	makeDust_spacing byte = iota
	makeDust_pos
	makeDust_pos2
)

func (sc makeDust) Run(c *Char, _ []int32) bool {
	StateControllerBase(sc).run(c, func(id byte, exp []BytecodeExp) bool {
		switch id {
		case makeDust_spacing:
			s := Max(1, exp[0].evalI(c))
			if c.time()%s != s-1 {
				return false
			}
		case makeDust_pos:
			x, y := exp[0].evalF(c), float32(0)
			if len(exp) > 1 {
				y = exp[1].evalF(c)
			}
			c.makeDust(x-float32(c.size.draw.offset[0]),
				y-float32(c.size.draw.offset[1]))
		case makeDust_pos2:
			x, y := exp[0].evalF(c), float32(0)
			if len(exp) > 1 {
				y = exp[1].evalF(c)
			}
			c.makeDust(x-float32(c.size.draw.offset[0]),
				y-float32(c.size.draw.offset[1]))
		}
		return true
	})
	return false
}

type attackDist StateControllerBase

const (
	attackDist_value byte = iota
)

func (sc attackDist) Run(c *Char, _ []int32) bool {
	StateControllerBase(sc).run(c, func(id byte, exp []BytecodeExp) bool {
		switch id {
		case attackDist_value:
			c.attackDist = exp[0].evalF(c)
		}
		return true
	})
	return false
}

type attackMulSet StateControllerBase

const (
	attackMulSet_value byte = iota
)

func (sc attackMulSet) Run(c *Char, _ []int32) bool {
	StateControllerBase(sc).run(c, func(id byte, exp []BytecodeExp) bool {
		switch id {
		case attackMulSet_value:
			c.attackMul = float32(c.gi().data.attack) / 100 * exp[0].evalF(c)
		}
		return true
	})
	return false
}

type defenceMulSet StateControllerBase

const (
	defenceMulSet_value byte = iota
)

func (sc defenceMulSet) Run(c *Char, _ []int32) bool {
	StateControllerBase(sc).run(c, func(id byte, exp []BytecodeExp) bool {
		switch id {
		case defenceMulSet_value:
			c.defenceMul = float32(c.gi().data.defence) / (exp[0].evalF(c) * 100)
		}
		return true
	})
	return false
}

type fallEnvShake StateControllerBase

const (
	fallEnvShake_ byte = iota
)

func (sc fallEnvShake) Run(c *Char, _ []int32) bool {
	StateControllerBase(sc).run(c, func(id byte, exp []BytecodeExp) bool {
		switch id {
		case fallEnvShake_:
			sys.envShake = EnvShake{time: c.ghv.fall.envshake_time,
				freq: c.ghv.fall.envshake_freq * math.Pi / 180,
				ampl: c.ghv.fall.envshake_ampl, phase: c.ghv.fall.envshake_phase}
			sys.envShake.setDefPhase()
		}
		return true
	})
	return false
}

type hitFallDamage StateControllerBase

const (
	hitFallDamage_ byte = iota
)

func (sc hitFallDamage) Run(c *Char, _ []int32) bool {
	StateControllerBase(sc).run(c, func(id byte, exp []BytecodeExp) bool {
		switch id {
		case hitFallDamage_:
			c.hitFallDamage()
		}
		return true
	})
	return false
}

type hitFallVel StateControllerBase

const (
	hitFallVel_ byte = iota
)

func (sc hitFallVel) Run(c *Char, _ []int32) bool {
	StateControllerBase(sc).run(c, func(id byte, exp []BytecodeExp) bool {
		switch id {
		case hitFallVel_:
			c.hitFallVel()
		}
		return true
	})
	return false
}

type hitFallSet StateControllerBase

const (
	hitFallSet_value byte = iota
	hitFallSet_xvel
	hitFallSet_yvel
)

func (sc hitFallSet) Run(c *Char, _ []int32) bool {
	f, xv, yv := int32(-1), float32(math.NaN()), float32(math.NaN())
	StateControllerBase(sc).run(c, func(id byte, exp []BytecodeExp) bool {
		switch id {
		case hitFallSet_value:
			f = exp[0].evalI(c)
			if len(c.ghv.hitBy) == 0 {
				return false
			}
		case hitFallSet_xvel:
			xv = exp[0].evalF(c)
		case hitFallSet_yvel:
			yv = exp[0].evalF(c)
		}
		return true
	})
	c.hitFallSet(f, xv, yv)
	return false
}

type varRangeSet StateControllerBase

const (
	varRangeSet_first byte = iota
	varRangeSet_last
	varRangeSet_value
	varRangeSet_fvalue
)

func (sc varRangeSet) Run(c *Char, _ []int32) bool {
	var first, last int32 = 0, 0
	StateControllerBase(sc).run(c, func(id byte, exp []BytecodeExp) bool {
		switch id {
		case varRangeSet_first:
			first = exp[0].evalI(c)
		case varRangeSet_last:
			last = exp[0].evalI(c)
		case varRangeSet_value:
			v := exp[0].evalI(c)
			if first >= 0 && last < int32(NumVar) {
				for i := first; i <= last; i++ {
					c.ivar[i] = v
				}
			}
		case varRangeSet_fvalue:
			fv := exp[0].evalF(c)
			if first >= 0 && last < int32(NumFvar) {
				for i := first; i <= last; i++ {
					c.fvar[i] = fv
				}
			}
		}
		return true
	})
	return false
}

type remapPal StateControllerBase

const (
	remapPal_source byte = iota
	remapPal_dest
)

func (sc remapPal) Run(c *Char, _ []int32) bool {
	src := [...]int32{-1, -1}
	StateControllerBase(sc).run(c, func(id byte, exp []BytecodeExp) bool {
		switch id {
		case remapPal_source:
			src[0] = exp[0].evalI(c)
			if len(exp) > 1 {
				src[1] = exp[1].evalI(c)
			}
		case remapPal_dest:
			dst := [...]int32{exp[0].evalI(c), -1}
			if len(exp) > 1 {
				dst[1] = exp[1].evalI(c)
			}
			c.remapPal(c.getPalfx(), src, dst)
		}
		return true
	})
	return false
}

type stopSnd StateControllerBase

const (
	stopSnd_channel byte = iota
)

func (sc stopSnd) Run(c *Char, _ []int32) bool {
	StateControllerBase(sc).run(c, func(id byte, exp []BytecodeExp) bool {
		switch id {
		case stopSnd_channel:
			if ch := Min(255, exp[0].evalI(c)); ch < 0 {
				sys.stopAllSound()
			} else if int(ch) < len(c.sounds) {
				c.sounds[ch].sound = nil
			}
		}
		return true
	})
	return false
}

type sndPan StateControllerBase

const (
	sndPan_channel byte = iota
	sndPan_pan
	sndPan_abspan
)

func (sc sndPan) Run(c *Char, _ []int32) bool {
	ch, pan, x := int32(-1), float32(0), &c.pos[0]
	StateControllerBase(sc).run(c, func(id byte, exp []BytecodeExp) bool {
		switch id {
		case sndPan_channel:
			ch = exp[0].evalI(c)
		case sndPan_pan:
			pan = exp[0].evalF(c)
		case sndPan_abspan:
			pan = exp[0].evalF(c)
			x = nil
		}
		return true
	})
	if ch <= 0 && int(ch) < len(c.sounds) {
		c.sounds[ch].SetPan(pan, x)
	}
	return false
}

type varRandom StateControllerBase

const (
	varRandom_v byte = iota
	varRandom_range
)

func (sc varRandom) Run(c *Char, _ []int32) bool {
	var v int32
	StateControllerBase(sc).run(c, func(id byte, exp []BytecodeExp) bool {
		switch id {
		case varRandom_v:
			v = exp[0].evalI(c)
		case varRandom_range:
			var min, max int32 = 0, exp[0].evalI(c)
			if len(exp) > 1 {
				min, max = max, exp[1].evalI(c)
			}
			c.varSet(v, RandI(min, max))
		}
		return true
	})
	return false
}

type gravity StateControllerBase

const (
	gravity_ byte = iota
)

func (sc gravity) Run(c *Char, _ []int32) bool {
	StateControllerBase(sc).run(c, func(id byte, exp []BytecodeExp) bool {
		switch id {
		case gravity_:
			c.gravity()
		}
		return true
	})
	return false
}

type bindToParent StateControllerBase

const (
	bindToParent_time byte = iota
	bindToParent_facing
	bindToParent_pos
)

func (sc bindToParent) Run(c *Char, _ []int32) bool {
	p := c.parent()
	if p == nil {
		return false
	}
	c.bindTime, c.bindPos = 1, [2]float32{}
	c.setBindToId(p)
	StateControllerBase(sc).run(c, func(id byte, exp []BytecodeExp) bool {
		switch id {
		case bindToParent_time:
			c.setBindTime(exp[0].evalI(c))
		case bindToParent_facing:
			if f := exp[0].evalI(c); f < 0 {
				c.bindFacing = -1
			} else if f > 0 {
				c.bindFacing = 1
			}
		case bindToParent_pos:
			c.bindPos[0] = exp[0].evalF(c)
			if len(exp) > 1 {
				c.bindPos[1] = exp[1].evalF(c)
			}
		}
		return true
	})
	return false
}

type bindToRoot bindToParent

func (sc bindToRoot) Run(c *Char, _ []int32) bool {
	r := c.root()
	if r == nil {
		return false
	}
	c.bindTime, c.bindPos = 1, [2]float32{}
	c.setBindToId(r)
	StateControllerBase(sc).run(c, func(id byte, exp []BytecodeExp) bool {
		switch id {
		case bindToParent_time:
			c.setBindTime(exp[0].evalI(c))
		case bindToParent_facing:
			if f := exp[0].evalI(c); f < 0 {
				c.bindFacing = -1
			} else if f > 0 {
				c.bindFacing = 1
			}
		case bindToParent_pos:
			c.bindPos[0] = exp[0].evalF(c)
			if len(exp) > 1 {
				c.bindPos[1] = exp[1].evalF(c)
			}
		}
		return true
	})
	return false
}

type removeExplod StateControllerBase

const (
	removeExplod_id byte = iota
)

func (sc removeExplod) Run(c *Char, _ []int32) bool {
	StateControllerBase(sc).run(c, func(id byte, exp []BytecodeExp) bool {
		switch id {
		case removeExplod_id:
			c.removeExplod(exp[0].evalI(c))
		}
		return true
	})
	return false
}

type explodBindTime StateControllerBase

const (
	explodBindTime_id byte = iota
	explodBindTime_time
)

func (sc explodBindTime) Run(c *Char, _ []int32) bool {
	var eid, time int32 = -1, 0
	StateControllerBase(sc).run(c, func(id byte, exp []BytecodeExp) bool {
		switch id {
		case explodBindTime_id:
			eid = exp[0].evalI(c)
		case explodBindTime_time:
			time = exp[0].evalI(c)
		}
		return true
	})
	c.explodBindTime(eid, time)
	return false
}

type moveHitReset StateControllerBase

const (
	moveHitReset_ byte = iota
)

func (sc moveHitReset) Run(c *Char, _ []int32) bool {
	StateControllerBase(sc).run(c, func(id byte, exp []BytecodeExp) bool {
		switch id {
		case moveHitReset_:
			c.clearMoveHit()
		}
		return true
	})
	return false
}

type hitAdd StateControllerBase

const (
	hitAdd_value byte = iota
)

func (sc hitAdd) Run(c *Char, _ []int32) bool {
	StateControllerBase(sc).run(c, func(id byte, exp []BytecodeExp) bool {
		switch id {
		case hitAdd_value:
			c.hitAdd(exp[0].evalI(c))
		}
		return true
	})
	return false
}

type offset StateControllerBase

const (
	offset_x byte = iota
	offset_y
)

func (sc offset) Run(c *Char, _ []int32) bool {
	StateControllerBase(sc).run(c, func(id byte, exp []BytecodeExp) bool {
		switch id {
		case offset_x:
			c.offset[0] = exp[0].evalF(c)
		case offset_y:
			c.offset[1] = exp[0].evalF(c)
		}
		return true
	})
	return false
}

type victoryQuote StateControllerBase

const (
	victoryQuote_value byte = iota
)

func (sc victoryQuote) Run(c *Char, _ []int32) bool {
	StateControllerBase(sc).run(c, func(id byte, exp []BytecodeExp) bool {
		switch id {
		case victoryQuote_value:
			c.victoryQuote(exp[0].evalI(c))
		}
		return true
	})
	return false
}

type zoom StateControllerBase

const (
	zoom_pos byte = iota
	zoom_scale
)

func (sc zoom) Run(c *Char, _ []int32) bool {
	sys.drawScale = sys.cam.Scale
	StateControllerBase(sc).run(c, func(id byte, exp []BytecodeExp) bool {
		switch id {
		case zoom_pos:
			sys.zoomPos[0] = exp[0].evalF(c)
			if len(exp) > 1 {
				sys.zoomPos[1] = exp[1].evalF(c)
			}
		case zoom_scale:
			sys.drawScale = exp[0].evalF(c)
		}
		return true
	})
	return false
}

type StateBytecode struct {
	stateType StateType
	moveType  MoveType
	physics   StateType
	playerNo  int
	stateDef  stateDef
	block     StateBlock
	ctrlsps   []int32
	numVars   int32
}

func newStateBytecode(pn int) *StateBytecode {
	sb := &StateBytecode{stateType: ST_S, moveType: MT_I, physics: ST_N,
		playerNo: pn, block: *newStateBlock()}
	return sb
}
func (sb *StateBytecode) init(c *Char) {
	if sb.stateType != ST_U {
		c.ss.stateType = sb.stateType
	}
	if sb.moveType != MT_U {
		c.ss.moveType = sb.moveType
	}
	if sb.physics != ST_U {
		c.ss.physics = sb.physics
	}
	sb.ctrlsps = make([]int32, len(sb.ctrlsps))
	sys.workingState = sb
	sb.stateDef.Run(c)
}
func (sb *StateBytecode) run(c *Char) (changeState bool) {
	sys.bcVar = sys.bcVarStack.Alloc(int(sb.numVars))
	sys.workingState = sb
	changeState = sb.block.Run(c, sb.ctrlsps)
	if len(sys.bcStack) != 0 {
		sys.errLog.Println(sys.cgi[sb.playerNo].def)
		for _, v := range sys.bcStack {
			sys.errLog.Printf("%+v\n", v)
		}
		c.panic()
	}
	sys.bcVarStack.Clear()
	return
}

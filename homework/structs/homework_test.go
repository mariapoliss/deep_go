package main

import (
	"math"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

type Option func(*GamePerson)

func WithName(name string) func(*GamePerson) {
	return func(person *GamePerson) {
		person.name = name
	}
}

func WithCoordinates(x, y, z int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.x = int32(x)
		person.y = int32(y)
		person.z = int32(z)
	}
}

func WithGold(gold int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.gold = int32(gold)
	}
}

func WithMana(mana int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.mana_health &= math.MaxUint32 - math.MaxUint16
		person.mana_health |= uint32(mana << 16)
	}
}

func WithHealth(health int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.mana_health &= math.MaxUint16 << 16
		person.mana_health |= uint32(health)
	}
}

func WithRespect(respect int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.respect_strength &= 0b00001111
		person.respect_strength |= byte(respect) << 4
	}
}

func WithStrength(strength int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.respect_strength &= 0b11110000
		person.respect_strength |= byte(strength)
	}
}

func WithExperience(experience int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.experience_level &= 0b00001111
		person.experience_level |= byte(experience) << 4
	}
}

func WithLevel(level int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.experience_level &= 0b11110000
		person.experience_level |= byte(level)
	}
}

func WithHouse() func(*GamePerson) {
	return func(person *GamePerson) {
		person.house_family_gun_type &= 0b11101111
		person.house_family_gun_type |= 1 << 4
	}
}

func WithGun() func(*GamePerson) {
	return func(person *GamePerson) {
		person.house_family_gun_type &= 0b11111011
		person.house_family_gun_type |= 1 << 2
	}
}

func WithFamily() func(*GamePerson) {
	return func(person *GamePerson) {
		person.house_family_gun_type &= 0b11110111
		person.house_family_gun_type |= 1 << 3
	}
}

func WithType(personType int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.house_family_gun_type &= 0b11111100
		person.house_family_gun_type |= byte(personType)
	}
}

const (
	BuilderGamePersonType = iota
	BlacksmithGamePersonType
	WarriorGamePersonType
)

type GamePerson struct {
	x, y, z               int32
	gold                  int32
	mana_health           uint32
	name                  string
	respect_strength      byte
	experience_level      byte
	house_family_gun_type byte // 000_1_1_1_11 (000_house_family_gun_type)
}

func NewGamePerson(options ...Option) GamePerson {
	person := GamePerson{}
	for _, o := range options {
		o(&person)
	}
	return person
}

func (p *GamePerson) Name() string {
	return p.name
}

func (p *GamePerson) X() int {
	return int(p.x)
}

func (p *GamePerson) Y() int {
	return int(p.y)
}

func (p *GamePerson) Z() int {
	return int(p.z)
}

func (p *GamePerson) Gold() int {
	return int(p.gold)
}

func (p *GamePerson) Mana() int {
	return int(p.mana_health >> 16)
}

func (p *GamePerson) Health() int {
	return int(p.mana_health & 65535)
}

func (p *GamePerson) Respect() int {
	return int(p.respect_strength >> 4)
}

func (p *GamePerson) Strength() int {
	buf := p.respect_strength & 0b00001111
	return int(buf)
}

func (p *GamePerson) Experience() int {
	return int(p.experience_level >> 4)
}

func (p *GamePerson) Level() int {
	return int(p.respect_strength & 0b00001111)
}

func (p *GamePerson) HasHouse() bool {
	return (p.house_family_gun_type & 0b00010000 >> 4) == 1
}

func (p *GamePerson) HasGun() bool {
	return (p.house_family_gun_type & 0b00000100 >> 2) == 1
}

func (p *GamePerson) HasFamilty() bool {
	return (p.house_family_gun_type & 0b00001000 >> 3) == 1
}

func (p *GamePerson) Type() int {
	return int(p.house_family_gun_type & 0b00000011)
}

func TestGamePerson(t *testing.T) {
	assert.LessOrEqual(t, unsafe.Sizeof(GamePerson{}), uintptr(64))

	const x, y, z = math.MinInt32, math.MaxInt32, 0
	const name = "aaaaaaaaaaaaa_bbbbbbbbbbbbb_cccccccccccccc"
	const personType = BlacksmithGamePersonType
	const gold = math.MaxInt32
	const mana = 1000
	const health = 1000
	const respect = 10
	const strength = 10
	const experience = 10
	const level = 10

	options := []Option{
		WithName(name),
		WithCoordinates(x, y, z),
		WithGold(gold),
		WithMana(mana),
		WithHealth(health),
		WithRespect(respect),
		WithStrength(strength),
		WithExperience(experience),
		WithLevel(level),
		WithHouse(),
		WithFamily(),
		WithType(personType),
	}

	person := NewGamePerson(options...)
	assert.Equal(t, name, person.Name())
	assert.Equal(t, x, person.X())
	assert.Equal(t, y, person.Y())
	assert.Equal(t, z, person.Z())
	assert.Equal(t, gold, person.Gold())
	assert.Equal(t, mana, person.Mana())
	assert.Equal(t, health, person.Health())
	assert.Equal(t, respect, person.Respect())
	assert.Equal(t, strength, person.Strength())
	assert.Equal(t, experience, person.Experience())
	assert.Equal(t, level, person.Level())
	assert.True(t, person.HasHouse())
	assert.True(t, person.HasFamilty())
	assert.False(t, person.HasGun())
	assert.Equal(t, personType, person.Type())
}

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
		person.mediumNumbers.Set(FirsNameByte, int(name[0]))
		copy(person.name[:], name[1:])
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
		person.gold = uint32(gold)
	}
}

func WithMana(mana int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.mediumNumbers.Set(Mana, mana)
	}
}

func WithHealth(health int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.mediumNumbers.Set(Health, health)
	}
}

func WithRespect(respect int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.smallNumbers.Set(Respect, respect)
	}
}

func WithStrength(strength int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.smallNumbers.Set(Strength, strength)
	}
}

func WithExperience(experience int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.smallNumbers.Set(Experience, experience)
	}
}

func WithLevel(level int) func(*GamePerson) {
	return func(person *GamePerson) {
		person.smallNumbers.Set(Level, level)
	}
}

func WithHouse() func(*GamePerson) {
	return func(person *GamePerson) {
		person.booleans.Set(HasHouse, true)
	}
}

func WithGun() func(*GamePerson) {
	return func(person *GamePerson) {
		person.booleans.Set(HasGun, true)
	}
}

func WithFamily() func(*GamePerson) {
	return func(person *GamePerson) {
		person.booleans.Set(HasFamily, true)
	}
}

func WithType(personType int) func(*GamePerson) {
	return func(person *GamePerson) {
		// reset all prev types
		person.booleans.Set(IsBlacksmith, false)
		person.booleans.Set(IsWarrior, false)
		person.booleans.Set(IsBuilder, false)

		switch personType {
		case BuilderGamePersonType:
			person.booleans.Set(IsBuilder, true)
		case BlacksmithGamePersonType:
			person.booleans.Set(IsBlacksmith, true)
		case WarriorGamePersonType:
			person.booleans.Set(IsWarrior, true)
		}
	}
}

const (
	BuilderGamePersonType = iota
	BlacksmithGamePersonType
	WarriorGamePersonType
)

type GamePerson struct {
	x    int32
	y    int32
	z    int32
	gold uint32

	mediumNumbers MediumNumbers
	smallNumbers  SmallNumbers
	booleans      Booleans
	name          [41]byte
}

func NewGamePerson(options ...Option) GamePerson {
	gp := GamePerson{}
	for i := range options {
		options[i](&gp)
	}
	return gp
}

func (p *GamePerson) Name() string {
	fb := byte(p.mediumNumbers.Get(FirsNameByte))
	return string(append([]byte{fb}, p.name[:]...))
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
	return p.mediumNumbers.Get(Mana)
}

func (p *GamePerson) Health() int {
	return p.mediumNumbers.Get(Health)
}

func (p *GamePerson) Respect() int {
	return p.smallNumbers.Get(Respect)
}

func (p *GamePerson) Strength() int {
	return p.smallNumbers.Get(Strength)
}

func (p *GamePerson) Experience() int {
	return p.smallNumbers.Get(Experience)
}

func (p *GamePerson) Level() int {
	return p.smallNumbers.Get(Level)
}

func (p *GamePerson) HasHouse() bool {
	return p.booleans.Get(HasHouse)
}

func (p *GamePerson) HasGun() bool {
	return p.booleans.Get(HasGun)
}

func (p *GamePerson) HasFamily() bool {
	return p.booleans.Get(HasFamily)
}

func (p *GamePerson) Type() int {
	switch {
	case p.booleans.Get(IsBuilder):
		return BuilderGamePersonType
	case p.booleans.Get(IsBlacksmith):
		return BlacksmithGamePersonType
	case p.booleans.Get(IsWarrior):
		return WarriorGamePersonType
	default:
		return BuilderGamePersonType
	}
}

func TestGamePerson(t *testing.T) {
	assert.LessOrEqual(t, unsafe.Sizeof(GamePerson{}), uintptr(64))

	const x, y, z = math.MinInt32, math.MaxInt32, 0
	const name = "aaaaaaaaaaaaa_bbbbbbbbbbbbb_cccccccccccccc"
	const personType = BuilderGamePersonType
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
	assert.True(t, person.HasFamily())
	assert.False(t, person.HasGun())
	assert.Equal(t, personType, person.Type())
}

package main

// SmallNumbers contains values from 0 to 10
type SmallNumbers uint16

func (p *SmallNumbers) Set(index SmallNumber, value int) {
	if value < 0 || value > 10 {
		panic("value out of range expected 0-10")
	}
	if !index.IsValid() {
		panic("invalid index")
	}
	shift := uint(index * 4)
	mask := SmallNumbers(0xF << shift)
	*p = (*p &^ mask) | SmallNumbers(value<<shift)
}

func (p *SmallNumbers) Get(index SmallNumber) int {
	if !index.IsValid() {
		panic("invalid index")
	}
	shift := SmallNumbers(index * 4)
	return int(*p>>shift) & 0xF
}

// MediumNumbers contains values from 0 to 1000
type MediumNumbers uint32

func (p *MediumNumbers) Set(index MediumNumber, value int) {
	if value < 0 || value > 1000 {
		panic("value out of range expected 0-1000")
	}
	if !index.IsValid() {
		panic("invalid index")
	}
	shift := index * 10
	mask := MediumNumbers(0x3FF << shift)
	*p = (*p &^ mask) | (MediumNumbers(value) << shift)
}

func (p *MediumNumbers) Get(index MediumNumber) int {
	if !index.IsValid() {
		panic("invalid index")
	}
	shift := MediumNumbers(index * 10)
	return int(*p>>shift) & 0x3FF
}

// Booleans contains values 0 to 1 values
type Booleans byte

func (p *Booleans) Set(index Boolean, value bool) {
	if !index.IsValid() {
		panic("invalid index")
	}
	if value {
		*p |= Booleans(index)
	} else {
		*p &^= Booleans(index)
	}
}

func (p *Booleans) Get(index Boolean) bool {
	if !index.IsValid() {
		panic("invalid index")
	}
	return *p&Booleans(index) != 0
}

// List of enumerations for attributes for each type

type SmallNumber byte

const (
	Level SmallNumber = iota
	Experience
	Respect
	Strength
)

func (a SmallNumber) IsValid() bool {
	switch a {
	case Level, Experience, Respect, Strength:
		return true
	default:
		return false
	}
}

type Boolean byte

const (
	HasHouse Boolean = 1 << iota
	HasGun
	HasFamily
	IsWarrior
	IsBlacksmith
	IsBuilder
)

func (a Boolean) IsValid() bool {
	switch a {
	case HasHouse, HasGun, HasFamily, IsWarrior, IsBlacksmith, IsBuilder:
		return true
	default:
		return false
	}
}

type MediumNumber byte

const (
	Health MediumNumber = iota
	Mana
	FirsNameByte
)

func (a MediumNumber) IsValid() bool {
	switch a {
	case Health, Mana, FirsNameByte:
		return true
	default:
		return false
	}
}

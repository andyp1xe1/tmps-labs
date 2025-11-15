package character

// Race constants - Available character races
const (
	RaceHuman = "Human"
	RaceElf   = "Elf"
	RaceDwarf = "Dwarf"
)

// Class constants - Available character classes
const (
	ClassFighter = "Fighter"
	ClassWizard  = "Wizard"
	ClassRogue   = "Rogue"
)

// Background constants - Available character backgrounds
const (
	BackgroundAcolyte  = "Acolyte"
	BackgroundCriminal = "Criminal"
	BackgroundScholar  = "Scholar"
)

// Creation method constants - Character creation methods
const (
	CreationMethodStandard = "standard"
	CreationMethodRandom   = "random"
	CreationMethodCustom   = "custom"
)

// Stat name constants - The six core D&D ability scores
const (
	StatStrength     = "Strength"
	StatDexterity    = "Dexterity"
	StatConstitution = "Constitution"
	StatIntelligence = "Intelligence"
	StatWisdom       = "Wisdom"
	StatCharisma     = "Charisma"
)

// Form key constants - Form field keys for UI
const (
	FormKeyStrength     = "strength"
	FormKeyDexterity    = "dexterity"
	FormKeyConstitution = "constitution"
	FormKeyIntelligence = "intelligence"
	FormKeyWisdom       = "wisdom"
	FormKeyCharisma     = "charisma"
	FormKeyMethod       = "method"
)

// Racial ability constants - Common racial abilities
const (
	AbilityDarkvision            = "Darkvision"
	AbilityKeenSenses            = "Keen Senses"
	AbilityFeyAncestry           = "Fey Ancestry"
	AbilityTrance                = "Trance"
	AbilityDwarvenResilience     = "Dwarven Resilience"
	AbilityDwarvenCombatTraining = "Dwarven Combat Training"
	AbilityStonecunning          = "Stonecunning"
	AbilityExtraLanguage         = "Extra Language"
	AbilityExtraSkill            = "Extra Skill"
	AbilityHumanDetermination    = "Human Determination"
)

// Class ability constants - Core class abilities
const (
	AbilityFightingStyle     = "Fighting Style"
	AbilitySecondWind        = "Second Wind"
	AbilityActionSurge       = "Action Surge"
	AbilitySpellcasting      = "Spellcasting"
	AbilityArcaneRecovery    = "Arcane Recovery"
	AbilityArcaneTradition   = "Arcane Tradition"
	AbilityExpertise         = "Expertise"
	AbilitySneakAttack       = "Sneak Attack"
	AbilityCunningAction     = "Cunning Action"
	AbilityExtraAttack       = "Extra Attack"
	AbilityImprovedAbilities = "Improved Abilities"
)

// Equipment constants - Common equipment items
const (
	EquipmentChainMail    = "Chain Mail"
	EquipmentShield       = "Shield"
	EquipmentSword        = "Sword"
	EquipmentSpellbook    = "Spellbook"
	EquipmentQuarterstaff = "Quarterstaff"
	EquipmentDagger       = "Dagger"
	EquipmentLeatherArmor = "Leather Armor"
	EquipmentShortsword   = "Shortsword"
	EquipmentThievesTools = "Thieves' Tools"
	EquipmentHolySymbol   = "Holy Symbol"
	EquipmentPrayerBook   = "Prayer Book"
	EquipmentCrowbar      = "Crowbar"
	EquipmentDarkClothes  = "Dark Clothes"
	EquipmentBooks        = "Books"
	EquipmentInkAndQuill  = "Ink and Quill"
	EquipmentMagicWeapon  = "+1 Weapon"
	EquipmentMagicArmor   = "Magic Armor"
)

// Personality trait constants - Common personality traits
const (
	PersonalityDevout      = "devout"
	PersonalityHelpful     = "helpful"
	PersonalitySecretive   = "secretive"
	PersonalityStreetSmart = "street-smart"
	PersonalityCurious     = "curious"
	PersonalityMethodical  = "methodical"
)

// Skill constants - D&D skills
const (
	SkillInsight       = "Insight"
	SkillReligion      = "Religion"
	SkillDeception     = "Deception"
	SkillStealth       = "Stealth"
	SkillInvestigation = "Investigation"
	SkillHistory       = "History"
)

package recipes

import (
	"encoding/json"
	"fmt"
	"time"
)

type Config struct {
	Storage Storer
}

type ScheduleConfig struct {
	TickerDuration time.Duration
	DailyRecipe    *Recipe // represents a cached version of the current day's recipe
	NextRecipe     *Recipe // represents a cached version of the next day's recipe to be updated/displayed by the scheduler
}

// ScheduleDailyRecipe is a job scheduler that updates the daily recipe based on a configured condition that is checked
// periodically by a configured duration.
func (sc *ScheduleConfig) ScheduleDailyRecipe() {
	done := make(chan bool)
	ticker := time.NewTicker(sc.TickerDuration) // TODO: pull this from config, update param list
	defer ticker.Stop()

	var t time.Time
	for {
		select {
		case <-done:
			panic("uhoh") // TODO: handle it! this is an error case only
		case t = <-ticker.C:
			if t.Hour() == 0 { // TODO: update this to reflect config condition. currently checks at midnight
				fmt.Println("<< okay...updating recipe")
				sc.DailyRecipe = sc.NextRecipe
			}
		}
	}
}

// UnitLength are units of length, e.g. inches, centimeters, etc.
type UnitLength int

const (
	Centimeters UnitLength = iota
	Inches
)

// UnitMass are units of mass/volume, e.g. grams, cups, etc.
type UnitMass int

const (
	Whole UnitMass = iota
	Grams
	Milliliters
)

var UnitMassToString = map[UnitMass]string{
	Whole:       "whole",
	Grams:       "g",
	Milliliters: "ml",
}

func (um UnitMass) String() string {
	if s, ok := UnitMassToString[um]; ok {
		return s
	}
	// TODO: handle it!
	panic("uhoh")
}

type SourceCategory int

const (
	CountryFlag SourceCategory = iota
	CookBook
	ChefSite
)

// Recipe describes a recipe to be displayed.
type Recipe struct {
	Attribution Source `db:"attribution" json:"attribution,omitempty"` // where the recipe comes from
	//IngredientSpecifications []IngredientSpecification `db:"ingredient_specification" json:"ingredient_specifications"`
	Components   []Component   `db:"component" json:"components,omitempty"`
	Cuisine      Source        `db:"cuisine" json:"cuisine,omitempty"` // the country cuisine this recipe is, e.g. "Indian"
	Emojis       []Source      `db:"emoji" json:"emojis,omitempty"`    // emoji images representing quick information about the recipe, e.g. a stopwatch emoji to indicate quick prep or cook time
	Instructions []Instruction `db:"instruction" json:"instructions"`
	EnglishName  string        `db:"english_name" json:"english_name"`
	NativeName   string        `db:"native_name" json:"native_name,omitempty"` // the recipe name in its native language
	Notes        []string      `db:"note" json:"notes,omitempty"`              // any notes about the recipe overall
}

func (r *Recipe) Scan(src any) error {
	return json.Unmarshal(src.([]byte), r)
}

// Component refers to a name to group the ingredients by, e.g. "dough", "filling"
type Component struct {
	Name                     string                    `db:"name" json:"name"`
	IngredientSpecifications []IngredientSpecification `db:"ingredient_specification" json:"ingredient_specifications"`
}

// Source is any link (path to image, URL, etc) and its description.
type Source struct {
	Description string         `db:"description" json:"description"`
	Location    string         `db:"location" json:"location"` // a path to an image, url, page number, etc.
	Category    SourceCategory `db:"category" json:"category"`
}

// IngredientSpecification is an ingredient, how much of the ingredient to be used, plus how it's prepared.
type IngredientSpecification struct {
	Ingredient          `db:"ingredient" json:"ingredient"`
	Note                string     `db:"note" json:"note,omitempty"`                                 // describes more about the ingredient, such as "for garnish"
	AmountQuantity      []float32  `db:"amount_quantity" json:"amount_quantity,omitempty"`           // for the situation of "#-# of things". should always be by weight or amount, but never volume. either this is set, or PreparationQuantity is set, but not both
	AmountMass          UnitMass   `db:"amount_mass" json:"amount_mass,omitempty"`                   // e.g. grams, whole
	PreparationQuantity float32    `db:"preparation_quantity" json:"preparation_quantity,omitempty"` // e.g. 1, 0.5, etc. either this is set, or AmountQuantity is set, but not both
	PreparationLength   UnitLength `db:"preparation_length" json:"preparation_length,omitempty"`     // e.g. centimeter, inch
	PreparationType     string     `db:"preparation_type" json:"preparation_type,omitempty"`         // e.g. cubes, strips, etc.
}

// Ingredient is a specific ingredient describing what it's called, plus where it can be purchased.
type Ingredient struct {
	Name           string `db:"english_name" json:"english_name"`                   // english name of the ingredient, e.g. Korean Pancake Mix, Diced Tomato
	Category       string `db:"english_category" json:"english_category,omitempty"` // english category of the ingredient, e.g. tomato. for quickly finding ingredients
	NativeName     string `db:"native_name" json:"native_name,omitempty"`           // the ingredient name in its native language, e.g. 부침가루
	TranslatedName string `db:"translated_name" json:"translated_name,omitempty"`   // the ingredient name translated to its english name, e.g. Buchim Garu
	ShoppingLink   string `db:"shopping_link" json:"shopping_link,omitempty"`       // a direct link to purchase the ingredient from
}

/*
Instruction groups the steps into a shared name, e.g. "Prep the Chicken".
In the event of a single grouping, part should be left empty.
*/
type Instruction struct {
	Part  string `db:"part" json:"part"`
	Steps []Step `db:"step" json:"steps"`
}

/*
Step contains the information necessary to execute an instruction.

IsParallel refers to whether a step differing from previous steps needs to be done simultaneously:

	> e.g. "In a separate saucepan, cook the quinoa in water"

Action contains information on how to execute the step:

	> e.g. "Heat the sesame oil in a large skillet over medium heat"

Notes are any additional/optional information to mention/clarify during the step:

	> e.g. "For high altitude, set oven to 375F"
*/
type Step struct {
	IsParallel bool          `db:"is_parallel" json:"is_parallel,omitempty"` // affects front-end formatting
	Action     Information   `db:"action" json:"action"`
	Notes      []Information `db:"note" json:"note,omitempty"`
}

// Information contains step-related information, which might include a temperature.
type Information struct {
	Message            string `db:"message" json:"message"`
	TemperatureCelsius *int   `db:"unit_temperature" json:"unit_temperature,omitempty"`
}

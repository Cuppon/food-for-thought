package recipes

import (
	"github.com/jmoiron/sqlx"
)

type Storer interface {
	//AddRecipe()
	GetRecipe(ID int) (Recipe, error)
	GetTotalRecipes() (int, error)
	SetNextRecipe(ID int) error
	UpdateRecipeUsageCount(ID int) (bool, error)
}

type PG struct {
	Conn *sqlx.DB
}

func (db *PG) GetRecipe(ID int) (Recipe, error) {
	row := db.Conn.QueryRow(dailyRecipe, ID)

	var r Recipe
	err := row.Scan(&r)
	if err != nil {
		return Recipe{}, err
	}

	return r, nil
}

func (db *PG) GetTotalRecipes() (int, error) {
	row := db.Conn.QueryRow(`SELECT COUNT(*) FROM recipe;`)

	var count int
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// TODO: create an endpoint for this
func (db *PG) SetNextRecipe(ID int) error {
	return nil
}

func (db *PG) UpdateRecipeUsageCount(ID int) (bool, error) {
	_, err := db.Conn.Exec(`UPDATE recipe SET usage_count = usage_count + 1 WHERE ID = $1;`, ID)
	if err != nil {
		return false, err
	}

	return true, nil
}

var dailyRecipe = `SELECT jsonb_build_object(
               'attribution', (SELECT jsonb_build_object(
                                  'id', r.attribution_source_id,
                                  'description', s.description,
                                  'location', s.location,
                                  'category', s.category
                               )
                               FROM source AS s
                               WHERE s.id = r.attribution_source_id
               ),
               'components', (SELECT jsonb_agg(jsonb_build_object(
                                'name', component,
                                'ingredient_specifications', ingredient_specifications
                                )) AS components
                               FROM (
                                   SELECT
                                       component,
                                       jsonb_agg(jsonb_build_object(
                                               'note', note,
                                               'component', component,
                                               'ingredient', jsonb_build_object(
                                                       'id', ingredient.id,
                                                       'native_name', ingredient.native_name,
                                                       'english_name', ingredient.english_name,
                                                       'shopping_link', ingredient.shopping_link,
                                                       'translated_name', ingredient.translated_name,
                                                       'english_category', ingredient.english_category
                                                             ),
                                               'amount_quantity', ins.amount_quantity,
                                               'amount_mass', ins.amount_mass,
                                               'preparation_quantity', ins.preparation_quantity,
                                               'preparation_type', ins.preparation_type,
                                               'preparation_length', ins.preparation_length
                                                 )) AS ingredient_specifications
                                   FROM ingredient_specification as ins
                                            JOIN ingredient ON ins.ingredient_id = ingredient.id
                                   GROUP BY component
                               ) AS comps
               ),
               'cuisine', (SELECT jsonb_build_object(
                              'id', r.cuisine_source_id,
                              'description', s.description,
                              'location', s.location,
                              'category', s.category
                           )
                           FROM source AS s
                           WHERE s.id = r.cuisine_source_id
               ),
               'emojis', (SELECT jsonb_agg(jsonb_build_object(
                            'id', rs.id,
                            'description', s.description,
                            'location', s.location,
                            'category', s.category
                          ))
                          FROM source AS s
                          INNER JOIN recipe_source AS rs ON s.id = rs.emoji_source_id
                          GROUP BY rs.recipe_id
               ),
               'instructions', r.instruction,
               'english_name', r.english_name,
               'native_name', r.native_name,
               'notes', r.note
       ) AS recipe
FROM recipe r
INNER JOIN ingredient_specification AS isp ON isp.recipe_id = r.id
WHERE r.id = $1
GROUP BY r.id;`
